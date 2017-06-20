package types

import (
	"crypto/rc4"
	"net"
	"time"

	pb "github.com/master-g/omgo/proto/grpc/game"
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

// Session holds the context of a client having conversation with agent
type Session struct {
	Usn               uint64                      // User serial number
	Token             string                      // Session token generate at http login
	Die               chan struct{}               // Session close signal
	Flag              int32                       // Session flag
	IP                net.IP                      // Client IP address
	Port              string                      // Client Port
	MQ                chan pb.Game_Frame          // Channel of async messages send back to client
	Encoder           *rc4.Cipher                 // Encrypt
	Decoder           *rc4.Cipher                 // Decrypt
	GSID              string                      // Game server ID
	Stream            pb.GameService_StreamClient // Data stream send to game server
	ConnectTime       time.Time                   // Timestamp of TCP connection established
	PacketTime        time.Time                   // Timestamp of current packet arrived
	LastPacketTime    time.Time                   // Timestamp of previous packet arrived
	PacketCount       uint32                      // Total packets received
	PacketCountPerMin int                         // Packets received per minute
	Mailbox           chan []byte                 // Mailbox for internal communication
}

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
