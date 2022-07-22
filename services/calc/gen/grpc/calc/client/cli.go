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

// BuildHistoricalCarbonEmissionsPayload builds the payload for the calc
// historical_carbon_emissions endpoint from CLI flags.
func BuildHistoricalCarbonEmissionsPayload(calcHistoricalCarbonEmissionsMessage string) (*calc.RequestPayload, error) {
	var err error
	var message calcpb.HistoricalCarbonEmissionsRequest
	{
		if calcHistoricalCarbonEmissionsMessage != "" {
			err = json.Unmarshal([]byte(calcHistoricalCarbonEmissionsMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"Duration\": {\n         \"EndTime\": \"2020-01-01T00:00:00Z\",\n         \"StartTime\": \"2020-01-01T00:00:00Z\"\n      },\n      \"FacilityID\": \"Facere reiciendis.\",\n      \"Interval\": \"daily\",\n      \"LocationID\": \"Facere reiciendis.\",\n      \"OrgID\": \"Facere reiciendis.\"\n   }'")
			}
		}
	}
	v := &calc.RequestPayload{
		OrgID:      calc.UUID(message.OrgId),
		FacilityID: calc.UUID(message.FacilityId),
		Interval:   message.Interval,
		LocationID: calc.UUID(message.LocationId),
	}
	if message.Duration != nil {
		v.Duration = protobufCalcpbPeriodToCalcPeriod(message.Duration)
	}

	return v, nil
}
