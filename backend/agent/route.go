package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/api"
	"github.com/master-g/omgo/kit/packet"
	"github.com/master-g/omgo/kit/utils"
	proto_common "github.com/master-g/omgo/proto/pb/common"
)

func getPacketBody(reader *packet.RawPacket) []byte {
	return reader.Data()[reader.Pos():]
}

// route client protocol
func route(session *api.Session, p []byte) []byte {
	start := time.Now()
	defer utils.PrintPanicStack(session, p)
	// decrypt
	if session.IsFlagEncryptedSet() {
		session.Decoder.XORKeyStream(p, p)
	}
	// packet reader
	reader := packet.NewRawPacketReader(p)

	// read cmd
	cmdValue, err := reader.ReadS32()
	if err != nil {
		log.Errorf("read packet cmd failed:%v", err)
		session.SetFlagKicked()
		return nil
	}
	cmd := proto_common.Cmd(cmdValue)

	// route message to different service by command code
	var ret []byte
	if cmd > proto_common.Cmd_CMD_COMMON_END {
		// messages forward to game server
		if err := forward(session, getPacketBody(reader)); err != nil {
			log.Errorf("service id:%v execute failed, error:%v", cmd, err)
			session.SetFlagKicked()
			return nil
		}
	} else { // messages handle by agent service
		if h := api.Handlers[cmdValue]; h != nil {
			ret = h(session, reader)
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
