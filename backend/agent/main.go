package main

import (
	"encoding/binary"
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/utils"
	"gopkg.in/urfave/cli.v2"
	"io"
	"net"
	"net/http"
	"os"
	"time"
)

// Config holds configuration for agent
type Config struct {
	listen        string        // address to listen
	readDeadline  time.Duration // read timeout
	sockBufSize   int           // socket buffer size
	txQueueLength int           // transmission queue length
}

const (
	profileAddress       = "0.0.0.0:6666"
	defaultETCD          = "http://127.0.0.1:2379"
	defaultRoot          = "/backends"
	defaultReadDeadLine  = 15 * time.Second
	defaultSockBufSize   = 32*1024 - 1
	defaultTxQueueLength = 128
	defaultRPMLimit      = 200
)

var (
	defaultServices = []string{"snowflake", "game"}
)

func main() {
	log.SetLevel(log.DebugLevel)
	defer utils.PrintPanicStack()

	// profiling
	go http.ListenAndServe(profileAddress, nil)

	// cli
	app := &cli.App{
		Name:    "agent",
		Usage:   "a gateway service for game server",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases: []string{"l"},
				Name:    "listen",
				Usage:   "listening address:port",
				Value:   ":8888",
			},
			&cli.StringSliceFlag{
				Aliases: []string{"e"},
				Name:    "etcdhosts",
				Usage:   "etcd hosts",
				Value:   cli.NewStringSlice(defaultETCD),
			},
			&cli.StringFlag{
				Aliases: []string{"r"},
				Name:    "etcdroot",
				Usage:   "services root path on etcd",
				Value:   defaultRoot,
			},
			&cli.StringSliceFlag{
				Aliases: []string{"s"},
				Name:    "services",
				Usage:   "service names",
				Value:   cli.NewStringSlice(defaultServices...),
			},
			&cli.DurationFlag{
				Aliases: []string{"d"},
				Name:    "deadline",
				Usage:   "read timeout per connection",
				Value:   defaultReadDeadLine,
			},
			&cli.IntFlag{
				Aliases: []string{"t"},
				Name:    "txqueuelen",
				Usage:   "transmission queue length per connection",
				Value:   defaultTxQueueLength,
			},
			&cli.IntFlag{
				Aliases: []string{"o"},
				Name:    "sockbufsize",
				Usage:   "TCP socket buffer size per connection",
				Value:   defaultSockBufSize,
			},
			&cli.IntFlag{
				Aliases: []string{"p"},
				Name:    "rpm",
				Usage:   "Packet limit per minute per connection",
				Value:   defaultRPMLimit,
			},
		},
		Action: func(c *cli.Context) error {
			etcdHosts := c.StringSlice("etcdhosts")
			etcdRoot := c.String("etcdroot")
			serviceNames := c.StringSlice("services")
			rpmLimit = c.Int("rpm")
			log.Println("listen:", c.String("listen"))
			log.Println("etcdhosts:", etcdHosts)
			log.Println("etcdroot:", etcdRoot)
			log.Println("services:", serviceNames)
			log.Println("deadline:", c.Duration("deadline"))
			log.Println("txqueuelen:", c.Int("txqueuelen"))
			log.Println("sockbufsize:", c.Int("sockbufsize"))
			log.Println("rpm:", rpmLimit)

			// create configuration
			config := &Config{
				listen:        c.String("listen"),
				readDeadline:  c.Duration("deadline"),
				txQueueLength: c.Int("txqueuelen"),
				sockBufSize:   c.Int("sockbufsize"),
			}

			startup(etcdRoot, etcdHosts, serviceNames)
			initTimer(rpmLimit)

			// listeners
			go tcpServer(config)

			// wait forever
			select {}
			return nil
		},
	}
	app.Run(os.Args)
}

func tcpServer(config *Config) {
	// resolve address & start listening
	tcpAddr, err := net.ResolveTCPAddr("tcp4", config.listen)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	log.Info("listening on:", listener.Addr())

	// loop accepting
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Warning("accept failed:", err)
			continue
		}
		// set socket read buffer size
		conn.SetReadBuffer(config.sockBufSize)
		// set socket write buffer size
		conn.SetWriteBuffer(config.sockBufSize)
		// start a goroutine for every incoming connection to read
		go handleClient(conn, config)
	}
}

// PIPELINE #1: handleClient
// the goroutine is used for reading incoming packets
// each packet is defined as:
// | size (2 bytes) | payload |
func handleClient(conn net.Conn, config *Config) {
	defer utils.PrintPanicStack()
	defer conn.Close()
	// header
	header := make([]byte, 2)
	// agent's input channel
	in := make(chan []byte)
	defer func() {
		close(in) // session will be closed
	}()

	// create a new session object for this connection
	// and record its IP address
	var session types.Session
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error("get remote address failed: ", err)
		return
	}
	session.IP = net.ParseIP(host)
	log.Infof("new connection from %v:%v", host, port)

	// session die signal, will be triggered by agent()
	session.Die = make(chan struct{})

	// create a write buffer
	out := newBuffer(conn, session.Die, config.txQueueLength)
	go out.start()

	// start agent for packet processing
	wg.Add(1)
	go agent(&session, in, out)

	// read loop
	for {
		// solve dead link problem:
		// physical disconnection without any communication between client and server
		// will cause the read to block FOREVER, so a timeout will save the day.
		conn.SetReadDeadline(time.Now().Add(config.readDeadline))

		// read size header
		n, err := io.ReadFull(conn, header)
		if err != nil {
			log.Warningf("%v read header failed: %v %v bytes read", session.IP, err, n)
			return
		}
		size := binary.BigEndian.Uint16(header)

		// data
		payload := make([]byte, size)
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			log.Warningf("%v read payload failed: %v expect: %v actual read: %v", session.IP, err, size, n)
			return
		}

		// deliver the payload to the input queue of agent
		select {
		case in <- payload:
		case <-session.Die:
			log.Warningf("%v connection closed by logic, flag: %v", session.IP, session.Flag)
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}
}
