package pollerapi

import (
	"context"
	"fmt"
	//"fmt"
	//"sync"
	"time"

	"github.com/crossnokaye/carbon/services/poller/clients/carbonara"
	"github.com/crossnokaye/carbon/services/poller/clients/storage"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

// Poller service example implementation.
// The example methods log the requests and return zero values.
type pollersrvc struct {
	csc carbonara.Client
	dbc storage.Client
	ctx                 context.Context
	cancel context.CancelFunc
	readDates []time.Time
}
var timeFormat = "2006-01-02T15:04:05-07:00"
    var dateFormat = "2006-01-02"
var regionstartdates map[string]string
var regions []string
var reportdurations []string
	//ensurepastdata
// NewPoller returns the Poller service implementation.
func NewPoller(ctx context.Context, csc carbonara.Client, dbc storage.Client) *pollersrvc {
	ctx, cancel := context.WithCancel(ctx)
	s := &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		readDates:		    []time.Time{},
	}
	regions = []string{ "AESO", "BPA", "CAISO", "ERCO", "IESO",
       "ISONE", "MISO",
        "NYSIO", "NYISO.NYCW",
         "NYISO.NYLI", "NYISO.NYUP",
          "PJM", "SPP", "EIA"} 

	regionstartdates = map[string]string{"AESO": "2020-05-15T16:00:00+00:00",
		"BPA": "2018-01-01T08:00:00+00:00",
		"CAISO": "2018-04-10T07:00:00+00:00",
		"ERCO": "2018-07-02T05:05:00+00:00",
		"IESO": "2017-12-31T05:00:00+00:00",
		"ISONE": "2015-01-01T05:00:00+00:00",
		"MISO": "2018-01-01T05:00:00+00:00",
		"NYSIO": "2017-12-01T05:05:00+00:00",
		"NYISO.NYCW": "2019-01-01T00:00:00+00:00",
		"NYISO.NYLI": "2019-01-01T00:00:00+00:00",
		"NYISO.NYUP": "2019-01-01T00:00:00+00:00",
		"PJM": "2017-07-01T04:05:00+00:00",
		"SPP": "2017-12-31T00:00:00+00:00",
		"EIA": "2019-01-01T05:00:00+00:00"}
	reportdurations = []string{"hourly, daily, weekly, monthly"}
	times := s.ensurepastdata(ctx, dbc)
	return &pollersrvc{csc: csc, dbc: dbc, ctx: ctx, cancel: cancel, readDates: times}
}

func (s *pollersrvc) ensurepastdata(ctx context.Context, dbc storage.Client) (startDates []time.Time) {
	//configure start dates for each region 
	var dates []time.Time
	for i, region := range regions {
		date, err := s.dbc.CheckDB(ctx, string(region[i]))
		if date != nil {
			dates = append(dates, date)
		} else if err != nil {
			//default date
			//parse regionstartdate
			layout := "2019-01-01T05:00:00+00:00"
			defaultDate := regionstartdates[string(region[i])]
			newtime, err := time.Parse(layout, defaultDate) //convert string -> time
			if (err != nil) {dates = append(dates, newtime)}
			
		}
	}
	return dates
}


// query api getting search data for carbon_intensity event
func (s *pollersrvc) CarbonEmissions(ctx context.Context, dates []time.Time) ([][]*genpoller.CarbonForecast, error) {
	var reports [][]*genpoller.CarbonForecast
	end, err := time.Parse(timeFormat, time.Now().GoString())
	if err != nil {
		fmt.Printf("time parse problem in carbon emissions")
	}
	for i, region := range regions {
		start := dates[i]
		carbonres, err := s.csc.getemissions(ctx, region, start, end)
        if err != nil { //handle errors when a region is not available??
            //instead of returning have a way marking that a region is not available
            //handle case when
            return err
		}
		reports = append(reports, carbonres)
		s.dbc.SaveCarbonReports(ctx, carbonres)
	}
	return reports, nil
}

// query api using a search call for a fuel event from Carbonara API
func (s *pollersrvc) Fuels(ctx context.Context, dates []time.Time) (err error) {
	return
}

// get the aggregate data for an event from clickhouse
func (s *pollersrvc) AggregateData(ctx context.Context, carbonreports [][]*genpoller.CarbonForecast,
	 fuelsreports [][]*genpoller.FuelsForecast) (err error) {
		timeNow, err := time.Parse(timeFormat, time.Now().GoString())
		if err != nil {
			return err
		}
	if carbonreports != nil {
		for i, region := range regions {
			var days []*genpoller.Period
			var weeks []*genpoller.Period
			var months []*genpoller.Period
			weeks = 
			start = dates[i]
			for j, reporttype := range reportdurations {
				aggregateres, err := s.dbc.get_aggregate_data(ctx, reporttype[j], start, end, region)
				s.dbc.SaveAggregateData(ctx, aggregateres)
			}
		}
	}
	
	return nil
}
func getdailydates(hourstartTimes []time.Time{}, carbonreports) ([]*genpoller.Period) {
	return nil
}
func getweeklydates(dailystartTimes []time.Time{}, carbonreports) ([]*genpoller.Period) {
	return nil
}
func getmonthlydates(monthlystartTimes []time.Time{}, carbonreports) ([]*genpoller.Period) {
	return nil
} 
//helper functions
