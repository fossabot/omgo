package handler

import (
	"crypto/rand"
	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/backend/agent/proto"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/security/ecdh"
	"os"
    "math/big"
)

func ProcHeartBeatReq(session *types.Session, reader *packet.RawPacket) []byte {
	tbl, _ := PacketReadAutoID(reader)
	p := packet.NewRawPacket()
	p.WriteS32(tbl.ID)
	return p.Data()
}

func ProcGetSeedReq(session *types.Session, reader *packet.RawPacket) []byte {
	tbl, _ := PacketReadSeedInfo(reader)
	curve := ecdh.NewCurve25519ECDH()
	x1, e1, err := curve.GenerateECKey(rand.Reader)
	checkErr(err)
	x2, e2, err := curve.GenerateECKey(rand.Reader)
	checkErr(err)

    key1, err := curve.gen

    return nil
}
