package api

import (
	"sync"

	"github.com/master-g/omgo/kit/packet"
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
	// client request handlers
	Handlers map[int32]func(*Session, *packet.RawPacket) []byte
	// client session registry
	Registry sync.Map
	// game server pool
	GameServerPool *services.Pool
	// data service pool
	DataServicePool *services.Pool
	// config
	config Config
)

func init() {
	Handlers = map[int32]func(*Session, *packet.RawPacket) []byte{
		int32(pc.Cmd_HEART_BEAT_REQ): ProcHeartBeatReq,
		int32(pc.Cmd_LOGIN_REQ):      ProcUserLoginReq,
		int32(pc.Cmd_GET_SEED_REQ):   ProcGetSeedReq,
		int32(pc.Cmd_OFFLINE_REQ):    ProcOfflineReq,
	}
}

func Init(cfg Config) {
	config = cfg
	GameServerPool = services.New(cfg.Root, cfg.GameServerKind, cfg.Hosts)
	DataServicePool = services.New(cfg.Root, cfg.DataServiceKind, cfg.Hosts)
}
