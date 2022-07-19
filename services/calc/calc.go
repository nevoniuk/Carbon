package calcapi

import (
	"context"
	"fmt"
	"time"
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
var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"
var reportdurations [5]string = [5]string{ "minute", "hourly", "daily", "weekly", "monthly"}

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
	var emissionReport *gencalc.EmissionsReport
	var carbonReports []*gencalc.CarbonReport
	var powerReports []*gencalc.ElectricalReport
	
	var dates []*gencalc.Period
	dates, err := s.getDates(ctx, *req.Interval.Kind, req.Duration)
	if err != nil {
		return nil, fmt.Errorf("Error parsing time in GetDates: %s\n", err)
	}
		
	var res *facilityconfig.Carbon
	res, err = s.getLocationData(ctx, string(req.OrgID), string(req.FacilityID), string(*req.LocationID))
	if err != nil {
		return nil, fmt.Errorf("Error from GetPowerControlPoint: %s\n", err)
	}
	var singularityRegion = res.Region
	var controlPointName = res.ControlPointName
	var formula = res.Formula
		
	carbonReports, err = s.getEmissions(ctx, dates, *req.Interval.Kind, singularityRegion)
	if err != nil {
		return nil, fmt.Errorf("Error from GetEmissions: %s\n", err)
	}
		
	endTime, timeError1 := time.Parse(timeFormat, req.Duration.EndTime)
	if timeError1 != nil {
		return nil, fmt.Errorf("parsing time err: %s\n", timeError1)
	}
	startTime, timeError2 := time.Parse(timeFormat, req.Duration.StartTime)
	if timeError2 != nil {
		return nil, fmt.Errorf("parsing time err: %s\n", timeError2)
	}
	difference := endTime.Sub(startTime)
	duration := difference.Nanoseconds()
	
	powerReports, err = s.getPower(ctx, string(req.OrgID), req.Duration, controlPointName, duration, *req.Interval.Kind, &formula)
	if err != nil {
		return nil, fmt.Errorf("Error in GetPower: %s\n", err)
	}

	emissionReport, err = calculateCarbonEmissionsReport(ctx, carbonReports, powerReports)
	if err != nil {
		return nil, fmt.Errorf("Error from Calculate Reports: %s\n", err)
	}
	return &gencalc.AllReports{CarbonIntensityReports: carbonReports, PowerReports: powerReports,TotalEmissionReport: emissionReport}, nil
}

//CalculateReports yields lbs of CO2 given CO2 intensity(CO2lbs/MWh) and Power(KWh) reports
//Assumes that each report has the same duration and interval type
func calculateCarbonEmissionsReport(ctx context.Context, carbonReports []*gencalc.CarbonReport, powerReports []*gencalc.ElectricalReport) (*gencalc.EmissionsReport, error) {
	var dataPoints []*gencalc.DataPoint
	var report *gencalc.EmissionsReport
	for _, report := range carbonReports {
		toKWh := report.GeneratedRate * 1000
		dataPoints = append(dataPoints, &gencalc.DataPoint{Time: report.Duration.StartTime, CarbonFootprint:toKWh})
	}
	report = &gencalc.EmissionsReport{Duration: report.Duration, Interval: report.Interval, Points: dataPoints}
	return report, nil
}

//getLocationData uses a facility config client to get the following input for past-values service: control point name for power meter, formula for power conversion, and region for carbonemissions
func (s *calcSvc) getLocationData(ctx context.Context, orgID string, facilityID string, locationID string) (*facilityconfig.Carbon, error) {
	if orgID == "" {
		return nil, fmt.Errorf("Org ID is null\n")
	}
	if facilityID == "" {
		return nil, fmt.Errorf("Facility ID is null\n")
	}

	var CarbonData *facilityconfig.Carbon
	CarbonData, err := s.fc.GetCarbonConfig(ctx, orgID, facilityID, locationID)
	if err != nil {
		return nil, err
	}

	return CarbonData, nil
}

//GetPower is a wrapper function for talking to the power client. Right now there is only a power meter
//at Oxnard so this will only work for that power meter
func (s *calcSvc) getPower(ctx context.Context, orgID string, dateRange *gencalc.Period, cpaliasname string, pastValInterval int64, reportInterval string, formula *string) ([]*gencalc.ElectricalReport, error) {
	var reports []*gencalc.ElectricalReport
	
	if orgID == "" {
		return nil, fmt.Errorf("Org ID is null\n")
	}
	if cpaliasname == "" {
		return nil, fmt.Errorf("No Control Point\n")
	}
	if formula == nil {
		return nil, fmt.Errorf("No Formula\n")
	}
	
	reports, err := s.psc.GetPower(ctx, orgID, cpaliasname, pastValInterval, dateRange.StartTime, dateRange.EndTime, formula)
	if err != nil {
		return nil, fmt.Errorf("Error from GetPower: %s\n", err)
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

	//var counter int
	var newDates []*gencalc.Period
	initialstart, err1 := time.Parse(timeFormat, duration.StartTime)
	end, err2 := time.Parse(timeFormat, duration.EndTime)
	
	if err1 != nil {
		return nil, err1
	}
	if err2 != nil {
		return nil, err2
	}

	var diff = initialstart.Sub(end)
	var datesCount float64
	var durationType int

	switch intervalType {
	case reportdurations[0]: //minute
		//counter = time.Time.Minute(initialstart)
		datesCount = diff.Minutes()
		durationType = int(time.Minute)
	case reportdurations[1]: //hour
		//counter = time.Time.Hour(initialstart)
		datesCount = diff.Hours()
		durationType = int(time.Hour)
	case reportdurations[2]: //daily
		//counter = time.Time.Day(initialstart)
		datesCount = diff.Hours() / 24
		durationType = int(time.Hour) * 24
	case reportdurations[3]: //weekly
		//counter = time.Time.Day(initialstart)
		datesCount = diff.Hours() / (24 * 7)
		durationType = int(time.Hour) * 24 * 7
	case reportdurations[4]: //monthly
		//var m time.Month = time.Time.Month(initialstart)
		//counter = int(m)
		datesCount = diff.Hours() / (24 * 29)
		durationType = int(time.Hour) * 24 * 29
	}

	var tempstart = initialstart
	var tempend time.Time

	for i := 0.0; i < datesCount; i++ {
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



