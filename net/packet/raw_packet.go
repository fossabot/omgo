package packet

import (
	"encoding/binary"
	"errors"
)

type RawPacket struct {
	buf []byte
	pos int
}

// NewRawPacket creates a empty RawPacket for writing
func NewRawPacket() *RawPacket {
	return &RawPacket{make([]byte, 0), 0}
}

// NewRawPacketReader creates a new RawPacket from bytes for reading
func NewRawPacketReader(buf []byte) *RawPacket {
	return &RawPacket{buf, 0}
}

// Data returns packet's internal buffer
func (p *RawPacket) Data() []byte {
	return p.buf
}

// Len returns length of the packet's internal buffer
func (p *RawPacket) Len() int {
	return len(p.buf)
}

func (p *RawPacket) Pos() int {
	return p.pos
}

// ----------------------------------------------------------------------------
// Read
// ----------------------------------------------------------------------------

func (p *RawPacket) ReadBool() (ret bool, err error) {
	b, err := p.ReadByte()
	if b != byte(0) {
		return true, err
	}

	return false, err
}

func (p *RawPacket) ReadByte() (ret byte, err error) {
	if p.pos >= len(p.buf) {
		err = errors.New("read byte failed")
		return
	}

	ret = p.buf[p.pos]
	p.pos++
	return
}

func (p *RawPacket) ReadBytes() (ret []byte, err error) {
	if p.pos+2 > len(p.buf) {
		err = errors.New("read bytes header failed")
		return
	}
	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.buf) {
		err = errors.New("read bytes data failed")
		return
	}

	ret = p.buf[p.pos : p.pos+int(size)]
	p.pos += int(size)
	return
}

func (p *RawPacket) ReadString() (ret string, err error) {
	if p.pos+2 > len(p.buf) {
		err = errors.New("read string length failed")
		return
	}

	size, _ := p.ReadU16()
	if p.pos+int(size) > len(p.buf) {
		err = errors.New("read string data failed")
		return
	}

	bytes := p.buf[p.pos : p.pos+int(size)]
	p.pos += int(size)
	ret = string(bytes)
	return
}

func (p *RawPacket) ReadS8() (ret int8, err error) {
	_ret, _err := p.ReadByte()
	ret = int8(_ret)
	err = _err
	return
}

func (p *RawPacket) ReadU16() (ret uint16, err error) {
	if p.pos+2 > len(p.buf) {
		err = errors.New("read uint16 failed")
		return
	}

	buf := p.buf[p.pos : p.pos+2]
	ret = binary.LittleEndian.Uint16(buf)
	p.pos += 2
	return
}

func (p *RawPacket) ReadS16() (ret int16, err error) {
	_ret, _err := p.ReadU16()
	ret = int16(_ret)
	err = _err
	return
}

func (p *RawPacket) ReadU32() (ret uint32, err error) {
	if p.pos+4 > len(p.buf) {
		err = errors.New("read uint32 failed")
		return
	}

	buf := p.buf[p.pos : p.pos+4]
	ret = binary.LittleEndian.Uint32(buf)
	p.pos += 4
	return
}

func (p *RawPacket) ReadS32() (ret int32, err error) {
	_ret, _err := p.ReadU32()
	ret = int32(_ret)
	err = _err
	return
}

func (p *RawPacket) ReadU64() (ret uint64, err error) {
	if p.pos+8 > len(p.buf) {
		err = errors.New("read uint64 failed")
		return
	}

	buf := p.buf[p.pos : p.pos+8]
	ret = binary.LittleEndian.Uint64(buf)
	p.pos += 8
	return
}

func (p *RawPacket) ReadS64() (ret int64, err error) {
	_ret, _err := p.ReadU64()
	ret = int64(_ret)
	err = _err
	return
}

// ----------------------------------------------------------------------------
// Write
// ----------------------------------------------------------------------------

func (p *RawPacket) WriteZeros(n int) *RawPacket {
	v := make([]byte, n)
	p.buf = append(p.buf, v...)
	return p
}

func (p *RawPacket) WriteBool(v bool) *RawPacket {
	if v {
		p.buf = append(p.buf, byte(1))
	} else {
		p.buf = append(p.buf, byte(0))
	}
	return p
}

func (p *RawPacket) WriteByte(v byte) *RawPacket {
	p.buf = append(p.buf, v)
	return p
}

func (p *RawPacket) WriteBytes(v []byte) *RawPacket {
	p.WriteU16(uint16(len(v)))
	p.buf = append(p.buf, v...)
	return p
}

func (p *RawPacket) WriteRawBytes(v []byte) *RawPacket {
	p.buf = append(p.buf, v...)
	return p
}

func (p *RawPacket) WriteString(v string) *RawPacket {
	bytes := []byte(v)
	p.WriteU16(uint16(len(bytes)))
	p.buf = append(p.buf, bytes...)
	return p
}

func (p *RawPacket) WriteS8(v int8) *RawPacket {
	p.WriteByte(byte(v))
	return p
}

func (p *RawPacket) WriteU16(v uint16) *RawPacket {
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes, v)
	p.buf = append(p.buf, bytes...)
	return p
}

func (p *RawPacket) WriteS16(v int16) *RawPacket {
	p.WriteU16(uint16(v))
	return p
}

func (p *RawPacket) WriteU32(v uint32) *RawPacket {
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, v)
	p.buf = append(p.buf, bytes...)
	return p
}

func (p *RawPacket) WriteS32(v int32) *RawPacket {
	p.WriteU32(uint32(v))
	return p
}

func (p *RawPacket) WriteU64(v uint64) *RawPacket {
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, v)
	p.buf = append(p.buf, bytes...)
	return p
}

func (p *RawPacket) WriteS64(v int64) *RawPacket {
	p.WriteU64(uint64(v))
	return p
}
