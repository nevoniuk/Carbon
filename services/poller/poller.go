package pollerapi

import (
	"context"
	"errors"
	"time"
	"fmt"
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
// common start date for regions
const regionstartdate = "2020-01-01T00:00:00+00:00"
const regiontestdate = "2022-01-01T00:00:00+00:00"
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

// EnsurePastData will query clickhouse for the most recent report date
func (s *pollersrvc) ensurePastData(ctx context.Context) (startDates []string) {
	var dates []string
	for i := 0; i < len(regions); i++ {
		date, err := s.dbc.CheckDB(ctx, string(regions[i]))
		if err == nil {
			dates = append(dates, date)
		} else {
			log.Info(ctx, log.KV{K: "error from checkdb %w", V: err})
			var defaultDate string
			if regions[i] == model.Aeso {
				defaultDate = AesoStartDate
			} else {defaultDate = regiontestdate}
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
		fmt.Println("start time")
		fmt.Println(startTime)
		fmt.Println("region")
		fmt.Println(regions[i])
		if err != nil {
			return mapAndLogError(ctx, err)
		}
		var endTime = finalEndTime
		var twoWeeksDuration = time.Hour * 24 * 14
		if startTime.Before(finalEndTime.Add(twoWeeksDuration * -1)) {//if the query is longer than two weeks
			endTime = startTime.Add(twoWeeksDuration + time.Second)
			fmt.Println("query over two weeks: new end time")
			fmt.Println(endTime)
		}

		for startTime.Before(endTime) {
			newEndTime := startTime.AddDate(0, 0, 7)
			if !newEndTime.Before(finalEndTime) {
				newEndTime = finalEndTime 
			}
			fmt.Println("in loop start time")
			fmt.Println(startTime)
			fmt.Println("in loop end time")
			fmt.Println(newEndTime)
			minreports, err := s.csc.GetEmissions(ctx, regions[i], startTime.Format(timeFormat), newEndTime.Format(timeFormat))
			fmt.Println("length of reports")
			fmt.Println(len(minreports))
			var noDataError carbonara.NoDataError
			if err != nil {
				if !errors.As(err, &noDataError) {
					return mapAndLogErrorf(ctx, "failed to get Carbon Intensity Reports:%w\n", err)
				}
				//newEndTime = newEndTime.AddDate(0, 0, 1)
				startTime = newEndTime
				continue
			}	
			err = s.dbc.SaveCarbonReports(ctx, minreports)
			if err != nil {
				return mapAndLogErrorf(ctx, "failed to Save Carbon Reports:%w\n", err)
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


