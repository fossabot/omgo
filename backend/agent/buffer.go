// PIPELINE #3: buffer

package main

import (
	"encoding/binary"
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/master-g/omgo/backend/agent/api"
	"github.com/master-g/omgo/kit/utils"
)

// Buffer manages sending packets to client
type Buffer struct {
	ctrl    chan struct{} // receive exit signal
	pending chan []byte   // pending packets
	conn    net.Conn      // connection
	cache   []byte        // for combined syscall write
}

// send data into buffer's channel
func (buf *Buffer) send(session *api.Session, pkg *api.OutgoingPacket) {
	// in case of empty packet
	if pkg == nil {
		return
	}

	pkg.Header.

	// encryption
	// (NOT_ENCRYPTED) -> KEYEXCG -> ENCRYPTED
	if session.IsFlagEncryptedSet() {
		// encryption is enabled
		session.Encoder.XORKeyStream(data, data)
	} else if session.IsFlagKeyExchangedSet() {
		// key is exchanged, encryption is not yet establish
		session.ClearFlagKeyExchanged()
		session.SetFlagEncrypted()
	}

	// queue the data for sending
	select {
	case buf.pending <- data:
	default:
		// packet will be dropped if it exceeds txQueueLength
		log.WithFields(log.Fields{"usn": session.Usn, "ip": session.IP}).Warning("pending full")
	}
	return
}

// packet sending goroutine
func (buf *Buffer) start() {
	defer utils.PrintPanicStack()
	for {
		select {
		case data := <-buf.pending:
			// dequeue data to send
			buf.rawSend(data)
		case <-buf.ctrl:
			// session control signal received
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

func newBuffer(conn net.Conn, ctrl chan struct{}, txQueueLen int) *Buffer {
	buf := Buffer{conn: conn, ctrl: ctrl}
	buf.pending = make(chan []byte, txQueueLen)
	buf.cache = make([]byte, 32*1024)
	return &buf
}
