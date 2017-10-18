package api

import (
	"github.com/master-g/omgo/agent/session"
	"github.com/master-g/omgo/kit/packet"
	pc "github.com/master-g/omgo/proto/pb/common"
)

var Handlers map[int32]func(*session.Session, *packet.RawPacket) []byte

func init() {
	Handlers = map[int32]func(*session.Session, *packet.RawPacket) []byte{
		int32(pc.Cmd_HEART_BEAT_REQ): ProcHeartBeatReq,
		int32(pc.Cmd_LOGIN_REQ):      ProcUserLoginReq,
		int32(pc.Cmd_GET_SEED_REQ):   ProcGetSeedReq,
		int32(pc.Cmd_OFFLINE_REQ):    ProcOfflineReq,
	}
}
