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

	// read client packet sequence number
	// every time client sends a packet, its sequence number must strictly increase by one
	seqNumber, err := reader.ReadU32()
	if err != nil {
		log.Error("read client timestamp failed:", err)
		session.SetFlagKicked()
		return nil
	}

	// sequence number verification
	if seqNumber != session.PacketCount {
		log.Errorf("illegal packet sequence id:%v should be %v size:%v", seqNumber, session.PacketCount, len(p)-6)
		session.SetFlagKicked()
		return nil
	}

	// read protocol number
	cmdValue, err := reader.ReadS32()
	if err != nil {
		log.Error("read protocol number failed.")
		session.SetFlagKicked()
		return nil
	}

	cmd := proto_common.Cmd(cmdValue)

	// route message to different service by protocol number
	var ret []byte
	if cmd < proto_common.Cmd_CMD_COMMON_END {
		if err := forward(session, p[4:]); err != nil {
			log.Errorf("service id:%v execute failed, error:%v", cmd, err)
			session.SetFlagKicked()
			return nil
		}
	} else {
		if h := handler.Handlers[cmdValue]; h != nil {
			ret = h(session, reader)
		} else {
			log.Errorf("service id:%v not bind", cmd)
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
