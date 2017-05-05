package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/etcdclient"
	pb "github.com/master-g/omgo/proto/grpc/snowflake"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
	"net"
	"os"
	"sort"
)

const (
	defaultETCD = "http://127.0.0.1:2379"
)

func main() {
	log.SetLevel(log.DebugLevel)
	defer utils.PrintPanicStack()

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Aliases:     []string{"p"},
				DefaultText: "random",
				Name:        "port",
				Usage:       "local port to listen",
				Value:       0,
			},
			&cli.StringSliceFlag{
				Aliases: []string{"e"},
				EnvVars: []string{"ETCD_HOST"},
				Name:    "etcd",
				Usage:   "etcd server address, if multiple hosts, -e host1 -e host2 ...",
				Value:   cli.NewStringSlice(defaultETCD),
			},
		},
		Name:    "snowflake",
		Usage:   "Twitter's UUID generator snowflake in golang",
		Version: "v1.0.0",
		Action: func(c *cli.Context) error {
			port := c.Int("port")
			etcdHosts := c.StringSlice("etcd")
			log.Infof("start snowflake with etcd hosts:%v", etcdHosts)
			startSnowflake(etcdHosts, port)
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)
}

func startSnowflake(endpoints []string, port int) {
	// etcd client
	etcdclient.Init(endpoints)

	// listen
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	log.Info("listening on ", listener.Addr())

	// register service
	s := grpc.NewServer()
	instance := &server{}
	instance.init()
	pb.RegisterSnowflakeServiceServer(s, instance)

	// Start service
	s.Serve(listener)
}
