package storage

import (
	"context"
	"fmt"
	"go/constant"
	"time"

	ch "github.com/ClickHouse/clickhouse-go/v2"

	//"github.com/crossnokaye/rates/clients/clickhouse"
	genpoller "github.com/crossnokaye/carbon/services/poller/gen/data"
)

type (


	Client interface {
		// Init development database
		Init(context.Context, bool) error
		// Retrieve last date there is available data in clickhouse
		CheckDB(ctx context.Context, region string) (genpoller.Period, error)
		// Save report for carbon intensity event only
		SaveCarbonReports(ctx context.Context, reports []*CarbonResponse)

		GetCarbonReports(ctx context.Context, string, genpoller.Period, string) (reports []*CarbonResponse)

		//SaveFuelReports(ctx context.Context, reports []*FuelsForecast)

		Ping(ctx context.Context) error
		
		//query data then update tables with aggregate information
		get_aggregate_data(ctx context.Context, string, genpoller.Period, string)

		save_aggregate_data(ctx context.Context, genpoller.AggregateData)
	}

	const(
		regionstartdates := map[string]string{"AESO": "2020-05-15T16:00:00+00:00",
		"BPA": "2018-01-01T08:00:00+00:00",
		"CAISO": "2018-04-10T07:00:00+00:00",
		"ERCO": "2018-07-02T05:05:00+00:00",
		"IESO": "2017-12-31T05:00:00+00:00",
		"ISONE": "2015-01-01T05:00:00+00:00",
		"MISO": "2018-01-01T05:00:00+00:00",
		"NYSIO": "2017-12-01T05:05:00+00:00",
		"NYISO.NYCW": "2019-01-01T00:00:00+00:00",
		"NYISO.NYLI": "2019-01-01T00:00:00+00:00",
		"NYISO.NYUP": "2019-01-01T00:00:00+00:00",
		"PJM": "2017-07-01T04:05:00+00:00",
		"SPP": "2017-12-31T00:00:00+00:00",
		"EIA": "2019-01-01T05:00:00+00:00"
	}
	)
)

func (c *client) Ping(ctx context.Context) error {
	return c.chcon.Ping(ctx)
}

func (svc *Service) CheckDB(ctx context.Context, region string) (timeInterval genpoller.Period) {
	startDate = regionstartdates[region]
	//var date time.Time.UTC()
	if err := c.chcon.QueryRow(ctx, `
			SELECT
					start, max(end),
			FROM carbon_reports

			GROUP BY
			     start, duration
			`).Scan(&timeInterval{StartTime, EndTime}); err != nil {
				return time.Time, err
			}
	if timeInterval == nil {
		timeInterval = timeInterval{startDate, time.Now()}
	}
	return timeInterval
	//if result is null then set to startdate above
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
			CREATE TABLE IF NOT EXISTS carbon_reports (
					duration String,
					start DateTime,
					end DateTime,
					generatedrate Float64,
					marginalrate Float64,
					consumedrate Float64,
					generatedsource String,
					marginalsource String,
					consumedsource String,
					emissionfactor String,

					max Float64,
					min Float64,
					sum Float64
				) Engine =  MergeTree()
				ORDER BY (start, duration)
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

func (c *client) SaveCarbonReports(ctx context.Context, reports []*CarbonResponse) {
	res, err := c.chcon.PrepareBatch(ctx, `Insert INTO carbon_reports (duration, start,
		 end, generatedrate, marginalrate, consumedrate, generatedsource, marginalsource, consumedsource,
		  emissionfactor) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`)
	if err != nil {
		return err
	}
	for _, reporttype := range reports {
		for _, report := range reporttype {
			//append all data in carbon forecast report
			if err := stmt.Append(report.duration); err != nil {
				return err
			}
		}
	}
	

	return stmt.Send()
}

func (svc *Service) get_aggregate_data(ctx context.Context, durationtype string,
	 timeInterval genpoller.Period, region string) ([]*genpoller.AggregateData) {
	var aggdata []*genpoller.AggregateData
	var finalaggdata []*genpoller.AggregateData
	
	//should return aggregate date for hour reports, daily reports etc...
	var min float64
	var max float64
	var sum float64
	var count int

	rows, err := c.chcon.Query(ctx, `
		SELECT
				MIN(generatedrate)
				MAX(generatedrate)
				COUNT(generatedrate)
				SUM(generatedrate)
		FROM
				carbon_reports
		WHERE   
				region = $1 AND duration = $2
		GROUP BY
				duration
		ORDER BY
				duration
	`, region, durationtype).Scan(&min, &max, &count, &sum)
	return &finalaggdata{
		min: min,
		max: max,
		sum: sum,
		count: count
	}
}

func (c *client) save_aggregate_data([]*genpoller.AggregateData) {
	//save aggregate data in new table
}

func (c *client) GetCarbonReports(ctx context.Context, durationtype string,
	timeInterval genpoller.Period, region string) {
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




