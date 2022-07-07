package power
//this file will client the past-values service
//similar to the control points client
//method for getting control points
//method for using control points as input to service

//1. get control points using fc store
//2. make calls to past-values service
//3. store power data in structures defined here
import (
	"context"
	"time"
	"fmt"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
	"github.com/google/uuid"
	"strings"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	genvalues "github.com/crossnokaye/past-values/services/past-values/gen/past_values"
	genvaluesc "github.com/crossnokaye/past-values/services/past-values/gen/grpc/past_values/client"
)
var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"

type (
	// Client interface to past-values service
	Client interface {
		GetPower(context.Context, string, []genvalues.UUID, int64, time.Time, time.Time) ([]*genvaluesc.HistoricalValues, error)
	}

	// client implements the Client interface.
	client struct {
		getPower goa.Endpoint
	}

	// Rate is the energy rate for a given hour.
	
	

	// ErrNotFound is returned when an org or facility is not found.
	ErrNotFound struct{ Err error }
)

func New(conn *grpc.ClientConn) Client {
	c := genvaluesc.NewClient(conn, grpc.WaitForReady(true))
	//method to client
	return &client{
		getPower: c.GetValues(),
	}
}

func (c *client) GetPower(ctx context.Context, orgID string, controlPoints []genvalues.UUID, interval int64,
	 start time.Time, end time.Time) ([]*genvalues.HistoricalValues, error) {
	
	var startTime = start.Format(timeFormat)
	var endTime = end.Format(timeFormat)
	p := genvalues.ValuesQuery{
		OrgID: genvalues.UUID(orgID),
		PointIds: controlPoints,
		Start: startTime,
		End: endTime,
		Interval: interval,
	}

	res, err := c.getPower(ctx, p)
	//res is historical values
	//historical values = discrete points, analog points and structures
	if err != nil {
		return nil, err
	}

	return res, nil
}

func toPower(r interface{}) []*gencalc.ElectricalReport {
	//casting the response
	res := r.([]*genvalues.HistoricalValues)
	//making an array of the reports I want to return
	reports := make([]*gencalc.ElectricalReport, len(res))
	for i, r := range res {
		//read r.Analog[i].Values into stamps
		reports[i] = &gencalc.ElectricalReport{
			stamps := make([]*gencalc.PowerStamp, len(res))
			Stamp: r.Analog[i].Values
		}
	}

	return rates
}
