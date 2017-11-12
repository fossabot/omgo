package api

import (
	"fmt"
	"strings"

	"crypto/rand"
	"crypto/rc4"

	"io"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/kit/ecdh"
	pbdb "github.com/master-g/omgo/proto/grpc/db"
	pbgame "github.com/master-g/omgo/proto/grpc/game"
	pc "github.com/master-g/omgo/proto/pb/common"
	"golang.org/x/net/context"
	"google.golang.org/grpc/metadata"
)

func makeErrorResponse(sess *Session, msg string, statusCode pc.ResultCode, session *Session) []byte {
	rsp := &pc.S2CHandshakeRsp{}
	hdr := genRspHeader(pc.Cmd_HANDSHAKE_RSP)
	hdr.Msg = msg
	hdr.Status = int32(statusCode)
	session.SetFlagKicked()
	return makeResponse(sess, hdr, rsp)
}

// ProcHandshakeReq handles client handshake request
// 1. verify token via database service
// 2. kick previous client if exists
// 3. exchange cryption seed
// 4. connect to game server with gRPC stream
func ProcHandshakeReq(session *Session, inPacket []byte) []byte {
	rsp := &pc.S2CHandshakeRsp{}
	req := &pc.C2SHandshakeReq{}

	msg := ""
	header, err := unpackPacket(inPacket, req)
	if err != nil {
		msg = fmt.Sprintf("invalid protobuf: %v", err)
	} else if header.ClientInfo == nil {
		msg = "invalid header, client_info missing"
	} else if header.ClientInfo.Usn == 0 {
		msg = "invalid header, invalid usn"
	} else if req.Token == "" {
		msg = "invalid token"
	} else if len(req.RecvSeed) != 32 || len(req.SendSeed) != 32 {
		msg = "invalid seed"
	}

	if msg != "" {
		log.Warningf("Handshake proc error :%v", msg)
		return makeErrorResponse(session, msg, pc.ResultCode_RESULT_INVALID, session)
	}

	// validate user token
	dbConn := DataServicePool.NextClient()
	if dbConn == nil {
		msg = "dataservice not connected yet"
		log.Error(msg)
		return makeErrorResponse(session, msg, pc.ResultCode_RESULT_INTERNAL_ERROR, session)
	}

	usn := header.ClientInfo.Usn
	token := req.Token

	dbClient := pbdb.NewDBServiceClient(dbConn)
	userKey := &pbdb.DB_UserEntry{Usn: usn}
	dbRsp, err := dbClient.UserExtraInfoQuery(context.Background(), userKey)
	if err != nil {
		log.Errorf("error while query user extra info:%v", usn)
		return makeErrorResponse(session, "", pc.ResultCode_RESULT_INTERNAL_ERROR, session)
	}

	if dbRsp.Result.Status != int32(pbdb.DB_STATUS_OK) || dbRsp.User.Token == "" {
		log.Warningf("user extra info not found:%v", usn)
		return makeErrorResponse(session, "", pc.ResultCode_RESULT_INTERNAL_ERROR, session)
	}

	if strings.Compare(token, dbRsp.User.Token) != 0 {
		msg = "invalid token"
		log.Info(msg)
		return makeErrorResponse(session, msg, pc.ResultCode_RESULT_INVALID, session)
	}

	// kick previous session if existed
	p, ok := Registry.Load(usn)
	if ok {
		if prevSession, ok := p.(*Session); ok {
			kickNotify := &pc.S2CKickNotify{
				Reason: pc.KickReason_KICK_LOGIN_ELSEWHERE,
				Msg:    session.IP.String(),
			}
			kickHdr := genRspHeader(pc.Cmd_KICK_NOTIFY)
			prevSession.Mailbox <- makeResponse(session, kickHdr, kickNotify)
			prevSession.SetFlagKicked()
		}
	}

	Registry.Store(usn, session)

	// exchange seed
	curve := ecdh.NewCurve25519ECDH()
	x1, e1 := curve.GenerateECKeyBuf(rand.Reader)
	key1 := curve.GenerateSharedSecretBuf(x1, req.GetSendSeed())
	x2, e2 := curve.GenerateECKeyBuf(rand.Reader)
	key2 := curve.GenerateSharedSecretBuf(x2, req.GetRecvSeed())

	encoder, err := rc4.NewCipher(key2)
	if err != nil {
		log.Error(err)
		return makeErrorResponse(session, err.Error(), pc.ResultCode_RESULT_INTERNAL_ERROR, session)
	}
	decoder, err := rc4.NewCipher(key1)
	if err != nil {
		log.Error(err)
		return makeErrorResponse(session, err.Error(), pc.ResultCode_RESULT_INTERNAL_ERROR, session)
	}
	session.Encoder = encoder
	session.Decoder = decoder
	session.SetFlagKeyExchanged()

	rsp.SendSeed = e1
	rsp.RecvSeed = e2

	// connect to other services
	session.Usn = usn
	session.Token = token
	session.GameServerID = config.GameServerName
	session.SetFlagAuth()

	conn := GameServerPool.GetClient(config.GameServerName)
	if conn == nil {
		msg = fmt.Sprintf("cannot get game service:%v", session.GameServerID)
		log.Error(msg)
		return makeErrorResponse(session, msg, pc.ResultCode_RESULT_INTERNAL_ERROR, session)
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

	// all ok
	rspHeader := genRspHeader(pc.Cmd_HANDSHAKE_RSP)
	return makeResponse(session, rspHeader, rsp)
}
