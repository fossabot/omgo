package main

import (
	"errors"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/api"
	pb "github.com/master-g/omgo/proto/grpc/game"
)

var (
	// ErrorStreamNotOpen indicates error while opening gRPC stream
	ErrorStreamNotOpen = errors.New("stream not opened yet")
)

// forward messages to game server
func forward(session *api.Session, msg []byte) error {
	frame := &pb.Game_Frame{
		Type:    pb.Game_Message,
		Message: msg,
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
