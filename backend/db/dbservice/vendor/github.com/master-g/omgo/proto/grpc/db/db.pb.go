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

// user query key
type DB_UserKey struct {
	Usn   uint64 `protobuf:"fixed64,1,opt,name=usn" json:"usn,omitempty"`
	Uid   uint64 `protobuf:"varint,2,opt,name=uid" json:"uid"`
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

// user extra information
type DB_UserExtraInfo struct {
	Secret []byte `protobuf:"bytes,1,opt,name=secret,proto3" json:"secret,omitempty"`
	Token  string `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
}

func (m *DB_UserExtraInfo) Reset()                    { *m = DB_UserExtraInfo{} }
func (m *DB_UserExtraInfo) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserExtraInfo) ProtoMessage()               {}
func (*DB_UserExtraInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 1} }

func (m *DB_UserExtraInfo) GetSecret() []byte {
	if m != nil {
		return m.Secret
	}
	return nil
}

func (m *DB_UserExtraInfo) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type DB_UserQueryResponse struct {
	Result *proto_common.RspHeader     `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Info   *proto_common.UserBasicInfo `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
}

func (m *DB_UserQueryResponse) Reset()                    { *m = DB_UserQueryResponse{} }
func (m *DB_UserQueryResponse) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserQueryResponse) ProtoMessage()               {}
func (*DB_UserQueryResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 2} }

func (m *DB_UserQueryResponse) GetResult() *proto_common.RspHeader {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *DB_UserQueryResponse) GetInfo() *proto_common.UserBasicInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

type DB_UserRegisterRequest struct {
	Info   *proto_common.UserBasicInfo `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
	Secret []byte                      `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
}

func (m *DB_UserRegisterRequest) Reset()                    { *m = DB_UserRegisterRequest{} }
func (m *DB_UserRegisterRequest) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserRegisterRequest) ProtoMessage()               {}
func (*DB_UserRegisterRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 3} }

func (m *DB_UserRegisterRequest) GetInfo() *proto_common.UserBasicInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *DB_UserRegisterRequest) GetSecret() []byte {
	if m != nil {
		return m.Secret
	}
	return nil
}

type DB_UserRegisterResponse struct {
	Result *proto_common.RspHeader     `protobuf:"bytes,1,opt,name=result" json:"result,omitempty"`
	Info   *proto_common.UserBasicInfo `protobuf:"bytes,2,opt,name=info" json:"info,omitempty"`
	Token  string                      `protobuf:"bytes,3,opt,name=token" json:"token,omitempty"`
}

func (m *DB_UserRegisterResponse) Reset()                    { *m = DB_UserRegisterResponse{} }
func (m *DB_UserRegisterResponse) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserRegisterResponse) ProtoMessage()               {}
func (*DB_UserRegisterResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 4} }

func (m *DB_UserRegisterResponse) GetResult() *proto_common.RspHeader {
	if m != nil {
		return m.Result
	}
	return nil
}

func (m *DB_UserRegisterResponse) GetInfo() *proto_common.UserBasicInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *DB_UserRegisterResponse) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type DB_UserLoginRequest struct {
	Info   *proto_common.UserBasicInfo `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
	Secret []byte                      `protobuf:"bytes,2,opt,name=secret,proto3" json:"secret,omitempty"`
}

func (m *DB_UserLoginRequest) Reset()                    { *m = DB_UserLoginRequest{} }
func (m *DB_UserLoginRequest) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserLoginRequest) ProtoMessage()               {}
func (*DB_UserLoginRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 5} }

func (m *DB_UserLoginRequest) GetInfo() *proto_common.UserBasicInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *DB_UserLoginRequest) GetSecret() []byte {
	if m != nil {
		return m.Secret
	}
	return nil
}

type DB_UserLoginResponse struct {
	Info  *proto_common.UserBasicInfo `protobuf:"bytes,1,opt,name=info" json:"info,omitempty"`
	Token string                      `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
}

func (m *DB_UserLoginResponse) Reset()                    { *m = DB_UserLoginResponse{} }
func (m *DB_UserLoginResponse) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserLoginResponse) ProtoMessage()               {}
func (*DB_UserLoginResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 6} }

func (m *DB_UserLoginResponse) GetInfo() *proto_common.UserBasicInfo {
	if m != nil {
		return m.Info
	}
	return nil
}

func (m *DB_UserLoginResponse) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type DB_UserLogoutRequest struct {
	Usn   uint64 `protobuf:"varint,1,opt,name=usn" json:"usn"`
	Token string `protobuf:"bytes,2,opt,name=token" json:"token,omitempty"`
}

func (m *DB_UserLogoutRequest) Reset()                    { *m = DB_UserLogoutRequest{} }
func (m *DB_UserLogoutRequest) String() string            { return proto1.CompactTextString(m) }
func (*DB_UserLogoutRequest) ProtoMessage()               {}
func (*DB_UserLogoutRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0, 7} }

func (m *DB_UserLogoutRequest) GetUsn() uint64 {
	if m != nil {
		return m.Usn
	}
	return 0
}

func (m *DB_UserLogoutRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto1.RegisterType((*DB)(nil), "proto.DB")
	proto1.RegisterType((*DB_UserKey)(nil), "proto.DB.UserKey")
	proto1.RegisterType((*DB_UserExtraInfo)(nil), "proto.DB.UserExtraInfo")
	proto1.RegisterType((*DB_UserQueryResponse)(nil), "proto.DB.UserQueryResponse")
	proto1.RegisterType((*DB_UserRegisterRequest)(nil), "proto.DB.UserRegisterRequest")
	proto1.RegisterType((*DB_UserRegisterResponse)(nil), "proto.DB.UserRegisterResponse")
	proto1.RegisterType((*DB_UserLoginRequest)(nil), "proto.DB.UserLoginRequest")
	proto1.RegisterType((*DB_UserLoginResponse)(nil), "proto.DB.UserLoginResponse")
	proto1.RegisterType((*DB_UserLogoutRequest)(nil), "proto.DB.UserLogoutRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for DBService service

type DBServiceClient interface {
	// query user info
	UserQuery(ctx context.Context, in *DB_UserKey, opts ...grpc.CallOption) (*DB_UserQueryResponse, error)
	// update user info
	UserUpdateInfo(ctx context.Context, in *proto_common.UserBasicInfo, opts ...grpc.CallOption) (*proto_common.RspHeader, error)
	// register
	UserRegister(ctx context.Context, in *DB_UserRegisterRequest, opts ...grpc.CallOption) (*DB_UserRegisterResponse, error)
	// login
	UserLogin(ctx context.Context, in *DB_UserLoginRequest, opts ...grpc.CallOption) (*DB_UserLoginResponse, error)
	// logout
	UserLogout(ctx context.Context, in *DB_UserLogoutRequest, opts ...grpc.CallOption) (*proto_common.RspHeader, error)
}

type dBServiceClient struct {
	cc *grpc.ClientConn
}

func NewDBServiceClient(cc *grpc.ClientConn) DBServiceClient {
	return &dBServiceClient{cc}
}

func (c *dBServiceClient) UserQuery(ctx context.Context, in *DB_UserKey, opts ...grpc.CallOption) (*DB_UserQueryResponse, error) {
	out := new(DB_UserQueryResponse)
	err := grpc.Invoke(ctx, "/proto.DBService/UserQuery", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) UserUpdateInfo(ctx context.Context, in *proto_common.UserBasicInfo, opts ...grpc.CallOption) (*proto_common.RspHeader, error) {
	out := new(proto_common.RspHeader)
	err := grpc.Invoke(ctx, "/proto.DBService/UserUpdateInfo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) UserRegister(ctx context.Context, in *DB_UserRegisterRequest, opts ...grpc.CallOption) (*DB_UserRegisterResponse, error) {
	out := new(DB_UserRegisterResponse)
	err := grpc.Invoke(ctx, "/proto.DBService/UserRegister", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) UserLogin(ctx context.Context, in *DB_UserLoginRequest, opts ...grpc.CallOption) (*DB_UserLoginResponse, error) {
	out := new(DB_UserLoginResponse)
	err := grpc.Invoke(ctx, "/proto.DBService/UserLogin", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dBServiceClient) UserLogout(ctx context.Context, in *DB_UserLogoutRequest, opts ...grpc.CallOption) (*proto_common.RspHeader, error) {
	out := new(proto_common.RspHeader)
	err := grpc.Invoke(ctx, "/proto.DBService/UserLogout", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for DBService service

type DBServiceServer interface {
	// query user info
	UserQuery(context.Context, *DB_UserKey) (*DB_UserQueryResponse, error)
	// update user info
	UserUpdateInfo(context.Context, *proto_common.UserBasicInfo) (*proto_common.RspHeader, error)
	// register
	UserRegister(context.Context, *DB_UserRegisterRequest) (*DB_UserRegisterResponse, error)
	// login
	UserLogin(context.Context, *DB_UserLoginRequest) (*DB_UserLoginResponse, error)
	// logout
	UserLogout(context.Context, *DB_UserLogoutRequest) (*proto_common.RspHeader, error)
}

func RegisterDBServiceServer(s *grpc.Server, srv DBServiceServer) {
	s.RegisterService(&_DBService_serviceDesc, srv)
}

func _DBService_UserQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DB_UserKey)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).UserQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DBService/UserQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).UserQuery(ctx, req.(*DB_UserKey))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_UserUpdateInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(proto_common.UserBasicInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).UserUpdateInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DBService/UserUpdateInfo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).UserUpdateInfo(ctx, req.(*proto_common.UserBasicInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_UserRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DB_UserRegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).UserRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DBService/UserRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).UserRegister(ctx, req.(*DB_UserRegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_UserLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DB_UserLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).UserLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DBService/UserLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).UserLogin(ctx, req.(*DB_UserLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DBService_UserLogout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DB_UserLogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DBServiceServer).UserLogout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.DBService/UserLogout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DBServiceServer).UserLogout(ctx, req.(*DB_UserLogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _DBService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.DBService",
	HandlerType: (*DBServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UserQuery",
			Handler:    _DBService_UserQuery_Handler,
		},
		{
			MethodName: "UserUpdateInfo",
			Handler:    _DBService_UserUpdateInfo_Handler,
		},
		{
			MethodName: "UserRegister",
			Handler:    _DBService_UserRegister_Handler,
		},
		{
			MethodName: "UserLogin",
			Handler:    _DBService_UserLogin_Handler,
		},
		{
			MethodName: "UserLogout",
			Handler:    _DBService_UserLogout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "db.proto",
}

func init() { proto1.RegisterFile("db.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 419 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xbc, 0x53, 0xd1, 0x8e, 0xd2, 0x40,
	0x14, 0x4d, 0x0b, 0x54, 0x7b, 0x41, 0x03, 0x23, 0x51, 0x32, 0x44, 0x43, 0x7c, 0xe2, 0x09, 0x12,
	0x7c, 0x24, 0xbe, 0xd4, 0x6a, 0x34, 0xea, 0x83, 0x63, 0x78, 0xd1, 0xc4, 0xa4, 0xb4, 0x17, 0xd2,
	0x08, 0x9d, 0x3a, 0xd3, 0x9a, 0xe5, 0x27, 0xf6, 0x33, 0x37, 0xd9, 0xbf, 0xd8, 0xcc, 0x4c, 0x97,
	0xb6, 0x2c, 0x6c, 0x76, 0x93, 0xcd, 0x3e, 0xcd, 0xcc, 0xe9, 0x99, 0x33, 0xf7, 0xdc, 0x7b, 0x0a,
	0x4f, 0xa3, 0xe5, 0x24, 0x15, 0x3c, 0xe3, 0xa4, 0xa5, 0x17, 0xda, 0x09, 0xf9, 0x76, 0xcb, 0x13,
	0x03, 0xbe, 0xbd, 0x68, 0x81, 0xed, 0x7b, 0xf4, 0x03, 0x3c, 0x59, 0x48, 0x14, 0x5f, 0x71, 0x47,
	0xba, 0xd0, 0xc8, 0x65, 0x32, 0xb0, 0x46, 0xd6, 0xd8, 0x61, 0x6a, 0xab, 0x91, 0x38, 0x1a, 0xd8,
	0x23, 0x6b, 0xdc, 0x64, 0x6a, 0x4b, 0xfa, 0xd0, 0xc2, 0x6d, 0x10, 0x6f, 0x06, 0x8d, 0x91, 0x35,
	0x76, 0x99, 0x39, 0xd0, 0xf7, 0xf0, 0x4c, 0x89, 0x7c, 0x3c, 0xcb, 0x44, 0xf0, 0x25, 0x59, 0x71,
	0xf2, 0x12, 0x1c, 0x89, 0xa1, 0xc0, 0x4c, 0xab, 0x75, 0x58, 0x71, 0x52, 0xd7, 0x33, 0xfe, 0x17,
	0x13, 0x2d, 0xe9, 0x32, 0x73, 0xa0, 0x39, 0xf4, 0xd4, 0xf5, 0x1f, 0x39, 0x8a, 0x1d, 0x43, 0x99,
	0xf2, 0x44, 0x22, 0x99, 0x82, 0x23, 0x50, 0xe6, 0x1b, 0x23, 0xd1, 0x9e, 0xbd, 0x32, 0x75, 0x4f,
	0x0a, 0x13, 0x4c, 0xa6, 0x9f, 0x31, 0x88, 0x50, 0xb0, 0x82, 0x46, 0xa6, 0xd0, 0x8c, 0x93, 0x15,
	0xd7, 0xd2, 0xed, 0xd9, 0xb0, 0x4e, 0x57, 0xfa, 0x5e, 0x20, 0xe3, 0x50, 0x95, 0xc7, 0x34, 0x91,
	0xfe, 0x81, 0x17, 0x0a, 0x66, 0xb8, 0x8e, 0x65, 0xa6, 0xd6, 0x7f, 0x39, 0xca, 0x52, 0xc7, 0xba,
	0xa3, 0x4e, 0xc5, 0xac, 0x5d, 0x35, 0x4b, 0xcf, 0x2d, 0xe8, 0xd7, 0x1f, 0x78, 0x2c, 0x6b, 0x65,
	0x9f, 0x1b, 0xd5, 0x3e, 0xff, 0x86, 0xae, 0x22, 0x7f, 0xe3, 0xeb, 0x38, 0x79, 0x70, 0xb7, 0xbf,
	0xcc, 0x10, 0x0b, 0xf1, 0xbd, 0xd3, 0x7b, 0xaa, 0x1f, 0x0f, 0xc8, 0x7c, 0xaf, 0xcd, 0xf3, 0xec,
	0xba, 0xf2, 0x4a, 0x5c, 0x9b, 0x26, 0xae, 0x47, 0x2f, 0xcf, 0x2e, 0x6d, 0x70, 0x7d, 0xef, 0x27,
	0x8a, 0xff, 0x71, 0x88, 0x64, 0x0e, 0xee, 0x3e, 0x6b, 0xa4, 0x57, 0x14, 0xe4, 0x7b, 0x93, 0xe2,
	0x27, 0xa0, 0xc3, 0x3a, 0x54, 0xcf, 0xe4, 0x27, 0x78, 0xae, 0xc0, 0x45, 0x1a, 0x05, 0x19, 0xea,
	0xa0, 0xdf, 0x66, 0x89, 0x9e, 0x9a, 0x2b, 0xf9, 0x0e, 0x9d, 0x6a, 0x30, 0xc8, 0xeb, 0xfa, 0xa3,
	0x07, 0x89, 0xa4, 0x6f, 0x4e, 0x7d, 0x2e, 0xca, 0xf2, 0x8d, 0x27, 0xdd, 0x7a, 0x42, 0xeb, 0xe4,
	0xea, 0xb0, 0x0f, 0xcd, 0xd5, 0x67, 0xe5, 0x01, 0x94, 0x4d, 0x26, 0x37, 0xa9, 0x65, 0xeb, 0x4f,
	0x1a, 0x5b, 0x3a, 0x1a, 0x7f, 0x77, 0x15, 0x00, 0x00, 0xff, 0xff, 0x2e, 0x96, 0xad, 0x73, 0x7c,
	0x04, 0x00, 0x00,
}
