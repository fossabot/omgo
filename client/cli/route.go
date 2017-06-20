package main

import (
	"crypto/rc4"
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"
	"github.com/master-g/omgo/net/packet"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/security/ecdh"
	"github.com/prometheus/common/log"
)

var Handlers map[int32]func(*Session, *packet.RawPacket) []byte

func init() {
	Handlers = map[int32]func(*Session, *packet.RawPacket) []byte{
		int32(pc.Cmd_HEART_BEAT_RSP): ProcHeartBeatRsp,
		int32(pc.Cmd_LOGIN_RSP):      ProcLoginRsp,
		int32(pc.Cmd_GET_SEED_RSP):   ProcGetSeedRsp,
		int32(pc.Cmd_KICK_NOTIFY):    ProcKickNotify,
	}
}

func readPacketBody(packet *packet.RawPacket) []byte {
	body, err := packet.ReadBytes()
	if err != nil {
		log.Fatalf("error while reading buffer from packet:%v", err)
		return nil
	}

	return body
}

func ProcHeartBeatRsp(session *Session, packet *packet.RawPacket) []byte {
	log.Info("receive server heartbeat response")
	return nil
}

func ProcLoginRsp(session *Session, packet *packet.RawPacket) []byte {
	return nil
}

func ProcGetSeedRsp(session *Session, packet *packet.RawPacket) []byte {
	rspBody := readPacketBody(packet)
	rsp := &pc.S2CGetSeedRsp{}
	err := proto.Unmarshal(rspBody, rsp)
	if err != nil {
		log.Errorf("error while parsing proto:%v", err)
		return nil
	}

	curve := ecdh.NewCurve25519ECDH()
	keySend := curve.GenerateSharedSecretBuf(session.privateSend, rsp.GetSendSeed())
	keyRecv := curve.GenerateSharedSecretBuf(session.privateRecv, rsp.GetRecvSeed())

	session.Encoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, keySend)))
	if err != nil {
		log.Fatalf("error while creating encoder:%v", err)
	}
	session.Decoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, keyRecv)))
	if err != nil {
		log.Fatalf("error while creating decoder:%v", err)
	}

	log.Infof("encoder seed:%v", strings.ToUpper(hex.EncodeToString(keySend)))
	log.Infof("decoder seed:%v", strings.ToUpper(hex.EncodeToString(keyRecv)))

	session.SetFlagEncrypted()

	return nil
}

func ProcKickNotify(session *Session, packet *packet.RawPacket) []byte {
	rspBody := readPacketBody(packet)
	rsp := &pc.S2CKickNotify{}
	err := proto.Unmarshal(rspBody, rsp)
	if err != nil {
		log.Errorf("error while parsing proto:%v", err)
		return nil
	}
	log.Warnf("kicked by server, msg:%v reason:%v", rsp.Msg, rsp.Reason)
	session.SetFlagKicked()
	return nil
}
