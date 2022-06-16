// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Data endpoints
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package data

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "Data" service endpoints.
type Endpoints struct {
	CarbonEmissions goa.Endpoint
	Fuels           goa.Endpoint
	AggregateData   goa.Endpoint
}

// NewEndpoints wraps the methods of the "Data" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		CarbonEmissions: NewCarbonEmissionsEndpoint(s),
		Fuels:           NewFuelsEndpoint(s),
		AggregateData:   NewAggregateDataEndpoint(s),
	}
}

// Use applies the given middleware to all the "Data" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.CarbonEmissions = m(e.CarbonEmissions)
	e.Fuels = m(e.Fuels)
	e.AggregateData = m(e.AggregateData)
}

// NewCarbonEmissionsEndpoint returns an endpoint function that calls the
// method "carbon_emissions" of service "Data".
func NewCarbonEmissionsEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.CarbonEmissions(ctx)
	}
}

// NewFuelsEndpoint returns an endpoint function that calls the method "fuels"
// of service "Data".
func NewFuelsEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.Fuels(ctx)
	}
}

// NewAggregateDataEndpoint returns an endpoint function that calls the method
// "aggregate_data" of service "Data".
func NewAggregateDataEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.AggregateData(ctx)
	}
}
