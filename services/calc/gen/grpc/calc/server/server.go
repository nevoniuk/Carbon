// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc gRPC server
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design

package server

import (
	"context"
	"errors"

	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc/codes"
)

// Server implements the calcpb.CalcServer interface.
type Server struct {
	HistoricalCarbonEmissionsH goagrpc.UnaryHandler
	calcpb.UnimplementedCalcServer
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the expr.
type ErrorNamer interface {
	ErrorName() string
}

// New instantiates the server struct with the calc service endpoints.
func New(e *calc.Endpoints, uh goagrpc.UnaryHandler) *Server {
	return &Server{
		HistoricalCarbonEmissionsH: NewHistoricalCarbonEmissionsHandler(e.HistoricalCarbonEmissions, uh),
	}
}

// NewHistoricalCarbonEmissionsHandler creates a gRPC handler which serves the
// "calc" service "historical_carbon_emissions" endpoint.
func NewHistoricalCarbonEmissionsHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeHistoricalCarbonEmissionsRequest, EncodeHistoricalCarbonEmissionsResponse)
	}
	return h
}

// HistoricalCarbonEmissions implements the "HistoricalCarbonEmissions" method
// in calcpb.CalcServer interface.
func (s *Server) HistoricalCarbonEmissions(ctx context.Context, message *calcpb.HistoricalCarbonEmissionsRequest) (*calcpb.HistoricalCarbonEmissionsResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "historical_carbon_emissions")
	ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
	resp, err := s.HistoricalCarbonEmissionsH.Handle(ctx, message)
	if err != nil {
		var en ErrorNamer
		if errors.As(err, &en) {
			switch en.ErrorName() {
			case "not_found":
				return nil, goagrpc.NewStatusError(codes.NotFound, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*calcpb.HistoricalCarbonEmissionsResponse), nil
}
