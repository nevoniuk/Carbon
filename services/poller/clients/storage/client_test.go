package storage

import (
	"context"
	"reflect"
	"testing"

	"github.com/crossnokaye/carbon/clients/clickhouse"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

func Test_client_SaveCarbonReports(t *testing.T) {
	/*
	errors:
	IncorrectReportsError
	no reports error or null
	*/
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

func Test_client_GetAggregateReports(t *testing.T) {
	/* 
		errors:
		no reports error
	
	*/
	type fields struct {
		chcon clickhouse.Conn
	}
	type args struct {
		ctx      context.Context
		periods  []*genpoller.Period
		region   string
		duration string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*genpoller.CarbonForecast
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				chcon: tt.fields.chcon,
			}
			got, err := c.GetAggregateReports(tt.args.ctx, tt.args.periods, tt.args.region, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("client.GetAggregateReports() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("client.GetAggregateReports() = %v, want %v", got, tt.want)
			}
		})
	}
}
