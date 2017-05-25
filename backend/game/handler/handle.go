package handler

import (
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/backend/game/types"
	"github.com/master-g/omgo/net/packet"
	proto_common "github.com/master-g/omgo/proto/pb/common"
)

func response(cmd proto_common.Cmd, msg proto.Message) []byte {
	p := packet.NewRawPacket()
	p.WriteS32(int32(cmd))
	rspBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}
	p.WriteBytes(rspBytes)
	return p.Data()
}

func genRspHeader() *proto_common.RspHeader {
	header := &proto_common.RspHeader{
		Status:    proto_common.ResultCode_RESULT_OK,
		Timestamp: uint64(time.Now().Unix()),
	}

	return header
}

func ProcPingReq(session *types.Session, reader *packet.RawPacket) []byte {
	p := packet.NewRawPacket()
	p.WriteS32(int32(proto_common.Cmd_PING_RSP))
	return p.Data()
}
