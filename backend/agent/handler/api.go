package handler

import (
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
)

const (
	ProtocolAuthEnd = proto_common.Cmd_
)

var Handlers map[int16]func(*types.Session, *packet.RawPacket) []byte

func init() {
	Handlers = map[int16]func(*types.Session, *packet.RawPacket) []byte{
		0:  ProcHeartBeatReq,
		10: ProcUserLoginReq,
		30: ProcGetSeedReq,
	}
}
