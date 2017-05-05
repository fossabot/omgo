package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/proto"
	"github.com/master-g/omgo/backend/game/registry"
	"github.com/master-g/omgo/backend/game/types"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc/metadata"
	"io"
)

const (
	DefaultIPCChannelSize   = 16
	ErrorIncorrectFrameType = errors.New("incorrect frame type")
	ErrorServiceNotBound    = errors.New("service not bound")
)

type server struct{}

// PIPELINE #1 stream receiver
// this function is to make the stream receiving SELECTABLE
func (s *server) recv(stream proto.GameService_StreamServer, chSessDie chan struct{}) chan *proto.Game_Frame {
	ch := make(chan *proto.Game_Frame, 1)
	go func() {
		defer func() {
			close(ch)
		}()

		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// client closed
				return
			}

			if err != nil {
				log.Error(err)
				return
			}

			select {
			case ch <- in:
			case <-chSessDie:
			}
		}
	}()
	return ch
}

// PIPELINE #2 stream processing
// the center of game logic
func (s *server) Stream(stream proto.GameService_StreamServer) error {
	defer utils.PrintPanicStack()
	// session init
	var sess types.Session
	chSessDie := make(chan struct{})
	chAgent := s.recv(stream, chSessDie)
	chIPC := make(chan *proto.Game_Frame, DefaultIPCChannelSize)

	defer func() {
		registry.Unregister(sess.UserID, chIPC)
		close(chSessDie)
		log.Debug("stream end:", sess.UserID)
	}()

	// read metadata from context
	md, ok := metadata.FromContext(stream.Context())
	if !ok {
		log.Error("cannot read metadata from context")
		return ErrorIncorrectFrameType
	}
}
