package storage

import (
	"context"
	"testing"

	"github.com/crossnokaye/carbon/clients/clickhouse"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

func Test_client_SaveCarbonReports(t *testing.T) {
	type fields struct {
		chcon clickhouse.Conn
	}
	type args struct {
		ctx     context.Context
		reports []*genpoller.CarbonForecast
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				chcon: tt.fields.chcon,
			}
			if err := c.SaveCarbonReports(tt.args.ctx, tt.args.reports); (err != nil) != tt.wantErr {
				t.Errorf("client.SaveCarbonReports() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
