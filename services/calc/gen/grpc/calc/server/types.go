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

// NewHistoricalCarbonEmissionsPayload builds the payload of the
// "Historical_Carbon_Emissions" endpoint of the "calc" service from the gRPC
// request type.
func NewHistoricalCarbonEmissionsPayload(message *calcpb.HistoricalCarbonEmissionsRequest) *calc.RequestPayload {
	v := &calc.RequestPayload{
		OrgID:      calc.UUID(message.OrgId),
		FacilityID: calc.UUID(message.FacilityId),
	}
	if message.LocationId != "" {
		locationIDptr := calc.UUID(message.LocationId)
		v.LocationID = &locationIDptr
	}
	if message.Duration != nil {
		v.Duration = protobufCalcpbPeriodToCalcPeriod(message.Duration)
	}
	if message.Interval != nil {
		v.Interval = protobufCalcpbIntervalTypeToCalcIntervalType(message.Interval)
	}
	return v
}

// NewProtoHistoricalCarbonEmissionsResponse builds the gRPC response type from
// the result of the "Historical_Carbon_Emissions" endpoint of the "calc"
// service.
func NewProtoHistoricalCarbonEmissionsResponse(result *calc.AllReports) *calcpb.HistoricalCarbonEmissionsResponse {
	message := &calcpb.HistoricalCarbonEmissionsResponse{}
	if result.CarbonIntensityReports != nil {
		message.CarbonIntensityReports = make([]*calcpb.CarbonReport, len(result.CarbonIntensityReports))
		for i, val := range result.CarbonIntensityReports {
			message.CarbonIntensityReports[i] = &calcpb.CarbonReport{
				GeneratedRate: val.GeneratedRate,
			}
			if val.Duration != nil {
				message.CarbonIntensityReports[i].Duration = svcCalcPeriodToCalcpbPeriod(val.Duration)
			}
			if val.Interval != nil {
				message.CarbonIntensityReports[i].Interval = svcCalcIntervalTypeToCalcpbIntervalType(val.Interval)
			}
			if val.Region != nil {
				message.CarbonIntensityReports[i].Region = svcCalcRegionNameToCalcpbRegionName(val.Region)
			}
		}
	}
	if result.PowerReports != nil {
		message.PowerReports = make([]*calcpb.ElectricalReport, len(result.PowerReports))
		for i, val := range result.PowerReports {
			message.PowerReports[i] = &calcpb.ElectricalReport{
				Power: val.Power,
			}
			if val.Duration != nil {
				message.PowerReports[i].Duration = svcCalcPeriodToCalcpbPeriod(val.Duration)
			}
			if val.Interval != nil {
				message.PowerReports[i].Interval = svcCalcIntervalTypeToCalcpbIntervalType(val.Interval)
			}
		}
	}
	if result.TotalEmissionReport != nil {
		message.TotalEmissionReport = svcCalcEmissionsReportToCalcpbEmissionsReport(result.TotalEmissionReport)
	}
	return message
}

// ValidateHistoricalCarbonEmissionsRequest runs the validations defined on
// HistoricalCarbonEmissionsRequest.
func ValidateHistoricalCarbonEmissionsRequest(message *calcpb.HistoricalCarbonEmissionsRequest) (err error) {
	if message.Duration == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Duration", "message"))
	}
	if message.Interval == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("Interval", "message"))
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message", string(message.OrgId), goa.FormatUUID))

	if message.Duration != nil {
		if err2 := ValidatePeriod(message.Duration); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message", string(message.FacilityId), goa.FormatUUID))

	if message.Interval != nil {
		if err2 := ValidateIntervalType(message.Interval); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message", string(message.LocationId), goa.FormatUUID))

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

// ValidateIntervalType runs the validations defined on IntervalType.
func ValidateIntervalType(message *calcpb.IntervalType) (err error) {
	if message.Kind != "" {
		if !(message.Kind == "minute" || message.Kind == "hourly" || message.Kind == "daily" || message.Kind == "weekly" || message.Kind == "monthly") {
			err = goa.MergeErrors(err, goa.InvalidEnumValueError("message.Kind", message.Kind, []interface{}{"minute", "hourly", "daily", "weekly", "monthly"}))
		}
	}
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

// protobufCalcpbIntervalTypeToCalcIntervalType builds a value of type
// *calc.IntervalType from a value of type *calcpb.IntervalType.
func protobufCalcpbIntervalTypeToCalcIntervalType(v *calcpb.IntervalType) *calc.IntervalType {
	res := &calc.IntervalType{}
	if v.Kind != "" {
		res.Kind = &v.Kind
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

// svcCalcIntervalTypeToCalcpbIntervalType builds a value of type
// *calcpb.IntervalType from a value of type *calc.IntervalType.
func svcCalcIntervalTypeToCalcpbIntervalType(v *calc.IntervalType) *calcpb.IntervalType {
	res := &calcpb.IntervalType{}
	if v.Kind != nil {
		res.Kind = *v.Kind
	}

	return res
}

// svcCalcRegionNameToCalcpbRegionName builds a value of type
// *calcpb.RegionName from a value of type *calc.RegionName.
func svcCalcRegionNameToCalcpbRegionName(v *calc.RegionName) *calcpb.RegionName {
	res := &calcpb.RegionName{}
	if v.Region != nil {
		res.Region = *v.Region
	}

	return res
}

// svcCalcEmissionsReportToCalcpbEmissionsReport builds a value of type
// *calcpb.EmissionsReport from a value of type *calc.EmissionsReport.
func svcCalcEmissionsReportToCalcpbEmissionsReport(v *calc.EmissionsReport) *calcpb.EmissionsReport {
	res := &calcpb.EmissionsReport{
		OrgId:      string(v.OrgID),
		FacilityId: string(v.FacilityID),
		LocationId: string(v.LocationID),
	}
	if v.Duration != nil {
		res.Duration = svcCalcPeriodToCalcpbPeriod(v.Duration)
	}
	if v.Interval != nil {
		res.Interval = svcCalcIntervalTypeToCalcpbIntervalType(v.Interval)
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
	if v.Region != nil {
		res.Region = svcCalcRegionNameToCalcpbRegionName(v.Region)
	}

	return res
}

// protobufCalcpbRegionNameToCalcRegionName builds a value of type
// *calc.RegionName from a value of type *calcpb.RegionName.
func protobufCalcpbRegionNameToCalcRegionName(v *calcpb.RegionName) *calc.RegionName {
	res := &calc.RegionName{}
	if v.Region != "" {
		res.Region = &v.Region
	}

	return res
}

// protobufCalcpbEmissionsReportToCalcEmissionsReport builds a value of type
// *calc.EmissionsReport from a value of type *calcpb.EmissionsReport.
func protobufCalcpbEmissionsReportToCalcEmissionsReport(v *calcpb.EmissionsReport) *calc.EmissionsReport {
	res := &calc.EmissionsReport{
		OrgID:      calc.UUID(v.OrgId),
		FacilityID: calc.UUID(v.FacilityId),
		LocationID: calc.UUID(v.LocationId),
	}
	if v.Duration != nil {
		res.Duration = protobufCalcpbPeriodToCalcPeriod(v.Duration)
	}
	if v.Interval != nil {
		res.Interval = protobufCalcpbIntervalTypeToCalcIntervalType(v.Interval)
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
	if v.Region != nil {
		res.Region = protobufCalcpbRegionNameToCalcRegionName(v.Region)
	}

	return res
}
