package pollerapi

import (
	"context"
	"errors"
	"fmt"
	"goa.design/clue/log"
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
}


//timeFormat is used to parse times in order to store time as ISO8601 format
var timeFormat = "2006-01-02T15:04:05-07:00"

//regions maintains all the regions that Singularity returns CO2 intensity for
var regions [13]string = [13]string{"CAISO", "AESO", "BPA", "ERCO", "IESO",
"ISONE", "MISO",
 "NYISO", "NYISO.NYCW",
  "NYISO.NYLI", "NYISO.NYUP",
   "PJM", "SPP"} 

//reportdurations maintains the interval length of each report
var reportdurations [5]string = [5]string{ "minute", "hourly", "daily", "weekly", "monthly"}
var regionstartdate = "2020-01-1T00:00:00+00:00"
var AesoStartDate = "2020-05-15T16:00:00+00:00"

// NewPoller returns the Poller service implementation.
func NewPoller(ctx context.Context, csc carbonara.Client, dbc storage.Client) *pollersrvc {
	ctx, cancel := context.WithCancel(ctx)
	s := &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		startDates:		    []string{},
	}

	regions = [...]string{ "CAISO", "AESO", "BPA", "ERCO", "IESO",
       "ISONE", "MISO",
        "NYISO", "NYISO.NYCW",
         "NYISO.NYLI", "NYISO.NYUP",
          "PJM", "SPP"} 
	
	times := s.EnsurePastData(ctx)
	s = &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		startDates:		    times,
	}
	return s
}

//EnsurePastData will query clickhouse for the most recent report date
func (s *pollersrvc) EnsurePastData(ctx context.Context) (startDates []string) {
	var dates []string
	for i := 0; i < len(regions); i++ {
		date, err := s.dbc.CheckDB(ctx, string(regions[i]))
		if err == nil {
			dates = append(dates, date)
		} else if err != nil {
			fmt.Println(err)
			var defaultDate string
			if regions[i] == "AESO" {
				defaultDate = AesoStartDate
			} else {defaultDate = regionstartdate}
			dates = append(dates, defaultDate)
		}
	}
	return dates
}

//Update will fetch the latest reports for all regions
/* 
Errors: server error, download error, dates error, clickhouse error
*/
//calculate days to backfill data in Newpoller

func (s *pollersrvc) Update(ctx context.Context) error {
	var timeNow = time.Now()
	for i := 0; i < len(regions); i++ {
		minutereports, err := s.CarbonEmissions(ctx, s.startDates[i], timeNow.Format(timeFormat), regions[i]) //returns single array of forecasts
		var NoDataError carbonara.NoDataError
		if !errors.As(err, NoDataError.Err) {
			return mapAndLogErrorf(ctx, "failed to get Carbon Intensity Reports:%w\n", err)
		} else if errors.As(err, NoDataError.Err) {
			continue
		}
		dateConfigs, err := GetDates(ctx, minutereports)
		if err != nil {
			log.Error(ctx, err)
			return fmt.Errorf("error in GetDates: %v\n", err)
		}
		err = s.dbc.SaveCarbonReports(ctx, minutereports)
		if err != nil {
			log.Error(ctx, err)
			return fmt.Errorf("error in SaveCarbonReports: %v\n", err)
		}
		for j := 0; j < len(dateConfigs); j++ {
			if dateConfigs[j] != nil {
				aggErr := s.AggregateData(ctx, regions[i], dateConfigs[j], reportdurations[j])
				if aggErr != nil {
					log.Error(ctx, err)
					fmt.Errorf("error from aggregate data %s\n", aggErr)
					return aggErr
				}
			}
		}
	}
	return nil
}


//R&D can use this function to obtain CO2 intensity reports for a specific region
func (ser *pollersrvc) GetEmissionsForRegion(ctx context.Context, input *genpoller.CarbonPayload) ([]*genpoller.CarbonForecast, error) {
	var start = *input.Start
	var end = *input.End
	var region = *input.Region
	var reports []*genpoller.CarbonForecast
	reports, emissionsErr := ser.csc.GetEmissions(ctx, region, start, end, reports)
	if emissionsErr != nil {
		return nil, emissionsErr
	}
	return reports, emissionsErr
}

//CarbonEmissions is a helper function to calculate API calls in 7 day intervals
//Errors: server error, download error(general error)
func (ser *pollersrvc) CarbonEmissions(ctx context.Context, start string, end string, region string) ([]*genpoller.CarbonForecast, error) {
	var reports []*genpoller.CarbonForecast
	startTime, err := time.Parse(timeFormat, start)
	if err != nil {
		return nil, err
	}
	finalendTime, err := time.Parse(timeFormat, end)
	if err != nil {
		return nil, err
	}
	for startTime.Before(finalendTime) {
		t, err := time.Parse(timeFormat, start)
		if err != nil {
			return nil, err
		}
		newEndTime := t.AddDate(0, 0, 6)
		if !newEndTime.Before(finalendTime) {
			newEndTime = finalendTime 
		}
		var emissionsErr error
		reports, emissionsErr = ser.csc.GetEmissions(ctx, region, start, newEndTime.Format(timeFormat), reports)
		if emissionsErr != nil {
			return nil, emissionsErr
		}
		newEndTime = newEndTime.AddDate(0, 0, 1)
		start = newEndTime.Format(timeFormat)
	}
	return reports, nil
}

//AggregateData gets aggregate reports for all report dates returned by GetDates and store them in clickhouse
func (ser *pollersrvc) AggregateData(ctx context.Context, region string, dates []*genpoller.Period, duration string) (error) {
	aggregateres, getErr := ser.dbc.GetAggregateReports(ctx, dates, region, duration)
	if getErr != nil {
		return getErr
	}
	saveErr := ser.dbc.SaveCarbonReports(ctx, aggregateres)
	if saveErr != nil {
		return saveErr
	}
	return nil
}

//GetDates gets all the report dates that are used as input to clickhouse queries and obtain aggregate CO2 intensity data
func GetDates(ctx context.Context, minutereports []*genpoller.CarbonForecast) ([][]*genpoller.Period, error) {
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
	//dont write any reports unless complete
//counter to maintain start of each interval
	var hourstart = initialstart
	var daystart = initialstart
	var weekstart = initialstart
	year, month, day := initialstart.Date()
	var monthstart time.Time
	if day != 1 {
		monthstart = time.Date(year, month, 1 ,0,0,0,0, initialstart.Location())
	} else {
		monthstart = initialstart
	}
	var previous = initialstart
	//could compare to previous instead and eliminate need for these counters
	//var hourcounter = time.Time.Hour(initialstart)
	//var daycounter = time.Time.Day(initialstart)
	var weekcounter = 0
	//var monthcounter = time.Time.Month(initialstart)
	for _, event := range minutereports {
		var time, err = time.Parse(timeFormat, event.Duration.StartTime)
		if err != nil {
			return nil, fmt.Errorf("parsing error")
		}
		var month = time.Month()
		var day = time.Day()
		var hour = time.Hour()
		if hour != previous.Hour() {
			if month != previous.Month() {
				//monthcounter = month
				monthlyDates = append(monthlyDates, &genpoller.Period{monthstart.Format(timeFormat), previous.Format(timeFormat)})
				monthstart = time
			}
			if day != previous.Day() {
				//daycounter = day
				weekcounter += 1
				dailyDates = append(dailyDates, &genpoller.Period{daystart.Format(timeFormat), previous.Format(timeFormat)})
				daystart = time
				if weekcounter == 7 {
					weeklyDates = append(weeklyDates, &genpoller.Period{weekstart.Format(timeFormat), previous.Format(timeFormat)})
					weekstart = time
					weekcounter = 0
				}
			}
			//hourcounter = hour
			//bug: pretty sure that start == end date
			hourlyDates = append(hourlyDates, &genpoller.Period{hourstart.Format(timeFormat), previous.Format(timeFormat)})
			hourstart = time
		}
		previous = time
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





