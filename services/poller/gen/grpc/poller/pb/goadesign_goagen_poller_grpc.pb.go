// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pollerpb

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

// PollerClient is the client API for Poller service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PollerClient interface {
	// query api getting search data for carbon_intensity event
	CarbonEmissions(ctx context.Context, in *CarbonEmissionsRequest, opts ...grpc.CallOption) (*CarbonEmissionsResponse, error)
	// get the aggregate data for an event from clickhouse
	AggregateDataEndpoint(ctx context.Context, in *AggregateDataRequest, opts ...grpc.CallOption) (*AggregateDataResponse, error)
}

type pollerClient struct {
	cc grpc.ClientConnInterface
}

func NewPollerClient(cc grpc.ClientConnInterface) PollerClient {
	return &pollerClient{cc}
}

func (c *pollerClient) CarbonEmissions(ctx context.Context, in *CarbonEmissionsRequest, opts ...grpc.CallOption) (*CarbonEmissionsResponse, error) {
	out := new(CarbonEmissionsResponse)
	err := c.cc.Invoke(ctx, "/poller.Poller/CarbonEmissions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pollerClient) AggregateDataEndpoint(ctx context.Context, in *AggregateDataRequest, opts ...grpc.CallOption) (*AggregateDataResponse, error) {
	out := new(AggregateDataResponse)
	err := c.cc.Invoke(ctx, "/poller.Poller/AggregateDataEndpoint", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PollerServer is the server API for Poller service.
// All implementations must embed UnimplementedPollerServer
// for forward compatibility
type PollerServer interface {
	// query api getting search data for carbon_intensity event
	CarbonEmissions(context.Context, *CarbonEmissionsRequest) (*CarbonEmissionsResponse, error)
	// get the aggregate data for an event from clickhouse
	AggregateDataEndpoint(context.Context, *AggregateDataRequest) (*AggregateDataResponse, error)
	mustEmbedUnimplementedPollerServer()
}

// UnimplementedPollerServer must be embedded to have forward compatible implementations.
type UnimplementedPollerServer struct {
}

func (UnimplementedPollerServer) CarbonEmissions(context.Context, *CarbonEmissionsRequest) (*CarbonEmissionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CarbonEmissions not implemented")
}
func (UnimplementedPollerServer) AggregateDataEndpoint(context.Context, *AggregateDataRequest) (*AggregateDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AggregateDataEndpoint not implemented")
}
func (UnimplementedPollerServer) mustEmbedUnimplementedPollerServer() {}

// UnsafePollerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PollerServer will
// result in compilation errors.
type UnsafePollerServer interface {
	mustEmbedUnimplementedPollerServer()
}

func RegisterPollerServer(s grpc.ServiceRegistrar, srv PollerServer) {
	s.RegisterService(&Poller_ServiceDesc, srv)
}

func _Poller_CarbonEmissions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CarbonEmissionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).CarbonEmissions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poller.Poller/CarbonEmissions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).CarbonEmissions(ctx, req.(*CarbonEmissionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Poller_AggregateDataEndpoint_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AggregateDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PollerServer).AggregateDataEndpoint(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poller.Poller/AggregateDataEndpoint",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PollerServer).AggregateDataEndpoint(ctx, req.(*AggregateDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Poller_ServiceDesc is the grpc.ServiceDesc for Poller service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Poller_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "poller.Poller",
	HandlerType: (*PollerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CarbonEmissions",
			Handler:    _Poller_CarbonEmissions_Handler,
		},
		{
			MethodName: "AggregateDataEndpoint",
			Handler:    _Poller_AggregateDataEndpoint_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "goadesign_goagen_poller.proto",
}