package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/api"
	"github.com/master-g/omgo/kit/utils"
	proto_common "github.com/master-g/omgo/proto/pb/common"
)

// route client protocol
func route(session *api.Session, inPacket *api.IncomingPacket) *api.OutgoingPacket {
	start := time.Now()
	defer utils.PrintPanicStack(session, inPacket)
	// decrypt
	if session.IsFlagEncryptedSet() {
		session.Decoder.XORKeyStream(inPacket.Body, inPacket.Body)
	}

	// read cmd
	cmdValue := inPacket.Header.Cmd
	cmd := proto_common.Cmd(cmdValue)

	// route message to different service by command code
	var ret *api.OutgoingPacket
	if cmd > proto_common.Cmd_CMD_COMMON_END {
		// messages forward to other servers
		if err := forward(session, inPacket); err != nil {
			log.Errorf("service id:%v execute failed, error:%v", cmd, err)
			session.SetFlagKicked()
			return nil
		}
	} else { // messages handle by agent service
		if h := api.Handlers[cmdValue]; h != nil {
			ret = h(session, inPacket)
		} else {
			log.Errorf("no handler for cmd:%v", cmd)
			session.SetFlagKicked()
			return nil
		}
	}

	elapsed := time.Now().Sub(start)
	//if cmd != proto_common.Cmd_HEART_BEAT_REQ {
	log.WithFields(log.Fields{
		"cost": elapsed,
		"code": cmd,
	}).Debug("REQ")
	//}

	return ret
}
