package carbonara

import (
	"context"
	"testing"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	mock "goa.design/clue/mock"
)

type(

	GetEmissionsFunc func(context.Context, string, string, string) ([]*genpoller.CarbonForecast, error)
	
	Mock struct {
		m *mock.Mock
		t *testing.T
	}
)
func NewMock(t *testing.T) *Mock {
	return &Mock{mock.New(), t}
}

// AddGetEmissionsFunc adds a GetEmissionsFunc to the mock sequence.
func (m *Mock) AddGetEmissionsFunc(f GetEmissionsFunc) { m.m.Add("GetEmissions", f) }
// SetGetEmissionsFunc sets a permanent GetEmissions mock implementation.
func (m *Mock) SetGetEmissionsFunc(f GetEmissionsFunc) { m.m.Set("GetEmissions", f) }

// GetEmissions provides a mock function with given fields: _a0, _a1, _a2, _a3, _a4
func (m *Mock) GetEmissions(ctx context.Context, region string, start string, end string) ([]*genpoller.CarbonForecast, error) {
	if f := m.m.Next("GetEmissions"); f != nil {
		return f.(GetEmissionsFunc)(ctx, region, start, end)
	}
	m.t.Error("unexpected GetEmissions call")
	return nil, nil
}
