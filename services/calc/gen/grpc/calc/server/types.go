// Code generated by goa v3.7.6, DO NOT EDIT.
//
// calc gRPC server types
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design -o services/calc

package server

import (
	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	goa "goa.design/goa/v3/pkg"
)

// NewHandleRequestsPayload builds the payload of the "handle_requests"
// endpoint of the "calc" service from the gRPC request type.
func NewHandleRequestsPayload(message *calcpb.HandleRequestsRequest) *calc.RequestPayload {
	v := &calc.RequestPayload{
		OrgID:      calc.UUID(message.OrgId),
		AgentID:    message.AgentId,
		FacilityID: message.FacilityId,
		Interval:   message.Interval,
	}
	if message.Duration != nil {
		v.Duration = protobufCalcpbPeriodToCalcPeriod(message.Duration)
	}
	return v
}

// NewProtoHandleRequestsResponse builds the gRPC response type from the result
// of the "handle_requests" endpoint of the "calc" service.
func NewProtoHandleRequestsResponse(result *calc.AllReports) *calcpb.HandleRequestsResponse {
	message := &calcpb.HandleRequestsResponse{}
	if result.CarbonIntensityReports != nil {
		message.CarbonIntensityReports = make([]*calcpb.CarbonReport, len(result.CarbonIntensityReports))
		for i, val := range result.CarbonIntensityReports {
			message.CarbonIntensityReports[i] = &calcpb.CarbonReport{
				GeneratedRate: val.GeneratedRate,
				DurationType:  val.DurationType,
				Region:        val.Region,
			}
			if val.Duration != nil {
				message.CarbonIntensityReports[i].Duration = svcCalcPeriodToCalcpbPeriod(val.Duration)
			}
		}
	}
	if result.PowerReports != nil {
		message.PowerReports = make([]*calcpb.ElectricalReport, len(result.PowerReports))
		for i, val := range result.PowerReports {
			message.PowerReports[i] = &calcpb.ElectricalReport{
				OrgId:         string(val.OrgID),
				AgentId:       val.AgentID,
				GeneratedRate: val.GeneratedRate,
				IntervalType:  val.IntervalType,
				FacilityId:    val.FacilityID,
			}
			if val.Duration != nil {
				message.PowerReports[i].Duration = svcCalcPeriodToCalcpbPeriod(val.Duration)
			}
		}
	}
	if result.TotalEmissionReport != nil {
		message.TotalEmissionReport = svcCalcEmissionsReportToCalcpbEmissionsReport(result.TotalEmissionReport)
	}
	return message
}

// NewProtoGetCarbonReportResponse builds the gRPC response type from the
// result of the "get_carbon_report" endpoint of the "calc" service.
func NewProtoGetCarbonReportResponse() *calcpb.GetCarbonReportResponse {
	message := &calcpb.GetCarbonReportResponse{}
	return message
}

// ValidateHandleRequestsRequest runs the validations defined on
// HandleRequestsRequest.
func ValidateHandleRequestsRequest(message *calcpb.HandleRequestsRequest) (err error) {
	if message.Duration == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Duration", "message"))
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message", string(message.OrgId), goa.FormatUUID))

	if message.Duration != nil {
		if err2 := ValidatePeriod(message.Duration); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	return
}

// ValidateUUID runs the validations defined on UUID.
func ValidateUUID(message string) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("message", message, goa.FormatUUID))

	return
}

// ValidatePeriod runs the validations defined on Period.
func ValidatePeriod(message *calcpb.Period) (err error) {
	err = goa.MergeErrors(err, goa.ValidateFormat("message.StartTime", message.StartTime, goa.FormatDateTime))

	err = goa.MergeErrors(err, goa.ValidateFormat("message.EndTime", message.EndTime, goa.FormatDateTime))

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

// svcCalcEmissionsReportToCalcpbEmissionsReport builds a value of type
// *calcpb.EmissionsReport from a value of type *calc.EmissionsReport.
func svcCalcEmissionsReportToCalcpbEmissionsReport(v *calc.EmissionsReport) *calcpb.EmissionsReport {
	res := &calcpb.EmissionsReport{
		DurationType: v.DurationType,
		OrgId:        string(v.OrgID),
		AgentId:      v.AgentID,
		FacilityId:   v.FacilityID,
	}
	if v.Duration != nil {
		res.Duration = svcCalcPeriodToCalcpbPeriod(v.Duration)
	}
	if v.Points != nil {
		res.Points = make([]*calcpb.DataPoint, len(v.Points))
		for i, val := range v.Points {
			res.Points[i] = &calcpb.DataPoint{
				Time:            val.Time,
				CarbonFootprint: val.CarbonFootprint,
			}
		}
	}

	return res
}

// protobufCalcpbEmissionsReportToCalcEmissionsReport builds a value of type
// *calc.EmissionsReport from a value of type *calcpb.EmissionsReport.
func protobufCalcpbEmissionsReportToCalcEmissionsReport(v *calcpb.EmissionsReport) *calc.EmissionsReport {
	res := &calc.EmissionsReport{
		DurationType: v.DurationType,
		OrgID:        calc.UUID(v.OrgId),
		AgentID:      v.AgentId,
		FacilityID:   v.FacilityId,
	}
	if v.Duration != nil {
		res.Duration = protobufCalcpbPeriodToCalcPeriod(v.Duration)
	}
	if v.Points != nil {
		res.Points = make([]*calc.DataPoint, len(v.Points))
		for i, val := range v.Points {
			res.Points[i] = &calc.DataPoint{
				Time:            val.Time,
				CarbonFootprint: val.CarbonFootprint,
			}
		}
	}

	return res
}
