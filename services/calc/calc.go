package calcapi
import (
	"context"
	"fmt"
	"math"
	"time"
	"errors"
	"goa.design/clue/log"
	goa "goa.design/goa/v3/pkg"
	"github.com/crossnokaye/carbon/model"
	"github.com/crossnokaye/carbon/services/calc/clients/facilityconfig"
	"github.com/crossnokaye/carbon/services/calc/clients/power"
	"github.com/crossnokaye/carbon/services/calc/clients/storage"
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
	uuidArray []struct{}
)
// timeFormat is used to parse times in order to store time as ISO8601 format
const timeFormat = "2006-01-02T15:04:05-07:00"
// constants below are used to log errors that aren't goa service errors
const (
	FailedToGetCarbonReports = "failed to get carbon reports from clickhouse"
	FailedToGetLocationData = "failed to get location data"
	FailedToGetPowerReports = "failed to get power data"
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

// HistoricalCarbonEmissions will output the CO2 intensity, Power Meter, and resulting CO2 emission reports
func (s *calcSvc) HistoricalCarbonEmissions(ctx context.Context, req *gencalc.RequestPayload) (*gencalc.AllReports, error) {
	dates, err := s.getDates(ctx, req.Interval, req.Duration)
	if err != nil {
		log.Errorf(ctx, err, "error parsing time in getDates:%w\n", err)
		return nil, err
	}

	var res *facilityconfig.Carbon
	res, err = s.getLocationData(ctx, string(req.OrgID), string(req.FacilityID), string(req.LocationID))
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetLocationData, err)
	}
	singularityRegion, controlPointName, formula, agentName := res.Region, res.ControlPointName, res.Formula, res.AgentName

	carbonReports, err := s.getCarbonIntensityData(ctx, dates, req.Interval, singularityRegion)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetCarbonReports, err)
	}
		
	endTime, _ := time.Parse(timeFormat, req.Duration.EndTime)
	startTime, _ := time.Parse(timeFormat, req.Duration.StartTime)
	difference := endTime.Sub(startTime)
	duration := difference.Nanoseconds()
	
	powerReports, err := s.getPower(ctx, string(req.OrgID), req.Duration, controlPointName, duration, req.Interval, &formula, agentName)
	if err != nil {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetPowerReports, err)
	}
	
	emissionReport, err := calculateCarbonEmissionsReport(ctx, carbonReports, powerReports)
	if err != nil {
		log.Errorf(ctx, err, "error calcualting carbon emission reports:%w\n", err)
		return nil, err
	}
	return &gencalc.AllReports{CarbonIntensityReports: carbonReports, PowerReports: powerReports,TotalEmissionReport: emissionReport}, nil
}

// calculateCarbonEmissionsReport yields lbs of CO2 given CO2 intensity(CO2lbs/MWh) and Power(KWh) reports
func calculateCarbonEmissionsReport(ctx context.Context, carbonReports []*gencalc.CarbonReport, powerReports []*gencalc.ElectricalReport) (*gencalc.EmissionsReport, error) {
	var dataPoints []*gencalc.DataPoint
	for i, r := range carbonReports {
		toKWh := r.GeneratedRate * 1000 //convert mwh->kwh
		carbonemissions := toKWh * powerReports[i].Power
		dataPoints = append(dataPoints, &gencalc.DataPoint{Time: r.Duration.StartTime, CarbonFootprint:carbonemissions})
	}
	return &gencalc.EmissionsReport{Duration: carbonReports[0].Duration, Interval: carbonReports[0].Interval, Points: dataPoints}, nil
}

// getLocationData uses a facility config client to get the following input for past-values service: control point name for power meter, formula for power conversion, and region for carbonemissions
func (s *calcSvc) getLocationData(ctx context.Context, orgID string, facilityID string, locationID string) (*facilityconfig.Carbon, error) {
	carbonData, err := s.fc.GetCarbonConfig(ctx, orgID, facilityID, locationID)
	if err != nil {
		return nil, err
	}
	return carbonData, nil
}

/**
getPower is a wrapper function for talking to the power client. Right now there is only a power meter
at Oxnard so this will only work for that power meter
*/
func (s *calcSvc) getPower(ctx context.Context, orgID string, dateRange *gencalc.Period, cpaliasname string, pastValInterval int64, reportInterval string, formula *string, agentname string) ([]*gencalc.ElectricalReport, error) {
	var reports []*gencalc.ElectricalReport
	payloadDuration := &gencalc.Period{StartTime: dateRange.StartTime, EndTime: dateRange.EndTime}
	payload := &gencalc.PastValPayload{OrgID: orgID, Duration: payloadDuration, PastValInterval: pastValInterval, ControlPoint: cpaliasname, Formula: formula, AgentName: agentname}
	reports, err := s.psc.GetPower(ctx, payload)
	if err != nil {
		return nil, err
	}

	return reports, nil
}

// getCarbonIntensityData is a wrapper function for talking to storage client
func (s *calcSvc) getCarbonIntensityData(ctx context.Context, dates []*gencalc.Period, interval string, region string) ([]*gencalc.CarbonReport, error) {
	reports, err := s.dbc.GetCarbonReports(ctx, dates, interval, region)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

// getDates returns an array of dates for the storage client in order to correctly query for carbon reports
func (s *calcSvc) getDates(ctx context.Context, intervalType string, duration *gencalc.Period) ([]*gencalc.Period, error) {
	var newDates []*gencalc.Period
	initialstart, _ := time.Parse(timeFormat, duration.StartTime)
	end, _ := time.Parse(timeFormat, duration.EndTime)
	var diff = initialstart.Sub(end)
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
	}

	var tempstart = initialstart
	var tempend time.Time
	for i := 0; i < datesCount; i++ {
		tempend = initialstart.Add(time.Duration(durationType))
		if tempend.After(end) {
			break
		}
		var startString = tempstart.Format(timeFormat)
		var endString = tempend.Format(timeFormat)
		newDates = append(newDates, &gencalc.Period{StartTime: startString, EndTime: endString})
		tempstart = tempend
	}
	return newDates, nil
}

// mapAndLogErrorf maps client errors to a coordinator service error responses using
// the format string.
func mapAndLogErrorf(ctx context.Context, format string, a ...interface{}) error {
	return mapAndLogError(ctx, fmt.Errorf(format, a...))
}

// mapAndLogError maps client errors to a coordinator service error responses.
// It logs the error using the context if it does not map to a design error
// (i.e. is unexpected).
func mapAndLogError(ctx context.Context, err error) error {
	var gerr *goa.ServiceError
	if errors.As(err, &gerr) {
		if gerr.Name == "not_found" {
			return gencalc.MakeNotFound(gerr)
		}
	}
	var carbonreportsNotFound storage.ErrNotFound
	if errors.As(err, &carbonreportsNotFound) {
		return gencalc.MakeNotFound(carbonreportsNotFound)
	}
	var fNotFound *facilityconfig.ErrFacilityNotFound
	if errors.As(err, &fNotFound) {
		return gencalc.MakeNotFound(fNotFound)
	}

	var lNotFound *facilityconfig.ErrLocationNotFound
	if errors.As(err, &lNotFound) {
		return gencalc.MakeNotFound(lNotFound)
	}

	var powerreportsNotFound *power.ErrPowerReportsNotFound
	if errors.As(err, &powerreportsNotFound) {
		return gencalc.MakeNotFound(powerreportsNotFound)
	}
	log.Error(ctx, err)
	return err
}



