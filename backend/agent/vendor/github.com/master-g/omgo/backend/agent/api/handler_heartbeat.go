package api

import (
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/kit/packet"
	pc "github.com/master-g/omgo/proto/pb/common"
)

// ProcHeartBeatReq process client heartbeat packet
func ProcHeartBeatReq(session *Session, inPacket *IncomingPacket) []byte {
	if !session.IsFlagAuthSet() {
		log.Errorf("heartbeat from unauthenticated session:%v", session)
		session.SetFlagKicked()
		return nil
	}

	p := packet.NewRawPacket()
	p.WriteS32(int32(pc.Cmd_HEART_BEAT_RSP))
	return p.Data()
}
