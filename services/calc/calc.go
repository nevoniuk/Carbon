package calcapi
//something
import (
	"context"
	"fmt"
	//"fmt"
	//"sync"
	"time"
	"github.com/google/uuid"
	"github.com/satori/go.uuid"
	"github.com/crossnokaye/facilityconfig"
	"github.com/crossnokaye/carbon/services/calc/clients/power"
	"github.com/crossnokaye/carbon/services/calc/clients/power_server"
	"github.com/crossnokaye/carbon/services/calc/clients/storage"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
)
//only works for oxnard and riverside
type calcSvc struct {
	psc power.Client
	dbc storage.Client
	psr power_server.Repository
	ctx context.Context
	cancel context.CancelFunc
	
}
var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"


var reportdurations [6]string
func NewCalc(ctx context.Context, psc power.Client, dbc storage.Client, psr power_server.Repository) *calcSvc {
	ctx, cancel := context.WithCancel(ctx)
	s := &calcSvc{
		psc:				psc,
		dbc:				dbc,
		psr:                psr,
		ctx:                ctx,
		cancel: 			cancel,
	}
	reportdurations = [...]string{ "minute", "hourly", "daily", "weekly", "monthly", "yearly"}
	
	
	return s
}

//calculates a report given a carbon report and electrical report
func CalculateReports(context.Context, *gencalc.CarbonReport, *gencalc.ElectricalReport) (*gencalc.TotalReport, error) {
	//input =carbon reports, 
	//1 MWh = 1000 KWh
	//1.convert from MWh to KWh
	//2.
}

//uses store to get input for past-values service
func (s *calcSvc) GetControlPoints(context.Context, *gencalc.PastValuesPayload) ([]string, error) {
	
}

//wrapper function for talking to power client
//power meters at riverside, oxnard
//0's more than a minute resemble blackout


func (s *calcSvc) GetPower(context.Context, *gencalc.GetPowerPayload) (*gencalc.ElectricalReport, error) {
	//power client returns minute reports
	//1.store minute reports in clickhouse
	//2.get dates for averages
	//3. 
}

//wrapper function for talking to storage client
func (s *calcSvc) GetEmissions(context.Context, *gencalc.RequestPayload) ([]*gencalc.CarbonReport, error) {

}

//should maybe also take facility as an input
func (s *calcSvc) HandleRequests(ctx context.Context, req *gencalc.RequestPayload) (error) {
//based on the time period and time interval type that a client wants, get thre respective data from clickhouse
//need region
var payload = &gencalc.EmissionsPayload{Period: req.Period, Interval: &req.Interval}
carbonReports, err := s.GetEmissions()
}

//R&D method
func (s *calcSvc) Carbonreport(context.Context) (err error) {
	//gets reports in carbon forecasts

}

func getdates(ctx context.Context, minutereports []*gencalc.PowerStamp) ([][]*gencalc.Period, error) {
	
	if minutereports == nil {
		var err = fmt.Errorf("no reports for get dates")
		return nil, err
	}

	var initialstart, err = time.Parse(timeFormat, minutereports[0].Period.StartTime)

	if err != nil {
		return nil, err
	}
	var finalDates [][]*gencalc.Period

	var hourlyDates []*gencalc.Period
	var dailyDates []*gencalc.Period
	var weeklyDates []*gencalc.Period
	var monthlyDates []*gencalc.Period
	

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
		
		var time, err = time.Parse(timeFormat, event.Period.StartTime)
		if err != nil {
			return nil, fmt.Errorf("parsing error")
		}
	
		
		var month = time.Month()
		var day = time.Day()
		var hour = time.Hour()
		

		if hour != hourcounter {
			
			if month != monthcounter {
				
				monthcounter = month
				monthlyDates = append(monthlyDates, &gencalc.Period{monthstart.Format(timeFormat), previous.Format(timeFormat)})
				monthstart = time
			}

			if day != daycounter {
				daycounter = day
				weekcounter += 1
				dailyDates = append(dailyDates, &gencalc.Period{daystart.Format(timeFormat), previous.Format(timeFormat)})
				
				daystart = time
				
				if weekcounter == 7 {
					weeklyDates = append(weeklyDates, &gencalc.Period{weekstart.Format(timeFormat), previous.Format(timeFormat)})
					weekstart = time
					weekcounter = 0
				}
			}

			hourcounter = hour
			hourlyDates = append(hourlyDates, &gencalc.Period{hourstart.Format(timeFormat), previous.Format(timeFormat)})
			hourstart = time
		}

		previous = time

	}

	
	if daycounter == time.Time.Day(initialstart) {
		
		dailyDates = append(dailyDates, &gencalc.Period{daystart.Format(timeFormat),previous.Format(timeFormat)})
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

