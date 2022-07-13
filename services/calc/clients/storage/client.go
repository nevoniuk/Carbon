package storage
import (
	"context"
	"fmt"
	"time"
	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/crossnokaye/carbon/clients/clickhouse"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
)

type (


	Client interface {
		Name() string
		Init(context.Context, bool) error
		Ping(context.Context) error
		GetCarbonReport(context.Context, []*gencalc.Period, string, string) ([]*gencalc.CarbonReport, error)
		
	}

	client struct {
		chcon clickhouse.Conn
	}
)
const(
	timeFormat = "2006-01-02T15:04:05-07:00"
	dateFormat = "2006-01-02"
)

var reportdurations [5]string = [5]string{ "minute", "hourly", "daily", "weekly", "monthly"}
func (c *client) Name() string {
	var name = "Clickhouse"
	return name
}

func New(chcon clickhouse.Conn) Client {
	return &client{chcon}
}

//ping is the time it takes for a small data set to be transmitted from the device to a server on the internet 
func (c *client) Ping(ctx context.Context) error {
	return c.chcon.Ping(ctx)
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

	var err error 
	err = c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS carbondb.power_reports (
				start DateTime,
				end DateTime,
				generatedrate Float64,
				) Engine =  MergeTree()
				ORDER BY (start)
	`) 
	
	return err	
}

func (c *client) GetCarbonReport(ctx context.Context, duration []*gencalc.Period, intervalType string, region string) ([]*gencalc.CarbonReport, error) {
	var reports []*gencalc.CarbonReport
	var report *gencalc.CarbonReport
	var averagegen float64

	for _, period := range duration {
		var newstart, starterr = time.Parse(timeFormat, period.StartTime)
		if starterr != nil {
			return nil, starterr
		}
		var newend, enderr = time.Parse(timeFormat, period.EndTime)
		if enderr != nil {
			return nil, enderr
		}
		
		rows := c.chcon.QueryRow(ctx,`
		SELECT
			AVG(generatedrate) AS generatedate,
			AVG(marginalrate) AS marginalrate,
			AVG(consumedrate) AS consumedrate
		FROM 
			carbondb.carbon_reports
		WHERE
			region = $1 AND start >= $2 AND end <= $3 AND duration = $4
		GROUP BY region
				`, region, newstart.UTC(), newend.UTC(), intervalType)
		err := rows.Scan(&averagegen)

		if err != nil {
			return nil, fmt.Errorf("Error could not get report [%s]", err)
		}
		var duration = &gencalc.Period{StartTime: period.StartTime, EndTime: period.EndTime}
		report = &gencalc.CarbonReport{GeneratedRate: averagegen, Duration: duration, DurationType: intervalType, Region: region}
		reports = append(reports, report)	
	}

	return reports, nil
}

