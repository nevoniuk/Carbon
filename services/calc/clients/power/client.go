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
		GetPower(context.Context, string, []uuid.UUID, int64, string, string) (*gencalc.ElectricalReport, error)
	}
	client struct {
		getPower goa.Endpoint
	}
	
)

func New(conn *grpc.ClientConn) Client {
	c := genvaluesc.NewClient(conn, grpc.WaitForReady(true))
	//method to client
	return &client{
		getPower: c.GetValues(),
	}
}

func (c *client) GetPower(ctx context.Context, orgID string, controlPoints []uuid.UUID, interval int64,
	 start string, end string) (*gencalc.ElectricalReport, error) {
	var cps []genvalues.UUID
	for _,point := range controlPoints {
		newPoint := genvalues.UUID(point.ID())
		cps = append(cps, newPoint)
	}

	p := genvalues.ValuesQuery{
		OrgID: genvalues.UUID(orgID),
		PointIds: cps,
		Start: start,
		End: end,
		Interval: interval,
	}

	res, err := c.getPower(ctx, &p)
	//res is historical values
	//historical values = discrete points, analog points and structures

	if err != nil {
		return nil, fmt.Errorf("Error in GetPower: %s\n", err)
	}
	newRes, err := toPower(res)
	//analyze error
	//wrong ordID
	if err != nil {
		return nil, fmt.Errorf("Error in GetPower: %s\n", err)
	}
	return newRes, nil
}

//ToPower will cast the response from GetValues and return 5 minute interval reports to match the ones
//returned from the Poller service
func toPower(r interface{}) (*gencalc.ElectricalReport, error) {
	//knowing that the client name and agent name were already passed in
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



