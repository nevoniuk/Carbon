// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc gRPC client encoders and decoders
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design

package client

import (
	"context"

	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// BuildHandleRequestsFunc builds the remote method to invoke for "calc"
// service "handle_requests" endpoint.
func BuildHandleRequestsFunc(grpccli calcpb.CalcClient, cliopts ...grpc.CallOption) goagrpc.RemoteFunc {
	return func(ctx context.Context, reqpb interface{}, opts ...grpc.CallOption) (interface{}, error) {
		for _, opt := range cliopts {
			opts = append(opts, opt)
		}
		if reqpb != nil {
			return grpccli.HandleRequests(ctx, reqpb.(*calcpb.HandleRequestsRequest), opts...)
		}
		return grpccli.HandleRequests(ctx, &calcpb.HandleRequestsRequest{}, opts...)
	}
}

// EncodeHandleRequestsRequest encodes requests sent to calc handle_requests
// endpoint.
func EncodeHandleRequestsRequest(ctx context.Context, v interface{}, md *metadata.MD) (interface{}, error) {
	payload, ok := v.(*calc.RequestPayload)
	if !ok {
		return nil, goagrpc.ErrInvalidType("calc", "handle_requests", "*calc.RequestPayload", v)
	}
	return NewProtoHandleRequestsRequest(payload), nil
}

// DecodeHandleRequestsResponse decodes responses from the calc handle_requests
// endpoint.
func DecodeHandleRequestsResponse(ctx context.Context, v interface{}, hdr, trlr metadata.MD) (interface{}, error) {
	message, ok := v.(*calcpb.HandleRequestsResponse)
	if !ok {
		return nil, goagrpc.ErrInvalidType("calc", "handle_requests", "*calcpb.HandleRequestsResponse", v)
	}
	if err := ValidateHandleRequestsResponse(message); err != nil {
		return nil, err
	}
	res := NewHandleRequestsResult(message)
	return res, nil
}

// BuildGetCarbonReportFunc builds the remote method to invoke for "calc"
// service "get_carbon_report" endpoint.
func BuildGetCarbonReportFunc(grpccli calcpb.CalcClient, cliopts ...grpc.CallOption) goagrpc.RemoteFunc {
	return func(ctx context.Context, reqpb interface{}, opts ...grpc.CallOption) (interface{}, error) {
		for _, opt := range cliopts {
			opts = append(opts, opt)
		}
		if reqpb != nil {
			return grpccli.GetCarbonReport(ctx, reqpb.(*calcpb.GetCarbonReportRequest), opts...)
		}
		return grpccli.GetCarbonReport(ctx, &calcpb.GetCarbonReportRequest{}, opts...)
	}
}
