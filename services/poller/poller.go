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
	now	time.Time
}
// timeFormat is used to parse times in order to store time as ISO8601 format
var timeFormat = "2006-01-02T15:04:05-07:00"
//timeNow is used as the end time for all search queries
var timeNow = time.Now
// regions maintains all the valid regions that Singularity will calculate carbon intensity
var regions [13]string = [13]string{model.Caiso, model.Aeso, model.Bpa, model.Erco, model.Ieso,
model.Isone, model.Miso,
 model.Nyiso, model.Nyiso_nycw,
  model.Nyiso_nyli, model.Nyiso_nyup,
   model.Pjm, model.Spp} 
// reportdurations maintains the interval length of each report using constants from the model directory
var reportdurations [5]string = [5]string{model.Minute, model.Hourly, model.Daily, model.Weekly, model.Monthly}
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
		now:				time.Time{},
	}
	current := timeNow()
	s = &pollersrvc{
		csc:				csc,
		dbc:				dbc,
		ctx:                ctx,
		cancel: 			cancel,
		now:				current,
	}
	return s
}
//TODO will ensure that only minute reports are picked up
// EnsurePastData will query clickhouse for the most recent report date
func (s *pollersrvc) ensurePastData(ctx context.Context) (startDates []string) {
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
func (s *pollersrvc) Update(ctx context.Context) error {
	times := s.ensurePastData(ctx)
	finalEndTime, err := time.Parse(timeFormat, s.now.Format(timeFormat))
	if err != nil {
		return mapAndLogError(ctx, err)
	}
	for i := 0; i < len(regions); i++ {
		startTime, err := time.Parse(timeFormat, times[i])
		if err != nil {
			return mapAndLogError(ctx, err)
		}
		region := regions[i]
		
		log.Info(ctx, log.KV{K: "region", V: regions[i]}, 
		log.KV{K: "startTime", V: startTime},
		log.KV{K: "endTime", V: finalEndTime})

		for startTime.Before(finalEndTime) {
			newEndTime := startTime.AddDate(0, 0, 7)
			if !newEndTime.Before(finalEndTime) {
				newEndTime = finalEndTime 
			}
			minreports, err := s.csc.GetEmissions(ctx, regions[i], startTime.Format(timeFormat), newEndTime.Format(timeFormat))
			var NoDataError carbonara.NoDataError
			if err != nil {
				if !errors.As(err, &NoDataError) {
					return mapAndLogErrorf(ctx, "failed to get Carbon Intensity Reports:%w\n", err)
				}
				newEndTime = newEndTime.AddDate(0, 0, 1)
				startTime = newEndTime
				continue
			}	
			log.Info(ctx, log.KV{K: "reports length", V: len(minreports)}, 
			log.KV{K: "startTime", V: startTime},
			log.KV{K: "endTime", V: newEndTime},
			log.KV{K: "report type", V: model.Minute})
			dateConfigs, err := getDates(ctx, minreports)
			
			if err != nil {
				log.Error(ctx, err)
				newEndTime = newEndTime.AddDate(0, 0, 1)
				startTime = newEndTime
				continue
			}
			log.Info(ctx, log.KV{K: "length of hourly dates", V: len(dateConfigs[0])}, 
					log.KV{K: "length of daily dates", V: len(dateConfigs[1])},
					log.KV{K: "length of weekly dates", V: len(dateConfigs[2])},
					log.KV{K: "length of monthly dates", V: len(dateConfigs[3])})
			err = s.dbc.SaveCarbonReports(ctx, minreports)
			if err != nil {
				return mapAndLogErrorf(ctx, "failed to Save Carbon Reports:%w\n", err)
			}
			for j := 0; j < len(dateConfigs); j++ {
				if dateConfigs[j] != nil {
					res, aggErr := s.aggregateData(ctx, region, dateConfigs[j], reportdurations[(j+1)])
					log.Info(ctx, log.KV{K: "reports length", V: len(res)}, 
					log.KV{K: "startTime", V: startTime},
					log.KV{K: "endTime", V: newEndTime},
					log.KV{K: "report type", V: reportdurations[(j + 1)]})
					if aggErr != nil {
						return mapAndLogErrorf(ctx,  "failed to get Average Carbon Reports:%w\n", aggErr)
					}
					if res == nil {
						log.Error(ctx, fmt.Errorf("No aggregate reports returned for region %s and interval type %s\n", regions[i], reportdurations[(j + 1)]))
					}
				}
			}
			startTime = newEndTime
		}
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

// aggregateData gets aggregate reports for all report dates returned by GetDates and store them in clickhouse
func (ser *pollersrvc) aggregateData(ctx context.Context, region string, dates []*genpoller.Period, duration string) ([]*genpoller.CarbonForecast, error) {
	aggregateres, getErr := ser.dbc.GetAggregateReports(ctx, dates, region, duration)
	if getErr != nil {
		return nil, getErr
	}
	saveErr := ser.dbc.SaveCarbonReports(ctx, aggregateres)
	if saveErr != nil {
		return nil, saveErr
	}
	return aggregateres, nil
}

// getDates gets all the report dates that are used as input to clickhouse queries and obtain aggregate CO2 intensity data
func getDates(ctx context.Context, minutereports []*genpoller.CarbonForecast) ([][]*genpoller.Period, error) {
	if minutereports == nil {
		return nil, fmt.Errorf("no reports for get dates")
	}
	var initialstart, _ = time.Parse(timeFormat, minutereports[0].Duration.StartTime)

	var finalDates [][]*genpoller.Period
	var hourlyDates []*genpoller.Period = nil
	var dailyDates []*genpoller.Period = nil
	var weeklyDates []*genpoller.Period = nil
	var monthlyDates []*genpoller.Period = nil
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
	var addedMonthlyReport = false
	for i := 0; i < len(minutereports); i++ {
		var startTime, _ = time.Parse(timeFormat, minutereports[i].Duration.StartTime)
		var endTime, _ = time.Parse(timeFormat, minutereports[i].Duration.EndTime) //make this the new previous/end time
		if endTime.Before(startTime) {
			log.Error(ctx, fmt.Errorf("invalid date"))
			continue
		}
		if startTime.Before(previous) {
			log.Error(ctx, fmt.Errorf("invalid date"))
			continue
		}
		var month = endTime.Month()
		var day = endTime.Day()
		var hour = endTime.Hour()
		var year = endTime.Year()

		if int(month) < int(previous.Month()) && year == previous.Year() {
			log.Error(ctx, fmt.Errorf("invalid date"))
			continue
		}
		if day < previous.Day() && int(month) == int(previous.Month())  {
			log.Error(ctx, fmt.Errorf("invalid date"))
			continue
		}
		if hour < previous.Hour() && day == previous.Day() {
			log.Error(ctx, fmt.Errorf("invalid date"))
			continue
		}

		if hour != previous.Hour() {
			if previous.Equal(initialstart) {
				previous = endTime
			}
			if startTime.AddDate(0, 0, 1).Month() != startTime.Month() && !addedMonthlyReport {
				addedMonthlyReport = true
				monthlyDates = append(monthlyDates, &genpoller.Period{monthstart.Format(timeFormat), previous.Format(timeFormat)})
				monthstart = startTime
			} 
			if day == 1 && (startTime.AddDate(0, 0, 1).Month() == month) {
				addedMonthlyReport = false
			}
			if day != previous.Day() {
				weekcounter += int(endTime.Sub(daystart).Hours() / 24)
				dailyDates = append(dailyDates, &genpoller.Period{daystart.Format(timeFormat), previous.Format(timeFormat)})
				daystart = endTime
				if weekcounter == 7 {
					weeklyDates = append(weeklyDates, &genpoller.Period{weekstart.Format(timeFormat), previous.Format(timeFormat)})
					weekstart = endTime
					weekcounter = 0
				}
			}
			
			hourlyDates = append(hourlyDates, &genpoller.Period{hourstart.Format(timeFormat), previous.Format(timeFormat)})
			hourstart = endTime
		}
		previous = endTime
	}
	finalDates = append(finalDates, hourlyDates)
	finalDates = append(finalDates, dailyDates)
	finalDates = append(finalDates, weeklyDates)
	finalDates = append(finalDates, monthlyDates)
	return finalDates, nil
}





