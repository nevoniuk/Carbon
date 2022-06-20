// Code generated with goa v3.7.6, DO NOT EDIT.
//
// Poller protocol buffer definition
//
// Command:
// $ goa gen github.com/crossnokaye/carbon/services/poller/design

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: goadesign_goagen_poller.proto

package pollerpb

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

type CarbonEmissionsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CarbonEmissionsRequest) Reset() {
	*x = CarbonEmissionsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarbonEmissionsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarbonEmissionsRequest) ProtoMessage() {}

func (x *CarbonEmissionsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_poller_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarbonEmissionsRequest.ProtoReflect.Descriptor instead.
func (*CarbonEmissionsRequest) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{0}
}

type CarbonEmissionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field []*ArrayOfCarbonForecast `protobuf:"bytes,1,rep,name=field,proto3" json:"field,omitempty"`
}

func (x *CarbonEmissionsResponse) Reset() {
	*x = CarbonEmissionsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarbonEmissionsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarbonEmissionsResponse) ProtoMessage() {}

func (x *CarbonEmissionsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_poller_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarbonEmissionsResponse.ProtoReflect.Descriptor instead.
func (*CarbonEmissionsResponse) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{1}
}

func (x *CarbonEmissionsResponse) GetField() []*ArrayOfCarbonForecast {
	if x != nil {
		return x.Field
	}
	return nil
}

type ArrayOfCarbonForecast struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field []*CarbonForecast `protobuf:"bytes,1,rep,name=field,proto3" json:"field,omitempty"`
}

func (x *ArrayOfCarbonForecast) Reset() {
	*x = ArrayOfCarbonForecast{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ArrayOfCarbonForecast) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ArrayOfCarbonForecast) ProtoMessage() {}

func (x *ArrayOfCarbonForecast) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_poller_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ArrayOfCarbonForecast.ProtoReflect.Descriptor instead.
func (*ArrayOfCarbonForecast) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{2}
}

func (x *ArrayOfCarbonForecast) GetField() []*CarbonForecast {
	if x != nil {
		return x.Field
	}
	return nil
}

// Emissions Forecast
type CarbonForecast struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// generated_rate
	GeneratedRate float64 `protobuf:"fixed64,1,opt,name=generated_rate,json=generatedRate,proto3" json:"generated_rate,omitempty"`
	// marginal_rate
	MarginalRate float64 `protobuf:"fixed64,2,opt,name=marginal_rate,json=marginalRate,proto3" json:"marginal_rate,omitempty"`
	// consumed_rate
	ConsumedRate float64 `protobuf:"fixed64,3,opt,name=consumed_rate,json=consumedRate,proto3" json:"consumed_rate,omitempty"`
	// duration
	Duration *Period `protobuf:"bytes,4,opt,name=duration,proto3" json:"duration,omitempty"`
	// marginal_source
	MarginalSource string `protobuf:"bytes,5,opt,name=marginal_source,json=marginalSource,proto3" json:"marginal_source,omitempty"`
	// consumed_source
	ConsumedSource string `protobuf:"bytes,6,opt,name=consumed_source,json=consumedSource,proto3" json:"consumed_source,omitempty"`
	// generated_source
	GeneratedSource string `protobuf:"bytes,7,opt,name=generated_source,json=generatedSource,proto3" json:"generated_source,omitempty"`
	// emission_factor
	EmissionFactor string `protobuf:"bytes,8,opt,name=emission_factor,json=emissionFactor,proto3" json:"emission_factor,omitempty"`
	// region
	Region string `protobuf:"bytes,9,opt,name=region,proto3" json:"region,omitempty"`
}

func (x *CarbonForecast) Reset() {
	*x = CarbonForecast{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarbonForecast) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarbonForecast) ProtoMessage() {}

func (x *CarbonForecast) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_poller_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CarbonForecast.ProtoReflect.Descriptor instead.
func (*CarbonForecast) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{3}
}

func (x *CarbonForecast) GetGeneratedRate() float64 {
	if x != nil {
		return x.GeneratedRate
	}
	return 0
}

func (x *CarbonForecast) GetMarginalRate() float64 {
	if x != nil {
		return x.MarginalRate
	}
	return 0
}

func (x *CarbonForecast) GetConsumedRate() float64 {
	if x != nil {
		return x.ConsumedRate
	}
	return 0
}

func (x *CarbonForecast) GetDuration() *Period {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *CarbonForecast) GetMarginalSource() string {
	if x != nil {
		return x.MarginalSource
	}
	return ""
}

func (x *CarbonForecast) GetConsumedSource() string {
	if x != nil {
		return x.ConsumedSource
	}
	return ""
}

func (x *CarbonForecast) GetGeneratedSource() string {
	if x != nil {
		return x.GeneratedSource
	}
	return ""
}

func (x *CarbonForecast) GetEmissionFactor() string {
	if x != nil {
		return x.EmissionFactor
	}
	return ""
}

func (x *CarbonForecast) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

// Period of time from start to end of Forecast
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
		mi := &file_goadesign_goagen_poller_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Period) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Period) ProtoMessage() {}

func (x *Period) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_poller_proto_msgTypes[4]
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
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{4}
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

type AggregateDataRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AggregateDataRequest) Reset() {
	*x = AggregateDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregateDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregateDataRequest) ProtoMessage() {}

func (x *AggregateDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_poller_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregateDataRequest.ProtoReflect.Descriptor instead.
func (*AggregateDataRequest) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{5}
}

type AggregateDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field []*AggregateData `protobuf:"bytes,1,rep,name=field,proto3" json:"field,omitempty"`
}

func (x *AggregateDataResponse) Reset() {
	*x = AggregateDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregateDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregateDataResponse) ProtoMessage() {}

func (x *AggregateDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_poller_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregateDataResponse.ProtoReflect.Descriptor instead.
func (*AggregateDataResponse) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{6}
}

func (x *AggregateDataResponse) GetField() []*AggregateData {
	if x != nil {
		return x.Field
	}
	return nil
}

// aggregate data
type AggregateData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// average
	Average float64 `protobuf:"fixed64,1,opt,name=average,proto3" json:"average,omitempty"`
	// min
	Min float64 `protobuf:"fixed64,2,opt,name=min,proto3" json:"min,omitempty"`
	// max
	Max float64 `protobuf:"fixed64,3,opt,name=max,proto3" json:"max,omitempty"`
	// sum
	Sum float64 `protobuf:"fixed64,4,opt,name=sum,proto3" json:"sum,omitempty"`
	// count
	Count int32 `protobuf:"zigzag32,5,opt,name=count,proto3" json:"count,omitempty"`
	// duration
	Duration *Period `protobuf:"bytes,6,opt,name=duration,proto3" json:"duration,omitempty"`
	// report_type
	ReportType string `protobuf:"bytes,7,opt,name=report_type,json=reportType,proto3" json:"report_type,omitempty"`
}

func (x *AggregateData) Reset() {
	*x = AggregateData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregateData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregateData) ProtoMessage() {}

func (x *AggregateData) ProtoReflect() protoreflect.Message {
	mi := &file_goadesign_goagen_poller_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AggregateData.ProtoReflect.Descriptor instead.
func (*AggregateData) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{7}
}

func (x *AggregateData) GetAverage() float64 {
	if x != nil {
		return x.Average
	}
	return 0
}

func (x *AggregateData) GetMin() float64 {
	if x != nil {
		return x.Min
	}
	return 0
}

func (x *AggregateData) GetMax() float64 {
	if x != nil {
		return x.Max
	}
	return 0
}

func (x *AggregateData) GetSum() float64 {
	if x != nil {
		return x.Sum
	}
	return 0
}

func (x *AggregateData) GetCount() int32 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *AggregateData) GetDuration() *Period {
	if x != nil {
		return x.Duration
	}
	return nil
}

func (x *AggregateData) GetReportType() string {
	if x != nil {
		return x.ReportType
	}
	return ""
}

var File_goadesign_goagen_poller_proto protoreflect.FileDescriptor

var file_goadesign_goagen_poller_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x67, 0x6f, 0x61, 0x64, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x67, 0x6f, 0x61, 0x67,
	0x65, 0x6e, 0x5f, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x22, 0x18, 0x0a, 0x16, 0x43, 0x61, 0x72, 0x62, 0x6f,
	0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x22, 0x4e, 0x0a, 0x17, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x33, 0x0a, 0x05,
	0x66, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x70, 0x6f,
	0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x41, 0x72, 0x72, 0x61, 0x79, 0x4f, 0x66, 0x43, 0x61, 0x72, 0x62,
	0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x22, 0x45, 0x0a, 0x15, 0x41, 0x72, 0x72, 0x61, 0x79, 0x4f, 0x66, 0x43, 0x61, 0x72, 0x62,
	0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x6f, 0x6c, 0x6c,
	0x65, 0x72, 0x2e, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73,
	0x74, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x22, 0xeb, 0x02, 0x0a, 0x0e, 0x43, 0x61, 0x72,
	0x62, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73, 0x74, 0x12, 0x25, 0x0a, 0x0e, 0x67,
	0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x0d, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x52, 0x61,
	0x74, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x6d, 0x61, 0x72, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x72,
	0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x6d, 0x61, 0x72, 0x67, 0x69,
	0x6e, 0x61, 0x6c, 0x52, 0x61, 0x74, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x63, 0x6f, 0x6e, 0x73, 0x75,
	0x6d, 0x65, 0x64, 0x5f, 0x72, 0x61, 0x74, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c,
	0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x64, 0x52, 0x61, 0x74, 0x65, 0x12, 0x2a, 0x0a, 0x08,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e,
	0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x52, 0x08,
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x27, 0x0a, 0x0f, 0x6d, 0x61, 0x72, 0x67,
	0x69, 0x6e, 0x61, 0x6c, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x6d, 0x61, 0x72, 0x67, 0x69, 0x6e, 0x61, 0x6c, 0x53, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x75, 0x6d, 0x65, 0x64, 0x5f, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6e, 0x73,
	0x75, 0x6d, 0x65, 0x64, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x67, 0x65,
	0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x53,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x65, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x5f, 0x66, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e,
	0x65, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x46, 0x61, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x16,
	0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06,
	0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x22, 0x42, 0x0a, 0x06, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64,
	0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12,
	0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x07, 0x65, 0x6e, 0x64, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x16, 0x0a, 0x14, 0x41, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x22, 0x44, 0x0a, 0x15, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2b, 0x0a, 0x05, 0x66,
	0x69, 0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x6f, 0x6c,
	0x6c, 0x65, 0x72, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x22, 0xc2, 0x01, 0x0a, 0x0d, 0x41, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x76,
	0x65, 0x72, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x61, 0x76, 0x65,
	0x72, 0x61, 0x67, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x69, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x03, 0x6d, 0x69, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x6d, 0x61, 0x78, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x03, 0x6d, 0x61, 0x78, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x75, 0x6d, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x73, 0x75, 0x6d, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x18, 0x05, 0x20, 0x01, 0x28, 0x11, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x12, 0x2a, 0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x50, 0x65, 0x72, 0x69,
	0x6f, 0x64, 0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x1f, 0x0a, 0x0b,
	0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x54, 0x79, 0x70, 0x65, 0x32, 0xb2, 0x01,
	0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x52, 0x0a, 0x0f, 0x43, 0x61, 0x72, 0x62,
	0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x1e, 0x2e, 0x70, 0x6f,
	0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x6f,
	0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73,
	0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x54, 0x0a, 0x15,
	0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x64,
	0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x41,
	0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x41, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x0b, 0x5a, 0x09, 0x2f, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_goadesign_goagen_poller_proto_rawDescOnce sync.Once
	file_goadesign_goagen_poller_proto_rawDescData = file_goadesign_goagen_poller_proto_rawDesc
)

func file_goadesign_goagen_poller_proto_rawDescGZIP() []byte {
	file_goadesign_goagen_poller_proto_rawDescOnce.Do(func() {
		file_goadesign_goagen_poller_proto_rawDescData = protoimpl.X.CompressGZIP(file_goadesign_goagen_poller_proto_rawDescData)
	})
	return file_goadesign_goagen_poller_proto_rawDescData
}

var file_goadesign_goagen_poller_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_goadesign_goagen_poller_proto_goTypes = []interface{}{
	(*CarbonEmissionsRequest)(nil),  // 0: poller.CarbonEmissionsRequest
	(*CarbonEmissionsResponse)(nil), // 1: poller.CarbonEmissionsResponse
	(*ArrayOfCarbonForecast)(nil),   // 2: poller.ArrayOfCarbonForecast
	(*CarbonForecast)(nil),          // 3: poller.CarbonForecast
	(*Period)(nil),                  // 4: poller.Period
	(*AggregateDataRequest)(nil),    // 5: poller.AggregateDataRequest
	(*AggregateDataResponse)(nil),   // 6: poller.AggregateDataResponse
	(*AggregateData)(nil),           // 7: poller.AggregateData
}
var file_goadesign_goagen_poller_proto_depIdxs = []int32{
	2, // 0: poller.CarbonEmissionsResponse.field:type_name -> poller.ArrayOfCarbonForecast
	3, // 1: poller.ArrayOfCarbonForecast.field:type_name -> poller.CarbonForecast
	4, // 2: poller.CarbonForecast.duration:type_name -> poller.Period
	7, // 3: poller.AggregateDataResponse.field:type_name -> poller.AggregateData
	4, // 4: poller.AggregateData.duration:type_name -> poller.Period
	0, // 5: poller.Poller.CarbonEmissions:input_type -> poller.CarbonEmissionsRequest
	5, // 6: poller.Poller.AggregateDataEndpoint:input_type -> poller.AggregateDataRequest
	1, // 7: poller.Poller.CarbonEmissions:output_type -> poller.CarbonEmissionsResponse
	6, // 8: poller.Poller.AggregateDataEndpoint:output_type -> poller.AggregateDataResponse
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_goadesign_goagen_poller_proto_init() }
func file_goadesign_goagen_poller_proto_init() {
	if File_goadesign_goagen_poller_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_goadesign_goagen_poller_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CarbonEmissionsRequest); i {
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
		file_goadesign_goagen_poller_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CarbonEmissionsResponse); i {
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
		file_goadesign_goagen_poller_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ArrayOfCarbonForecast); i {
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
		file_goadesign_goagen_poller_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CarbonForecast); i {
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
		file_goadesign_goagen_poller_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_goadesign_goagen_poller_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregateDataRequest); i {
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
		file_goadesign_goagen_poller_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregateDataResponse); i {
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
		file_goadesign_goagen_poller_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AggregateData); i {
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
			RawDescriptor: file_goadesign_goagen_poller_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_goadesign_goagen_poller_proto_goTypes,
		DependencyIndexes: file_goadesign_goagen_poller_proto_depIdxs,
		MessageInfos:      file_goadesign_goagen_poller_proto_msgTypes,
	}.Build()
	File_goadesign_goagen_poller_proto = out.File
	file_goadesign_goagen_poller_proto_rawDesc = nil
	file_goadesign_goagen_poller_proto_goTypes = nil
	file_goadesign_goagen_poller_proto_depIdxs = nil
}
