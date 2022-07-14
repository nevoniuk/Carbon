// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller service
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design -o
// services/poller

package poller

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Service that provides forecasts to clickhouse from Carbonara API
type Service interface {
	// query Singularity's search endpoint and convert 5 min interval reports into
	// averages
	Update(context.Context) (err error)
	// query search endpoint for a region.
	GetEmissionsForRegion(context.Context, *CarbonPayload) (res []*CarbonForecast, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "Poller"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [2]string{"update", "get_emissions_for_region"}

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
	// duration_type
	DurationType string
	// region
	Region string
}

// CarbonPayload is the payload type of the Poller service
// get_emissions_for_region method.
type CarbonPayload struct {
	// region
	Region *string
	// start
	Start *string
	// end
	End *string
}

// Period of time from start to end of Forecast
type Period struct {
	// Start time
	StartTime string
	// End time
	EndTime string
}

// MakeServerError builds a goa.ServiceError from an error.
func MakeServerError(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "server_error",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeNoData builds a goa.ServiceError from an error.
func MakeNoData(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "no_data",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeRegionNotFound builds a goa.ServiceError from an error.
func MakeRegionNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "region_not_found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
