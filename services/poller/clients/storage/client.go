package storage
import (
	"context"
	"fmt"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/crossnokaye/carbon/clients/clickhouse"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	"goa.design/clue/log"
)
type (
	Client interface {
		// returns the name of the DB
		Name() string
		// initialize clickhouse DB
		Init(context.Context, bool) error
		// Pings clickhouse
		Ping(ctx context.Context) error
		// Checks for the last report available in clickhouse for a given region
		CheckDB(context.Context, string) (string, error)
		// Saves carbon intensity reports in clickhouse
		SaveCarbonReports(context.Context, []*genpoller.CarbonForecast) (error)
	}
	client struct {
		chcon clickhouse.Conn
	}
	// NoReportsError is returned when no carbon intensity reports are given or available to find
	NoReportsError struct{ Err error }
	// IncorrectReportsError is returned when reports cannot be saved in the correct format
	IncorrectReportsError struct{ Err error }
)
// timeFormat is used to parse times in order to store time as ISO8601 format
const timeFormat = "2006-01-02T15:04:05-07:00"
func (c *client) Name() string {
	var name = "Clickhouse"
	return name
}
func New(chcon clickhouse.Conn) Client {
	return &client{chcon}
}
func (c *client) Ping(ctx context.Context) error {
	return c.chcon.Ping(ctx)
}
// CheckDB returns a time if previous reports are found, otherwise nil
func (c *client) CheckDB(ctx context.Context, region string) (string, error) {
	var start time.Time
	var err error
	if region == "" {
		return "", fmt.Errorf("error initializing database: %w", err)
	}
	if err = c.chcon.QueryRow(ctx, `
			SELECT
					MAX(end) as max_end
			FROM 
					carbondb.carbon_intensity_reports
			WHERE
					region = $1
			`, region).Scan(&start); err != nil {
				return "", err
			}
	// clickhouse returns a 1900s year if nothing is found
	if start.Year() < 2000 {
		log.Info(ctx, log.KV{K: "No logs for region:", V: region})
		return "", fmt.Errorf("no records for region")
	}
	start = start.UTC()
	fmt.Println("start date found by check DB:")
	fmt.Println(start)
	fmt.Println("for region:")
	fmt.Println(region)
	return start.Format(timeFormat), err
}

func (c *client) Init(ctx context.Context, test bool) error {
	if err := c.chcon.Ping(ctx); err != nil {
		if exception, ok := err.(*ch.Exception); ok {
			return fmt.Errorf("[%d] %s", exception.Code, exception.Message)
		}
		return err
	}
	if err := c.chcon.Exec(ctx, `CREATE DATABASE IF NOT EXISTS carbondb;`); err != nil {
		log.Errorf(ctx, err, "error initializing database: %w", err)
		return err
	}

	if err := c.chcon.Exec(ctx, `
	CREATE TABLE IF NOT EXISTS carbondb.carbon_intensity_reports (
		start DateTime,
		end DateTime,
		generatedrate Float64,
		marginalrate Float64,
		consumedrate Float64,
		region LowCardinality(String)
	)
	ENGINE = ReplicatedMergeTree('/clickhouse/{cluster}/tables/{shard}/{database}/{table}', '{replica}') 
	PARTITION BY (toYYYYMM(start), ignore(end))
	ORDER BY (region, start)
	SETTINGS index_granularity = 8192
	`); err != nil {
		return fmt.Errorf("error initializing clickhouse[%w]", err)
	}
	return nil
}

// SaveCarbonReports saves CO2 intensity reports in clickhouse
func (c *client) SaveCarbonReports(ctx context.Context, reports []*genpoller.CarbonForecast) (error) {
	res, err := c.chcon.PrepareBatch(ctx, `Insert INTO carbondb.carbon_intensity_reports (start,
		 end, generatedrate, marginalrate, consumedrate, region) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		return IncorrectReportsError{Err: fmt.Errorf("error in save carbon reports [%w]\n", err)}
	}
	if reports == nil {
		return NoReportsError{Err: fmt.Errorf("no reports given\n")}
	}
	for _, report := range reports {
		var startTime, _ = time.Parse(timeFormat, report.Duration.StartTime)
		var endTime, _ = time.Parse(timeFormat, report.Duration.EndTime)
		if err := res.Append(startTime.UTC(),
			 endTime.UTC(), report.GeneratedRate, report.MarginalRate,
			  report.ConsumedRate, report.Region); err != nil {
				return IncorrectReportsError{Err: fmt.Errorf("reports could not be saved in the specified format: [%w]", err)}
			}
	}
	return res.Send()
}
func (err NoReportsError) Error() string { return err.Err.Error() }
func (err IncorrectReportsError) Error() string { return err.Err.Error() }





