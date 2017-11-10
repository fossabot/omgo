package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/etcd/clientv3"
	pb "github.com/master-g/omgo/proto/grpc/snowflake"
	"golang.org/x/net/context"
)

const (
	envMachineID   = "MACHINE_ID" // Specific machine id
	pathETCD       = "seqs/"
	uuidKey        = "seqs/snowflake-uuid"
	reduction      = 100  // Max reduction delay millisecond
	concurrentETCD = 128  // Max concurrent connections to ETCD
	uuidQueueSize  = 1024 // UUID process queue
	requestTimeout = 5 * time.Second
)

const (
	tsMask        = 0x1FFFFFFFFFF // 41bit
	snMask        = 0xFFF         // 12bit
	machineIDMask = 0x3FF         // 10bit
)

type server struct {
	machineID  uint64 // 10-bit machine id
	clientPool chan *clientv3.Client
	chProc     chan chan uint64
}

var (
	etcdCfg clientv3.Config
)

func (s *server) init(endpoints []string) {
	etcdCfg = clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: requestTimeout,
	}

	s.clientPool = make(chan *clientv3.Client, concurrentETCD)
	s.chProc = make(chan chan uint64, uuidQueueSize)

	// Init client pool
	for i := 0; i < concurrentETCD; i++ {
		cli, err := clientv3.New(etcdCfg)
		if err != nil {
			log.Fatal(err)
		}
		s.clientPool <- cli
	}

	// Check if user specified machine id is set
	if env := os.Getenv(envMachineID); env != "" {
		if id, err := strconv.Atoi(env); err == nil {
			s.machineID = (uint64(id) & machineIDMask) << 12
			log.Info("machine id specified:", id)
		} else {
			log.Panic(err)
			os.Exit(-1)
		}
	} else {
		s.initMachineID()
	}

	go s.uuidTask()
}

func (s *server) initMachineID() {
	client := <-s.clientPool
	defer func() { s.clientPool <- client }()

	for {
		// Get the key
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		resp, err := client.Get(ctx, uuidKey)
		cancel()
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}

		if len(resp.Kvs) == 0 {
			log.Panic("uuid key missing")
			os.Exit(-1)
		}

		kv := resp.Kvs[0]
		// Get prevValue & prevIndex
		prevValue, err := strconv.Atoi(string(kv.Value))
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}

		// CompareAndSwap
		_, err = clientv3.NewKV(client).Txn(context.Background()).If(
			clientv3.Compare(clientv3.ModRevision(uuidKey), "=", kv.ModRevision),
		).Then(
			clientv3.OpPut(uuidKey, fmt.Sprint(prevValue+1)),
			clientv3.OpGet(uuidKey),
		).Commit()

		// newer version exist
		if err != nil {
			casDelay()
			continue
		}

		// record serial number of this service, already shifted
		s.machineID = (uint64(prevValue+1) & machineIDMask) << 12
		return
	}
}

// Get next value of a key, like auto-increment in mysql
func (s *server) Next(ctx context.Context, in *pb.Snowflake_Key) (*pb.Snowflake_Value, error) {
	client := <-s.clientPool
	defer func() { s.clientPool <- client }()
	key := pathETCD + in.GetName()
	for {
		// Get the key
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		resp, err := client.Get(ctx, key)
		cancel()
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}

		if len(resp.Kvs) == 0 {
			log.Panic("uuid key missing")
			os.Exit(-1)
		}

		kv := resp.Kvs[0]
		// Get prevValue & prevIndex
		prevValue, err := strconv.Atoi(string(kv.Value))
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}

		// CompareAndSwap
		_, err = clientv3.NewKV(client).Txn(context.Background()).If(
			clientv3.Compare(clientv3.ModRevision(key), "=", kv.ModRevision),
		).Then(
			clientv3.OpPut(key, fmt.Sprint(prevValue+1)),
			clientv3.OpGet(key),
		).Commit()

		// newer version exist
		if err != nil {
			casDelay()
			continue
		}

		return &pb.Snowflake_Value{Value: int64(prevValue + 1)}, nil
	}
}

// Generate an user id
func (s *server) Next2(ctx context.Context, param *pb.Snowflake_Param) (*pb.Snowflake_Value, error) {
	client := <-s.clientPool
	defer func() { s.clientPool <- client }()
	key := pathETCD + param.GetName()
	for {
		// Get the key
		ctx, cancel := context.WithTimeout(context.Background(), requestTimeout)
		resp, err := client.Get(ctx, key)
		cancel()
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}

		if len(resp.Kvs) == 0 {
			log.Panic("uuid key missing")
			os.Exit(-1)
		}

		kv := resp.Kvs[0]
		// Get prevValue & prevIndex
		prevValue, err := strconv.Atoi(string(kv.Value))
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}

		currentValue := int64(0)
		if param.Step != 0 {
			currentValue = int64(prevValue) + param.Step
		} else {
			currentValue = int64(prevValue + rand.Intn(2048) + 1)
		}

		// CompareAndSwap
		_, err = clientv3.NewKV(client).Txn(context.Background()).If(
			clientv3.Compare(clientv3.ModRevision(key), "=", kv.ModRevision),
		).Then(
			clientv3.OpPut(key, fmt.Sprint(currentValue)),
			clientv3.OpGet(key),
		).Commit()

		// newer version exist
		if err != nil {
			casDelay()
			continue
		}

		return &pb.Snowflake_Value{Value: currentValue}, nil
	}
}

// Generate an unique uuid
func (s *server) GetUUID(context.Context, *pb.Snowflake_NullRequest) (*pb.Snowflake_UUID, error) {
	req := make(chan uint64, 1)
	s.chProc <- req
	return &pb.Snowflake_UUID{Uuid: <-req}, nil
}

// UUID generator
func (s *server) uuidTask() {
	var sn uint64    // 12-bit serial no
	var lastTs int64 // last timestamp
	for {
		ret := <-s.chProc
		// get a correct serial number
		t := ts()
		if t < lastTs { // clock shift backward
			log.Error("clock shift happened, waiting until the clock moving to the next millisecond.")
			t = s.waitMilliseconds(lastTs)
		}

		if lastTs == t { // same millisecond
			sn = (sn + 1) & snMask
			if sn == 0 { // serial number overflows, wait until next ms
				t = s.waitMilliseconds(lastTs)
			}
		} else { // new millisecond, reset serial number to 0
			sn = 0
		}
		// remember last timestamp
		lastTs = t

		// generate uuid, format:
		//
		// 0		0.................0		0..............0	0........0
		// 1-bit	41bit timestamp			10bit machine-id	12bit sn
		var uuid uint64
		uuid |= (uint64(t) & tsMask) << 22
		uuid |= s.machineID
		uuid |= sn
		ret <- uuid
	}
}

// wait_ms will spin wait till next millisecond.
func (s *server) waitMilliseconds(lastTs int64) int64 {
	t := ts()
	for t <= lastTs {
		t = ts()
	}
	return t
}

////////////////////////////////////////////////////////////////////////////////
// CompareAndSwap delay
func casDelay() {
	<-time.After(time.Duration(rand.Int63n(reduction)) * time.Millisecond)
}

// get timestamp
func ts() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
