// Code generated with goa v3.7.6, DO NOT EDIT.
//
// Calc protocol buffer definition
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/calc/design -o services/calc

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.21.4
// source: goadesign_goagen_calc.proto

package calcpb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type HistoricalCarbonEmissionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// OrgID
	OrgId string `protobuf:"bytes,1,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	// Duration
	Duration *Period `protobuf:"bytes,2,opt,name=duration,proto3" json:"duration,omitempty"`
	// FacilityID
	FacilityId string `protobuf:"bytes,3,opt,name=facility_id,json=facilityId,proto3" json:"facility_id,omitempty"`
	Interval   string `protobuf:"bytes,4,opt,name=interval,proto3" json:"interval,omitempty"`
	// LocationID
	LocationId string `protobuf:"bytes,5,opt,name=location_id,json=locationId,proto3" json:"location_id,omitempty"`
}

func (x *HistoricalCarbonEmissionsRequest) Reset() {
	*x = HistoricalCarbonEmissionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_calc_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistoricalCarbonEmissionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoricalCarbonEmissionsRequest) ProtoMessage() {}

func (x *HistoricalCarbonEmissionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_calc_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoricalCarbonEmissionsRequest.ProtoReflect.Descriptor instead.
func (*HistoricalCarbonEmissionsRequest) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_calc_proto_rawDescGZIP(), []int{0}
}

func (x *HistoricalCarbonEmissionsRequest) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

func (x *HistoricalCarbonEmissionsRequest) GetDuration() *Period {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *HistoricalCarbonEmissionsRequest) GetFacilityId() string {
	if x != nil {
		return x.FacilityId
	}
	return ""
}

func (x *HistoricalCarbonEmissionsRequest) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

func (x *HistoricalCarbonEmissionsRequest) GetLocationId() string {
	if x != nil {
		return x.LocationId
	}
	return ""
}

// Period of time from start to end for any report type
type Period struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Start time
	StartTime string `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// End time
	EndTime string `protobuf:"bytes,2,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
}

func (x *Period) Reset() {
	*x = Period{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_calc_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Period) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Period) ProtoMessage() {}

func (x *Period) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_calc_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Period.ProtoReflect.Descriptor instead.
func (*Period) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_calc_proto_rawDescGZIP(), []int{1}
}

func (x *Period) GetStartTime() string {
	if x != nil {
		return x.StartTime
	}
	return ""
}

func (x *Period) GetEndTime() string {
	if x != nil {
		return x.EndTime
	}
	return ""
}

type HistoricalCarbonEmissionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// TotalEmissionReport
	TotalEmissionReport *EmissionsReport `protobuf:"bytes,1,opt,name=total_emission_report,json=totalEmissionReport,proto3" json:"total_emission_report,omitempty"`
	// CarbonIntensityReports
	CarbonIntensityReports *CarbonReport `protobuf:"bytes,2,opt,name=carbon_intensity_reports,json=carbonIntensityReports,proto3" json:"carbon_intensity_reports,omitempty"`
	// PowerReports
	PowerReports *ElectricalReport `protobuf:"bytes,3,opt,name=power_reports,json=powerReports,proto3" json:"power_reports,omitempty"`
}

func (x *HistoricalCarbonEmissionsResponse) Reset() {
	*x = HistoricalCarbonEmissionsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_calc_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HistoricalCarbonEmissionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoricalCarbonEmissionsResponse) ProtoMessage() {}

func (x *HistoricalCarbonEmissionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_calc_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoricalCarbonEmissionsResponse.ProtoReflect.Descriptor instead.
func (*HistoricalCarbonEmissionsResponse) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_calc_proto_rawDescGZIP(), []int{2}
}

func (x *HistoricalCarbonEmissionsResponse) GetTotalEmissionReport() *EmissionsReport {
	if x != nil {
		return x.TotalEmissionReport
	}
	return nil
}

func (x *HistoricalCarbonEmissionsResponse) GetCarbonIntensityReports() *CarbonReport {
	if x != nil {
		return x.CarbonIntensityReports
	}
	return nil
}

func (x *HistoricalCarbonEmissionsResponse) GetPowerReports() *ElectricalReport {
	if x != nil {
		return x.PowerReports
	}
	return nil
}

// Carbon/Energy Generation Report
type EmissionsReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Duration
	Duration *Period `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	Interval string  `protobuf:"bytes,2,opt,name=interval,proto3" json:"interval,omitempty"`
	// Points
	Points []*DataPoint `protobuf:"bytes,3,rep,name=points,proto3" json:"points,omitempty"`
	// OrgID
	OrgId string `protobuf:"bytes,4,opt,name=org_id,json=orgId,proto3" json:"org_id,omitempty"`
	// FacilityID
	FacilityId string `protobuf:"bytes,5,opt,name=facility_id,json=facilityId,proto3" json:"facility_id,omitempty"`
	// LocationID
	LocationId string `protobuf:"bytes,6,opt,name=location_id,json=locationId,proto3" json:"location_id,omitempty"`
	Region     string `protobuf:"bytes,7,opt,name=region,proto3" json:"region,omitempty"`
}

func (x *EmissionsReport) Reset() {
	*x = EmissionsReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_calc_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EmissionsReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EmissionsReport) ProtoMessage() {}

func (x *EmissionsReport) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_calc_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EmissionsReport.ProtoReflect.Descriptor instead.
func (*EmissionsReport) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_calc_proto_rawDescGZIP(), []int{3}
}

func (x *EmissionsReport) GetDuration() *Period {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *EmissionsReport) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

func (x *EmissionsReport) GetPoints() []*DataPoint {
	if x != nil {
		return x.Points
	}
	return nil
}

func (x *EmissionsReport) GetOrgId() string {
	if x != nil {
		return x.OrgId
	}
	return ""
}

func (x *EmissionsReport) GetFacilityId() string {
	if x != nil {
		return x.FacilityId
	}
	return ""
}

func (x *EmissionsReport) GetLocationId() string {
	if x != nil {
		return x.LocationId
	}
	return ""
}

func (x *EmissionsReport) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

// Contains carbon emissions in terms of DataPoints, which can be used as
// points for a time/CO2 emissions graph
type DataPoint struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Time
	Time string `protobuf:"bytes,1,opt,name=time,proto3" json:"time,omitempty"`
	// either a carbon footprint(lbs of Co2) in a CarbonEmissions struct or power
	// stamp(KW) in an Electrical Report
	Value float64 `protobuf:"fixed64,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *DataPoint) Reset() {
	*x = DataPoint{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_calc_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DataPoint) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DataPoint) ProtoMessage() {}

func (x *DataPoint) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_calc_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DataPoint.ProtoReflect.Descriptor instead.
func (*DataPoint) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_calc_proto_rawDescGZIP(), []int{4}
}

func (x *DataPoint) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

func (x *DataPoint) GetValue() float64 {
	if x != nil {
		return x.Value
	}
	return 0
}

// Carbon Report from clickhouse
type CarbonReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Values are in units of (lbs of CO2/MWh)
	IntensityPoints []*DataPoint `protobuf:"bytes,1,rep,name=intensity_points,json=intensityPoints,proto3" json:"intensity_points,omitempty"`
	// Duration
	Duration *Period `protobuf:"bytes,2,opt,name=duration,proto3" json:"duration,omitempty"`
	Interval string  `protobuf:"bytes,3,opt,name=interval,proto3" json:"interval,omitempty"`
	Region   string  `protobuf:"bytes,4,opt,name=region,proto3" json:"region,omitempty"`
}

func (x *CarbonReport) Reset() {
	*x = CarbonReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_calc_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarbonReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarbonReport) ProtoMessage() {}

func (x *CarbonReport) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_calc_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarbonReport.ProtoReflect.Descriptor instead.
func (*CarbonReport) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_calc_proto_rawDescGZIP(), []int{5}
}

func (x *CarbonReport) GetIntensityPoints() []*DataPoint {
	if x != nil {
		return x.IntensityPoints
	}
	return nil
}

func (x *CarbonReport) GetDuration() *Period {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *CarbonReport) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

func (x *CarbonReport) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

// Energy Generation Report from the Past values function GetValues
type ElectricalReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Duration
	Duration *Period `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	// Power meter data in KWh
	PowerStamps []*DataPoint `protobuf:"bytes,2,rep,name=power_stamps,json=powerStamps,proto3" json:"power_stamps,omitempty"`
	Interval    string       `protobuf:"bytes,3,opt,name=interval,proto3" json:"interval,omitempty"`
}

func (x *ElectricalReport) Reset() {
	*x = ElectricalReport{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_calc_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ElectricalReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ElectricalReport) ProtoMessage() {}

func (x *ElectricalReport) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_calc_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ElectricalReport.ProtoReflect.Descriptor instead.
func (*ElectricalReport) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_calc_proto_rawDescGZIP(), []int{6}
}

func (x *ElectricalReport) GetDuration() *Period {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *ElectricalReport) GetPowerStamps() []*DataPoint {
	if x != nil {
		return x.PowerStamps
	}
	return nil
}

func (x *ElectricalReport) GetInterval() string {
	if x != nil {
		return x.Interval
	}
	return ""
}

var File_goadesign_goagen_calc_proto protoreflect.FileDescriptor

var file_goadesign_goagen_calc_proto_rawDesc = []byte{
	0x0a, 0x1b, 0x67, 0x6f, 0x61, 0x64, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x67, 0x6f, 0x61, 0x67,
	0x65, 0x6e, 0x5f, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63,
	0x61, 0x6c, 0x63, 0x22, 0xc1, 0x01, 0x0a, 0x20, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x63,
	0x61, 0x6c, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x72, 0x67, 0x5f,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49, 0x64, 0x12,
	0x28, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x52,
	0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x61, 0x63,
	0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x66, 0x61, 0x63, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e,
	0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6c, 0x6f, 0x63,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x22, 0x42, 0x0a, 0x06, 0x50, 0x65, 0x72, 0x69, 0x6f,
	0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0xf9, 0x01, 0x0a, 0x21,
	0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x63, 0x61, 0x6c, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e,
	0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x49, 0x0a, 0x15, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x65, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x15, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x13, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x45, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x4c, 0x0a, 0x18,
	0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x79,
	0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12,
	0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x52, 0x16, 0x63, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x6e, 0x73,
	0x69, 0x74, 0x79, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x12, 0x3b, 0x0a, 0x0d, 0x70, 0x6f,
	0x77, 0x65, 0x72, 0x5f, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x16, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x45, 0x6c, 0x65, 0x63, 0x74, 0x72, 0x69,
	0x63, 0x61, 0x6c, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x52, 0x0c, 0x70, 0x6f, 0x77, 0x65, 0x72,
	0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x73, 0x22, 0xf1, 0x01, 0x0a, 0x0f, 0x45, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x28, 0x0a, 0x08, 0x64,
	0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e,
	0x63, 0x61, 0x6c, 0x63, 0x2e, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x52, 0x08, 0x64, 0x75, 0x72,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61,
	0x6c, 0x12, 0x27, 0x0a, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69,
	0x6e, 0x74, 0x52, 0x06, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12, 0x15, 0x0a, 0x06, 0x6f, 0x72,
	0x67, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x72, 0x67, 0x49,
	0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x66, 0x61, 0x63, 0x69, 0x6c, 0x69, 0x74, 0x79, 0x5f, 0x69, 0x64,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x66, 0x61, 0x63, 0x69, 0x6c, 0x69, 0x74, 0x79,
	0x49, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x69,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x6c, 0x6f, 0x63, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x35, 0x0a, 0x09, 0x44,
	0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0xa8, 0x01, 0x0a, 0x0c, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x52, 0x65, 0x70,
	0x6f, 0x72, 0x74, 0x12, 0x3a, 0x0a, 0x10, 0x69, 0x6e, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x79,
	0x5f, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x63, 0x61, 0x6c, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x52, 0x0f,
	0x69, 0x6e, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x79, 0x50, 0x6f, 0x69, 0x6e, 0x74, 0x73, 0x12,
	0x28, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x52,
	0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74,
	0x65, 0x72, 0x76, 0x61, 0x6c, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x8c, 0x01,
	0x0a, 0x10, 0x45, 0x6c, 0x65, 0x63, 0x74, 0x72, 0x69, 0x63, 0x61, 0x6c, 0x52, 0x65, 0x70, 0x6f,
	0x72, 0x74, 0x12, 0x28, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x50, 0x65, 0x72, 0x69,
	0x6f, 0x64, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x32, 0x0a, 0x0c,
	0x70, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x73, 0x18, 0x02, 0x20, 0x03,
	0x28, 0x0b, 0x32, 0x0f, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x50, 0x6f,
	0x69, 0x6e, 0x74, 0x52, 0x0b, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61, 0x6d, 0x70, 0x73,
	0x12, 0x1a, 0x0a, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x76, 0x61, 0x6c, 0x32, 0x74, 0x0a, 0x04,
	0x43, 0x61, 0x6c, 0x63, 0x12, 0x6c, 0x0a, 0x19, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x63,
	0x61, 0x6c, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x26, 0x2e, 0x63, 0x61, 0x6c, 0x63, 0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69,
	0x63, 0x61, 0x6c, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x27, 0x2e, 0x63, 0x61, 0x6c, 0x63,
	0x2e, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x69, 0x63, 0x61, 0x6c, 0x43, 0x61, 0x72, 0x62, 0x6f,
	0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x09, 0x5a, 0x07, 0x2f, 0x63, 0x61, 0x6c, 0x63, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_goadesign_goagen_calc_proto_rawDescOnce sync.Once
	file_goadesign_goagen_calc_proto_rawDescData = file_goadesign_goagen_calc_proto_rawDesc
)

func file_goadesign_goagen_calc_proto_rawDescGZIP() []byte {
	file_goadesign_goagen_calc_proto_rawDescOnce.Do(func() {
		file_goadesign_goagen_calc_proto_rawDescData = protoimpl.X.CompressGZIP(file_goadesign_goagen_calc_proto_rawDescData)
	})
	return file_goadesign_goagen_calc_proto_rawDescData
}

var file_goadesign_goagen_calc_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_goadesign_goagen_calc_proto_goTypes = []interface{}{
	(*HistoricalCarbonEmissionsRequest)(nil),  // 0: calc.HistoricalCarbonEmissionsRequest
	(*Period)(nil),                            // 1: calc.Period
	(*HistoricalCarbonEmissionsResponse)(nil), // 2: calc.HistoricalCarbonEmissionsResponse
	(*EmissionsReport)(nil),                   // 3: calc.EmissionsReport
	(*DataPoint)(nil),                         // 4: calc.DataPoint
	(*CarbonReport)(nil),                      // 5: calc.CarbonReport
	(*ElectricalReport)(nil),                  // 6: calc.ElectricalReport
}
var file_goadesign_goagen_calc_proto_depIdxs = []int32{
	1,  // 0: calc.HistoricalCarbonEmissionsRequest.duration:type_name -> calc.Period
	3,  // 1: calc.HistoricalCarbonEmissionsResponse.total_emission_report:type_name -> calc.EmissionsReport
	5,  // 2: calc.HistoricalCarbonEmissionsResponse.carbon_intensity_reports:type_name -> calc.CarbonReport
	6,  // 3: calc.HistoricalCarbonEmissionsResponse.power_reports:type_name -> calc.ElectricalReport
	1,  // 4: calc.EmissionsReport.duration:type_name -> calc.Period
	4,  // 5: calc.EmissionsReport.points:type_name -> calc.DataPoint
	4,  // 6: calc.CarbonReport.intensity_points:type_name -> calc.DataPoint
	1,  // 7: calc.CarbonReport.duration:type_name -> calc.Period
	1,  // 8: calc.ElectricalReport.duration:type_name -> calc.Period
	4,  // 9: calc.ElectricalReport.power_stamps:type_name -> calc.DataPoint
	0,  // 10: calc.Calc.HistoricalCarbonEmissions:input_type -> calc.HistoricalCarbonEmissionsRequest
	2,  // 11: calc.Calc.HistoricalCarbonEmissions:output_type -> calc.HistoricalCarbonEmissionsResponse
	11, // [11:12] is the sub-list for method output_type
	10, // [10:11] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_goadesign_goagen_calc_proto_init() }
func file_goadesign_goagen_calc_proto_init() {
	if File_goadesign_goagen_calc_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_goadesign_goagen_calc_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HistoricalCarbonEmissionsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_goadesign_goagen_calc_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Period); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_goadesign_goagen_calc_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HistoricalCarbonEmissionsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_goadesign_goagen_calc_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EmissionsReport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_goadesign_goagen_calc_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DataPoint); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_goadesign_goagen_calc_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CarbonReport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_goadesign_goagen_calc_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ElectricalReport); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_goadesign_goagen_calc_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_goadesign_goagen_calc_proto_goTypes,
		DependencyIndexes: file_goadesign_goagen_calc_proto_depIdxs,
		MessageInfos:      file_goadesign_goagen_calc_proto_msgTypes,
	}.Build()
	File_goadesign_goagen_calc_proto = out.File
	file_goadesign_goagen_calc_proto_rawDesc = nil
	file_goadesign_goagen_calc_proto_goTypes = nil
	file_goadesign_goagen_calc_proto_depIdxs = nil
}
