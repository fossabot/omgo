package api

import (
	log "github.com/Sirupsen/logrus"
	pc "github.com/master-g/omgo/proto/pb/common"
)

// ProcHeartBeatReq process client heartbeat packet
func ProcHeartBeatReq(session *Session, inPacket []byte) []byte {
	if !session.IsFlagAuthSet() {
		log.Errorf("heartbeat from unauthenticated session:%v", session)
		session.SetFlagKicked()
		return nil
	}

	hdr := genRspHeader(pc.Cmd_HEART_BEAT_RSP)
	return makeResponse(session, hdr, nil)
}
