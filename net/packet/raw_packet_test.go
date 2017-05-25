package packet

import (
	"testing"
)

func TestRawPacketWriter(t *testing.T) {
	p := NewRawPacket()
	tByte := byte(0xFF)
	tU16 := uint16(0xFF00)
	tU32 := uint32(0xFF0000)
	tU64 := uint64(0xFF000000000000)
	tS8 := int8(-1)

	p.WriteBool(true)
	p.WriteBool(false)
	p.WriteByte(tByte)
	p.WriteU16(tU16)
	p.WriteU32(tU32)
	p.WriteU64(tU64)
	p.WriteS8(tS8)

	str := "hello world!"

	p.WriteString(str)
	p.WriteBytes([]byte(str))
	var nilBytes []byte
	p.WriteBytes(nilBytes)

	reader := NewRawPacketReader(p.Data())

	boolean, _ := reader.ReadBool()
	if boolean != true {
		t.Error("packet readbool mismatch")
	}
	boolean, _ = reader.ReadBool()
	if boolean != false {
		t.Error("packet readbool mismatch")
	}

	rByte, _ := reader.ReadByte()
	if rByte != tByte {
		t.Error("packet readByte mismatch")
	}

	rU16, _ := reader.ReadU16()
	if rU16 != tU16 {
		t.Error("packet readU16 mismatch")
	}

	rU32, _ := reader.ReadU32()
	if rU32 != tU32 {
		t.Error("packet readU32 mismatch")
	}

	rU64, _ := reader.ReadU64()
	if rU64 != tU64 {
		t.Error("packet readU64 mismatch")
	}

	rS8, _ := reader.ReadS8()
	if rS8 != tS8 {
		t.Error("packet readS8 mismatch")
	}

	rStr, _ := reader.ReadString()
	if rStr != str {
		t.Error("packet read string mismatch")
	}

	rBytes, _ := reader.ReadBytes()
	if rBytes[0] != str[0] {
		t.Error("packet read bytes mismatch")
	}

	rNil, _ := reader.ReadBytes()
	if len(rNil) != 0 {
		t.Error("packet read nil bytes mismatch")
	}

	_, err := reader.ReadByte()
	if err == nil {
		t.Error("overflow check failed")
	}
}

func BenchmarkNewRawPacketWriter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := NewRawPacket()
		p.WriteU16(128)
		p.WriteBool(true)
		p.WriteS32(-16)
		p.WriteString("A")
		p.WriteU32(16)
		p.WriteBytes([]byte{1, 2, 3})
	}
}
