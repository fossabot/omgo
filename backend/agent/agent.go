package main

import (
	pb "github.com/master-g/omgo/backend/agent/proto"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/utils"
	"time"
)

const (
	defaultMQSize = 512
)

// PIPELINE #2: agent
// all the packets from handleClient() will be handled here
func agent(session *types.Session, in chan []byte, out *Buffer) {
	defer wg.Done() // will decrease waitgroup by one, useful for manual server shutdown
	defer utils.PrintPanicStack()

	// init session
	session.MQ = make(chan pb.Game_Frame, defaultMQSize)
	session.ConnectTime = time.Now()
	session.LastPacketTime = time.Now()
	// minute timer
	minTimer := time.After(time.Minute)

	// cleanup
	defer func() {
		close(session.Die)
		if session.Stream != nil {
			session.Stream.CloseSend()
		}
	}()

	// **** MAIN MESSAGE LOOP ****
	// handles 4 types of message:
	//  1. from client
	//  2. from game service
	//  3. timer
	//  4. server shutdown signal
	for {
		select {
		case msg, ok := <-in: // packet from network
			if !ok {
				return
			}

			session.PacketCount++
			session.PacketCountPerMin++
			session.PacketTime = time.Now()

			if result := route(session, msg); result != nil {
				out.send(session, result)
			}
			session.LastPacketTime = session.PacketTime
		case frame := <-session.MQ: // packets from frame
			switch frame.Type {
			case pb.Game_Message:
				out.send(session, frame.Message)
			case pb.Game_Kick:
				session.SetFlagKicked()
			}
		case <-minTimer: // minutes timer
			timerWork(session, out)
			minTimer = time.After(time.Minute)
		case <-die: // server is shutting down
			session.SetFlagKicked()
		}

		// see if user should be kicked out
		if session.IsFlagKickedSet() {
			return
		}
	}
}
