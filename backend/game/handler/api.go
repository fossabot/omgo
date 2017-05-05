package handler

import (
	"github.com/master-g/omgo/backend/game/types"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/proto/pb"
)

var Handlers map[int32]func(*types.Session, *packet.RawPacket) []byte

func init() {
	Handlers = map[int32]func(session *types.Session, rawPacket *packet.RawPacket) []byte{
		int32(proto_common.Cmd_PING_REQ): ProcPingReq,
	}
}
