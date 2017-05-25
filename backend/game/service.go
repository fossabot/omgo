package main

import (
	"errors"
	"io"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/game/handler"
	"github.com/master-g/omgo/backend/game/registry"
	"github.com/master-g/omgo/backend/game/types"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/proto/grpc/game"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc/metadata"
)

var (
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
	var session types.Session
	chSessionDie := make(chan struct{})
	chAgent := s.recv(stream, chSessionDie)
	chIPC := make(chan *proto.Game_Frame, DefaultIPCChannelSize)

	defer func() {
		registry.Unregister(session.UserID, chIPC)
		close(chSessionDie)
		log.Debug("stream end:", session.UserID)
	}()

	// read metadata from context
	md, ok := metadata.FromContext(stream.Context())
	if !ok {
		log.Error("cannot read metadata from context")
		return ErrorIncorrectFrameType
	}
	// read key
	if len(md["userid"]) == 0 {
		log.Error("cannot read key:userid from metadata")
		return ErrorIncorrectFrameType
	}
	// parse userID
	userID, err := strconv.Atoi(md["userid"][0])
	if err != nil {
		log.Error(err)
		return ErrorIncorrectFrameType
	}

	// register user
	session.UserID = int32(userID)
	registry.Register(session.UserID, chIPC)
	log.Debug("userid", session.UserID, "logged in")

	// *** main message loop ***
	for {
		select {
		case frame, ok := <-chAgent:
			// frame from agent
			if !ok {
				return nil
			}
			switch frame.Type {
			case proto.Game_Message:
				// passthrough message from client->agent
				reader := packet.NewRawPacketReader(frame.Message)
				c, err := reader.ReadS32()
				if err != nil {
					log.Error(err)
					return err
				}
				// handle request
				h := handler.Handlers[c]
				if h == nil {
					log.Error("service not bound for:", c)
					return ErrorServiceNotBound
				}
				ret := h(&session, reader)

				// construct frame and return message from logic
				if ret != nil {
					if err := stream.Send(&proto.Game_Frame{Type: proto.Game_Message, Message: ret}); err != nil {
						log.Error(err)
						return err
					}
				}

				// session control by logic
				if session.IsFlagKickedSet() {
					// logic kick out
					if err := stream.Send(&proto.Game_Frame{Type: proto.Game_Kick}); err != nil {
						log.Error(err)
						return err
					}
					return nil
				}
			case proto.Game_Ping:
				if err := stream.Send(&proto.Game_Frame{Type: proto.Game_Ping, Message: frame.Message}); err != nil {
					log.Error(err)
					return err
				}
				log.Debug("ping respond")
			default:
				log.Error("incorrect frame type:", frame.Type)
				return ErrorIncorrectFrameType
			}
		case frame := <-chIPC:
			// forward async messages from interprocess(goroutines) communication
			if err := stream.Send(frame); err != nil {
				log.Error(err)
				return err
			}
		}
	}
}
