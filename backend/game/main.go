package main

import (
	"net"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/proto/grpc/game"
	"github.com/master-g/omgo/services"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
)

const (
	profileAddress = "0.0.0.0:6666"
	defaultETCD    = "http://127.0.0.1:2379"
	defaultRoot    = "backends"
	defaultListen  = ":10000"
	defaultSID     = "game-0"
)

var (
	defaultServices = []string{"snowflake", "agent", "dbs"}
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
				Usage:   "id of this service",
				Value:   defaultSID,
			},
			&cli.StringFlag{
				Aliases: []string{"l"},
				Name:    "listen",
				Usage:   "listening address:port",
				Value:   defaultListen,
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
		Action: func(c *cli.Context) error {
			cfgID := c.String("id")
			cfgListen := c.String("listen")
			cfgETCDHosts := c.StringSlice("etcdhosts")
			cfgETCDRoot := c.String("etcdroot")
			cfgServices := c.StringSlice("services")

			log.Println("id:", cfgID)
			log.Println("listen:", cfgListen)
			log.Println("etcd-hosts:", cfgETCDHosts)
			log.Println("etcd-root:", cfgETCDRoot)
			log.Println("services:", cfgServices)
			// listen
			lis, err := net.Listen("tcp", cfgListen)
			if err != nil {
				log.Panic(err)
				os.Exit(-1)
			}
			log.Info("listening on ", lis.Addr())

			// register services
			s := grpc.NewServer()
			ins := new(server)
			pb.RegisterGameServiceServer(s, ins)

			// initialize services
			services.Init(cfgETCDRoot, cfgETCDHosts, cfgServices)

			return s.Serve(lis)
		},
	}
	app.Run(os.Args)
}
