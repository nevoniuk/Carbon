// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller service
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package poller

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Service that provides forecasts to clickhouse from Carbonara API
type Service interface {
	// query api getting search data for carbon_intensity event
	CarbonEmissions(context.Context, []string) (res [][]*CarbonForecast, err error)
	// get the aggregate data for an event from clickhouse
	AggregateDataEndpoint(context.Context) (res [][]*AggregateData, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "Poller"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"carbon_emissions", "aggregate_data"}

type AggregateData struct {
	// average
	Average float64
	// min
	Min float64
	// max
	Max float64
	// sum
	Sum float64
	// count
	Count int
	// duration
	Duration *Period
	// report_type
	ReportType string
}

// Emissions Forecast
type CarbonForecast struct {
	// generated_rate
	GeneratedRate float64
	// marginal_rate
	MarginalRate float64
	// consumed_rate
	ConsumedRate float64
	// Duration
	Duration *Period
	// generated_source
	GeneratedSource string
	// region
	Region string
}

// Period of time from start to end of Forecast
type Period struct {
	// Start time
	StartTime string
	// End time
	EndTime string
}

// MakeMissingRequiredParameter builds a goa.ServiceError from an error.
func MakeMissingRequiredParameter(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "missing-required-parameter",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
