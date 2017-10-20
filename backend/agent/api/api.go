package api

import (
	"github.com/master-g/omgo/kit/packet"
	"github.com/master-g/omgo/kit/services"
	pc "github.com/master-g/omgo/proto/pb/common"
)

var (
	Handlers           map[int32]func(*Session, *packet.RawPacket) []byte
	gameServerPool     *services.Pool
	gameServerFullPath string
	gameServerName     string
)

func init() {
	Handlers = map[int32]func(*Session, *packet.RawPacket) []byte{
		int32(pc.Cmd_HEART_BEAT_REQ): ProcHeartBeatReq,
		int32(pc.Cmd_LOGIN_REQ):      ProcUserLoginReq,
		int32(pc.Cmd_GET_SEED_REQ):   ProcGetSeedReq,
		int32(pc.Cmd_OFFLINE_REQ):    ProcOfflineReq,
	}
}

func Init(root, kind, name string, etcdHosts []string) {
	gameServerName = name
	gameServerFullPath = services.GenPath(root, kind, name)
	gameServerPool = services.New(root, kind, etcdHosts)
}
