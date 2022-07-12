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
	UpdateEndpoint                goa.Endpoint
	GetEmissionsForRegionEndpoint goa.Endpoint
}

// NewClient initializes a "Poller" service client given the endpoints.
func NewClient(update, getEmissionsForRegion goa.Endpoint) *Client {
	return &Client{
		UpdateEndpoint:                update,
		GetEmissionsForRegionEndpoint: getEmissionsForRegion,
	}
}

// Update calls the "update" endpoint of the "Poller" service.
// Update may return the following errors:
//	- "server_error" (type *goa.ServiceError): Error with Singularity Server.
//	- error: internal error
func (c *Client) Update(ctx context.Context) (err error) {
	_, err = c.UpdateEndpoint(ctx, nil)
	return
}

// GetEmissionsForRegion calls the "get_emissions_for_region" endpoint of the
// "Poller" service.
// GetEmissionsForRegion may return the following errors:
//	- "server_error" (type *goa.ServiceError): Error with Singularity Server.
//	- "no_data" (type *goa.ServiceError): No new data available for any region
//	- "region_not_found" (type *goa.ServiceError): The given region is not represented by Singularity
//	- error: internal error
func (c *Client) GetEmissionsForRegion(ctx context.Context, p *CarbonPayload) (res []*CarbonForecast, err error) {
	var ires interface{}
	ires, err = c.GetEmissionsForRegionEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.([]*CarbonForecast), nil
}
