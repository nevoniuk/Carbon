// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Calc gRPC server encoders and decoders
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design -o services/calc

package server

import (
	"context"

	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	goagrpc "goa.design/goa/v3/grpc"
	"google.golang.org/grpc/metadata"
)

// EncodeHistoricalCarbonEmissionsResponse encodes responses from the "Calc"
// service "historical_carbon_emissions" endpoint.
func EncodeHistoricalCarbonEmissionsResponse(ctx context.Context, v interface{}, hdr, trlr *metadata.MD) (interface{}, error) {
	result, ok := v.(*calc.AllReports)
	if !ok {
		return nil, goagrpc.ErrInvalidType("Calc", "historical_carbon_emissions", "*calc.AllReports", v)
	}
	resp := NewProtoHistoricalCarbonEmissionsResponse(result)
	return resp, nil
}

// DecodeHistoricalCarbonEmissionsRequest decodes requests sent to "Calc"
// service "historical_carbon_emissions" endpoint.
func DecodeHistoricalCarbonEmissionsRequest(ctx context.Context, v interface{}, md metadata.MD) (interface{}, error) {
	var (
		message *calcpb.HistoricalCarbonEmissionsRequest
		ok      bool
	)
	{
		if message, ok = v.(*calcpb.HistoricalCarbonEmissionsRequest); !ok {
			return nil, goagrpc.ErrInvalidType("Calc", "historical_carbon_emissions", "*calcpb.HistoricalCarbonEmissionsRequest", v)
		}
		if err := ValidateHistoricalCarbonEmissionsRequest(message); err != nil {
			return nil, err
		}
	}
	var payload *calc.RequestPayload
	{
		payload = NewHistoricalCarbonEmissionsPayload(message)
	}
	return payload, nil
}
