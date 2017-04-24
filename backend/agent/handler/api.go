package handler

import (
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
)

const (
	ProtocolAuthStart = 0
	ProtocolAuthEnd   = 2000
)

var Code = map[string]int16{
	"heart_beat_req":         0,    // heart beat request
	"heart_beat_rsp":         1,    // heart beat response
	"user_login_req":         10,   // user login request
	"user_login_succeed_rsp": 11,   // user login succeed response
	"user_login_failed_rsp":  13,   // user login failed response
	"client_error_rsp":       15,   // user login error
	"get_seed_req":           30,   // socket seed request
	"get_seed_rsp":           31,   // socket seed response
	"proto_ping_req":         1000, // ping request
	"proto_ping_rsp":         1001, // ping response

}

var RCode = map[int16]string{
	0:    "heart_beat_req",
	1:    "heart_beat_rsp",
	10:   "user_login_req",
	11:   "user_login_succeed_rsp",
	13:   "user_login_failed_rsp",
	15:   "client_error_rsp",
	30:   "get_seed_req",
	31:   "get_seed_rsp",
	1000: "proto_ping_req",
	1001: "proto_ping_rsp",
}

var Handlers map[int16]func(*types.Session, *packet.RawPacket) []byte

func init() {
	Handlers = map[int16]func(*types.Session, *packet.RawPacket) []byte{
		0:  ProcHeartBeatReq,
		10: ProcUserLoginReq,
		30: ProcGetSeedReq,
	}
}
