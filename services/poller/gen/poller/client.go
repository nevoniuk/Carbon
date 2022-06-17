// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller client
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package poller

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "Poller" service client.
type Client struct {
	CarbonEmissionsEndpoint       goa.Endpoint
	FuelsEndpoint                 goa.Endpoint
	AggregateDataEndpointEndpoint goa.Endpoint
}

// NewClient initializes a "Poller" service client given the endpoints.
func NewClient(carbonEmissions, fuels, aggregateDataEndpoint goa.Endpoint) *Client {
	return &Client{
		CarbonEmissionsEndpoint:       carbonEmissions,
		FuelsEndpoint:                 fuels,
		AggregateDataEndpointEndpoint: aggregateDataEndpoint,
	}
}

// CarbonEmissions calls the "carbon_emissions" endpoint of the "Poller"
// service.
// CarbonEmissions may return the following errors:
//	- "data_not_available" (type *goa.ServiceError): The data is not available or server error
//	- "missing-required-parameter" (type *goa.ServiceError): missing-required-parameter
//	- error: internal error
func (c *Client) CarbonEmissions(ctx context.Context) (res []*CarbonForecast, err error) {
	var ires interface{}
	ires, err = c.CarbonEmissionsEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.([]*CarbonForecast), nil
}

// Fuels calls the "fuels" endpoint of the "Poller" service.
// Fuels may return the following errors:
//	- "data_not_available" (type *goa.ServiceError): The data is not available or server error
//	- "missing-required-parameter" (type *goa.ServiceError): missing-required-parameter
//	- error: internal error
func (c *Client) Fuels(ctx context.Context) (res []*FuelsForecast, err error) {
	var ires interface{}
	ires, err = c.FuelsEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.([]*FuelsForecast), nil
}

// AggregateDataEndpoint calls the "aggregate_data" endpoint of the "Poller"
// service.
// AggregateDataEndpoint may return the following errors:
//	- "data_not_available" (type *goa.ServiceError): The data is not available or server error
//	- "missing-required-parameter" (type *goa.ServiceError): missing-required-parameter
//	- error: internal error
func (c *Client) AggregateDataEndpoint(ctx context.Context) (res []*AggregateData, err error) {
	var ires interface{}
	ires, err = c.AggregateDataEndpointEndpoint(ctx, nil)
	if err != nil {
		return
	}
	return ires.([]*AggregateData), nil
}
