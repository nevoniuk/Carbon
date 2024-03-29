// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC server encoders and decoders
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design -o
// services/poller

package server

import (
	"context"

	pollerpb "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/pb"
	poller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc/metadata"
)

// EncodeUpdateResponse encodes responses from the "Poller" service "update"
// endpoint.
func EncodeUpdateResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	resp := NewProtoUpdateResponse()
	return resp, nil
}

// EncodeGetEmissionsForRegionResponse encodes responses from the "Poller"
// service "get_emissions_for_region" endpoint.
func EncodeGetEmissionsForRegionResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.([]*poller.CarbonForecast)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Poller", "get_emissions_for_region", "[]*poller.CarbonForecast", v)
	}
	resp := NewProtoGetEmissionsForRegionResponse(result)
	return resp, nil
}

// DecodeGetEmissionsForRegionRequest decodes requests sent to "Poller" service
// "get_emissions_for_region" endpoint.
func DecodeGetEmissionsForRegionRequest(ctx context.Context, v interface{}, md metadata.MD) (interface{}, error) {
	var (
		message *pollerpb.GetEmissionsForRegionRequest
		ok      bool
	)
	{
		if message, ok = v.(*pollerpb.GetEmissionsForRegionRequest); !ok {
			return nil, goagrpc.ErrInvalidType("Poller", "get_emissions_for_region", "*pollerpb.GetEmissionsForRegionRequest", v)
		}
		if err := ValidateGetEmissionsForRegionRequest(message); err != nil {
			return nil, err
		}
	}
	var payload *poller.CarbonPayload
	{
		payload = NewGetEmissionsForRegionPayload(message)
	}
	return payload, nil
}
