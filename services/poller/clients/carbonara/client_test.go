package carbonara

import (
	"context"
	"net/http"
	"reflect"
	"testing"
	"strings"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

//TODO: figure out how to encode API key - may come with client
var cs_url = "https://api.singularity.energy/v1/"
var (
	validreq  = `https://api.singularity.energy/v1/region_events/search?region=CAISO&event_type=carbon_intensity&start=2022-01-01T00:00:00-00:00&end=2022-01-02T00:00:00-00:00&' --header 'Content-Type: application/json' --header 'X-Api-Key: key`
	invalidreq = `https://api.singularity.energy/v1/region_events/search?region=CAISO&event_type=carbon_intensty&start=2022-01-01T00:00:00-00:00&end=2022-01-02T00:00:00-00:00&' --header 'Content-Type: application/json' --header 'X-Api-Key: key`
)

func TestGetEmissions(t *testing.T) {
	type fields struct {
		c   *http.Client
		key string
	}
	goodreq, err := http.NewRequest("GET", validreq, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	goodreq.Close = true
	goodreq.Header.Add("Content-Type", "application/json")
	goodreq.Header.Add("X-Api-Key", fields.key)
	badreq, err := http.NewRequest("GET", invalidreq, nil)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	badreq.Close = true
	badreq.Header.Add("Content-Type", "application/json")
	badreq.Header.Add("X-Api-Key", key)
	

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
			fields: fields{
				key: "blahblahblah",
			},
		},
		{
			name:        "invalid",
			roundTripFn: downloadCarbonReport(t, invalidreq),
			expectedErr: "record on line 1: wrong number of fields",
		},
		{
			name:        "error",
			roundTripFn: downloadErr(t),
			expectedErr: `download of CAISO report for 01/01/2018 failed: http://oasis.caiso.com/oasisapi/SingleZip?queryname=PRC_LMP&startdatetime=20171231T08:00-0000&enddatetime=20180101T08:00-0000&version=1&market_run_id=DAM&resultformat=6&grp_type=ALL: 500`,
		},
		{
			name:        "retry",
			roundTripFn: downloadReportRetry(t),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var key = tt.fields.key
			cl := New(&http.Client{Transport: roundTripFunc(tt.roundTripFn)}, key).(*client)
			ctx := context.Background()
			got, err := cl.GetEmissions(ctx, tt.args.region, tt.args.startime, tt.args.endtime, tt.args.reports)
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
				t.Fatal("report is nil")
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
		if req.URL.Path != "/oasisapi/SingleZip" {
			t.Errorf("got path %s, want /oasisapi/SingleZip", req.URL.Path)
		}
		req, err := http.NewRequest("GET", content, nil)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		req.Close = true
		req.Header.Add("Content-Type", "application/json")
		req.Header.Add("X-Api-Key", key)
		var buf bytes.Buffer
		w := zip.NewWriter(&buf)
		if content != "" {
			f, err := w.Create("LMP.xml")
			if err != nil {
				t.Fatalf("Create: %v", err)
			}
			_, err = f.Write([]byte(content))
			if err != nil {
				t.Fatalf("Write: %v", err)
			}
		}
		w.Close()
		return &http.Response{
			StatusCode: http.StatusOK,
			Body:       ioutil.NopCloser(&buf),
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

func downloadReportRetry(t *testing.T) func(*http.Request) *http.Response {
	retried := false
	return func(req *http.Request) *http.Response {
		if retried {
			return downloadCarbonReport(t, validreq)(req) //create way of determining if XML is complete
		}
		retried = true
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
