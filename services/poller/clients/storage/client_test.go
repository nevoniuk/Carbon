package storage

import (
	"context"
	"reflect"
	"testing"
	"time"
	"errors"
	"os"
	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/crossnokaye/carbon/clients/clickhouse"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)
var ctx = context.Background()

func setupClickhouse(t *testing.T) *client {
	retried := false
try:
	con, err := ch.Open(&ch.Options{
		Addr: []string{"localhost:8088"},
		Auth: ch.Auth{
			Username: os.Getenv("CLICKHOUSE_USER"),
			Password: os.Getenv("CLICKHOUSE_PASSWORD"),
		}})
	if err != nil {
		if retried {
			t.Fatalf("could not connect to clickhouse: %v", err)
		}
		// Give docker-compose a few seconds to come up for tests
		time.Sleep(10 * time.Second)
		retried = true
		goto try
	}
	return New(clickhouse.New(con)).(*client)
}

func TestInit(t *testing.T) {
	c := setupClickhouse(t)
	//defer cleanupClickhouse(t, c)
	ctx := context.Background()

	if err := c.Init(ctx, true); err != nil {
		t.Errorf("could not initialize clickhouse: %v", err)
		return
	}
	rows, err := c.chcon.Query(ctx, "DESCRIBE carbondb.carbon_reports;")
	if err != nil {
		t.Errorf("could not describe carbon intensity table: %v", err)
	}
	defer rows.Close()

	if !rows.Next() {
		t.Errorf("could not read start column")
	}
	var n, type_, def, defexpr, co, coexpr, ttl string
	if err := rows.Scan(&n, &type_, &def, &defexpr, &co, &coexpr, &ttl); err != nil {
		t.Errorf("could not scan start column: %v", err)
	}
	if n != "start" || type_ != "DateTime" || def != "" || defexpr != "" || co != "" || coexpr != "" || ttl != "" {
		t.Errorf("unexpected start column details: %s %s %s %s %s %s %s", n, type_, def, defexpr, co, coexpr, ttl)
	}
	if !rows.Next() {
		t.Errorf("could not read end column")
	}
	if err := rows.Scan(&n, &type_, &def, &defexpr, &co, &coexpr, &ttl); err != nil {
		t.Errorf("could not scan end column: %v", err)
	}
	if n != "end" || type_ != "DateTime" || def != "" || defexpr != "" || co != "" || coexpr != "" || ttl != "" {
		t.Errorf("unexpected end column details: %s %s %s %s %s %s %s", n, type_, def, defexpr, co, coexpr, ttl)
	}
	if !rows.Next() {
		t.Errorf("could not read generatedrate column")
	}
	if err := rows.Scan(&n, &type_, &def, &defexpr, &co, &coexpr, &ttl); err != nil {
		t.Errorf("could not scan generatedrate column: %v", err)
	}
	if n != "generatedrate" || type_ != "Float64" || def != "" || defexpr != "" || co != "" || coexpr != "" || ttl != "" {
		t.Errorf("unexpected generatedrate column details: %s %s %s %s %s %s %s", n, type_, def, defexpr, co, coexpr, ttl)
	}
	if !rows.Next() {
		t.Errorf("could not read marginalrate column")
	}
	if err := rows.Scan(&n, &type_, &def, &defexpr, &co, &coexpr, &ttl); err != nil {
		t.Errorf("could not scan marginalrate column: %v", err)
	}
	if n != "marginalrate" || type_ != "Float64" || def != "" || defexpr != "" || co != "" || coexpr != "" || ttl != "" {
		t.Errorf("unexpected marginalrate column details: %s %s %s %s %s %s %s", n, type_, def, defexpr, co, coexpr, ttl)
	}
	if !rows.Next() {
		t.Errorf("could not read consumedrate column")
	}
	if err := rows.Scan(&n, &type_, &def, &defexpr, &co, &coexpr, &ttl); err != nil {
		t.Errorf("could not scan consumedrate column: %v", err)
	}
	if n != "consumedrate" || type_ != "Float64" || def != "" || defexpr != "" || co != "" || coexpr != "" || ttl != "" {
		t.Errorf("unexpected consumedrate column details: %s %s %s %s %s %s %s", n, type_, def, defexpr, co, coexpr, ttl)
	}

	if !rows.Next() {
		t.Errorf("could not read region column")
	}
	if err := rows.Scan(&n, &type_, &def, &defexpr, &co, &coexpr, &ttl); err != nil {
		t.Errorf("could not scan region column: %v", err)
	}
	if n != "region" || type_ != "String" || def != "" || defexpr != "" || co != "" || coexpr != "" || ttl != "" {
		t.Errorf("unexpected region column details: %s %s %s %s %s %s %s", n, type_, def, defexpr, co, coexpr, ttl)
	}

	if !rows.Next() {
		t.Errorf("could not read duration column")
	}
	if err := rows.Scan(&n, &type_, &def, &defexpr, &co, &coexpr, &ttl); err != nil {
		t.Errorf("could not scan duration column: %v", err)
	}
	if n != "duration" || type_ != "String" || def != "" || defexpr != "" || co != "" || coexpr != "" || ttl != "" {
		t.Errorf("unexpected duration column details: %s %s %s %s %s %s %s", n, type_, def, defexpr, co, coexpr, ttl)
	}

	if rows.Next() {
		t.Errorf("found unexpected additional column")
	}

	if err := c.Init(ctx, true); err != nil {
		t.Errorf("got error initializing twice: %v", err)
	}
}

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
	var startTime = time.Date(2021, time.June, 1, 0, 0, 0, 0, nil)
	var invalidStartTime = time.Date(2021, time.June, 1, 1, 10, 0, 0, nil)
	var endTime = time.Date(2021, time.June, 1, 1, 10, 0, 0, nil)
	var invalidEndTime = time.Date(2021, time.February, 1, 1, 20, 0, 0, nil)
	valid := genpoller.Period{StartTime: startTime.Format(timeFormat), EndTime: endTime.Format(timeFormat)}
	invalid := genpoller.Period{StartTime: invalidStartTime.Format(timeFormat), EndTime: invalidEndTime.Format(timeFormat)}
	var periods []*genpoller.Period
	nilReportsErr := errors.New("no reports error")
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
		noReportsErrors  error
		expectedError error
	}{
		{name: "success", noReportsErrors: nil, expectedError: nil},
		{name: "success", noReportsErrors: nilReportsErr, expectedError: nilReportsErr},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := New()
			c := &client{
				chcon: tt.fields.chcon,
			}
			dbc := storage.NewMock(c.chcon)
			if tt.noReportsErrors != nil {
				
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
