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

	// region
	Region string `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	// start
	Start string `protobuf:"bytes,2,opt,name=start,proto3" json:"start,omitempty"`
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

func (x *CarbonEmissionsRequest) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *CarbonEmissionsRequest) GetStart() string {
	if x != nil {
		return x.Start
	}
	return ""
}

type CarbonEmissionsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Field []*CarbonForecast `protobuf:"bytes,1,rep,name=field,proto3" json:"field,omitempty"`
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

func (x *CarbonEmissionsResponse) GetField() []*CarbonForecast {
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
	// Duration
	Duration *Period `protobuf:"bytes,4,opt,name=duration,proto3" json:"duration,omitempty"`
	// duration_type
	DurationType string `protobuf:"bytes,5,opt,name=duration_type,json=durationType,proto3" json:"duration_type,omitempty"`
	// generated_source
	GeneratedSource string `protobuf:"bytes,6,opt,name=generated_source,json=generatedSource,proto3" json:"generated_source,omitempty"`
	// region
	Region string `protobuf:"bytes,7,opt,name=region,proto3" json:"region,omitempty"`
}

func (x *CarbonForecast) Reset() {
	*x = CarbonForecast{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CarbonForecast) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CarbonForecast) ProtoMessage() {}

func (x *CarbonForecast) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use CarbonForecast.ProtoReflect.Descriptor instead.
func (*CarbonForecast) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{2}
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

func (x *CarbonForecast) GetDurationType() string {
	if x != nil {
		return x.DurationType
	}
	return ""
}

func (x *CarbonForecast) GetGeneratedSource() string {
	if x != nil {
		return x.GeneratedSource
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
		mi := &file_goadesign_goagen_poller_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Period) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Period) ProtoMessage() {}

func (x *Period) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use Period.ProtoReflect.Descriptor instead.
func (*Period) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{3}
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

	// region
	Region string `protobuf:"bytes,1,opt,name=region,proto3" json:"region,omitempty"`
	// periods
	Periods []*Period `protobuf:"bytes,2,rep,name=periods,proto3" json:"periods,omitempty"`
	// duration
	Duration string `protobuf:"bytes,3,opt,name=duration,proto3" json:"duration,omitempty"`
}

func (x *AggregateDataRequest) Reset() {
	*x = AggregateDataRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregateDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregateDataRequest) ProtoMessage() {}

func (x *AggregateDataRequest) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use AggregateDataRequest.ProtoReflect.Descriptor instead.
func (*AggregateDataRequest) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{4}
}

func (x *AggregateDataRequest) GetRegion() string {
	if x != nil {
		return x.Region
	}
	return ""
}

func (x *AggregateDataRequest) GetPeriods() []*Period {
	if x != nil {
		return x.Periods
	}
	return nil
}

func (x *AggregateDataRequest) GetDuration() string {
	if x != nil {
		return x.Duration
	}
	return ""
}

type AggregateDataResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *AggregateDataResponse) Reset() {
	*x = AggregateDataResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_goadesign_goagen_poller_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AggregateDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AggregateDataResponse) ProtoMessage() {}

func (x *AggregateDataResponse) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use AggregateDataResponse.ProtoReflect.Descriptor instead.
func (*AggregateDataResponse) Descriptor() ([]byte, []int) {
	return file_goadesign_goagen_poller_proto_rawDescGZIP(), []int{5}
}

var File_goadesign_goagen_poller_proto protoreflect.FileDescriptor

var file_goadesign_goagen_poller_proto_rawDesc = []byte{
	0x0a, 0x1d, 0x67, 0x6f, 0x61, 0x64, 0x65, 0x73, 0x69, 0x67, 0x6e, 0x5f, 0x67, 0x6f, 0x61, 0x67,
	0x65, 0x6e, 0x5f, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x06, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x22, 0x46, 0x0a, 0x16, 0x43, 0x61, 0x72, 0x62, 0x6f,
	0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x73, 0x74, 0x61,
	0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x22,
	0x47, 0x0a, 0x17, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2c, 0x0a, 0x05, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x70, 0x6f, 0x6c, 0x6c,
	0x65, 0x72, 0x2e, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x46, 0x6f, 0x72, 0x65, 0x63, 0x61, 0x73,
	0x74, 0x52, 0x05, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x22, 0x95, 0x02, 0x0a, 0x0e, 0x43, 0x61, 0x72,
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
	0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x23, 0x0a, 0x0d, 0x64, 0x75, 0x72, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x79, 0x70, 0x65, 0x12, 0x29, 0x0a,
	0x10, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x5f, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74,
	0x65, 0x64, 0x53, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x67, 0x69,
	0x6f, 0x6e, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e,
	0x22, 0x42, 0x0a, 0x06, 0x50, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x73, 0x74,
	0x61, 0x72, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09,
	0x73, 0x74, 0x61, 0x72, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x19, 0x0a, 0x08, 0x65, 0x6e, 0x64,
	0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x65, 0x6e, 0x64,
	0x54, 0x69, 0x6d, 0x65, 0x22, 0x74, 0x0a, 0x14, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74,
	0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06,
	0x72, 0x65, 0x67, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x72, 0x65,
	0x67, 0x69, 0x6f, 0x6e, 0x12, 0x28, 0x0a, 0x07, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x73, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x50,
	0x65, 0x72, 0x69, 0x6f, 0x64, 0x52, 0x07, 0x70, 0x65, 0x72, 0x69, 0x6f, 0x64, 0x73, 0x12, 0x1a,
	0x0a, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x64, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x17, 0x0a, 0x15, 0x41, 0x67,
	0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x32, 0xaa, 0x01, 0x0a, 0x06, 0x50, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x52,
	0x0a, 0x0f, 0x43, 0x61, 0x72, 0x62, 0x6f, 0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e,
	0x73, 0x12, 0x1e, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x43, 0x61, 0x72, 0x62, 0x6f,
	0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x43, 0x61, 0x72, 0x62, 0x6f,
	0x6e, 0x45, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x4c, 0x0a, 0x0d, 0x41, 0x67, 0x67, 0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x1c, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x41, 0x67, 0x67,
	0x72, 0x65, 0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x2e, 0x41, 0x67, 0x67, 0x72, 0x65,
	0x67, 0x61, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x42, 0x0b, 0x5a, 0x09, 0x2f, 0x70, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
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

var file_goadesign_goagen_poller_proto_msgTypes = make([]protoimpl.MessageInfo, 6)
var file_goadesign_goagen_poller_proto_goTypes = []interface{}{
	(*CarbonEmissionsRequest)(nil),  // 0: poller.CarbonEmissionsRequest
	(*CarbonEmissionsResponse)(nil), // 1: poller.CarbonEmissionsResponse
	(*CarbonForecast)(nil),          // 2: poller.CarbonForecast
	(*Period)(nil),                  // 3: poller.Period
	(*AggregateDataRequest)(nil),    // 4: poller.AggregateDataRequest
	(*AggregateDataResponse)(nil),   // 5: poller.AggregateDataResponse
}
var file_goadesign_goagen_poller_proto_depIdxs = []int32{
	2, // 0: poller.CarbonEmissionsResponse.field:type_name -> poller.CarbonForecast
	3, // 1: poller.CarbonForecast.duration:type_name -> poller.Period
	3, // 2: poller.AggregateDataRequest.periods:type_name -> poller.Period
	0, // 3: poller.Poller.CarbonEmissions:input_type -> poller.CarbonEmissionsRequest
	4, // 4: poller.Poller.AggregateData:input_type -> poller.AggregateDataRequest
	1, // 5: poller.Poller.CarbonEmissions:output_type -> poller.CarbonEmissionsResponse
	5, // 6: poller.Poller.AggregateData:output_type -> poller.AggregateDataResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
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
		file_goadesign_goagen_poller_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
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
		file_goadesign_goagen_poller_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
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
		file_goadesign_goagen_poller_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
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
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_goadesign_goagen_poller_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   6,
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
