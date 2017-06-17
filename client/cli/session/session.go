package session

import (
	"crypto/rand"
	"crypto/rc4"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/net/packet"
	pc "github.com/master-g/omgo/proto/pb/common"
	"github.com/master-g/omgo/security/ecdh"
	"github.com/master-g/omgo/utils"
)

type Session struct {
	IsConnected bool
	conn        net.Conn
	encrypted   bool
	encoder     *rc4.Cipher
	decoder     *rc4.Cipher
}

const (
	Salt = "DH"
)

func makePacket(cmd pc.Cmd) *packet.RawPacket {
	pk := packet.NewRawPacket()
	pk.WriteS32(int32(cmd))
	return pk
}

func NewSession(address string) *Session {
	return new(Session)
}

func (s *Session) Connect(addr string) {
	log.Infof("connecting to %v", addr)
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("could not connect to server:%v, error:%v", addr, err)
		return
	}
	s.conn = conn
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Errorf("get remote addr failed:%v", err)
		return
	}

	s.IsConnected = true
	log.Infof("server %v%v connected", host, port)
}

func (s *Session) Close() {
	if s.conn != nil {
		log.Infof("disconnecting from %v", s.conn.RemoteAddr().String())
		s.conn.Close()
		s.conn = nil
	}
	s.IsConnected = false
}

func (s *Session) Send(data []byte) {
	if s.encrypted {
		s.encoder.XORKeyStream(data, data)
	}

	size := len(data)
	cache := make([]byte, 2+size)
	binary.BigEndian.PutUint16(cache, uint16(size))
	copy(cache[2:], data)
	_, err := s.conn.Write(cache[:size+2])
	if err != nil {
		log.Fatalf("error while sending data %v", err)
		return
	}
	log.Infof("--> %v bytes", size+2)
}

func (s *Session) Recv() []byte {
	header := make([]byte, 2)
	_, err := io.ReadFull(s.conn, header)
	if err != nil {
		log.Fatalf("error while reading header:%v", err)
		return nil
	}

	size := binary.BigEndian.Uint16(header)
	// data
	payload := make([]byte, size)
	n, err := io.ReadFull(s.conn, payload)
	if err != nil {
		log.Fatalf("read payload failed:%v expect:%v actual read:%v", err, size, n)
		return nil
	}

	if s.decoder != nil {
		s.decoder.XORKeyStream(payload, payload)
	}

	log.Infof("<-- %v bytes", size+2)

	return payload
}

func (s *Session) Heartbeat() {
	log.Info("sending heartbeat")
	reqPacket := makePacket(pc.Cmd_HEART_BEAT_REQ)
	s.Send(reqPacket.Data())
	rspPacket := s.Recv()
	if rspPacket != nil {
		reader := packet.NewRawPacketReader(rspPacket)
		cmd, err := reader.ReadS32()
		if err != nil {
			log.Fatalf("read cmd failed:%v", err)
			return
		}
		if cmd != int32(pc.Cmd_HEART_BEAT_RSP) {
			log.Fatalf("expect %v got %v", pc.Cmd_HEART_BEAT_RSP, cmd)
			return
		}
		log.Info("recv heartbeat response from server")
	}
}

func (s *Session) ExchangeKey() {
	log.Info("about to exchange key")
	reqPacket := makePacket(pc.Cmd_GET_SEED_REQ)
	req := &pc.C2SGetSeedReq{}

	curve := ecdh.NewCurve25519ECDH()
	x1, e1 := curve.GenerateECKeyBuf(rand.Reader)
	x2, e2 := curve.GenerateECKeyBuf(rand.Reader)

	req.SendSeed = e1
	req.RecvSeed = e2

	data, err := proto.Marshal(req)
	if err != nil {
		log.Fatalf("error while create request:%v", err)
	}
	reqPacket.WriteBytes(data)
	s.Send(reqPacket.Data())

	rspBody := s.Recv()
	rsp := &pc.S2CGetSeedRsp{}
	reader := packet.NewRawPacketReader(rspBody)
	cmd, err := reader.ReadS32()
	buf, err := reader.ReadBytes()
	if cmd != int32(pc.Cmd_GET_SEED_RSP) || err != nil {
		log.Fatalf("error while parsing response cmd:%v error%v", cmd, err)
	}

	err = proto.Unmarshal(buf, rsp)
	if err != nil {
		log.Fatalf("error while parsing proto:%v", err)
	}

	key1 := curve.GenerateSharedSecretBuf(x1, rsp.GetSendSeed())
	key2 := curve.GenerateSharedSecretBuf(x2, rsp.GetRecvSeed())

	s.encoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key1)))
	if err != nil {
		log.Fatalf("error while creating encoder:%v", err)
	}
	s.decoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key2)))
	if err != nil {
		log.Fatalf("error while creating decoder:%v", err)
	}

	log.Infof("encoder seed:%v", strings.ToUpper(hex.EncodeToString(key1)))
	log.Infof("decoder seed:%v", strings.ToUpper(hex.EncodeToString(key2)))

	s.encrypted = true
}

func (s *Session) Login(usn uint64, token string) {
	log.Info("about to login")
	reqPacket := makePacket(pc.Cmd_LOGIN_REQ)
	req := &pc.C2SLoginReq{
		Timestamp: utils.Timestamp(),
		Usn:       usn,
		Token:     token,
	}
	data, err := proto.Marshal(req)
	if err != nil {
		log.Fatalf("error while create request:%v", err)
	}
	reqPacket.WriteBytes(data)
	s.Send(reqPacket.Data())
}

func (s *Session) Bye() {
	log.Info("sending bye")
	reqPacket := makePacket(pc.Cmd_OFFLINE_REQ)
	s.Send(reqPacket.Data())
	s.Close()
}
