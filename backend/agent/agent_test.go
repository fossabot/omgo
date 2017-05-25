package main

import (
	"crypto/rand"
	"crypto/rc4"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/security/ecdh"
	"github.com/master-g/omgo/utils"
)

var (
	address string
	encoder *rc4.Cipher
	decoder *rc4.Cipher
)

const (
	Salt = "DH"
)

func init() {
	address = utils.GetLocalIP() + ":8888"
}

func connect(t *testing.T) net.Conn {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		t.Fatalf("could not connect to server %v, error: %v", address, err)
	}

	return conn
}

func send(conn net.Conn, data []byte, t *testing.T) {
	size := len(data)
	cache := make([]byte, 2+size)
	binary.BigEndian.PutUint16(cache, uint16(size))
	copy(cache[2:], data)
	_, err := conn.Write(cache[:size+2])
	if err != nil {
		t.Fatalf("error while sending data: %v", err)
	}
}

func recv(conn net.Conn, t *testing.T) []byte {
	header := make([]byte, 2)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		t.Fatalf("error while reading header:%v", err)
	}
	size := binary.BigEndian.Uint16(header)
	// data
	payload := make([]byte, size)
	n, err := io.ReadFull(conn, payload)
	if err != nil {
		t.Fatalf("read payload failed: %v expect: %v actual read: %v", err, size, n)
	}

	return payload
}

func TestHeartBeat(t *testing.T) {
	conn := connect(t)
	defer conn.Close()

	reqPacket := packet.NewRawPacket()
	reqPacket.WriteS32(int32(proto_common.Cmd_HEART_BEAT_REQ))
	send(conn, reqPacket.Data(), t)

	recv(conn, t)
}

func TestGetSeed(t *testing.T) {
	conn := connect(t)
	defer conn.Close()

	reqPacket := packet.NewRawPacket()
	reqPacket.WriteS32(int32(proto_common.Cmd_GET_SEED_REQ))

	req := &proto_common.C2SGetSeedReq{}

	curve := ecdh.NewCurve25519ECDH()
	x1, e1 := curve.GenerateECKeyBuf(rand.Reader)
	x2, e2 := curve.GenerateECKeyBuf(rand.Reader)

	req.SendSeed = e1
	req.RecvSeed = e2

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatalf("error while create request:%v", err)
	}
	reqPacket.WriteBytes(data)
	send(conn, reqPacket.Data(), t)

	rspBody := recv(conn, t)
	rsp := &proto_common.S2CGetSeedRsp{}
	reader := packet.NewRawPacketReader(rspBody)
	cmd, err := reader.ReadS32()
	buf, err := reader.ReadBytes()
	if cmd != int32(proto_common.Cmd_GET_SEED_RSP) || err != nil {
		t.Fatalf("error while parsing response cmd:%v error:%v", cmd, err)
	}

	err = proto.Unmarshal(buf, rsp)
	if err != nil {
		t.Fatalf("error while parsing proto:%v", err)
	}

	key1 := curve.GenerateSharedSecretBuf(x1, rsp.GetSendSeed())
	key2 := curve.GenerateSharedSecretBuf(x2, rsp.GetRecvSeed())

	encoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key2)))
	if err != nil {
		t.Fatalf("error while creating encoder:%v", err)
	}
	decoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key1)))
	if err != nil {
		t.Fatalf("error while creating decoder:%v", err)
	}
}
