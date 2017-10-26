package main

import (
	"fmt"
	"net"
	"os"
	"sort"

	"context"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
	pb "github.com/master-g/omgo/proto/grpc/snowflake"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc"
	"gopkg.in/urfave/cli.v2"
)

const (
	defaultETCD = "http://127.0.0.1:2379"
	defaultSID  = "backends/snowflake/snowflake-0"
	defaultHOST = "127.0.0.1:40001"
)

func main() {
	log.SetLevel(log.DebugLevel)
	defer utils.PrintPanicStack()

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.IntFlag{
				DefaultText: "random",
				Name:        "port",
				Usage:       "local port to listen",
				Value:       0,
			},
			&cli.StringSliceFlag{
				EnvVars: []string{"ETCD_HOST"},
				Name:    "etcd",
				Usage:   "etcd server address, if multiple hosts, -e host1 -e host2 ...",
				Value:   cli.NewStringSlice(defaultETCD),
			},
			&cli.StringFlag{
				DefaultText: defaultSID,
				Name:        "service-key",
				Usage:       "service key",
			},
			&cli.StringFlag{
				DefaultText: defaultHOST,
				Name:        "service-host",
				Usage:       "service host",
			},
		},
		Name:    "snowflake",
		Usage:   "Twitter's UUID generator snowflake in golang",
		Version: "v1.0.0",
		Action: func(c *cli.Context) error {
			port := c.Int("port")
			etcdHosts := c.StringSlice("etcd")
			key := c.String("service-key")
			host := c.String("service-host")
			log.Infof("start snowflake with etcd hosts:%v", etcdHosts)
			log.Infof("service key:%v host:%v", key, host)

			setupETCD(etcdHosts, key, host, port)
			// start snowflake service
			startSnowflake(etcdHosts, port)
			return nil
		},
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	app.Run(os.Args)
}

func setupETCD(endpoints []string, key, host string, port int) {
	// connect to etcd
	log.Infof("connecting to ETCD: %v", endpoints)
	etcdCli, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: time.Second * 5,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer etcdCli.Close()

	// register snowflake to etcd
	log.Infof("register self to ETCD : %v @ %v:%v", key, host, port)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	_, err = etcdCli.Put(ctx, key, fmt.Sprintf("%s:%d", host, port))
	cancel()
	if err != nil {
		log.Error(err)
	}

	// setup snowflake key-values on etcd
	casPut(etcdCli, "seqs/snowflake-uuid", "0")
	casPut(etcdCli, "seqs/test_key", "0")
	casPut(etcdCli, "seqs/userid", "10000")
}

func casPut(client *clientv3.Client, key, value string) {
	_, err := clientv3.NewKV(client).Txn(context.Background()).If(
		clientv3.Compare(clientv3.ModRevision(key), "=", 0),
	).Then(
		clientv3.OpPut(key, value),
	).Commit()

	if err != nil {
		log.Error(err)
	}
}

func startSnowflake(endpoints []string, port int) {
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
	instance.init(endpoints)
	pb.RegisterSnowflakeServiceServer(s, instance)

	// start service
	s.Serve(listener)
}
