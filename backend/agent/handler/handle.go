package handler

import (
	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/backend/agent/proto"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
)

func ProcHeartBeatReq(session *types.Session, reader *packet.RawPacket) []byte {
	tbl, _ := PacketReadAutoID(reader)
	p := packet.NewRawPacket()
	p.WriteS32(tbl.ID)
	return p.Data()
}

func ProcGetSeedReq(session *types.Session, reader *packet.RawPacket) []byte {
    tbl, _ := PacketReadSeedInfo(reader)
    x1, e1 := ecdh
}
