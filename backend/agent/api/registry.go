package api

import "sync"

// store client sessions
var (
	Registry sync.Map
)
