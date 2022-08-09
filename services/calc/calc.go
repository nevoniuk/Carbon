package calcapi
import (
	"context"
	"fmt"
	"math"
	"time"
	"goa.design/clue/log"
	"github.com/crossnokaye/carbon/model"
	facilityconfig "github.com/crossnokaye/carbon/services/calc/clients/facilityconfig"
	power "github.com/crossnokaye/carbon/services/calc/clients/power"
	storage "github.com/crossnokaye/carbon/services/calc/clients/storage"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
)
type (
	calcSvc struct {
		psc power.Client
		dbc storage.Client
		fc facilityconfig.Client
		ctx context.Context
		cancel context.CancelFunc
	}
)
// timeFormat is used to parse times in order to store time as ISO8601 format
const timeFormat = "2006-01-02T15:04:05-07:00"
// constants below are used to log errors that aren't goa service errors
const (
	FailedToGetCarbonReports = "failed to get carbon reports from clickhouse"
	FailedToGetLocationData = "failed to get location data"
	FailedToGetPowerReports = "failed to get power data"
	FailedToConfigureDates = "failed to configure dates from input"
	FailedToCalculateEmissionReports = "failed to calculate co2 emission reports"
)

func NewCalc(ctx context.Context, psc power.Client, dbc storage.Client, fc facilityconfig.Client) *calcSvc {
	ctx, cancel := context.WithCancel(ctx)
	s := &calcSvc{
		psc:				psc,
		dbc:				dbc,
		fc:                 fc,
		ctx:                ctx,
		cancel: 			cancel,
	}
	return s
}
//note: need to keep UUID's as such in design because this maintains their format
// HistoricalCarbonEmissions will output the CO2 intensity, Power Meter, and resulting CO2 emission reports
func (s *calcSvc) HistoricalCarbonEmissions(ctx context.Context, req *gencalc.RequestPayload) (*gencalc.AllReports, error) {
	log.Info(ctx, log.KV{K: "orgID", V: req.OrgID}, 
		log.KV{K: "facilityID", V: req.FacilityID},
		log.KV{K: "locationID", V: req.LocationID},
		log.KV{K: "start", V: req.Duration.StartTime},
		log.KV{K: "end", V: req.Duration.EndTime},
		log.KV{K: "type", V: req.Interval},
	)
	dates, err := s.getDates(ctx, req.Interval, req.Duration)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToConfigureDates, err)
	}
	log.Info(ctx, log.KV{K: "length of dates", V: len(dates)})
	var intervalduration time.Duration
	switch req.Interval{
	case model.Minute:
		intervalduration = time.Minute
	case model.Hourly:
		intervalduration = time.Hour
	case model.Daily:
		intervalduration = time.Hour * 24
	case model.Weekly:
		intervalduration = time.Hour * 24 * 7
	case model.Monthly:
		intervalduration = time.Hour * 24 * 29
	}
	carbonData, err := s.fc.GetCarbonConfig(ctx, string(req.OrgID), string(req.FacilityID), string(req.LocationID))
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetLocationData, err)
	}
	log.Info(ctx, log.KV{K: "control point name from facility config", V: carbonData.ControlPointName})
	log.Info(ctx, log.KV{K: "agent name from facility config", V: carbonData.AgentName})
	log.Info(ctx, log.KV{K: "formula from facility config", V: carbonData.Formula})
	singularityRegion, controlPointName, formula, agentName := carbonData.Region, carbonData.ControlPointName, carbonData.Formula, carbonData.AgentName
	carbonReport, err := s.dbc.GetCarbonIntensityReports(ctx, dates, req.Interval, singularityRegion)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetCarbonReports, err)
	}
	endTime, _ := time.Parse(timeFormat, req.Duration.EndTime)
	startTime, _ := time.Parse(timeFormat, req.Duration.StartTime)
	difference := endTime.Sub(startTime)
	duration := difference.Nanoseconds()
	powerReport, err := s.psc.GetPower(ctx, string(req.OrgID), req.Duration, controlPointName, duration, intervalduration, &formula, agentName)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetPowerReports, err)
	}
	emissionReport, err := calculateCarbonEmissionsReport(ctx, carbonReport, powerReport, intervalduration)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToCalculateEmissionReports, err)
	}
	return &gencalc.AllReports{TotalEmissionReport: emissionReport}, nil
}

// calculateCarbonEmissionsReport yields lbs of CO2 given CO2 intensity(CO2lbs/MWh) and Power(KWh) reports
func calculateCarbonEmissionsReport(ctx context.Context, carbonReport *gencalc.CarbonReport, powerReport *gencalc.ElectricalReport, intervalType time.Duration) (*gencalc.EmissionsReport, error) {
	var dataPoints []*gencalc.DataPoint
	var powerreportCounter = 0
	var intenreportCounter = 0
	for intenreportCounter < len(carbonReport.IntensityPoints) && powerreportCounter < len(powerReport.PowerStamps) {
		log.Info(ctx, log.KV{K: "Intensity Point", V: carbonReport.IntensityPoints[intenreportCounter]})
		log.Info(ctx, log.KV{K: "Power Point", V: powerReport.PowerStamps[powerreportCounter]})
		powerT,_ := time.Parse(timeFormat, powerReport.PowerStamps[powerreportCounter].Time)
		carbonT,_ := time.Parse(timeFormat, carbonReport.IntensityPoints[intenreportCounter].Time)
		var difference time.Duration
		var leastTime = true
		if powerT.Before(carbonT) {
			difference = carbonT.Sub(powerT)
			if difference > intervalType {
				powerreportCounter += 1
				continue
			}
		} else {
			leastTime = false
			difference = powerT.Sub(carbonT)
			if difference > intervalType {
				intenreportCounter += 1
				continue
			}
		}
		toKWh := carbonReport.IntensityPoints[intenreportCounter].Value / 1000 //convert mwh->kwh
		carbonemissions := toKWh * powerReport.PowerStamps[powerreportCounter].Value
		var time = carbonReport.IntensityPoints[intenreportCounter].Time
		if leastTime {
			time = powerReport.PowerStamps[powerreportCounter].Time
		} 
		log.Info(ctx, log.KV{K: "Emissions Point", V: &gencalc.DataPoint{Time: time, Value: carbonemissions}})
		fmt.Println(&gencalc.DataPoint{Time: time, Value: carbonemissions})
		dataPoints = append(dataPoints, &gencalc.DataPoint{Time: time, Value: carbonemissions})
		intenreportCounter += 1
		powerreportCounter += 1
	}
	return &gencalc.EmissionsReport{Duration: powerReport.Duration, Interval: powerReport.Interval, Points: dataPoints}, nil
}

// getDates returns an array of dates for the storage client in order to correctly query for carbon reports
func (s *calcSvc) getDates(ctx context.Context, intervalType string, duration *gencalc.Period) ([]*gencalc.Period, error) {
	var newDates []*gencalc.Period
	initialstart, err := time.Parse(timeFormat, duration.StartTime)
	if err != nil {
		return nil, gencalc.MakeReportsNotFound(fmt.Errorf("incorrect start date given: %w", err))
	}
  	end, err:= time.Parse(timeFormat, duration.EndTime)
	  if err != nil {
		return nil, gencalc.MakeReportsNotFound(fmt.Errorf("incorrect end date given: %w", err))
	}
	var diff = end.Sub(initialstart)
	var datesCount int
	var durationType int
	switch intervalType {
	case model.Minute: 
		datesCount = int(math.Ceil(diff.Minutes()))
		durationType = int(time.Minute)
	case model.Hourly:
		datesCount = int(math.Ceil(diff.Hours()))
		durationType = int(time.Hour)
	case model.Daily: 
		datesCount = int(math.Ceil(diff.Hours())) / 24
		durationType = int(time.Hour) * 24
	case model.Weekly:
		datesCount = int(math.Ceil(diff.Hours())) / (24 * 7)
		durationType = int(time.Hour) * 24 * 7
	case model.Monthly: 
		datesCount = int(math.Ceil(diff.Hours())) / (24 * 29)
		durationType = int(time.Hour) * 24 * 29
	default:
		datesCount = 1
	}
	var tempstart = initialstart
	var tempend time.Time
	for i := 0; i < datesCount; i++ {
		tempend = tempstart.Add(time.Duration(durationType))
		if tempend.After(end) { 
			break
		}
		var startString = tempstart.Format(timeFormat)
		var endString = tempend.Format(timeFormat)
		newDates = append(newDates, &gencalc.Period{StartTime: startString, EndTime: endString})
		tempstart = tempend
	}
	if len(newDates) == 0 {
		return nil, gencalc.MakeReportsNotFound(fmt.Errorf("no dates for period %s and interval %s: %w", duration, intervalType, err))
	}
	return newDates, nil
}



