// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc gRPC client CLI support package
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design

package client

import (
	"encoding/json"
	"fmt"

	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
)

// BuildHandleRequestsPayload builds the payload for the calc handle_requests
// endpoint from CLI flags.
func BuildHandleRequestsPayload(calcHandleRequestsMessage string) (*calc.RequestPayload, error) {
	var err error
	var message calcpb.HandleRequestsRequest
	{
		if calcHandleRequestsMessage != "" {
			err = json.Unmarshal([]byte(calcHandleRequestsMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"Period\": {\n         \"endTime\": \"2020-01-01T00:00:00Z\",\n         \"startTime\": \"2020-01-01T00:00:00Z\"\n      },\n      \"building\": \"Facere reiciendis.\",\n      \"interval\": \"hours, days, weeks, months, years\",\n      \"org\": \"Facere reiciendis.\"\n   }'")
			}
		}
	}
	v := &calc.RequestPayload{
		Org:      calc.UUID(message.Org),
		Building: calc.UUID(message.Building),
		Interval: message.Interval,
	}
	if message.Period != nil {
		v.Period = protobufCalcpbPeriodToCalcPeriod(message.Period)
	}

	return v, nil
}
