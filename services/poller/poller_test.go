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

var testRegion = regions[0]

func Test_pollersrvc_CarbonEmissions(t *testing.T) {
	//testing 10 minute period
	downloadError := errors.New("download error")
	serverError := errors.New("server error")
	var mockReports []*genpoller.CarbonForecast
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, nil)
	var endTime = time.Date(2021, time.June, 1, 0, 10, 0, 0, nil)
	var fiveamt = time.Minute * 5
	var mockReportOne = &genpoller.CarbonForecast{
		Duration:     &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: endTime.Add(-fiveamt).Format(timeFormat)},
		DurationType: reportdurations[0],
	}
	var mockReportTwo = &genpoller.CarbonForecast{
		Duration:     &genpoller.Period{StartTime: startTime.Add(fiveamt).Format(timeFormat), EndTime: endTime.Format(timeFormat)},
		DurationType: reportdurations[0],
		Region:       testRegion,
	}
	mockReports = append(mockReports, mockReportOne)
	mockReports = append(mockReports, mockReportTwo)
	tests := []struct {
		name           string
		apiErr         error
		serverErr      error
		expectedOutput []*genpoller.CarbonForecast
	}{
		{name: "success", apiErr: nil, serverErr: nil, expectedOutput: mockReports},
		{name: "server error", apiErr: nil, serverErr: serverError, expectedOutput: nil},
		{name: "download error", apiErr: downloadError, serverErr: nil, expectedOutput: nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			carbonarac := carbonara.NewMock(t)
			if tt.apiErr != nil {
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					return nil, tt.apiErr
				})
			}
			if tt.serverErr != nil {
				endTime = time.Date(2000, time.June, 1, 0, 0, 0, 0, nil)
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					return nil, tt.serverErr
				})
			}
			carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
				return mockReports, nil
			})

			ctx := context.Background()
			var svc *pollersrvc
			svc = NewPoller(ctx, carbonarac, nil)
			res, err := svc.CarbonEmissions(ctx, startTime.Format(timeFormat), endTime.Format(timeFormat), testRegion)
			if tt.expectedOutput != nil {
				if err == nil {
					//validate output here

					t.Errorf("carbonEmissions did not return an error")
				} else if errors.As(err, tt.apiErr.Error()) {
					t.Errorf("carbonEmissions returned error: %v, want %v", err, tt.apiErr.Error())
				}
			} else { //no expected output
				if errors.As(err, tt.apiErr.Error()) {
					t.Errorf("carbonEmissions returned error: %v, want %v", err, tt.apiErr.Error())
				}
			}
		})
	}
}
//BUG: does get dates assume dates are sequential
func TestGetDates(t *testing.T) {
	var mockReports []*genpoller.CarbonForecast
	var mockReps []*genpoller.Period
	dateErr := errors.New("incorrect date error")
	nilReportsErr := errors.New("no reports error")
	var startTimeOne = time.Date(2021, time.June, 1, 0, 0, 0, 0, nil)
	var startTimeTwo = time.Date(2021, time.June, 1, 1, 10, 0, 0, nil)
	var endTimeOne = time.Date(2021, time.June, 1, 1, 10, 0, 0, nil)
	var endTimeTwo = time.Date(2021, time.June, 1, 1, 20, 0, 0, nil)
	var invalidEndTime = time.Date(2021, time.February, 1, 1, 20, 0, 0, nil)
	var mockReportOne = &genpoller.CarbonForecast {
		Duration: &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: endTimeOne.Format(timeFormat)},
	}
	var mockReportTwo = &genpoller.CarbonForecast {
		Duration: &genpoller.Period{StartTime: startTimeTwo.Format(timeFormat), EndTime: endTimeTwo.Format(timeFormat)},
	}
	var invalidReport = &genpoller.CarbonForecast {
		Duration: &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: invalidEndTime.Format(timeFormat)},
	}
	mockRes := &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: startTimeOne.Format(timeFormat)}
	mockReps = append(mockReps, mockRes)
	tests := []struct {
		name           string
		datesError     error
		nilReportsErr  error
		expectedOutput []*genpoller.Period
		expectedError  error
	}{
		{name: "success", datesError: nil, nilReportsErr: nil, expectedOutput: mockReps, expectedError: nil},
		{name: "date Error", datesError: dateErr, nilReportsErr: nil, expectedOutput: nil, expectedError: dateErr},
		{name: "nil Error", datesError: nil, nilReportsErr: nilReportsErr, expectedOutput: nil, expectedError: nilReportsErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if dateErr != nil {
				mockReports = append(mockReports, mockReportOne)
				mockReports = append(mockReports, invalidReport)
			} else if nilReportsErr != nil {
				mockReports = nil
			} else {
				mockReports = append(mockReports, mockReportOne)
				mockReports = append(mockReports, mockReportTwo)
			}
			ctx := context.Background()
			got, err := GetDates(ctx, mockReports)
			if !errors.As(err, tt.expectedError) {
				t.Errorf("GetDates() error = %v, wantErr %v", err, tt.expectedError)
				return
			}
			if len(got) != 1 {
				t.Errorf("GetDates generated redundant reports")
				return
			}
			for i, rep := range got[0] { //0 because only hour reports were returned
				if (rep.StartTime != mockReps[i].StartTime) || (rep.EndTime != mockReps[i].EndTime) {
					t.Errorf("Result does not equal mock")
				}
			}
		})
	}
}
func Test_pollersrvc_Update(t *testing.T) {
	//this function won't actually call Update() because its runs on an unknown input(i.e time since last download)
	//testCase start and end times:
	var startTimeOne = time.Date(2021, time.June, 1, 0, 0, 0, 0, nil)
	var endTimeTwo = time.Date(2021, time.June, 1, 1, 20, 0, 0, nil)
	//this will return 16 5 min reports
	//check length after get emissions
	var mockReports []*genpoller.CarbonForecast
	var mockReps []*genpoller.Period

	var startTimeOne = time.Date(2021, time.June, 1, 0, 0, 0, 0, nil)
	var startTimeTwo = time.Date(2021, time.June, 1, 1, 10, 0, 0, nil)
	var endTimeOne = time.Date(2021, time.June, 1, 1, 10, 0, 0, nil)
	
	var invalidEndTime = time.Date(2021, time.February, 1, 1, 20, 0, 0, nil)

	var mockReportOne = &genpoller.CarbonForecast {
		Duration: &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: endTimeOne.Format(timeFormat)},
	}
	var mockReportTwo = &genpoller.CarbonForecast {
		Duration: &genpoller.Period{StartTime: startTimeTwo.Format(timeFormat), EndTime: endTimeTwo.Format(timeFormat)},
	}
	var invalidReport = &genpoller.CarbonForecast {
		Duration: &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: invalidEndTime.Format(timeFormat)},
	}
	mockRes := &genpoller.Period{StartTime: startTimeOne.Format(timeFormat), EndTime: startTimeOne.Format(timeFormat)}
	mockReps = append(mockReps, mockRes)
	downloadError := errors.New("download error")
	serverError := errors.New("server error")
	clickhouseError := errors.New("server error")
	
	tests := []struct {
		name            string
		downloadError   error
		serverError     error
		datesError      error
		clickhouseError error
		nilReportsError error
		expectedError   error
	}{
		{name: "no error", downloadError: nil, serverError: nil, datesError: nil, clickhouseError: nil, expectedError: nil},
		{name: "server error", downloadError: nil, serverError: serverError, datesError: nil, clickhouseError: nil, expectedError: serverError},
		{name: "download error", downloadError: downloadError, serverError: nil, datesError: nil, clickhouseError: nil, expectedError: downloadError},
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
				endTimeTwo = invalidEndTime
				carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
					return nil, serverError
				})
			}

			stc := storage.NewMock(t)
			if clickhouseError != nil {
				//throw an error in client if bad reports are about to be saved
				//give bad input to save reports and give bad dates for save carbon reports

				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					return nil, clickhouseError
				})
				stc.AddSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					return clickhouseError
				})
			}
			//called once
			carbonarac.AddGetEmissionsFunc(func(ctx context.Context, r string, s string, e string, reps []*genpoller.CarbonForecast) ([]*genpoller.CarbonForecast, error) {
				return mockMinuteReports, nil
			})
			stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
				return nil, nil
			})
			stc.AddSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
				return nil
			})
			ctx := context.Background()
			svc := NewPoller(ctx, carbonarac, stc)
			for i := 0; i < 13; i++ {
				minreps, err := svc.CarbonEmissions(ctx, startTimeOne.Format(timeFormat), endTimeTwo.Format(timeFormat), testRegion)
				if !errors.As(err, tt.expectedError) {

				} else if err != nil && len(minreps) != 16 {
					//not enough reports returned
				}

				dateConfigs, datesErr := GetDates(ctx, minreps)
				if !errors.As(datesErr, tt.expectedError) {

				} else if err != nil {
					//check date configs
				}
				//minreps
				err = svc.dbc.SaveCarbonReports(ctx, minreps)
				if !errors.As(err, tt.expectedError) {

				}
			for j := 0; j < len(dateConfigs); j++ {
				if dateConfigs[j] != nil {
					fmt.Printf("j is %d\n", j)
					aggErr := s.AggregateData(ctx, regions[i], dateConfigs[j], reportdurations[j])
					if aggErr != nil {
						fmt.Errorf("Error from aggregate data %s\n", aggErr)
						return aggErr
					}
				}
			}
				
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

func Test_pollersrvc_AggregateData(t *testing.T) {
	nilReportsError := errors.New("no reports error")
	clickhouseError := errors.New("clickhouse error")
	var mockReports []*genpoller.CarbonForecast
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, nil)
	var endTime = time.Date(2021, time.June, 1, 0, 10, 0, 0, nil)
	var invalidEndTime = time.Date(2021, time.February, 1, 0, 10, 0, 0, nil)
	var invalidTestPeriod = &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: invalidEndTime.Format(timeFormat)}
	var testPeriod = &genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: endTime.Format(timeFormat)}
	tests := []struct {
		name    string
		nilreportsError error
		clickhouseError error
		expectedError error
	}{
		{name: "success", nilreportsError: nil, clickhouseError: nil, expectedError: nil},
		{name: "no reports error", nilreportsError: nilReportsError, clickhouseError: nil, expectedError: nil},
		{name: "clickhouse error", nilreportsError: nil, clickhouseError: clickhouseError, expectedError: clickhouseError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stc := storage.NewMock(t)
			if clickhouseError != nil {
				testPeriod = invalidTestPeriod
				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					return nil, clickhouseError
				})
				stc.AddSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					return clickhouseError
				})
			}
			if nilReportsError != nil {
				stc.AddGetAggregateReportsFunc(func(ctx context.Context, dates []*genpoller.Period, r string, duration string) ([]*genpoller.CarbonForecast, error) {
					return nil, nilReportsError
				})
				stc.AddSaveCarbonReportsFunc(func(ctx context.Context, reps []*genpoller.CarbonForecast) error {
					return nilReportsError
				})
			}
			var mockReport = &genpoller.CarbonForecast{
				GeneratedRate: 1234.00,
				ConsumedRate: 12321.00,
				MarginalRate: 1232.00,
				Duration: testPeriod,
				DurationType: reportdurations[0],
				Region: testRegion,
			}
			ctx := context.Background()
			svc := NewPoller(ctx, nil, stc)
			var mockPeriods []*genpoller.Period
			mockPeriods = append(mockPeriods, testPeriod)
			mockReports = append(mockReports, mockReport)
			err := svc.AggregateData(ctx, testRegion, mockPeriods, reportdurations[0])
			 if err != nil {
				if err != tt.expectedError {
					t.Errorf("pollersrvc.AggregateData() error = %v, wantErr %v", err, tt.expectedError)
				}
			 }
			 
		})
	}
}
