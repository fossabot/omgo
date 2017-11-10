package main

import (
	"github.com/master-g/omgo/kit/packet"
	pc "github.com/master-g/omgo/proto/pb/common"
)

func makePacket(cmd pc.Cmd) *packet.RawPacket {
	pk := packet.NewRawPacket()
	pk.WriteS32(int32(cmd))
	return pk
}
