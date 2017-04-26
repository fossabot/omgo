// Code generated by protoc-gen-go.
// source: common.proto
// DO NOT EDIT!

/*
Package proto_common is a generated protocol buffer package.

It is generated from these files:
	common.proto

It has these top-level messages:
	RspHeader
	C2SGetSeedReq
	S2CGetSeedRsp
*/
package proto_common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// commands
// S -> C : even
// C -> S : odd
type Cmd int32

const (
	Cmd_RESERVED       Cmd = 0
	Cmd_HEART_BEAT_REQ Cmd = 1
	Cmd_HEART_BEAT_RSP Cmd = 2
	Cmd_LOGIN_REQ      Cmd = 17
	Cmd_LOGIN_RSP      Cmd = 18
	Cmd_GET_SEED_REQ   Cmd = 49
	Cmd_GET_SEED_RSP   Cmd = 50
	Cmd_PING_REQ       Cmd = 4097
	Cmd_PING_RSP       Cmd = 4098
	Cmd_CMD_COMMON_END Cmd = 65536
)

var Cmd_name = map[int32]string{
	0:     "RESERVED",
	1:     "HEART_BEAT_REQ",
	2:     "HEART_BEAT_RSP",
	17:    "LOGIN_REQ",
	18:    "LOGIN_RSP",
	49:    "GET_SEED_REQ",
	50:    "GET_SEED_RSP",
	4097:  "PING_REQ",
	4098:  "PING_RSP",
	65536: "CMD_COMMON_END",
}
var Cmd_value = map[string]int32{
	"RESERVED":       0,
	"HEART_BEAT_REQ": 1,
	"HEART_BEAT_RSP": 2,
	"LOGIN_REQ":      17,
	"LOGIN_RSP":      18,
	"GET_SEED_REQ":   49,
	"GET_SEED_RSP":   50,
	"PING_REQ":       4097,
	"PING_RSP":       4098,
	"CMD_COMMON_END": 65536,
}

func (x Cmd) String() string {
	return proto.EnumName(Cmd_name, int32(x))
}
func (Cmd) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ResultCode int32

const (
	ResultCode_RESULT_OK             ResultCode = 0
	ResultCode_RESULT_TIMEOUT        ResultCode = 1
	ResultCode_RESULT_UNAUTH         ResultCode = 2
	ResultCode_RESULT_UNKNOWN_CMD    ResultCode = 3
	ResultCode_RESULT_INTERNAL_ERROR ResultCode = 4
)

var ResultCode_name = map[int32]string{
	0: "RESULT_OK",
	1: "RESULT_TIMEOUT",
	2: "RESULT_UNAUTH",
	3: "RESULT_UNKNOWN_CMD",
	4: "RESULT_INTERNAL_ERROR",
}
var ResultCode_value = map[string]int32{
	"RESULT_OK":             0,
	"RESULT_TIMEOUT":        1,
	"RESULT_UNAUTH":         2,
	"RESULT_UNKNOWN_CMD":    3,
	"RESULT_INTERNAL_ERROR": 4,
}

func (x ResultCode) String() string {
	return proto.EnumName(ResultCode_name, int32(x))
}
func (ResultCode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type RspHeader struct {
	Status    ResultCode `protobuf:"varint,1,opt,name=status,enum=proto.common.ResultCode" json:"status,omitempty"`
	Timestamp uint64     `protobuf:"fixed64,2,opt,name=timestamp" json:"timestamp,omitempty"`
	Msg       string     `protobuf:"bytes,3,opt,name=msg" json:"msg,omitempty"`
}

func (m *RspHeader) Reset()                    { *m = RspHeader{} }
func (m *RspHeader) String() string            { return proto.CompactTextString(m) }
func (*RspHeader) ProtoMessage()               {}
func (*RspHeader) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RspHeader) GetStatus() ResultCode {
	if m != nil {
		return m.Status
	}
	return ResultCode_RESULT_OK
}

func (m *RspHeader) GetTimestamp() uint64 {
	if m != nil {
		return m.Timestamp
	}
	return 0
}

func (m *RspHeader) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type C2SGetSeedReq struct {
	SendSeed []byte `protobuf:"bytes,1,opt,name=send_seed,json=sendSeed,proto3" json:"send_seed,omitempty"`
	RecvSeed []byte `protobuf:"bytes,2,opt,name=recv_seed,json=recvSeed,proto3" json:"recv_seed,omitempty"`
}

func (m *C2SGetSeedReq) Reset()                    { *m = C2SGetSeedReq{} }
func (m *C2SGetSeedReq) String() string            { return proto.CompactTextString(m) }
func (*C2SGetSeedReq) ProtoMessage()               {}
func (*C2SGetSeedReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *C2SGetSeedReq) GetSendSeed() []byte {
	if m != nil {
		return m.SendSeed
	}
	return nil
}

func (m *C2SGetSeedReq) GetRecvSeed() []byte {
	if m != nil {
		return m.RecvSeed
	}
	return nil
}

type S2CGetSeedRsp struct {
	Header   *RspHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	SendSeed []byte     `protobuf:"bytes,2,opt,name=send_seed,json=sendSeed,proto3" json:"send_seed,omitempty"`
	RecvSeed []byte     `protobuf:"bytes,3,opt,name=recv_seed,json=recvSeed,proto3" json:"recv_seed,omitempty"`
}

func (m *S2CGetSeedRsp) Reset()                    { *m = S2CGetSeedRsp{} }
func (m *S2CGetSeedRsp) String() string            { return proto.CompactTextString(m) }
func (*S2CGetSeedRsp) ProtoMessage()               {}
func (*S2CGetSeedRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *S2CGetSeedRsp) GetHeader() *RspHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *S2CGetSeedRsp) GetSendSeed() []byte {
	if m != nil {
		return m.SendSeed
	}
	return nil
}

func (m *S2CGetSeedRsp) GetRecvSeed() []byte {
	if m != nil {
		return m.RecvSeed
	}
	return nil
}

func init() {
	proto.RegisterType((*RspHeader)(nil), "proto.common.RspHeader")
	proto.RegisterType((*C2SGetSeedReq)(nil), "proto.common.C2SGetSeedReq")
	proto.RegisterType((*S2CGetSeedRsp)(nil), "proto.common.S2CGetSeedRsp")
	proto.RegisterEnum("proto.common.Cmd", Cmd_name, Cmd_value)
	proto.RegisterEnum("proto.common.ResultCode", ResultCode_name, ResultCode_value)
}

func init() { proto.RegisterFile("common.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 406 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x91, 0xcf, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0xb1, 0x5d, 0x45, 0xf1, 0x60, 0x47, 0xdb, 0x11, 0x7f, 0x82, 0xe0, 0x60, 0xe5, 0x14,
	0xf5, 0x10, 0xc0, 0x3c, 0x41, 0xb0, 0x57, 0x89, 0xd5, 0x78, 0x6d, 0x66, 0x1d, 0x38, 0xae, 0x42,
	0xbd, 0x02, 0x24, 0x5c, 0x87, 0xac, 0xc3, 0x81, 0x53, 0xe1, 0x8d, 0x78, 0x43, 0x64, 0xc7, 0x25,
	0x4d, 0x0f, 0x3d, 0xed, 0xce, 0xef, 0xfb, 0x34, 0xdf, 0x8c, 0x06, 0xbc, 0xab, 0xba, 0xaa, 0xea,
	0xeb, 0xd9, 0x76, 0x57, 0x37, 0x35, 0x7a, 0xdd, 0x33, 0x3b, 0xb0, 0x49, 0x05, 0x2e, 0x99, 0xed,
	0x52, 0x6f, 0x4a, 0xbd, 0xc3, 0x37, 0x30, 0x30, 0xcd, 0xa6, 0xd9, 0x9b, 0xb1, 0x15, 0x58, 0xd3,
	0x51, 0x38, 0x9e, 0xdd, 0xf5, 0xce, 0x48, 0x9b, 0xfd, 0xf7, 0x26, 0xaa, 0x4b, 0x4d, 0xbd, 0x0f,
	0x5f, 0x81, 0xdb, 0x7c, 0xab, 0xb4, 0x69, 0x36, 0xd5, 0x76, 0x6c, 0x07, 0xd6, 0x74, 0x40, 0x47,
	0x80, 0x0c, 0x9c, 0xca, 0x7c, 0x19, 0x3b, 0x81, 0x35, 0x75, 0xa9, 0xfd, 0x4e, 0x12, 0xf0, 0xa3,
	0x50, 0x2e, 0x74, 0x23, 0xb5, 0x2e, 0x49, 0xff, 0xc0, 0x97, 0xe0, 0x1a, 0x7d, 0x5d, 0x2a, 0xa3,
	0x75, 0xd9, 0xa5, 0x7a, 0x34, 0x6c, 0x41, 0xab, 0xb7, 0xe2, 0x4e, 0x5f, 0xfd, 0x3c, 0x88, 0xf6,
	0x41, 0x6c, 0x41, 0x2b, 0x4e, 0x7e, 0x81, 0x2f, 0xc3, 0xe8, 0xb6, 0x95, 0xd9, 0xe2, 0x6b, 0x18,
	0x7c, 0xed, 0xf6, 0xe8, 0xfa, 0x3c, 0x0e, 0x9f, 0xdf, 0x9b, 0xfe, 0x76, 0x4d, 0xea, 0x6d, 0xa7,
	0xd9, 0xf6, 0x43, 0xd9, 0xce, 0x69, 0xf6, 0xc5, 0x5f, 0x0b, 0x9c, 0xa8, 0x2a, 0xd1, 0x83, 0x21,
	0x71, 0xc9, 0xe9, 0x23, 0x8f, 0xd9, 0x23, 0x44, 0x18, 0x2d, 0xf9, 0x9c, 0x0a, 0xf5, 0x9e, 0xcf,
	0x0b, 0x45, 0xfc, 0x03, 0xb3, 0xee, 0x33, 0x99, 0x33, 0x1b, 0x7d, 0x70, 0x57, 0xd9, 0x22, 0x11,
	0x9d, 0xe5, 0xfc, 0x4e, 0x29, 0x73, 0x86, 0xc8, 0xc0, 0x5b, 0xf0, 0x42, 0x49, 0xce, 0xe3, 0xce,
	0xf0, 0xf6, 0x94, 0xc8, 0x9c, 0x85, 0xe8, 0xc3, 0x30, 0x4f, 0xc4, 0xa2, 0xd3, 0x7f, 0x07, 0xc7,
	0x52, 0xe6, 0xec, 0x4f, 0x80, 0x4f, 0x60, 0x14, 0xa5, 0xb1, 0x8a, 0xb2, 0x34, 0xcd, 0x84, 0xe2,
	0x22, 0x66, 0x37, 0x37, 0x67, 0x17, 0x7b, 0x80, 0xe3, 0x01, 0xdb, 0x50, 0xe2, 0x72, 0xbd, 0x2a,
	0x54, 0x76, 0x79, 0x18, 0xbd, 0x2f, 0x8b, 0x24, 0xe5, 0xd9, 0xba, 0x60, 0x16, 0x9e, 0x83, 0xdf,
	0xb3, 0xb5, 0x98, 0xaf, 0x8b, 0x25, 0xb3, 0xf1, 0x19, 0xe0, 0x7f, 0x74, 0x29, 0xb2, 0x4f, 0x42,
	0x45, 0x69, 0xcc, 0x1c, 0x7c, 0x01, 0x4f, 0x7b, 0x9e, 0x88, 0x82, 0x93, 0x98, 0xaf, 0x14, 0x27,
	0xca, 0x88, 0x9d, 0x7d, 0x1e, 0x74, 0x47, 0x78, 0xf7, 0x2f, 0x00, 0x00, 0xff, 0xff, 0x32, 0x95,
	0xfa, 0x15, 0x85, 0x02, 0x00, 0x00,
}
