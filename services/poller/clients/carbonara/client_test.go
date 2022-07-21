package carbonara

import (
	"context"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
	"time"
	"errors"

	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

//TODO: figure out how to encode API key - may come with client

func TestGetEmissions(t *testing.T) {
	type fields struct {
		c   *http.Client
		key string
	}
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, nil).Format(timeFormat)
	var endTime = time.Date(2021, time.June, 1, 0, 10, 0, 0, nil).Format(timeFormat)
	var invalidEndTime = time.Date(2021, time.February, 1, 0, 10, 0, 0, nil).Format(timeFormat)

	var validreq = strings.Join([]string{"https://api.singularity.energy/v1/region_events/search?region=CAISO&event_type=carbon_intensity&start=", startTime,"&end=", endTime, "&' --header 'Content-Type: application/json' --header 'X-Api-Key:"},"")
	var invalidreq = strings.Join([]string{"https://api.singularity.energy/v1/region_events/search?region=CAISO&event_type=carbon_intensity&start=", startTime,"&end=", invalidEndTime, "&' --header 'Content-Type: application/json' --header 'X-Api-Key:"},"")
	downloadErr := errors.New("server error")
	type args struct {
		ctx      context.Context
		region   string
		startime string
		endtime  string
		reports  []*genpoller.CarbonForecast
	}

	//define end result reports here and error
	tests := []struct {
		name    string
		fields fields
		roundTripFn func(req *http.Request) *http.Response
		expectedErr string
	}{ 
		{
			name:        "valid",
			roundTripFn: downloadCarbonReport(t, validreq),
		},
		{
			name:        "invalid",
			roundTripFn: downloadCarbonReport(t, invalidreq),
			expectedErr: downloadErr.Error(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key = tt.fields.key
			cl := New(&http.Client{Transport: roundTripFunc(tt.roundTripFn)}, key).(*client)
			ctx := context.Background()
			testRegion := "CAISO"
			got, err := cl.GetEmissions(ctx, testRegion, startTime, endTime, nil)
			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("GetEmissions did not return an error")
				} else if err.Error() != tt.expectedErr {
					t.Errorf("client.GetEmissions() error = %s, wantErr %s", err.Error(), tt.expectedErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("GetEmissions returned error: %v", err)
			}
			if got == nil {
				t.Fatal("reports are nil")
			}
			if len(got) != 2 {
				t.Errorf("len(carbonreports) == %v, want %v", len(got), 2)
			}
			if got[0].Duration.StartTime != startTime {
				t.Errorf("carbon reports start time is %s, want %s",got[0].Duration.StartTime , startTime)
			}
			//possible other checks
		})
	}
}


func downloadCarbonReport(t *testing.T, content string) func(*http.Request) *http.Response {
	return func(req *http.Request) *http.Response {
		if req.URL.Scheme != "http" {
			t.Errorf("got scheme %s, want http", req.URL.Scheme)
		}
		if req.URL.Host != "https://api.singularity.energy/v1/" {
			t.Errorf("got host %s, want https://api.singularity.energy/v1/", req.URL.Host)
		}
		if req.URL.Path != "region_events/search?" {
			t.Errorf("got path %s, want /oasisapi/SingleZip", req.URL.Path)
		}
		if !strings.Contains(req.URL.RawQuery, "application/json") {
			t.Errorf("got path %s, want ", req.URL.Path)
		}
		r := strings.NewReader(content)
		return &http.Response {
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(r),
		}
	}
}
//need some kind of downloadCarbonReport


func downloadErr(t *testing.T) func(*http.Request) *http.Response {
	return func(req *http.Request) *http.Response {
		return &http.Response{
			StatusCode: http.StatusInternalServerError,
		}
	}
}

type roundTripFunc func(req *http.Request) *http.Response

// RoundTrip .
func (f roundTripFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req), nil
}
