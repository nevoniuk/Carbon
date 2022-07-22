package pollerapi

import (
	"context"
	"errors"
	"fmt"
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
var regionstartdate = "2020-01-01T00:00:00+00:00"
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
		} else {
			var defaultDate string
			if regions[i] == "AESO" {
				defaultDate = AesoStartDate
			} else {defaultDate = regionstartdate}
			dates = append(dates, defaultDate)
		}
	}
	return dates
}

//Update will fetch the latest reports for all regions and return either a server or no-data error
func (s *pollersrvc) Update(ctx context.Context) error {
	var finalEndTime = time.Now()
	for i := 0; i < len(regions); i++ {
		startTime, err := time.Parse(timeFormat, s.startDates[i])
		if err != nil {
			return fmt.Errorf("timestamp %s in observation %v could not be parsed into a time correctly error: [%s]",s.startDates[i],
			i, err)
		}
		for startTime.Before(finalEndTime) {
			newEndTime := startTime.AddDate(0, 0, 6)
			if !newEndTime.Before(finalEndTime) {
				newEndTime = finalEndTime 
			}
			minreports, err := s.csc.GetEmissions(ctx, regions[i], startTime.Format(timeFormat), newEndTime.Format(timeFormat))
			var NoDataError carbonara.NoDataError
			if err != nil {
				if !errors.As(err, &NoDataError) {
					return mapAndLogErrorf(ctx, "failed to get Carbon Intensity Reports:%w\n", err)
				}
			}
			dateConfigs, err := GetDates(ctx, minreports)
			if err != nil {
				log.Error(ctx, err)
				continue
			}
			err = s.dbc.SaveCarbonReports(ctx, minreports)
			if err != nil {
				return mapAndLogErrorf(ctx, "failed to Save Carbon Reports:%w\n", err)
			}
			for j := 0; j < len(dateConfigs); j++ {
				if dateConfigs[j] != nil {
					aggErr := s.AggregateData(ctx, regions[i], dateConfigs[j], reportdurations[j])
					if aggErr != nil {
						return mapAndLogErrorf(ctx,  "failed to get Average Carbon Reports:%w\n", aggErr)
					}
				}
			}
			newEndTime = newEndTime.AddDate(0, 0, 1)
			startTime = newEndTime
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
	reports, err := ser.csc.GetEmissions(ctx, region, start, end)
	var serverError carbonara.ServerError
	if err != nil {
		if !errors.As(err, &serverError) {
			return nil, mapAndLogErrorf(ctx, "failed to get Carbon Intensity Reports:%w\n", err)
		}
		return nil, fmt.Errorf("error in GetEmissionsForRegion: %v\n", err)
	}
	return reports, err
}

//CarbonEmissions is a helper function to calculate API calls in 7 day intervals
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
		t, errr := time.Parse(timeFormat, start)
		if errr != nil {
			return nil, errr
		}
		newEndTime := t.AddDate(0, 0, 6)
		if !newEndTime.Before(finalendTime) {
			newEndTime = finalendTime 
		}
		var err error
		reports, err = ser.csc.GetEmissions(ctx, region, start, newEndTime.Format(timeFormat))
		if err != nil {
			return nil, mapAndLogError(ctx, err)
		}
		err = ser.dbc.SaveCarbonReports(ctx, reports)
		if err != nil {
			return nil, mapAndLogError(ctx, err)
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
		return nil, fmt.Errorf("no reports for get dates")
	}
	var initialstart, err = time.Parse(timeFormat, minutereports[0].Duration.StartTime)
	if err != nil {
		return nil, fmt.Errorf("timestamp %s in observation %v could not be parsed into a time correctly error: [%s]", minutereports[0].Duration.StartTime,
		0, err)
	}
	var finalDates [][]*genpoller.Period
	var hourlyDates []*genpoller.Period
	var dailyDates []*genpoller.Period
	var weeklyDates []*genpoller.Period
	var monthlyDates []*genpoller.Period
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
	var weekcounter = 0
	for i, event := range minutereports {
		var time, err = time.Parse(timeFormat, event.Duration.StartTime)
		if err != nil {
			return nil, fmt.Errorf("timestamp %s in observation %v could not be parsed into a time correctly error: [%s]", event.Duration.StartTime,
			i, err)
		}
		var month = time.Month()
		var day = time.Day()
		var hour = time.Hour()
		if hour != previous.Hour() {
			if month != previous.Month() || (time.AddDate(0, 0, 1).Month() != time.Month()) { 
				monthlyDates = append(monthlyDates, &genpoller.Period{monthstart.Format(timeFormat), previous.Format(timeFormat)})
				monthstart = time
			}
			if day != previous.Day() {
				weekcounter += 1
				dailyDates = append(dailyDates, &genpoller.Period{daystart.Format(timeFormat), previous.Format(timeFormat)})
				daystart = time
				if weekcounter == 7 {
					weeklyDates = append(weeklyDates, &genpoller.Period{weekstart.Format(timeFormat), previous.Format(timeFormat)})
					weekstart = time
					weekcounter = 0
				}
			}
			hourlyDates = append(hourlyDates, &genpoller.Period{hourstart.Format(timeFormat), previous.Format(timeFormat)})
			hourstart = time
		}
		previous = time
	}
	finalDates = append(finalDates, hourlyDates)
	finalDates = append(finalDates, dailyDates)
	finalDates = append(finalDates, weeklyDates)
	finalDates = append(finalDates, monthlyDates)
	return finalDates, nil
}





