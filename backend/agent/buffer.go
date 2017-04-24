package main

import (
	"encoding/binary"
	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/utils"
	"net"
    "github.com/master-g/omgo/net/packet"
)

// PIPELINE #3: buffer
// controls the packet send to clients
type Buffer struct {
	ctrl    chan struct{} // receive exit signal
	pending chan []byte   // pending packets
	conn    net.Conn      // connection
	cache   []byte        // for combined syscall write
}

// packet sending procedure
func (buf *Buffer) send(session *Session, data []byte) {
	// in case of empty packet
	if data == nil {
		return
	}

	// encryption
	// (NOT_ENCRYPTED) -> KEYEXCG -> ENCRYPTED
	if session.Flag&SESS_ENCRYPT != 0 {
		// encryption is enabled
		session.Encoder.XORKeyStream(data, data)
	} else if session.Flag&SESS_KEYEXCG != 0 {
		// key is exchanged, encryption is not yet enabled
		session.Flag &^= SESS_KEYEXCG
		session.Flag |= SESS_ENCRYPT
	}

	// queue the data for sending
	select {
	case buf.pending <- data:
	default:
		// packet will be dropped if it exceeds txQueueLength
		log.WithFields(log.Fields{"userid": session.UserID, "ip": session.IP}).Warning("pending full")
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
		log.Warningf("Error send reply data, bytes:%v reason:%v", n, err)
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
