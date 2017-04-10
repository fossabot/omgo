package packet

type Reader interface {
	ReadBool() (bool, error)
	ReadByte() (byte, error)
	ReadBytes() ([]byte, error)
	ReadString() (string, error)
	ReadS8() (int8, error)
	ReadU16() (uint16, error)
	ReadS16() (int16, error)
	ReadU32() (uint32, error)
	ReadS32() (int32, error)
	ReadU64() (uint64, error)
	ReadS64() (int64, error)
}

type Writer interface {
	WriteZeros(n int)
	WriteBool(v bool)
	WriteByte(v byte)
	WriteBytes(v []byte)
	WriteRawBytes(v []byte)
	WriteString(v string)
	WriteS8(v int8)
	WriteU16(v uint16)
	WriteS16(v int16)
	WriteU32(v uint32)
	WriteS32(v int32)
	WriteU64(v uint64)
	WriteS64(v int64)
}

type Packet interface {
	Data() []byte
	Len() int
	Reader
	Writer
}
