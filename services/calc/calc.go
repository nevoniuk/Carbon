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
	//remove after PR is merged into legacy
	//remove after testing
	/**
	var res *facilityconfig.Carbon
	carbonData, err := s.fc.GetCarbonConfig(ctx, orgID, facilityID, locationID)
	res, err = s.getLocationData(ctx, string(req.OrgID), string(req.FacilityID), string(req.LocationID))
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetLocationData, err)
	}
	singularityRegion, controlPointName, formula, agentName := res.Region, res.ControlPointName, res.Formula, res.AgentName
*/
	formula := "0.6"
	agentName := "office Lineage Oxnard Building 4"
	singularityRegion := model.Caiso
	controlPointName := "energy_meter_4_pulse_val"
	carbonReport, err := s.dbc.GetCarbonIntensityReports(ctx, dates, req.Interval, singularityRegion)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetCarbonReports, err)
	}
	endTime, _ := time.Parse(timeFormat, req.Duration.EndTime)
	startTime, _ := time.Parse(timeFormat, req.Duration.StartTime)
	difference := endTime.Sub(startTime)
	duration := difference.Nanoseconds()
	powerReport, err := s.psc.GetPower(ctx, string(req.OrgID), req.Duration, controlPointName, duration, req.Interval, &formula, agentName)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetPowerReports, err)
	}
	emissionReport, err := calculateCarbonEmissionsReport(ctx, carbonReport, powerReport)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToCalculateEmissionReports, err)
	}
	return &gencalc.AllReports{TotalEmissionReport: emissionReport}, nil
}

// calculateCarbonEmissionsReport yields lbs of CO2 given CO2 intensity(CO2lbs/MWh) and Power(KWh) reports
func calculateCarbonEmissionsReport(ctx context.Context, carbonReport *gencalc.CarbonReport, powerReport *gencalc.ElectricalReport) (*gencalc.EmissionsReport, error) {
	var dataPoints []*gencalc.DataPoint
	for i, r := range carbonReport.IntensityPoints {
		toKWh := r.Value * 1000 //convert mwh->kwh
		carbonemissions := toKWh * powerReport.PowerStamps[i].Value
		fmt.Println("Data Point")
		fmt.Println(&gencalc.DataPoint{Time: powerReport.Duration.StartTime, Value: carbonemissions})
		dataPoints = append(dataPoints, &gencalc.DataPoint{Time: powerReport.Duration.StartTime, Value: carbonemissions})
	}
	return &gencalc.EmissionsReport{Duration: powerReport.Duration, Interval: powerReport.Interval, Points: dataPoints}, nil
}

// getLocationData uses a facility config client to get the following input for past-values service: control point name for power meter, formula for power conversion, and region for carbonemissions
func (s *calcSvc) getLocationData(ctx context.Context, orgID string, facilityID string, locationID string) (*facilityconfig.Carbon, error) {
	carbonData, err := s.fc.GetCarbonConfig(ctx, orgID, facilityID, locationID)
	if err != nil {
		return nil, err
	}
	return carbonData, nil
}

// getDates returns an array of dates for the storage client in order to correctly query for carbon reports
func (s *calcSvc) getDates(ctx context.Context, intervalType string, duration *gencalc.Period) ([]*gencalc.Period, error) {
	var newDates []*gencalc.Period
	initialstart, err := time.Parse(timeFormat, duration.StartTime)
	if err != nil {
		return nil, gencalc.MakeReportsNotFound(fmt.Errorf("incorrect start date given: %w", err))
	}
	fmt.Println(initialstart)
  	end, err:= time.Parse(timeFormat, duration.EndTime)
	  if err != nil {
		return nil, gencalc.MakeReportsNotFound(fmt.Errorf("incorrect end date given: %w", err))
	}
	fmt.Println(end)
	var diff = end.Sub(initialstart)
	fmt.Println(diff)
	var datesCount int
	var durationType int
	switch intervalType {
	case model.Minute: 
		datesCount = int(math.Ceil(diff.Minutes()))
		durationType = int(time.Minute)
	case model.Hourly:
		fmt.Println("hourly average")
		fmt.Println(diff.Hours())
		fmt.Println(math.Ceil(diff.Hours()))
		datesCount = int(math.Ceil(diff.Hours()))
		fmt.Println(datesCount)
		durationType = int(time.Hour)
		fmt.Println(durationType)
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
		fmt.Println(tempstart)
		fmt.Println(tempend)
		tempend = tempstart.Add(time.Duration(durationType))
		if tempend.After(end) { 
			break
		}
		//implement logic to iterate until the end of the month
		var startString = tempstart.Format(timeFormat)
		var endString = tempend.Format(timeFormat)
		fmt.Println(&gencalc.Period{StartTime: startString, EndTime: endString})
		newDates = append(newDates, &gencalc.Period{StartTime: startString, EndTime: endString})
		tempstart = tempend
	}
	if len(newDates) == 0 {
		return nil, gencalc.MakeReportsNotFound(fmt.Errorf("no dates for period %s and interval %s: %w", duration, intervalType, err))
	}
	return newDates, nil
}



