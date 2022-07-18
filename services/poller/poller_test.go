package pollerapi

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/crossnokaye/carbon/services/poller/clients/carbonara"
	"github.com/crossnokaye/carbon/services/poller/clients/storage"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

//TODO:
/*
1. create mock reports
2. create mock reports for client_test.go */
var startTime = time.Date(2021, time.January, 1, 0, 0, 0, 0, nil)
var endTime = time.Date(2021, time.January, 2, 0, 0, 0, 0, nil)

var mockReport = &genpoller.CarbonForecast{
	GeneratedRate: 1290.00,
	MarginalRate:  123.00,
	ConsumedRate:  123.88,
	Duration:      &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: endTime.Format(timeFormat)},
	DurationType:  reportdurations[1],
	Region:        regions[0],
}

func Test_pollersrvc_CarbonEmissions(t *testing.T) {
	carbonaraErrNotFound := carbonara.ErrNotFound{Err: fmt.Errorf("Error not found")}
	carbonaraServerError := carbonara.ServerError{Err: fmt.Errorf("Server error")}
	carbonaraNoDataError := carbonara.NoData{Err: fmt.Errorf("No Data error")}
	//test input
	start := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC)
	region := regions[0]
	period := &genpoller.Period{StartTime: start.Format(timeFormat), EndTime: end.Format(timeFormat)}
	expectedOutput := []*genpoller.CarbonForecast{
		{GeneratedRate: 170.0, MarginalRate: 124.0, ConsumedRate: 120.0, Duration: period, DurationType: reportdurations[0], Region: regions[0]},
	}

	type fields struct {
		csc        carbonara.Client
		dbc        storage.Client
		ctx        context.Context
		cancel     context.CancelFunc
		startDates []string
	}
	type args struct {
		ctx    context.Context
		start  string
		end    string
		region string
	}
	tests := []struct {
		name        string
		apiErr      error
		expectedErr error
	}{
		//test cases
		{"success", nil, nil},
		{"server error", carbonaraServerError.Err, genpoller.MakeServerError(carbonaraServerError.Err)},
		{"no data", carbonaraServerError.Err, genpoller.MakeNoData(carbonaraNoDataError.Err)},
		{"other error", fmt.Errorf("other error"), fmt.Errorf("other error")},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carbonarac := carbonara.NewMock(t)

			if tt.name == "success" {

			}
			if tt.apiErr != nil {
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					//TODO check input fields to make sure that they are valid
					return nil, tt.apiErr

				})
			}
			for i := 0; i < 13; i++ {
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					return mockReport, nil

				})
			}
			ctx := context.Background()
			var svc *pollersrvc
			svc = NewPoller(ctx, carbonarac, nil)

			svc.CarbonEmissions(ctx, start.Format(timeFormat), region)
		})
	}
}

func Test_pollersrvc_AggregateData(t *testing.T) {
	/*
		give a bad time period and expected error
	*/
	type fields struct {
		csc           carbonara.Client
		dbc           storage.Client
		ctx           context.Context
		cancel        context.CancelFunc
		startDates    []string
		minuteReports [][]*genpoller.CarbonForecast
	}
	type args struct {
		ctx      context.Context
		region   string
		dates    []*genpoller.Period
		duration string
	}
	tests := []struct {
		name    string
		wantErr bool
	}{}
	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {
			ser := &pollersrvc{
				csc:           tt.fields.csc,
				dbc:           tt.fields.dbc,
				ctx:           tt.fields.ctx,
				cancel:        tt.fields.cancel,
				startDates:    tt.fields.startDates,
				minuteReports: tt.fields.minuteReports,
			}
			if err := ser.AggregateData(tt.args.ctx, tt.args.region, tt.args.dates, tt.args.duration); (err != nil) != tt.wantErr {
				t.Errorf("pollersrvc.AggregateData() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestGetDates(t *testing.T) {
	//Need to test this function because:
	//1. reports may be null
	//2. test the time durations of the reports
	//cant do anythign else because there may be missing data
	type args struct {
		ctx           context.Context
		minutereports []*genpoller.CarbonForecast
	}
	dateErr := errors.New("incorrect date error")
	//date1: ""
	testDates := &genpoller.Period{}
	tests := []struct {
		name    string
		dateErr error
	}{
		{name: "success", dateErr: nil},
		{name: "date Error", dateErr: dateErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDates(tt.args.ctx, tt.args.minutereports)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDates() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetDates() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_pollersrvc_Update(t *testing.T) {
	carbonaraErrNotFound := carbonara.ErrNotFound{Err: fmt.Errorf("Error not found")}
	carbonaraServerError := carbonara.ServerError{Err: fmt.Errorf("Server error")}
	carbonaraNoDataError := carbonara.NoData{Err: fmt.Errorf("No Data error")}
	
	type fields struct {
		csc        carbonara.Client
		dbc        storage.Client
		ctx        context.Context
		cancel     context.CancelFunc
		startDates []string
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{}
	}
	for _, tt := range tests {
		carbonarac := carbonara.NewMock(t)
		for i := 0; i < 13; i++ {
			carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
				return mockReport, nil

			})
		}



		t.Run(tt.name, func(t *testing.T) {
			s := &pollersrvc{
				csc:        tt.fields.csc,
				dbc:        tt.fields.dbc,
				ctx:        tt.fields.ctx,
				cancel:     tt.fields.cancel,
				startDates: tt.fields.startDates,
			}
			if err := s.Update(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("pollersrvc.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
