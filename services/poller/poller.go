package pollerapi

import (
	"context"
	"errors"
	"fmt"
	"time"
	"goa.design/clue/log"
	"github.com/crossnokaye/carbon/services/poller/clients/carbonara"
	"github.com/crossnokaye/carbon/services/poller/clients/storage"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	"github.com/crossnokaye/carbon/model"
)
type pollersrvc struct {
	csc carbonara.Client
	dbc storage.Client
	ctx                 context.Context
	cancel context.CancelFunc
	startDates []string
}
// timeFormat is used to parse times in order to store time as ISO8601 format
var timeFormat = "2006-01-02T15:04:05-07:00"

// regions maintains all the valid regions that Singularity will calculate carbon intensity
var regions [13]string = [13]string{model.Caiso, model.Aeso, model.Bpa, model.Erco, model.Ieso,
model.Isone, model.Miso,
 model.Nyiso, model.Nyiso_nycw,
  model.Nyiso_nyli, model.Nyiso_nyup,
   model.Pjm, model.Spp} 
// reportdurations maintains the interval length of each report using constants from the model directory
var reportdurations [5]string = [5]string{ model.Minute, model.Hourly, model.Hourly, model.Weekly, model.Monthly}
// common start date for regions
const regionstartdate = "2020-01-01T00:00:00+00:00"
// The AESO region start date is earlier than other region start dates
const AesoStartDate = "2020-05-15T16:00:00+00:00"

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
	times := s.EnsurePastData(ctx)
	s = &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		startDates:		    times,
	}
	// kronjob only needs to create a 'NewPoller' instead of running Update
	end := time.Now().Format(timeFormat)
	for i := 0; i < len(regions); i++ {
		payload := &genpoller.UpdatePayload{StartTime: times[i], EndTime: end, Region: regions[i]}
		s.Update(ctx, payload)
	}
	return s
}

// EnsurePastData will query clickhouse for the most recent report date
func (s *pollersrvc) EnsurePastData(ctx context.Context) (startDates []string) {
	var dates []string
	for i := 0; i < len(regions); i++ {
		date, err := s.dbc.CheckDB(ctx, string(regions[i]))
		if err == nil {
			dates = append(dates, date)
		} else {
			var defaultDate string
			if regions[i] == model.Aeso {
				defaultDate = AesoStartDate
			} else {defaultDate = regionstartdate}
			dates = append(dates, defaultDate)
		}
	}
	return dates
}

// Update will fetch the latest reports for all regions and return either a server or no-data error
func (s *pollersrvc) Update(ctx context.Context, payload *genpoller.UpdatePayload) error {
	finalEndTime, _ := time.Parse(timeFormat, payload.EndTime)
	startTime, _ := time.Parse(timeFormat, payload.StartTime)
	region := payload.Region
	for startTime.Before(finalEndTime) {
		newEndTime := startTime.AddDate(0, 0, 6)
		if !newEndTime.Before(finalEndTime) {
			newEndTime = finalEndTime 
		}
		minreports, err := s.csc.GetEmissions(ctx, region, startTime.Format(timeFormat), newEndTime.Format(timeFormat))
		var NoDataError carbonara.NoDataError
		if err != nil {
			if !errors.As(err, &NoDataError) {
				return mapAndLogErrorf(ctx, "failed to get Carbon Intensity Reports:%w\n", err)
			}
			continue
		}
		dateConfigs, err := getDates(ctx, minreports)
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
				aggErr := s.aggregateData(ctx, region, dateConfigs[j], reportdurations[j])
				if aggErr != nil {
					return mapAndLogErrorf(ctx,  "failed to get Average Carbon Reports:%w\n", aggErr)
				}
			}
		}
		newEndTime = newEndTime.AddDate(0, 0, 1)
		startTime = newEndTime
	}
	return nil
}


// R&D can use this function to obtain CO2 intensity reports for a specific region
func (ser *pollersrvc) GetEmissionsForRegion(ctx context.Context, input *genpoller.CarbonPayload) ([]*genpoller.CarbonForecast, error) {
	var start = input.Start
	var end = input.End
	var region = input.Region
	var reports []*genpoller.CarbonForecast
	reports, err := ser.csc.GetEmissions(ctx, region, start, end)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "failed to get Carbon Intensity Reports:%w\n", err)
	}
	return reports, err
}

// AggregateData gets aggregate reports for all report dates returned by GetDates and store them in clickhouse
func (ser *pollersrvc) aggregateData(ctx context.Context, region string, dates []*genpoller.Period, duration string) (error) {
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

// GetDates gets all the report dates that are used as input to clickhouse queries and obtain aggregate CO2 intensity data
func getDates(ctx context.Context, minutereports []*genpoller.CarbonForecast) ([][]*genpoller.Period, error) {
	if minutereports == nil {
		return nil, fmt.Errorf("no reports for get dates")
	}
	var initialstart, _ = time.Parse(timeFormat, minutereports[0].Duration.StartTime)
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
	for _, event := range minutereports {
		var time, _ = time.Parse(timeFormat, event.Duration.StartTime)
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





