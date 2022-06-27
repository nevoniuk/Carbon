// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller endpoints
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package poller

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "Poller" service endpoints.
type Endpoints struct {
	CarbonEmissions goa.Endpoint
	AggregateData   goa.Endpoint
}

// NewEndpoints wraps the methods of the "Poller" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		CarbonEmissions: NewCarbonEmissionsEndpoint(s),
		AggregateData:   NewAggregateDataEndpoint(s),
	}
}

// Use applies the given middleware to all the "Poller" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.CarbonEmissions = m(e.CarbonEmissions)
	e.AggregateData = m(e.AggregateData)
}

// NewCarbonEmissionsEndpoint returns an endpoint function that calls the
// method "carbon_emissions" of service "Poller".
func NewCarbonEmissionsEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.CarbonEmissions(ctx)
	}
}

// NewAggregateDataEndpoint returns an endpoint function that calls the method
// "aggregate_data" of service "Poller".
func NewAggregateDataEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.AggregateData(ctx)
	}
}
