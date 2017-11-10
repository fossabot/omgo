package main

import (
	"crypto/rand"
	"crypto/rc4"
	"encoding/binary"
	"io"
	"net"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/golang/protobuf/proto"
	"github.com/master-g/omgo/kit/ecdh"
	"github.com/master-g/omgo/kit/packet"
	"github.com/master-g/omgo/kit/utils"
	pc "github.com/master-g/omgo/proto/pb/common"
)

type Session struct {
	Usn         uint64
	Token       string
	Die         chan struct{}
	Mailbox     chan []byte
	Encoder     *rc4.Cipher
	Decoder     *rc4.Cipher
	Flag        int
	Conn        net.Conn
	Out         *Buffer
	privateSend []byte
	privateRecv []byte
	Seq         uint32
}

const (
	// FlagConnected  indicates the connection status of the session
	FlagConnected = 0x01
	// FlagKeyExchanged indicates the key exchange process has completed
	FlagKeyExchanged = 0x2
	// FlagEncrypted indicates the transmission of this session is encrypted
	FlagEncrypted = 0x4
	// FlagKicked indicates the client has been kicked out
	FlagKicked = 0x8
	// FlagAuth indicates the session has been authorized
	FlagAuth = 0x10
)

// SetFlagConnected sets the connected bit
func (s *Session) SetFlagConnected() *Session {
	s.Flag |= FlagConnected
	return s
}

// ClearFlagConnected clears the connected bit
func (s *Session) ClearFlagConnected() *Session {
	s.Flag &^= FlagConnected
	return s
}

// IsFlagConnectedSet return true if the connected bit is set
func (s *Session) IsFlagConnectedSet() bool {
	return s.Flag&FlagConnected != 0
}

// SetFlagKeyExchanged sets the connected bit
func (s *Session) SetFlagKeyExchanged() *Session {
	s.Flag |= FlagKeyExchanged
	return s
}

// ClearFlagKeyExchanged clears the key exchanged bit
func (s *Session) ClearFlagKeyExchanged() *Session {
	s.Flag &^= FlagKeyExchanged
	return s
}

// IsFlagKeyExchangedSet return true if the key exchanged bit is set
func (s *Session) IsFlagKeyExchangedSet() bool {
	return s.Flag&FlagKeyExchanged != 0
}

// SetFlagEncrypted sets the encrypted bit
func (s *Session) SetFlagEncrypted() *Session {
	s.Flag |= FlagEncrypted
	return s
}

// ClearFlagEncrypted clears the encrypted bit
func (s *Session) ClearFlagEncrypted() *Session {
	s.Flag &^= FlagEncrypted
	return s
}

// IsFlagEncryptedSet returns true if the encrypted bit is set
func (s *Session) IsFlagEncryptedSet() bool {
	return s.Flag&FlagEncrypted != 0
}

// SetFlagKicked sets the kicked bit
func (s *Session) SetFlagKicked() *Session {
	s.Flag |= FlagKicked
	return s
}

// ClearFlagKicked clears the kicked bit
func (s *Session) ClearFlagKicked() *Session {
	s.Flag &^= FlagKicked
	return s
}

// IsFlagKickedSet returns true if the kicked bit is set
func (s *Session) IsFlagKickedSet() bool {
	return s.Flag&FlagKicked != 0
}

// SetFlagAuth sets the auth bit
func (s *Session) SetFlagAuth() *Session {
	s.Flag |= FlagAuth
	return s
}

// ClearFlagAuth clears the auth bit
func (s *Session) ClearFlagAuth() *Session {
	s.Flag &^= FlagAuth
	return s
}

// IsFlagAuthSet returns true if the auth bit is set
func (s *Session) IsFlagAuthSet() bool {
	return s.Flag&FlagAuth != 0
}

//------------------------------------------------------------------------------
// goroutines
//------------------------------------------------------------------------------

func (s *Session) Loop(in chan []byte) {
	defer utils.PrintPanicStack()

	minTimer := time.After(time.Minute)

	defer func() {
		close(s.Die)
	}()

	for {
		select {
		case msg, ok := <-in:
			if !ok {
				return
			}

			if result := s.Route(msg); result != nil {
				s.Out.send(s, result)
			}
		case mail, ok := <-s.Mailbox:
			if ok {
				s.Out.send(s, mail)
			}
		case <-minTimer:
			s.TimeWork()
			minTimer = time.After(time.Minute)
		}

		if s.IsFlagKickedSet() {
			return
		}
	}

	s.Flag = 0
}

func (s *Session) startLoop() {
	defer utils.PrintPanicStack()
	defer s.Conn.Close()
	header := make([]byte, 2)
	in := make(chan []byte)
	defer func() {
		close(in)
	}()

	s.Die = make(chan struct{})
	s.Mailbox = make(chan []byte, 128)
	s.Out = newBuffer(s.Conn, s.Die, 128)

	go s.Out.start()

	go s.Loop(in)

	for {
		n, err := io.ReadFull(s.Conn, header)
		if err != nil {
			log.Warningf("read header failed: %v %v bytes read", err, n)
			return
		}
		size := binary.BigEndian.Uint16(header)

		payload := make([]byte, size)
		n, err = io.ReadFull(s.Conn, payload)
		if err != nil {
			log.Warningf("read payload failed: %v expect: %v actual read: %v", err, size, n)
			return
		}

		select {
		case in <- payload:
		case <-s.Die:
			log.Warningf("connection closed by logic, flag: %v", s.Flag)
			return
		}
	}
}

func (s *Session) Route(msg []byte) []byte {
	defer utils.PrintPanicStack()
	// decrypt
	if s.IsFlagEncryptedSet() {
		s.Decoder.XORKeyStream(msg, msg)
	}
	// packet reader
	reader := packet.NewRawPacketReader(msg)

	// read cmd
	cmdValue, err := reader.ReadS32()
	if err != nil {
		log.Errorf("read packet cmd failed:%v", err)
		s.SetFlagKicked()
		return nil
	}
	cmd := pc.Cmd(cmdValue)

	// route message
	var ret []byte
	if cmd > pc.Cmd_CMD_COMMON_END {
		log.Info("stream function not implemented yet")
		return nil
	} else {
		shell.ShowPrompt(false)
		shell.Println("")
		if h := Handlers[cmdValue]; h != nil {
			ret = h(s, reader)
		} else {
			log.Errorf("no handler for cmd:%v", cmd)
			return nil
		}
		shell.ShowPrompt(true)
	}
	return ret
}

func (s *Session) TimeWork() {
	shell.ShowPrompt(false)
	s.Heartbeat()
	shell.ShowPrompt(true)
}

//------------------------------------------------------------------------------
// Interface
//------------------------------------------------------------------------------

func (s *Session) makeHeader(cmd pc.Cmd) *pc.Header {
	s.Seq++
	header := &pc.Header{
		Version: 1,
		Cmd:     int32(cmd),
		Seq:     s.Seq,
		ClientInfo: &pc.ClientInfo{
			Usn:       s.Usn,
			Timestamp: utils.Timestamp(),
		},
	}

	return header
}

func (s *Session) send(header *pc.Header, msg proto.Message) {
	var msgBuf []byte
	var err error
	if msg != nil {
		msgBuf, err = proto.Marshal(msg)
		if err != nil {
			log.Errorf("error while marshal %v error:%v", msg, err)
			return
		}
	}

	header.BodySize = int32(len(msgBuf))
	headerBuf, err := proto.Marshal(header)
	if err != nil {
		log.Errorf("error while marshal header, error:%v", err)
		return
	}

	headerSize := len(headerBuf)
	buf := make([]byte, headerSize+int(header.BodySize)+2)
	binary.BigEndian.PutUint16(buf, uint16(headerSize))
	copy(buf[2:], headerBuf)
	if msg != nil {
		copy(buf[2+headerSize:], msgBuf)
	}
	log.Infof("sending %v %v", header, msg)
	s.Out.send(s, buf)
}

// Connect to agent server
func (s *Session) Connect(address string) (sess *Session) {
	if s.IsFlagConnectedSet() {
		log.Error("already connected")
		return
	}

	sess = s
	log.Infof("connecting to %v", address)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Errorf("could not connect to server:%v error:%v", address, err)
		return
	}
	s.Conn = conn
	host, port, err := net.SplitHostPort(conn.RemoteAddr().String())
	if err != nil {
		log.Errorf("get remote addr failed:%v", err)
		return
	}
	s.SetFlagConnected()
	log.Infof("server %v:%v connected", host, port)

	go s.startLoop()

	return
}

// Close connection to agent server
func (s *Session) Close() *Session {
	s.SetFlagKicked()
	return s
}

func (s *Session) Heartbeat() {
	log.Info("sending heartbeat")
	header := s.makeHeader(pc.Cmd_HEART_BEAT_REQ)
	s.send(header, nil)
}

func (s *Session) Handshake() {
	log.Info("about to login")
	header := s.makeHeader(pc.Cmd_HANDSHAKE_REQ)

	curve := ecdh.NewCurve25519ECDH()
	x1, e1 := curve.GenerateECKeyBuf(rand.Reader)
	x2, e2 := curve.GenerateECKeyBuf(rand.Reader)

	s.privateSend = x1
	s.privateRecv = x2

	req := &pc.C2SHandshakeReq{
		SendSeed: e1,
		RecvSeed: e2,
		Token:    s.Token,
	}

	s.send(header, req)
}

func (s *Session) Bye() {
	log.Info("sending bye")
	header := s.makeHeader(pc.Cmd_OFFLINE_REQ)
	s.send(header, nil)
	s.Close()
}
