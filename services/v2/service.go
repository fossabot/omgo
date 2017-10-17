package service

import (
	"context"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync/atomic"
)

// Service modification event
type Event struct {
	Type  int
	Value string
}

// Single gRPC connection client
type Client struct {
	Fullpath string           // full path of the service, root/type/name
	Name     string           //  service name
	Address  string           // service host address, IP:PORT
	Conn     *grpc.ClientConn // gRPC connection
}

// Service pool manages clients
type Pool struct {
	sync.RWMutex                    // sync mutex
	Root         string             // service root
	Kind         string             // service type
	clientArray  []*Client          // client list
	clientMap    map[string]*Client // client map
	idx          uint32             // index for round-robin
	etcdCfg      clientv3.Config    // ETCD client config
	callbacks    []chan Event       // callback list
}

const (
	SEP             = "/"
	RootIndex       = 0
	KindIndex       = 1
	NameIndex       = 2
	EventTypeAdd    = 0
	EventTypeRemove = 1
)

var (
	once           sync.Once
	defaultTimeout = 5 * time.Second
)

// GenPath concat arguments with '/'
func GenPath(arg ...string) string {
	return strings.Join(arg, SEP)
}

// GetRangeKey generate a ranged key from a given key
func GetRangeKey(key string) string {
	rangeKey := make([]byte, len([]byte(key)))
	copy(rangeKey[:], key)
	rangeKey[len(rangeKey)-1]++
	return string(rangeKey)
}

func GetRoot(fullPath string) string {
	return getCompAtIndex(fullPath, RootIndex)
}

func GetKind(fullPath string) string {
	return getCompAtIndex(fullPath, KindIndex)
}

func GetName(fullPath string) string {
	return getCompAtIndex(fullPath, NameIndex)
}

func getCompAtIndex(fullPath string, index int) string {
	ret := strings.Split(fullPath, SEP)
	if index >= len(ret) {
		return ""
	} else {
		return ret[index]
	}
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
		clientArray: make([]*Client, 1),
		clientMap:   make(map[string]*Client),
		callbacks:   make([]chan Event, 1),
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
		p.addService(string(v.Key), string(v.Value))
	}
	log.Println("services added")

	go p.watch()
}

func (p *Pool) addService(fullPath, address string) {
	p.Lock()
	defer p.Unlock()

	client, ok := p.clientMap[fullPath]
	if !ok {
		// prepare client if not exists
		client = &Client{
			Fullpath: fullPath,
			Name:     GetName(fullPath),
			Address:  address,
		}
		p.clientMap[fullPath] = client
	} else if client.Conn.GetState() != connectivity.Shutdown {
		log.Warnf("service already added: %v", fullPath)
		return
	}

	// service client not exists or has been shutdown

	if conn, err := grpc.Dial(address, grpc.WithBlock(), grpc.WithInsecure()); err == nil {
		client.Conn = conn
		p.clientArray = append(p.clientArray, client)
		event := Event{
			Type:  EventTypeAdd,
			Value: fullPath,
		}
		for k := range p.callbacks {
			select {
			case p.callbacks[k] <- event:
			default:
			}
		}
	} else {
		log.Errorf("unable to connect service:%v err:%v", fullPath+" -->"+address, err)
	}

}

func (p *Pool) removeService(fullPath string) {
	p.Lock()
	defer p.Unlock()

	client, ok := p.clientMap[fullPath]
	if ok {
		if client.Conn != nil && client.Conn.GetState() != connectivity.Shutdown {
			client.Conn.Close()
		}
		delete(p.clientMap, fullPath)

		for k, v := range p.clientArray {
			if v == client {
				p.clientArray = append(p.clientArray[:k], p.clientArray[k+1:]...)
				log.Infof("service removed:%v", fullPath)
				break
			}
		}

		// notify callbacks
		event := Event{
			Type:  EventTypeRemove,
			Value: fullPath,
		}
		for k := range p.callbacks {
			select {
			case p.callbacks[k] <- event:
			default:
			}
		}
	} else {
		log.Infof("unable to remove service:%v not exist", fullPath)
	}
}

func (p *Pool) watch() {
	cli, err := clientv3.New(p.etcdCfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()

	rch := cli.Watch(context.Background(), p.Root, p.getRangePathKey())
	for watchRsp := range rch {
		for _, ev := range watchRsp.Events {
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

func (p *Pool) RegisterCallback(callback chan Event) *Pool {
	p.Lock()
	defer p.Unlock()
	if p.callbacks == nil {
		p.callbacks = make([]chan Event, 1)
	}

	p.callbacks = append(p.callbacks, callback)
	return p
}

func (p *Pool) RegisterService(fullPath, address string) {
	cli, err := clientv3.New(p.etcdCfg)
	if err != nil {
		log.Error(err)
		return
	}
	defer cli.Close()

	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	_, err = cli.Put(ctx, fullPath, address)
	cancel()
	if err != nil {
		log.Error(err)
		return
	}
	log.Infof("put key %v, value %v", fullPath, address)
}

func (p *Pool) NextClient() (conn *grpc.ClientConn, key string) {
	p.RLock()
	defer p.RUnlock()
	if len(p.clientArray) == 0 {
		return nil, ""
	}

	idx := int(atomic.AddUint32(&p.idx, 1)) % len(p.clientArray)
	return p.clientArray[idx].Conn, p.clientArray[idx].Fullpath
}

func main() {
	endpoints := make([]string, 1)
	endpoints[0] = "192.168.0.22:2379"
	pool := New("backends", "test", endpoints)

}
