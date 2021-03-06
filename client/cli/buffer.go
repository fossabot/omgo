package main

import (
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/kit/utils"
)

// Buffer controls the packet send to server
type Buffer struct {
	ctrl    chan struct{}
	pending chan []byte
	conn    net.Conn
}

// packet sending procedure
func (buf *Buffer) send(session *Session, data []byte) {
	// we don't send empty package
	if data == nil {
		return
	}

	// encryption
	// NOT_ENCRYPTED -> KEY_EXCHANGED -> ENCRYPTED
	if session.IsFlagEncryptedSet() {
		// encryption is enabled
		session.Encoder.XORKeyStream(data, data)
	} else if session.IsFlagKeyExchangedSet() {
		// key is exchanged, encryption is not yet enabled
		session.ClearFlagKeyExchanged()
		session.SetFlagEncrypted()
	}

	// queue data for sending
	select {
	case buf.pending <- data:
	default:
		// packet will be dropped if it exceeds txQueueLength
		log.Warning("pending full")
	}
	return
}

// packet sending goroutine
func (buf *Buffer) start() {
	defer utils.PrintPanicStack()
	for {
		select {
		case data := <-buf.pending:
			buf.rawSend(data)
		case <-buf.ctrl:
			// receive session end signal
			return
		}
	}
}

func (buf *Buffer) rawSend(data []byte) bool {
	// write data
	n, err := buf.conn.Write(data)
	if err != nil {
		log.Warningf("Error while sending reply data, bytes:%v reason:%v", n, err)
		return false
	}
	return true
}

func newBuffer(conn net.Conn, ctrl chan struct{}, txqueuelen int) *Buffer {
	buf := Buffer{conn: conn, ctrl: ctrl}
	buf.pending = make(chan []byte, txqueuelen)
	return &buf
}
