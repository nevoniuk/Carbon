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
func CalculateReports(ctx context.Context, carbonReports []*gencalc.CarbonReport, powerReports []*gencalc.ElectricalReport) ([]*gencalc.EmissionsReport, error) {
	var reports []*gencalc.EmissionsReport
	//assume that each report has the same duration and interval type
	for i, report := range powerReports {
		
	}
	//input =carbon reports, 
	//1 MWh = 1000 KWh
	//1.convert from MWh to KWh
	//2.
	return report, nil
}

//uses store to get input for past-values service
func (s *calcSvc) GetPowerControlPoint(ctx context.Context, org uuid.UUID, agent string) ([]uuid.UUID, error) {
	var temp []uuid.UUID
	if org == uuid.Nil {
		return temp, fmt.Errorf("Org ID is null\n")
	}
	
	if agent == "" {
		return temp, fmt.Errorf("Agent ID is null\n")
	}

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
func (s *calcSvc) GetEmissions(ctx context.Context, dateRange *gencalc.Period, interval string, region string) ([]*gencalc.CarbonReport, error) {
	var reports []*gencalc.CarbonReport
	return reports, nil
}

//HandleRequests will output the CO2 intensity, Power Meter, and resulting CO2 emission reports
func (s *calcSvc) HandleRequests(ctx context.Context, req *gencalc.RequestPayload) (*gencalc.AllReports, error) {
	var emissionReports []*gencalc.EmissionsReport
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
		//find region from facility config client
		//dummy
		var region = ""
		carbonReports, err = s.GetEmissions(ctx, req.Duration, req.Interval, region)
		if err != nil {
			return nil, fmt.Errorf("Error from GetEmissions: %s\n", err)
		}
		//should not need to know the point name
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

		emissionReports, err = CalculateReports(ctx, carbonReports, powerReports)
		if err != nil {
			return nil, fmt.Errorf("Error from Calculate Reports: %s\n", err)
		}
	}
	return &gencalc.AllReports{CarbonIntensityReports: carbonReports, PowerReports: powerReports,TotalEmissionReports: emissionReports}, nil
}

//R&D method
func (s *calcSvc) GetCarbonReport(ctx context.Context) (error) {
	return nil
}



