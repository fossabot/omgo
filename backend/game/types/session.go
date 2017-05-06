package types

const (
	// FlagKicked indicates the client has been kicked out
	FlagKicked = 0x4
)

// Session holds the context of a client having conversation with game
type Session struct {
	Flag   int32 // Session flag
	UserID int32 // User ID
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
