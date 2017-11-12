package api

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	pc "github.com/master-g/omgo/proto/pb/common"
)

// generate a common.RspHeader
func genRspHeader(cmd pc.Cmd) *pc.RspHeader {
	header := &pc.RspHeader{
		Cmd:       int32(cmd),
		Status:    int32(pc.ResultCode_RESULT_OK),
		Timestamp: uint64(time.Now().Unix()),
	}

	return header
}

// unpack incoming protobuf packet
func unpackPacket(pkg []byte, message proto.Message) (*pc.Header, error) {
	header := &pc.Header{}
	err := proto.Unmarshal(pkg, header)
	if err != nil {
		return nil, err
	}

	if message != nil {
		err = proto.Unmarshal(header.Body, message)
		if err != nil {
			return nil, err
		}
	}

	return header, nil
}

func makeResponse(sess *Session, header *pc.RspHeader, msg proto.Message) []byte {
	sess.serverSeq++
	header.Seq = sess.serverSeq
	if msg != nil {
		body, err := proto.Marshal(msg)
		if err != nil {
			logrus.Errorf("error while marshaling message error:%v", err)
			return nil
		}

		header.Body = body
	}

	ret, err := proto.Marshal(header)
	if err != nil {
		logrus.Errorf("error while marshaling message error:%v", err)
		return nil
	}

	return ret
}
