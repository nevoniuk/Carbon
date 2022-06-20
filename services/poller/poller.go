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
	carbonReports		[][]*genpoller.CarbonForecast
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
		carbonReports:		[][]*genpoller.CarbonForecast{},
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
	//rates poller service uses the command "go"
	times := s.ensurepastdata(ctx)
	carbonReports, err := s.CarbonEmissions(ctx)
	if err != nil {
		fmt.Printf("could not retrieve co2 emissions")
	}
	return &pollersrvc{csc: csc, dbc: dbc, ctx: ctx, cancel: cancel, readDates: times, carbonReports: carbonReports}
}

func (s *pollersrvc) ensurepastdata(ctx context.Context) (startDates []time.Time) {
	//configure start dates for each region 
	var dates []time.Time
	for i, region := range regions {
		date, err := s.dbc.CheckDB(ctx, string(region[i]))
		if err == nil {
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
func (s *pollersrvc) CarbonEmissions(ctx context.Context) (res [][]*genpoller.CarbonForecast, err error) {
	var dates = s.readDates
	var reports [][]*genpoller.CarbonForecast
	end, err := time.Parse(timeFormat, time.Now().GoString())
	if err != nil {
		fmt.Printf("time parse problem in carbon emissions")
	}
	for i, region := range regions {
		start := dates[i]

		carbonres, err := s.csc.GetEmissions(ctx, region, start, end)
        if err != nil { //handle errors when a region is not available??
            //instead of returning have a way marking that a region is not available
            //handle case when
            return nil, err
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
func (s *pollersrvc) AggregateDataEndpoint(ctx context.Context) (res []*genpoller.AggregateData, err error) {
		//timeNow, err := time.Parse(timeFormat, time.Now().GoString())
		var carbonreports = s.carbonReports
		var dates = s.readDates
		if err != nil {
			return nil, err
		}
	if carbonreports != nil {
		for i, region := range regions {
			var days []*genpoller.Period
			var months []*genpoller.Period
			var years []*genpoller.Period
			days, months, years = getdates(ctx, dates[i], carbonreports[i])
			if days != nil {
				aggregateres, err := s.dbc.GetAggregateReports(ctx, days, region, reportdurations[1])
				if err == nil {
					s.dbc.SaveAggregateReports(ctx, aggregateres)
				}
			}
			if months != nil {
				aggregateres, err := s.dbc.GetAggregateReports(ctx, months, region, reportdurations[2])
				if err == nil {
					s.dbc.SaveAggregateReports(ctx, aggregateres)
				}
			}

			if years != nil {
				aggregateres, err := s.dbc.GetAggregateReports(ctx, years, region, reportdurations[3])
				if err == nil {
					s.dbc.SaveAggregateReports(ctx, aggregateres)
				}
			}
			
			
		}
	}
	
	return nil, err
}
//TODO make function return periods for daily, weekly and monthly reports
func getdates(ctx context.Context, initialstart time.Time, hourlyreports []*genpoller.CarbonForecast) ([]*genpoller.Period, []*genpoller.Period, []*genpoller.Period) {
	//var datescounter := 0
	var dailyDates []*genpoller.Period
	var monthlyDates []*genpoller.Period
	var yearlyDates []*genpoller.Period
	var daystart time.Time
	var dayend time.Time
	var monthstart time.Time
	var monthend time.Time
	var yearstart time.Time
	var yearend time.Time


	var newreport = true

	//TODO dont know if function below returns the right values
	var daycounter = time.Time.Day(initialstart)
	var monthcounter = time.Time.Month(initialstart)
	var yearcounter = time.Time.Year(initialstart)

	daystart = initialstart
	for _, event := range hourlyreports {
		var time, err = time.Parse(timeFormat, event.Duration.StartTime)
		if err != nil {
			fmt.Errorf("parsing error")
		}
	
		var year = time.Year()
		var month = time.Month()
		var day = time.Day()
		

		if (month >= monthcounter && year >= yearcounter) || (month <= monthcounter && year >= yearcounter) {
			if day > daycounter || (day < daycounter && month == monthcounter) {
				daycounter = day
				newreport = true 
			}
		}

		if newreport == true {
			newreport = false
			if month != monthcounter {
				monthcounter = month
				monthlyDates = append(monthlyDates, &genpoller.Period{monthstart.GoString(), monthend.GoString()})
				monthstart = time
			}
			if year != yearcounter {
				yearcounter = year
				yearlyDates = append(yearlyDates, &genpoller.Period{yearstart.GoString(), yearend.GoString()})
				yearstart = time
			}
			dailyDates = append(dailyDates, &genpoller.Period{daystart.GoString(), dayend.GoString()})
			daystart = time
		}
		dayend = time
		monthend = time
		yearend = time

	}
	return dailyDates, monthlyDates, yearlyDates
}


