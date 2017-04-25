package handler

import (
	"context"
	"crypto/rand"
	"crypto/rc4"
	"fmt"
	log "github.com/Sirupsen/logrus"
	pb "github.com/master-g/omgo/backend/agent/proto"
	"github.com/master-g/omgo/backend/agent/types"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/security/ecdh"
	"github.com/master-g/omgo/services"
	"google.golang.org/grpc/metadata"
	"io"
)

func ProcHeartBeatReq(session *types.Session, reader *packet.RawPacket) []byte {
	p := packet.NewRawPacket()
	p.WriteS16(Code["heart_beat_rsp"])
	p.WriteU32(tbl.ID)
	return p.Data()
}

func ProcGetSeedReq(session *types.Session, reader *packet.RawPacket) []byte {
	tbl, _ := PacketReadSeedInfo(reader)
	curve := ecdh.NewCurve25519ECDH()
	x1, e1 := curve.GenerateECKeyBuf(rand.Reader)
	key1 := curve.GenerateSharedSecretBuf(x1, tbl.ClientSendSeed)
	x2, e2 := curve.GenerateECKeyBuf(rand.Reader)
	key2 := curve.GenerateSharedSecretBuf(x2, tbl.ClientRecvSeed)

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

	p := packet.NewRawPacket()
	p.WriteS16(Code["get_seed_rsp"])
	p.WriteBytes(e1)
	p.WriteBytes(e2)

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
	p.WriteS16(Code["user_login_succeed_rsp"])
	p.WriteU32(session.UserID)

	return p.Data()
}

func ProcPingReq(session *types.Session, reader *packet.RawPacket) []byte {
	return nil
}
