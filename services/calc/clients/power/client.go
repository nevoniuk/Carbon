package power

import (
	"context"
	"fmt"
	"time"
    "strconv"
	"github.com/crossnokaye/carbon/model"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	genvaluesc "github.com/crossnokaye/past-values/services/past-values/gen/grpc/past_values/client"
	genvalues "github.com/crossnokaye/past-values/services/past-values/gen/past_values"
	"github.com/google/uuid"
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
    pointID, err := c.getControlPointID(ctx, orgID, agentname, cpaliasname)
    if err != nil {
        return nil, &ErrPowerReportsNotFound{fmt.Errorf("control point id not found for name %s for agent %s with err: %w", cpaliasname, agentname, err)}
    }
    p := &genvalues.ValuesQuery {
        OrgID: genvalues.UUID(orgID),
        PointIds: []genvalues.UUID{pointID},
        Start: dateRange.StartTime,
        End: dateRange.EndTime,
        Interval: pastValInterval,
    }
    fmt.Println(p)
    res, err := c.getValues(ctx, p)
    fmt.Println(res)
    if err != nil {
        return nil, &ErrPowerReportsNotFound{Err: fmt.Errorf("err in getvalues: %w\n", err)}
    }
    analogValues, err := toPower(res)
    if err != nil {
        return nil, &ErrPowerReportsNotFound{fmt.Errorf("values not found for org: %s for pointID %s with err: %w", orgID, pointID, err)}
    }
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
    fmt.Println("duration type")
    fmt.Println(durationType)
	kwhPoints, err := convertToPower(analogValues, formula, durationType)
    fmt.Println("length of KWH points")
    fmt.Println(len(kwhPoints))
    if err != nil {
        return nil, ErrPowerReportsNotFound{Err: fmt.Errorf("err converting to KWh: %w", err)}
    }
    duration := &gencalc.Period{StartTime: dateRange.StartTime, EndTime: dateRange.EndTime}
    return &gencalc.ElectricalReport{Duration: duration, PowerStamps:  kwhPoints, Interval: reportInterval}, nil
}

//ToPower will cast the response from GetValues and return 1 hour interval reports to match the ones
//returned from the Poller service. It will read the values from the input control point and convert them to Power in KW utilizing the formula
func toPower(r interface{}) ([]*genvalues.AnalogPoint, error) {
    res := r.(*genvalues.GetValuesResult)
    var analogPoints = res.Values.Analog
    fmt.Println("length of analog points")
    fmt.Println(len(analogPoints))
    if len(analogPoints) != 1 {
        return nil, fmt.Errorf("incorrect analog points returned")
    }
    var analogForCP = analogPoints[0]
    if analogForCP == nil {
        return nil, fmt.Errorf("analog points are null")
    }
    fmt.Println(analogForCP)
    analogVals := analogForCP.Values
    if len(analogVals) == 0 {
        return nil, fmt.Errorf("analog points are null")
    }
    fmt.Println("# of energy pulses")
    fmt.Println(len(analogVals))
    for _, p := range analogVals {
        fmt.Println(p.Timestamp)
        fmt.Println(p.Value) //error its 1
    }
	return analogVals, nil
}

// getControlPointID will use the past values function getControlPointConfigByName to get the point ID
func (c *client) getControlPointID(ctx context.Context, orgID string, agentName string, pointName string) (genvalues.UUID, error) {
    payload := genvalues.PointNameQuery{OrgID: genvalues.UUID(orgID), ClientName: agentName, PointName: pointName}
    fmt.Println("point payload")
    fmt.Println(payload)
    res, err := c.findControlPointConfigsByName(ctx, &payload)
    if err != nil {
        return genvalues.UUID(uuid.Nil.String()), err
    }
    newres, err := toControlPointID(res)
    if err != nil {
        return genvalues.UUID(uuid.Nil.String()), err
    }
    return newres, nil
}
// toControlPointID will cast the response from getControlPointConfigByName to a point ID
func toControlPointID(r interface{}) (genvalues.UUID, error) {
    res := r.(*genvalues.FindControlPointConfigsByNameResult)
    values := res.Values
    fmt.Println("length of point values:")
    fmt.Println(len(values))
    if len(values) > 1 || len(values) == 0 {
        return genvalues.UUID(uuid.Nil.String()), fmt.Errorf("more control points returned than input")
    }
    fmt.Println("point id:")
    fmt.Println(genvalues.UUID(values[0].ID))
    return genvalues.UUID(values[0].ID), nil
}
func (err ErrPowerReportsNotFound) Error() string { return err.Err.Error() }



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
	sfinalEndTime := analogPoints[totalPoints].Timestamp
	sstartTime := analogPoints[0].Timestamp
	start, err := time.Parse(timeFormat, sstartTime)
    if err != nil {
        return nil, err
    }
	end, err := time.Parse(timeFormat, sfinalEndTime)
    if err != nil {
        return nil, err
    }
	var points []*gencalc.DataPoint
	var reportCounter int
	var previousReport = *analogPoints[0]
    //nil formula check in facility config client
    mult, err := strconv.ParseFloat(*formula, 64)
    if err != nil {
        return nil, err
    }
	for start.Before(end) {
		if reportCounter == totalPoints {
			return points, nil
		}
        analogPoint := analogPoints[reportCounter]
        fmt.Println("analog point")
        fmt.Println(analogPoint)
		if analogPoint == nil || analogPoints[reportCounter].Value == 0 {
			reportCounter += 1
			continue
		}

		reportTime, err := time.Parse(timeFormat, analogPoint.Timestamp)
        if err != nil {
            return nil, err
        }
		if reportTime.Sub(start) >= durationtype {
			power := (analogPoint.Value - previousReport.Value) * mult
			point := &gencalc.DataPoint{Time: reportTime.Format(timeFormat), Value: power}
            fmt.Println("power stamp")
            fmt.Println(point)
			points = append(points, point)
			previousReport = *analogPoint
			start = reportTime
		}
		reportCounter += 1
	}
    return points, nil
}



