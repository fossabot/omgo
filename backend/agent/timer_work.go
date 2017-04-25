package main

import (
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
			"userid": session.UserID,
			"rate":   session.PacketCountPerMin,
			"total":  session.PacketCount,
		}).Error("RPM")
		return
	}
}
