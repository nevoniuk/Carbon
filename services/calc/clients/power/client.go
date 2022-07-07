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
	//"fmt"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
	//"github.com/google/uuid"
	//"strings"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	genvalues "github.com/crossnokaye/past-values/services/past-values/gen/past_values"
	genvaluesc "github.com/crossnokaye/past-values/services/past-values/gen/grpc/past_values/client"
)
var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"

type (
	// Client interface to past-values service
	Client interface {
		GetPower(context.Context, string, []genvalues.UUID, int64, time.Time, time.Time) (*gencalc.ElectricalReport, error)
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
	 start time.Time, end time.Time) (*gencalc.ElectricalReport, error) {
	
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
	newRes, err := toPower(res)
	if err != nil {
		return nil, err
	}
	return newRes, nil
}

func toPower(r interface{}) (*gencalc.ElectricalReport, error) {
	//casting the response
	res := r.([]*genvalues.HistoricalValues)
	//making an array of the reports I want to return
	var report *gencalc.ElectricalReport
	stamps := make([]*gencalc.PowerStamp, len(res))
	//each analog point contains an ID and an array of values
	var startReport time.Time
	var endReport time.Time

	var analogPoint int
	//iterate over historical values
	for i, historicalValues := range res {
		//i == 1 should be analog points
		if i == 1 {

			for j, analog := range historicalValues.Analog {

				if j == analogPoint {

					var minCounter time.Time
					var previousPoint float64
					var counter = 0

					for p, point := range analog.Values {

						if p == 0 {
							var err error
							startReport, err = time.Parse(timeFormat, *point.Timestamp)
							if err != nil {
								return nil , err
							}
							minCounter = startReport
						}

						if p == (len(analog.Values) - 1) {
							var err error
							endReport, err = time.Parse(timeFormat, *point.Timestamp)
							if err != nil {
								return nil , err
							}
							
						}
						//if new time is 1 minute  or more greater than the counter and is not 0
						
						pointTime, errpoint := time.Parse(timeFormat, *point.Timestamp)
						if errpoint != nil  {
							return nil, errpoint
						}

						 if pointTime.After(minCounter.Add(time.Minute)) && (previousPoint != 0) {

							var temp = minCounter.Add(time.Minute).Format(timeFormat)

							stamps[counter] = &gencalc.PowerStamp{
								Time: &temp,
								GenRate: &previousPoint,
							}
							counter += 1
							//set to point time in case that elapsed time is greater than a minute
							minCounter = pointTime
						 }
						//read in values that aren't 0
						//values are per second
						previousPoint = *point.Value
						
					}
					
				}
				
			}
		}
		
	}

	report = &gencalc.ElectricalReport{
		Period: &gencalc.Period{StartTime: startReport.Format(timeFormat), EndTime: endReport.Format(timeFormat)},
		Stamp: stamps,
		IntervalType: "minute",
	}

	return report, nil
}
