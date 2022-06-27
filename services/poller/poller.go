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
	readDates []string
	carbonReports		[][]*genpoller.CarbonForecast
}
var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"
//var regionstartdates map[string]string
var regions [13]string
var reportdurations [3]string
	//ensurepastdata
// NewPoller returns the Poller service implementation.
func NewPoller(ctx context.Context, csc carbonara.Client, dbc storage.Client) *pollersrvc {
	ctx, cancel := context.WithCancel(ctx)
	s := &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		readDates:		    []string{},
		carbonReports:		[][]*genpoller.CarbonForecast{},
	}
	regions = [...]string{ "AESO", "BPA", "CAISO", "ERCO", "IESO",
       "ISONE", "MISO",
        "NYISO", "NYISO.NYCW",
         "NYISO.NYLI", "NYISO.NYUP",
          "PJM", "SPP"} 
	var regionstartdates = map[string]string{
		
		"AESO": "2020-05-15T16:00:00+00:00",
		"BPA": "2018-01-01T08:00:00+00:00",
		"CAISO":"2018-04-10T07:00:00+00:00",
		"ERCO" : "2018-07-02T05:05:00+00:00",
		"IESO": "2017-12-31T05:00:00+00:00",
		"ISONE": "2015-01-01T05:00:00+00:00",
		"MISO": "2018-01-01T05:00:00+00:00",
		"NYISO": "2017-12-01T05:05:00+00:00",
		"NYISO.NYCW": "2019-01-01T00:00:00+00:00", //wont add to array
		"NYISO.NYLI": "2019-01-01T00:00:00+00:00",
		"NYISO.NYUP": "2019-01-01T00:00:00+00:00",
		"PJM": "2017-07-01T04:05:00+00:00",
		"SPP": "2017-12-31T00:00:00+00:00",
	}
	reportdurations = [...]string{"daily", "weekly", "monthly"}
	//rates poller service uses the command "go"
	times := s.Ensurepastdata(ctx, regionstartdates)
	//carbonReports, err := s.CarbonEmissions(ctx, times)
	//aggregateReports, err := s.AggregateDataEndpoint(ctx, times)
	s = &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		readDates:		    times,
	}
	
	return s
}

func (s *pollersrvc) Start(ctx context.Context) error {
	err1 := s.CarbonEmissions(ctx)
	if err1 != nil {
		return err1
	}
	err2 := s.AggregateDataEndpoint(ctx)
	return err2
}
func (s *pollersrvc) Ensurepastdata(ctx context.Context, regionstartdates map[string]string) (startDates []string) {
	//configure start dates for each region 
	var dates []string
	for i := 0; i < len(regions); i++ {
		//failing at checkDB because clickhouse connection is refused
		date, err := s.dbc.CheckDB(ctx, string(regions[i]))
		if err == nil {
			dates = append(dates, date)
		} else if err != nil {
			defaultDate := regionstartdates[string(regions[i])]
			//fmt.Printf("date is %s\n", defaultDate)
			if (err != nil) {dates = append(dates, defaultDate)}
		}
	}
	return dates
}

// query api getting search data for carbon_intensity event
func (ser *pollersrvc) CarbonEmissions(ctx context.Context) (error) {
	var dates = ser.readDates
	var reports [][]*genpoller.CarbonForecast
	end, err := time.Parse(timeFormat, time.Now().GoString())
	if err != nil {
		fmt.Printf("time parse problem in carbon emissions")
	}
	//TODO for testing only V
	for i := 0; i < 2; i++ {
		fmt.Printf("region is %s\n", regions[i])
		carbonres, err := ser.csc.GetEmissions(ctx, regions[i], dates[i], end.GoString())
        if err != nil {
            return err
		}
		if carbonres != nil {
			reports = append(reports, carbonres)
			ser.dbc.SaveCarbonReports(ctx, carbonres)	
		}
	}
	if reports == nil {
		return fmt.Errorf("No reports found")
	}
	ser.carbonReports = reports
	return nil
}

// query api using a search call for a fuel event from Carbonara API
func (s *pollersrvc) Fuels(ctx context.Context, dates []time.Time) (err error) {
	return
}

// get the aggregate data for an event from clickhouse
func (ser *pollersrvc) AggregateDataEndpoint(ctx context.Context) (error) {

	var carbonreports = ser.carbonReports
	var dates = ser.readDates

	//var dates = s.readDates //the start dates for each carbon report per region
	
	if carbonreports != nil {
		//TODO: change 
		for i := 0; i < 2; i++ {

			var days []*genpoller.Period
			var months []*genpoller.Period
			var years []*genpoller.Period

			var initialstart, err = time.Parse(timeFormat, dates[i])

			if err != nil {
				fmt.Errorf("error parsing time")
			}

			days, months, years = getdates(ctx, initialstart, carbonreports[i])

			if days != nil {

				aggregateres, err := ser.dbc.GetAggregateReports(ctx, days, regions[i], reportdurations[0])
				
				fmt.Println("AGGREGATE REPORTS")
				fmt.Println(aggregateres)
				if err != nil {
					return err
				}
				ser.dbc.SaveAggregateReports(ctx, aggregateres)

			}
			if months != nil {
				aggregateres, err := ser.dbc.GetAggregateReports(ctx, months, regions[i], reportdurations[1])
				if err != nil {
					return err
				}
				ser.dbc.SaveAggregateReports(ctx, aggregateres)
			}

			if years != nil {
				aggregateres, err := ser.dbc.GetAggregateReports(ctx, years, regions[i], reportdurations[2])
				if err != nil {
					return err
				}
				ser.dbc.SaveAggregateReports(ctx, aggregateres)
			}
			
			
		}
	}
	return nil
}

func getdates(ctx context.Context, initialstart time.Time, hourlyreports []*genpoller.CarbonForecast) ([]*genpoller.Period, []*genpoller.Period, []*genpoller.Period) {
	
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

	
	var daycounter = time.Time.Day(initialstart)
	var monthcounter = time.Time.Month(initialstart)
	var yearcounter = time.Time.Year(initialstart)

	daystart = initialstart

	for _, event := range hourlyreports {
		fmt.Printf("here\n")
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


