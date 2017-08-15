package main

import (
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"sync"
	"time"

	"context"
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"sync/atomic"
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
	p.root = root
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
	for _, v := range resp.Kvs {
		p.addService(string(v.Key), string(v.Value))
	}
	log.Println("services added")

	go p.watcher()
}

// watcher for data change in ETCD
func (p *servicePool) watcher() {
	cli, err := clientv3.New(p.etcdCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), p.root, clientv3.WithFromKey())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			log.Printf("%s %q : %q\n", ev.Type, ev.Kv.Key, ev.Kv.Value)
			switch ev.Type {
			case clientv3.EventTypePut:
				p.addService(string(ev.Kv.Key), string(ev.Kv.Value))
			case clientv3.EventTypeDelete:
				p.removeService(string(ev.Kv.Key))
			}
		}
	}
}

// add a service
func (p *servicePool) addService(servicePath, address string) {
	p.Lock()
	defer p.Unlock()
	// name check
	serviceKind := filepath.Dir(servicePath)
	if p.namesProvided && !p.names[serviceKind] {
		return
	}

	// try new service kind init
	if p.services[serviceKind] == nil {
		p.services[serviceKind] = &service{}
	}

	// create service connections
	service := p.services[serviceKind]
	if conn, err := grpc.Dial(address, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
		service.clients = append(service.clients, client{servicePath, conn})
		log.Println("service added:", servicePath, "-->", address)
		for k := range p.callbacks[serviceKind] {
			select {
			case p.callbacks[serviceKind][k] <- servicePath:
			default:
			}
		}
	} else {
		log.Println("did not connect:", servicePath, "-->", address, "error:", err)
	}
}

// remove service
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

// getServiceWithID returns a connection to a specific service with the id given
// eg:
// path: backends/snowflake, id: s1
// the full canonical path for this service is :
// backends/snowflake/s1
func (p *servicePool) getServiceWithID(path, id string) *grpc.ClientConn {
	p.RLock()
	defer p.RUnlock()
	// check
	service := p.services[path]
	if service == nil {
		return nil
	}
	if len(service.clients) == 0 {
		return nil
	}

	fullPath := string(path) + pathSep + id
	for k := range service.clients {
		if service.clients[k].key == fullPath {
			return service.clients[k].conn
		}
	}

	return nil
}

// getService returns a service in round-robin style
func (p *servicePool) getService(path string) (conn *grpc.ClientConn, key string) {
	p.RLock()
	defer p.RUnlock()
	// check
	service := p.services[path]
	if service == nil {
		return nil, ""
	}

	if len(service.clients) == 0 {
		return nil, ""
	}

	idx := int(atomic.AddUint32(&service.idx, 1)) % len(service.clients)
	return service.clients[idx].conn, service.clients[idx].key
}

// registerCallback will add callback to specific service when its added or removed
func (p *servicePool) registerCallback(path string, callback chan string) {
	p.Lock()
	defer p.Unlock()
	if p.callbacks == nil {
		p.callbacks = make(map[string][]chan string)
	}

	p.callbacks[path] = append(p.callbacks[path], callback)
	if s, ok := p.services[path]; ok {
		for k := range s.clients {
			callback <- s.clients[k].key
		}
	}
	log.Println("register callback on:", path)
}

func (p *servicePool) registerService(path, address string) {
	// get etcd v3 client
	cli, err := clientv3.New(p.etcdCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err = cli.Put(ctx, path, address)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	log.Infof("service %v --> %v registered", path, address)
}

// GetService finds gRPC service with path in service pool
func GetService(path string) *grpc.ClientConn {
	conn, _ := defaultPool.getService(defaultPool.root + pathSep + path)
	return conn
}

// GetServiceAndKey finds gRPC service and its key with path in service pool
func GetServiceAndKey(path string) (*grpc.ClientConn, string) {
	conn, key := defaultPool.getService(defaultPool.root + pathSep + path)
	return conn, key
}

// GetServiceWithID finds gRPC service with path and its ID
func GetServiceWithID(path, id string) *grpc.ClientConn {
	return defaultPool.getServiceWithID(defaultPool.root+pathSep+path, id)
}

// RegisterCallback register callback at give path, once changes are made, callbacks will be invoke
func RegisterCallback(path string, callback chan string) {
	defaultPool.registerCallback(defaultPool.root+pathSep+path, callback)
}

func RegisterService(path string, address string) {
	defaultPool.registerService(path, address)
}
