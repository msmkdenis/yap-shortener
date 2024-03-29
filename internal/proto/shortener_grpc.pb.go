// Code generated by protoc-gen-go-grpchandlers. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpchandlers v1.3.0
// - protoc             v4.25.3
// source: internal/proto/shortener.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpchandlers package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	URLShortener_GetListURLs_FullMethodName        = "/proto.URLShortener/GetListURLs"
	URLShortener_PostURL_FullMethodName            = "/proto.URLShortener/PostURL"
	URLShortener_PostBatchURLs_FullMethodName      = "/proto.URLShortener/PostBatchURLs"
	URLShortener_GetURL_FullMethodName             = "/proto.URLShortener/GetURL"
	URLShortener_Ping_FullMethodName               = "/proto.URLShortener/Ping"
	URLShortener_DeleteAllURLs_FullMethodName      = "/proto.URLShortener/DeleteAllURLs"
	URLShortener_GetURLsByUserID_FullMethodName    = "/proto.URLShortener/GetURLsByUserID"
	URLShortener_DeleteURLsByUserID_FullMethodName = "/proto.URLShortener/DeleteURLsByUserID"
	URLShortener_GetStats_FullMethodName           = "/proto.URLShortener/GetStats"
)

// URLShortenerClient is the client API for URLShortener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type URLShortenerClient interface {
	GetListURLs(ctx context.Context, in *GetListURLsRequest, opts ...grpc.CallOption) (*GetListURLsResponse, error)
	PostURL(ctx context.Context, in *PostURLRequest, opts ...grpc.CallOption) (*PostURLResponse, error)
	PostBatchURLs(ctx context.Context, in *PostBatchURLRequest, opts ...grpc.CallOption) (*PostBatchURLResponse, error)
	GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error)
	Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error)
	DeleteAllURLs(ctx context.Context, in *DeleteAllURLsRequest, opts ...grpc.CallOption) (*DeleteAllURLsResponse, error)
	GetURLsByUserID(ctx context.Context, in *GetURLsByUserIDRequest, opts ...grpc.CallOption) (*GetURLsByUserIDResponse, error)
	DeleteURLsByUserID(ctx context.Context, in *DeleteURLsByUserIDRequest, opts ...grpc.CallOption) (*DeleteURLsByUserIDResponse, error)
	GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error)
}

type uRLShortenerClient struct {
	cc grpc.ClientConnInterface
}

func NewURLShortenerClient(cc grpc.ClientConnInterface) URLShortenerClient {
	return &uRLShortenerClient{cc}
}

func (c *uRLShortenerClient) GetListURLs(ctx context.Context, in *GetListURLsRequest, opts ...grpc.CallOption) (*GetListURLsResponse, error) {
	out := new(GetListURLsResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetListURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) PostURL(ctx context.Context, in *PostURLRequest, opts ...grpc.CallOption) (*PostURLResponse, error) {
	out := new(PostURLResponse)
	err := c.cc.Invoke(ctx, URLShortener_PostURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) PostBatchURLs(ctx context.Context, in *PostBatchURLRequest, opts ...grpc.CallOption) (*PostBatchURLResponse, error) {
	out := new(PostBatchURLResponse)
	err := c.cc.Invoke(ctx, URLShortener_PostBatchURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetURL(ctx context.Context, in *GetURLRequest, opts ...grpc.CallOption) (*GetURLResponse, error) {
	out := new(GetURLResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetURL_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) Ping(ctx context.Context, in *PingRequest, opts ...grpc.CallOption) (*PingResponse, error) {
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, URLShortener_Ping_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) DeleteAllURLs(ctx context.Context, in *DeleteAllURLsRequest, opts ...grpc.CallOption) (*DeleteAllURLsResponse, error) {
	out := new(DeleteAllURLsResponse)
	err := c.cc.Invoke(ctx, URLShortener_DeleteAllURLs_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetURLsByUserID(ctx context.Context, in *GetURLsByUserIDRequest, opts ...grpc.CallOption) (*GetURLsByUserIDResponse, error) {
	out := new(GetURLsByUserIDResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetURLsByUserID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) DeleteURLsByUserID(ctx context.Context, in *DeleteURLsByUserIDRequest, opts ...grpc.CallOption) (*DeleteURLsByUserIDResponse, error) {
	out := new(DeleteURLsByUserIDResponse)
	err := c.cc.Invoke(ctx, URLShortener_DeleteURLsByUserID_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *uRLShortenerClient) GetStats(ctx context.Context, in *GetStatsRequest, opts ...grpc.CallOption) (*GetStatsResponse, error) {
	out := new(GetStatsResponse)
	err := c.cc.Invoke(ctx, URLShortener_GetStats_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// URLShortenerServer is the server API for URLShortener service.
// All implementations must embed UnimplementedURLShortenerServer
// for forward compatibility
type URLShortenerServer interface {
	GetListURLs(context.Context, *GetListURLsRequest) (*GetListURLsResponse, error)
	PostURL(context.Context, *PostURLRequest) (*PostURLResponse, error)
	PostBatchURLs(context.Context, *PostBatchURLRequest) (*PostBatchURLResponse, error)
	GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error)
	Ping(context.Context, *PingRequest) (*PingResponse, error)
	DeleteAllURLs(context.Context, *DeleteAllURLsRequest) (*DeleteAllURLsResponse, error)
	GetURLsByUserID(context.Context, *GetURLsByUserIDRequest) (*GetURLsByUserIDResponse, error)
	DeleteURLsByUserID(context.Context, *DeleteURLsByUserIDRequest) (*DeleteURLsByUserIDResponse, error)
	GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error)
	mustEmbedUnimplementedURLShortenerServer()
}

// UnimplementedURLShortenerServer must be embedded to have forward compatible implementations.
type UnimplementedURLShortenerServer struct {
}

func (UnimplementedURLShortenerServer) GetListURLs(context.Context, *GetListURLsRequest) (*GetListURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetListURLs not implemented")
}
func (UnimplementedURLShortenerServer) PostURL(context.Context, *PostURLRequest) (*PostURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostURL not implemented")
}
func (UnimplementedURLShortenerServer) PostBatchURLs(context.Context, *PostBatchURLRequest) (*PostBatchURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PostBatchURLs not implemented")
}
func (UnimplementedURLShortenerServer) GetURL(context.Context, *GetURLRequest) (*GetURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURL not implemented")
}
func (UnimplementedURLShortenerServer) Ping(context.Context, *PingRequest) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedURLShortenerServer) DeleteAllURLs(context.Context, *DeleteAllURLsRequest) (*DeleteAllURLsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAllURLs not implemented")
}
func (UnimplementedURLShortenerServer) GetURLsByUserID(context.Context, *GetURLsByUserIDRequest) (*GetURLsByUserIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetURLsByUserID not implemented")
}
func (UnimplementedURLShortenerServer) DeleteURLsByUserID(context.Context, *DeleteURLsByUserIDRequest) (*DeleteURLsByUserIDResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteURLsByUserID not implemented")
}
func (UnimplementedURLShortenerServer) GetStats(context.Context, *GetStatsRequest) (*GetStatsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStats not implemented")
}
func (UnimplementedURLShortenerServer) mustEmbedUnimplementedURLShortenerServer() {}

// UnsafeURLShortenerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to URLShortenerServer will
// result in compilation errors.
type UnsafeURLShortenerServer interface {
	mustEmbedUnimplementedURLShortenerServer()
}

func RegisterURLShortenerServer(s grpc.ServiceRegistrar, srv URLShortenerServer) {
	s.RegisterService(&URLShortener_ServiceDesc, srv)
}

func _URLShortener_GetListURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetListURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetListURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetListURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetListURLs(ctx, req.(*GetListURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_PostURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).PostURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_PostURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).PostURL(ctx, req.(*PostURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_PostBatchURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PostBatchURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).PostBatchURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_PostBatchURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).PostBatchURLs(ctx, req.(*PostBatchURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetURL(ctx, req.(*GetURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).Ping(ctx, req.(*PingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_DeleteAllURLs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAllURLsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).DeleteAllURLs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_DeleteAllURLs_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).DeleteAllURLs(ctx, req.(*DeleteAllURLsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetURLsByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetURLsByUserIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetURLsByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetURLsByUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetURLsByUserID(ctx, req.(*GetURLsByUserIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_DeleteURLsByUserID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteURLsByUserIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).DeleteURLsByUserID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_DeleteURLsByUserID_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).DeleteURLsByUserID(ctx, req.(*DeleteURLsByUserIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _URLShortener_GetStats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStatsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(URLShortenerServer).GetStats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: URLShortener_GetStats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(URLShortenerServer).GetStats(ctx, req.(*GetStatsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// URLShortener_ServiceDesc is the grpc.ServiceDesc for URLShortener service.
// It's only intended for direct use with grpchandlers.RegisterService,
// and not to be introspected or modified (even as a copy)
var URLShortener_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.URLShortener",
	HandlerType: (*URLShortenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetListURLs",
			Handler:    _URLShortener_GetListURLs_Handler,
		},
		{
			MethodName: "PostURL",
			Handler:    _URLShortener_PostURL_Handler,
		},
		{
			MethodName: "PostBatchURLs",
			Handler:    _URLShortener_PostBatchURLs_Handler,
		},
		{
			MethodName: "GetURL",
			Handler:    _URLShortener_GetURL_Handler,
		},
		{
			MethodName: "Ping",
			Handler:    _URLShortener_Ping_Handler,
		},
		{
			MethodName: "DeleteAllURLs",
			Handler:    _URLShortener_DeleteAllURLs_Handler,
		},
		{
			MethodName: "GetURLsByUserID",
			Handler:    _URLShortener_GetURLsByUserID_Handler,
		},
		{
			MethodName: "DeleteURLsByUserID",
			Handler:    _URLShortener_DeleteURLsByUserID_Handler,
		},
		{
			MethodName: "GetStats",
			Handler:    _URLShortener_GetStats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/proto/shortener.proto",
}
