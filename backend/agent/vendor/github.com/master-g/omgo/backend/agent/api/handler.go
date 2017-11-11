package api

import (
	"time"

	pc "github.com/master-g/omgo/proto/pb/common"
)

// generate a common.RspHeader
func genRspHeader() *pc.RspHeader {
	header := &pc.RspHeader{
		Status:    int32(pc.ResultCode_RESULT_OK),
		Timestamp: uint64(time.Now().Unix()),
	}

	return header
}
