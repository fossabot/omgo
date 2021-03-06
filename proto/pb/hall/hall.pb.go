// Code generated by protoc-gen-go. DO NOT EDIT.
// source: hall.proto

/*
Package hall is a generated protocol buffer package.

It is generated from these files:
	hall.proto

It has these top-level messages:
	RoomInfo
	RoomConfigElement
	C2SGetRoomConfigReq
	S2CGetRoomConfigRsp
	C2SEnterRoomReq
	S2CEnterRoomRsp
	C2SCreateRoomReq
	S2CCreateRoomRsp
*/
package hall

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import proto_common "github.com/master-g/omgo/proto/pb/common"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type RoomType int32

const (
	RoomType_ROOM_RESERVED RoomType = 0
	RoomType_ROOM_NORMAL   RoomType = 1
	RoomType_ROOM_PRIVATE  RoomType = 2
)

var RoomType_name = map[int32]string{
	0: "ROOM_RESERVED",
	1: "ROOM_NORMAL",
	2: "ROOM_PRIVATE",
}
var RoomType_value = map[string]int32{
	"ROOM_RESERVED": 0,
	"ROOM_NORMAL":   1,
	"ROOM_PRIVATE":  2,
}

func (x RoomType) String() string {
	return proto.EnumName(RoomType_name, int32(x))
}
func (RoomType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type BillingMode int32

const (
	BillingMode_BILLING_RESERVED BillingMode = 0
	BillingMode_BILLING_TIME     BillingMode = 1
	BillingMode_BILLING_ROUND    BillingMode = 2
)

var BillingMode_name = map[int32]string{
	0: "BILLING_RESERVED",
	1: "BILLING_TIME",
	2: "BILLING_ROUND",
}
var BillingMode_value = map[string]int32{
	"BILLING_RESERVED": 0,
	"BILLING_TIME":     1,
	"BILLING_ROUND":    2,
}

func (x BillingMode) String() string {
	return proto.EnumName(BillingMode_name, int32(x))
}
func (BillingMode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type GameMode int32

const (
	GameMode_GAME_MODE_RESERVED GameMode = 0
	GameMode_GAME_MODE_RACE     GameMode = 1
	GameMode_GAME_MODE_CALL     GameMode = 2
)

var GameMode_name = map[int32]string{
	0: "GAME_MODE_RESERVED",
	1: "GAME_MODE_RACE",
	2: "GAME_MODE_CALL",
}
var GameMode_value = map[string]int32{
	"GAME_MODE_RESERVED": 0,
	"GAME_MODE_RACE":     1,
	"GAME_MODE_CALL":     2,
}

func (x GameMode) String() string {
	return proto.EnumName(GameMode_name, int32(x))
}
func (GameMode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type EnterRoomResult int32

const (
	EnterRoomResult_ENTER_ROOM_OK                  EnterRoomResult = 0
	EnterRoomResult_ENTER_ROOM_INSUFFICIENT_CREDIT EnterRoomResult = 1
	EnterRoomResult_ENTER_ROOM_WRONG_PASSWORD      EnterRoomResult = 2
	EnterRoomResult_ENTER_ROOM_NOT_EXIST           EnterRoomResult = 3
	EnterRoomResult_ENTER_ROOM_FULL                EnterRoomResult = 4
	EnterRoomResult_ENTER_ROOM_FAIL_NO_OBSERVE     EnterRoomResult = 5
)

var EnterRoomResult_name = map[int32]string{
	0: "ENTER_ROOM_OK",
	1: "ENTER_ROOM_INSUFFICIENT_CREDIT",
	2: "ENTER_ROOM_WRONG_PASSWORD",
	3: "ENTER_ROOM_NOT_EXIST",
	4: "ENTER_ROOM_FULL",
	5: "ENTER_ROOM_FAIL_NO_OBSERVE",
}
var EnterRoomResult_value = map[string]int32{
	"ENTER_ROOM_OK":                  0,
	"ENTER_ROOM_INSUFFICIENT_CREDIT": 1,
	"ENTER_ROOM_WRONG_PASSWORD":      2,
	"ENTER_ROOM_NOT_EXIST":           3,
	"ENTER_ROOM_FULL":                4,
	"ENTER_ROOM_FAIL_NO_OBSERVE":     5,
}

func (x EnterRoomResult) String() string {
	return proto.EnumName(EnterRoomResult_name, int32(x))
}
func (EnterRoomResult) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type CreateRoomResult int32

const (
	CreateRoomResult_CREATE_ROOM_OK                CreateRoomResult = 0
	CreateRoomResult_CREATE_ROOM_INVALID_PARAM     CreateRoomResult = 1
	CreateRoomResult_CREATE_ROOM_INSUFFICIENT_CARD CreateRoomResult = 2
	CreateRoomResult_CREATE_ROOM_ALREADY_EXIST     CreateRoomResult = 3
)

var CreateRoomResult_name = map[int32]string{
	0: "CREATE_ROOM_OK",
	1: "CREATE_ROOM_INVALID_PARAM",
	2: "CREATE_ROOM_INSUFFICIENT_CARD",
	3: "CREATE_ROOM_ALREADY_EXIST",
}
var CreateRoomResult_value = map[string]int32{
	"CREATE_ROOM_OK":                0,
	"CREATE_ROOM_INVALID_PARAM":     1,
	"CREATE_ROOM_INSUFFICIENT_CARD": 2,
	"CREATE_ROOM_ALREADY_EXIST":     3,
}

func (x CreateRoomResult) String() string {
	return proto.EnumName(CreateRoomResult_name, int32(x))
}
func (CreateRoomResult) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type RoomInfo struct {
	Id            uint64                 `protobuf:"fixed64,1,opt,name=id" json:"id,omitempty"`
	Owner         uint64                 `protobuf:"fixed64,2,opt,name=owner" json:"owner,omitempty"`
	RoomType      RoomType               `protobuf:"varint,3,opt,name=room_type,json=roomType,enum=proto.hall.RoomType" json:"room_type"`
	GameMode      GameMode               `protobuf:"varint,4,opt,name=game_mode,json=gameMode,enum=proto.hall.GameMode" json:"game_mode"`
	BillMode      BillingMode            `protobuf:"varint,5,opt,name=bill_mode,json=billMode,enum=proto.hall.BillingMode" json:"bill_mode"`
	Title         string                 `protobuf:"bytes,6,opt,name=title" json:"title,omitempty"`
	Desc          string                 `protobuf:"bytes,7,opt,name=desc" json:"desc,omitempty"`
	Since         uint64                 `protobuf:"fixed64,8,opt,name=since" json:"since,omitempty"`
	Observable    bool                   `protobuf:"varint,9,opt,name=observable" json:"observable"`
	BoomLimit     int32                  `protobuf:"varint,10,opt,name=boom_limit,json=boomLimit" json:"boom_limit"`
	BasePoint     int32                  `protobuf:"varint,11,opt,name=base_point,json=basePoint" json:"base_point"`
	TotalDuration int32                  `protobuf:"varint,12,opt,name=total_duration,json=totalDuration" json:"total_duration"`
	LeftDuration  int32                  `protobuf:"varint,13,opt,name=left_duration,json=leftDuration" json:"left_duration"`
	TotalRounds   int32                  `protobuf:"varint,14,opt,name=total_rounds,json=totalRounds" json:"total_rounds"`
	LeftRounds    int32                  `protobuf:"varint,15,opt,name=left_rounds,json=leftRounds" json:"left_rounds"`
	CallTime      int32                  `protobuf:"varint,16,opt,name=call_time,json=callTime" json:"call_time"`
	HandTime      int32                  `protobuf:"varint,17,opt,name=hand_time,json=handTime" json:"hand_time"`
	Players       []uint64               `protobuf:"fixed64,18,rep,packed,name=players" json:"players,omitempty"`
	Secure        bool                   `protobuf:"varint,19,opt,name=secure" json:"secure"`
	Location      *proto_common.Location `protobuf:"bytes,20,opt,name=location" json:"location,omitempty"`
}

func (m *RoomInfo) Reset()                    { *m = RoomInfo{} }
func (m *RoomInfo) String() string            { return proto.CompactTextString(m) }
func (*RoomInfo) ProtoMessage()               {}
func (*RoomInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RoomInfo) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *RoomInfo) GetOwner() uint64 {
	if m != nil {
		return m.Owner
	}
	return 0
}

func (m *RoomInfo) GetRoomType() RoomType {
	if m != nil {
		return m.RoomType
	}
	return RoomType_ROOM_RESERVED
}

func (m *RoomInfo) GetGameMode() GameMode {
	if m != nil {
		return m.GameMode
	}
	return GameMode_GAME_MODE_RESERVED
}

func (m *RoomInfo) GetBillMode() BillingMode {
	if m != nil {
		return m.BillMode
	}
	return BillingMode_BILLING_RESERVED
}

func (m *RoomInfo) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *RoomInfo) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

func (m *RoomInfo) GetSince() uint64 {
	if m != nil {
		return m.Since
	}
	return 0
}

func (m *RoomInfo) GetObservable() bool {
	if m != nil {
		return m.Observable
	}
	return false
}

func (m *RoomInfo) GetBoomLimit() int32 {
	if m != nil {
		return m.BoomLimit
	}
	return 0
}

func (m *RoomInfo) GetBasePoint() int32 {
	if m != nil {
		return m.BasePoint
	}
	return 0
}

func (m *RoomInfo) GetTotalDuration() int32 {
	if m != nil {
		return m.TotalDuration
	}
	return 0
}

func (m *RoomInfo) GetLeftDuration() int32 {
	if m != nil {
		return m.LeftDuration
	}
	return 0
}

func (m *RoomInfo) GetTotalRounds() int32 {
	if m != nil {
		return m.TotalRounds
	}
	return 0
}

func (m *RoomInfo) GetLeftRounds() int32 {
	if m != nil {
		return m.LeftRounds
	}
	return 0
}

func (m *RoomInfo) GetCallTime() int32 {
	if m != nil {
		return m.CallTime
	}
	return 0
}

func (m *RoomInfo) GetHandTime() int32 {
	if m != nil {
		return m.HandTime
	}
	return 0
}

func (m *RoomInfo) GetPlayers() []uint64 {
	if m != nil {
		return m.Players
	}
	return nil
}

func (m *RoomInfo) GetSecure() bool {
	if m != nil {
		return m.Secure
	}
	return false
}

func (m *RoomInfo) GetLocation() *proto_common.Location {
	if m != nil {
		return m.Location
	}
	return nil
}

type RoomConfigElement struct {
	Id        int32 `protobuf:"varint,1,opt,name=id" json:"id"`
	CardCount int32 `protobuf:"varint,2,opt,name=card_count,json=cardCount" json:"card_count"`
	Duration  int32 `protobuf:"varint,3,opt,name=duration" json:"duration"`
	Rounds    int32 `protobuf:"varint,4,opt,name=rounds" json:"rounds"`
}

func (m *RoomConfigElement) Reset()                    { *m = RoomConfigElement{} }
func (m *RoomConfigElement) String() string            { return proto.CompactTextString(m) }
func (*RoomConfigElement) ProtoMessage()               {}
func (*RoomConfigElement) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *RoomConfigElement) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *RoomConfigElement) GetCardCount() int32 {
	if m != nil {
		return m.CardCount
	}
	return 0
}

func (m *RoomConfigElement) GetDuration() int32 {
	if m != nil {
		return m.Duration
	}
	return 0
}

func (m *RoomConfigElement) GetRounds() int32 {
	if m != nil {
		return m.Rounds
	}
	return 0
}

type C2SGetRoomConfigReq struct {
}

func (m *C2SGetRoomConfigReq) Reset()                    { *m = C2SGetRoomConfigReq{} }
func (m *C2SGetRoomConfigReq) String() string            { return proto.CompactTextString(m) }
func (*C2SGetRoomConfigReq) ProtoMessage()               {}
func (*C2SGetRoomConfigReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type S2CGetRoomConfigRsp struct {
	Header   *proto_common.RspHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
	Elements []*RoomConfigElement    `protobuf:"bytes,2,rep,name=elements" json:"elements,omitempty"`
}

func (m *S2CGetRoomConfigRsp) Reset()                    { *m = S2CGetRoomConfigRsp{} }
func (m *S2CGetRoomConfigRsp) String() string            { return proto.CompactTextString(m) }
func (*S2CGetRoomConfigRsp) ProtoMessage()               {}
func (*S2CGetRoomConfigRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *S2CGetRoomConfigRsp) GetHeader() *proto_common.RspHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *S2CGetRoomConfigRsp) GetElements() []*RoomConfigElement {
	if m != nil {
		return m.Elements
	}
	return nil
}

type C2SEnterRoomReq struct {
	Type     RoomType `protobuf:"varint,1,opt,name=type,enum=proto.hall.RoomType" json:"type"`
	RoomId   uint64   `protobuf:"fixed64,2,opt,name=room_id,json=roomId" json:"room_id,omitempty"`
	Password string   `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
}

func (m *C2SEnterRoomReq) Reset()                    { *m = C2SEnterRoomReq{} }
func (m *C2SEnterRoomReq) String() string            { return proto.CompactTextString(m) }
func (*C2SEnterRoomReq) ProtoMessage()               {}
func (*C2SEnterRoomReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *C2SEnterRoomReq) GetType() RoomType {
	if m != nil {
		return m.Type
	}
	return RoomType_ROOM_RESERVED
}

func (m *C2SEnterRoomReq) GetRoomId() uint64 {
	if m != nil {
		return m.RoomId
	}
	return 0
}

func (m *C2SEnterRoomReq) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

type S2CEnterRoomRsp struct {
	Header *proto_common.RspHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
}

func (m *S2CEnterRoomRsp) Reset()                    { *m = S2CEnterRoomRsp{} }
func (m *S2CEnterRoomRsp) String() string            { return proto.CompactTextString(m) }
func (*S2CEnterRoomRsp) ProtoMessage()               {}
func (*S2CEnterRoomRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *S2CEnterRoomRsp) GetHeader() *proto_common.RspHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

type C2SCreateRoomReq struct {
	BillingMode BillingMode            `protobuf:"varint,1,opt,name=billing_mode,json=billingMode,enum=proto.hall.BillingMode" json:"billing_mode"`
	CardNum     int32                  `protobuf:"varint,2,opt,name=card_num,json=cardNum" json:"card_num"`
	Password    string                 `protobuf:"bytes,3,opt,name=password" json:"password,omitempty"`
	Title       string                 `protobuf:"bytes,4,opt,name=title" json:"title,omitempty"`
	Desc        string                 `protobuf:"bytes,5,opt,name=desc" json:"desc,omitempty"`
	Observable  bool                   `protobuf:"varint,6,opt,name=observable" json:"observable"`
	BoomLimit   int32                  `protobuf:"varint,7,opt,name=boom_limit,json=boomLimit" json:"boom_limit"`
	Secure      bool                   `protobuf:"varint,8,opt,name=secure" json:"secure"`
	Location    *proto_common.Location `protobuf:"bytes,9,opt,name=location" json:"location,omitempty"`
}

func (m *C2SCreateRoomReq) Reset()                    { *m = C2SCreateRoomReq{} }
func (m *C2SCreateRoomReq) String() string            { return proto.CompactTextString(m) }
func (*C2SCreateRoomReq) ProtoMessage()               {}
func (*C2SCreateRoomReq) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *C2SCreateRoomReq) GetBillingMode() BillingMode {
	if m != nil {
		return m.BillingMode
	}
	return BillingMode_BILLING_RESERVED
}

func (m *C2SCreateRoomReq) GetCardNum() int32 {
	if m != nil {
		return m.CardNum
	}
	return 0
}

func (m *C2SCreateRoomReq) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *C2SCreateRoomReq) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *C2SCreateRoomReq) GetDesc() string {
	if m != nil {
		return m.Desc
	}
	return ""
}

func (m *C2SCreateRoomReq) GetObservable() bool {
	if m != nil {
		return m.Observable
	}
	return false
}

func (m *C2SCreateRoomReq) GetBoomLimit() int32 {
	if m != nil {
		return m.BoomLimit
	}
	return 0
}

func (m *C2SCreateRoomReq) GetSecure() bool {
	if m != nil {
		return m.Secure
	}
	return false
}

func (m *C2SCreateRoomReq) GetLocation() *proto_common.Location {
	if m != nil {
		return m.Location
	}
	return nil
}

type S2CCreateRoomRsp struct {
	Header *proto_common.RspHeader `protobuf:"bytes,1,opt,name=header" json:"header,omitempty"`
}

func (m *S2CCreateRoomRsp) Reset()                    { *m = S2CCreateRoomRsp{} }
func (m *S2CCreateRoomRsp) String() string            { return proto.CompactTextString(m) }
func (*S2CCreateRoomRsp) ProtoMessage()               {}
func (*S2CCreateRoomRsp) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *S2CCreateRoomRsp) GetHeader() *proto_common.RspHeader {
	if m != nil {
		return m.Header
	}
	return nil
}

func init() {
	proto.RegisterType((*RoomInfo)(nil), "proto.hall.RoomInfo")
	proto.RegisterType((*RoomConfigElement)(nil), "proto.hall.RoomConfigElement")
	proto.RegisterType((*C2SGetRoomConfigReq)(nil), "proto.hall.C2SGetRoomConfigReq")
	proto.RegisterType((*S2CGetRoomConfigRsp)(nil), "proto.hall.S2CGetRoomConfigRsp")
	proto.RegisterType((*C2SEnterRoomReq)(nil), "proto.hall.C2SEnterRoomReq")
	proto.RegisterType((*S2CEnterRoomRsp)(nil), "proto.hall.S2CEnterRoomRsp")
	proto.RegisterType((*C2SCreateRoomReq)(nil), "proto.hall.C2SCreateRoomReq")
	proto.RegisterType((*S2CCreateRoomRsp)(nil), "proto.hall.S2CCreateRoomRsp")
	proto.RegisterEnum("proto.hall.RoomType", RoomType_name, RoomType_value)
	proto.RegisterEnum("proto.hall.BillingMode", BillingMode_name, BillingMode_value)
	proto.RegisterEnum("proto.hall.GameMode", GameMode_name, GameMode_value)
	proto.RegisterEnum("proto.hall.EnterRoomResult", EnterRoomResult_name, EnterRoomResult_value)
	proto.RegisterEnum("proto.hall.CreateRoomResult", CreateRoomResult_name, CreateRoomResult_value)
}

func init() { proto.RegisterFile("hall.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1023 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x41, 0x6f, 0xe3, 0x44,
	0x14, 0x5e, 0xa7, 0x49, 0xea, 0xbc, 0xa4, 0x89, 0x3b, 0x0d, 0x5d, 0xb7, 0x28, 0x25, 0x1b, 0x04,
	0x8a, 0x2a, 0xd1, 0x88, 0xc0, 0x05, 0x4e, 0xb8, 0x8e, 0xdb, 0x35, 0x38, 0x76, 0x35, 0x4e, 0xbb,
	0xc0, 0xc5, 0x72, 0xe2, 0x69, 0x6a, 0xc9, 0xf6, 0x78, 0x6d, 0x67, 0x57, 0x3d, 0x72, 0xe2, 0x3f,
	0xf0, 0x43, 0x38, 0xf0, 0xeb, 0xd0, 0x8c, 0x9d, 0xc4, 0x29, 0x50, 0x56, 0x7b, 0xb2, 0xdf, 0xf7,
	0x7d, 0x33, 0xf3, 0xde, 0xcc, 0xfb, 0x1e, 0xc0, 0x83, 0x1b, 0x04, 0x17, 0x71, 0x42, 0x33, 0x8a,
	0x80, 0x7f, 0x2e, 0x18, 0x72, 0xda, 0x5a, 0xd0, 0x30, 0xa4, 0x51, 0xce, 0x0c, 0xfe, 0xa8, 0x81,
	0x88, 0x29, 0x0d, 0xf5, 0xe8, 0x9e, 0xa2, 0x36, 0x54, 0x7c, 0x4f, 0x16, 0xfa, 0xc2, 0xb0, 0x8e,
	0x2b, 0xbe, 0x87, 0xba, 0x50, 0xa3, 0xef, 0x23, 0x92, 0xc8, 0x15, 0x0e, 0xe5, 0x01, 0xfa, 0x1a,
	0x1a, 0x09, 0xa5, 0xa1, 0x93, 0x3d, 0xc6, 0x44, 0xde, 0xeb, 0x0b, 0xc3, 0xf6, 0xb8, 0x7b, 0xb1,
	0x3d, 0xe0, 0x82, 0x6d, 0x37, 0x7b, 0x8c, 0x09, 0x16, 0x93, 0xe2, 0x8f, 0x2d, 0x59, 0xba, 0x21,
	0x71, 0x42, 0xea, 0x11, 0xb9, 0xfa, 0xcf, 0x25, 0xd7, 0x6e, 0x48, 0xa6, 0xd4, 0x23, 0x58, 0x5c,
	0x16, 0x7f, 0xe8, 0x5b, 0x68, 0xcc, 0xfd, 0x20, 0xc8, 0x97, 0xd4, 0xf8, 0x92, 0x97, 0xe5, 0x25,
	0x97, 0x7e, 0x10, 0xf8, 0xd1, 0x32, 0x5f, 0xc5, 0x94, 0x7c, 0x55, 0x17, 0x6a, 0x99, 0x9f, 0x05,
	0x44, 0xae, 0xf7, 0x85, 0x61, 0x03, 0xe7, 0x01, 0x42, 0x50, 0xf5, 0x48, 0xba, 0x90, 0xf7, 0x39,
	0xc8, 0xff, 0x99, 0x32, 0xf5, 0xa3, 0x05, 0x91, 0xc5, 0xbc, 0x36, 0x1e, 0xa0, 0x33, 0x00, 0x3a,
	0x4f, 0x49, 0xf2, 0xce, 0x9d, 0x07, 0x44, 0x6e, 0xf4, 0x85, 0xa1, 0x88, 0x4b, 0x08, 0xea, 0x01,
	0xcc, 0x59, 0xed, 0x81, 0x1f, 0xfa, 0x99, 0x0c, 0x7d, 0x61, 0x58, 0xc3, 0x0d, 0x86, 0x18, 0x0c,
	0xe0, 0xb4, 0x9b, 0x12, 0x27, 0xa6, 0x7e, 0x94, 0xc9, 0xcd, 0x82, 0x76, 0x53, 0x72, 0xc3, 0x00,
	0xf4, 0x05, 0xb4, 0x33, 0x9a, 0xb9, 0x81, 0xe3, 0xad, 0x12, 0x37, 0xf3, 0x69, 0x24, 0xb7, 0xb8,
	0xe4, 0x80, 0xa3, 0x93, 0x02, 0x44, 0x9f, 0xc3, 0x41, 0x40, 0xee, 0xb3, 0xad, 0xea, 0x80, 0xab,
	0x5a, 0x0c, 0xdc, 0x88, 0x5e, 0x41, 0x2b, 0xdf, 0x2b, 0xa1, 0xab, 0xc8, 0x4b, 0xe5, 0x36, 0xd7,
	0x34, 0x39, 0x86, 0x39, 0x84, 0x3e, 0x83, 0x26, 0xdf, 0xa7, 0x50, 0x74, 0xb8, 0x02, 0x18, 0x54,
	0x08, 0x3e, 0x85, 0xc6, 0xc2, 0x0d, 0x02, 0x27, 0xf3, 0x43, 0x22, 0x4b, 0x9c, 0x16, 0x19, 0x30,
	0xf3, 0x43, 0xc2, 0xc8, 0x07, 0x37, 0xf2, 0x72, 0xf2, 0x30, 0x27, 0x19, 0xc0, 0x49, 0x19, 0xf6,
	0xe3, 0xc0, 0x7d, 0x24, 0x49, 0x2a, 0xa3, 0xfe, 0xde, 0xb0, 0x8e, 0xd7, 0x21, 0x3a, 0x86, 0x7a,
	0x4a, 0x16, 0xab, 0x84, 0xc8, 0x47, 0xfc, 0xf6, 0x8a, 0x08, 0x8d, 0x41, 0x0c, 0xe8, 0x22, 0xaf,
	0xa7, 0xdb, 0x17, 0x86, 0xcd, 0xf1, 0x71, 0xf1, 0x9c, 0x45, 0x3f, 0x1a, 0x05, 0x8b, 0x37, 0xba,
	0xc1, 0x3b, 0x38, 0x64, 0xcd, 0xa4, 0xd2, 0xe8, 0xde, 0x5f, 0x6a, 0x01, 0x09, 0x49, 0x94, 0x95,
	0x9a, 0xb4, 0xc6, 0x9b, 0xb4, 0x07, 0xb0, 0x70, 0x13, 0xcf, 0x59, 0xd0, 0x55, 0x94, 0xf1, 0x4e,
	0xad, 0xe1, 0x06, 0x43, 0x54, 0x06, 0xa0, 0x53, 0x10, 0x37, 0xf7, 0xb8, 0x97, 0x57, 0xb1, 0x8e,
	0x59, 0xae, 0xc5, 0xdd, 0x54, 0x39, 0x53, 0x44, 0x83, 0x4f, 0xe0, 0x48, 0x1d, 0xdb, 0xd7, 0x24,
	0xdb, 0x9e, 0x8e, 0xc9, 0xdb, 0xc1, 0x6f, 0x02, 0x1c, 0xd9, 0x63, 0x75, 0x17, 0x4f, 0x63, 0x34,
	0x82, 0xfa, 0x03, 0x71, 0x3d, 0x92, 0xf0, 0xac, 0x9a, 0x9b, 0x3e, 0x2d, 0x0a, 0xc3, 0x69, 0xfc,
	0x9a, 0xd3, 0xb8, 0x90, 0xa1, 0xef, 0x40, 0x24, 0x79, 0x35, 0xa9, 0x5c, 0xe9, 0xef, 0x0d, 0x9b,
	0xe3, 0xde, 0x53, 0x03, 0xed, 0xd4, 0x8c, 0x37, 0xf2, 0x41, 0x0c, 0x1d, 0x75, 0x6c, 0x6b, 0x51,
	0x46, 0x12, 0x26, 0xc3, 0xe4, 0x2d, 0x1a, 0x42, 0x95, 0x5b, 0x51, 0x78, 0xc6, 0x8a, 0x5c, 0x81,
	0x5e, 0xc2, 0x3e, 0x77, 0xae, 0xef, 0x15, 0x8e, 0xae, 0xb3, 0x50, 0xf7, 0xd8, 0x25, 0xc5, 0x6e,
	0x9a, 0xbe, 0xa7, 0x89, 0xc7, 0x2f, 0xa9, 0x81, 0x37, 0xf1, 0xe0, 0x12, 0x3a, 0xf6, 0x58, 0xdd,
	0x9e, 0xf8, 0x11, 0x05, 0x0f, 0xfe, 0xaa, 0x80, 0xa4, 0x8e, 0x6d, 0x35, 0x21, 0x6e, 0x46, 0xd6,
	0x79, 0x7f, 0x0f, 0xad, 0x79, 0x6e, 0xe2, 0xdc, 0xe4, 0xc2, 0xf3, 0x26, 0x6f, 0xce, 0xb7, 0x01,
	0x3a, 0x01, 0x91, 0x3f, 0x7a, 0xb4, 0x0a, 0x8b, 0x27, 0xdf, 0x67, 0xb1, 0xb9, 0x0a, 0x9f, 0xab,
	0x65, 0x3b, 0x1e, 0xaa, 0xff, 0x36, 0x1e, 0x6a, 0xa5, 0xf1, 0xb0, 0x3b, 0x08, 0xea, 0xff, 0x33,
	0x08, 0xf6, 0x9f, 0x0e, 0x82, 0xad, 0x0b, 0xc4, 0xff, 0x74, 0x41, 0xe3, 0x03, 0x5d, 0xa0, 0x82,
	0x64, 0x8f, 0xd5, 0xd2, 0xdd, 0x7d, 0xc4, 0x0b, 0x9c, 0xff, 0x90, 0x8f, 0x79, 0x3e, 0x8d, 0x0f,
	0xe1, 0x00, 0x5b, 0xd6, 0xd4, 0xc1, 0x9a, 0xad, 0xe1, 0x3b, 0x6d, 0x22, 0xbd, 0x40, 0x1d, 0x68,
	0x72, 0xc8, 0xb4, 0xf0, 0x54, 0x31, 0x24, 0x01, 0x49, 0xd0, 0xe2, 0xc0, 0x0d, 0xd6, 0xef, 0x94,
	0x99, 0x26, 0x55, 0xce, 0x5f, 0x43, 0xb3, 0xf4, 0x1c, 0xa8, 0x0b, 0xd2, 0xa5, 0x6e, 0x18, 0xba,
	0x79, 0x5d, 0xde, 0x47, 0x82, 0xd6, 0x1a, 0x9d, 0xe9, 0x53, 0x4d, 0x12, 0xd8, 0x61, 0x1b, 0x9d,
	0x75, 0x6b, 0x4e, 0xa4, 0xca, 0xf9, 0x8f, 0x20, 0xae, 0x07, 0x3e, 0x3a, 0x06, 0x74, 0xad, 0x4c,
	0x35, 0x67, 0x6a, 0x4d, 0xb4, 0xf2, 0x46, 0x08, 0xda, 0x25, 0x5c, 0x51, 0xd9, 0x56, 0x3b, 0x98,
	0xaa, 0x18, 0x86, 0x54, 0x39, 0xff, 0x53, 0x80, 0x4e, 0xc9, 0x0d, 0xe9, 0x2a, 0xc8, 0xd8, 0x91,
	0x9a, 0x39, 0xd3, 0xb0, 0xc3, 0x2b, 0xb0, 0x7e, 0x92, 0x5e, 0xa0, 0x01, 0x9c, 0x95, 0x20, 0xdd,
	0xb4, 0x6f, 0xaf, 0xae, 0x74, 0x55, 0xd7, 0xcc, 0x99, 0xa3, 0x62, 0x6d, 0xa2, 0xcf, 0x24, 0x01,
	0xf5, 0xe0, 0xa4, 0xa4, 0x79, 0x83, 0x2d, 0xf3, 0xda, 0xb9, 0x51, 0x6c, 0xfb, 0x8d, 0x85, 0x27,
	0x52, 0x05, 0xc9, 0xd0, 0x2d, 0xd1, 0xa6, 0x35, 0x73, 0xb4, 0x9f, 0x75, 0x7b, 0x26, 0xed, 0xa1,
	0x23, 0xe8, 0x94, 0x98, 0xab, 0x5b, 0xc3, 0x90, 0xaa, 0xe8, 0x0c, 0x4e, 0xcb, 0xa0, 0xa2, 0x1b,
	0x8e, 0x69, 0x39, 0xd6, 0x25, 0xaf, 0x50, 0xaa, 0x9d, 0xff, 0x2e, 0x80, 0x54, 0xf6, 0x03, 0xcf,
	0x1c, 0x41, 0x5b, 0xc5, 0x9a, 0x32, 0xd3, 0x4a, 0xa9, 0xf7, 0xe0, 0xa4, 0x8c, 0xe9, 0xe6, 0x9d,
	0x62, 0xe8, 0x13, 0xe7, 0x46, 0xc1, 0xca, 0x54, 0x12, 0xd0, 0x2b, 0xe8, 0xed, 0xd2, 0xe5, 0xd2,
	0x14, 0x9e, 0xf9, 0x93, 0x1d, 0x14, 0x03, 0x6b, 0xca, 0xe4, 0x97, 0x75, 0xfa, 0x97, 0xc3, 0x5f,
	0xbf, 0x5c, 0xfa, 0xd9, 0xc3, 0x6a, 0xce, 0xba, 0x67, 0x14, 0xba, 0x69, 0x46, 0x92, 0xaf, 0x96,
	0x23, 0x1a, 0x2e, 0xe9, 0x88, 0x77, 0xd5, 0x28, 0x9e, 0x8f, 0x98, 0x1d, 0xe7, 0x75, 0x1e, 0x7e,
	0xf3, 0x77, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2f, 0xbe, 0xd7, 0xf4, 0x5b, 0x08, 0x00, 0x00,
}
