package carbonara

import (
	//"bytes"
	"context"
	"errors"
	//"errors"
	"fmt"

	//"fmt"
	"io"
	//"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
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
//TODO: figure out why download carbon report is not called and httprequest is called instead
//var response := ""{"data":[{"data":{"consumed_rate":447.52431491638515,"generated_rate":459.27507428132077,"marginal_rate":212.9772083348141},"dedup_key":"CAISO:carbon_intensity:EGRID_u2019:2021-06-01T00:00:00+00:00","event_type":"carbon_intensity","meta":{"consumed_emissions_source":"EGRID_u2019","consumed_rate_calculated_at":"2021-09-03T20:42:09.877028Z","consumed_source":"generated_fuel_mix:EIA.CISO:2021-06-01T00:00:00+00:00","generated_emissions_source":"EGRID_u2019","inserted_at":"2021-06-01T00:39:25.583183Z","marginal_emissions_source":"EGRID_u2019","marginal_fuel_mix_source":"INTV_DIFF","marginal_source":"CAISO:marginal_fuel_mix:2021-06-01T00:00:00+00:00","raw_start_date":"2021-06-01T00:00:00+00:00","source":"generated_fuel_mix:CAISO:2021-06-01T00:00:00+00:00","unit":"lbs/MWh"},"region":"CAISO","start_date":"2021-06-01T00:00:00+00:00"},{"data":{"generated_rate":333.92305541812004,"marginal_rate":149.8037603921239},"dedup_key":"CAISO:carbon_intensity:EGRID_2019_eq:2021-06-01T00:00:00+00:00","event_type":"carbon_intensity","meta":{"generated_emissions_source":"EGRID_2019_eq","inserted_at":"2021-06-01T00:39:25.399669Z","marginal_emissions_source":"EGRID_2019_eq","marginal_fuel_mix_source":"INTV_DIFF","marginal_source":"CAISO:marginal_fuel_mix:2021-06-01T00:00:00+00:00","raw_start_date":"2021-06-01T00:00:00+00:00","source":"generated_fuel_mix:CAISO:2021-06-01T00:00:00+00:00","unit":"lbs/MWh"},"region":"CAISO","start_date":"2021-06-01T00:00:00+00:00"},{"data":{"consumed_rate":382.24666200308013,"generated_rate":332.8066432185912,"marginal_rate":149.1695560903641},"dedup_key":"CAISO:carbon_intensity:EGRID_2019:2021-06-01T00:00:00+00:00","event_type":"carbon_intensity","meta":{"consumed_emissions_source":"EGRID_2019","consumed_rate_calculated_at":"2021-09-03T20:42:09.877017Z","consumed_source":"generated_fuel_mix:EIA.CISO:2021-06-01T00:00:00+00:00","generated_emissions_source":"EGRID_2019","inserted_at":"2021-06-01T00:39:25.905354Z","marginal_emissions_source":"EGRID_2019","marginal_fuel_mix_source":"INTV_DIFF","marginal_source":"CAISO:marginal_fuel_mix:2021-06-01T00:00:00+00:00","raw_start_date":"2021-06-01T00:00:00+00:00","source":"generated_fuel_mix:CAISO:2021-06-01T00:00:00+00:00","unit":"lbs/MWh"},"region":"CAISO","start_date":"2021-06-01T00:00:00+00:00"},{"data":{"generated_rate":341.1563315643191,"marginal_rate":151.7533168090365},"dedup_key":"CAISO:carbon_intensity:EGRID_2018_eq:2021-06-01T00:00:00+00:00","event_type":"carbon_intensity","meta":{"generated_emissions_source":"EGRID_2018_eq","inserted_at":"2021-06-01T00:01:27.553995Z","marginal_emissions_source":"EGRID_2018_eq","marginal_source":"CAISO:marginal_fuel_mix:2021-06-01T00:00:00+00:00","raw_start_date":"2021-06-01T00:00:00+00:00","source":"generated_fuel_mix:CAISO:2021-06-01T00:00:00+00:00","unit":"lbs/MWh"},"region":"CAISO","start_date":"2021-06-01T00:00:00+00:00"},{"data":{"generated_rate":340.04028546727193,"marginal_rate":151.12627797440229},"dedup_key":"CAISO:carbon_intensity:EGRID_2018:2021-06-01T00:00:00+00:00","event_type":"carbon_intensity","meta":{"generated_emissions_source":"EGRID_2018","inserted_at":"2021-06-01T00:01:27.641362Z","marginal_emissions_source":"EGRID_2018","marginal_source":"CAISO:marginal_fuel_mix:2021-06-01T00:00:00+00:00","raw_start_date":"2021-06-01T00:00:00+00:00","source":"generated_fuel_mix:CAISO:2021-06-01T00:00:00+00:00","unit":"lbs/MWh"},"region":"CAISO","start_date":"2021-06-01T00:00:00+00:00"},{"data":{"generated_rate":473.7744069785323,"marginal_rate":217.70461144377828},"dedup_key":"CAISO:carbon_intensity:2021-06-01T00:00:00+00:00","event_type":"carbon_intensity","meta":{"generated_emissions_source":"EGRID_u2018","inserted_at":"2021-06-01T00:01:27.714075Z","marginal_emissions_source":"EGRID_u2018","marginal_source":"CAISO:marginal_fuel_mix:2021-06-01T00:00:00+00:00","raw_start_date":"2021-06-01T00:00:00+00:00","source":"generated_fuel_mix:CAISO:2021-06-01T00:00:00+00:00","unit":"lbs/MWh"},"region":"CAISO","start_date":"2021-06-01T00:00:00+00:00"}],"meta":{"pagination":{"first":1,"last":1,"this":1}}}""
func TestGetEmissions(t *testing.T) {
	type fields struct {
		c   *http.Client
		key string
	}
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC).Format(timeFormat)
	var endTime = time.Date(2021, time.June, 1, 0, 10, 0, 0, time.UTC).Format(timeFormat)
	var invalidEndTime = time.Date(2021, time.February, 1, 0, 10, 0, 0, time.UTC).Format(timeFormat)
	//var validreq = strings.Join([]string{"https://api.singularity.energy/v1/region_events/search?region=CAISO&event_type=carbon_intensity&start=", startTime,"&end=", endTime, "&per_page=1000&page="},"")
	//var invalidreq = strings.Join([]string{"https://api.singularity.energy/v1/region_events/search?region=CAISO&event_type=carbon_intensity&start=", startTime,"&end=", invalidEndTime, "&per_page=1000&page="},"")
	downloadErr := errors.New("server error 400")
	
	type args struct {
		ctx      context.Context
		region   string
		startime string
		endtime  string
	}

	//define end result reports here and error
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
			testRegion := "CAISO"
			if tt.expectedErr != nil {
				endTime = invalidEndTime
			}
			got, err := cl.GetEmissions(ctx, testRegion, startTime, endTime)
			fmt.Println("returning from get emissions")
			println(got)
			if err != nil {
				if  err == tt.expectedErr {
					t.Errorf("client.GetEmissions() error = %s, wantErr %s", err, tt.expectedErr)
					return
				}
			} 
			if err != nil && got != nil {
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
			t.Errorf("got scheme %s, want https", req.URL.Scheme)
		}
		if req.URL.Host != "api.singularity.energy" {
			t.Errorf("got host %s, want api.singularity.energy", req.URL.Host)
		}
		if req.URL.Path != "/v1/region_events/search" {
			t.Errorf("got path %s, want /v1/region_events/search", req.URL.Path)
		}
		if !strings.Contains(req.URL.RawQuery, "region=CAISO&event_type=carbon_intensity&") {
			t.Errorf("got path %s, want ", req.URL.RawQuery)
		}
		if !strings.Contains(req.URL.RawQuery, "start=") {
			t.Errorf("got path %s, want ", req.URL.RawQuery)
		}
		if !strings.Contains(req.URL.RawQuery, "end=") {
			t.Errorf("got path %s, want ", req.URL.RawQuery)
		}
		start, _ := time.Parse(timeFormat, start)
		end, _ := time.Parse(timeFormat, end)
		downloadErr := ServerError{fmt.Errorf("server error")}
		if !start.Before(end) {
			err = downloadErr
			t.Errorf("invalid request with start %s and end %s", start, end)
		}
		var badRequest = "badrequest"
		invalidr := strings.NewReader(badRequest)
		if err != nil {
			fmt.Println("returning bad req")
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
/**
func downloadReportRetry(t *testing.T) func(*http.Request) *http.Response {
	retried := false
	return func(req *http.Request) *http.Response {
		if retried {
			return downloadCarbonReport(t, content, )(req)
		}
		retried = true
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}
	}
}
*/
