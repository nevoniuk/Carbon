// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller endpoints
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design -o
// services/poller

package poller

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Endpoints wraps the "Poller" service endpoints.
type Endpoints struct {
	Update                goa.Endpoint
	GetEmissionsForRegion goa.Endpoint
}

// NewEndpoints wraps the methods of the "Poller" service with endpoints.
func NewEndpoints(s Service) *Endpoints {
	return &Endpoints{
		Update:                NewUpdateEndpoint(s),
		GetEmissionsForRegion: NewGetEmissionsForRegionEndpoint(s),
	}
}

// Use applies the given middleware to all the "Poller" service endpoints.
func (e *Endpoints) Use(m func(goa.Endpoint) goa.Endpoint) {
	e.Update = m(e.Update)
	e.GetEmissionsForRegion = m(e.GetEmissionsForRegion)
}

// NewUpdateEndpoint returns an endpoint function that calls the method
// "update" of service "Poller".
func NewUpdateEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, s.Update(ctx)
	}
}

// NewGetEmissionsForRegionEndpoint returns an endpoint function that calls the
// method "get_emissions_for_region" of service "Poller".
func NewGetEmissionsForRegionEndpoint(s Service) goa.Endpoint {
	return func(ctx context.Context, req interface{}) (interface{}, error) {
		p := req.(*CarbonPayload)
		return s.GetEmissionsForRegion(ctx, p)
	}
}
