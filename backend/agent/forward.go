package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/backend/agent/proto"
	"github.com/master-g/omgo/backend/agent/types"
)

var (
	ErrorStreamNotOpen = errors.New("stream not opened yet")
)

// forward messages to game server
func forward(session *types.Session, p []byte) error {
	frame := &pb.Game_Frame{
		Type:    pb.Game_Message,
		Message: p,
	}

	// check stream
	if session.Stream == nil {
		return ErrorStreamNotOpen
	}

	// forward the frame to game
	if err := session.Stream.Send(frame); err != nil {
		log.Error(err)
		return err
	}

	return nil
}
