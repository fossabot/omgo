package main

import (
	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/proto/grpc/db"
	proto_common "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"gopkg.in/urfave/cli.v2"
	"net"
	"os"
	"sort"
	"time"
)

const (
	// gRPC
	defaultListen = ":60001"
	// mongodb
	defaultMongoHost     = "127.0.0.1:37017"
	defaultMongoTimeout  = 60 * time.Second
	defaultMongoDatabase = "master"
	defaultMongoUserName = "admin"
	defaultMongoPassword = "admin"
	// redis
	defaultRedisHost        = "127.0.0.1:6379" // DO NOT use 6379 in production environment
	defaultRedisDB          = 0
	defaultRedisMaxIdle     = 80
	defaultRedisMaxActive   = 1024
	defaultRedisIdleTimeout = 180 * time.Second
)

func testMongoDB(dialInfo *mgo.DialInfo) {
	mongoSession, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		log.Fatal(err)
	}
	session := mongoSession.Copy()
	defer session.Close()

	c := session.DB("master").C("users")
	if c == nil {
		log.Fatal("shit happens")
	}

	userInfo := proto_common.UserBasicInfo{
		Usn:      10001,
		Uid:      10002,
		Birthday: 12345,
		Gender:   proto_common.Gender_GENDER_FEMALE,
		Nickname: "Anna",
		Email:    "anna@acme.com",
		Avatar:   "gg",
		Country:  "cn",
	}

	v, err := c.Upsert(bson.M{"usn": userInfo.Usn}, userInfo)
	if err != nil {
		log.Fatal(err)
	}
	log.Info(v)
}

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
				Value:   []string{defaultMongoHost},
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

			redisCfg := &redisConfig{
				host:        redisHost,
				db:          redisDB,
				maxIdle:     redisMaxIdle,
				maxActive:   redisMaxActive,
				idleTimeout: redisIdleTimeout,
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
			instance.init(mongoCfg, redisCfg)
			pb.RegisterDBServiceServer(s, instance)

			// start service
			s.Serve(lis)

			testMongoDB(mongoCfg)

			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)
}
