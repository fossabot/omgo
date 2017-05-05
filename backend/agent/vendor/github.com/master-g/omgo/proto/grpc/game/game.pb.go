// Code generated by protoc-gen-go.
// source: game.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	game.proto

It has these top-level messages:
	Game
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

type Game_FrameType int32

const (
	Game_Message Game_FrameType = 0
	Game_Kick    Game_FrameType = 1
	Game_Ping    Game_FrameType = 2
)

var Game_FrameType_name = map[int32]string{
	0: "Message",
	1: "Kick",
	2: "Ping",
}
var Game_FrameType_value = map[string]int32{
	"Message": 0,
	"Kick":    1,
	"Ping":    2,
}

func (x Game_FrameType) String() string {
	return proto1.EnumName(Game_FrameType_name, int32(x))
}
func (Game_FrameType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

type Game struct {
}

func (m *Game) Reset()                    { *m = Game{} }
func (m *Game) String() string            { return proto1.CompactTextString(m) }
func (*Game) ProtoMessage()               {}
func (*Game) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Game_Frame struct {
	Type    Game_FrameType `protobuf:"varint,1,opt,name=Type,enum=proto.Game_FrameType" json:"Type,omitempty"`
	Message []byte         `protobuf:"bytes,2,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (m *Game_Frame) Reset()                    { *m = Game_Frame{} }
func (m *Game_Frame) String() string            { return proto1.CompactTextString(m) }
func (*Game_Frame) ProtoMessage()               {}
func (*Game_Frame) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *Game_Frame) GetType() Game_FrameType {
	if m != nil {
		return m.Type
	}
	return Game_Message
}

func (m *Game_Frame) GetMessage() []byte {
	if m != nil {
		return m.Message
	}
	return nil
}

func init() {
	proto1.RegisterType((*Game)(nil), "proto.Game")
	proto1.RegisterType((*Game_Frame)(nil), "proto.Game.Frame")
	proto1.RegisterEnum("proto.Game_FrameType", Game_FrameType_name, Game_FrameType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for GameService service

type GameServiceClient interface {
	Stream(ctx context.Context, opts ...grpc.CallOption) (GameService_StreamClient, error)
}

type gameServiceClient struct {
	cc *grpc.ClientConn
}

func NewGameServiceClient(cc *grpc.ClientConn) GameServiceClient {
	return &gameServiceClient{cc}
}

func (c *gameServiceClient) Stream(ctx context.Context, opts ...grpc.CallOption) (GameService_StreamClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_GameService_serviceDesc.Streams[0], c.cc, "/proto.GameService/Stream", opts...)
	if err != nil {
		return nil, err
	}
	x := &gameServiceStreamClient{stream}
	return x, nil
}

type GameService_StreamClient interface {
	Send(*Game_Frame) error
	Recv() (*Game_Frame, error)
	grpc.ClientStream
}

type gameServiceStreamClient struct {
	grpc.ClientStream
}

func (x *gameServiceStreamClient) Send(m *Game_Frame) error {
	return x.ClientStream.SendMsg(m)
}

func (x *gameServiceStreamClient) Recv() (*Game_Frame, error) {
	m := new(Game_Frame)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Server API for GameService service

type GameServiceServer interface {
	Stream(GameService_StreamServer) error
}

func RegisterGameServiceServer(s *grpc.Server, srv GameServiceServer) {
	s.RegisterService(&_GameService_serviceDesc, srv)
}

func _GameService_Stream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(GameServiceServer).Stream(&gameServiceStreamServer{stream})
}

type GameService_StreamServer interface {
	Send(*Game_Frame) error
	Recv() (*Game_Frame, error)
	grpc.ServerStream
}

type gameServiceStreamServer struct {
	grpc.ServerStream
}

func (x *gameServiceStreamServer) Send(m *Game_Frame) error {
	return x.ServerStream.SendMsg(m)
}

func (x *gameServiceStreamServer) Recv() (*Game_Frame, error) {
	m := new(Game_Frame)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _GameService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.GameService",
	HandlerType: (*GameServiceServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Stream",
			Handler:       _GameService_Stream_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "game.proto",
}

func init() { proto1.RegisterFile("game.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 175 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4a, 0x4f, 0xcc, 0x4d,
	0xd5, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x4a, 0x4d, 0x8c, 0x5c, 0x2c, 0xee,
	0x89, 0xb9, 0xa9, 0x52, 0x3e, 0x5c, 0xac, 0x6e, 0x45, 0x89, 0xb9, 0xa9, 0x42, 0x9a, 0x5c, 0x2c,
	0x21, 0x95, 0x05, 0xa9, 0x12, 0x8c, 0x0a, 0x8c, 0x1a, 0x7c, 0x46, 0xa2, 0x10, 0xe5, 0x7a, 0x20,
	0x35, 0x7a, 0x60, 0x05, 0x20, 0xc9, 0x20, 0xb0, 0x12, 0x21, 0x09, 0x2e, 0x76, 0xdf, 0xd4, 0xe2,
	0xe2, 0xc4, 0xf4, 0x54, 0x09, 0x26, 0x05, 0x46, 0x0d, 0x9e, 0x20, 0x18, 0x57, 0x49, 0x87, 0x8b,
	0x13, 0xae, 0x58, 0x88, 0x1b, 0xae, 0x4c, 0x80, 0x41, 0x88, 0x83, 0x8b, 0xc5, 0x3b, 0x33, 0x39,
	0x5b, 0x80, 0x11, 0xc4, 0x0a, 0xc8, 0xcc, 0x4b, 0x17, 0x60, 0x32, 0x72, 0xe4, 0xe2, 0x06, 0x99,
	0x1f, 0x9c, 0x5a, 0x54, 0x96, 0x99, 0x9c, 0x2a, 0x64, 0xc4, 0xc5, 0x16, 0x5c, 0x52, 0x94, 0x9a,
	0x98, 0x2b, 0x24, 0x88, 0x61, 0xbb, 0x14, 0xa6, 0x90, 0x06, 0xa3, 0x01, 0x63, 0x12, 0x1b, 0x58,
	0xd4, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0xd5, 0xea, 0xd1, 0xd7, 0xe3, 0x00, 0x00, 0x00,
}
