// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC client types
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

package client

import (
	pollerpb "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/pb"
	poller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	goa "goa.design/goa/v3/pkg"
)

// NewProtoCarbonEmissionsRequest builds the gRPC request type from the payload
// of the "carbon_emissions" endpoint of the "Poller" service.
func NewProtoCarbonEmissionsRequest(payload *poller.CarbonPayload) *pollerpb.CarbonEmissionsRequest {
	message := &pollerpb.CarbonEmissionsRequest{}
	if payload.Region != nil {
		message.Region = *payload.Region
	}
	if payload.Start != nil {
		message.Start = *payload.Start
	}
	return message
}

// NewCarbonEmissionsResult builds the result type of the "carbon_emissions"
// endpoint of the "Poller" service from the gRPC response type.
func NewCarbonEmissionsResult(message *pollerpb.CarbonEmissionsResponse) []*poller.CarbonForecast {
	result := make([]*poller.CarbonForecast, len(message.Field))
	for i, val := range message.Field {
		result[i] = &poller.CarbonForecast{
			GeneratedRate:   val.GeneratedRate,
			MarginalRate:    val.MarginalRate,
			ConsumedRate:    val.ConsumedRate,
			GeneratedSource: val.GeneratedSource,
			Region:          val.Region,
		}
		if val.Duration != nil {
			result[i].Duration = protobufPollerpbPeriodToPollerPeriod(val.Duration)
		}
	}
	return result
}

// NewProtoAggregateDataRequest builds the gRPC request type from the payload
// of the "aggregate_data" endpoint of the "Poller" service.
func NewProtoAggregateDataRequest(payload *poller.AggregatePayload) *pollerpb.AggregateDataRequest {
	message := &pollerpb.AggregateDataRequest{}
	if payload.Region != nil {
		message.Region = *payload.Region
	}
	if payload.Periods != nil {
		message.Periods = make([]*pollerpb.Period, len(payload.Periods))
		for i, val := range payload.Periods {
			message.Periods[i] = &pollerpb.Period{
				StartTime: val.StartTime,
				EndTime:   val.EndTime,
			}
		}
	}
	return message
}

// ValidateCarbonEmissionsResponse runs the validations defined on
// CarbonEmissionsResponse.
func ValidateCarbonEmissionsResponse(message *pollerpb.CarbonEmissionsResponse) (err error) {
	for _, e := range message.Field {
		if e != nil {
			if err2 := ValidateCarbonForecast(e); err2 != nil {
				err = goa.MergeErrors(err, err2)
			}
		}
	}
	return
}

// ValidateCarbonForecast runs the validations defined on CarbonForecast.
func ValidateCarbonForecast(message *pollerpb.CarbonForecast) (err error) {
	if message.Duration == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Duration", "message"))
	}
	if message.Duration != nil {
		if err2 := ValidatePeriod(message.Duration); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidatePeriod runs the validations defined on Period.
func ValidatePeriod(message *pollerpb.Period) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("message.startTime", message.StartTime, goa.FormatDateTime))

	err = goa.MergeErrors(err, goa.ValidateFormat("message.endTime", message.EndTime, goa.FormatDateTime))

	return
}

// svcPollerPeriodToPollerpbPeriod builds a value of type *pollerpb.Period from
// a value of type *poller.Period.
func svcPollerPeriodToPollerpbPeriod(v *poller.Period) *pollerpb.Period {
	res := &pollerpb.Period{
		StartTime: v.StartTime,
		EndTime:   v.EndTime,
	}

	return res
}

// protobufPollerpbPeriodToPollerPeriod builds a value of type *poller.Period
// from a value of type *pollerpb.Period.
func protobufPollerpbPeriodToPollerPeriod(v *pollerpb.Period) *poller.Period {
	res := &poller.Period{
		StartTime: v.StartTime,
		EndTime:   v.EndTime,
	}

	return res
}
