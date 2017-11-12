package main

import (
	"crypto/rc4"
	"encoding/hex"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/kit/ecdh"
	pc "github.com/master-g/omgo/proto/pb/common"
)

var Handlers map[int32]func(*Session, *pc.RspHeader) []byte

func init() {
	Handlers = map[int32]func(*Session, *pc.RspHeader) []byte{
		int32(pc.Cmd_HEART_BEAT_RSP): ProcHeartBeatRsp,
		int32(pc.Cmd_HANDSHAKE_RSP):  ProcHandshakeRsp,
		int32(pc.Cmd_KICK_NOTIFY):    ProcKickNotify,
	}
}

func ProcHeartBeatRsp(session *Session, header *pc.RspHeader) []byte {
	log.Info("receive server heartbeat response")
	return nil
}

func ProcHandshakeRsp(session *Session, header *pc.RspHeader) []byte {
	rsp := &pc.S2CHandshakeRsp{}
	err := proto.Unmarshal(header.Body, rsp)
	if err != nil {
		log.Errorf("error while parsing proto:%v", err)
		return nil
	}

	curve := ecdh.NewCurve25519ECDH()
	keySend := curve.GenerateSharedSecretBuf(session.privateSend, rsp.GetSendSeed())
	keyRecv := curve.GenerateSharedSecretBuf(session.privateRecv, rsp.GetRecvSeed())

	session.Encoder, err = rc4.NewCipher(keySend)
	if err != nil {
		log.Fatalf("error while creating encoder:%v", err)
	}
	session.Decoder, err = rc4.NewCipher(keyRecv)
	if err != nil {
		log.Fatalf("error while creating decoder:%v", err)
	}

	log.Infof("encoder seed:%v", strings.ToUpper(hex.EncodeToString(keySend)))
	log.Infof("decoder seed:%v", strings.ToUpper(hex.EncodeToString(keyRecv)))

	session.SetFlagEncrypted()

	return nil
}

func ProcKickNotify(session *Session, header *pc.RspHeader) []byte {
	rsp := &pc.S2CKickNotify{}
	err := proto.Unmarshal(header.Body, rsp)
	if err != nil {
		log.Errorf("error while parsing proto:%v", err)
		return nil
	}
	log.Warnf("kicked by server, msg:%v reason:%v", rsp.Msg, rsp.Reason)
	session.SetFlagKicked()
	close(session.Mailbox)
	return nil
}
