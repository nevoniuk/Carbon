// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC client encoders and decoders
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design -o
// services/poller

package client

import (
	"context"

	pollerpb "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/pb"
	poller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// BuildUpdateFunc builds the remote method to invoke for "Poller" service
// "update" endpoint.
func BuildUpdateFunc(grpccli pollerpb.PollerClient, cliopts ...grpc.CallOption) goagrpc.RemoteFunc {
	return func(ctx context.Context, reqpb interface{}, opts ...grpc.CallOption) (interface{}, error) {
		for _, opt := range cliopts {
			opts = append(opts, opt)
		}
		if reqpb != nil {
			return grpccli.Update(ctx, reqpb.(*pollerpb.UpdateRequest), opts...)
		}
		return grpccli.Update(ctx, &pollerpb.UpdateRequest{}, opts...)
	}
}

// BuildGetEmissionsForRegionFunc builds the remote method to invoke for
// "Poller" service "get_emissions_for_region" endpoint.
func BuildGetEmissionsForRegionFunc(grpccli pollerpb.PollerClient, cliopts ...grpc.CallOption) goagrpc.RemoteFunc {
	return func(ctx context.Context, reqpb interface{}, opts ...grpc.CallOption) (interface{}, error) {
		for _, opt := range cliopts {
			opts = append(opts, opt)
		}
		if reqpb != nil {
			return grpccli.GetEmissionsForRegion(ctx, reqpb.(*pollerpb.GetEmissionsForRegionRequest), opts...)
		}
		return grpccli.GetEmissionsForRegion(ctx, &pollerpb.GetEmissionsForRegionRequest{}, opts...)
	}
}

// EncodeGetEmissionsForRegionRequest encodes requests sent to Poller
// get_emissions_for_region endpoint.
func EncodeGetEmissionsForRegionRequest(ctx context.Context, v interface{}, md *metadata.MD) (interface{}, error) {
	payload, ok := v.(*poller.CarbonPayload)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Poller", "get_emissions_for_region", "*poller.CarbonPayload", v)
	}
	return NewProtoGetEmissionsForRegionRequest(payload), nil
}

// DecodeGetEmissionsForRegionResponse decodes responses from the Poller
// get_emissions_for_region endpoint.
func DecodeGetEmissionsForRegionResponse(ctx context.Context, v interface{}, hdr, trlr metadata.MD) (interface{}, error) {
	message, ok := v.(*pollerpb.GetEmissionsForRegionResponse)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Poller", "get_emissions_for_region", "*pollerpb.GetEmissionsForRegionResponse", v)
	}
	if err := ValidateGetEmissionsForRegionResponse(message); err != nil {
		return nil, err
	}
	res := NewGetEmissionsForRegionResult(message)
	return res, nil
}
