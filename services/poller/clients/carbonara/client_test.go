package carbonara

import (
	"context"
	"net/http"
	"reflect"
	"testing"

	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

//TODO: figure out how to encode API key - may come with client
var (
	validreq  = `INTERVALSTARTTIME_GMT,INTERVALENDTIME_GMT,OPR_DT,OPR_HR,OPR_INTERVAL,NODE_ID_XML,NODE_ID,NODE,MARKET_RUN_ID,LMP_TYPE,XML_DATA_ITEM,PNODE_RESMRID,GRP_TYPE,POS,MW,GROUP
2022-02-01T15:00:00-00:00,2022-02-01T16:00:00-00:00,2022-02-01,8,0,0096WD_7_N001,0096WD_7_N001,0096WD_7_N001,DAM,LMP,LMP_PRC,0096WD_7_N001,ALL,1,74.63,1`
	invalidreq = `INTERVALSTARTTIME_GMT,INTERVALENDTIME_GMT,OPR_DT,OPR_HR,OPR_INTERVAL,NODE_ID_XML,NODE_ID,NODE,MARKET_RUN_ID,LMP_TYPE,XML_DATA_ITEM,PNODE_RESMRID,POS,MW,GROUP
2022-02-01T15:00:00-00:00,2022-02-01T16:00:00-00:00,2022-02-01,8,0,0096WD_7_N001,0096WD_7_N001,0096WD_7_N001,DAM,LMP,LMP_PRC,0096WD_7_N001,ALL,1,74.63,1`
)

func TestGetEmissions(t *testing.T) {

	type fields struct {
		c   *http.Client
		key string
	}

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
		fields  fields
		args    args
		roundTripFn func(req *http.Request) *http.Response
		want    []*genpoller.CarbonForecast
		expectedErr string
		wantErr bool
	}{ 
		{
			name:        "valid",
			roundTripFn: downloadCarbonReport(t, validreq),
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
			cl := New(&http.Client{Transport: roundTripFunc(tt.roundTripFn)}, tt.fields.key).(*client)
			got, err := cl.GetEmissions(tt.args.ctx, tt.args.region, tt.args.startime, tt.args.endtime, tt.args.reports)
			
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.GetEmissions() = %v, want %v", got, tt.want)
			}
			if tt.expectedErr != "" {
				if err == nil {
					t.Errorf("DownloadCAISO did not return an error")
				}
				if (err != nil) != tt.wantErr {
					t.Errorf("client.GetEmissions() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
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
		//TODO: ask raphael what to do here
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
