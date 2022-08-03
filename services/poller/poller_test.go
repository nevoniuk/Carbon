package pollerapi

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/crossnokaye/carbon/model"
	"github.com/crossnokaye/carbon/services/poller/clients/carbonara"
	"github.com/crossnokaye/carbon/services/poller/clients/storage"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

var testRegion = regions[0]
func TestGetDates(t *testing.T) {
	var mockReports []*genpoller.CarbonForecast
	var mockDates []*genpoller.Period
	dateErr := errors.New("incorrect date error")
	nilReportsErr := errors.New("no reports error")
	var startTimeOne = time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
	var endTimeOne = time.Date(2021, time.June, 1, 1, 10, 0, 0, time.UTC)
	var startTimeTwo = time.Date(2021, time.June, 1, 1, 10, 0, 0, time.UTC)
	var endTimeTwo = time.Date(2021, time.June, 1, 1, 20, 0, 0, time.UTC)
	var invalidEndTime = time.Date(2021, time.February, 1, 1, 20, 0, 0, time.UTC)

	var mockReportOne = &genpoller.CarbonForecast{
		Duration: &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: endTimeOne.Format(timeFormat)},
	}
	var mockReportTwo = &genpoller.CarbonForecast{
		Duration: &genpoller.Period{StartTime: startTimeTwo.Format(timeFormat), EndTime: endTimeTwo.Format(timeFormat)},
	}
	var invalidReport = &genpoller.CarbonForecast{
		Duration: &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: invalidEndTime.Format(timeFormat)},
	}
	mockRes := &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: endTimeOne.Format(timeFormat)}
	mockDates = append(mockDates, mockRes)
	tests := []struct {
		name           string
		datesError     error
		nilReportsErr  error
		expectedOutput []*genpoller.Period
		expectedError  error
	}{
		{name: "success", datesError: nil, nilReportsErr: nil, expectedOutput: mockDates, expectedError: nil},
		{name: "invalid date", datesError: dateErr, nilReportsErr: nil, expectedOutput: mockDates, expectedError: dateErr}, //no actual error thrown
		{name: "nil Error", datesError: nil, nilReportsErr: nilReportsErr, expectedOutput: nil, expectedError: nilReportsErr},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockReports = nil
			if tt.datesError != nil {
				mockReports = append(mockReports, invalidReport)
				mockReports = append(mockReports, mockReportOne)
			} else if tt.nilReportsErr!= nil {
				mockReports = nil
			} else {
				mockReports = append(mockReports, mockReportOne)
				mockReports = append(mockReports, mockReportTwo)
			}
			ctx := context.Background()
			got, err := getDates(ctx, mockReports, model.Hourly)
			valid := validateDates(got, mockDates)
			if tt.nilReportsErr != nil && !valid {
				t.Errorf("GetDates() error = %v, wantErr %v", err, tt.expectedError)
				return
			} else if tt.datesError != nil && !valid {
				t.Errorf("GetDates() error = %v, wantErr %v", tt.expectedError, dateErr)
			} else {
				fmt.Printf("success for test %d\n", i)
			}
		})
	}
}

func validateDates(dates []*genpoller.Period, mockDates []*genpoller.Period) bool {
	if dates == nil {
		return true
	}
	counter := 0
	for _, rep := range dates {
		if (rep.StartTime != mockDates[counter].StartTime) || (rep.EndTime != mockDates[counter].EndTime) {
			return false
		}
		counter += 1
	}
	return true
}

func Test_pollersrvc_Update(t *testing.T) {
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
	var endTime = time.Date(2021, time.June, 1, 1, 0, 0, 0, time.UTC)
	var invalidEndTime = time.Date(2021, time.February, 1, 1, 20, 0, 0, time.UTC)
	var mockReportOne = &genpoller.CarbonForecast{
		Duration:     &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: endTime.Format(timeFormat)},
		DurationType: reportdurations[1],
	}
	var mockReports []*genpoller.CarbonForecast
	serverError := carbonara.ServerError{Err: fmt.Errorf("server error 400")}
	noDataError := carbonara.NoDataError{Err: fmt.Errorf("no data for Region")}
	saveReportsError := storage.NoReportsError{Err: fmt.Errorf("no reports\n")}
	getReportsError := storage.NoReportsError{Err: fmt.Errorf("no reports\n")}

	tests := []struct {
		name             string
		serverError      error
		noDataError 	 error
		saveReportsError error
		getReportsError  error
		expectedError    error
	}{
		{name: "success", serverError: nil, noDataError: nil, saveReportsError: nil, getReportsError: nil, expectedError: nil},
		{name: "server error", serverError: serverError, noDataError: nil, saveReportsError: nil, getReportsError: nil, expectedError: serverError},
		{name: "error saving reports", serverError: nil, noDataError: nil, saveReportsError: saveReportsError, getReportsError: nil, expectedError: saveReportsError},
		{name: "no data error", serverError: nil, noDataError: noDataError, saveReportsError: nil, getReportsError: nil, expectedError: noDataError},
		{name: "error getting reports from clickhouse", serverError: nil, noDataError: nil, saveReportsError: nil, getReportsError: getReportsError, expectedError: getReportsError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carbonarac := carbonara.NewMock(t)
			stc := storage.NewMock(t)
			mockReports = append(mockReports, mockReportOne)
			ctx := context.Background()
			stc.SetCheckDBFunc(func(ctx context.Context, r string) (string, error) {
				return startTime.Format(timeFormat), nil
			})
			if tt.expectedError == nil {
				carbonarac.SetGetEmissionsFunc(func(ctx context.Context, r string, s string, e string) ([]*genpoller.CarbonForecast, error) {
					return mockReports, nil
				})
				
				stc.SetSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					reps = mockReports
					return nil
				})
				stc.SetGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					return mockReports, nil
				})
			}

			if tt.serverError != nil {
				mockReports = nil
				carbonarac.SetGetEmissionsFunc(func(ctx context.Context, r string, s string, e string) ([]*genpoller.CarbonForecast, error) {
					e = invalidEndTime.Format(timeFormat)
					return nil, serverError
				})
			}
			if tt.noDataError != nil {
					carbonarac.SetGetEmissionsFunc(func(ctx context.Context, r string, s string, e string) ([]*genpoller.CarbonForecast, error) {
						return nil, noDataError
					})
			}
			if tt.saveReportsError != nil {
				carbonarac.SetGetEmissionsFunc(func(ctx context.Context, r string, s string, e string) ([]*genpoller.CarbonForecast, error) {
					return mockReports, nil
				})
				
				stc.SetSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					reps = nil
					return saveReportsError
				})
			}
			if tt.getReportsError != nil {
				carbonarac.SetGetEmissionsFunc(func(ctx context.Context, r string, s string, e string) ([]*genpoller.CarbonForecast, error) {
					return mockReports, nil
				})
				
				stc.SetSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					reps = mockReports
					return nil
				})
				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					var newdates []*genpoller.Period
					dates = nil
					newdates = append(newdates, &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: invalidEndTime.Format(timeFormat)})
					dates = newdates
					return nil, getReportsError
				})
			}
			
			svc := NewPoller(ctx, carbonarac, stc)
			svc.now = endTime
			err := svc.Update(ctx)
			mockReports = nil
			if err != nil {
				if !strings.Contains(err.Error(), tt.expectedError.Error()){
					t.Errorf("pollersrvc.Update() error = %v, wantErr %v", err, tt.expectedError)
				}
			} else {
				fmt.Println("success")
			}
			svc.now = timeNow()
		})
	}
}

func Test_pollersrvc_AggregateData(t *testing.T) {
	nilReportsError := errors.New("no reports error")
	invalidReportsError := errors.New("incorrect reports error")
	invalidTimesError := errors.New("invalid time periods error")
	var mockReports []*genpoller.CarbonForecast
	var mockDates []*genpoller.Period
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, time.UTC)
	var endTime = time.Date(2021, time.June, 1, 0, 10, 0, 0, time.UTC)
	var invalidEndTime = time.Date(2021, time.February, 1, 0, 10, 0, 0, time.UTC)
	var invalidTestPeriod = &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: invalidEndTime.Format(timeFormat)}
	var testPeriod = &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: endTime.Format(timeFormat)}
	invalidReport := &genpoller.CarbonForecast{GeneratedRate: 0, ConsumedRate: 0, MarginalRate: 0, Duration: invalidTestPeriod, DurationType: reportdurations[0], Region: testRegion}
	validReport := &genpoller.CarbonForecast{GeneratedRate: 0, ConsumedRate: 0, MarginalRate: 0, Duration: testPeriod, DurationType: reportdurations[0], Region: testRegion}
	tests := []struct {
		name                string
		nilreportsError     error
		invalidTimesError   error
		invalidReportsError error
		expectedError       error
	}{
		{name: "success", nilreportsError: nil, invalidTimesError: nil, invalidReportsError: nil, expectedError: nil},
		{name: "invalid times error", nilreportsError: nil, invalidTimesError: invalidTimesError, invalidReportsError: nil, expectedError: invalidTimesError},
		{name: "null reports error", nilreportsError: nilReportsError, invalidTimesError: nil, invalidReportsError: nil, expectedError: nilReportsError},
		{name: "invalid reports error", nilreportsError: nil, invalidTimesError: nil, invalidReportsError: invalidReportsError, expectedError: invalidReportsError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stc := storage.NewMock(t)
			mockDates = append(mockDates, testPeriod)
			mockReports = append(mockReports, validReport)

			if tt.invalidTimesError != nil {
				mockDates = nil
				mockDates = append(mockDates, invalidTestPeriod)
				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					dates = mockDates
					return nil, invalidTimesError
				})

			}
			if tt.invalidReportsError != nil {
				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					return mockReports, nil
				})
				mockReports = nil
				mockReports = append(mockReports, invalidReport)
				stc.AddSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					reps = mockReports
					return invalidReportsError
				})
			}
			if tt.nilreportsError != nil {
				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					return mockReports, nil
				})
				mockReports = nil
				stc.AddSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					reps = nil
					return nilReportsError
				})
			}
			if tt.expectedError == nil {
				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					return mockReports, nil
				})
				stc.AddSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					return nil
				})
			}
			ctx := context.Background()
			svc := NewPoller(ctx, nil, stc)
			res, err := svc.aggregateData(ctx, testRegion, mockDates, reportdurations[0])
			if err != nil && res == nil {
				if !strings.Contains(err.Error(), tt.expectedError.Error()) {
					t.Errorf("pollersrvc.AggregateData() error = %v, wantErr %v", err, tt.expectedError)
				}
			}
			mockDates = nil
			mockReports = nil
		})
	}
}

