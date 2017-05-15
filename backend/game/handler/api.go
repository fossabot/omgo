package handler

import (
	"github.com/master-g/omgo/backend/game/types"
	"github.com/master-g/omgo/net/packet"
	proto_common "github.com/master-g/omgo/proto/pb/common"
)

var Handlers map[int32]func(*types.Session, *packet.RawPacket) []byte

func init() {
	Handlers = map[int32]func(session *types.Session, rawPacket *packet.RawPacket) []byte{
		int32(proto_common.Cmd_PING_REQ): ProcPingReq,
	}
}
