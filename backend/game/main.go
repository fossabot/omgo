package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/utils"
	"gopkg.in/urfave/cli.v2"
	"net"
	"net/http"
	"os"
)

const (
	profileAddress  = "0.0.0.0:6666"
	defaultETCD     = "http://127.0.0.1:2379"
	defaultRoot     = "/backends"
	defaultServices = []string{"snowflake", "agent"}
)

func main() {
	log.SetLevel(log.DebugLevel)
	defer utils.PrintPanicStack()

	// profiling
	go http.ListenAndServe(profileAddress, nil)

	// cli
	app := &cli.App{
		Name:    "game",
		Usage:   "a stream processor based game",
		Version: "2.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases: []string{"i"},
				Name:    "id",
				Value:   "game1",
				Usage:   "id of this service",
			},
			&cli.StringFlag{
				Aliases: []string{"l"},
				Name:    "listen",
				Usage:   "listening address:port",
				Value:   "10000",
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
		},
		Action: func(c *cli.Context) {
			log.Println("id:", c.String("id"))

			// listen
			lis, err := net.Listen("tcp", c.String("listen"))
			if err != nil {
				log.Panic(err)
				os.Exit(-1)
			}
			log.Info("listening on ", lis.Addr())

			// register services
		},
	}

	app.Run(os.Args)
}
