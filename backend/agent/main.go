package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/utils"
	"gopkg.in/urfave/cli.v2"
	"net/http"
	"os"
	"time"
)

// Config holds configuration for agent
type Config struct {
	listen         string        // address to listen
	readDeadline   time.Duration // read timeout
	sockBufferSize int           // socket buffer size
	udpBufferSize  int           // UDP buffer size
	txQueueLength  int           // transmission queue length
	dscp           int           // Differentiated Services Code Point https://www.tucny.com/Home/dscp-tos
	sendWindowSize int           // send window size
	recvWindowSize int           // receive window size
	mtu            int           // MTU maximum transmission unit
	nodelay        int           // TCP no delay flag
}

var (
	profileAddr = "0.0.0.0:6666"
)

func main() {
	log.SetLevel(log.DebugLevel)
	defer utils.PrintPanicStack()

	// profiling
	go http.ListenAndServe(profileAddr, nil)
	app := &cli.App{
		Name:    "agent",
		Usage:   "a gateway service for game server",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "listen",
				Value: ":8888",
				Usage: "listening address:port",
			},
			&cli.StringSliceFlag{
				Name:  "etcdhosts",
				Value: cli.NewStringSlice("http://127.0.0.1:2379"),
				Usage: "etcd hosts",
			},
		},
		Action: func(c *cli.Context) error {
			log.Println("listen:", c.String("listen"))
			log.Println("etcdhosts:", c.StringSlice("etcdhosts"))
			return nil
		},
	}
	app.Run(os.Args)
}
