package power

import (
	"context"
	"fmt"
	"time"
    "strconv"
    "goa.design/clue/log"
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
        GetPower(ctx context.Context, orgID string, dateRange *gencalc.Period, cpaliasname string, pastValInterval int64, reportInterval time.Duration, formula *string, agentname string) (*gencalc.ElectricalReport, error)
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

// GetPower will call the Past Value functions "FindControlPointConfigsByName" and "GetValues" to get control point ID's and power data
func (c *client) GetPower(ctx context.Context, orgID string, dateRange *gencalc.Period, cpaliasname string, pastValInterval int64, reportInterval time.Duration, formula *string, agentname string) (*gencalc.ElectricalReport, error) {
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
    res, err := c.getValues(ctx, p)
    if err != nil {
        return nil, &ErrPowerReportsNotFound{Err: fmt.Errorf("err in getvalues: %w\n", err)}
    }
    analogValues, err := toPower(res)
    if err != nil {
        return nil, &ErrPowerReportsNotFound{fmt.Errorf("values not found for org: %s for pointID %s with err: %w", orgID, pointID, err)}
    }
	kwhPoints, err := convertToPower(ctx, analogValues, formula, reportInterval)
    if err != nil {
        return nil, ErrPowerReportsNotFound{Err: fmt.Errorf("err converting to KWh: %w", err)}
    }
    duration := &gencalc.Period{StartTime: dateRange.StartTime, EndTime: dateRange.EndTime}
    return &gencalc.ElectricalReport{Duration: duration, PowerStamps:  kwhPoints}, nil
}

//ToPower will cast the response from GetValues and return 1 hour interval reports to match the ones
//returned from the Poller service. It will read the values from the input control point and convert them to Power in KW utilizing the formula
func toPower(r interface{}) ([]*genvalues.AnalogPoint, error) {
    res := r.(*genvalues.GetValuesResult)
    var analogPoints = res.Values.Analog
    if len(analogPoints) != 1 {
        return nil, fmt.Errorf("incorrect analog points returned")
    }
    var analogForCP = analogPoints[0]
    if analogForCP == nil {
        return nil, fmt.Errorf("analog points are null")
    }
    analogVals := analogForCP.Values
    if len(analogVals) == 0 {
        return nil, fmt.Errorf("no analog points")
    }
	return analogVals, nil
}

// getControlPointID will use the past values function getControlPointConfigByName to get the point ID
func (c *client) getControlPointID(ctx context.Context, orgID string, agentName string, pointName string) (genvalues.UUID, error) {
    payload := genvalues.PointNameQuery{OrgID: genvalues.UUID(orgID), ClientName: agentName, PointName: pointName}
    log.Info(ctx, log.KV{K: "payload to getcontrolpoint", V: payload})
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
    if len(values) > 1 || len(values) == 0 {
        return genvalues.UUID(uuid.Nil.String()), fmt.Errorf("more control points returned than input")
    }
    return genvalues.UUID(values[0].ID), nil
}
func (err ErrPowerReportsNotFound) Error() string { return err.Err.Error() }

func  convertToPower(ctx context.Context, analogPoints []*genvalues.AnalogPoint, formula *string, durationtype time.Duration) ([]*gencalc.DataPoint, error) {
	endTime := analogPoints[(len(analogPoints) - 1)].Timestamp
	startTime := analogPoints[0].Timestamp
	start, err := time.Parse(time.RFC3339, startTime)
    if err != nil {
        return nil, err
    }
	end, err := time.Parse(time.RFC3339, endTime)
    if err != nil {
        return nil, err
    }
	var points []*gencalc.DataPoint
	var reportCounter = 0
	var previousReport = *analogPoints[0]
    mult, err := strconv.ParseFloat(*formula, 64)
    if err != nil {
        return nil, err
    }
	for start.Before(end) {
		if reportCounter == len(analogPoints) {
			return points, nil
		}
        analogPoint := analogPoints[reportCounter]
		if analogPoint == nil || analogPoint.Value == 0 {
			reportCounter += 1
			continue
		}
		reportTime, err := time.Parse(time.RFC3339, analogPoint.Timestamp)
        if err != nil {
            return nil, err
        }
		if reportTime.Sub(start) >= durationtype {
			power := (analogPoint.Value - previousReport.Value) * mult
            timeInISO := time.Date(reportTime.Year(), reportTime.Month(), reportTime.Day(),
             reportTime.Hour(), reportTime.Minute(), reportTime.Second(), reportTime.Nanosecond(), reportTime.Location())
             point := &gencalc.DataPoint{Time: timeInISO.Format(timeFormat), Value: power}
            log.Info(ctx, log.KV{K: "power point", V: point})
			points = append(points, point)
			previousReport = *analogPoint
			start = reportTime
		}
		reportCounter += 1
	}
    return points, nil
}





