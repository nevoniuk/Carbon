package carbonara

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	//"strconv"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	//"os"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/data"

	"goa.design/clue/log"
)

//methods need to be created for different event types. they all use the same search endpoint
type (

	
	Client interface {
		get_emissions(ctx context.Context, region string, timeRange genpoller.Period) (res *genpoller.CarbonForecast, err error)
		get_fuels(ctx context.Context, region string, timeRange genpoller.Period) (res *genpoller.FuelsForecast, err error)
		//^for weekly reports only
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
		emission_factor             string
	}

	fuelmetadata struct {
		inserted_at    string
		raw_start_date string
		scraped_at     string
		source         string
	}
)

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
func (c *client) get_emissions(ctx context.Context, region string, timeRange genpoller.Period) (res []*genpoller.CarbonForecast, err error) {
	timeRange.StartTime = "2022-01-06T15:00:00-00:00" //for testing
	timeRange.EndTime = "2022-05-06T15:00:00-00:00" //testing
	carbonUrl := strings.Join([]string{cs_url, "region_events/search?", "region=", region, "?event_type=carbon_intensity&start=",
		timeRange.StartTime, "&end=", timeRange.EndTime}, "/") //for testing
	//TODO: add io reader instead of nil
	req, err := http.NewRequest(http.MethodGet, carbonUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", "52f0a90b3a2747dcb651f508b63e002c")
	carbonresp, err := c.httpGetRequestCall(ctx, req)
	defer carbonresp.Body.Close()

	carbonData := carbonreport{}
	err = json.NewDecoder(carbonresp.Body).Decode(&carbonData)
	if err != nil {
		log.Errorf(ctx, err, "cs client Carbon API JSON error")
		return nil, err
	}
	//we now have a collection of 5 minute reports of type "carbonevent"
	//make weekly reports
	//reports will contain 3 forecasts: hourly, daily, weekly
	//var reports []*genpoller.CarbonForecast
	//reportCounter := 0
	carbonDatahourly, carbonDatadaily := getdayhourlyreports(ctx, carbonData)
	carbonDataweekly := getweeklycarbonreport(ctx, carbonDatadaily)
	carbonDatamonthly := getmonthlycarbonreport(ctx, carbonDataweekly)
	return &genpoller.CarbonResponse{carbonDatahourly, carbonDataweekly, carbonDatamonthly}
}

func (c *client) get_fuels(ctx context.Context, region string, timeRange genpoller.Period) (res *genpoller.FuelsForecast, err error) {
	fuelUrl := strings.Join([]string{cs_url, "region_events/search", "?region=", region, "?event_type=generated_fuel_mix", "&start=",
		timeRange.StartTime, "&end=", timeRange.EndTime}, "/")
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
	return nil, nil
}

func parseDateTime(ctx context.Context, datetime string) (year int, month int, day int, hours int, minutes int) {
	//"2022-01-06T15:00:00-00:00"
	res1 := strings.Split(datetime, "T")
	res2 := strings.Split(res1[0], "-") //date
	res3 := strings.Split(res1[1], ":") //time

	year, err := strconv.Atoi(res2[0])
	if (err != nil) {
		log.Errorf(ctx, err, "error converting")
	}
	month, err2 := strconv.Atoi(res2[1])
	if (err2 != nil) {
		log.Errorf(ctx, err, "error converting")
	}

	day, err3 := strconv.Atoi(res2[2])
	if (err3 != nil) {
		log.Errorf(ctx, err, "error converting")
	}

	hours, err4 := strconv.Atoi(res3[0])
	if (err4 != nil) {
		log.Errorf(ctx, err, "error converting")
	}

	minutes, err5 := strconv.Atoi(res3[1])
	if (err5 != nil) {
		log.Errorf(ctx, err, "error converting")
	}
	return year, month, day, hours, minutes
}

func getdayhourlyreports(ctx context.Context, minutereports carbonreport) (hourlyreports []*genpoller.HourlyCarbonReports,
	 dailyreports []*genpoller.DailyCarbonReports) {
	newreport := true
	hourcounter := 0
	minutecounter := -1
	addtoreport := false
	reportcounter := 0
	//get averages of all minute report for a given hour
	//var hourlyreports []*genpoller.HourlyCarbonReports
	for _, event := range minutereports.reports {
		//get data for each event
		year, month, day, hours, minutes := parseDateTime(ctx, event.date_taken)
		//don't want to read report with the same start time as the last
		
		if minutes > minutecounter && !newreport {
			minutecounter = minutes
			addtoreport = true
		}
		if hours > hourcounter {
			hourcounter = hours 
			newreport = true
			addtoreport = false
			minutecounter = 0
		}
		if addtoreport {
			addtoreport = false
			hourlyreports[reportcounter].ConsumedRate = event.data.consumed_rate
			hourlyreports[reportcounter].GeneratedRate = event.data.generated_rate
			hourlyreports[reportcounter].MarginalRate = event.data.marginal_rate
			hourlyreports[reportcounter].Duration.EndTime = event.date_taken //overwrite end time each report
		}
		if newreport {
			newreport = false
			hourlyreports[reportcounter].Duration.StartTime = event.date_taken
			//assuming they arent 0
			hourlyreports[reportcounter].ConsumedRate = event.data.consumed_rate
			hourlyreports[reportcounter].GeneratedRate = event.data.generated_rate
			hourlyreports[reportcounter].MarginalRate = event.data.marginal_rate
			hourlyreports[reportcounter].ConsumedSource = event.metadata.consumed_emissions_source
			hourlyreports[reportcounter].MarginalSource = event.metadata.marginal_emissions_source
			hourlyreports[reportcounter].GeneratedSource = event.metadata.generated_emissions_source
			hourlyreports[reportcounter].EmissionFactor = event.metadata.emission_factor
			hourlyreports[reportcounter].Region = event.region
			reportcounter += 1
		}
		//fmt.Printf("%s\n", m.region)
	}
	reportcounter = 0
	daycounter := -1
	for i, event := range hourlyreports {
		year, month, day, hours, minutes := parseDateTime(ctx, event.date_taken)

		if day > daycounter {
			daycounter = day
			//newreport = false
			dailyreports[reportcounter].Duration.StartTime = event[i].Duration.StartTime
			//assuming they arent 0
			dailyreports[reportcounter].ConsumedRate = event[i].ConsumedRate
			dailyreports[reportcounter].GeneratedRate = event[i].GeneratedRate
			dailyreports[reportcounter].MarginalRate = event[i].MarginalRate
			dailyreports[reportcounter].ConsumedSource = event[i].metadata.consumed_emissions_source
			dailyreports[reportcounter].MarginalSource = event[i].MarginalSource
			dailyreports[reportcounter].GeneratedSource = event[i].GeneratedSource
			dailyreports[reportcounter].EmissionFactor = event[i].EmissionFactor
			dailyreports[reportcounter].Region = event[i].Region
			reportcounter +=1
		} else { //add hours
			dailyreports[reportcounter].ConsumedRate = event[i].ConsumedRate
			dailyreports[reportcounter].GeneratedRate = event[i].GeneratedRate
			dailyreports[reportcounter].MarginalRate = event[i].MarginalRate
			dailyreports[reportcounter].Duration.EndTime = event[i].Duration.StartTime //overwrite end time each report
		}
	}
	return hourlyreports, dailyreports
}

func getweeklycarbonreport(ctx context.Context, dailyreports []*genpoller.DailyCarbonReports) (weeklyreports []*genpoller.WeeklyCarbonReports) {
	//newreport := true
	//may be multiple weeks
	counter := 0 //used to keep track of 7 days in a week
	year, month, day, hours, minutes := parseDateTime(ctx, dailyreports[0].Duration.StartTime)
	startday := day
	reportcounter := 0

	for i, event := range dailyreports {
		year, month, day, hours, minutes := parseDateTime(ctx, event[i].Duration.StartTime)
		if counter < (day - startday) {
			counter = day - startday
		}
		if counter == 0 {
			//new report
			counter = 7
			weeklyreports[i].Duration.StartTime = event[i].Duration.StartTime
			weeklyreports[reportcounter].ConsumedRate = event[i].ConsumedRate
			weeklyreports[reportcounter].GeneratedRate = event[i].GeneratedRate
			weeklyreports[reportcounter].MarginalRate = event[i].MarginalRate
			weeklyreports[reportcounter].ConsumedSource = event[i].metadata.consumed_emissions_source
			weeklyreports[reportcounter].MarginalSource = event[i].MarginalSource
			weeklyreports[reportcounter].GeneratedSource = event[i].GeneratedSource
			weeklyreports[reportcounter].EmissionFactor = event[i].EmissionFactor
			weeklyreports[reportcounter].Region = event[i].Region
			
		} else {
			weeklyreports[reportcounter].ConsumedRate += event[i].ConsumedRate
			weeklyreports[reportcounter].GeneratedRate += event[i].GeneratedRate
			weeklyreports[reportcounter].MarginalRate += event[i].MarginalRate
			counter += 1
		}

	}
	return weeklyreports
}

func getmonthlycarbonreport(ctx context.Context, weeklyreports []*genpoller.WeeklyCarbonReports) (monthlyreports []*genpoller.MonthlyCarbonReports) {
	year, month, day, hours, minutes := parseDateTime(ctx, weeklyreports[0].Duration.StartTime)
	startmonth := month
	reportcounter := 0

	for i, event := range weeklyreports {
		year, month, day, hours, minutes := parseDateTime(ctx, event[i].Duration.StartTime)
		
		if month != startmonth {
			//new report
			monthlyreports[reportcounter].Duration.StartTime = event[i].Duration.StartTime
			monthlyreports[reportcounter].ConsumedRate = event[i].ConsumedRate
			monthlyreports[reportcounter].GeneratedRate = event[i].GeneratedRate
			monthlyreports[reportcounter].MarginalRate = event[i].MarginalRate
			monthlyreports[reportcounter].ConsumedSource = event[i].metadata.consumed_emissions_source
			monthlyreports[reportcounter].MarginalSource = event[i].MarginalSource
			monthlyreports[reportcounter].GeneratedSource = event[i].GeneratedSource
			monthlyreports[reportcounter].EmissionFactor = event[i].EmissionFactor
			monthlyreports[reportcounter].Region = event[i].Region
			
		} else {
			monthlyreports[reportcounter].ConsumedRate += event[i].ConsumedRate
			monthlyreports[reportcounter].GeneratedRate += event[i].GeneratedRate
			monthlyreports[reportcounter].MarginalRate += event[i].MarginalRate
			reportcounter += 1
		}

	} 
	return monthlyreports
}

