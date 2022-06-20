package carbonara

import (
	"context"
	//"bytes"
	"encoding/json"
	"net/http"
	"strconv"

	//"strconv"
	"fmt"
	"io/ioutil"
	"strings"
	"time"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	"goa.design/clue/log"
)

//methods need to be created for different event types. they all use the same search endpoint
type (

	Client interface {
		//Init()
		//HttpGetRequestCall(context.Context, *http.Request) (*http.Response, error)
		GetEmissions(context.Context, string, time.Time, time.Time) ([]*genpoller.CarbonForecast, error)
	}

	client struct {
		c *http.Client
	}

	//some numbers and events are estimated**
	carbonreport struct {
		reports []carbonevent
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
	
)

const (
	timeFormat = "2006-01-02T15:04:05-07:00"
	dateFormat = "2006-01-02"
	cs_url     = "https://api.singularity.energy/v1/"
)

func (c *client) Init() {
	fmt.Printf("initialized")
}

func New(c *http.Client) Client {
	c.Timeout = 10 * time.Second
	return &client{c}
}

func (c *client) HttpGetRequestCall(ctx context.Context, req *http.Request) (*http.Response, error) {
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


func (c *client) GetEmissions(ctx context.Context, region string, startime time.Time, endtime time.Time) ([]*genpoller.CarbonForecast, error) {
	//ignore starttime and endtime for now

	start := "2022-01-06T15:00:00-00:00" //for testing
	end := "2022-05-06T15:00:00-00:00" //testing
	carbonUrl := strings.Join([]string{cs_url, "region_events/search?", "region=", region, "?event_type=carbon_intensity&start=",
		start, "&end=", end}, "/") //for testing
	//TODO: add io reader instead of nil
	req, err := http.NewRequest(http.MethodGet, carbonUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Api-Key", "52f0a90b3a2747dcb651f508b63e002c")
	carbonresp, err := c.HttpGetRequestCall(ctx, req)
	defer carbonresp.Body.Close()

	carbonData := carbonreport{}
	err = json.NewDecoder(carbonresp.Body).Decode(&carbonData)
	if err != nil {
		log.Errorf(ctx, err, "cs client Carbon API JSON error")
		return nil, err
	}
	//carbonDatahourly, carbonDatadaily := getdayhourlyreports(ctx, carbonData)
	//carbonDataweekly := getweeklycarbonreport(ctx, carbonDatadaily)
	//carbonDatamonthly := getmonthlycarbonreport(ctx, carbonDataweekly)
	//return &genpoller.CarbonResponse{carbonDatahourly, carbonDatadaily, carbonDataweekly, carbonDatamonthly}, err
	return gethourlyreports(ctx, carbonData), nil
}


func gethourlyreports(ctx context.Context, minutereports carbonreport) ([]*genpoller.CarbonForecast) {
	//get averages of all minute report for a given hour
	newreport := true
	addtoreport := false

	hourcounter := 0
	minutecounter := -1
	daycounter := -1
	monthcounter := -1
	yearcounter := -1

	
	
	//keep track of which reports have consumed, gen, and marg data
	var consumedcounter float64 = 0
	var gencounter float64 = 0
	var margcounter float64 = 0

	var hourlyreports []*genpoller.CarbonForecast
	var hourlyreport *genpoller.CarbonForecast

	year, month, day, hours, minutes := parseDateTime(ctx, minutereports.reports[0].date_taken)
	yearcounter = year
	monthcounter = month
	daycounter = day
	minutecounter = minutes
	hourcounter = hours

	//for the first report only
	hourlyreport.Duration.StartTime = minutereports.reports[0].date_taken
	for _, event := range minutereports.reports {
		year, month, day, hours, minutes := parseDateTime(ctx, event.date_taken)

		//don't want to read report with the same start time as the last
		if minutes > minutecounter {
			minutecounter = minutes
			addtoreport = true
		}
		//have to make sure we dont read data from the previous day because it will mess up the minute counter
		if hours > hourcounter {
			hourcounter = hours 
			newreport = true
			addtoreport = false
			minutecounter = -1
		}
		if month < monthcounter && year <= yearcounter {
			addtoreport = false
		} else if day < daycounter && month <= monthcounter {
			addtoreport = false
		} else if hours < hourcounter && day <= daycounter {
			addtoreport = false
		}
 
		if addtoreport {
			addtoreport = false
			if (event.data.consumed_rate != 0) {
				hourlyreport.ConsumedRate += event.data.consumed_rate
				consumedcounter += 1
			}
			if (event.data.generated_rate != 0) {
				hourlyreport.GeneratedRate += event.data.generated_rate
				gencounter += 1
			}
			if (event.data.marginal_rate != 0) {
				hourlyreport.MarginalRate += event.data.marginal_rate
				margcounter += 1
			}
			hourlyreport.Duration.EndTime = event.date_taken //overwrite end time each report
		}

		if newreport {
			newreport = false
			//previous report
			hourlyreport.ConsumedRate = (hourlyreport.ConsumedRate / consumedcounter)
			hourlyreport.GeneratedRate = (hourlyreport.GeneratedRate / gencounter)
			hourlyreport.MarginalRate = (hourlyreport.MarginalRate / margcounter)
			hourlyreport.ConsumedSource = event.metadata.consumed_emissions_source
			hourlyreport.MarginalSource = event.metadata.marginal_emissions_source
			hourlyreport.GeneratedSource = event.metadata.generated_emissions_source
			hourlyreport.EmissionFactor = event.metadata.emission_factor
			hourlyreport.Region = event.region
			//append report
			hourlyreports = append(hourlyreports, hourlyreport)
			//newreport
			hourlyreport.Duration.StartTime = event.date_taken
		}
	}
	return hourlyreports
}

/**
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
	for _, event := range hourlyreports {
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
*/
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

/*
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
	
	fuelreport struct {
		reports []fuelevent
	}
	fuelmetadata struct {
		inserted_at    string
		raw_start_date string
		scraped_at     string
		source         string
	}
	fuelevent struct {
		data       fueldata
		metadata   fuelmetadata
		region     string
		date_taken string
	}

	func (c *client) getfuels(ctx context.Context, region string, timeRange genpoller.Period) (res *genpoller.FuelsForecast, err error) {
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
	*/
