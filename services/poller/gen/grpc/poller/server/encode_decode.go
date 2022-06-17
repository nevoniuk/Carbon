// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC server encoders and decoders
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package server

import (
	"context"

	poller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc/metadata"
)

// EncodeCarbonEmissionsResponse encodes responses from the "Poller" service
// "carbon_emissions" endpoint.
func EncodeCarbonEmissionsResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.([]*poller.CarbonForecast)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Poller", "carbon_emissions", "[]*poller.CarbonForecast", v)
	}
	resp := NewProtoCarbonEmissionsResponse(result)
	return resp, nil
}

// EncodeFuelsResponse encodes responses from the "Poller" service "fuels"
// endpoint.
func EncodeFuelsResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.([]*poller.FuelsForecast)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Poller", "fuels", "[]*poller.FuelsForecast", v)
	}
	resp := NewProtoFuelsResponse(result)
	return resp, nil
}

// EncodeAggregateDataEndpointResponse encodes responses from the "Poller"
// service "aggregate_data" endpoint.
func EncodeAggregateDataEndpointResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.([]*poller.AggregateData)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Poller", "aggregate_data", "[]*poller.AggregateData", v)
	}
	resp := NewProtoAggregateDataResponse(result)
	return resp, nil
}
