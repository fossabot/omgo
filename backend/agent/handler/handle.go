package handler

import (
	"context"
	"crypto/rand"
	"crypto/rc4"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
	pb "github.com/master-g/omgo/proto/grpc/game"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/security/ecdh"
	"github.com/master-g/omgo/services"
	"google.golang.org/grpc/metadata"
)

// convert proto message into packet
func response(cmd pc.Cmd, msg proto.Message) []byte {
	p := packet.NewRawPacket()
	p.WriteS32(int32(cmd))
	rspBytes, err := proto.Marshal(msg)
	if err != nil {
		return nil
	}
	p.WriteBytes(rspBytes)
	return p.Data()
}

// generate a common.RspHeader
func genRspHeader() *pc.RspHeader {
	header := &pc.RspHeader{
		Status:    pc.ResultCode_RESULT_OK,
		Timestamp: uint64(time.Now().Unix()),
	}

	return header
}

// ProcHeartBeatReq process client heartbeat packet
// TODO: reset client timeout timer
func ProcHeartBeatReq(session *types.Session, reader *packet.RawPacket) []byte {
	p := packet.NewRawPacket()
	p.WriteS32(int32(pc.Cmd_HEART_BEAT_RSP))
	return p.Data()
}

// ProcGetSeedReq exchange secret with client via ECDH algorithm
// TODO: optimize performance
func ProcGetSeedReq(session *types.Session, reader *packet.RawPacket) []byte {
	rsp := &pc.S2CGetSeedRsp{Header: genRspHeader()}
	req := &pc.C2SGetSeedReq{}
	marshalPb, _ := reader.ReadBytes()

	if err := proto.Unmarshal(marshalPb, req); err != nil {
		log.Errorf("invalid protobuf :%v", err)
		rsp.Header.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		return response(pc.Cmd_GET_SEED_RSP, rsp)
	}

	curve := ecdh.NewCurve25519ECDH()
	x1, e1 := curve.GenerateECKeyBuf(rand.Reader)
	key1 := curve.GenerateSharedSecretBuf(x1, req.GetSendSeed())
	x2, e2 := curve.GenerateECKeyBuf(rand.Reader)
	key2 := curve.GenerateSharedSecretBuf(x2, req.GetRecvSeed())

	encoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key2)))
	if err != nil {
		log.Error(err)
		rsp.Header.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		return response(pc.Cmd_GET_SEED_RSP, rsp)
	}
	decoder, err := rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key1)))
	if err != nil {
		log.Error(err)
		rsp.Header.Status = pc.ResultCode_RESULT_INTERNAL_ERROR
		return response(pc.Cmd_GET_SEED_RSP, rsp)
	}
	session.Encoder = encoder
	session.Decoder = decoder
	session.SetFlagKeyExchanged()

	rsp.SendSeed = e1
	rsp.RecvSeed = e2

	return response(pc.Cmd_GET_SEED_RSP, rsp)
}

// ProcUserLoginReq process user login request
func ProcUserLoginReq(session *types.Session, reader *packet.RawPacket) []byte {
	if !session.IsFlagEncryptedSet() {
		log.Errorf("session login without encryption:%v", session)
		session.SetFlagKicked()
		return nil
	}

	req := &pc.C2SLoginReq{}
	marshalPb, _ := reader.ReadBytes()

	if err := proto.Unmarshal(marshalPb, req); err != nil {
		log.Errorf("invalid protobuf :%v", err)
		session.SetFlagKicked()
		return nil
	}

	session.Usn = 1
	session.GSID = DefaultGSID

	conn := services.GetServiceWithID("game", session.GSID)
	if conn == nil {
		log.Error("cannot get game service:", session.GSID)
		return nil
	}
	cli := pb.NewGameServiceClient(conn)

	ctx := metadata.NewContext(context.Background(), metadata.New(map[string]string{"usn": fmt.Sprint(session.Usn)}))
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
	p.WriteS32(int32(pc.Cmd_LOGIN_RSP))
	p.WriteU64(session.Usn)

	return p.Data()
}

func ProcOfflineReq(session *types.Session, reader *packet.RawPacket) []byte {
	session.SetFlagKicked()
	return nil
}
