package main

import (
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"sync"
	"time"

	"context"
	log "github.com/sirupsen/logrus"
	"strings"
)

// single connection
type client struct {
	key  string
	conn *grpc.ClientConn
}

// serivce
type service struct {
	clients []client
	idx     uint32
}

// all services
type servicePool struct {
	sync.RWMutex
	root          string
	names         map[string]bool
	services      map[string]*service
	namesProvided bool
	callbacks     map[string][]chan string
	etcdCfg       clientv3.Config
}

var (
	defaultPool    servicePool
	once           sync.Once
	pathSep        = "/"
	defaultTimeout = 5 * time.Second
)

// Init service pool with given service root on ETCD hosts
// {root}/{services}/{service-endpoints}
func Init(root string, hosts, services []string) {
	once.Do(func() {
		defaultPool.init(root, hosts, services)
	})
}

func (p *servicePool) init(root string, hosts, services []string) {
	// init etcd client
	p.etcdCfg = clientv3.Config{
		Endpoints:   hosts,
		DialTimeout: defaultTimeout,
	}

	// init
	p.services = make(map[string]*service)
	p.names = make(map[string]bool)

	// names init
	names := services
	if len(names) > 0 {
		p.namesProvided = true
	}

	log.Println("all service names:", names)
	for _, v := range names {
		p.names[p.root+pathSep+strings.TrimSpace(v)] = true
	}

	// start connection
	p.connectAll(p.root)
}

// connect to all services
func (p *servicePool) connectAll(directory string) {
	// get etcd v3 client
	cli, err := clientv3.New(p.etcdCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	resp, err := cli.Get(ctx, directory, clientv3.WithFromKey())
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range resp.Kvs {
		log.Println(k, v)
	}
}
