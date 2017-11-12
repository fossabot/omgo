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
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	_ "net/http/pprof"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/backend/agent/api"
	"github.com/master-g/omgo/kit/services"
	"github.com/master-g/omgo/kit/utils"
	pc "github.com/master-g/omgo/proto/pb/common"
	"gopkg.in/urfave/cli.v2"
)

// Config holds configuration for agent
type Config struct {
	port          int           // port to listen
	readDeadline  time.Duration // read timeout
	sockBufSize   int           // socket buffer size
	txQueueLength int           // transmission queue length
}

const (
	profileAddress       = "localhost:6060"
	defaultPort          = 8888
	defaultKind          = "agent"
	defaultName          = "agent-0"
	defaultHost          = "localhost"
	defaultETCD          = "http://127.0.0.1:2379"
	defaultRoot          = "backends"
	defaultGameServerID  = "game-0"
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
			&cli.IntFlag{
				Name:  "listen",
				Usage: "listening port",
				Value: defaultPort,
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
			&cli.StringFlag{
				Name:  "service-host",
				Usage: "agent service host",
				Value: defaultHost,
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
				Value: defaultGameServerID,
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
			listenPort := c.String("listen")
			serviceRoot := c.String("service-root")
			agentKind := c.String("service-kind")
			agentName := c.String("service-name")
			agentHost := c.String("service-host")
			gameServerName := c.String("gameserver-name")

			log.Info("--------------------------------------------------")
			log.Infof("listen on:%v", listenPort)
			log.Infof("service-root:%v", serviceRoot)
			log.Infof("service-kind:%v", agentKind)
			log.Infof("service-name:%v", agentName)
			log.Infof("service-host:%v", agentHost)
			log.Infof("etcd-hosts:%v", etcdHosts)
			log.Infof("services:%v", serviceNames)
			log.Infof("deadline:%v", c.Duration("deadline"))
			log.Infof("txqueuelen:%v", c.Int("txqueuelen"))
			log.Infof("sockbufsize:%v", c.Int("sockbufsize"))
			log.Infof("rpm:%v", rpmLimit)
			log.Info("--------------------------------------------------")

			// create configuration
			config := &Config{
				port:          c.Int("listen"),
				readDeadline:  c.Duration("deadline"),
				txQueueLength: c.Int("txqueuelen"),
				sockBufSize:   c.Int("sockbufsize"),
			}

			// register agent to ETCD
			agentFullPath := services.GenPath(serviceRoot, agentKind, agentName)
			agentFullHost := fmt.Sprintf("%v:%v", agentHost, listenPort)
			services.RegisterService(etcdHosts, agentFullPath, agentFullHost)
			// connect to other services
			srvConfig := api.Config{
				Root:            serviceRoot,
				Hosts:           etcdHosts,
				GameServerKind:  "game",
				GameServerName:  gameServerName,
				DataServiceKind: "dataservice",
			}
			api.Init(srvConfig)

			// server exit callback
			exitCallback := func() {
				services.UnregisterService(etcdHosts, agentFullPath)
			}

			// capture UNIX SIGTERM signal
			go sigHandler(exitCallback)

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
	tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf(":%v", config.port))
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
	// header size
	headerSize := make([]byte, 2)
	// agent's input channel
	in := make(chan *api.IncomingPacket)
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

		// read header size
		n, err := io.ReadFull(conn, headerSize)
		if err != nil {
			log.Warningf("%v read header size failed: %v %v bytes read", session.IP, err, n)
			return
		}
		size := binary.BigEndian.Uint16(headerSize)

		// header message
		headerData := make([]byte, size)
		n, err = io.ReadFull(conn, headerData)
		if err != nil {
			log.Warningf("%v read header message failed: %v expect: %v actual read: %v", session.IP, err, size, n)
			return
		}

		headerMsg := &pc.Header{}
		err = proto.Unmarshal(headerData, headerMsg)
		if err != nil {
			log.Warningf("%v invalid header: %v", session.IP, err)
			return
		}

		payloadSize := headerMsg.BodySize
		if payloadSize == 0 || payloadSize > defaultSockBufSize {
			log.Warningf("%v payload size error %v", session.IP, payloadSize)
		}

		// payload
		payload := make([]byte, payloadSize)
		n, err = io.ReadFull(conn, payload)
		if err != nil {
			log.Warningf("%v read payload failed: %v expect: %v actual read: %v", session.IP, err, size, n)
			return
		}

		inPacket := &api.IncomingPacket{
			Header: headerMsg,
			Body:   payload,
		}

		// deliver the payload to the input queue of agent
		select {
		case in <- inPacket:
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
