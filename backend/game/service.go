package main

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/proto"
	"github.com/master-g/omgo/backend/game/registry"
	"github.com/master-g/omgo/backend/game/types"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/utils"
	"google.golang.org/grpc/metadata"
	"io"
	"strconv"
	"github.com/master-g/omgo/backend/game/handler"
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
	sess.UserID = int32(userID)
	registry.Register(sess.UserID, chIPC)
	log.Debug("userid", sess.UserID, "logged in")

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
				ret := h(&sess, reader)

				// construct frame and return message from logic
				if ret != nil {
					if err := stream.Send(&proto.Game_Frame{Type:proto.Game_Message, Message:ret}); err != nil {
						log.Error(err)
						return err
					}
				}

				// session control by logic
				if sess.Flag & types.FlagKicked != 0 {
					// logic kick out
					if err := stream.Send(&proto.Game_Frame{Type:proto.Game_Kick})
				}
			}
		}
	}
}
