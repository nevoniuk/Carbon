package carbonara

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
	"github.com/crossnokaye/carbon/model"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)
type (
	Client interface {
		// GetEmissions talks to the Singularity 'Search' endpoint and returns carbon intensity reports in 5 minute intervals
		GetEmissions(context.Context, string, string, string) ([]*genpoller.CarbonForecast, error)
	}
	client struct {
		c *http.Client
		key string
	}
	Outermoststruct struct {
		Data []struct {
			Data struct {
				GeneratedRate float64 `json:"generated_rate"`
				MarginalRate  float64 `json:"marginal_rate"`
				ConsumedRate  float64 `json:"consumed_rate"`
			}`json:"data"`
			StartDate string `json:"start_date"`
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
	// timeFormat is used to parse times in order to store time as ISO8601 format
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
// GetEmissions gets 5 min interval reports from the Carbonara API with pagination
func (c *client) GetEmissions(ctx context.Context, region string, startime string, endtime string) ([]*genpoller.CarbonForecast, error) {
	var reports []*genpoller.CarbonForecast
	var page = 1
	var last = 100
	for page <= last {
		carbonUrl := strings.Join([]string{cs_url, "region_events/search?", "region=", region, "&event_type=carbon_intensity&start=",
		startime, "&end=", endtime, "&per_page=1000", "&page=", strconv.Itoa(page)}, "")
		req, err := http.NewRequest("GET", carbonUrl, nil)
		if err != nil {
			return nil, fmt.Errorf("Error Making Request: %w\n", err)
		}
		req.Close = true
		req.Header.Add("Content-Type", "application/json")
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
				return nil, NoDataError{Err: fmt.Errorf("no data for Region %s", region)}
			}
			return nil, fmt.Errorf("Error Decoding JSON Response: %w[%d]\n", err, http.StatusBadRequest)
		}
		
		last = carbonData.Meta.Pagination.Last
		var start = carbonData.Data[0].StartDate
		for idx := 1; idx < len(carbonData.Data); idx++ {
			if carbonData.Data == nil {
				fmt.Errorf("nil carbon data element at index %d", idx)
				continue
			}
			data := carbonData.Data[idx]
			if data.StartDate == start {
				continue
			}
			end := data.StartDate
			reportperiod := &genpoller.Period{StartTime: start, EndTime: end}
			start = end
			report := &genpoller.CarbonForecast{GeneratedRate: data.Data.GeneratedRate, MarginalRate: data.Data.MarginalRate,
					ConsumedRate: data.Data.ConsumedRate, Duration: reportperiod, DurationType: model.Minute, Region: data.Region}
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