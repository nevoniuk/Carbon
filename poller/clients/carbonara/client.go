package carbonara

import (
	"context"
	"encoding/json"
	"net/http"

	//"strconv"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	//"os"
	gencarbon "gen/data"

	"goa.design/clue/log"
)

//methods need to be created for different event types. they all use the same search endpoint
type (
	Client interface {
		get_emissions(ctx context.Context, region string, timeRange gencarbon.Period) (res *gencarbon.CarbonForecast, err error)
		get_fuels(ctx context.Context, region string, timeRange gencarbon.Period) (res *gencarbon.FuelsForecast, err error)
		get_aggregate_data(ctx context.Context, region string, timeRange gencarbon.Period, event string) (res *gencarbon.AggregateData, err error)
	}

	client struct {
		c *http.Client
	}

	//some numbers and events are estimated**
	fuelreport struct {
		reports []fuelevent
	}
	carbonreport struct {
		reports []carbonevent
	}
	fuelevent struct {
		data       fueldata
		metadata   fuelmetadata
		region     string
		date_taken string
	}
	carbonevent struct {
		data       carbondata
		metadata   carbonmetadata
		date_taken string
		region     string
	}
	carbondata struct {
		generated_rate float64
		marginal_rate  float64
		consumed_rate  float64
	}
	fueldata struct {
		coal_mw        int
		hydro_mw       int
		natural_gas_mw int
		nuclear_mw     int
		other_mw       int
		petroluem_mw   int
		solar_mw       int
		wind_mw        int
	}
	carbonmetadata struct {
		consumed_emissions_source   string
		consumed_rate_calculated_at string
		consumed_source             string
		generated_emissions_source  string
		inserted_at                 string
		marginal_emissions_source   string
		raw_start_date              string
		source                      string
		unit                        string
		updated_at                  string
	}

	fuelmetadata struct {
		inserted_at    string
		raw_start_date string
		scraped_at     string
		source         string
	}
)

/**
type initialrates struct {
	generated_rate float64
	marginal_rate float64
}

type CarbonIntensity struct {
	carbonrates CarbonIntensityData
}

type CarbonIntensityData struct {
	Minutes int
	generated_rate float64
	marginal_rate float64
}
type CarbonMetaData struct {
	emissionfactors ArrayOf(Emissionfactor)
	source string
}
type Forecast struct {
	carbon
	ArrayOf(timeReport)
}

type carbonReport struct {
	CarbonIntensityData
	dedupkey string
	eventtype string
	data carbonreportData
	region string
	start_date string
}

type fuelReport struct {
	FuelIntensityData
	dedupkey string
	eventtype string
	data carbonreportData
	region string
	start_date string
}

type fueldmixreport struct {
	data mixdata
	meta fuelReportData
	start_date string
}


type carbonreportData Struct {
	timeInterval int
	genemissionssource string
	inserted_date string
	margemissionssource string
	margsource string
	source string
	unit string
	updated_date string
}


type FuelIntensityData struct {
	coalmw int
	timeinterval int
	naturalgasmw int
	nuclearmw int
	othermw int
	solarmw int
	windmw int
}

type mixdata struct {
	coalmw int
	naturalgasmw int
	nuclearmw int
	othermw int
	solarmw int
	windmw int
}

type fuelReportData struct {
	source string
}
*/
const (
	timeFormat = "2006-01-02T15:04:05-07:00"
	dateFormat = "2006-01-02"
	cs_url     = "https://api.singularity.energy/v1/"
)

func New(c *http.Client) Client {
	c.Timeout = 10 * time.Second
	return &client{c}
}

func (c *client) httpGetRequestCall(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		retries := 0
		for (err != nil || resp.StatusCode != http.StatusOK) && retries < 3 {
			time.Sleep(time.Duration(retries) * time.Second)
			resp, err = http.DefaultClient.Do(req)
			retries++
		}
	}
	if err != nil {
		log.Errorf(ctx, err, "carbon client API Get error")
		return resp, err
	}
	if resp.StatusCode != http.StatusOK {
		log.Errorf(ctx, err, "%d", resp.StatusCode)
		return resp, err
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	fmt.Printf("response body: %s\n", resBody)
	return resp, nil
}

//getforecast method should make an http call for each region
func (c *client) get_emissions(ctx context.Context, region string, timeRange gencarbon.Period) (res *gencarbon.Forecast, err error) {
	timeRange.startTime = "2022-01-06T15:00:00-00:00"
	timeRange.endTime = "2022-05-06T15:00:00-00:00"
	carbonUrl := strings.Join([]string{cs_url, "region_events/search?", "region=", region, "?event_type=carbon_intensity&start=",
		timeRange.startTime, "&end=", timeRange.endTime}, "/")
	//TODO: add io reader instead of nil
	req, err := http.NewRequest(http.MethodGet, carbonUrl, nil)
	if err != nil {
		//return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", "52f0a90b3a2747dcb651f508b63e002c")
	carbonresp, err := c.httpGetRequestCall(ctx, req)
	defer carbonresp.Body.Close()

	carbonData := carbonreport{}
	err = json.NewDecoder(carbonresp.Body).Decode(&carbonData)
	if err != nil {
		log.Errorf(ctx, err, "cs client Carbon API JSON error")
		//return nil, err
	}

	for {
		var m carbonevent
		fmt.Printf("%s\n", m.region)
	}
	return nil
}

func (c *client) get_fuels(ctx context.Context, region string, timeRange gencarbon.Period) (res *gencarbon.Forecast, err error) {
	fuelUrl := strings.Join([]string{cs_url, "region_events/search", "?region=", region, "?event_type=generated_fuel_mix", "&start=",
		timeRange.startTime, "&end=", timeRange.endTime}, "/")
	//TODO: add io reader instead of nil
	req, err := http.NewRequest(http.MethodGet, fuelUrl, nil)
	if err != nil {
		//return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", "52f0a90b3a2747dcb651f508b63e002c")
	fuelresp, err := c.httpGetRequestCall(ctx, req)
	defer fuelresp.Body.Close()

	fuelData := fuelreport{}
	err = json.NewDecoder(fuelresp.Body).Decode(&fuelData)
	if err != nil {
		log.Errorf(ctx, err, "cs client Fuel API JSON error")
		//return nil, err

	}
	return nil
}

func (c *client) get_aggregate_data(ctx context.Context, region string, timeRange gencarbon.Period) (res *gencarbon.aggregateData, err error) {
	return nil
}
