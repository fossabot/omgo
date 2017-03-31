package services

import "google.golang.org/grpc"

// single connection
type client struct {
	key  string
	conn *grpc.ClientConn
}

type service struct {
    clients  []client
    idx uint32
}

type servicePool struct {
    root string
    names map[string]bool
}
