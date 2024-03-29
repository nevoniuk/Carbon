// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Calc service
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design -o services/calc

package calc

import (
	"context"

	goa "goa.design/goa/v3/pkg"
)

// Service to interpret CO2 emissions through KW and carbon intensity data.
// Offers the endpoint Historical Carbon Emissions
type Service interface {
	// This endpoint is used by a front end service to return carbon emission
	// reports
	HistoricalCarbonEmissions(context.Context, *RequestPayload) (res *AllReports, err error)
}

// ServiceName is the name of the service as defined in the design. This is the
// same value that is set in the endpoint request contexts under the ServiceKey
// key.
const ServiceName = "Calc"

// MethodNames lists the service method names as defined in the design. These
// are the same values that are set in the endpoint request contexts under the
// MethodKey key.
var MethodNames = [1]string{"historical_carbon_emissions"}

// AllReports is the result type of the Calc service
// historical_carbon_emissions method.
type AllReports struct {
	// TotalEmissionReport
	TotalEmissionReport *EmissionsReport
	// CarbonIntensityReports
	CarbonIntensityReports *CarbonReport
	// PowerReports
	PowerReports *ElectricalReport
}

// Carbon Report from clickhouse
type CarbonReport struct {
	// Values are in units of (lbs of CO2/MWh)
	IntensityPoints []*DataPoint
	// Duration
	Duration *Period
	Interval string
	Region   string
}

// Contains carbon emissions in terms of DataPoints, which can be used as
// points for a time/CO2 emissions graph
type DataPoint struct {
	// Time
	Time string
	// either a carbon footprint(lbs of Co2) in a CarbonEmissions struct or power
	// stamp(KW) in an Electrical Report
	Value float64
}

// Energy Generation Report from the Past values function GetValues
type ElectricalReport struct {
	// Duration
	Duration *Period
	// Power meter data in KWh
	PowerStamps []*DataPoint
	Interval    string
}

// Carbon/Energy Generation Report
type EmissionsReport struct {
	// Duration
	Duration *Period
	Interval string
	// Points
	Points []*DataPoint
	// OrgID
	OrgID UUID
	// FacilityID
	FacilityID UUID
	// LocationID
	LocationID UUID
	Region     string
}

// Period of time from start to end for any report type
type Period struct {
	// Start time
	StartTime string
	// End time
	EndTime string
}

// RequestPayload is the payload type of the Calc service
// historical_carbon_emissions method.
type RequestPayload struct {
	// OrgID
	OrgID UUID
	// Duration
	Duration *Period
	// FacilityID
	FacilityID UUID
	Interval   string
	// LocationID
	LocationID UUID
}

// Universally unique identifier
type UUID string

// MakeReportsNotFound builds a goa.ServiceError from an error.
func MakeReportsNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "reports_not_found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}

// MakeFacilityNotFound builds a goa.ServiceError from an error.
func MakeFacilityNotFound(err error) *goa.ServiceError {
	return &goa.ServiceError{
		Name:    "facility_not_found",
		ID:      goa.NewErrorID(),
		Message: err.Error(),
	}
}
