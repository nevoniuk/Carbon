// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc gRPC server types
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design

package server

import (
	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	goa "goa.design/goa/v3/pkg"
)

// NewCalculateReportsPayload builds the payload of the "calculate_reports"
// endpoint of the "calc" service from the gRPC request type.
func NewCalculateReportsPayload(message *calcpb.CalculateReportsRequest) *calc.CarbonReport {
	v := &calc.CarbonReport{
		GeneratedRate: message.GeneratedRate,
		DurationType:  message.DurationType,
		Region:        message.Region,
	}
	if message.Duration != nil {
		v.Duration = protobufCalcpbPeriodToCalcPeriod(message.Duration)
	}
	return v
}

// NewProtoCalculateReportsResponse builds the gRPC response type from the
// result of the "calculate_reports" endpoint of the "calc" service.
func NewProtoCalculateReportsResponse(result *calc.TotalReport) *calcpb.CalculateReportsResponse {
	message := &calcpb.CalculateReportsResponse{
		DurationType: result.DurationType,
		Facility:     result.Facility,
	}
	if result.Duration != nil {
		message.Duration = svcCalcPeriodToCalcpbPeriod(result.Duration)
	}
	if result.Point != nil {
		message.Point = make([]*calcpb.DataPoint, len(result.Point))
		for i, val := range result.Point {
			message.Point[i] = &calcpb.DataPoint{
				Time:       val.Time,
				CarbonRate: val.CarbonRate,
			}
		}
	}
	return message
}

// NewGetControlPointsPayload builds the payload of the "get_control_points"
// endpoint of the "calc" service from the gRPC request type.
func NewGetControlPointsPayload(message *calcpb.GetControlPointsRequest) *calc.PastValuesPayload {
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
	return v
}

// NewProtoGetControlPointsResponse builds the gRPC response type from the
// result of the "get_control_points" endpoint of the "calc" service.
func NewProtoGetControlPointsResponse(result []string) *calcpb.GetControlPointsResponse {
	message := &calcpb.GetControlPointsResponse{}
	message.Field = make([]string, len(result))
	for i, val := range result {
		message.Field[i] = val
	}
	return message
}

// NewGetPowerPayload builds the payload of the "get_power" endpoint of the
// "calc" service from the gRPC request type.
func NewGetPowerPayload(message *calcpb.GetPowerRequest) *calc.GetPowerPayload {
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
	return v
}

// NewProtoGetPowerResponse builds the gRPC response type from the result of
// the "get_power" endpoint of the "calc" service.
func NewProtoGetPowerResponse(result *calc.ElectricalReport) *calcpb.GetPowerResponse {
	message := &calcpb.GetPowerResponse{
		Postalcode:   result.Postalcode,
		Facility:     result.Facility,
		Building:     result.Building,
		IntervalType: result.IntervalType,
	}
	if result.Period != nil {
		message.Period = svcCalcPeriodToCalcpbPeriod(result.Period)
	}
	if result.Stamp != nil {
		message.Stamp = make([]*calcpb.PowerStamp, len(result.Stamp))
		for i, val := range result.Stamp {
			message.Stamp[i] = &calcpb.PowerStamp{}
			if val.Time != nil {
				message.Stamp[i].Time = *val.Time
			}
			if val.GenRate != nil {
				message.Stamp[i].GenRate = *val.GenRate
			}
		}
	}
	return message
}

// NewGetEmissionsPayload builds the payload of the "get_emissions" endpoint of
// the "calc" service from the gRPC request type.
func NewGetEmissionsPayload(message *calcpb.GetEmissionsRequest) *calc.EmissionsPayload {
	v := &calc.EmissionsPayload{}
	if message.Interval != "" {
		v.Interval = &message.Interval
	}
	if message.Period != nil {
		v.Period = protobufCalcpbPeriodToCalcPeriod(message.Period)
	}
	return v
}

// NewProtoGetEmissionsResponse builds the gRPC response type from the result
// of the "get_emissions" endpoint of the "calc" service.
func NewProtoGetEmissionsResponse(result *calc.CarbonReport) *calcpb.GetEmissionsResponse {
	message := &calcpb.GetEmissionsResponse{
		GeneratedRate: result.GeneratedRate,
		DurationType:  result.DurationType,
		Region:        result.Region,
	}
	if result.Duration != nil {
		message.Duration = svcCalcPeriodToCalcpbPeriod(result.Duration)
	}
	return message
}

// NewHandleRequestsPayload builds the payload of the "handle_requests"
// endpoint of the "calc" service from the gRPC request type.
func NewHandleRequestsPayload(message *calcpb.HandleRequestsRequest) *calc.RequestPayload {
	v := &calc.RequestPayload{
		Org:      message.Org,
		Building: message.Building,
		Interval: message.Interval,
	}
	if message.Period != nil {
		v.Period = protobufCalcpbPeriodToCalcPeriod(message.Period)
	}
	return v
}

// NewProtoHandleRequestsResponse builds the gRPC response type from the result
// of the "handle_requests" endpoint of the "calc" service.
func NewProtoHandleRequestsResponse() *calcpb.HandleRequestsResponse {
	message := &calcpb.HandleRequestsResponse{}
	return message
}

// NewProtoCarbonreportResponse builds the gRPC response type from the result
// of the "carbonreport" endpoint of the "calc" service.
func NewProtoCarbonreportResponse() *calcpb.CarbonreportResponse {
	message := &calcpb.CarbonreportResponse{}
	return message
}

// ValidateCalculateReportsRequest runs the validations defined on
// CalculateReportsRequest.
func ValidateCalculateReportsRequest(message *calcpb.CalculateReportsRequest) (err error) {
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
func ValidatePeriod(message *calcpb.Period) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("message.startTime", message.StartTime, goa.FormatDateTime))

	err = goa.MergeErrors(err, goa.ValidateFormat("message.endTime", message.EndTime, goa.FormatDateTime))

	return
}

// ValidateGetControlPointsRequest runs the validations defined on
// GetControlPointsRequest.
func ValidateGetControlPointsRequest(message *calcpb.GetControlPointsRequest) (err error) {
	if message.Org != "" {
		err = goa.MergeErrors(err, goa.ValidateFormat("message.org", message.Org, goa.FormatUUID))
	}
	if message.Period != nil {
		if err2 := ValidatePeriod(message.Period); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	if message.Building != "" {
		err = goa.MergeErrors(err, goa.ValidateFormat("message.building", message.Building, goa.FormatUUID))
	}
	if message.Client != "" {
		err = goa.MergeErrors(err, goa.ValidateFormat("message.client", message.Client, goa.FormatUUID))
	}
	return
}

// ValidateGetPowerRequest runs the validations defined on GetPowerRequest.
func ValidateGetPowerRequest(message *calcpb.GetPowerRequest) (err error) {
	if message.Period == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Period", "message"))
	}
	if message.Cps == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("cps", "message"))
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message.org", message.Org, goa.FormatUUID))

	if message.Period != nil {
		if err2 := ValidatePeriod(message.Period); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateGetEmissionsRequest runs the validations defined on
// GetEmissionsRequest.
func ValidateGetEmissionsRequest(message *calcpb.GetEmissionsRequest) (err error) {
	if message.Period != nil {
		if err2 := ValidatePeriod(message.Period); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateHandleRequestsRequest runs the validations defined on
// HandleRequestsRequest.
func ValidateHandleRequestsRequest(message *calcpb.HandleRequestsRequest) (err error) {
	if message.Period == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Period", "message"))
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message.org", message.Org, goa.FormatUUID))

	if message.Period != nil {
		if err2 := ValidatePeriod(message.Period); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message.building", message.Building, goa.FormatUUID))

	return
}

// protobufCalcpbPeriodToCalcPeriod builds a value of type *calc.Period from a
// value of type *calcpb.Period.
func protobufCalcpbPeriodToCalcPeriod(v *calcpb.Period) *calc.Period {
	res := &calc.Period{
		StartTime: v.StartTime,
		EndTime:   v.EndTime,
	}

	return res
}

// svcCalcPeriodToCalcpbPeriod builds a value of type *calcpb.Period from a
// value of type *calc.Period.
func svcCalcPeriodToCalcpbPeriod(v *calc.Period) *calcpb.Period {
	res := &calcpb.Period{
		StartTime: v.StartTime,
		EndTime:   v.EndTime,
	}

	return res
}
