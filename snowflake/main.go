/*
 * This is a standard gRPC server implementation
 * the Snowflake.ServiceServer is implement in service.go
 */
package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/golandlord/snowflake/etcdclient"
	pb "github.com/master-g/golandlord/snowflake/proto"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
	"net"
	"os"
	"sort"
	"strings"
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
			&cli.StringFlag{
				Name:    "etcd",
				Value:   "",
				Aliases: []string{"e"},
				Usage:   "etcd server address, if there are multiple servers, use ';' to separate",
				EnvVars: []string{"ETCD_HOST"},
			},
		},
		Name:    "snowflake",
		Usage:   "Twitter's UUID generator snowflake in golang",
		Version: "v1.0.0",
		Action: func(c *cli.Context) error {
			port := c.Int("port")
			etcdHosts := c.String("etcd")
			endpoints := []string{DEFAULT_ETCD}
			if etcdHosts != "" {
				endpoints = strings.Split(etcdHosts, ";")
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
