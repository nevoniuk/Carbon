package calcapi

import (
	"context"
	"fmt"
	"math"
	"time"
	"errors"
	"goa.design/clue/log"
	goa "goa.design/goa/v3/pkg"
	//errors as
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
const (timeFormat = "2006-01-02T15:04:05-07:00")
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

//HandleRequests will output the CO2 intensity, Power Meter, and resulting CO2 emission reports
func (s *calcSvc) HistoricalCarbonEmissions(ctx context.Context, req *gencalc.RequestPayload) (*gencalc.AllReports, error) {
	if req.Interval == "" {
		return nil, fmt.Errorf("requested interval is null\n",)
	}

	dates, err := s.getDates(ctx, req.Interval, req.Duration)
	if err != nil {
		log.Errorf(ctx, err, "error parsing time in getDates:%s\n", err)
		return nil, err
	}

	var res *facilityconfig.Carbon
	var facilitynotFound facilityconfig.ErrFacilityNotFound
	var locationnotFound facilityconfig.ErrLocationNotFound

	res, err = s.getLocationData(ctx, string(req.OrgID), string(req.FacilityID), string(req.LocationID))
	if !errors.As(err, &facilitynotFound) && !errors.As(err, &locationnotFound)  {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetLocationData, err)
	} else if err != nil {
		mapAndLogError(ctx, err)
	}

	singularityRegion, controlPointName, formula := res.Region, res.ControlPointName, res.Formula

	var carbonnotFound storage.ErrNotFound
	carbonReports, err := s.getEmissions(ctx, dates, req.Interval, singularityRegion)
	if !errors.As(err, &carbonnotFound) {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetCarbonReports, err)
	} else if err != nil {
		mapAndLogError(ctx, err)
	}
		
	endTime, err := time.Parse(timeFormat, req.Duration.EndTime)
	if err != nil {
		return nil, fmt.Errorf("parsing time err: %s\n", err)
	}
	startTime, err:= time.Parse(timeFormat, req.Duration.StartTime)
	if err != nil {
		return nil, fmt.Errorf("parsing time err: %s\n", err)
	}
	difference := endTime.Sub(startTime)
	duration := difference.Nanoseconds()
	
	powerReports, err := s.getPower(ctx, string(req.OrgID), req.Duration, controlPointName, duration, req.Interval, &formula)
	var powernotFound power.ErrPowerReportsNotFound
	if !errors.As(err, &powernotFound) {
		return nil, mapAndLogErrorf(ctx, "%s: %w", FailedToGetPowerReports, err)
	} else if err != nil {
		mapAndLogError(ctx, err)
	}
	
	emissionReport, err := calculateCarbonEmissionsReport(ctx, carbonReports, powerReports)
	if err != nil {
		return nil, fmt.Errorf("error from CalculateReports: %s\n", err)
	}
	return &gencalc.AllReports{CarbonIntensityReports: carbonReports, PowerReports: powerReports,TotalEmissionReport: emissionReport}, nil
}

//CalculateReports yields lbs of CO2 given CO2 intensity(CO2lbs/MWh) and Power(KWh) reports
//Assumes that each report has the same duration and interval type
func calculateCarbonEmissionsReport(ctx context.Context, carbonReports []*gencalc.CarbonReport, powerReports []*gencalc.ElectricalReport) (*gencalc.EmissionsReport, error) {
	var dataPoints []*gencalc.DataPoint
	for i, r := range carbonReports {
		toKWh := r.GeneratedRate * 1000 //convert mwh->kwh
		carbonemissions := toKWh * powerReports[i].Power
		dataPoints = append(dataPoints, &gencalc.DataPoint{Time: r.Duration.StartTime, CarbonFootprint:carbonemissions})
	}
	emissionsreport := &gencalc.EmissionsReport{Duration: carbonReports[0].Duration, Interval: carbonReports[0].Interval, Points: dataPoints}
	return emissionsreport, nil
}

//getLocationData uses a facility config client to get the following input for past-values service: control point name for power meter, formula for power conversion, and region for carbonemissions
func (s *calcSvc) getLocationData(ctx context.Context, orgID string, facilityID string, locationID string) (*facilityconfig.Carbon, error) {
	carbonData, err := s.fc.GetCarbonConfig(ctx, orgID, facilityID, locationID)
	if err != nil {
		return nil, err
	}
	return carbonData, nil
}

//GetPower is a wrapper function for talking to the power client. Right now there is only a power meter
//at Oxnard so this will only work for that power meter
func (s *calcSvc) getPower(ctx context.Context, orgID string, dateRange *gencalc.Period, cpaliasname string, pastValInterval int64, reportInterval string, formula *string) ([]*gencalc.ElectricalReport, error) {
	var reports []*gencalc.ElectricalReport
	if cpaliasname == "" {
		return nil, fmt.Errorf("no Control Point\n")
	}
	if formula == nil {
		return nil, fmt.Errorf("no Formula\n")
	}
	
	reports, err := s.psc.GetPower(ctx, orgID, cpaliasname, pastValInterval, dateRange.StartTime, dateRange.EndTime, formula)
	if err != nil {
		return nil, fmt.Errorf("error from GetPower: %s\n", err)
	}

	return reports, nil
}

//GetEmissions is a wrapper function for talking to storage client
func (s *calcSvc) getEmissions(ctx context.Context, dates []*gencalc.Period, interval string, region string) ([]*gencalc.CarbonReport, error) {
	reports, err := s.dbc.GetCarbonReports(ctx, dates, interval, region)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

//GetDates returns an array of dates for the storage client in order to correctly query for carbon reports
func (s *calcSvc) getDates(ctx context.Context, intervalType string, duration *gencalc.Period) ([]*gencalc.Period, error) {
	var newDates []*gencalc.Period
	initialstart, err1 := time.Parse(timeFormat, duration.StartTime)
	if err1 != nil {
		log.Error(ctx, err1, log.KV{K: "msg", V: "failed to parse time"})
		return nil, err1
	}
	end, err2 := time.Parse(timeFormat, duration.EndTime)
	if err2 != nil {
		log.Error(ctx, err2, log.KV{K: "msg", V: "failed to parse time"})
		return nil, err2
	}

	var diff = initialstart.Sub(end)
	var datesCount int
	var durationType int
	
	switch intervalType{
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


func mapAndLogErrorf(ctx context.Context, format string, a ...interface{}) error {
	return mapAndLogError(ctx, fmt.Errorf(format, a...))
}
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



