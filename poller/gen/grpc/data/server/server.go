// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Data gRPC server
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/poller/design

package server

import (
	"context"
	"errors"

	data "github.com/crossnokaye/carbon/gen/data"
	datapb "github.com/crossnokaye/carbon/gen/grpc/data/pb"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc/codes"
)

// Server implements the datapb.DataServer interface.
type Server struct {
	CarbonEmissionsH goagrpc.UnaryHandler
	FuelsH           goagrpc.UnaryHandler
	AggregateDataH   goagrpc.UnaryHandler
	datapb.UnimplementedDataServer
}

// ErrorNamer is an interface implemented by generated error structs that
// exposes the name of the error as defined in the expr.
type ErrorNamer interface {
	ErrorName() string
}

// New instantiates the server struct with the Data service endpoints.
func New(e *data.Endpoints, uh goagrpc.UnaryHandler) *Server {
	return &Server{
		CarbonEmissionsH: NewCarbonEmissionsHandler(e.CarbonEmissions, uh),
		FuelsH:           NewFuelsHandler(e.Fuels, uh),
		AggregateDataH:   NewAggregateDataHandler(e.AggregateData, uh),
	}
}

// NewCarbonEmissionsHandler creates a gRPC handler which serves the "Data"
// service "carbon_emissions" endpoint.
func NewCarbonEmissionsHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeCarbonEmissionsResponse)
	}
	return h
}

// CarbonEmissions implements the "CarbonEmissions" method in datapb.DataServer
// interface.
func (s *Server) CarbonEmissions(ctx context.Context, message *datapb.CarbonEmissionsRequest) (*datapb.CarbonEmissionsResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "carbon_emissions")
	ctx = context.WithValue(ctx, goa.ServiceKey, "Data")
	resp, err := s.CarbonEmissionsH.Handle(ctx, message)
	if err != nil {
		var en ErrorNamer
		if errors.As(err, &en) {
			switch en.ErrorName() {
			case "data_not_available":
				return nil, goagrpc.NewStatusError(codes.DataLoss, err, goagrpc.NewErrorResponse(err))
			case "missing-required-parameter":
				return nil, goagrpc.NewStatusError(codes.NotFound, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*datapb.CarbonEmissionsResponse), nil
}

// NewFuelsHandler creates a gRPC handler which serves the "Data" service
// "fuels" endpoint.
func NewFuelsHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeFuelsResponse)
	}
	return h
}

// Fuels implements the "Fuels" method in datapb.DataServer interface.
func (s *Server) Fuels(ctx context.Context, message *datapb.FuelsRequest) (*datapb.FuelsResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "fuels")
	ctx = context.WithValue(ctx, goa.ServiceKey, "Data")
	resp, err := s.FuelsH.Handle(ctx, message)
	if err != nil {
		var en ErrorNamer
		if errors.As(err, &en) {
			switch en.ErrorName() {
			case "data_not_available":
				return nil, goagrpc.NewStatusError(codes.DataLoss, err, goagrpc.NewErrorResponse(err))
			case "missing-required-parameter":
				return nil, goagrpc.NewStatusError(codes.NotFound, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*datapb.FuelsResponse), nil
}

// NewAggregateDataHandler creates a gRPC handler which serves the "Data"
// service "aggregate_data" endpoint.
func NewAggregateDataHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeAggregateDataResponse)
	}
	return h
}

// AggregateData implements the "AggregateData" method in datapb.DataServer
// interface.
func (s *Server) AggregateData(ctx context.Context, message *datapb.AggregateDataRequest) (*datapb.AggregateDataResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "aggregate_data")
	ctx = context.WithValue(ctx, goa.ServiceKey, "Data")
	resp, err := s.AggregateDataH.Handle(ctx, message)
	if err != nil {
		var en ErrorNamer
		if errors.As(err, &en) {
			switch en.ErrorName() {
			case "data_not_available":
				return nil, goagrpc.NewStatusError(codes.DataLoss, err, goagrpc.NewErrorResponse(err))
			case "missing-required-parameter":
				return nil, goagrpc.NewStatusError(codes.NotFound, err, goagrpc.NewErrorResponse(err))
			}
		}
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*datapb.AggregateDataResponse), nil
}
