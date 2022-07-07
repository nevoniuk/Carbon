// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc gRPC server
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design

package server

import (
	"context"

	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	goagrpc "goa.design/goa/v3/grpc"
	goa "goa.design/goa/v3/pkg"
)

// Server implements the calcpb.CalcServer interface.
type Server struct {
	CalculateReportsH goagrpc.UnaryHandler
	GetControlPointsH goagrpc.UnaryHandler
	GetPowerH         goagrpc.UnaryHandler
	GetEmissionsH     goagrpc.UnaryHandler
	HandleRequestsH   goagrpc.UnaryHandler
	CarbonreportH     goagrpc.UnaryHandler
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
		CalculateReportsH: NewCalculateReportsHandler(e.CalculateReports, uh),
		GetControlPointsH: NewGetControlPointsHandler(e.GetControlPoints, uh),
		GetPowerH:         NewGetPowerHandler(e.GetPower, uh),
		GetEmissionsH:     NewGetEmissionsHandler(e.GetEmissions, uh),
		HandleRequestsH:   NewHandleRequestsHandler(e.HandleRequests, uh),
		CarbonreportH:     NewCarbonreportHandler(e.Carbonreport, uh),
	}
}

// NewCalculateReportsHandler creates a gRPC handler which serves the "calc"
// service "calculate_reports" endpoint.
func NewCalculateReportsHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeCalculateReportsResponse)
	}
	return h
}

// CalculateReports implements the "CalculateReports" method in
// calcpb.CalcServer interface.
func (s *Server) CalculateReports(ctx context.Context, message *calcpb.CalculateReportsRequest) (*calcpb.CalculateReportsResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "calculate_reports")
	ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
	resp, err := s.CalculateReportsH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*calcpb.CalculateReportsResponse), nil
}

// NewGetControlPointsHandler creates a gRPC handler which serves the "calc"
// service "get_control_points" endpoint.
func NewGetControlPointsHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeGetControlPointsRequest, EncodeGetControlPointsResponse)
	}
	return h
}

// GetControlPoints implements the "GetControlPoints" method in
// calcpb.CalcServer interface.
func (s *Server) GetControlPoints(ctx context.Context, message *calcpb.GetControlPointsRequest) (*calcpb.GetControlPointsResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "get_control_points")
	ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
	resp, err := s.GetControlPointsH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*calcpb.GetControlPointsResponse), nil
}

// NewGetPowerHandler creates a gRPC handler which serves the "calc" service
// "get_power" endpoint.
func NewGetPowerHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeGetPowerRequest, EncodeGetPowerResponse)
	}
	return h
}

// GetPower implements the "GetPower" method in calcpb.CalcServer interface.
func (s *Server) GetPower(ctx context.Context, message *calcpb.GetPowerRequest) (*calcpb.GetPowerResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "get_power")
	ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
	resp, err := s.GetPowerH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*calcpb.GetPowerResponse), nil
}

// NewGetEmissionsHandler creates a gRPC handler which serves the "calc"
// service "get_emissions" endpoint.
func NewGetEmissionsHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeGetEmissionsResponse)
	}
	return h
}

// GetEmissions implements the "GetEmissions" method in calcpb.CalcServer
// interface.
func (s *Server) GetEmissions(ctx context.Context, message *calcpb.GetEmissionsRequest) (*calcpb.GetEmissionsResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "get_emissions")
	ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
	resp, err := s.GetEmissionsH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*calcpb.GetEmissionsResponse), nil
}

// NewHandleRequestsHandler creates a gRPC handler which serves the "calc"
// service "handle_requests" endpoint.
func NewHandleRequestsHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, DecodeHandleRequestsRequest, EncodeHandleRequestsResponse)
	}
	return h
}

// HandleRequests implements the "HandleRequests" method in calcpb.CalcServer
// interface.
func (s *Server) HandleRequests(ctx context.Context, message *calcpb.HandleRequestsRequest) (*calcpb.HandleRequestsResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "handle_requests")
	ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
	resp, err := s.HandleRequestsH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*calcpb.HandleRequestsResponse), nil
}

// NewCarbonreportHandler creates a gRPC handler which serves the "calc"
// service "carbonreport" endpoint.
func NewCarbonreportHandler(endpoint goa.Endpoint, h goagrpc.UnaryHandler) goagrpc.UnaryHandler {
	if h == nil {
		h = goagrpc.NewUnaryHandler(endpoint, nil, EncodeCarbonreportResponse)
	}
	return h
}

// Carbonreport implements the "Carbonreport" method in calcpb.CalcServer
// interface.
func (s *Server) Carbonreport(ctx context.Context, message *calcpb.CarbonreportRequest) (*calcpb.CarbonreportResponse, error) {
	ctx = context.WithValue(ctx, goa.MethodKey, "carbonreport")
	ctx = context.WithValue(ctx, goa.ServiceKey, "calc")
	resp, err := s.CarbonreportH.Handle(ctx, message)
	if err != nil {
		return nil, goagrpc.EncodeError(err)
	}
	return resp.(*calcpb.CarbonreportResponse), nil
}
