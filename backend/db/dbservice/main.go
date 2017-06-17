package main

import (
	"net"
	"os"
	"sort"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/db/dbservice/driver"
	pb "github.com/master-g/omgo/proto/grpc/db"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"gopkg.in/urfave/cli.v2"
)

const (
	// gRPC
	defaultListen = ":60001"
	// mongodb
	defaultMongoHost          = ":37017"
	defaultMongoTimeout       = 60 * time.Second
	defaultMongoSocketTimeout = 10 * time.Second
	defaultMongoDatabase      = "master"
	defaultMongoUserName      = "admin"
	defaultMongoPassword      = "admin"
	defaultMongoConcurrent    = 128
	// redis
	defaultRedisHost        = ":6379" // FIXME: DO NOT use 6379 in production environment
	defaultRedisDB          = 0
	defaultRedisMaxIdle     = 80
	defaultRedisMaxActive   = 1024
	defaultRedisIdleTimeout = 180 * time.Second
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
				Name:    "redis-host",
				Usage:   "redis host",
				Value:   defaultRedisHost,
			},
			&cli.IntFlag{
				Aliases: []string{"d"},
				Name:    "redis-db",
				Usage:   "redis db",
				Value:   defaultRedisDB,
			},
			&cli.IntFlag{
				Aliases: []string{"i"},
				Name:    "idle",
				Usage:   "max idle connection to redis",
				Value:   defaultRedisMaxIdle,
			},
			&cli.IntFlag{
				Aliases: []string{"a"},
				Name:    "active",
				Usage:   "max active connection to redis",
				Value:   defaultRedisMaxActive,
			},
			&cli.DurationFlag{
				Aliases: []string{"t"},
				Name:    "timeout",
				Usage:   "idle connection timeout duration",
				Value:   defaultRedisIdleTimeout,
			},
			&cli.StringSliceFlag{
				Aliases: []string{"m"},
				Name:    "mongo-host",
				Usage:   "mongodb host",
				Value:   cli.NewStringSlice(defaultMongoHost),
			},
			&cli.DurationFlag{
				Aliases: []string{"o"},
				Name:    "mongo-timeout",
				Usage:   "mongodb connect timeout",
				Value:   defaultMongoTimeout,
			},
			&cli.StringFlag{
				Aliases: []string{"b"},
				Name:    "mongo-database",
				Usage:   "mongodb-database",
				Value:   defaultMongoDatabase,
			},
			&cli.StringFlag{
				Aliases: []string{"u"},
				Name:    "mongo-username",
				Usage:   "mongodb user name",
				Value:   defaultMongoUserName,
			},
			&cli.StringFlag{
				Aliases: []string{"p"},
				Name:    "mongo-password",
				Usage:   "mongodb password",
				Value:   defaultMongoPassword,
			},
			&cli.IntFlag{
				Aliases: []string{"c"},
				Name:    "mongo-concurrent",
				Usage:   "mongodb concurrent pool size",
				Value:   defaultMongoConcurrent,
			},
			&cli.DurationFlag{
				Aliases: []string{"s"},
				Name:    "mongo-socket-timeout",
				Usage:   "mongodb socket timeout",
				Value:   defaultMongoSocketTimeout,
			},
		},
		Action: func(c *cli.Context) error {
			// gRPC
			listen := c.String("listen")
			// redis
			redisHost := c.String("redis-host")
			redisDB := c.Int("redis-db")
			redisMaxIdle := c.Int("idle")
			redisMaxActive := c.Int("active")
			redisIdleTimeout := c.Duration("timeout")
			// mongoDB
			mongoHost := c.StringSlice("mongo-host")
			mongoTimeout := c.Duration("mongo-timeout")
			mongoDatabase := c.String("mongo-database")
			mongoUsername := c.String("mongo-username")
			mongoPassword := c.String("mongo-password")
			mongoConcurrent := c.Int("mongo-concurrent")
			mongoSocketTimeout := c.Duration("mongo-socket-timeout")

			log.Println("listen:", listen)
			log.Println("redis host:", redisHost)
			log.Println("redis db:", redisDB)
			log.Println("redis max idle:", redisMaxIdle)
			log.Println("redis max active:", redisMaxActive)
			log.Println("redis timeout:", redisIdleTimeout)
			log.Println("mongo host:", mongoHost)
			log.Println("mongo timeout:", mongoTimeout)
			log.Println("mongo database:", mongoDatabase)
			log.Println("mongo username:", mongoUsername)
			log.Println("mongo password:", mongoPassword)
			log.Println("mongo concurrent:", mongoConcurrent)
			log.Println("mongo socket timeout:", mongoSocketTimeout)

			redisCfg := &driver.Config{
				Host:        redisHost,
				DB:          redisDB,
				MaxIdle:     redisMaxIdle,
				MaxActive:   redisMaxActive,
				IdleTimeout: redisIdleTimeout,
			}

			mongoCfg := &mgo.DialInfo{
				Addrs:    mongoHost,
				Timeout:  mongoTimeout,
				Database: mongoDatabase,
				Username: mongoUsername,
				Password: mongoPassword,
			}

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
			instance.init(mongoCfg, mongoConcurrent, mongoSocketTimeout, redisCfg)
			pb.RegisterDBServiceServer(s, instance)

			// start service
			s.Serve(lis)

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)
}
