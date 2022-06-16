package poller

import (
	"context"
	"goa.design/clue/log"
	"database/sql"
	"time"
	//why would we need a facility config object but not for the other clients
	//fc "github.com/crossnokaye/rates/services/weather/clients/facilityconfig"
	"github.com/crossnokaye/carbon/poller/clients/carbonara"
	"github.com/crossnokaye/carbon/poller/clients/storage"
	genpoller "github.com/crossnokaye/carbon/poller/gen/data"
)

type (
	// Weather service
	//cs == carbonservice
	Service struct {
		csc carbonara.Client
		//dbc db.Client
	}
)


	var timeFormat = "2006-01-02T15:04:05-07:00"
	var dateFormat = "2006-01-02"


//add db
func New(csc carbonara.Client) *Service {
	return &Service{csc: csc}
}
//return date to read from


//input is the array of regions to get forecasts for(past 24 hours)
//returns an array of forecasts for the corresponding input array of regions
//dont technically need regions as a parameter
func (svc *Service) CarbonEmissions(ctx context.Context) (res []*genpoller.Forecas, err error) {
	regions := [14]string{ "AESO", "BPA", "CAISO", "ERCO", "IESO",
	   "ISONE", "MISO",
	    "NYSIO", "NYISO.NYCW",
		 "NYISO.NYLI", "NYISO.NYUP",
		  "PJM", "SPP", "EIA"}

	for i := 0; i < len(regions); i++ {
		//tempres is type carbon response for each region/postal code
		//checkDB
		tempres, err := svc.csc.get_emissions(ctx, regions[i], timeInterval)
		if err != nil { //handle errors when a region is not available??
			//instead of returning have a way marking that a region is not available
			//handle case when
			return nil, err
		}
		res = append(res, tempres)
	}
	//TODO: write to clickhouse
	//configure if previous reports need to be read
	//getcarbonreports
	err = SaveCarbonReports(ctx, res)
	//get aggregate
	//save aggregate
	return res, nil
}

func (svc *Service) Fuels(ctx context.Context) (res []*genpoller.Forecast2, err error) {
	for i := 0; i < len(regions); i++ {
		tempres, err := svc.csc.get_fuels(ctx, regions[i], timeInterval)
		if err != nil { //handle errors when a region is not available??
			//instead of returning have a way marking that a region is not available
			return nil, err
		}
		res = append(res, tempres)
	}
	//TODO: write to clickhouse
	return res, nil
}

//used as a helper function above to calculate sums for report
func (svc *Service) aggregateReport(ctx context.Context, carbonreport genpoller.CarbonForecast,
	 fuelreport genpoller.FuelsForecast, event string) (genpoller.AggregateData) {
		 //check db to get date where to last find sums for aggregate data
	//call get aggregate
	//call save aggregate
}

