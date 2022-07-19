package power

import (
	"context"
	"fmt"
	goa "goa.design/goa/v3/pkg"
	"google.golang.org/grpc"
	"github.com/google/uuid"
	gencalc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	genvalues "github.com/crossnokaye/past-values/services/past-values/gen/past_values"
	genvaluesc "github.com/crossnokaye/past-values/services/past-values/gen/grpc/past_values/client"
)

var timeFormat = "2006-01-02T15:04:05-07:00"
var dateFormat = "2006-01-02"

type (
	Client interface {
		GetPower(ctx context.Context, orgID string, controlPoint string, durationInterval int64, start string, end string, formula *string) ([]*gencalc.ElectricalReport, error)
	}
	client struct {
		getPower goa.Endpoint
		getControlPointID goa.Endpoint
	}
	
)

func New(conn *grpc.ClientConn) Client {
	c := genvaluesc.NewClient(conn, grpc.WaitForReady(true))
	return &client{
		getPower: c.GetValues(),
		getControlPointID: c.FindControlPointConfigsByName(),
	}
}


func (c *client) GetPower(ctx context.Context, orgID string, controlPoint string, durationInterval int64,
	 start string, end string, formula *string) ([]*gencalc.ElectricalReport, error) {
	var newOrg = genvalues.UUID(orgID)
	var pointIDs []genvalues.UUID
	pointIDs = append(pointIDs, genvalues.UUID(controlPoint))
	//Make call to getControlPointID here
	p := genvalues.ValuesQuery{
		OrgID: newOrg,
		PointIds: pointIDs,
		Start: start,
		End: end,
		Interval: durationInterval,
	}

	res, err := c.getPower(ctx, &p) //res is value *genvalues.HistoricalValues: discrete points, analog points and structures
	
	if err != nil {
		return nil, fmt.Errorf("Error in GetPower: %s\n", err)
	}
	newRes, err := toPower(res)
	//TODO: Roman has to implement errors in his design file so i can implement error handling
	if err != nil {
		return nil, fmt.Errorf("Error in GetPower: %s\n", err)
	}
	return newRes, nil
}

//ToPower will cast the response from GetValues and return 5 minute interval reports to match the ones
//returned from the Poller service. It will read the values from the input control point and convert them to Power in KW utilizing the formula
func toPower(r interface{}) ([]*gencalc.ElectricalReport, error) {
	//TODO implement this function
	res := r.([]*genvalues.HistoricalValues)
	var report *gencalc.ElectricalReport
	var analogPoints = res[1] //Array of Analog Points
	var analog = analogPoints.Analog
	for _, cp := range analog {
		fmt.Printf("control point id is %d\n", cp.ID)
	}
 	//historical values->ArrayOf(Analog points){id, array of analogpoint}->analogpoint: timestamp and value
	//historical values->ArrayOf(devices)->array of control points per device->each control point contains a timestamp and a value
	return report, nil
}

//getControlPointID will obtain the control point id from the below input in order to use the GetPower function with valid input
func getControlPointID(ctx context.Context, orgID string, agentName string, pointName string) (uuid.UUID, error) {
	return nil, nil
}




