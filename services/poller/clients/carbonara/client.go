package carbonara

import (
	"context"
	"strconv"
	//"bytes"
	"encoding/json"
	"net/http"

	//"strconv"

	//"strconv"
	"fmt"
	//"math"
	//"io"
	//"io/ioutil"
	"strings"
	"time"

	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	"goa.design/clue/log"
)
var reportdurations [6]string
//methods need to be created for different event types. they all use the same search endpoint
type (

	Client interface {
		GetEmissions(context.Context, string, string, string, []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error)
	}

	client struct {
		c *http.Client
	}
	
	
	//some numbers and events are estimated**
	Outermoststruct struct {
		Data []struct {
			Data struct {
				Generated_rate float64 `json:"generated_rate"`
				Marginal_rate  float64 `json:"marginal_rate"`
				Consumed_rate  float64 `json:"consumed_rate"`
			}`json:"data"`
			Meta struct {
				Generated_emissions_source  string `json:"generated_emissions_source"`
			}`json:"meta"`
			Start_date string `json:"start_date"`
			Region     string `json:"region"`
		}`json:"data"`
		Meta struct {
			Pagination struct {
				Last int `json:"last"`
				This int `json:"this"`
			}`json:pagination`
		}`json:"meta"`
	}
)

const (
	//format is year-month-daysThours:minutes:seconds:something:something
	timeFormat = "2006-01-02T15:04:05-07:00"
	dateFormat = "2006-01-02"
	cs_url     = "https://api.singularity.energy/v1/"
)

func (c *client) Init() {
	fmt.Printf("initialized")
	reportdurations = [...]string{ "minute", "hourly", "daily", "weekly", "monthly", "yearly"}
}

func New(c *http.Client) Client {
	c.Timeout = 10 * time.Second
	return &client{c}
}

func (c *client) HttpGetRequestCall(ctx context.Context, req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	//retry
	if err != nil || resp.StatusCode != http.StatusOK {
		retries := 0
		for (err != nil || resp.StatusCode != http.StatusOK) && retries < 3 {
			time.Sleep(time.Duration(retries) * time.Second)
			resp, err = http.DefaultClient.Do(req)
			retries++
		}
	}
	//return error if the "DO" action fails
	if err != nil {
		log.Errorf(ctx, err, "carbon client API Get error")
		return resp, err
	}
	//return the exact error code
	if resp.StatusCode != http.StatusOK {
		log.Errorf(ctx, err, "%d", resp.StatusCode)
		return resp, err
	}

	return resp, nil
}

//Bug last hourly report is not retrieved
func (c *client) GetEmissions(ctx context.Context, region string, startime string, endtime string, reports []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
	//ignore starttime and endtime for now
	//fmt.Println("start is")
	//fmt.Println(startime)
	//fmt.Println("end is")
	//fmt.Println(endtime)
	
	var page = 1
	var last = 100
	var report *genpoller.CarbonForecast
	var reportperiod *genpoller.Period

	//var reports []*genpoller.CarbonForecast

	for page <= last {
		carbonUrl := strings.Join([]string{cs_url, "region_events/search?", "region=", region, "&event_type=carbon_intensity&start=",
		startime, "&end=", endtime, "&per_page=1000", "&page=", strconv.Itoa(page)}, "") //for testing
	
		fmt.Println(carbonUrl)
		//TODO: add io reader instead of nil
		req, err := http.NewRequest("GET", carbonUrl, nil)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		//close request to prevent EOF
		req.Close = true
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Api-Key", "f74f6d5bafa942dcab07dd8e2a0af564")

		carbonresp, err := c.HttpGetRequestCall(ctx, req)
		if err != nil {
			fmt.Errorf("ERROR FROM GET REQUEST: %s", err)
			return nil, err
		}
		//nil report
		if carbonresp.ContentLength < 100 {
			return nil, fmt.Errorf("No data available for region %s\n", region)
		}

		defer carbonresp.Body.Close()

		var carbonData Outermoststruct
		//var finalcarbonData []Outermoststruct
		err = json.NewDecoder(carbonresp.Body).Decode(&carbonData)
		if err != nil {
			log.Errorf(ctx, err, "cs client Carbon API JSON error")
			return nil, err
		}
		//fmt.Print("response")
		//fmt.Println(carbonData)
		var count = 0
		last = carbonData.Meta.Pagination.Last
		//fmt.Printf("there are %d pages\n", carbonData.Meta.Pagination.Last)
		var start = carbonData.Data[0].Start_date
		var end string
		//iterate though the page returned to make carbon forecasts
		//maybe not do this part or only read half the reports
		for count < len(carbonData.Data) {
			
			end = carbonData.Data[count].Start_date
			if start != end {
				reportperiod = &genpoller.Period{StartTime: start, EndTime: end}
				start = end
				report = &genpoller.CarbonForecast{GeneratedRate: carbonData.Data[count].Data.Generated_rate, MarginalRate: carbonData.Data[count].Data.Marginal_rate,
					ConsumedRate: carbonData.Data[count].Data.Consumed_rate, Duration: reportperiod, DurationType: reportdurations[0], GeneratedSource: carbonData.Data[count].Meta.Generated_emissions_source, Region: carbonData.Data[count].Region}
				//fmt.Println(report.Duration)
					reports = append(reports, report)
			}
			count += 1
		}
		//fmt.Printf("THIS PAGE IS %d\n", carbonData.Meta.Pagination.This)
		if carbonData.Meta.Pagination.This == carbonData.Meta.Pagination.Last {
			//fmt.Println("reached last report")
			return reports, nil
		}
		page += 1
	}
	return reports, nil
}

/**
func gethourlyreports(ctx context.Context, minutereports Outermoststruct, hourlyreports []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast) {
	//get averages of all minute report for a given hour

	newreport := false
	addtoreport := false
	
	//keep track of which reports have consumed, gen, and marg data
	var consumedcounter float64
	var gencounter float64
	var margcounter float64

	//var hourlyreports []*genpoller.CarbonForecast

	var hourlyreport *genpoller.CarbonForecast
	var hourlyreportperiod *genpoller.Period
	var t, err = time.Parse(timeFormat, minutereports.Data[0].Start_date)

	if err != nil {
		fmt.Printf("error %p", t)
	}
	
	//initialize counters
	var yearcounter = t.Year()
	
	var monthcounter = t.Month()
	
	var daycounter = t.Day()
	
	var minutecounter = t.Minute()
	
	var hourcounter = t.Hour()
	
	//for the first report only
	//null deference
	
	var Start = minutereports.Data[0].Start_date
	var End string

	var consumed float64
	var marginal float64
	var generated float64

	var GeneratedSource string
	var Region string

	for i, event := range minutereports.Data {

		var t, err = time.Parse(timeFormat, event.Start_date)
		if err != nil {
			fmt.Printf("error %p", t)
		}
		var year = t.Year()
		var month = t.Month()
		var day =  t.Day()
		var hours = t.Hour()
		var minutes = t.Minute()

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

			if (event.Data.Consumed_rate != float64(0)) {
				consumed += event.Data.Consumed_rate
				consumedcounter += 1.0
			}

			if (event.Data.Generated_rate != float64(0)) {
				generated += event.Data.Generated_rate
				gencounter += 1.0
			}

			if (event.Data.Marginal_rate != float64(0)) {
				marginal += event.Data.Marginal_rate

				margcounter += 1.0
			}
			End = event.Start_date //overwrite end time each report
		}

		if newreport || (i == len(minutereports.Data) - 1) {
			newreport = false
			//previous report
			if consumed != float64(0) && consumedcounter != float64(0) {
				consumed = consumed / consumedcounter
			}
			//fmt.Printf("%f", ConsumedRate)
			if generated != float64(0) && gencounter != float64(0) {
				generated = generated / gencounter
			}
			if marginal != float64(0) && margcounter != float64(0) {
				marginal = float64(marginal / margcounter)
			}

			GeneratedSource = event.Meta.Generated_emissions_source
			Region = event.Region
			//append report

			hourlyreportperiod = &genpoller.Period{StartTime: Start, EndTime: End}
			fmt.Printf(" report period is %+v\n", hourlyreportperiod)
			hourlyreport = &genpoller.CarbonForecast{GeneratedRate: generated, MarginalRate: marginal,
			ConsumedRate: consumed, Duration: hourlyreportperiod, GeneratedSource: GeneratedSource, Region: Region}
			fmt.Printf("hourly report is %+v\n", hourlyreport)
			hourlyreports = append(hourlyreports, hourlyreport)

			//reset values
			consumed = float64(0)
			consumedcounter = float64(0)
			margcounter = float64(0)
			marginal = float64(0)
			gencounter = float64(0)
			generated = float64(0)

			//newreport
			Start = event.Start_date
		}
	}

	return hourlyreports
}
*/
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
