package api

import (
	"crypto/rc4"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/master-g/omgo/kit/ecdh"
	"google.golang.org/grpc/metadata"

	"crypto/rand"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/kit/packet"
	"github.com/master-g/omgo/kit/utils"
	pbdb "github.com/master-g/omgo/proto/grpc/db"
	pbgame "github.com/master-g/omgo/proto/grpc/game"
	pc "github.com/master-g/omgo/proto/pb/common"
	"golang.org/x/net/context"
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
		Status:    int32(pc.ResultCode_RESULT_OK),
		Timestamp: uint64(time.Now().Unix()),
	}

	return header
}

// ProcHeartBeatReq process client heartbeat packet
// TODO: reset client timeout timer
func ProcHeartBeatReq(session *Session, reader *packet.RawPacket) []byte {
	if !session.IsFlagAuthSet() {
		log.Errorf("heartbeat from unauth session:%v", session)
		session.SetFlagKicked()
		return nil
	}

	p := packet.NewRawPacket()
	p.WriteS32(int32(pc.Cmd_HEART_BEAT_RSP))
	return p.Data()
}

// ProcGetSeedReq exchange secret with client via ECDH algorithm
// TODO: optimize performance
func ProcGetSeedReq(session *Session, reader *packet.RawPacket) []byte {
	rsp := &pc.S2CGetSeedRsp{Header: genRspHeader()}
	req := &pc.C2SGetSeedReq{}
	marshalPb, _ := reader.ReadBytes()

	if err := proto.Unmarshal(marshalPb, req); err != nil {
		log.Errorf("invalid protobuf :%v", err)
		rsp.Header.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		return response(pc.Cmd_GET_SEED_RSP, rsp)
	}

	curve := ecdh.NewCurve25519ECDH()
	x1, e1 := curve.GenerateECKeyBuf(rand.Reader)
	key1 := curve.GenerateSharedSecretBuf(x1, req.GetSendSeed())
	x2, e2 := curve.GenerateECKeyBuf(rand.Reader)
	key2 := curve.GenerateSharedSecretBuf(x2, req.GetRecvSeed())

	encoder, err := rc4.NewCipher(key2)
	if err != nil {
		log.Error(err)
		rsp.Header.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		return response(pc.Cmd_GET_SEED_RSP, rsp)
	}
	decoder, err := rc4.NewCipher(key1)
	if err != nil {
		log.Error(err)
		rsp.Header.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
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
func ProcUserLoginReq(session *Session, reader *packet.RawPacket) []byte {
	rsp := &pc.S2CLoginRsp{Header: genRspHeader()}
	rsp.Header.Timestamp = utils.Timestamp()
	rsp.Header.Status = int32(pc.ResultCode_RESULT_INVALID)

	// can only login after key exchange
	if !session.IsFlagEncryptedSet() {
		log.Errorf("session login without encryption:%v", session)
		session.SetFlagKicked()
		return response(pc.Cmd_LOGIN_RSP, rsp)
	}

	// parse login request
	req := &pc.C2SLoginReq{}
	marshalPb, _ := reader.ReadBytes()

	if err := proto.Unmarshal(marshalPb, req); err != nil {
		log.Errorf("invalid protobuf:%v", err)
		session.SetFlagKicked()
		return response(pc.Cmd_LOGIN_RSP, rsp)
	}

	usn := req.GetUsn()
	token := req.GetToken()

	if usn == 0 || token == "" {
		log.Errorf("invalid usn:%v or token:%v", usn, token)
		session.SetFlagKicked()
		return response(pc.Cmd_LOGIN_RSP, rsp)
	}

	// validate user token
	dbConn := DataServicePool.NextClient()
	if dbConn == nil {
		log.Errorf("dataservice not connected yet")
		return response(pc.Cmd_LOGIN_RSP, rsp)
	}

	dbClient := pbdb.NewDBServiceClient(dbConn)
	userKey := &pbdb.DB_UserEntry{Usn: usn}
	dbRsp, err := dbClient.UserExtraInfoQuery(context.Background(), userKey)
	if err != nil {
		log.Errorf("error while query user extra info:%v", usn)
		rsp.Header.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		return response(pc.Cmd_LOGIN_RSP, rsp)
	}

	if dbRsp.Result.Status != pbdb.DB_STATUS_OK || dbRsp.User.Token == "" {
		log.Errorf("user extra info not found:%v", usn)
		rsp.Header.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		return response(pc.Cmd_LOGIN_RSP, rsp)
	}

	if strings.Compare(token, dbRsp.User.Token) != 0 {
		log.Infof("invalid token")
		session.SetFlagKicked()
		return response(pc.Cmd_LOGIN_RSP, rsp)
	}

	// kick previous session if existed
	p, ok := Registry.Load(usn)
	if ok {
		if prevSession, ok := p.(*Session); ok {
			kickNotify := &pc.S2CKickNotify{
				Timestamp: utils.Timestamp(),
				Reason:    pc.KickReason_KICK_LOGIN_ELSEWHERE,
				Msg:       session.IP.String(),
			}
			prevSession.Mailbox <- response(pc.Cmd_KICK_NOTIFY, kickNotify)
			prevSession.SetFlagKicked()
		}
	}

	Registry.Store(usn, session)

	// connection to game server
	session.Usn = usn
	session.Token = token
	session.GSID = config.GameServerName
	session.SetFlagAuth()

	conn := GameServerPool.GetClient(config.GameServerName)
	if conn == nil {
		log.Error("cannot get game service:", session.GSID)
		rsp.Header.Status = int32(pc.ResultCode_RESULT_INTERNAL_ERROR)
		return response(pc.Cmd_LOGIN_RSP, rsp)
	}
	cli := pbgame.NewGameServiceClient(conn)

	// open game server stream
	md := metadata.New(map[string]string{"usn": fmt.Sprint(session.Usn)})
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	stream, err := cli.Stream(ctx)
	if err != nil {
		log.Error(err)
		return nil
	}
	session.Stream = stream

	// read message returned by game server
	fetcherTask := func(session *Session) {
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

	// login success
	rsp.Header.Status = int32(pc.ResultCode_RESULT_OK)
	return response(pc.Cmd_LOGIN_RSP, rsp)
}

func ProcOfflineReq(session *Session, reader *packet.RawPacket) []byte {
	session.SetFlagKicked()
	return nil
}
