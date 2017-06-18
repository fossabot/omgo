package main

import (
	"net"

	"encoding/binary"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/net/packet"
	"github.com/master-g/omgo/utils"
)

type Buffer struct {
	ctrl    chan struct{}
	pending chan []byte
	conn    net.Conn
	cache   []byte
}

func (buf *Buffer) send(session *Session, data []byte) {
	if data == nil {
		return
	}

	if session.IsFlagEncryptedSet() {
		session.Encoder.XORKeyStream(data, data)
	} else if session.IsFlagKeyExchangedSet() {
		session.ClearFlagKeyExchanged()
		session.SetFlagEncrypted()
	}

	select {
	case buf.pending <- data:
	default:
		log.Warning("pending full")
	}
	return
}

func (buf *Buffer) start() {
	defer utils.PrintPanicStack()
	for {
		select {
		case data := <-buf.pending:
			buf.rawSend(data)
		case <-buf.ctrl:
			return
		}
	}
}

func (buf *Buffer) rawSend(data []byte) bool {
	size := len(data)
	binary.BigEndian.PutUint16(buf.cache, uint16(size))
	copy(buf.cache[2:], data)

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
