package storage

import (
	"context"
	"fmt"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/crossnokaye/carbon/clients/clickhouse"
	"github.com/crossnokaye/carbon/model"
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
		// Obtains carbon intensity reports for a given report interval, time duration, and region
		GetAggregateReports(context.Context, []*genpoller.Period, string, string) ([]*genpoller.CarbonForecast, error)
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
//should maybe check for the end instead and check that all reports were written(success)
func (c *client) CheckDB(ctx context.Context, region string) (string, error) {
	var start time.Time
	var err error
	if region == "" {
		log.Error(ctx, fmt.Errorf("err: no region provided in CheckDB\n"))
		return "", nil
	}
	if err = c.chcon.QueryRow(ctx, `
			SELECT
					MAX(end) as max_end
			FROM 
					carbondb.carbon_reports
			WHERE
					region = $1 AND duration = $2
			`, region, model.Weekly).Scan(&start); err != nil {

				return "", NoReportsError{Err: fmt.Errorf("error in checkDB [%w]\n", err)}
			}
	// clickhouse returns a 1900s year if nothing is found
	if start.Year() < 2000 {
		err = fmt.Errorf("No records for given region")
		return "", NoReportsError{Err: fmt.Errorf("error in checkDB [%w]\n", err)}
	}
	start = start.UTC()
	return start.Format(timeFormat), err
}

func (c *client) Init(ctx context.Context, test bool) error {
	if err := c.chcon.Ping(ctx); err != nil {
		if exception, ok := err.(*ch.Exception); ok {
			return fmt.Errorf("[%d] %s", exception.Code, exception.Message)
		}
		return err
	}
	fmt.Println("intitialized clickhouse")
	if err := c.chcon.Exec(ctx, `CREATE DATABASE IF NOT EXISTS carbondb;`); err != nil {
		log.Errorf(ctx, err, "error initializing database: %w", err)
		return err
	}
	fmt.Println("intitialized clickhouse")
	if err := c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS carbondb.carbon_reports (
					start DateTime,
					end DateTime,
					generatedrate Float64,
					marginalrate Float64,
					consumedrate Float64,
					region String,
					duration String
				) Engine =  MergeTree()
				ORDER BY (start)
	`); err != nil {
		log.Errorf(ctx, err, "error initializing database: %w", err)
		return fmt.Errorf("Error initializing clickhouse[%w]", err)
	}
	fmt.Println("intitialized clickhouse")
	return nil
}

// SaveCarbonReports saves CO2 intensity reports in clickhouse
func (c *client) SaveCarbonReports(ctx context.Context, reports []*genpoller.CarbonForecast) (error) {
	res, err := c.chcon.PrepareBatch(ctx, `Insert INTO carbondb.carbon_reports (start,
		 end, generatedrate, marginalrate, consumedrate, region, duration) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
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
			  report.ConsumedRate, report.Region, report.DurationType); err != nil {
				return IncorrectReportsError{Err: fmt.Errorf("reports could not be saved in the specified format: [%w]", err)}
			}
	}
	return res.Send()
}

// GetAggregateReports queries clickhouse for average CO2 intensity data
func (c *client) GetAggregateReports(ctx context.Context,
	 periods []*genpoller.Period, region string, duration string) ([]*genpoller.CarbonForecast, error) {
	var finalaggdata []*genpoller.CarbonForecast
	var aggdata *genpoller.CarbonForecast
	var averagegen float64
	var averagemarg float64
	var averagecons float64
	for _, period := range periods {
		newstart, _ := time.Parse(timeFormat, period.StartTime)
		newend, _ := time.Parse(timeFormat, period.EndTime)
		rows := c.chcon.QueryRow(ctx,`
		SELECT
			AVG(generatedrate) AS generatedrate,
			AVG(marginalrate) AS marginalrate,
			AVG(consumedrate) AS consumedrate
		FROM 
			carbondb.carbon_reports
		WHERE
			region = $1 AND start >= $2 AND end <= $3
		GROUP BY region
				`, region, newstart.UTC(), newend.UTC())
		err := rows.Scan(&averagegen, &averagemarg, &averagecons)
		if err != nil {
			return nil, NoReportsError{Err: fmt.Errorf("no data for Region %s and start %s and end %s with err: %w", region, period.StartTime, period.EndTime, err)}
		}
		aggdata = &genpoller.CarbonForecast{GeneratedRate: averagegen, MarginalRate: averagemarg, ConsumedRate: averagecons,
			Duration: period, DurationType: duration, Region: region}
		finalaggdata = append(finalaggdata, aggdata)	
	}
	return finalaggdata, nil
}

func (err NoReportsError) Error() string { return err.Err.Error() }
func (err IncorrectReportsError) Error() string { return err.Err.Error() }





