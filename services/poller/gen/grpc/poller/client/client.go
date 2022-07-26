// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC client
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design -o
// services/poller

package client

import (
	"context"

	pollerpb "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/pb"
	goagrpc "goa.design/goa/v3/grpc"
	goapb "goa.design/goa/v3/grpc/pb"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
)

// Client lists the service endpoint gRPC clients.
type Client struct {
	grpccli pollerpb.PollerClient
	opts    []grpc.CallOption
}

// NewClient instantiates gRPC client for all the Poller service servers.
func NewClient(cc *grpc.ClientConn, opts ...grpc.CallOption) *Client {
	return &Client{
		grpccli: pollerpb.NewPollerClient(cc),
		opts:    opts,
	}
}

// Update calls the "Update" function in pollerpb.PollerClient interface.
func (c *Client) Update() goa.Endpoint {
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		inv := goagrpc.NewInvoker(
			BuildUpdateFunc(c.grpccli, c.opts...),
			nil,
			nil)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			resp := goagrpc.DecodeError(err)
			switch message := resp.(type) {
			case *goapb.ErrorResponse:
				return nil, goagrpc.NewServiceError(message)
			default:
				return nil, goa.Fault(err.Error())
			}
		}
		return res, nil
	}
}

// GetEmissionsForRegion calls the "GetEmissionsForRegion" function in
// pollerpb.PollerClient interface.
func (c *Client) GetEmissionsForRegion() goa.Endpoint {
	return func(ctx context.Context, v interface{}) (interface{}, error) {
		inv := goagrpc.NewInvoker(
			BuildGetEmissionsForRegionFunc(c.grpccli, c.opts...),
			EncodeGetEmissionsForRegionRequest,
			DecodeGetEmissionsForRegionResponse)
		res, err := inv.Invoke(ctx, v)
		if err != nil {
			resp := goagrpc.DecodeError(err)
			switch message := resp.(type) {
			case *goapb.ErrorResponse:
				return nil, goagrpc.NewServiceError(message)
			default:
				return nil, goa.Fault(err.Error())
			}
		}
		return res, nil
	}
}
