// Code generated by mockery v2.14.0. DO NOT EDIT.

package storage

import (
	"context"
	"testing"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	mock "goa.design/clue/mock"
)

type(
	//CheckDBFunc mocks the storage Client checkDB method
	CheckDBFunc func(context.Context, string) (string, error)
	//SaveCarbonReportsFunc mocks the storage Client SaveCarbonReports method
	SaveCarbonReportsFunc func(context.Context, []*genpoller.CarbonForecast) (error)
	// PingFunc mocks the storage client Ping method.
	PingFunc func(context.Context) error
	//GetAggregateReportsFunc mocks the storage client GetAggregateReports method
	GetAggregateReportsFunc func(context.Context, []*genpoller.Period, string, string) ([]*genpoller.CarbonForecast, error)
	
	Mock struct {
		m *mock.Mock
		t *testing.T
	}
)
var _ Client = (*Mock)(nil)
func NewMock(t *testing.T) *Mock {
	return &Mock{mock.New(), t}
}

// AddCheckDBFunc adds a AddCheckDBFunc to the mock sequence.
func (m *Mock) AddCheckDBFunc(f CheckDBFunc) { m.m.Add("CheckDB", f) }
// SetCheckDBFunc adds a SetCheckDBFunc to the mock sequence.
func (m *Mock) SetCheckDBFunc(f CheckDBFunc) { m.m.Set("CheckDB", f) }
// GetAggregateReportsFunc  adds a GetAggregateReportsFunc to the mock sequence.
func (m *Mock) AddGetAggregateReportsFunc(f GetAggregateReportsFunc) { m.m.Add("GetAggregateReports", f) }
// SetAggregateReportsFunc  adds a SetAggregateReportsFunc  to the mock sequence.
func (m *Mock) SetGetAggregateReportsFunc(f GetAggregateReportsFunc) { m.m.Set("GetAggregateReports", f) }

//AddPingFunc adds a AddPingFunc to the mock sequence.
func (m *Mock) AddPingFunc(f PingFunc) { m.m.Add("Ping", f) }
// setPingFunc adds a setPingFunc to the mock sequence.
func (m *Mock) SetPingFunc(f PingFunc) { m.m.Set("Ping", f) }
// GetAggregateReportsFunc  adds a GetAggregateReportsFunc to the mock sequence.
func (m *Mock) AddSaveCarbonReportsFunc(f SaveCarbonReportsFunc) { m.m.Add("SaveCarbonReports", f) }
// SetAggregateReportsFunc  adds a SetAggregateReportsFunc  to the mock sequence.
func (m *Mock) SetSaveCarbonReportsFunc(f SaveCarbonReportsFunc) { m.m.Set("SaveCarbonReports", f) }


// CheckDB provides a mock function with given fields: context, region
func (m *Mock) CheckDB(ctx context.Context, region string) (string, error) {
	if f := m.m.Next("CheckDB"); f != nil {
		return f.(CheckDBFunc)(ctx, region)
	}
	m.t.Error("unexpected CheckDB call")
	return "", nil
}

// GetAggregateReportsFunc provides a mock function with given fields: context, dates, region ,duration
func (m *Mock) GetAggregateReports(ctx context.Context, periods []*genpoller.Period, region string, duration string) ([]*genpoller.CarbonForecast, error) {
	if f := m.m.Next("GetAggregateReports"); f != nil {
		return f.(GetAggregateReportsFunc)(ctx, periods, region, duration)
	}
	m.t.Error("unexpected GetAggregateReports call")
	return nil, nil
}

// Ping provides a mock function with given fields: ctx
func (m *Mock) Ping(ctx context.Context) error {
	if f := m.m.Next("Ping"); f != nil {
		return f.(PingFunc)(ctx)
	}
	m.t.Error("unexpected Ping call")
	return nil
}

// SaveCarbonReports provides a mock function with given fields: context, reports
func (m *Mock) SaveCarbonReports(ctx context.Context, reports []*genpoller.CarbonForecast) error {
	if f := m.m.Next("SaveCarbonReports"); f != nil {
		return f.(SaveCarbonReportsFunc)(ctx, reports)
	}
	m.t.Error("unexpected SaveCarbonReports call")
	return nil
}
func (m *Mock) Name() string               { return "mock" }
func (m *Mock) Init(context.Context, bool) error { return nil }

