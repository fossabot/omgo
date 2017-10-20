package main

import (
	"encoding/binary"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/api"
	"github.com/master-g/omgo/kit/services"
	"github.com/master-g/omgo/kit/utils"
	"gopkg.in/urfave/cli.v2"
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
	defaultListen        = ":8888"
	defaultKind          = "agent"
	defaultName          = "agent-0"
	defaultETCD          = "http://127.0.0.1:2379"
	defaultRoot          = "backends"
	defaultGameServerId  = "game-0"
	defaultReadDeadLine  = 3 * time.Minute
	defaultSockBufSize   = 32*1024 - 1
	defaultTxQueueLength = 128
	defaultRPMLimit      = 200
)

var (
	defaultServices = []string{"snowflake", "dbservice", "gameservice"}
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
				Value:   defaultListen,
			},
			&cli.StringFlag{
				Aliases: []string{"k"},
				Name:    "kind",
				Usage:   "agent service kind",
				Value:   defaultKind,
			},
			&cli.StringFlag{
				Aliases: []string{"n"},
				Name:    "name",
				Usage:   "agent service name",
				Value:   defaultName,
			},
			&cli.StringSliceFlag{
				Aliases: []string{"e"},
				Name:    "etcdhosts",
				Usage:   "ETCD endpoint addresses",
				Value:   cli.NewStringSlice(defaultETCD),
			},
			&cli.StringFlag{
				Aliases: []string{"r"},
				Name:    "etcdroot",
				Usage:   "services root path on ETCD",
				Value:   defaultRoot,
			},
			&cli.StringSliceFlag{
				Aliases: []string{"s"},
				Name:    "services",
				Usage:   "service with kinds to connect to",
				Value:   cli.NewStringSlice(defaultServices...),
			},
			&cli.StringFlag{
				Aliases: []string{"g"},
				Name:    "game",
				Usage:   "game server name",
				Value:   defaultGameServerId,
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
			listenOn := c.String("listen")
			agentKind := c.String("kind")
			agentName := c.String("name")

			log.Infof("listen:%v", listenOn)
			log.Infof("kind:%v", agentKind)
			log.Infof("name:%v", agentName)
			log.Infof("etcdhosts:%v", etcdHosts)
			log.Infof("etcdroot:%v", etcdRoot)
			log.Infof("services:%v", serviceNames)
			log.Infof("deadline:%v", c.Duration("deadline"))
			log.Infof("txqueuelen:%v", c.Int("txqueuelen"))
			log.Infof("sockbufsize:%v", c.Int("sockbufsize"))
			log.Infof("rpm:%v", rpmLimit)

			// create configuration
			config := &Config{
				listen:        c.String("listen"),
				readDeadline:  c.Duration("deadline"),
				txQueueLength: c.Int("txqueuelen"),
				sockBufSize:   c.Int("sockbufsize"),
			}

			// capture UNIX SYSTERM signal
			go sigHandler()
			// register agent to ETCD
			agentFullPath := services.GenPath(etcdRoot, agentKind, agentName)
			services.RegisterService(etcdHosts, agentFullPath, listenOn)
			// connect to other services
			// TODO, create service pools

			// start timer worker
			initTimer(rpmLimit)

			// listen to client connections
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
	var session api.Session
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Error("get remote address failed: ", err)
		return
	}
	session.IP = net.ParseIP(host)
	session.Port = port
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
