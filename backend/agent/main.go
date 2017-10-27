/*
Package agent is designed for:
1. Manage client connections and sessions.
2. Passthrough data to game server (via gRPC streaming).
3. Allow backends to reboot/deploy without close the connection.
4. Isolate core services from exposure.

The basic work flow of agent is:
1. main.go      extract arguments from command line via urfave/cli.v2 package
2. signal.go    start a goroutine to capture UNIX SIGTERM signal
3. api.go       connect to data and game gRPC services via ETCD
4. main.go      start a tcpServer goroutine for handling incoming connections
5. main.go      for each connection, spawns a handleClient goroutine

    handleClient() {
        // init client session and other context
        createSession()
        // create buffer object and start a goroutine for sending packet
        go bufferOut()
        // create agent instance for this client, process 4 types of message:
        // 1. incoming packages
        // 2. game stream frames
        // 3. timer (rpm limit, heartbeat, etc.)
        // 4. server shutdown
        go agent()
        for {
            // read from TCP and feed to agent via channel
        }
    }

*/
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
	defaultRPCPort       = 30001
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
				Name:  "listen",
				Usage: "listening address:port",
				Value: defaultListen,
			},
			&cli.StringFlag{
				Name:  "service-root",
				Usage: "services root path on ETCD",
				Value: defaultRoot,
			},
			&cli.StringFlag{
				Name:  "service-kind",
				Usage: "agent service kind",
				Value: defaultKind,
			},
			&cli.StringFlag{
				Name:  "service-name",
				Usage: "agent service name",
				Value: defaultName,
			},
			&cli.IntFlag{
				Name:  "service-port",
				Usage: "agent rpc service port",
				Value: defaultRPCPort,
			},
			&cli.StringSliceFlag{
				Name:  "etcd-host",
				Usage: "ETCD endpoint addresses",
				Value: cli.NewStringSlice(defaultETCD),
			},
			&cli.StringSliceFlag{
				Name:  "add-service",
				Usage: "service with kinds to connect to",
				Value: cli.NewStringSlice(defaultServices...),
			},
			&cli.StringFlag{
				Name:  "gameserver-name",
				Usage: "game server name",
				Value: defaultGameServerId,
			},
			&cli.DurationFlag{
				Name:  "deadline",
				Usage: "read timeout per connection",
				Value: defaultReadDeadLine,
			},
			&cli.IntFlag{
				Name:  "txqueuelen",
				Usage: "transmission queue length per connection",
				Value: defaultTxQueueLength,
			},
			&cli.IntFlag{
				Name:  "sockbufsize",
				Usage: "TCP socket buffer size per connection",
				Value: defaultSockBufSize,
			},
			&cli.IntFlag{
				Name:  "rpm",
				Usage: "packet limit per minute",
				Value: defaultRPMLimit,
			},
		},
		Action: func(c *cli.Context) error {
			etcdHosts := c.StringSlice("etcd-host")
			serviceNames := c.StringSlice("add-service")
			rpmLimit := c.Int("rpm")
			listenOn := c.String("listen")
			serviceRoot := c.String("service-root")
			agentKind := c.String("service-kind")
			agentName := c.String("service-name")
			agentRPCPort := c.Int("service-port")
			gameServerName := c.String("gameserver-name")

			log.Info("--------------------------------------------------")
			log.Infof("listen on:%v", listenOn)
			log.Infof("service-root:%v", serviceRoot)
			log.Infof("service-kind:%v", agentKind)
			log.Infof("service-name:%v", agentName)
			log.Infof("service-port:%v", agentRPCPort)
			log.Infof("etcd-hosts:%v", etcdHosts)
			log.Infof("services:%v", serviceNames)
			log.Infof("deadline:%v", c.Duration("deadline"))
			log.Infof("txqueuelen:%v", c.Int("txqueuelen"))
			log.Infof("sockbufsize:%v", c.Int("sockbufsize"))
			log.Infof("rpm:%v", rpmLimit)
			log.Info("--------------------------------------------------")

			// create configuration
			config := &Config{
				listen:        c.String("listen"),
				readDeadline:  c.Duration("deadline"),
				txQueueLength: c.Int("txqueuelen"),
				sockBufSize:   c.Int("sockbufsize"),
			}

			// capture UNIX SIGTERM signal
			go sigHandler()
			// register agent to ETCD
			agentFullPath := services.GenPath(serviceRoot, agentKind, agentName)
			services.RegisterService(etcdHosts, agentFullPath, listenOn)
			// connect to other services
			srvConfig := api.Config{
				Root:            serviceRoot,
				Hosts:           etcdHosts,
				GameServerKind:  "game",
				GameServerName:  gameServerName,
				DataServiceKind: "dataservice",
			}
			api.Init(srvConfig)

			// setup session time parameters
			api.SetReadDeadLine(defaultReadDeadLine)
			api.SetRPMLimit(rpmLimit)

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
