package packet

// Packet wraps a byte slice as a binary buffer
type Packet struct {
	buf []byte
	pos int
}

// Data returns packet's internal buffer
func (p *Packet) Data() []byte {
	return p.buf
}

// Len returns length of the packet's internal buffer
func (p *Packet) Len() int {
	return len(p.buf)
}

// ----------------------------------------------------------------------------
// Read
// ----------------------------------------------------------------------------
