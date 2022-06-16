package pollerapi

import (
	"context"
	"log"

	poller "github.com/crossnokaye/carbon/gen/poller"
)

// Poller service example implementation.
// The example methods log the requests and return zero values.
type pollersrvc struct {
	logger *log.Logger
}

// NewPoller returns the Poller service implementation.
func NewPoller(logger *log.Logger) poller.Service {
	return &pollersrvc{logger}
}

// query api getting search data for carbon_intensity event
func (s *pollersrvc) CarbonEmissions(ctx context.Context) (err error) {
	s.logger.Print("poller.carbon_emissions")
	return
}

// query api using a search call for a fuel event from Carbonara API
func (s *pollersrvc) Fuels(ctx context.Context) (err error) {
	s.logger.Print("poller.fuels")
	return
}

// get the aggregate data for an event from clickhouse
func (s *pollersrvc) AggregateData(ctx context.Context) (err error) {
	s.logger.Print("poller.aggregate_data")
	return
}
