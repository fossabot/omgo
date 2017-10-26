// agent gRPC service
package main

import (
	pb "github.com/master-g/omgo/proto/grpc/agent"
	"golang.org/x/net/context"
)

type server struct{}

// KickUser from agent server
func (s *server) KickUser(ctx context.Context, userEntry *pb.Agent_UserEntry) (*pb.Agent_Result, error) {
	return nil, nil
}
