package power

import (
	"context"
	"fmt"
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
		GetPower(ctx context.Context, payload *gencalc.PastValPayload) ([]*gencalc.ElectricalReport, error)
	}
	client struct {
		getPower goa.Endpoint
		getControlPointByName goa.Endpoint
	}
	// ErrNotFound is returned when a facility config is not found.
	ErrPowerReportsNotFound struct{ Err error }
)

func New(conn *grpc.ClientConn) Client {
	c := genvaluesc.NewClient(conn, grpc.WaitForReady(true))
	return &client{
		getPower: c.GetValues(),
		getControlPointByName: c.FindControlPointConfigsByName(),
	}
}

// GetPower will call the Past Value functions "FindControlPointConfigsByName" and "GetValues" to get control point ID's and power data
func (c *client) GetPower(ctx context.Context, payload *gencalc.PastValPayload) ([]*gencalc.ElectricalReport, error) {
	pointIDs, err := c.getControlPointID(ctx, payload.OrgID, payload.AgentName, payload.ControlPoint)
	if err != nil {
		return nil, &ErrPowerReportsNotFound{fmt.Errorf("control point id not found for name %s for agent %s with err: %w\n", payload.ControlPoint, payload.AgentName, err)}
	}
	cpIDs := make([]genvalues.UUID, len(pointIDs))
	for i, p := range pointIDs {
    	cpIDs[i] = *p.ID
	}
	p := &genvalues.ValuesQuery{
		OrgID: genvalues.UUID(payload.OrgID),
		PointIds: cpIDs,
		Start: payload.Duration.StartTime,
		End: payload.Duration.EndTime,
		Interval: payload.PastValInterval,
	}
	res, err := c.getPower(ctx, &p)
	if err != nil {
		return nil, &ErrPowerReportsNotFound{Err: fmt.Errorf("err in getvalues: %w\n", err)}
	}
	newRes, err := toPower(res)
	//TODO: Roman has to implement errors in his design file so i can implement error handling
	if err != nil {
		return nil, ErrPowerReportsNotFound{Err: fmt.Errorf("err casting getvalues response: %w", err)}
	}
	return newRes, nil
}

//ToPower will cast the response from GetValues and return 5 minute interval reports to match the ones
//returned from the Poller service. It will read the values from the input control point and convert them to Power in KW utilizing the formula
func toPower(r interface{}) ([]*gencalc.ElectricalReport, error) {
	//TODO implement this function
	res := r.([]*genvalues.HistoricalValues)
	var reports []*gencalc.ElectricalReport
	var analogPoints = res[1] //Array of Analog Points
	var analog = analogPoints.Analog
	for _, cp := range analog {
		fmt.Printf("control point id is %d\n", cp.ID)
	}
 	//historical values->ArrayOf(Analog points){id, array of analogpoint}->analogpoint: timestamp and value
	//historical values->ArrayOf(devices)->array of control points per device->each control point contains a timestamp and a value
	return reports, nil
}

// getControlPointID will use the past values function getControlPointConfigByName to get the point ID
func (c *client) getControlPointID(ctx context.Context, orgID string, agentName string, pointName string) ([]*genvalues.ControlPointConfig, error) {
	payload := genvalues.PointNameQuery{OrgID: genvalues.UUID(orgID), ClientName: agentName, PointName: pointName}
	res, err := c.getControlPointByName(ctx, &payload)
	if err != nil {
		return nil, err
	}
	newres, err := toControlPointID(res)
	if err != nil {
		return nil, err
	}
	return newres, nil
}

// toControlPointID will cast the response from getControlPointConfigByName to a point ID
func toControlPointID(r interface{}) ([]*genvalues.ControlPointConfig, error) {
	res := r.([]*genvalues.ControlPointConfig)
	//TODO implement
	return res, nil
}
func (err ErrPowerReportsNotFound) Error() string { return err.Err.Error() }




