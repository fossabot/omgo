package main

import (
	"crypto/rc4"
	pb "github.com/master-g/omgo/backend/agent/proto"
	"net"
	"time"
)

const (
	SESS_KEYEXCG = 0x1 // if key is exchanged
	SESS_ENCRYPT = 0x2 // if encryption is available
	SESS_KICKED  = 0x4 // kick out
	SESS_AUTHED  = 0x8 // if authorized
)

type Session struct {
	IP      net.IP                      // client IP address
	MQ      chan pb.Game_Frame          // channel of async messages send back to client
	Encoder *rc4.Cipher                 // encrypt
	Decoder *rc4.Cipher                 // decrypt
	UserID  int32                       // user ID
	GSID    string                      // game server ID
	Stream  pb.GameService_StreamClient // data stream send to game server
	Die     chan struct{}               // session close signal

	Flag int32 // session flag

	ConnectTime    time.Time // timestamp of TCP connection established
	PacketTime     time.Time // timestamp of current packet arrived
	LastPacketTime time.Time // timestamp of previous packet arrived

	PacketCount       uint32 // total packets received
	PacketCountPerMin int    // packets received per minute
}
