package handler

import (
	"context"
	"crypto/rand"
	"crypto/rc4"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	pb "github.com/master-g/omgo/backend/agent/proto"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/proto/pb"
	"github.com/master-g/omgo/security/ecdh"
	"github.com/master-g/omgo/services"
	"google.golang.org/grpc/metadata"
	"io"
	"time"
)

func errorRsp(cmd proto_common.Cmd, msg proto.Message) {

}

func genRspHeader() *proto_common.RspHeader {
	header := &proto_common.RspHeader{
		Status:    proto_common.ResultCode_RESULT_OK,
		Timestamp: uint64(time.Now().Unix()),
	}

	return header
}

func ProcHeartBeatReq(session *types.Session, reader *packet.RawPacket) []byte {
	p := packet.NewRawPacket()
	p.WriteS32(int32(proto_common.Cmd_HEART_BEAT_RSP))
	return p.Data()
}

func ProcGetSeedReq(session *types.Session, reader *packet.RawPacket) []byte {
	rsp := &proto_common.S2CGetSeedRsp{Header: genRspHeader()}

	req := &proto_common.C2SGetSeedReq{}
	if err := proto.Unmarshal(reader.Data(), req); err != nil {
		log.Errorf("invalid protobuf :%v", err)
		rsp.Header.Status = proto_common.ResultCode_RESULT_INTERNAL_ERROR
	}

	curve := ecdh.NewCurve25519ECDH()
	x1, e1 := curve.GenerateECKeyBuf(rand.Reader)
	key1 := curve.GenerateSharedSecretBuf(x1, req.GetSendSeed())
	x2, e2 := curve.GenerateECKeyBuf(rand.Reader)
	key2 := curve.GenerateSharedSecretBuf(x2, req.GetRecvSeed())

	encoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key2)))
	if err != nil {
		log.Error(err)
		return nil
	}
	decoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key1)))
	if err != nil {
		log.Error(err)
		return nil
	}
	session.Encoder = encoder
	session.Decoder = decoder
	session.SetFlagKeyExchanged()

	rsp.SendSeed = e1
	rsp.RecvSeed = e2
	rspBytes, err := proto.Marshal(rsp)
	if err != nil {
		log.Errorf("failed to marshal response, error:%v", err)
		return nil
	}

Response:
	p := packet.NewRawPacket()
	p.WriteS32(int32(proto_common.Cmd_GET_SEED_RSP))
	p.WriteBytes(rspBytes)

	return p.Data()
}

func ProcUserLoginReq(session *types.Session, reader *packet.RawPacket) []byte {
	session.UserID = 1
	session.GSID = DefaultGSID

	conn := services.GetServiceWithID("game-1000", session.GSID)
	if conn == nil {
		log.Error("cannot get game service:", session.GSID)
		return nil
	}
	cli := pb.NewGameServiceClient(conn)

	ctx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"userid": fmt.Sprint(session.UserID)}))
	stream, err := cli.Stream(ctx)
	if err != nil {
		log.Error(err)
		return nil
	}
	session.Stream = stream

	fetcherTask := func(session *types.Session) {
		for {
			in, err := session.Stream.Recv()
			if err == io.EOF {
				log.Debug(err)
				return
			}
			if err != nil {
				log.Error(err)
				return
			}
			select {
			case session.MQ <- *in:
			case <-session.Die:
			}
		}
	}
	go fetcherTask(session)
	p := packet.NewRawPacket()
	p.WriteS32(int32(proto_common.Cmd_LOGIN_RSP))
	p.WriteU32(session.UserID)

	return p.Data()
}
