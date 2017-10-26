// agent gRPC service
package main

import (
	"github.com/master-g/omgo/backend/agent/api"
	"github.com/master-g/omgo/kit/utils"
	pb "github.com/master-g/omgo/proto/grpc/agent"
	pc "github.com/master-g/omgo/proto/pb/common"
	"golang.org/x/net/context"
)

type server struct{}

// KickUser from agent server
func (s *server) KickUser(ctx context.Context, userEntry *pb.Agent_UserEntry) (*pb.Agent_Result, error) {
	value, ok := api.Registry.Load(userEntry.Usn)
	if ok {
		if session, ok := value.(*api.Session); ok {
			kickNotify := &pc.S2CKickNotify{
				Timestamp: utils.Timestamp(),
				Reason:    pc.KickReason_KICK_NO_REASON,
				Msg:       session.IP.String(),
			}
			session.Mailbox <- api.MakeResponse(pc.Cmd_KICK_NOTIFY, kickNotify)
			session.SetFlagKicked()
		}
		return &pb.Agent_Result{Status: pb.Agent_STATUS_OK}, nil
	} else {
		return &pb.Agent_Result{Status: pb.Agent_STATUS_NOT_FOUND, Msg: "not found"}, nil
	}
}
