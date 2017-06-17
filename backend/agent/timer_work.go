package main

import (
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/types"
)

var (
	rpmLimit int
)

func initTimer(limit int) {
	rpmLimit = limit
}

func timerWork(session *types.Session, out *Buffer) {
	defer func() {
		session.PacketCountPerMin = 0
	}()

	// rpm control
	if session.PacketCountPerMin > rpmLimit {
		session.SetFlagKicked()
		log.WithFields(log.Fields{
			"usn":   session.Usn,
			"rate":  session.PacketCountPerMin,
			"total": session.PacketCount,
		}).Error("RPM")
		return
	}

	// heartbeat check
	elapsed := time.Since(session.LastPacketTime)
	if time.Minute > elapsed {
		session.SetFlagKicked()
		log.WithFields(log.Fields{
			"usn":         session.Usn,
			"lastpkgtime": session.LastPacketTime,
		}).Error("TIMEOUT")
		return
	}
}
