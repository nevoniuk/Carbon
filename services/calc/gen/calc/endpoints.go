// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc endpoints
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design

package calc

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "calc" service endpoints.
type Endpoints struct {
	HandleRequests  goa.Endpoint
	GetCarbonReport goa.Endpoint
}

// NewEndpoints wraps the methods of the "calc" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		HandleRequests:  NewHandleRequestsEndpoint(s),
		GetCarbonReport: NewGetCarbonReportEndpoint(s),
	}
}

// Use applies the given middleware to all the "calc" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.HandleRequests = m(e.HandleRequests)
	e.GetCarbonReport = m(e.GetCarbonReport)
}

// NewHandleRequestsEndpoint returns an endpoint function that calls the method
// "handle_requests" of service "calc".
func NewHandleRequestsEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*RequestPayload)
		return s.HandleRequests(ctx, p)
	}
}

// NewGetCarbonReportEndpoint returns an endpoint function that calls the
// method "get_carbon_report" of service "calc".
func NewGetCarbonReportEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.GetCarbonReport(ctx)
	}
}
