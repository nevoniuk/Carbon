package calcapi
import (
	"context"
	"fmt"
	"time"
	"github.com/crossnokaye/carbon/services/calc/clients/power"
	"github.com/crossnokaye/carbon/services/calc/clients/power_server"
	"github.com/crossnokaye/carbon/services/calc/clients/storage"
	"github.com/crossnokaye/facilityconfig"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	"github.com/google/uuid"
)

type (
	calcSvc struct {
	psc power.Client
	dbc storage.Client
	psr power_server.Repository
	ctx context.Context
	cancel context.CancelFunc
	}
	uuidArray []struct{}
)
var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"

var reportdurations [5]string = [5]string{ "minute", "hourly", "daily", "weekly", "monthly"}
func NewCalc(ctx context.Context, psc power.Client, dbc storage.Client, psr power_server.Repository) *calcSvc {
	ctx, cancel := context.WithCancel(ctx)
	s := &calcSvc{
		psc:				psc,
		dbc:				dbc,
		psr:                psr,
		ctx:                ctx,
		cancel: 			cancel,
	}
	return s
}

//CalculateReports yields lbs of CO2 given CO2 intensity(CO2lbs/MWh) and Power(KWh) reports
//Assumes that each report has the same duration and interval type
func CalculateEmissionsReport(ctx context.Context, carbonReports []*gencalc.CarbonReport, powerReports []*gencalc.ElectricalReport) (*gencalc.EmissionsReport, error) {
	var dataPoints []*gencalc.DataPoint
	var report *gencalc.EmissionsReport
	var org = powerReports[0].OrgID
	var agent = powerReports[0].AgentID
	for _, report := range carbonReports {
		toKWh := report.GeneratedRate * 1000
		dataPoints = append(dataPoints, &gencalc.DataPoint{Time: report.Duration.StartTime, CarbonFootprint:toKWh})
	}
	report = &gencalc.EmissionsReport{Duration: report.Duration, DurationType: report.DurationType, Points: dataPoints, OrgID: org, AgentID: agent}
	return report, nil
}

//GetPowerControlPoint uses a facility config store to get the following input for past-values service: pointname for power meter, facility data, building data
//It will get those values with the following input from the HandleRequest function: OrgID, AgentID, FacilityID
func (s *calcSvc) GetPowerControlPoint(ctx context.Context, org uuid.UUID, agent string) ([]uuid.UUID, error) {
	var temp []uuid.UUID
	if org == uuid.Nil {
		return temp, fmt.Errorf("Org ID is null\n")
	}
	
	if agent == "" {
		return temp, fmt.Errorf("Agent ID is null\n")
	}
	pointName, err := s.psr.FindControlPointName(org uuid.UUID, agent string, facility uuid.UUID)
	//TODO: find this point name 
	point, err := s.psr.FindControlPointIDsByName(org, agent, pointName)
	if err != nil {
		return temp, fmt.Errorf("Error finding control point: [%s]\n", err)
	}
	return point, nil
}

//GetPower is a wrapper function for talking to the power client. Right now there is only a power meter
//at Oxnard so this will only work for that power meter
func (s *calcSvc) GetPower(ctx context.Context, org uuid.UUID, dateRange *gencalc.Period, cps []uuid.UUID, pastValInterval int64, reportInterval string) ([]*gencalc.ElectricalReport, error) {
	var reports []*gencalc.ElectricalReport
	//nullid := uuid.Nil
	if org == uuid.Nil {
		return nil, fmt.Errorf("Org ID is null\n")
	}
	if cps[0] == uuid.Nil {
		return nil, fmt.Errorf("No Control Points\n")
	}
	//interval has to be in nanoseconds
	return reports, nil
}

//GetEmissions is a wrapper function for talking to storage client
func (s *calcSvc) GetEmissions(ctx context.Context, dates []*gencalc.Period, interval string, region string) ([]*gencalc.CarbonReport, error) {
	reports, err := s.dbc.GetCarbonReports(ctx, dates, interval, region)
	if err != nil {
		return nil, err
	}
	return reports, nil
}

//HandleRequests will output the CO2 intensity, Power Meter, and resulting CO2 emission reports
func (s *calcSvc) HandleRequests(ctx context.Context, req *gencalc.RequestPayload) (*gencalc.AllReports, error) {
	var emissionReport *gencalc.EmissionsReport
	var carbonReports []*gencalc.CarbonReport
	var powerReports []*gencalc.ElectricalReport
	var validInterval bool
	for _, b := range reportdurations {
        if b == req.Interval {
            validInterval = true
        }
    }

	var err error
	if validInterval {
		var dates []*gencalc.Period
		dates, err = s.GetDates(ctx, req.Interval, req.Duration)
		if err != nil {
			return nil, fmt.Errorf("Error parsing time in GetDates: %s\n", err)
		}
		//find region from facility config client
		//dummy
		var region = ""
		carbonReports, err = s.GetEmissions(ctx, dates, req.Interval, region)
		if err != nil {
			return nil, fmt.Errorf("Error from GetEmissions: %s\n", err)
		}
		
		var orgID uuid.UUID
		orgID, err = uuid.Parse(string(req.Org))
		if err != nil {
			return nil, fmt.Errorf("Error from parsing org id in HandleRequests: %s\n", err)
		}
		
		var controlPoints, err = s.GetPowerControlPoint(ctx, orgID, req.Agent)
		if err != nil {
			return nil, fmt.Errorf("Error from GetPowerControlPoint() in HandleRequests: %s\n", err)
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
		powerReports, err = s.GetPower(ctx, orgID, req.Duration, controlPoints, duration, req.Interval)
		if err != nil {
			return nil, fmt.Errorf("Error in GetPower: %s\n", err)
		}

		emissionReport, err = CalculateEmissionsReport(ctx, carbonReports, powerReports)
		if err != nil {
			return nil, fmt.Errorf("Error from Calculate Reports: %s\n", err)
		}
	}
	return &gencalc.AllReports{CarbonIntensityReports: carbonReports, PowerReports: powerReports,TotalEmissionReport: emissionReport}, nil
}

//R&D method, will implement this later
func (s *calcSvc) GetCarbonReport(ctx context.Context) (error) {
	return nil
}


//GetDates returns an array of dates for the storage client in order to correctly query for carbon reports
func (s *calcSvc) GetDates(ctx context.Context, intervalType string, duration *gencalc.Period) ([]*gencalc.Period, error) {

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



