package api

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/kit/services"
	pc "github.com/master-g/omgo/proto/pb/common"
)

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
	Handlers map[int32]func(*Session, []byte) *OutgoingPacket
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
	Handlers = map[int32]func(*Session, []byte) *OutgoingPacket{
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
func MakeResponse(hdr *pc.RspHeader, msg proto.Message) *OutgoingPacket {
	p := &OutgoingPacket{
		Header: hdr,
	}
	rspBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}
	p.Body = rspBytes
	return p
}
