package storage

import (
	"context"
	"fmt"
	"time"
	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/crossnokaye/carbon/clients/clickhouse"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

type (
	Client interface {
		Name() string
		Init(context.Context, bool) error
		Ping(ctx context.Context) error
		CheckDB(context.Context, string) (string, error)
		SaveCarbonReports(context.Context, []*genpoller.CarbonForecast) (error)
		GetAggregateReports(context.Context, []*genpoller.Period, string, string) ([]*genpoller.CarbonForecast, error)
	}

	client struct {
		chcon clickhouse.Conn
	}
	NoReportsError struct{ Err error }
	IncorrectReportsError struct{ Err error }
)
const(
	//timeFormat is used to parse times in order to store time as ISO8601 format
	timeFormat = "2006-01-02T15:04:05-07:00"
)

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
//CheckDB returns a time if previous reports are found, otherwise nil
func (c *client) CheckDB(ctx context.Context, region string) (string, error) {
	var start time.Time
	var err error
	if err = c.chcon.QueryRow(ctx, `
			SELECT
					MAX(start) as max_start
			FROM 
					carbondb.carbon_reports
			WHERE
					region = $1
			`, region).Scan(&start); err != nil {

				return "", fmt.Errorf("Error in CheckDB: [%s]\n")
			}
	if start.Year() < 2000 {
		err = fmt.Errorf("No records for given region")
		return "", err
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
	if err := c.chcon.Exec(ctx, `CREATE DATABASE IF NOT EXISTS carbondb;`); err != nil {
		return err
	}

	/**
	err := c.chcon.Exec(ctx, `
			DROP TABLE carbondb.carbon_reports
	`)
	if err != nil {
		return fmt.Errorf("Error initializing clickhouse[%s]", err)
	}
	*/
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
		return fmt.Errorf("Error initializing clickhouse[%s]", err)
	}
	return nil
}

//SaveCarbonReports saves CO2 intensity reports in clickhouse
func (c *client) SaveCarbonReports(ctx context.Context, reports []*genpoller.CarbonForecast) (error) {
	res, err := c.chcon.PrepareBatch(ctx, `Insert INTO carbondb.carbon_reports (start,
		 end, generatedrate, marginalrate, consumedrate, region, duration) VALUES ($1, $2, $3, $4, $5, $6, $7)`)

	if err != nil {
		var invaliderr = IncorrectReportsError{Err: fmt.Errorf("error in save carbon reports [%w]\n", err)}
		return invaliderr
	}
	var noRepErr = NoReportsError{Err: fmt.Errorf("no reports given in save carbon reports\n")}
	if reports == nil {
		return noRepErr
	}

	for _, report := range reports {
		var startTime, err1 = time.Parse(timeFormat, report.Duration.StartTime)
		if err1 != nil {
			return fmt.Errorf("timestamp %s in observation %v could not be parsed into a time correctly error: [%s]",
			report.Duration.StartTime, report, err1)
		}
	
		var endTime, err2 = time.Parse(timeFormat, report.Duration.EndTime)
		if err2 != nil {
			return fmt.Errorf("timestamp %s in observation %v could not be parsed into a time correctly error: [%s]",
			report.Duration.EndTime, report, err2)
		}

		if err := res.Append(startTime.UTC(),
			 endTime.UTC(), report.GeneratedRate, report.MarginalRate,
			  report.ConsumedRate, report.Region, report.DurationType); err != nil {
				return IncorrectReportsError{Err: fmt.Errorf("reports could not be saved in the specified format: [%w]", err)}
			}
	}
	return res.Send()
}

//GetAggregateReports queries clickhouse for average CO2 intensity data
func (c *client) GetAggregateReports(ctx context.Context,
	 periods []*genpoller.Period, region string, duration string) ([]*genpoller.CarbonForecast, error) {
	var finalaggdata []*genpoller.CarbonForecast
	var aggdata *genpoller.CarbonForecast
	var averagegen float64
	var averagemarg float64
	var averagecons float64
	for i, period := range periods {
		newstart, err := time.Parse(timeFormat, period.StartTime)
		if err != nil {
			return nil, fmt.Errorf("timestamp %s in observation %v could not be parsed into a time correctly error: [%s]",
			period.StartTime, i, err)
		}
		newend, err := time.Parse(timeFormat, period.EndTime)
		if err != nil {
			return nil, fmt.Errorf("timestamp %s in observation %v could not be parsed into a time correctly error: [%s]",
			period.EndTime, i, err)
		}
		rows := c.chcon.QueryRow(ctx,`
		SELECT
			AVG(generatedrate) AS generatedate,
			AVG(marginalrate) AS marginalrate,
			AVG(consumedrate) AS consumedrate
		FROM 
			carbondb.carbon_reports
		WHERE
			region = $1 AND start >= $2 AND end <= $3
		GROUP BY region
				`, region, newstart.UTC(), newend.UTC())
		err = rows.Scan(&averagegen, &averagemarg, &averagecons)
		if err != nil {
			return nil, NoReportsError{Err: fmt.Errorf("no data for Region %s and start %s and end %s", region, period.StartTime, period.EndTime)}
		}
		aggdata = &genpoller.CarbonForecast{GeneratedRate: averagegen, MarginalRate: averagemarg, ConsumedRate: averagecons,
			Duration: period, DurationType: duration, Region: region}
		finalaggdata = append(finalaggdata, aggdata)	
	}
	return finalaggdata, nil
}

func (err NoReportsError) Error() string { return err.Err.Error() }
func (err IncorrectReportsError) Error() string { return err.Err.Error() }





