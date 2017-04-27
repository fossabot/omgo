package main

import (
	"encoding/binary"
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

func send(conn net.Conn, data []byte, t *testing.T) {
	size := len(data)
	cache := make([]byte, 2+size)
	binary.BigEndian.PutUint16(cache, uint16(size))
	copy(cache[2:], data)
	_, err := conn.Write(cache[:size+2])
	if err != nil {
		t.Fatalf("error while sending data: %v", err)
	}
}

func TestHeartBeat(t *testing.T) {
	conn := connect(t)
	defer conn.Close()

	reqPacket := packet.NewRawPacket()
	reqPacket.WriteS32(int32(proto_common.Cmd_HEART_BEAT_REQ))
	send(conn, reqPacket.Data(), t)
}
