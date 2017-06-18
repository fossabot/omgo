package main

import (
	"crypto/rc4"
	"github.com/master-g/omgo/utils"
	"net"
	"time"
	"unicode"
)

type Session struct {
	Usn     uint64
	Token   string
	Die     chan struct{}
	Encoder *rc4.Cipher
	Decoder *rc4.Cipher
	Flag    int
	Conn    net.Conn
}

const (
	Salt = "DH"
)

const (
	// FlagKeyExchanged indicates the key exchange process has completed
	FlagKeyExchanged = 0x1
	// FlagEncrypted indicates the transmission of this session is encrypted
	FlagEncrypted = 0x2
	// FlagKicked indicates the client has been kicked out
	FlagKicked = 0x4
	// FlagAuth indicates the session has been authorized
	FlagAuth = 0x8
)

// SetFlagKeyExchanged sets the key exchanged bit
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

func (s *Session) Loop(in chan []byte, out *Buffer) {
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
				out.send(s, result)
			}
		case <-minTimer:
			s.TimeWork(out)
			minTimer = time.After(time.Minute)
		}

		if s.IsFlagKickedSet() {
			return
		}
	}
}

func (s *Session) Route(msg []byte) []byte {
	return nil
}

func (s *Session) TimeWork(out *Buffer) {

}
