package api

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/kit/packet"
	"github.com/master-g/omgo/kit/services"
	pc "github.com/master-g/omgo/proto/pb/common"
)

// IncomingPacket contains Header and payload bytes
type IncomingPacket struct {
	Header  *pc.Header
	Payload []byte
}

// Config of ETCD and services
type Config struct {
	Root            string   // service root
	Hosts           []string // ETCD hosts
	GameServerName  string   // unlike other service, game server should be specific
	GameServerKind  string   // game server kind, default 'game'
	DataServiceKind string   // data service kind, default 'dataservice'
}

var (
	// Handlers stores request handlers
	Handlers map[int32]func(*Session, *IncomingPacket) []byte
	// Registry stores client session registry
	Registry sync.Map
	// GameServerPool game service pool
	GameServerPool *services.Pool
	// DataServicePool data service pool
	DataServicePool *services.Pool
	// config
	config Config
)

func init() {
	Handlers = map[int32]func(*Session, *IncomingPacket) []byte{
		int32(pc.Cmd_HEART_BEAT_REQ): ProcHeartBeatReq,
		int32(pc.Cmd_HANDSHAKE_REQ):  ProcHandshakeReq,
		int32(pc.Cmd_OFFLINE_REQ):    ProcOfflineReq,
	}
}

// Init services needed by agent
func Init(cfg Config) {
	config = cfg
	GameServerPool = services.New(cfg.Root, cfg.GameServerKind, cfg.Hosts)
	DataServicePool = services.New(cfg.Root, cfg.DataServiceKind, cfg.Hosts)
}

// MakeResponse convert proto message into packet
func MakeResponse(cmd pc.Cmd, msg proto.Message) []byte {
	p := packet.NewRawPacket()
	p.WriteS32(int32(cmd))
	rspBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}
	p.WriteBytes(rspBytes)
	return p.Data()
}
