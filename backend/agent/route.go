package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/handler"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/proto/pb"
	"github.com/master-g/omgo/utils"
	"time"
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
	if cmd < proto_common.Cmd_CMD_COMMON_END {
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
	if cmd != 0 {
		log.WithFields(log.Fields{
			"cost": elapsed,
			"api":  proto_common.Cmd_name[cmdValue],
			"code": cmd,
		}).Debug("REQ")
	}

	return ret
}
