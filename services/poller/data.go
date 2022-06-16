package dataserviceapi

import (
	"context"
	"log"

	data "github.com/crossnokaye/carbon/gen/data"
)

// Data service example implementation.
// The example methods log the requests and return zero values.
type datasrvc struct {
	logger *log.Logger
}

// NewData returns the Data service implementation.
func NewData(logger *log.Logger) data.Service {
	return &datasrvc{logger}
}

// query api getting search data for carbon_intensity event
func (s *datasrvc) CarbonEmissions(ctx context.Context) (err error) {
	s.logger.Print("data.carbon_emissions")
	return
}

// query api using a search call for a fuel event from Carbonara API
func (s *datasrvc) Fuels(ctx context.Context) (err error) {
	s.logger.Print("data.fuels")
	return
}

// get the aggregate data for an event from clickhouse
func (s *datasrvc) AggregateData(ctx context.Context) (err error) {
	s.logger.Print("data.aggregate_data")
	return
}
