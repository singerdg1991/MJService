// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: tikka/protobuf/tikka.proto

package protobuf

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// OTPServiceClient is the client API for OTPService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OTPServiceClient interface {
	SendOTP(ctx context.Context, in *OTPRequest, opts ...grpc.CallOption) (*OTPResponse, error)
}

type oTPServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOTPServiceClient(cc grpc.ClientConnInterface) OTPServiceClient {
	return &oTPServiceClient{cc}
}

func (c *oTPServiceClient) SendOTP(ctx context.Context, in *OTPRequest, opts ...grpc.CallOption) (*OTPResponse, error) {
	out := new(OTPResponse)
	err := c.cc.Invoke(ctx, "/protobuf.OTPService/SendOTP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OTPServiceServer is the server API for OTPService service.
// All implementations must embed UnimplementedOTPServiceServer
// for forward compatibility
type OTPServiceServer interface {
	SendOTP(context.Context, *OTPRequest) (*OTPResponse, error)
	mustEmbedUnimplementedOTPServiceServer()
}

// UnimplementedOTPServiceServer must be embedded to have forward compatible implementations.
type UnimplementedOTPServiceServer struct {
}

func (UnimplementedOTPServiceServer) SendOTP(context.Context, *OTPRequest) (*OTPResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendOTP not implemented")
}
func (UnimplementedOTPServiceServer) mustEmbedUnimplementedOTPServiceServer() {}

// UnsafeOTPServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OTPServiceServer will
// result in compilation errors.
type UnsafeOTPServiceServer interface {
	mustEmbedUnimplementedOTPServiceServer()
}

func RegisterOTPServiceServer(s grpc.ServiceRegistrar, srv OTPServiceServer) {
	s.RegisterService(&OTPService_ServiceDesc, srv)
}

func _OTPService_SendOTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OTPRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OTPServiceServer).SendOTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.OTPService/SendOTP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OTPServiceServer).SendOTP(ctx, req.(*OTPRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OTPService_ServiceDesc is the grpc.ServiceDesc for OTPService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OTPService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.OTPService",
	HandlerType: (*OTPServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendOTP",
			Handler:    _OTPService_SendOTP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tikka/protobuf/tikka.proto",
}

// EmailServiceClient is the client API for EmailService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EmailServiceClient interface {
	SendEmail(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error)
	SendEmailSMTP(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error)
}

type emailServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEmailServiceClient(cc grpc.ClientConnInterface) EmailServiceClient {
	return &emailServiceClient{cc}
}

func (c *emailServiceClient) SendEmail(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, "/protobuf.EmailService/SendEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *emailServiceClient) SendEmailSMTP(ctx context.Context, in *EmailRequest, opts ...grpc.CallOption) (*EmailResponse, error) {
	out := new(EmailResponse)
	err := c.cc.Invoke(ctx, "/protobuf.EmailService/SendEmailSMTP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EmailServiceServer is the server API for EmailService service.
// All implementations must embed UnimplementedEmailServiceServer
// for forward compatibility
type EmailServiceServer interface {
	SendEmail(context.Context, *EmailRequest) (*EmailResponse, error)
	SendEmailSMTP(context.Context, *EmailRequest) (*EmailResponse, error)
	mustEmbedUnimplementedEmailServiceServer()
}

// UnimplementedEmailServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEmailServiceServer struct {
}

func (UnimplementedEmailServiceServer) SendEmail(context.Context, *EmailRequest) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmail not implemented")
}
func (UnimplementedEmailServiceServer) SendEmailSMTP(context.Context, *EmailRequest) (*EmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEmailSMTP not implemented")
}
func (UnimplementedEmailServiceServer) mustEmbedUnimplementedEmailServiceServer() {}

// UnsafeEmailServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EmailServiceServer will
// result in compilation errors.
type UnsafeEmailServiceServer interface {
	mustEmbedUnimplementedEmailServiceServer()
}

func RegisterEmailServiceServer(s grpc.ServiceRegistrar, srv EmailServiceServer) {
	s.RegisterService(&EmailService_ServiceDesc, srv)
}

func _EmailService_SendEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).SendEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.EmailService/SendEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).SendEmail(ctx, req.(*EmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EmailService_SendEmailSMTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EmailServiceServer).SendEmailSMTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protobuf.EmailService/SendEmailSMTP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EmailServiceServer).SendEmailSMTP(ctx, req.(*EmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EmailService_ServiceDesc is the grpc.ServiceDesc for EmailService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EmailService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "protobuf.EmailService",
	HandlerType: (*EmailServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SendEmail",
			Handler:    _EmailService_SendEmail_Handler,
		},
		{
			MethodName: "SendEmailSMTP",
			Handler:    _EmailService_SendEmailSMTP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "tikka/protobuf/tikka.proto",
}
