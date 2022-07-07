package calcapi
//something
import (
	"context"
	"fmt"
	//"fmt"
	//"sync"
	"time"

	"github.com/crossnokaye/carbon/services/calc/clients/power"
	"github.com/crossnokaye/carbon/services/calc/clients/power_server"
	"github.com/crossnokaye/carbon/services/calc/clients/storage"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
)

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
func GetControlPoints(context.Context, *gencalc.PastValuesPayload) ([]string, error) {

}

//wrapper function for talking to power client
//power meters at riverside, oxnard
//0's more than a minute resemble blackout


func GetPower(context.Context, *gencalc.GetPowerPayload) (*gencalc.ElectricalReport, error) {

}

//wrapper function for talking to storage client
func GetEmissions(context.Context, *gencalc.RequestPayload) (*gencalc.CarbonReport, error) {

}

func (s *calcSvc) HandleRequests(context.Context, *gencalc.RequestPayload) (error) {

}

//R&D method
func (s *calcSvc) Carbonreport(context.Context) (err error) {
	//gets reports in carbon forecasts

}