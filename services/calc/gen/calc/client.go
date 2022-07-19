// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc client
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design

package calc

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Client is the "calc" service client.
type Client struct {
	HistoricalCarbonEmissionsEndpoint goa.Endpoint
}

// NewClient initializes a "calc" service client given the endpoints.
func NewClient(historicalCarbonEmissions goa.Endpoint) *Client {
	return &Client{
		HistoricalCarbonEmissionsEndpoint: historicalCarbonEmissions,
	}
}

// HistoricalCarbonEmissions calls the "Historical_Carbon_Emissions" endpoint
// of the "calc" service.
func (c *Client) HistoricalCarbonEmissions(ctx context.Context, p *RequestPayload) (res *AllReports, err error) {
	var ires interface{}
	ires, err = c.HistoricalCarbonEmissionsEndpoint(ctx, p)
	if err != nil {
		return
	}
	return ires.(*AllReports), nil
}
