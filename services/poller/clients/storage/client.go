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
		GetAggregateReports(context.Context, []*genpoller.Period, string, string) ([]*genpoller.CarbonForecast, error)

		SaveAggregateReports(context.Context, []*genpoller.AggregateData) (error)
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
	//1.check if any hourly reports are missing
	if err = c.chcon.QueryRow(ctx, `
			SELECT
					MAX(start)
			FROM carbondb.carbon_reports
			WHERE region = $1
			`, region).Scan(&start); err != nil {
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
			return fmt.Errorf("[%d] %s", exception.Code, exception.Message)
		}
		return err
	}

	if err := c.chcon.Exec(ctx, `CREATE DATABASE IF NOT EXISTS carbondb;`); err != nil {
		return err
	}

	var err error 
		//research replication engine
	err = c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS carbondb.carbon_reports (
					start DateTime,
					end DateTime,
					generatedrate Float64,
					marginalrate Float64,
					consumedrate Float64,
					generatedsource String
				) Engine =  MergeTree()
				ORDER BY (start)
	`) 
	
	/**
	err = c.chcon.Exec(ctx, `
			CREATE TABLE IF NOT EXISTS carbondb.aggregate_reports (
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
	*/
	return err
	
		
}

func (c *client) SaveCarbonReports(ctx context.Context, reports []*genpoller.CarbonForecast) (error) {
	res, err := c.chcon.PrepareBatch(ctx, `Insert INTO carbondb.carbon_reports (start,
		 end, generatedrate, marginalrate, consumedrate, generatedsource) VALUES ($1, $2, $3, $4, $5, $6)`)
	if err != nil {
		fmt.Println("save carbon reports error preparing report")
		return err
	}

	for _, report := range reports {
		var startTime, err1 = time.Parse(timeFormat, report.Duration.StartTime)
		fmt.Println("starttime is ")
		fmt.Println(startTime)
		if err1 != nil {
			fmt.Errorf("Timestamp %s in observation %v could not be parsed into a time correctly")
			continue
		}
		
		var endTime, err2 = time.Parse(timeFormat, report.Duration.EndTime)
		fmt.Println("endtime is ")
		fmt.Println(endTime)
		if err2 != nil {
			fmt.Errorf("Timestamp %s in observation %v could not be parsed into a time correctly")
			continue
		}
		if err := res.Append(startTime,
			 endTime, report.GeneratedRate, report.MarginalRate,
			  report.ConsumedRate, report.GeneratedSource); err != nil {
				return err
			}
	}
	fmt.Printf("reports sent\n")
	return res.Send()
}

func (c *client) GetAggregateReports(ctx context.Context,
	 periods []*genpoller.Period, region string, duration string) ([]*genpoller.CarbonForecast, error) {
	
	var finalaggdata []*genpoller.CarbonForecast
	
	var aggdata *genpoller.CarbonForecast
	
	//should return aggregate date for hour reports, daily reports etc...
	startTime := time.Time{}
	endTime := time.Time{}
	var averagegen float64
	var averagemarg float64
	var averagecons float64
	var tempregion string
	
	for _, period := range periods {
		rows := c.chcon.QueryRow(ctx, `
			SELECT
				AVG(generatedrate),
				AVG(marginalrate),
				AVG(consumedrate)
			FROM
				(SELECT *
				FROM 
					carbondb.carbon_reports
				WHERE
					start >= $1 AND end <= $2 AND region == $3
				ORDER BY
					start) 
				`, period.StartTime, period.EndTime, region)
		//aggdata := *genpoller.CarbonForecast{}
		if err := rows.Scan(&startTime, &endTime, &tempregion, &averagegen, &averagemarg, &averagecons); err != nil {
			return nil, err
		}
		fmt.Printf(" g average is %f\n", averagegen)
		fmt.Printf(" m average is %f\n", averagemarg)
		fmt.Printf(" c average is %s\n", averagecons)
	
	aggdata = &genpoller.CarbonForecast{GeneratedRate: averagegen, MarginalRate: averagemarg, ConsumedRate: averagecons,
	Duration: periods[0], GeneratedSource: "", Region: region}
	finalaggdata = append(finalaggdata, aggdata)
	}

	return finalaggdata, nil
}


func (c *client) SaveAggregateReports(ctx context.Context, aggData []*genpoller.AggregateData) (error) {
	//save aggregate data in new table
	batch, err := c.chcon.PrepareBatch(ctx, `INSERT INTO carbondb.aggregate_reports (duration, start, end, average, min, max, sum, count) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)
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




