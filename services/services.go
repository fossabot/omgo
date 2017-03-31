package services

import (
	etcdclient "github.com/coreos/etcd/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// single connection
type client struct {
	key  string
	conn *grpc.ClientConn
}

// service
type service struct {
	clients []client
	idx     uint32
}

// all service
type servicePool struct {
	sync.RWMutex
	root          string
	names         map[string]bool
	services      map[string]*service
	namesProvided bool
	etcdClient    etcdclient.Client
	callbacks     map[string][]chan string
}

var (
	defaultPool servicePool
	once        sync.Once
)

func Init(root string, hosts, services []string) {
	once.Do(func() {
		defaultPool.init(root, hosts, services)
	})
}

func (p *servicePool) init(root string, hosts, services []string) {
	// init etcd client
	cfg := etcdclient.Config{
		Endpoints: hosts,
		Transport: etcdclient.DefaultTransport,
	}
	etcdcli, err := etcdclient.New(cfg)
	if err != nil {
		log.Panic(err)
		os.Exit(-1)
	}
	p.etcdClient = etcdcli
	p.root = root

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
		p.names[p.root+"/"+strings.TrimSpace(v)] = true
	}

	// start connection
	p.connectAll(p.root)
}

// connect to all services
func (p *servicePool) connectAll(directory string) {
	kAPI := etcdclient.NewKeysAPI(p.etcdClient)
	// get the keys under directory
	log.Println("connecting services under:", directory)
	resp, err := kAPI.Get(context.Background(), directory, &etcdclient.GetOptions{Recursive: true})
	if err != nil {
		log.Println(err)
		return
	}

	// validation
	if !resp.Node.Dir {
		log.Println("not a director")
		return
	}

	for _, node := range resp.Node.Nodes {
		if node.Dir {
			for _, service := range node.Nodes {
				p.addService(service.Key, service.Value)
			}
		}
	}
	log.Println("services added")

	go p.watcher()
}

// watcher for data change in etcd directory
func (p *servicePool) watcher() {
	kAPI := etcdclient.NewKeysAPI(p.etcdClient)
	w := kAPI.Watcher(p.root, &etcdclient.WatcherOptions{Recursive: true})
	for {
		resp, err := w.Next(context.Background())
		if err != nil {
			log.Println(err)
			continue
		}
		if resp.Node.Dir {
			continue
		}

		switch resp.Action {
		case "set", "create", "update", "compareAndSwap":
			p.addService(resp.Node.Key, resp.Node.Value)
		case "delete":
			p.removeService(resp.PrevNode.Key)
		}
	}
}

// add a service
func (p *servicePool) addService(key, value string) {
	p.Lock()
	defer p.Unlock()
	// name check
	serviceName := filepath.Dir(key)
	if p.namesProvided && !p.names[serviceName] {
		return
	}

	// try new service kind init
	if p.services[serviceName] == nil {
		p.services[serviceName] = &service{}
	}

	// create service connection
	service := p.services[serviceName]
	if conn, err := grpc.Dial(value, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
		service.clients = append(service.clients, client{key, conn})
		log.Println("service added:", key, "-->", value)
		for k := range p.callbacks[serviceName] {
			select {
			case p.callbacks[serviceName][k] <- key:
			default:
			}
		}
	} else {
		log.Println("did not connect:", key, "-->", value, "error:", err)
	}
}

// remove a service
func (p *servicePool) removeService(key string) {
	p.Lock()
	defer p.Unlock()
	// name check
	serviceName := filepath.Dir(key)
	if p.namesProvided && !p.names[serviceName] {
		return
	}

	// check service kind
	service := p.services[serviceName]
	if service == nil {
		log.Println("no such service:", serviceName)
		return
	}

	// remove service
	for k := range service.clients {
		if service.clients[k].key == key {
			service.clients[k].conn.Close()
			service.clients = append(service.clients[:k], service.clients[k+1:]...)
			log.Println("service removed:", key)
			return
		}
	}
}
