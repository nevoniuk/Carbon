package storage
//method to retrieve values from clickhouse
import (
	"context"
	"fmt"
	//"time"
	ch "github.com/ClickHouse/clickhouse-go/v2"
	"github.com/crossnokaye/carbon/clients/clickhouse"
	//gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
)

type (


	Client interface {
		// Init development database
		Name() string

		Init(context.Context, bool) error

		Ping(context.Context) error
		
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
			CREATE TABLE IF NOT EXISTS carbondb.power_reports (
				start DateTime,
				end DateTime,
				generatedrate Float64,
				
				) Engine =  MergeTree()
				ORDER BY (start)
	`) 
	
	return err
	
		
}
