package main

import (
	log "github.com/Sirupsen/logrus"
)

var (
	rpmLimit int
)

func initTimer(limit int) {
	rpmLimit = limit
}

func timerWork(session *Session, out *Buffer) {
	defer func() {
		session.PacketCountPerMin = 0
	}()

	// rpm control
	if session.PacketCountPerMin > rpmLimit {
		session.Flag |= SESS_KICKED
		log.WithFields(log.Fields{
			"userid": session.UserID,
			"rate":   session.PacketCountPerMin,
			"total":  session.PacketCount,
		}).Error("RPM")
		return
	}
}
