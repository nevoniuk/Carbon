package carbonara

import (
	"context"
	"strconv"
	"encoding/json"
	"net/http"
	"fmt"
	"strings"
	"time"
	"errors"
	"io"

	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	"goa.design/clue/log"
)

//reportdurations maintains the interval length of each report
var reportdurations [6]string = [6]string{ "minute", "hourly", "daily", "weekly", "monthly", "yearly"}

type (
	Client interface {
		GetEmissions(context.Context, string, string, string) ([]*genpoller.CarbonForecast, error)
	}
	client struct {
		c *http.Client
		key string
	}
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
	ServerError struct{ Err error }
	NoDataError struct{ Err error }
)

const (
	//timeFormat is used to parse times in order to store time as ISO8601 format
	timeFormat = "2006-01-02T15:04:05-07:00"
	cs_url     = "https://api.singularity.energy/v1/"
)

func New(c *http.Client, key string) Client {
	c.Timeout = 10 * time.Second
	return &client{c, key}
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
		var serverError = ServerError{Err: fmt.Errorf("server error %d", resp.StatusCode)}
		return resp, serverError.Err
	}
	
	if resp.StatusCode != http.StatusOK {
		var serverError = ServerError{Err: fmt.Errorf("server error %d", resp.StatusCode)}
		return resp, serverError.Err
	}

	return resp, nil
}
//GetEmissions gets 5 min interval reports from the Carbonara API with pagination
func (c *client) GetEmissions(ctx context.Context, region string, startime string, endtime string) ([]*genpoller.CarbonForecast, error) {
	var reports []*genpoller.CarbonForecast
	var page = 1
	var last = 100 //dummy value
	for page <= last {
		carbonUrl := strings.Join([]string{cs_url, "region_events/search?", "region=", region, "&event_type=carbon_intensity&start=",
		startime, "&end=", endtime, "&per_page=1000", "&page=", strconv.Itoa(page)}, "")
		req, err := http.NewRequest("GET", carbonUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("Error Making Request: %s\n", err)
		}
		req.Close = true
		req.Header.Add("Content-Type", "application/json")
		//replace with key after testing
		req.Header.Add("X-Api-Key", c.key)
		carbonresp, err := c.HttpGetRequestCall(ctx, req)
		if err != nil {
			return nil, err
		}
		defer carbonresp.Body.Close()
		var carbonData Outermoststruct
		err = json.NewDecoder(carbonresp.Body).Decode(&carbonData)
		if err != nil {
			if errors.Is(err, io.EOF) {
				var noDataError = NoDataError{Err: fmt.Errorf("no data for Region %s", region)}
				return nil, noDataError.Err
			}
			return nil, fmt.Errorf("Error Decoding JSON Response: %s[%d]\n", err, http.StatusBadRequest)
		}
		
		last = carbonData.Meta.Pagination.Last
		var start = carbonData.Data[0].Start_date
		for idx := 1; idx < len(carbonData.Data); idx++ {
			if carbonData.Data == nil {
				log.Infof(ctx, "nil carbon data element at index %d", idx)
				continue
			}
			data := carbonData.Data[idx]
			if data.Start_date == start {
				continue
			}
			end := data.Start_date
			reportperiod := &genpoller.Period{StartTime: start, EndTime: end}
			start = end
			report := &genpoller.CarbonForecast{GeneratedRate: data.Data.Generated_rate, MarginalRate: data.Data.Marginal_rate,
					ConsumedRate: data.Data.Consumed_rate, Duration: reportperiod, DurationType: reportdurations[0], Region: data.Region}
			reports = append(reports, report)
			
		}
		if carbonData.Meta.Pagination.This == carbonData.Meta.Pagination.Last {
			return reports, nil
		}
		page += 1
	}
	return reports, nil
}

func (err ServerError) Error() string { return err.Err.Error() }
func (err NoDataError) Error() string { return err.Err.Error() }