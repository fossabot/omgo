package main

import (
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/proto/pb"
	"github.com/master-g/omgo/utils"
	"net"
	"testing"
)

var address string

func init() {
	address = utils.GetLocalIP() + ":8888"
}

func connect(t *testing.T) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("could not connect to server %v, error: %v", address, err)
	}

	return conn
}

func TestConnect(t *testing.T) {
	// Setup a connection to the agent server.
	conn := connect(t)
	defer conn.Close()

}

func TestHeartBeat(t *testing.T) {
	conn := connect(t)
	defer conn.Close()

	reqPacket := packet.NewRawPacket()
	reqPacket.WriteS32(int32(proto_common.Cmd_HEART_BEAT_REQ))
}
