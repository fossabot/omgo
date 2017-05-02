package main

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	etcd "github.com/coreos/etcd/client"
	pb "github.com/master-g/omgo/backend/snowflake/proto"
	"github.com/master-g/omgo/etcdclient"
	"golang.org/x/net/context"
	"math/rand"
	"os"
	"strconv"
	"time"
)

const (
	envMachineID = "MACHINE_ID" // Specific machine id
	etcdPath     = "/seqs/"
	uuidKey      = "/seqs/snowflake-uuid"
	backoff      = 100  // Max backoff delay millisecond
	concurrent   = 128  // Max concurrent connections to etcd
	uuidQueue    = 1024 // UUID process queue
)

const (
	tsMask        = 0x1FFFFFFFFFF // 41bit
	snMask        = 0xFFF         // 12bit
	machineIDmask = 0x3FF         // 10bit
)

type server struct {
	machineID  uint64 // 10-bit machine id
	clientPool chan etcd.KeysAPI
	chProc     chan chan uint64
}

func (s *server) init() {
	s.clientPool = make(chan etcd.KeysAPI, concurrent)
	s.chProc = make(chan chan uint64, uuidQueue)

	// Init client pool
	for i := 0; i < concurrent; i++ {
		s.clientPool <- etcdclient.KeysAPI()
	}

	// Check if user specified machine id is set
	if env := os.Getenv(envMachineID); env != "" {
		if id, err := strconv.Atoi(env); err == nil {
			s.machineID = (uint64(id) & machineIDmask) << 12
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
		resp, err := client.Get(context.Background(), uuidKey, nil)
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}

		// Get prevValue & prevIndex
		prevValue, err := strconv.Atoi(resp.Node.Value)
		if err != nil {
			log.Panic(err)
			os.Exit(-1)
		}
		prevIndex := resp.Node.ModifiedIndex

		// CompareAndSwap
		resp, err = client.Set(context.Background(), uuidKey, fmt.Sprint(prevValue+1), &etcd.SetOptions{PrevIndex: prevIndex})
		if err != nil {
			casDelay()
			continue
		}

		// record serial number of this service, already shifted
		s.machineID = (uint64(prevValue+1) & machineIDmask) << 12
		return
	}
}

// Get next value of a key, like auto-increment in mysql
func (s *server) Next(ctx context.Context, in *pb.Snowflake_Key) (*pb.Snowflake_Value, error) {
	client := <-s.clientPool
	defer func() { s.clientPool <- client }()
	key := etcdPath + in.Name
	for {
		// Get the key
		resp, err := client.Get(context.Background(), key, nil)
		if err != nil {
			log.Error(err)
			return nil, fmt.Errorf("Key:%v not exists, need to create first", key)
		}

		// Get prevValue & prevIndex
		prevValue, err := strconv.Atoi(resp.Node.Value)
		if err != nil {
			log.Error(err)
			return nil, errors.New("Marlformed value")
		}
		prevIndex := resp.Node.ModifiedIndex

		// CompareAndSwap
		resp, err = client.Set(context.Background(), key, fmt.Sprint(prevValue+1), &etcd.SetOptions{PrevIndex: prevIndex})
		if err != nil {
			casDelay()
			continue
		}
		return &pb.Snowflake_Value{int64(prevValue + 1)}, nil
	}
}

// Generate an unique uuid
func (s *server) GetUUID(context.Context, *pb.Snowflake_NullRequest) (*pb.Snowflake_UUID, error) {
	req := make(chan uint64, 1)
	s.chProc <- req
	return &pb.Snowflake_UUID{<-req}, nil
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
// random delay
func casDelay() {
	<-time.After(time.Duration(rand.Int63n(backoff)) * time.Millisecond)
}

// get timestamp
func ts() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}
