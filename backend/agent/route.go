package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/handler"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
	proto_common "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/utils"
)

func getPacketBody(reader *packet.RawPacket) []byte {
	return reader.Data()[reader.Pos():]
}

// route client protocol
func route(session *types.Session, p []byte) []byte {
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
		log.Error("read packet cmd failed:", err)
		session.SetFlagKicked()
		return nil
	}
	cmd := proto_common.Cmd(cmdValue)

	// route message to different service by command code
	var ret []byte
	if cmd > proto_common.Cmd_CMD_COMMON_END {
		if err := forward(session, getPacketBody(reader)); err != nil {
			log.Errorf("service id:%v execute failed, error:%v", cmd, err)
			session.SetFlagKicked()
			return nil
		}
	} else {
		if h := handler.Handlers[cmdValue]; h != nil {
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
