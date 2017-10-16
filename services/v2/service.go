package service

import (
	"context"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
)

// Single gRPC connection client
type Client struct {
	Fullpath string           // full path of the service, root/type/name
	Name     string           //  service name
	Address  string           // service host address, IP:PORT
	Conn     *grpc.ClientConn // gRPC connection
}

// Service pool manages clients
type Pool struct {
	sync.RWMutex                   // sync mutex
	Root         string            // service root
	Kind         string            // service type
	clientArray  []Client          // client list
	clientMap    map[string]Client // client map
	listeners    []chan string     // listener list
	idx          uint32            // index for round-robin
	etcdCfg      clientv3.Config   // ETCD client config
}

var (
	once           sync.Once
	defaultTimeout = 5 * time.Second
)

func GenPath(arg ...string) string {
	return strings.Join(arg, "/")
}

// GetRangeKey generate a ranged key from a given key
func GetRangeKey(key string) string {
	rangeKey := make([]byte, len([]byte(key)))
	copy(rangeKey[:], key)
	rangeKey[len(rangeKey)-1]++
	return string(rangeKey)
}

// New will create a new service pool instance
func New(root, kind string, hosts []string) Pool {
	etcdCfg := clientv3.Config{
		Endpoints:   hosts,
		DialTimeout: defaultTimeout,
	}

	pool := Pool{
		Root:        root,
		Kind:        kind,
		clientArray: make([]Client, 1),
		clientMap:   make(map[string]Client),
		listeners:   make([]chan string, 1),
		idx:         0,
		etcdCfg:     etcdCfg,
	}

	pool.connectAll()

	return pool
}

func (p *Pool) getRangePathKey() clientv3.OpOption {
	return clientv3.WithRange(GetRangeKey(GenPath(p.Root, p.Kind)))
}

// connect all the service that is under the root/kind/
func (p *Pool) connectAll() {
	// get etcd v3 client
	cli, err := clientv3.New(p.etcdCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	resp, err := cli.Get(ctx, p.Root, p.getRangePathKey())
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range resp.Kvs {
		p.addService(v.Key, v.Value)
	}
	log.Println("services added")

	go p.watch()
}

func (p *Pool) addService(fullPath, address string) {

}

func (p *Pool) removeService(fullPath string) {

}

func (p *Pool) watch() {
	cli, err := clientv3.New(p.etcdCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), p.Root, p.getRangePathKey())
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
