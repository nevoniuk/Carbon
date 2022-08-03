package carbonara

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/crossnokaye/carbon/model"
)
var validResponse = `{
	"data": [		
		{
			"data": {
				"generated_rate": 473.7744069785323,
				"marginal_rate": 217.70461144377828
			},
			"dedup_key": "CAISO:carbon_intensity:2021-06-01T00:00:00+00:00",
			"event_type": "carbon_intensity",
			"meta": {
				"generated_emissions_source": "EGRID_u2018",
				"inserted_at": "2021-06-01T00:01:27.714075Z",
				"marginal_emissions_source": "EGRID_u2018",
				"marginal_source": "CAISO:marginal_fuel_mix:2021-06-01T00:00:00+00:00",
				"raw_start_date": "2021-06-01T00:00:00+00:00",
				"source": "generated_fuel_mix:CAISO:2021-06-01T00:00:00+00:00",
				"unit": "lbs/MWh"
			},
			"region": "CAISO",
			"start_date": "2021-06-01T00:00:00+00:00"
		}, 
		{
			"data": {
				"generated_rate":333.92305541812004,
				"marginal_rate":149.8037603921239
			},
			"dedup_key":"CAISO:carbon_intensity:EGRID_2019_eq:2021-06-01T00:00:00+00:00",
			"event_type":"carbon_intensity",
			"meta": {
				"generated_emissions_source":"EGRID_2019_eq",
				"inserted_at":"2021-06-01T00:39:25.399669Z",
				"marginal_emissions_source":"EGRID_2019_eq",
				"marginal_fuel_mix_source":"INTV_DIFF",
				"marginal_source":"CAISO:marginal_fuel_mix:2021-06-01T00:00:00+00:00",
				"raw_start_date":"2021-06-01T00:00:00+00:00",
				"source":"generated_fuel_mix:CAISO:2021-06-01T00:00:00+00:00",
				"unit":"lbs/MWh"
			},
			"region":"CAISO",
			"start_date":"2021-06-01T00:10:00+00:00"
		}
	],
	"meta": {
		"pagination": {
			"first": 1,
			"last": 1,
			"this": 1
		}
	}
}`


func TestGetEmissions(t *testing.T) {
	type fields struct {
		c   *http.Client
		key string
	}
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC).Format(timeFormat)
	var endTime = time.Date(2021, time.June, 1, 0, 10, 0, 0, time.UTC).Format(timeFormat)
	var invalidEndTime = time.Date(2021, time.February, 1, 0, 10, 0, 0, time.UTC).Format(timeFormat)
	
	downloadErr := errors.New("server error 400")
	type args struct {
		ctx      context.Context
		region   string
		startime string
		endtime  string
	}
	tests := []struct {
		name    string
		roundTripFn func(req *http.Request) *http.Response
		expectedErr error
		key string
	}{ 
		{
			name:        "valid",
			roundTripFn: downloadCarbonReport(t, validResponse, startTime, endTime),
			expectedErr: nil,
			key: os.Getenv("SINGULARITY_API_KEY"),
		},
		{
			name:        "invalid",
			roundTripFn: downloadCarbonReport(t, "", startTime, invalidEndTime),
			expectedErr: downloadErr,
			key: os.Getenv("SINGULARITY_API_KEY"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key = tt.key
			cl := New(&http.Client{Transport: roundTripFunc(tt.roundTripFn)}, key).(*client)
			ctx := context.Background()
			testRegion := model.Caiso
			if tt.expectedErr != nil {
				endTime = invalidEndTime
			}
			got, err := cl.GetEmissions(ctx, testRegion, startTime, endTime)
			if err != nil {
				if  !strings.Contains(err.Error(), tt.expectedErr.Error()) {
					t.Errorf("client.GetEmissions() error = %s, wantErr %s", err, tt.expectedErr)
					return
				}
			} 
			if got != nil {
				if len(got) != 1 {
					t.Errorf("len(carbonreports) == %v, want %v", len(got), 2)
				}
				if got[0].Duration.StartTime != startTime {
					t.Errorf("carbon reports start time is %s, want %s",got[0].Duration.StartTime , startTime)
				}
			}
		})
	}
}


func downloadCarbonReport(t *testing.T, content string, start string, end string) func(*http.Request) (*http.Response) {
	return func(req *http.Request) (*http.Response) {
		var err error = nil
		if req.URL.Scheme != "https" {
			t.Logf("got scheme %s, want https", req.URL.Scheme)
		}
		if req.URL.Host != "api.singularity.energy" {
			t.Logf("got host %s, want api.singularity.energy", req.URL.Host)
		}
		if req.URL.Path != "/v1/region_events/search" {
			t.Logf("got path %s, want /v1/region_events/search", req.URL.Path)
		}
		if !strings.Contains(req.URL.RawQuery, "region=CAISO&event_type=carbon_intensity&") {
			t.Logf("got path %s, want ", req.URL.RawQuery)
		}
		if !strings.Contains(req.URL.RawQuery, "start=") {
			t.Logf("got path %s, want ", req.URL.RawQuery)
		}
		if !strings.Contains(req.URL.RawQuery, "end=") {
			t.Logf("got path %s, want ", req.URL.RawQuery)
		}
		start, _ := time.Parse(timeFormat, start)
		end, _ := time.Parse(timeFormat, end)
		downloadErr := ServerError{fmt.Errorf("server error")}
		if !start.Before(end) {
			err = downloadErr
			t.Logf("invalid request with start %s and end %s", start, end)
		}
		var badRequest = "badrequest"
		invalidr := strings.NewReader(badRequest)
		if err != nil {
			return &http.Response {
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(invalidr),
			}
		}
		r := strings.NewReader(content)
		return &http.Response {
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(r),
		}
	}
}


type roundTripFunc func(req *http.Request) (*http.Response)

func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}

