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
	startDates []string
	minuteReports		[][]*genpoller.CarbonForecast
	
}
var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"
//var regionstartdates map[string]string
var regions [13]string
var reportdurations [5]string
	//ensurepastdata
// NewPoller returns the Poller service implementation.
func NewPoller(ctx context.Context, csc carbonara.Client, dbc storage.Client) *pollersrvc {
	ctx, cancel := context.WithCancel(ctx)
	s := &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		startDates:		    []string{},
		minuteReports:		[][]*genpoller.CarbonForecast{},
		
	}
	regions = [...]string{ "CAISO", "AESO", "BPA", "ERCO", "IESO",
       "ISONE", "MISO",
        "NYISO", "NYISO.NYCW",
         "NYISO.NYLI", "NYISO.NYUP",
          "PJM", "SPP"} 
	var regionstartdates = map[string]string{
		"CAISO":"2018-04-10T07:00:00+00:00",
		"AESO": "2020-05-15T16:00:00+00:00",
		"BPA": "2018-01-01T08:00:00+00:00",
		
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
	reportdurations = [...]string{ "minute", "hourly", "daily", "monthly", "yearly"}
	//return weird date
	times := s.Ensurepastdata(ctx, regionstartdates)
	fmt.Println(times)
	//carbonReports, err := s.CarbonEmissions(ctx, times)
	//aggregateReports, err := s.AggregateDataEndpoint(ctx, times)
	s = &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		startDates:		    times,
	}
	
	return s
}

func (s *pollersrvc) Ensurepastdata(ctx context.Context, regionstartdates map[string]string) (startDates []string) {
	//configure start dates for each region 
	var dates []string
	for i := 0; i < len(regions); i++ {
		//failing at checkDB because clickhouse connection is refused
		date, err := s.dbc.CheckDB(ctx, string(regions[i]))
		//TODO remove this
		err = fmt.Errorf("error")
		if err == nil {
			dates = append(dates, date)
		} else if err != nil {
			fmt.Println(err)
			defaultDate := regionstartdates[string(regions[i])]
			//fmt.Printf("date is %s\n", defaultDate)
			dates = append(dates, defaultDate)
		}
	}
	return dates
}

func (s *pollersrvc) Start(ctx context.Context) error {
	//1.loop through regions here
	//loop through len(regions) but for now only one
	for i := 0; i < 1; i++ {
		var payload = genpoller.CarbonPayload{Region: &regions[i], Start: &s.startDates[i]}
		minutereports, err1 := s.CarbonEmissions(ctx, &payload) //returns single array of forecasts
		if err1 != nil {
			return err1
		}
		//dates used as input for clickhouse queries to get averages
		dateConfigs, err := getdates(ctx, minutereports)

		if err != nil {
			return err
		}

		s.dbc.SaveCarbonReports(ctx, minutereports)
		//loop through hourly, weekly, monthly, yearly periods to create reports
		for j := 0; j < len(dateConfigs); j++ {
			if dateConfigs[j] != nil {
				var payload = genpoller.AggregatePayload{Region: &regions[i], Periods: dateConfigs[j]}
				s.AggregateData(ctx, &payload)
			}
			
		}
	}
	return nil
}





func (ser *pollersrvc) CarbonEmissions(ctx context.Context, input *genpoller.CarbonPayload) ([]*genpoller.CarbonForecast, error) {
	var start = *input.Start
	var region = *input.Region
	var reports []*genpoller.CarbonForecast
	lastDate := time.Now().Format(timeFormat)
	
	//search endpoint wont take requests over 7 day
	//REMOVE TEMP COUNTER
	var testcounter = 0
	//loop through every week from stand -> lastDate
	for testcounter < 5 {
		if lastDate == start {
			return reports, nil
		}
		t, err := time.Parse(timeFormat, start)
		if err != nil {
			return nil, err
		}
		newTime := t.AddDate(0, 0, 6)
		var end = newTime.Format(timeFormat)
		reports, err = ser.csc.GetEmissions(ctx, region, start, end, reports)
		
		newTime = newTime.AddDate(0, 0, 1)
		end = newTime.Format(timeFormat)
		start = end
		testcounter += 1
	}

	return reports, nil
}

// get the aggregate data for an event from clickhouse
func (ser *pollersrvc) AggregateData(ctx context.Context, input *genpoller.AggregatePayload) (error) {

	var region = *input.Region
	
	var dates = input.Periods
	//loop through period array
	aggregateres, err := ser.dbc.GetAggregateReports(ctx, dates, region, reportdurations[0])

	fmt.Println(aggregateres)
	if err != nil {
		return err
	}
	ser.dbc.SaveCarbonReports(ctx, aggregateres)

	return nil
}

//NOTE: need function to get the dates from minute reports because some data may have not been available
func getdates(ctx context.Context, minutereports []*genpoller.CarbonForecast) ([][]*genpoller.Period, error) {
	
	var initialstart, err = time.Parse(timeFormat, minutereports[0].Duration.StartTime)

	if err != nil {
		return nil, err
	}
	var finalDates [][]*genpoller.Period

	var hourlyDates []*genpoller.Period
	var dailyDates []*genpoller.Period
	var weeklyDates []*genpoller.Period
	var monthlyDates []*genpoller.Period
	var yearlyDates []*genpoller.Period

	var hourstart = initialstart
	//var hourend = initialstart

	var daystart = initialstart
	//var dayend = initialstart

	var weekstart = initialstart
	//var weekend = initialstart

	var monthstart = initialstart
	//var monthend = initialstart

	var yearstart = initialstart
	//var yearend = initialstart
	
	var previous = initialstart
	

	var hourcounter = time.Time.Hour(initialstart)
	var daycounter = time.Time.Day(initialstart)
	var weekcounter = 0
	//fmt.Printf("day counter is %d\n", daycounter)
	var monthcounter = time.Time.Month(initialstart)
	var yearcounter = time.Time.Year(initialstart)

	
	for _, event := range minutereports {
		
		var time, err = time.Parse(timeFormat, event.Duration.StartTime)
		if err != nil {
			return nil, fmt.Errorf("parsing error")
		}
	
		var year = time.Year()
		var month = time.Month()
		var day = time.Day()
		var hour = time.Hour()
		

		if hour != hourcounter {
			
			if month != monthcounter {
				
				monthcounter = month
				monthlyDates = append(monthlyDates, &genpoller.Period{monthstart.Format(timeFormat), previous.Format(timeFormat)})
				//fmt.Println("MONTHLY DATE")
				//fmt.Println(&genpoller.Period{monthstart.Format(timeFormat), previous.Format(timeFormat)})
				monthstart = time
			}

			if year != yearcounter {
				
				yearcounter = year
				yearlyDates = append(yearlyDates, &genpoller.Period{yearstart.Format(timeFormat), previous.Format(timeFormat)})
				yearstart = time
			}

			if day != daycounter {
				daycounter = day
				weekcounter += 1
				//fmt.Println("DAILY DATE")
				dailyDates = append(dailyDates, &genpoller.Period{daystart.Format(timeFormat), previous.Format(timeFormat)})
				///fmt.Println(&genpoller.Period{daystart.Format(timeFormat), previous.Format(timeFormat)})
				
				daystart = time
				//might be 8 instead
				if weekcounter == 7 {
					weeklyDates = append(weeklyDates, &genpoller.Period{weekstart.Format(timeFormat), previous.Format(timeFormat)})
					//fmt.Println("WEEKLYDATE")
					//fmt.Println(&genpoller.Period{weekstart.Format(timeFormat), previous.Format(timeFormat)})
					weekstart = time
					weekcounter = 0
				}
			}

			hourcounter = hour
			hourlyDates = append(hourlyDates, &genpoller.Period{hourstart.Format(timeFormat), previous.Format(timeFormat)})
			//fmt.Println("HOURLYDATE")
			//fmt.Println(&genpoller.Period{hourstart.Format(timeFormat), previous.Format(timeFormat)})
			//fmt.Println(&genpoller.Period{hourstart.Format(timeFormat), previous.Format(timeFormat)})
			hourstart = time
		}

		previous = time

	}
	//handle the case where only one day was returned
	if daycounter == time.Time.Day(initialstart) {
		//fmt.Println("day counter is the same")
		dailyDates = append(dailyDates, &genpoller.Period{daystart.Format(timeFormat),previous.Format(timeFormat)})
		//fmt.Println(dailyDates[0])
	}
	finalDates = append(finalDates, hourlyDates, dailyDates, weeklyDates, monthlyDates, yearlyDates)
	//fmt.Printf("DATES")
	//fmt.Println(hourlyDates)
	//fmt.Println(dailyDates)
	//fmt.Println(weeklyDates)
	//fmt.Println()
	return finalDates, nil
}


