package main

import (
	"net"

	"encoding/binary"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/utils"
)

// Buffer controls the packet send to server
type Buffer struct {
	ctrl    chan struct{}
	pending chan []byte
	conn    net.Conn
	cache   []byte
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
		// pakcet will be dropped if it exceeds txQueueLength
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
	size := len(data)
	binary.BigEndian.PutUint16(buf.cache, uint16(size))
	copy(buf.cache[2:], data)

	// write data
	n, err := buf.conn.Write(buf.cache[:size+2])
	if err != nil {
		log.Warningf("Error while sending reply data, bytes:%v reason:%v", n, err)
		return false
	}
	return true
}

func newBuffer(conn net.Conn, ctrl chan struct{}, txqueuelen int) *Buffer {
	buf := Buffer{conn: conn, ctrl: ctrl}
	buf.pending = make(chan []byte, txqueuelen)
	buf.cache = make([]byte, packet.MaximumPacketSize)
	return &buf
}
