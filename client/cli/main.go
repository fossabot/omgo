package main

import (
	"crypto/rand"
	"crypto/rc4"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/abiosoft/ishell"
	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/proto/pb"
	"github.com/master-g/omgo/security/ecdh"
	"github.com/master-g/omgo/utils"
	"gopkg.in/urfave/cli.v2"
	"io"
	"net"
	"os"
	"strings"
)

var (
	address   string
	conn      net.Conn
	encrypted bool
	encoder   *rc4.Cipher
	decoder   *rc4.Cipher
)

const (
	Salt = "DH"
)

func connect(addr string) {
	log.Infof("connecting to %v", addr)

	_conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Fatalf("could not connect to server %v, error: %v", address, err)
	} else {
		address = addr
		conn = _conn
		host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
		if err != nil {
			log.Error("get remote address failed: ", err)
			return
		} else {
			log.Infof("server %v:%v connected", host, port)
		}
	}
}

func disconnect() bool {
	if conn != nil {
		log.Infof("disconnecting from %v", address)
		conn.Close()
		conn = nil
		return true
	}
	return false
}

func send(data []byte) {
	if encrypted {
		encoder.XORKeyStream(data, data)
	}

	size := len(data)
	cache := make([]byte, 2+size)
	binary.BigEndian.PutUint16(cache, uint16(size))
	copy(cache[2:], data)
	_, err := conn.Write(cache[:size+2])
	if err != nil {
		log.Fatalf("error while sending data: %v", err)
	}

	log.Infof("--> %v bytes", size+2)
}

func recv() []byte {
	header := make([]byte, 2)
	_, err := io.ReadFull(conn, header)
	if err != nil {
		log.Fatalf("error while reading header:%v", err)
		return nil
	}

	size := binary.BigEndian.Uint16(header)
	// data
	payload := make([]byte, size)
	n, err := io.ReadFull(conn, payload)
	if err != nil {
		log.Fatalf("read payload failed: %v expect: %v actual read: %v", err, size, n)
		return nil
	}

	if decoder != nil {
		decoder.XORKeyStream(payload, payload)
	}

	log.Infof("<-- %v bytes", size+2)

	return payload
}

func heartbeat() {
	log.Info("sending heartbeat")
	reqPacket := packet.NewRawPacket()
	reqPacket.WriteS32(int32(proto_common.Cmd_HEART_BEAT_REQ))
	send(reqPacket.Data())
	rspPacket := recv()
	if rspPacket != nil {
		reader := packet.NewRawPacketReader(rspPacket)
		cmd, err := reader.ReadS32()
		if err != nil {
			log.Fatalf("read cmd failed:%v", err)
			return
		}
		if cmd != int32(proto_common.Cmd_HEART_BEAT_RSP) {
			log.Fatalf("expect %v got %v", proto_common.Cmd_HEART_BEAT_RSP, cmd)
			return
		}
		log.Info("recv heartbeat response from server")
	}
}

func keyExchange() {
	log.Info("about to exchange key")
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
		log.Fatalf("error while create request:%v", err)
	}
	reqPacket.WriteBytes(data)
	send(reqPacket.Data())

	rspBody := recv()
	rsp := &proto_common.S2CGetSeedRsp{}
	reader := packet.NewRawPacketReader(rspBody)
	cmd, err := reader.ReadS32()
	buf, err := reader.ReadBytes()
	if cmd != int32(proto_common.Cmd_GET_SEED_RSP) || err != nil {
		log.Fatalf("error while parsing response cmd:%v error:%v", cmd, err)
	}

	err = proto.Unmarshal(buf, rsp)
	if err != nil {
		log.Fatalf("error while parsing proto:%v", err)
	}

	key1 := curve.GenerateSharedSecretBuf(x1, rsp.GetSendSeed())
	key2 := curve.GenerateSharedSecretBuf(x2, rsp.GetRecvSeed())

	encoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key1)))
	if err != nil {
		log.Fatalf("error while creating encoder:%v", err)
	}
	decoder, err = rc4.NewCipher([]byte(fmt.Sprintf("%v%v", Salt, key2)))
	if err != nil {
		log.Fatalf("error while creating decoder:%v", err)
	}

	log.Infof("encoder seed:%v", strings.ToUpper(hex.EncodeToString(key2)))
	log.Infof("decoder seed:%v", strings.ToUpper(hex.EncodeToString(key1)))

	encrypted = true
}

func main() {
	log.SetLevel(log.DebugLevel)
	defer disconnect()
	defer utils.PrintPanicStack()

	app := &cli.App{
		Name:    "client",
		Usage:   "a cli-client for testing omgo",
		Version: "1.0",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases: []string{"a"},
				Name:    "address",
				Usage:   "connect to address:port",
				Value:   ":8888",
			},
		},
		Action: func(c *cli.Context) error {
			address = c.String("address")
			log.Info(address)
			return nil
		},
	}
	app.Run(os.Args)

	shell := ishell.New()
	shell.AddCmd(&ishell.Cmd{
		Name: "conn",
		Help: "conn address:port",
		Func: func(c *ishell.Context) {
			if len(c.Args) > 0 {
				address = c.Args[0]
			}
			disconnect()
			connect(address)
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "disconn",
		Help: "disconnect from server",
		Func: func(c *ishell.Context) {
			if disconnect() {
				c.Println("disconnected from server")
			} else {
				c.Println("no connection")
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "heartbeat",
		Help: "sending heartbeat to server",
		Func: func(c *ishell.Context) {
			if conn == nil {
				c.Println("no connection")
				return
			} else {
				heartbeat()
			}
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "go",
		Help: "go through all tests",
		Func: func(c *ishell.Context) {
			disconnect()
			connect(address)
			heartbeat()
			keyExchange()
		},
	})
	shell.AddCmd(&ishell.Cmd{
		Name: "exchangekey",
		Help: "exchange public key with server",
		Func: func(c *ishell.Context) {
			if conn == nil {
				c.Println("no connection")
				return
			} else {
				keyExchange()
			}
		},
	})

	shell.Start()
}
