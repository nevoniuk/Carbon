// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Poller gRPC server types
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design -o
// services/poller

package server

import (
	pollerpb "github.com/crossnokaye/carbon/services/poller/gen/grpc/poller/pb"
	poller "github.com/crossnokaye/carbon/services/poller/gen/poller"
	goa "goa.design/goa/v3/pkg"
)

// NewUpdatePayload builds the payload of the "update" endpoint of the "Poller"
// service from the gRPC request type.
func NewUpdatePayload(message *pollerpb.UpdateRequest) *poller.UpdatePayload {
	v := &poller.UpdatePayload{
		StartTime: message.StartTime,
		EndTime:   message.EndTime,
		Region:    message.Region,
	}
	return v
}

// NewProtoUpdateResponse builds the gRPC response type from the result of the
// "update" endpoint of the "Poller" service.
func NewProtoUpdateResponse() *pollerpb.UpdateResponse {
	message := &pollerpb.UpdateResponse{}
	return message
}

// NewGetEmissionsForRegionPayload builds the payload of the
// "get_emissions_for_region" endpoint of the "Poller" service from the gRPC
// request type.
func NewGetEmissionsForRegionPayload(message *pollerpb.GetEmissionsForRegionRequest) *poller.CarbonPayload {
	v := &poller.CarbonPayload{
		Region: message.Region,
		Start:  message.Start,
		End:    message.End,
	}
	return v
}

// NewProtoGetEmissionsForRegionResponse builds the gRPC response type from the
// result of the "get_emissions_for_region" endpoint of the "Poller" service.
func NewProtoGetEmissionsForRegionResponse(result []*poller.CarbonForecast) *pollerpb.GetEmissionsForRegionResponse {
	message := &pollerpb.GetEmissionsForRegionResponse{}
	message.Field = make([]*pollerpb.CarbonForecast, len(result))
	for i, val := range result {
		message.Field[i] = &pollerpb.CarbonForecast{
			GeneratedRate: val.GeneratedRate,
			MarginalRate:  val.MarginalRate,
			ConsumedRate:  val.ConsumedRate,
			DurationType:  val.DurationType,
			Region:        val.Region,
		}
		if val.Duration != nil {
			message.Field[i].Duration = svcPollerPeriodToPollerpbPeriod(val.Duration)
		}
	}
	return message
}

// ValidateUpdateRequest runs the validations defined on UpdateRequest.
func ValidateUpdateRequest(message *pollerpb.UpdateRequest) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("message.start_time", message.StartTime, goa.FormatDateTime))

	err = goa.MergeErrors(err, goa.ValidateFormat("message.end_time", message.EndTime, goa.FormatDateTime))

	if !(message.Region == "CAISO" || message.Region == "AESO" || message.Region == "BPA" || message.Region == "ERCO" || message.Region == "IESO" || message.Region == "ISONE" || message.Region == "MISO" || message.Region == "NYISO" || message.Region == "NYISO.NYCW" || message.Region == "NYISO.NYLI" || message.Region == "NYISO.NYUP" || message.Region == "PJM" || message.Region == "SPP") {
		err = goa.MergeErrors(err, goa.InvalidEnumValueError("message.region", message.Region, []interface{}{"CAISO", "AESO", "BPA", "ERCO", "IESO", "ISONE", "MISO", "NYISO", "NYISO.NYCW", "NYISO.NYLI", "NYISO.NYUP", "PJM", "SPP"}))
	}
	return
}

// ValidateGetEmissionsForRegionRequest runs the validations defined on
// GetEmissionsForRegionRequest.
func ValidateGetEmissionsForRegionRequest(message *pollerpb.GetEmissionsForRegionRequest) (err error) {
	if !(message.Region == "CAISO" || message.Region == "AESO" || message.Region == "BPA" || message.Region == "ERCO" || message.Region == "IESO" || message.Region == "ISONE" || message.Region == "MISO" || message.Region == "NYISO" || message.Region == "NYISO.NYCW" || message.Region == "NYISO.NYLI" || message.Region == "NYISO.NYUP" || message.Region == "PJM" || message.Region == "SPP") {
		err = goa.MergeErrors(err, goa.InvalidEnumValueError("message.region", message.Region, []interface{}{"CAISO", "AESO", "BPA", "ERCO", "IESO", "ISONE", "MISO", "NYISO", "NYISO.NYCW", "NYISO.NYLI", "NYISO.NYUP", "PJM", "SPP"}))
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message.start", message.Start, goa.FormatDateTime))

	err = goa.MergeErrors(err, goa.ValidateFormat("message.end", message.End, goa.FormatDateTime))

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
