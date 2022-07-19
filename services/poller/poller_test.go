package pollerapi

import (
	"context"
	"errors"
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
var testPeriod = &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: endTime.Format(timeFormat)}
var testRegion = regions[0]
var mockReport = &genpoller.CarbonForecast{
	GeneratedRate: 1290.00,
	MarginalRate:  123.00,
	ConsumedRate:  123.88,
	Duration:      testPeriod,
	DurationType:  reportdurations[1],
	Region:        testRegion,
}
var mockReports []*genpoller.CarbonForecast
var mockMinuteReportOne = &genpoller.CarbonForecast{
	GeneratedRate: 1290.00,
	MarginalRate:  123.00,
	ConsumedRate:  123.88,
	Duration:      &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: startTime.Add(time.Minute * 4).Format(timeFormat)},
	DurationType:  reportdurations[0],
	Region:        testRegion,
}



func Test_pollersrvc_CarbonEmissions(t *testing.T) {
	downloadError := errors.New("download error")
	serverError := errors.New("server error")
	tests := []struct {
		name           string
		apiErr         error
		expectedOutput []*genpoller.CarbonForecast
	}{
		//test cases
		{name: "success", apiErr: nil, expectedOutput: mockReports},
		{name: "server error", apiErr: serverError, expectedOutput: nil},
		{name: "download error", apiErr: downloadError, expectedOutput: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carbonarac := carbonara.NewMock(t)
			mockReports = append(mockReports, mockReport)
			if tt.name == "success" {

			}
			if tt.apiErr != nil {
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					return nil, tt.apiErr

				})
			}
			carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
				return mockReports, nil
			})
			ctx := context.Background()
			var svc *pollersrvc
			svc = NewPoller(ctx, carbonarac, nil)
			svc.CarbonEmissions(ctx, testPeriod.StartTime, testRegion)
		})
	}
}

func TestGetDates(t *testing.T) {
	//Need to test this function because:
	//1. reports may be null
	//2. reports need to span the length of the total report duration
	dateErr := errors.New("incorrect date error")
	nilReportsErr := errors.New("no reports error")
	
	tests := []struct {
		name           string
		datesError     error
		nilReportsErr  error
		ctx            context.Context
		expectedOutput []*genpoller.Period
	}{ 
		{name: "success", datesError: nil},
		{name: "date Error", datesError: dateErr},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReports = append(mockReports, mockMinuteReportOne)
			got, err := GetDates(tt.ctx, mockReports)
			if got == nil {
				t.Errorf("GetDates() error = %v, wantErr %v", tt.nilReportsErr, nilReportsErr)
			}
			if err != nil {
				t.Errorf("GetDates() error = %v, wantErr %v", tt.datesError, dateErr)
				return
			}
			//check length of returned reports instead
			//not sure what this is
			if len(got) != 4 {

			}
		})
	}
}
func Test_pollersrvc_Update(t *testing.T) {
	downloadError := errors.New("download error")
	serverError := errors.New("server error")
	datesError := errors.New("download error")
	clickhouseError := errors.New("server error")
	tests := []struct {
		name            string
		downloadError   error
		serverError     error
		datesError      error
		clickhouseError error
		expectedError   error
	}{
		{name: "no error", downloadError: nil, serverError: nil, datesError: nil, clickhouseError: nil, expectedError: nil},
		{name: "server error", downloadError: nil, serverError: serverError, datesError: nil, clickhouseError: nil, expectedError: serverError},
		{name: "download error", downloadError: downloadError, serverError: nil, datesError: nil, clickhouseError: nil, expectedError: downloadError},
		{name: "dates error", downloadError: downloadError, serverError: serverError, datesError: datesError, clickhouseError: nil, expectedError: datesError},
		{name: "clickhouse error", downloadError: nil, serverError: nil, datesError: nil, clickhouseError: clickhouseError, expectedError: clickhouseError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carbonarac := carbonara.NewMock(t)
			if downloadError != nil {
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					return nil, downloadError
				})
			}
			if serverError != nil {
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					return nil, serverError
				})
			}
			if datesError != nil {
				if downloadError != nil {
					carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
						return nil, datesError
					})
				}
			}
			stc := storage.NewMock(t)
			if clickhouseError != nil {
				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					return nil, clickhouseError
				})
				stc.AddSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					return clickhouseError
				})
			}

			for i := 0; i < 13; i++ {
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					return mockReports, nil

				})
			}

			ctx := context.Background()
			svc := NewPoller(ctx, carbonarac, stc)
			err := svc.Update(ctx)
			if err != nil {
				if err != tt.expectedError {
					t.Errorf("pollersrvc.Update() error = %v, wantErr %v", err, tt.expectedError)
				}
			}

		})
	}
}


