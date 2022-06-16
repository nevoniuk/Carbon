// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC client types
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package client

import (
	pollerpb "github.com/crossnokaye/carbon/gen/grpc/poller/pb"
)

// NewProtoCarbonEmissionsRequest builds the gRPC request type from the payload
// of the "carbon_emissions" endpoint of the "Poller" service.
func NewProtoCarbonEmissionsRequest() *pollerpb.CarbonEmissionsRequest {
	message := &pollerpb.CarbonEmissionsRequest{}
	return message
}

// NewProtoFuelsRequest builds the gRPC request type from the payload of the
// "fuels" endpoint of the "Poller" service.
func NewProtoFuelsRequest() *pollerpb.FuelsRequest {
	message := &pollerpb.FuelsRequest{}
	return message
}

// NewProtoAggregateDataRequest builds the gRPC request type from the payload
// of the "aggregate_data" endpoint of the "Poller" service.
func NewProtoAggregateDataRequest() *pollerpb.AggregateDataRequest {
	message := &pollerpb.AggregateDataRequest{}
	return message
}
