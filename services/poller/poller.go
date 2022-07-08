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
		"CAISO":"2020-01-1T00:00:00+00:00",
		"AESO": "2020-05-15T16:00:00+00:00",
		"BPA": "2020-01-1T00:00:00+00:00",
		
		"ERCO" : "2020-01-1T00:00:00+00:00",
		"IESO": "2020-01-1T00:00:00+00:00",
		"ISONE": "2020-01-1T00:00:00+00:00",
		"MISO": "2020-01-1T00:00:00+00:00",
		"NYISO":"2020-01-1T00:00:00+00:00",
		"NYISO.NYCW": "2020-01-1T00:00:00+00:00", //wont add to array
		"NYISO.NYLI": "2020-01-1T00:00:00+00:00",
		"NYISO.NYUP": "2020-01-1T00:00:00+00:00",
		"PJM": "2020-01-1T00:00:00+00:00",
		"SPP": "2020-01-1T00:00:00+00:00",
	}

	reportdurations = [...]string{ "minute", "hourly", "daily", "weekly", "monthly"}
	
	times := s.Ensurepastdata(ctx, regionstartdates)
	fmt.Println("READ DATES")
	fmt.Println(times)
	
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
	//configure dates for each region where data is last available to make API call
	//if no data is found for a region, return a default date from regionstartdates
	var dates []string

	for i := 0; i < len(regions); i++ {
		
		date, err := s.dbc.CheckDB(ctx, string(regions[i]))
		if err == nil {
			dates = append(dates, date)
		} else if err != nil {
			fmt.Println(err)
			defaultDate := regionstartdates[string(regions[i])]
			dates = append(dates, defaultDate)
		}
	}
	return dates
}

func (s *pollersrvc) Start(ctx context.Context) error {
	
	for i := 0; i < len(regions); i++ {

		var payload = genpoller.CarbonPayload{Region: &regions[i], Start: &s.startDates[i]}
		minutereports, emissionsError := s.CarbonEmissions(ctx, &payload) //returns single array of forecasts
		if emissionsError != nil {
			fmt.Errorf("Error from carbon emissions %s\n", emissionsError)
			//keep reading
			//return emissionsError
			continue
		}
		//dates used as input for clickhouse queries to get averages
		
		dateConfigs, datesErr := getdates(ctx, minutereports)

		if datesErr != nil {
			return datesErr
		}
		fmt.Println("5 MINUTE REPORTS:")
		fmt.Println(len(minutereports))
		s.dbc.SaveCarbonReports(ctx, minutereports)

		//now create hourly, weekly, monthly averages
		var durationsCounter = 1
		for j := 0; j < len(dateConfigs); j++ {
			if dateConfigs[j] != nil {
				var payload = genpoller.AggregatePayload{Region: &regions[i], Periods: dateConfigs[j], Duration: &reportdurations[durationsCounter]}
				aggErr := s.AggregateData(ctx, &payload)
				if aggErr != nil {
					fmt.Errorf("Error from aggregate data %s\n", aggErr)
					return aggErr
				}
				durationsCounter += 1
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
	
	//loop through every week from stand -> lastDate

	var startTime, err1 = time.Parse(timeFormat, start)

	if err1 != nil {
		return nil, err1
	}

	var endTime, err2 = time.Parse(timeFormat, lastDate)

	if err2 != nil {
		return nil, err2
	}

	for startTime.Before(endTime) {
		fmt.Printf("REGION IS %s\n", region)
		fmt.Printf("START TIME IS %s\n", start)
		
		t, err := time.Parse(timeFormat, start)

		if err != nil {
			return nil, err
		}

		newTime := t.AddDate(0, 0, 6)
		//exceeded last date
		if !newTime.Before(endTime) {
			return reports, nil
		}

		var tempend = newTime.Format(timeFormat)
		fmt.Printf("END TIME IS %s\n", tempend)
		var emissionsErr error
		reports, emissionsErr = ser.csc.GetEmissions(ctx, region, start, tempend, reports)
		if emissionsErr != nil {
			return nil, emissionsErr
		}
		newTime = newTime.AddDate(0, 0, 1)
		
		start = newTime.Format(timeFormat)
		//testcounter += 1
	}

	return reports, nil
}

// get the aggregate data for an event from clickhouse
func (ser *pollersrvc) AggregateData(ctx context.Context, input *genpoller.AggregatePayload) (error) {

	var region = *input.Region
	var duration = *input.Duration
	var dates = input.Periods
	//loop through period array
	aggregateres, geterr := ser.dbc.GetAggregateReports(ctx, dates, region, duration)

	fmt.Println(aggregateres)
	if geterr != nil {
		fmt.Errorf("Error from get carbon reports: %s\n", geterr)
		return geterr
	}
	saveerr := ser.dbc.SaveCarbonReports(ctx, aggregateres)
	if saveerr != nil {
		fmt.Errorf("Error from save carbon reports: %s\n", saveerr)
		return saveerr
	}

	return nil
}

//NOTE: need function to get the dates from minute reports because some data may have not been available
func getdates(ctx context.Context, minutereports []*genpoller.CarbonForecast) ([][]*genpoller.Period, error) {
	
	if minutereports == nil {
		var err = fmt.Errorf("no reports for get dates")
		return nil, err
	}

	var initialstart, err = time.Parse(timeFormat, minutereports[0].Duration.StartTime)

	if err != nil {
		return nil, err
	}
	var finalDates [][]*genpoller.Period

	var hourlyDates []*genpoller.Period
	var dailyDates []*genpoller.Period
	var weeklyDates []*genpoller.Period
	var monthlyDates []*genpoller.Period
	

//will always be the begginning of the hour
	var hourstart = initialstart
	
//will always be the beginning of the day
	var daystart = initialstart
	
//will always be the beginning of the week
	var weekstart = initialstart

	year, month, day := initialstart.Date()

	var monthstart time.Time
	//adjust to be the beginning of the month
	if day != 1 {
		monthstart = time.Date(year, month, 1 ,0,0,0,0, initialstart.Location())
	} else {
		monthstart = initialstart
	}

	var previous = initialstart
	

	var hourcounter = time.Time.Hour(initialstart)
	var daycounter = time.Time.Day(initialstart)
	var weekcounter = 0
	
	var monthcounter = time.Time.Month(initialstart)
	

	
	for _, event := range minutereports {
		
		var time, err = time.Parse(timeFormat, event.Duration.StartTime)
		if err != nil {
			return nil, fmt.Errorf("parsing error")
		}
	
		
		var month = time.Month()
		var day = time.Day()
		var hour = time.Hour()
		

		if hour != hourcounter {
			
			if month != monthcounter {
				
				monthcounter = month
				monthlyDates = append(monthlyDates, &genpoller.Period{monthstart.Format(timeFormat), previous.Format(timeFormat)})
				monthstart = time
			}

			if day != daycounter {
				daycounter = day
				weekcounter += 1
				dailyDates = append(dailyDates, &genpoller.Period{daystart.Format(timeFormat), previous.Format(timeFormat)})
				
				daystart = time
				
				if weekcounter == 7 {
					weeklyDates = append(weeklyDates, &genpoller.Period{weekstart.Format(timeFormat), previous.Format(timeFormat)})
					weekstart = time
					weekcounter = 0
				}
			}

			hourcounter = hour
			hourlyDates = append(hourlyDates, &genpoller.Period{hourstart.Format(timeFormat), previous.Format(timeFormat)})
			hourstart = time
		}

		previous = time

	}

	
	if daycounter == time.Time.Day(initialstart) {
		
		dailyDates = append(dailyDates, &genpoller.Period{daystart.Format(timeFormat),previous.Format(timeFormat)})
	}
	
	finalDates = append(finalDates, hourlyDates)
	finalDates = append(finalDates, dailyDates)
	finalDates = append(finalDates, weeklyDates)
	finalDates = append(finalDates, monthlyDates)
	

	for _, date := range finalDates {
		fmt.Println("DATE")
		fmt.Println(date)
	}
	return finalDates, nil
}




