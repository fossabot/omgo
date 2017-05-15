// Code generated by protoc-gen-go.
// source: db.proto
// DO NOT EDIT!

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	db.proto

It has these top-level messages:
	DB
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import proto_common "github.com/master-g/omgo/proto/pb/common"

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

type DB struct {
}

func (m *DB) Reset()                    { *m = DB{} }
func (m *DB) String() string            { return proto1.CompactTextString(m) }
func (*DB) ProtoMessage()               {}
func (*DB) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type DB_UserKey struct {
	Usn   uint64 `protobuf:"fixed64,1,opt,name=usn" json:"usn,omitempty"`
	Uid   uint64 `protobuf:"varint,2,opt,name=uid" json:"uid,omitempty"`
	Email string `protobuf:"bytes,3,opt,name=email" json:"email,omitempty"`
}

func (m *DB_UserKey) Reset()                    { *m = DB_UserKey{} }
func (m *DB_UserKey) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserKey) ProtoMessage()               {}
func (*DB_UserKey) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 0} }

func (m *DB_UserKey) GetUsn() uint64 {
	if m != nil {
		return m.Usn
	}
	return 0
}

func (m *DB_UserKey) GetUid() uint64 {
	if m != nil {
		return m.Uid
	}
	return 0
}

func (m *DB_UserKey) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type DB_UserQueryResult struct {
	Status proto_common.ResultCode     `protobuf:"varint,1,opt,name=status,enum=proto.common.ResultCode" json:"status,omitempty"`
	Info   *proto_common.UserBasicInfo `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *DB_UserQueryResult) Reset()                    { *m = DB_UserQueryResult{} }
func (m *DB_UserQueryResult) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserQueryResult) ProtoMessage()               {}
func (*DB_UserQueryResult) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *DB_UserQueryResult) GetStatus() proto_common.ResultCode {
	if m != nil {
		return m.Status
	}
	return proto_common.ResultCode_RESULT_OK
}

func (m *DB_UserQueryResult) GetInfo() *proto_common.UserBasicInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func init() {
	proto1.RegisterType((*DB)(nil), "proto.DB")
	proto1.RegisterType((*DB_UserKey)(nil), "proto.DB.UserKey")
	proto1.RegisterType((*DB_UserQueryResult)(nil), "proto.DB.UserQueryResult")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DBService service

type DBServiceClient interface {
	QueryUser(ctx context.Context, in *DB_UserKey, opts ...grpc.CallOption) (*DB_UserQueryResult, error)
}

type dBServiceClient struct {
	cc *grpc.ClientConn
}

func NewDBServiceClient(cc *grpc.ClientConn) DBServiceClient {
	return &dBServiceClient{cc}
}

func (c *dBServiceClient) QueryUser(ctx context.Context, in *DB_UserKey, opts ...grpc.CallOption) (*DB_UserQueryResult, error) {
	out := new(DB_UserQueryResult)
	err := grpc.Invoke(ctx, "/proto.DBService/QueryUser", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DBService service

type DBServiceServer interface {
	QueryUser(context.Context, *DB_UserKey) (*DB_UserQueryResult, error)
}

func RegisterDBServiceServer(s *grpc.Server, srv DBServiceServer) {
	s.RegisterService(&_DBService_serviceDesc, srv)
}

func _DBService_QueryUser_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DB_UserKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).QueryUser(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DBService/QueryUser",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).QueryUser(ctx, req.(*DB_UserKey))
	}
	return interceptor(ctx, in, info, handler)
}

var _DBService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DBService",
	HandlerType: (*DBServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryUser",
			Handler:    _DBService_QueryUser_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "db.proto",
}

func init() { proto1.RegisterFile("db.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 229 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x48, 0x49, 0xd2, 0x2b,
	0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0x3c, 0xc9, 0xf9, 0xb9, 0xb9, 0xf9, 0x79,
	0x10, 0x41, 0xa5, 0xfd, 0x8c, 0x5c, 0x4c, 0x2e, 0x4e, 0x52, 0xce, 0x5c, 0xec, 0xa1, 0xc5, 0xa9,
	0x45, 0xde, 0xa9, 0x95, 0x42, 0x02, 0x5c, 0xcc, 0xa5, 0xc5, 0x79, 0x12, 0x8c, 0x0a, 0x8c, 0x1a,
	0x6c, 0x41, 0x20, 0x26, 0x58, 0x24, 0x33, 0x45, 0x82, 0x49, 0x81, 0x51, 0x83, 0x25, 0x08, 0xc4,
	0x14, 0x12, 0xe1, 0x62, 0x4d, 0xcd, 0x4d, 0xcc, 0xcc, 0x91, 0x60, 0x56, 0x60, 0xd4, 0xe0, 0x0c,
	0x82, 0x70, 0xa4, 0x4a, 0xb8, 0xf8, 0x41, 0x86, 0x04, 0x96, 0xa6, 0x16, 0x55, 0x06, 0xa5, 0x16,
	0x97, 0xe6, 0x94, 0x08, 0x19, 0x70, 0xb1, 0x15, 0x97, 0x24, 0x96, 0x94, 0x16, 0x83, 0xcd, 0xe3,
	0x33, 0x92, 0x80, 0x58, 0xab, 0x07, 0x75, 0x03, 0x44, 0x95, 0x73, 0x7e, 0x4a, 0x6a, 0x10, 0x54,
	0x9d, 0x90, 0x3e, 0x17, 0x4b, 0x66, 0x5e, 0x5a, 0x3e, 0xd8, 0x36, 0x6e, 0x23, 0x69, 0x54, 0xf5,
	0x20, 0xe3, 0x9d, 0x12, 0x8b, 0x33, 0x93, 0x3d, 0xf3, 0xd2, 0xf2, 0x83, 0xc0, 0x0a, 0x8d, 0xdc,
	0xb8, 0x38, 0x5d, 0x9c, 0x82, 0x53, 0x8b, 0xca, 0x32, 0x93, 0x53, 0x85, 0x2c, 0xb9, 0x38, 0xc1,
	0xd6, 0x83, 0x14, 0x0a, 0x09, 0x42, 0x35, 0xbb, 0x38, 0xe9, 0x41, 0x3d, 0x27, 0x25, 0x89, 0x2a,
	0x84, 0xe4, 0xd4, 0x24, 0x36, 0xb0, 0x8c, 0x31, 0x20, 0x00, 0x00, 0xff, 0xff, 0x00, 0x85, 0x67,
	0xd4, 0x31, 0x01, 0x00, 0x00,
}
