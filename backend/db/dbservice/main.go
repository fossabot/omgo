package main

import (
	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/proto/grpc/db"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
	"net"
	"os"
	"sort"
	"time"
)

const (
	defaultListen      = ":60001"
	defaultRedisHost   = "127.0.0.1:27017"
	defaultRedisDB     = 0
	defaultMaxIdle     = 80
	defaultMaxActive   = 1024
	defaultIdleTimeout = 180 * time.Second
)

func main() {
	log.SetLevel(log.DebugLevel)
	defer utils.PrintPanicStack()

	app := &cli.App{
		Name:    "dbservice",
		Usage:   "Database service",
		Version: "v1.0.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases: []string{"l"},
				Name:    "listen",
				Usage:   "listening address:port",
				Value:   defaultListen,
			},
			&cli.StringFlag{
				Aliases: []string{"r"},
				Name:    "host",
				Usage:   "redis host",
				Value:   defaultRedisHost,
			},
			&cli.IntFlag{
				Aliases: []string{"d"},
				Name:    "db",
				Usage:   "redis db",
				Value:   defaultRedisDB,
			},
			&cli.IntFlag{
				Aliases: []string{"i"},
				Name:    "idle",
				Usage:   "max idle connection to redis",
				Value:   defaultMaxIdle,
			},
			&cli.IntFlag{
				Aliases: []string{"a"},
				Name:    "active",
				Usage:   "max active connection to redis",
				Value:   defaultMaxActive,
			},
			&cli.DurationFlag{
				Aliases: []string{"t"},
				Name:    "timeout",
				Usage:   "idle connection timeout duration",
				Value:   defaultIdleTimeout,
			},
		},
		Action: func(c *cli.Context) error {
			listen := c.String("listen")
			host := c.String("host")
			db := c.Int("db")
			maxIdle := c.Int("idle")
			maxActive := c.Int("active")
			idleTimeout := c.Duration("timeout")
			log.Println("listen:", listen)
			log.Println("host:", host)
			log.Println("db:", db)
			log.Println("maxIdle:", maxIdle)
			log.Println("maxActive:", maxActive)
			log.Println("idleTimeout:", idleTimeout)

			// listen
			lis, err := net.Listen("tcp", listen)
			if err != nil {
				log.Panic(err)
				return err
			}
			log.Info("listening on ", lis.Addr())

			// register service
			s := grpc.NewServer()
			instance := &server{}
			instance.init(host, db, maxIdle, maxActive, idleTimeout)
			pb.RegisterDBServiceServer(s, instance)

			// start service
			s.Serve(lis)

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)
}
