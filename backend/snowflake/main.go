/*
 * This is a standard gRPC server implementation
 * the Snowflake.ServiceServer is implement in service.go
 */
package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/backend/snowflake/proto"
	"github.com/master-g/omgo/etcdclient"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
	"net"
	"os"
	"sort"
)

const (
	DEFAULT_ETCD = "http://127.0.0.1:2379"
)

func main() {
	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:        "port",
				Value:       0,
				Aliases:     []string{"p"},
				Usage:       "local port to listen",
				DefaultText: "random",
			},
			&cli.StringSliceFlag{
				Name:    "etcd",
				Value:   cli.NewStringSlice(DEFAULT_ETCD),
				Aliases: []string{"e"},
				Usage:   "etcd server address, if multiple hosts, -e host1 -e host2 ...",
				EnvVars: []string{"ETCD_HOST"},
			},
		},
		Name:    "snowflake",
		Usage:   "Twitter's UUID generator snowflake in golang",
		Version: "v1.0.0",
		Action: func(c *cli.Context) error {
			port := c.Int("port")
			etcdHosts := c.StringSlice("etcd")
			endpoints := []string{DEFAULT_ETCD}
			if etcdHosts == nil {
				endpoints = etcdHosts
			}
			startSnowflake(endpoints, port)
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

// GetLocalIP returns the non loopback local IP of the host
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return ""
}
