package storage

import (
	"context"
	"fmt"
	//"go/constant"
	"time"
	ch "github.com/ClickHouse/clickhouse-go/v2"

	"github.com/crossnokaye/carbon/clients/clickhouse"
	//"github.com/crossnokaye/rates/clients/clickhouse"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

type (


	Client interface {
		// Init development database
		Name() string

		Init(context.Context, bool) error
		// Retrieve last date there is available data in clickhouse
		CheckDB(context.Context, string) (string, error)
		// Save report for carbon intensity event only
		SaveCarbonReports(context.Context, []*genpoller.CarbonForecast) (error)

		//GetCarbonReports(context.Context, string, genpoller.Period, string) ([]*genpoller.CarbonForecast)

		//SaveFuelReports(ctx context.Context, reports []*FuelsForecast)
		
		Ping(context.Context) error
		
		//query data then update tables with aggregate information(generated data only)
		//carbon data only
		GetAggregateReports(context.Context, []*genpoller.Period, string, string) ([]*genpoller.AggregateData, error)

		SaveAggregateReports(context.Context, []*genpoller.AggregateData) (error)
	}

	client struct {
		chcon clickhouse.Conn
	}
)

var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"

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
	//need to convert time written to string
	if err = c.chcon.QueryRow(ctx, `
			SELECT
					MAX(start)
			FROM carbon_reports
			WHERE duration = $1
			GROUP BY
			     start, duration
			`, "hourly").Scan(&start); err != nil {
				fmt.Errorf("error reading time in CheckDB\n")
				return convertTimeString(ctx, start), err//time would be null
			}
	return convertTimeString(ctx, start), err
}

func convertTimeString(ctx context.Context, t time.Time) (string) {
	return t.Format("2006-01-02T15:04:05-07:00")
}


func (c *client) Init(ctx context.Context, test bool) error {
	if err := c.chcon.Ping(ctx); err != nil {
		if exception, ok := err.(*ch.Exception); ok {
			fmt.Printf("HERE")
			return fmt.Errorf("[%d] %s", exception.Code, exception.Message)
		}
		//returns this error
		fmt.Printf("PING ERROR")
		return err
	}
	fmt.Printf("HERE2")
	if err := c.chcon.Exec(ctx, `CREATE DATABASE IF NOT EXISTS carbondb;`); err != nil {
		return err
	}

	var err error 
		//research replication engine
	err = c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS carbon_reports (
					duration String,
					start DateTime,
					end DateTime,
					generatedrate Float64,
					marginalrate Float64,
					consumedrate Float64,
					generatedsource String
				) Engine =  MergeTree()
				ORDER BY (start, duration)
	`) 
	if err != nil {
		return err
	}

	err = c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS aggregate_reports (
					duration String,
					start DateTime,
					end DateTime,
					average Float64,
					min Float64,
					max Float64,
					sum Float64,
					count Float64
				) Engine =  MergeTree()
				ORDER BY (duration, start)
			`) 
	return err
	
			/*
			err = c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS fuel_report (
					duration String,
					start DateTime,
					end DateTime,
					eventtype, String
					fuel1name String,
					flue1value Float64,
					fuel2name String,
					flue2value Float64,
					fuel3name String,
					flue3value Float64,
					fuel4name String,
					flue4value Float64,
					fuel5name String,
					flue5value Float64,
					fuel6name String,
					flue6value Float64,
					fuel7name String,
					flue7value Float64,
					fuel8name String,
					flue8value Float64,
					generatedsource String,
					marginalsource String,
					emissionfactor String,
					average Float64,
					max Float64,
					min Float64,
					sum Float64
			)`)
			*/
}

func (c *client) SaveCarbonReports(ctx context.Context, reports []*genpoller.CarbonForecast) (error) {
	res, err := c.chcon.PrepareBatch(ctx, `Insert INTO carbon_reports (duration, start,
		 end, generatedrate, marginalrate, consumedrate, generatedsource, marginalsource, consumedsource,
		  emissionfactor) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`)
	if err != nil {
		return err
	}

	for _, report := range reports {
		var startTime, err1 = time.Parse(timeFormat, report.Duration.StartTime)
		if err1 != nil {
			fmt.Errorf("Timestamp %s in observation %v could not be parsed into a time correctly")
			continue
		}
		
		var endTime, err2 = time.Parse(timeFormat, report.Duration.EndTime)
		if err2 != nil {
			fmt.Errorf("Timestamp %s in observation %v could not be parsed into a time correctly")
			continue
		}
		if err := res.Append(report.Duration, startTime,
			 endTime, report.GeneratedRate, report.MarginalRate,
			  report.ConsumedRate, report.GeneratedSource); err != nil {
				return err
			}
	}
	
	return res.Send()
}

func (c *client) GetAggregateReports(ctx context.Context,
	 periods []*genpoller.Period, region string, duration string) ([]*genpoller.AggregateData, error) {
	
	var finalaggdata []*genpoller.AggregateData
	/**
	var aggdata *genpoller.AggregateData
	
	//should return aggregate date for hour reports, daily reports etc...
	var min float64
	var max float64
	var sum float64
	var count int
	var average float64
	
	for _, period := range periods {
		rows, err := c.chcon.Query(ctx, `
			SELECT
				start,
				end,
				region,
				MIN(generatedrate),
				MAX(generatedrate),
				COUNT(generatedrate),
				SUM(generatedrate),
				AVG(generatedrate)
			FROM
				aggregate_reports
			WHERE   
				start >= $1 AND end <= $2 AND region = $3
			GROUP BY
				duration
			ORDER BY
				duration`, period.StartTime, period.EndTime, region).Scan(&min, &max, &count, &sum, &average)
	aggdata = &genpoller.AggregateData{average, min, max, sum, count, period, duration}
	finalaggdata = append(finalaggdata, aggdata)
	}
	*/
	return finalaggdata, nil
}

func (c *client) SaveAggregateReports(ctx context.Context, aggData []*genpoller.AggregateData) (error) {
	//save aggregate data in new table
	batch, err := c.chcon.PrepareBatch(ctx, `INSERT INTO aggregate_reports (duration, start, end, average, min, max, sum, count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
	if err != nil {
		return err
	}
	for _, p := range aggData {
		if err := batch.Append(p.ReportType, p.Duration.StartTime, p.Duration.EndTime,
			 p.Average, p.Min, p.Max, p.Sum, p.Count); err != nil {
			return err
		}
	}
	return batch.Send()
}
/**
func (c *client) GetCarbonReports(ctx context.Context, durationtype string,
	timeInterval genpoller.Period, region string) ([]*genpoller.CarbonForecast) {
		var reports []*genpoller.CarbonForecast
	rows, err := c.chcon.Query(ctx, `
		SELECT
			*
		FROM weather_observations

		WHERE
			region = $1 AND duration = $2

		ORDER BY
			end
	`, region, durationtype)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		p := genpoller.CarbonForecast{}
		if err := rows.Scan(&p.duration, &p.region, &p.generated_rate...); err != nil {
			return nil, err
		}
		reports = append(reports, &p)
	}
	return reports, nil


}
*/



