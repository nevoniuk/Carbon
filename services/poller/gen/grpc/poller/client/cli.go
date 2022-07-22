// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC client CLI support package
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design -o
// services/poller

package client

import (
	"encoding/json"
	"fmt"

	pollerpb "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/pb"
	poller "github.com/crossnokaye/carbon/services/poller/gen/poller"
)

// BuildGetEmissionsForRegionPayload builds the payload for the Poller
// get_emissions_for_region endpoint from CLI flags.
func BuildGetEmissionsForRegionPayload(pollerGetEmissionsForRegionMessage string) (*poller.CarbonPayload, error) {
	var err error
	var message pollerpb.GetEmissionsForRegionRequest
	{
		if pollerGetEmissionsForRegionMessage != "" {
			err = json.Unmarshal([]byte(pollerGetEmissionsForRegionMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"end\": \"2020-01-01T00:00:00Z\",\n      \"region\": \"Exercitationem saepe aut sit inventore itaque est.\",\n      \"start\": \"2020-01-01T00:00:00Z\"\n   }'")
			}
		}
	}
	v := &poller.CarbonPayload{}
	if message.Region != "" {
		v.Region = &message.Region
	}
	if message.Start != "" {
		v.Start = &message.Start
	}
	if message.End != "" {
		v.End = &message.End
	}

	return v, nil
}
