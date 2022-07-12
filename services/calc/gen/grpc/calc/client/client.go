// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc gRPC client
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design

package client

import (
	"context"

	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
)

// Client lists the service endpoint gRPC clients.
type Client struct {
	grpccli calcpb.CalcClient
	opts    []grpc.CallOption
}

// NewClient instantiates gRPC client for all the calc service servers.
func NewClient(cc *grpc.ClientConn, opts ...grpc.CallOption) *Client {
	return &Client{
		grpccli: calcpb.NewCalcClient(cc),
		opts:    opts,
	}
}

// HandleRequests calls the "HandleRequests" function in calcpb.CalcClient
// interface.
func (c *Client) HandleRequests() goa.Endpoint {
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		inv := goagrpc.NewInvoker(
			BuildHandleRequestsFunc(c.grpccli, c.opts...),
			EncodeHandleRequestsRequest,
			DecodeHandleRequestsResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			return nil, goa.Fault(err.Error())
		}
		return res, nil
	}
}

// GetCarbonReport calls the "GetCarbonReport" function in calcpb.CalcClient
// interface.
func (c *Client) GetCarbonReport() goa.Endpoint {
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		inv := goagrpc.NewInvoker(
			BuildGetCarbonReportFunc(c.grpccli, c.opts...),
			nil,
			nil)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			return nil, goa.Fault(err.Error())
		}
		return res, nil
	}
}
