package power

import (
	"context"
	"fmt"
	"time"
	"github.com/crossnokaye/carbon/model"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	genvaluesc "github.com/crossnokaye/past-values/services/past-values/gen/grpc/past_values/client"
	genvalues "github.com/crossnokaye/past-values/services/past-values/gen/past_values"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
)

// timeFormat is used to parse times in order to store time as ISO8601 format
const timeFormat = "2006-01-02T15:04:05-07:00"

type (
    Client interface {
        // GetPower will call the Past Value functions "FindControlPointConfigsByName" and "GetValues" to get control point ID's and power data 
        GetPower(ctx context.Context, orgID string, dateRange *gencalc.Period, cpaliasname string, pastValInterval int64, reportInterval string, formula *string, agentname string) (*gencalc.ElectricalReport, error)
    }
    client struct {
        getValues goa.Endpoint
        findControlPointConfigsByName goa.Endpoint
    }
    // ErrNotFound is returned when a facility config is not found.
    ErrPowerReportsNotFound struct{ Err error }
)

func New(conn *grpc.ClientConn) Client {
    c := genvaluesc.NewClient(conn, grpc.WaitForReady(true))
    return &client{
        getValues: c.GetValues(),
        findControlPointConfigsByName: c.FindControlPointConfigsByName(),
    }
}

var unitTypes [3]string = [3]string{model.Kwh, model.Kwh_Min, model.Pulse_count}
type PowerPoint struct {
    Unit string
    Value float64
    StartTime string
    EndTime string
	IntervalType string
}
func newPoint(unit string, value float64) (*PowerPoint) {
    return &PowerPoint{Unit: unit, Value: value}
}

// GetPower will call the Past Value functions "FindControlPointConfigsByName" and "GetValues" to get control point ID's and power data
func (c *client) GetPower(ctx context.Context, orgID string, dateRange *gencalc.Period, cpaliasname string, pastValInterval int64, reportInterval string, formula *string, agentname string) (*gencalc.ElectricalReport, error) {
    var cpIDs []genvalues.UUID
    pointIDs, err := c.getControlPointID(ctx, orgID, agentname, cpaliasname)
    if err != nil {
        return nil, &ErrPowerReportsNotFound{fmt.Errorf("control point id not found for name %s for agent %s with err: %w\n", cpaliasname, agentname, err)}
    }
    for _, p := range pointIDs { //going to be one Point ID
        cpIDs = append(cpIDs, *p.ID)
    }
    p := &genvalues.ValuesQuery{
        OrgID: genvalues.UUID(orgID),
        PointIds: cpIDs,
        Start: dateRange.StartTime,
        End: dateRange.EndTime,
        Interval: pastValInterval,
    }
    res, err := c.getValues(ctx, &p)
    if err != nil {
        return nil, &ErrPowerReportsNotFound{Err: fmt.Errorf("err in getvalues: %w\n", err)}
    }
    analogValues, err := toPower(res)
	var durationType time.Duration
	switch reportInterval {
	case model.Minute:
		durationType = time.Minute * 5
	case model.Hourly:
		durationType = time.Hour
	case model.Daily:
		durationType = time.Hour * 24
	case model.Weekly:
		durationType = time.Hour * 24 * 7
	case model.Monthly:
		durationType = time.Hour * 24 * 29
	} 
	//newFormula, err := utility.ParseFormula(*formula)
    //newFormula.Evaluate()//figure out what this formula would look like
	kwhPoints, err := convertToPower(analogValues, formula, durationType)

    if err != nil {
        return nil, ErrPowerReportsNotFound{Err: fmt.Errorf("err casting getvalues response: %w", err)}
    }
    duration := &gencalc.Period{StartTime: dateRange.StartTime, EndTime: dateRange.EndTime}
    return &gencalc.ElectricalReport{Duration: duration, PowerStamps:  kwhPoints, Interval: reportInterval}, nil

    /*
        Format of newRes once casted to historical values:
        analog points(id and array of "values" or "analog point")->timestamp and values
    {
    "Analog": [
        {
            "ID": "aaa09388-98e4-11ec-b909-0242ac120002",
            "Values": [
                {
                    "Timestamp": "2022-07-20T11:00:00Z",
                    "Value": 8053882
                },
                {
                    "Timestamp": "2022-07-20T11:00:01Z",
                    "Value": 8053882
                },
                {
                    "Timestamp": "2022-07-20T11:00:02Z",
                    "Value": 8053882
                },
                {
                    "Timestamp": "2022-07-20T11:00:03Z",
                    "Value": 8053882
                },
                {
                    "Timestamp": "2022-07-20T11:00:04Z",
                    "Value": 8053883
                },
                …(more values)
            ]
        }
    ],
    "Discrete": [],
    "Structures": []
}
    */
}

//ToPower will cast the response from GetValues and return 1 hour interval reports to match the ones
//returned from the Poller service. It will read the values from the input control point and convert them to Power in KW utilizing the formula
func toPower(r interface{}) ([]*genvalues.AnalogPoint, error) {
    res := r.(*genvalues.HistoricalValues)
    var analogPoints = res.Analog
    var analogForCP = analogPoints[0]
    analogVals := analogForCP.Values //[{timestamp, value}, {timestamp, value} ...]
	return analogVals, nil
}

// getControlPointID will use the past values function getControlPointConfigByName to get the point ID
func (c *client) getControlPointID(ctx context.Context, orgID string, agentName string, pointName string) ([]*genvalues.ControlPointConfig, error) {
    payload := genvalues.PointNameQuery{OrgID: genvalues.UUID(orgID), ClientName: agentName, PointName: pointName}
    fmt.Println(payload)
    res, err := c.findControlPointConfigsByName(ctx, &payload)
    if err != nil {
        return nil, err
    }
    newres, err := toControlPointID(res)
    if err != nil {
        return nil, err
    }
    fmt.Println("control point ID")
    fmt.Println(newres[0].ID)
    return newres, nil
}

// toControlPointID will cast the response from getControlPointConfigByName to a point ID
func toControlPointID(r interface{}) ([]*genvalues.ControlPointConfig, error) {
    res := r.([]*genvalues.ControlPointConfig)
    return res, nil
}
func (err ErrPowerReportsNotFound) Error() string { return err.Err.Error() }


/* 
1. make call in terms of payload start and end time
2. every pulse represents a count which can be converted to a certain value with a formula
## Steps - To calculate kwh @ Oxnard

1. Query for pulse count over a period of time
    * Example: energy_meter_val t1: Friday at 2:00pm, t2: Saturday at 2:00pm 
    * Data: [Value, Timestamp]

2. t1 = is the first value returned (Farthest in the past), t2 = the last value returned (most recent)
    * pulse_count = t2 - t1 = The amount of pulses that occured in your time window
3. obtain t1 and t2 every hour
4. do this for other time intervals
5. 
5. pulse_count * [Formula from the location file] = Amount of energy used in kwh???? 
    * ex: Oxnards factor is 0.6kwh per pulse 
    * t2 - t1 = 187
    * 187 * 0.6 = total energy used between t1 and t


*/
/*
    Converter: input is an array of analog points
        use PowerPoint structs
            ->units(KW, KWh, KWh/min)
            ->scaling(formula to convert pulse count -> Kwh)
            1. convert raw pulses into # of pulses per minute
            2. use formula to convert # of pulses per minute -> KWh
            3. convert Kwh -> KW
        converts raw pulse values -> power
        GOAL: get KW: a rate
        KWh is the substance/energy
        pulse count -> energy(KWh)/min -> take derivative of energy wrt time to get rate of energy consumed->
        1kWh = 60kWmin(energy used by doing 1 KW for an hour)
        formula: convert pulse count from power meters
    
    */
func  convertToPower(analogPoints []*genvalues.AnalogPoint, formula *string, durationtype time.Duration) ([]*gencalc.DataPoint, error) {
    totalPoints := (len(analogPoints) - 1)
	sfinalEndTime := *analogPoints[totalPoints].Timestamp
	sstartTime := *analogPoints[0].Timestamp
	start, _ := time.Parse(timeFormat, sstartTime)
	end, _ := time.Parse(timeFormat, sfinalEndTime)
	var points []*gencalc.DataPoint
	var reportCounter int
	var previousReport = *analogPoints[0]
	//TODO parse formula
	for start.Before(end) {
		if reportCounter == totalPoints { //should not happen
			return points, nil
		}
		if *analogPoints[reportCounter].Value == 0 {
			reportCounter += 1
			continue
		}
		reportTime, _ := time.Parse(timeFormat, *analogPoints[reportCounter].Timestamp)
		if reportTime.Sub(start) >= durationtype {
			power := *previousReport.Value - *analogPoints[reportCounter].Value //TODO: convert this using formula
			point := &gencalc.DataPoint{Time: reportTime.Format(timeFormat), Value: power}
			points = append(points, point)
			previousReport = *analogPoints[reportCounter]
			start = reportTime
		}
		reportCounter += 1
	}
    return points, nil
}
/*
func  convertToPowerr(analogPoints []*genvalues.AnalogPoint, formula *string, times []*gencalc.Period) ([]*gencalc.DataPoint, error) {
    totalPoints := (len(analogPoints) - 1)
	sfinalEndTime := *analogPoints[totalPoints].Timestamp
	sstartTime := *analogPoints[0].Timestamp
	start, _ := time.Parse(timeFormat, sstartTime)
	end, _ := time.Parse(timeFormat, sfinalEndTime)
	var points []*gencalc.DataPoint
	var reportCounter = 0
	var previousReport = *analogPoints[0]
	
    for _, per := range times {
        analogPoint := analogPoints[reportCounter]
        power := 
        power := *previousReport.Value - *analogPoints[reportCounter].Value //TODO: convert this using formula
			point := &gencalc.DataPoint{Time: reportTime.Format(timeFormat), Value: power}
			points = append(points, point)
    }
	for start.Before(end) {
		if reportCounter == totalPoints { //should not happen
			return points, nil
		}
		if *analogPoints[reportCounter].Value == 0 {
			reportCounter += 1
			continue
		}
		reportTime, _ := time.Parse(timeFormat, *analogPoints[reportCounter].Timestamp)
		if reportTime.Sub(start) >= durationtype {
			power := *previousReport.Value - *analogPoints[reportCounter].Value //TODO: convert this using formula
			point := &gencalc.DataPoint{Time: reportTime.Format(timeFormat), Value: power}
			points = append(points, point)
			previousReport = *analogPoints[reportCounter]
			start = reportTime
		}
		reportCounter += 1
	}
    return points, nil
}

*/
