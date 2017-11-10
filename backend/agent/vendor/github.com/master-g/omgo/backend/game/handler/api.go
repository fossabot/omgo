package handler

import (
	"github.com/master-g/omgo/backend/game/types"
	"github.com/master-g/omgo/net/packet"
	pc "github.com/master-g/omgo/proto/pb/common"
)

var Handlers map[int32]func(session *types.Session, packet *packet.RawPacket) []byte

func init() {
	Handlers = map[int32]func(session *types.Session, rawPacket *packet.RawPacket) []byte{
		int32(pc.Cmd_PING_REQ): ProcPingReq,
	}
}
