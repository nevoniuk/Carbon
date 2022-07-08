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
		// Init development database
		Name() string

		Init(context.Context, bool) error
		// Retrieve last date there is available data in clickhouse
		CheckDB(context.Context, string) (string, error)
		
		SaveCarbonReports(context.Context, []*genpoller.CarbonForecast) (error)

		
		Ping(context.Context) error
		
		GetAggregateReports(context.Context, []*genpoller.Period, string, string) ([]*genpoller.CarbonForecast, error)

	}

	client struct {
		chcon clickhouse.Conn
	}
)
const(
	timeFormat = "2006-01-02T15:04:05-07:00"
	dateFormat = "2006-01-02"
)


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

//meant to return a "start" time for query to begin
//returns a time if previous reports are found, otherwise nil
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

				return "", fmt.Errorf("Error in checkDB: [%s]\n")
			}
	fmt.Println("START IS")
	fmt.Println(start)
	if start.Year() < 2000 {
		err = fmt.Errorf("No records for given region")
		return "", err
	}
	return convertTimeString(ctx, start.UTC()), err
}

func convertTimeString(ctx context.Context, t time.Time) (string) {
	return t.Format("2006-01-02T15:04:05-07:00")
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
		
	/**
	err = c.chcon.Exec(ctx, `
			DROP TABLE carbondb.carbon_reports
	`)
	*/

	err = c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS carbondb.carbon_reports (
					start DateTime,
					end DateTime,
					generatedrate Float64,
					marginalrate Float64,
					consumedrate Float64,
					generatedsource String,
					region String,
					duration String
				) Engine =  MergeTree()
				ORDER BY (start)
	`) 
	
	return err
	
		
}

func (c *client) SaveCarbonReports(ctx context.Context, reports []*genpoller.CarbonForecast) (error) {
	res, err := c.chcon.PrepareBatch(ctx, `Insert INTO carbondb.carbon_reports (start,
		 end, generatedrate, marginalrate, consumedrate, generatedsource, region, duration) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		return fmt.Errorf("Error in Save Carbon Reports [%s]", err)
	}

	for _, report := range reports {
		var startTime, err1 = time.Parse(timeFormat, report.Duration.StartTime)
		fmt.Println("save reports start time is")
		fmt.Println(startTime)
		if err1 != nil {

			return fmt.Errorf("Timestamp %s in observation %v could not be parsed into a time correctly error: [%s]",
			report.Duration.StartTime, report, err1)
		}
		
		var endTime, err2 = time.Parse(timeFormat, report.Duration.EndTime)
		fmt.Println("save reports endtime is ")
		fmt.Println(endTime)
		if err2 != nil {

			return fmt.Errorf("Timestamp %s in observation %v could not be parsed into a time correctly error: [%s]",
			report.Duration.EndTime, report, err2)
		}
		if err := res.Append(startTime.UTC(),
			 endTime.UTC(), report.GeneratedRate, report.MarginalRate,
			  report.ConsumedRate, report.GeneratedSource, report.Region, report.DurationType); err != nil {
				return fmt.Errorf("Error saving carbon reports: [%s]", err)
			}
	}

	return res.Send()
}

func (c *client) GetAggregateReports(ctx context.Context,
	 periods []*genpoller.Period, region string, duration string) ([]*genpoller.CarbonForecast, error) {
	
	var finalaggdata []*genpoller.CarbonForecast
	
	var aggdata *genpoller.CarbonForecast
	
	
	var averagegen float64
	var averagemarg float64
	var averagecons float64
	
	for _, period := range periods {

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
			region = $1 AND start >= $2 AND end <= $3
		GROUP BY region
				`, region, newstart.UTC(), newend.UTC())
		
		err := rows.Scan(&averagegen, &averagemarg, &averagecons)

		if err != nil {
			return nil, fmt.Errorf("Error could not get report [%s]", err)
		}

		aggdata = &genpoller.CarbonForecast{GeneratedRate: averagegen, MarginalRate: averagemarg, ConsumedRate: averagecons,
			Duration: period, DurationType: duration, GeneratedSource: "", Region: region}
		finalaggdata = append(finalaggdata, aggdata)	
	}

	return finalaggdata, nil
}





