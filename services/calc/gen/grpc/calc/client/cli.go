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

// BuildCalculateReportsPayload builds the payload for the calc
// calculate_reports endpoint from CLI flags.
func BuildCalculateReportsPayload(calcCalculateReportsMessage string) (*calc.CarbonReport, error) {
	var err error
	var message calcpb.CalculateReportsRequest
	{
		if calcCalculateReportsMessage != "" {
			err = json.Unmarshal([]byte(calcCalculateReportsMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"Duration\": {\n         \"endTime\": \"2020-01-01T00:00:00Z\",\n         \"startTime\": \"2020-01-01T00:00:00Z\"\n      },\n      \"duration_type\": \"Quis sed.\",\n      \"generated_rate\": 37.8267,\n      \"region\": \"MISO, ISO...\"\n   }'")
			}
		}
	}
	v := &calc.CarbonReport{
		GeneratedRate: message.GeneratedRate,
		DurationType:  message.DurationType,
		Region:        message.Region,
	}
	if message.Duration != nil {
		v.Duration = protobufCalcpbPeriodToCalcPeriod(message.Duration)
	}

	return v, nil
}

// BuildGetControlPointsPayload builds the payload for the calc
// get_control_points endpoint from CLI flags.
func BuildGetControlPointsPayload(calcGetControlPointsMessage string) (*calc.PastValuesPayload, error) {
	var err error
	var message calcpb.GetControlPointsRequest
	{
		if calcGetControlPointsMessage != "" {
			err = json.Unmarshal([]byte(calcGetControlPointsMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"Period\": {\n         \"endTime\": \"2020-01-01T00:00:00Z\",\n         \"startTime\": \"2020-01-01T00:00:00Z\"\n      },\n      \"building\": \"1265498D-5A84-134A-1C7A-ED5B4B92788E\",\n      \"client\": \"7D80331A-7620-D09D-7CCB-2EF87B797732\",\n      \"org\": \"5E3B665E-1239-9C12-9643-FFC1E6C04697\"\n   }'")
			}
		}
	}
	v := &calc.PastValuesPayload{}
	if message.Org != "" {
		v.Org = &message.Org
	}
	if message.Building != "" {
		v.Building = &message.Building
	}
	if message.Client != "" {
		v.Client = &message.Client
	}
	if message.Period != nil {
		v.Period = protobufCalcpbPeriodToCalcPeriod(message.Period)
	}

	return v, nil
}

// BuildGetPowerPayload builds the payload for the calc get_power endpoint from
// CLI flags.
func BuildGetPowerPayload(calcGetPowerMessage string) (*calc.GetPowerPayload, error) {
	var err error
	var message calcpb.GetPowerRequest
	{
		if calcGetPowerMessage != "" {
			err = json.Unmarshal([]byte(calcGetPowerMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"Period\": {\n         \"endTime\": \"2020-01-01T00:00:00Z\",\n         \"startTime\": \"2020-01-01T00:00:00Z\"\n      },\n      \"cps\": [\n         \"Labore voluptates sed voluptatibus.\",\n         \"Sed provident omnis quisquam aliquam.\",\n         \"Commodi itaque.\",\n         \"At culpa et et.\"\n      ],\n      \"interval\": 7061356915507293268,\n      \"org\": \"76FB876C-96AC-91E7-BD21-B0C2988DDF65\"\n   }'")
			}
		}
	}
	v := &calc.GetPowerPayload{
		Org:      message.Org,
		Interval: message.Interval,
	}
	if message.Period != nil {
		v.Period = protobufCalcpbPeriodToCalcPeriod(message.Period)
	}
	if message.Cps != nil {
		v.Cps = make([]string, len(message.Cps))
		for i, val := range message.Cps {
			v.Cps[i] = val
		}
	}

	return v, nil
}

// BuildGetEmissionsPayload builds the payload for the calc get_emissions
// endpoint from CLI flags.
func BuildGetEmissionsPayload(calcGetEmissionsMessage string) (*calc.EmissionsPayload, error) {
	var err error
	var message calcpb.GetEmissionsRequest
	{
		if calcGetEmissionsMessage != "" {
			err = json.Unmarshal([]byte(calcGetEmissionsMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"Period\": {\n         \"endTime\": \"2020-01-01T00:00:00Z\",\n         \"startTime\": \"2020-01-01T00:00:00Z\"\n      },\n      \"interval\": \"hours, days, weeks, months, years\"\n   }'")
			}
		}
	}
	v := &calc.EmissionsPayload{}
	if message.Interval != "" {
		v.Interval = &message.Interval
	}
	if message.Period != nil {
		v.Period = protobufCalcpbPeriodToCalcPeriod(message.Period)
	}

	return v, nil
}

// BuildHandleRequestsPayload builds the payload for the calc handle_requests
// endpoint from CLI flags.
func BuildHandleRequestsPayload(calcHandleRequestsMessage string) (*calc.RequestPayload, error) {
	var err error
	var message calcpb.HandleRequestsRequest
	{
		if calcHandleRequestsMessage != "" {
			err = json.Unmarshal([]byte(calcHandleRequestsMessage), &message)
			if err != nil {
				return nil, fmt.Errorf("invalid JSON for message, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"Period\": {\n         \"endTime\": \"2020-01-01T00:00:00Z\",\n         \"startTime\": \"2020-01-01T00:00:00Z\"\n      },\n      \"building\": \"4CCDE767-7648-444F-D09F-4B4FFE4EB36B\",\n      \"interval\": \"hours, days, weeks, months, years\",\n      \"org\": \"A129B534-C1FC-F09D-BF29-3DA5781E0ECB\"\n   }'")
			}
		}
	}
	v := &calc.RequestPayload{
		Org:      message.Org,
		Building: message.Building,
		Interval: message.Interval,
	}
	if message.Period != nil {
		v.Period = protobufCalcpbPeriodToCalcPeriod(message.Period)
	}

	return v, nil
}
