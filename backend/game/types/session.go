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
