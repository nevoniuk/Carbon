// Code generated by goa v3.7.6, DO NOT EDIT.
//
// Calc gRPC server types
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design -o services/calc

package server

import (
	calc "github.com/crossnokaye/carbon/services/calc/gen/calc"
	calcpb "github.com/crossnokaye/carbon/services/calc/gen/grpc/calc/pb"
	goa "goa.design/goa/v3/pkg"
)

// NewHistoricalCarbonEmissionsPayload builds the payload of the
// "historical_carbon_emissions" endpoint of the "Calc" service from the gRPC
// request type.
func NewHistoricalCarbonEmissionsPayload(message *calcpb.HistoricalCarbonEmissionsRequest) *calc.RequestPayload {
	v := &calc.RequestPayload{
		OrgID:      calc.UUID(message.OrgId),
		FacilityID: calc.UUID(message.FacilityId),
		Interval:   message.Interval,
		LocationID: calc.UUID(message.LocationId),
	}
	if message.Duration != nil {
		v.Duration = protobufCalcpbPeriodToCalcPeriod(message.Duration)
	}
	return v
}

// NewProtoHistoricalCarbonEmissionsResponse builds the gRPC response type from
// the result of the "historical_carbon_emissions" endpoint of the "Calc"
// service.
func NewProtoHistoricalCarbonEmissionsResponse(result *calc.AllReports) *calcpb.HistoricalCarbonEmissionsResponse {
	message := &calcpb.HistoricalCarbonEmissionsResponse{}
	if result.TotalEmissionReport != nil {
		message.TotalEmissionReport = svcCalcEmissionsReportToCalcpbEmissionsReport(result.TotalEmissionReport)
	}
	if result.CarbonIntensityReports != nil {
		message.CarbonIntensityReports = svcCalcCarbonReportToCalcpbCarbonReport(result.CarbonIntensityReports)
	}
	if result.PowerReports != nil {
		message.PowerReports = svcCalcElectricalReportToCalcpbElectricalReport(result.PowerReports)
	}
	return message
}

// ValidateHistoricalCarbonEmissionsRequest runs the validations defined on
// HistoricalCarbonEmissionsRequest.
func ValidateHistoricalCarbonEmissionsRequest(message *calcpb.HistoricalCarbonEmissionsRequest) (err error) {
	if message.Duration == nil {
		err = goa.MergeErrors(err, goa.MissingFieldError("duration", "message"))
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message", string(message.OrgId), goa.FormatUUID))

	if message.Duration != nil {
		if err2 := ValidatePeriod(message.Duration); err2 != nil {
			err = goa.MergeErrors(err, err2)
		}
	}
	err = goa.MergeErrors(err, goa.ValidateFormat("message", string(message.FacilityId), goa.FormatUUID))

	if !(message.Interval == "minute" || message.Interval == "hourly" || message.Interval == "daily" || message.Interval == "weekly" || message.Interval == "monthly") {
		err = goa.MergeErrors(err, goa.InvalidEnumValueError("message.interval", message.Interval, []interface{}{"minute", "hourly", "daily", "weekly", "monthly"}))
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
	err = goa.MergeErrors(err, goa.ValidateFormat("message.start_time", message.StartTime, goa.FormatDateTime))

	err = goa.MergeErrors(err, goa.ValidateFormat("message.end_time", message.EndTime, goa.FormatDateTime))

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
		Interval:   v.Interval,
		OrgId:      string(v.OrgID),
		FacilityId: string(v.FacilityID),
		LocationId: string(v.LocationID),
		Region:     v.Region,
	}
	if v.Duration != nil {
		res.Duration = svcCalcPeriodToCalcpbPeriod(v.Duration)
	}
	if v.Points != nil {
		res.Points = make([]*calcpb.DataPoint, len(v.Points))
		for i, val := range v.Points {
			res.Points[i] = &calcpb.DataPoint{
				Time:  val.Time,
				Value: val.Value,
			}
		}
	}

	return res
}

// svcCalcCarbonReportToCalcpbCarbonReport builds a value of type
// *calcpb.CarbonReport from a value of type *calc.CarbonReport.
func svcCalcCarbonReportToCalcpbCarbonReport(v *calc.CarbonReport) *calcpb.CarbonReport {
	if v == nil {
		return nil
	}
	res := &calcpb.CarbonReport{
		Interval: v.Interval,
		Region:   v.Region,
	}
	if v.IntensityPoints != nil {
		res.IntensityPoints = make([]*calcpb.DataPoint, len(v.IntensityPoints))
		for i, val := range v.IntensityPoints {
			res.IntensityPoints[i] = &calcpb.DataPoint{
				Time:  val.Time,
				Value: val.Value,
			}
		}
	}
	if v.Duration != nil {
		res.Duration = svcCalcPeriodToCalcpbPeriod(v.Duration)
	}

	return res
}

// svcCalcElectricalReportToCalcpbElectricalReport builds a value of type
// *calcpb.ElectricalReport from a value of type *calc.ElectricalReport.
func svcCalcElectricalReportToCalcpbElectricalReport(v *calc.ElectricalReport) *calcpb.ElectricalReport {
	if v == nil {
		return nil
	}
	res := &calcpb.ElectricalReport{
		Interval: v.Interval,
	}
	if v.Duration != nil {
		res.Duration = svcCalcPeriodToCalcpbPeriod(v.Duration)
	}
	if v.PowerStamps != nil {
		res.PowerStamps = make([]*calcpb.DataPoint, len(v.PowerStamps))
		for i, val := range v.PowerStamps {
			res.PowerStamps[i] = &calcpb.DataPoint{
				Time:  val.Time,
				Value: val.Value,
			}
		}
	}

	return res
}

// protobufCalcpbEmissionsReportToCalcEmissionsReport builds a value of type
// *calc.EmissionsReport from a value of type *calcpb.EmissionsReport.
func protobufCalcpbEmissionsReportToCalcEmissionsReport(v *calcpb.EmissionsReport) *calc.EmissionsReport {
	res := &calc.EmissionsReport{
		Interval:   v.Interval,
		OrgID:      calc.UUID(v.OrgId),
		FacilityID: calc.UUID(v.FacilityId),
		LocationID: calc.UUID(v.LocationId),
		Region:     v.Region,
	}
	if v.Duration != nil {
		res.Duration = protobufCalcpbPeriodToCalcPeriod(v.Duration)
	}
	if v.Points != nil {
		res.Points = make([]*calc.DataPoint, len(v.Points))
		for i, val := range v.Points {
			res.Points[i] = &calc.DataPoint{
				Time:  val.Time,
				Value: val.Value,
			}
		}
	}

	return res
}

// protobufCalcpbCarbonReportToCalcCarbonReport builds a value of type
// *calc.CarbonReport from a value of type *calcpb.CarbonReport.
func protobufCalcpbCarbonReportToCalcCarbonReport(v *calcpb.CarbonReport) *calc.CarbonReport {
	if v == nil {
		return nil
	}
	res := &calc.CarbonReport{
		Interval: v.Interval,
		Region:   v.Region,
	}
	if v.IntensityPoints != nil {
		res.IntensityPoints = make([]*calc.DataPoint, len(v.IntensityPoints))
		for i, val := range v.IntensityPoints {
			res.IntensityPoints[i] = &calc.DataPoint{
				Time:  val.Time,
				Value: val.Value,
			}
		}
	}
	if v.Duration != nil {
		res.Duration = protobufCalcpbPeriodToCalcPeriod(v.Duration)
	}

	return res
}

// protobufCalcpbElectricalReportToCalcElectricalReport builds a value of type
// *calc.ElectricalReport from a value of type *calcpb.ElectricalReport.
func protobufCalcpbElectricalReportToCalcElectricalReport(v *calcpb.ElectricalReport) *calc.ElectricalReport {
	if v == nil {
		return nil
	}
	res := &calc.ElectricalReport{
		Interval: v.Interval,
	}
	if v.Duration != nil {
		res.Duration = protobufCalcpbPeriodToCalcPeriod(v.Duration)
	}
	if v.PowerStamps != nil {
		res.PowerStamps = make([]*calc.DataPoint, len(v.PowerStamps))
		for i, val := range v.PowerStamps {
			res.PowerStamps[i] = &calc.DataPoint{
				Time:  val.Time,
				Value: val.Value,
			}
		}
	}

	return res
}
