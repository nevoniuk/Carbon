package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/crossnokaye/carbon/clients/clickhouse"
	"github.com/crossnokaye/carbon/model"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	"goa.design/clue/log"
)

type (
	Client interface {
		// Name returns the Name of the DB
		Name() string
		// Init initializes the clickhouse DB
		Init(context.Context, bool) error
		// Ping will ensure the DB connection is valid
		Ping(context.Context) error
		// GetCarbonReports will return carbon intensity reports from clickhouse. These reports were obtained from the poller service
		GetCarbonIntensityReports(context.Context, []*gencalc.Period, string, string) (*gencalc.CarbonReport, error)
	}
	client struct {
		chcon clickhouse.Conn
	}
	// errReportsNotFound is returned when carbon reports for given input are not found.
	ErrReportsNotFound struct{ Err error }
)
// timeFormat is used to parse times in order to store time as ISO8601 format
const timeFormat = "2006-01-02T15:04:05-07:00"
const sqlError = "sql: no rows in result set"
// reportdurations maintains the interval length of each report using constants from the model directory
var reportdurations [5]string = [5]string{model.Minute, model.Hourly, model.Daily, model.Weekly, model.Monthly}
func (c *client) Name() string {
	var name = "Clickhouse"
	return name
}

func New(chcon clickhouse.Conn) Client {
	return &client{chcon}
}

// ping is the time it takes for a small data set to be transmitted from the device to a server on the internet 
func (c *client) Ping(ctx context.Context) error {
	return c.chcon.Ping(ctx)
}
func (c *client) Init(ctx context.Context, test bool) error {
	if err := c.chcon.Ping(ctx); err != nil {
		if exception, ok := err.(*ch.Exception); ok {
			return ErrReportsNotFound{fmt.Errorf("[%d] %s", exception.Code, exception.Message)}
		}
		return ErrReportsNotFound{fmt.Errorf("database could not be found: %w", err)}
	}

	if err := c.chcon.Exec(ctx, `CREATE DATABASE IF NOT EXISTS carbondb;`); err != nil {
		return err
	}
	/**
	var err error 
	err = c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS carbondb.power_reports (
				start DateTime,
				end DateTime,
				generatedrate Float64,
				intervaltype String,

				) Engine =  MergeTree()
				ORDER BY (start)
	`) 
	
	if err != nil {
		return fmt.Errorf("failed to create power reports table")
	}
*/
	return nil
}

// GetCarbonReports will return carbon intensity reports from clickhouse for the given duration, region and intervaltype
func (c *client) GetCarbonIntensityReports(ctx context.Context, duration []*gencalc.Period, intervalType string, region string) (*gencalc.CarbonReport, error) {
	var intensityPoints []*gencalc.DataPoint
	var averagegen float64
	for _, period := range duration {
		var newstart, _ = time.Parse(timeFormat, period.StartTime)
		var newend, _ = time.Parse(timeFormat, period.EndTime)
		rows := c.chcon.QueryRow(ctx,`
		SELECT
			AVG(generatedrate) AS generatedrate
		FROM 
			carbondb.carbon_intensity_reports
		WHERE
			region = $1 AND start >= $2 AND end <= $3
		GROUP BY region
				`, region, newstart.UTC(), newend.UTC())
		err := rows.Scan(&averagegen)
		if err != nil {
			if strings.Contains(err.Error(), sqlError) {
				//continue because there may be holes in data where it was not available
				log.Errorf(ctx, err, "no data for start %s and end %s and region %s", newstart.UTC(), newend.UTC(), region)
				continue
			}
			return nil, ErrReportsNotFound{Err: fmt.Errorf("could not get carbon intensity report for start %s and end %s and region %s: %w", period.StartTime, period.EndTime, region, err)}
		}
		intensityPoint := &gencalc.DataPoint{Time: period.StartTime, Value: averagegen}
		fmt.Println("INTENSITY POINT")
		fmt.Println(intensityPoint)
		intensityPoints = append(intensityPoints, intensityPoint)	
	}
	var startTime = duration[0].StartTime
	var endTime = duration[(len(duration) - 1)].EndTime
	var per = &gencalc.Period{StartTime: startTime, EndTime: endTime}
	report := &gencalc.CarbonReport{IntensityPoints: intensityPoints, Duration: per, Interval: intervalType, Region: region}
	return report, nil
}
func (err ErrReportsNotFound) Error() string { return err.Err.Error() }

