package service

import (
	"sync"
	"time"

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

func (p *Pool) connectAll() {

}
