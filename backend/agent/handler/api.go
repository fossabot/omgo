package handler

import (
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/proto/pb"
)

const (
	ProtocolAuthEnd = proto_common.Cmd_COMMON_END
)

var Handlers map[int32]func(*types.Session, *packet.RawPacket) []byte

func init() {
	Handlers = map[int32]func(*types.Session, *packet.RawPacket) []byte{
		0:  ProcHeartBeatReq,
		10: ProcUserLoginReq,
		30: ProcGetSeedReq,
	}
}
