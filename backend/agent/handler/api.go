package handler

import (
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
	proto_common "github.com/master-g/omgo/proto/pb/common"
)

var Handlers map[int32]func(*types.Session, *packet.RawPacket) []byte

func init() {
	Handlers = map[int32]func(*types.Session, *packet.RawPacket) []byte{
		int32(proto_common.Cmd_HEART_BEAT_REQ): ProcHeartBeatReq,
		int32(proto_common.Cmd_LOGIN_REQ):      ProcUserLoginReq,
		int32(proto_common.Cmd_GET_SEED_REQ):   ProcGetSeedReq,
	}
}
