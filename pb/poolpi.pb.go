// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.13.0
// source: pb/poolpi.proto

package pb

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

type EventRequest_Event int32

const (
	EventRequest_KeyNone    EventRequest_Event = 0
	EventRequest_KeyRight   EventRequest_Event = 1
	EventRequest_KeyMenu    EventRequest_Event = 2
	EventRequest_KeyLeft    EventRequest_Event = 3
	EventRequest_KeyService EventRequest_Event = 4
	EventRequest_KeyMinus   EventRequest_Event = 5
	EventRequest_KeyPlus    EventRequest_Event = 6
	EventRequest_KeyPoolSpa EventRequest_Event = 7
	EventRequest_KeyFilter  EventRequest_Event = 8
	EventRequest_KeyLights  EventRequest_Event = 9
	EventRequest_KeyAux1    EventRequest_Event = 10
	EventRequest_KeyAux2    EventRequest_Event = 11
	EventRequest_KeyAux3    EventRequest_Event = 12
	EventRequest_KeyAux4    EventRequest_Event = 13
	EventRequest_KeyAux5    EventRequest_Event = 14
	EventRequest_KeyAux6    EventRequest_Event = 15
	EventRequest_KeyAux7    EventRequest_Event = 16
	EventRequest_KeyValve3  EventRequest_Event = 17
	EventRequest_KeyValve4  EventRequest_Event = 18
	EventRequest_KeyHeater  EventRequest_Event = 19
)

// Enum value maps for EventRequest_Event.
var (
	EventRequest_Event_name = map[int32]string{
		0:  "KeyNone",
		1:  "KeyRight",
		2:  "KeyMenu",
		3:  "KeyLeft",
		4:  "KeyService",
		5:  "KeyMinus",
		6:  "KeyPlus",
		7:  "KeyPoolSpa",
		8:  "KeyFilter",
		9:  "KeyLights",
		10: "KeyAux1",
		11: "KeyAux2",
		12: "KeyAux3",
		13: "KeyAux4",
		14: "KeyAux5",
		15: "KeyAux6",
		16: "KeyAux7",
		17: "KeyValve3",
		18: "KeyValve4",
		19: "KeyHeater",
	}
	EventRequest_Event_value = map[string]int32{
		"KeyNone":    0,
		"KeyRight":   1,
		"KeyMenu":    2,
		"KeyLeft":    3,
		"KeyService": 4,
		"KeyMinus":   5,
		"KeyPlus":    6,
		"KeyPoolSpa": 7,
		"KeyFilter":  8,
		"KeyLights":  9,
		"KeyAux1":    10,
		"KeyAux2":    11,
		"KeyAux3":    12,
		"KeyAux4":    13,
		"KeyAux5":    14,
		"KeyAux6":    15,
		"KeyAux7":    16,
		"KeyValve3":  17,
		"KeyValve4":  18,
		"KeyHeater":  19,
	}
)

func (x EventRequest_Event) Enum() *EventRequest_Event {
	p := new(EventRequest_Event)
	*p = x
	return p
}

func (x EventRequest_Event) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EventRequest_Event) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_poolpi_proto_enumTypes[0].Descriptor()
}

func (EventRequest_Event) Type() protoreflect.EnumType {
	return &file_pb_poolpi_proto_enumTypes[0]
}

func (x EventRequest_Event) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EventRequest_Event.Descriptor instead.
func (EventRequest_Event) EnumDescriptor() ([]byte, []int) {
	return file_pb_poolpi_proto_rawDescGZIP(), []int{5, 0}
}

type MonitorStateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MonitorStateRequest) Reset() {
	*x = MonitorStateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_poolpi_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MonitorStateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MonitorStateRequest) ProtoMessage() {}

func (x *MonitorStateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_poolpi_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MonitorStateRequest.ProtoReflect.Descriptor instead.
func (*MonitorStateRequest) Descriptor() ([]byte, []int) {
	return file_pb_poolpi_proto_rawDescGZIP(), []int{0}
}

type MonitorStateResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Heater1         *Indicator `protobuf:"bytes,1,opt,name=Heater1,proto3" json:"Heater1,omitempty"`
	Valve3          *Indicator `protobuf:"bytes,2,opt,name=Valve3,proto3" json:"Valve3,omitempty"`
	CheckSystem     *Indicator `protobuf:"bytes,3,opt,name=CheckSystem,proto3" json:"CheckSystem,omitempty"`
	Pool            *Indicator `protobuf:"bytes,4,opt,name=Pool,proto3" json:"Pool,omitempty"`
	Spa             *Indicator `protobuf:"bytes,5,opt,name=Spa,proto3" json:"Spa,omitempty"`
	Filter          *Indicator `protobuf:"bytes,6,opt,name=Filter,proto3" json:"Filter,omitempty"`
	Lights          *Indicator `protobuf:"bytes,7,opt,name=Lights,proto3" json:"Lights,omitempty"`
	Aux1            *Indicator `protobuf:"bytes,8,opt,name=Aux1,proto3" json:"Aux1,omitempty"`
	Aux2            *Indicator `protobuf:"bytes,9,opt,name=Aux2,proto3" json:"Aux2,omitempty"`
	Service         *Indicator `protobuf:"bytes,10,opt,name=Service,proto3" json:"Service,omitempty"`
	Aux3            *Indicator `protobuf:"bytes,11,opt,name=Aux3,proto3" json:"Aux3,omitempty"`
	Aux4            *Indicator `protobuf:"bytes,12,opt,name=Aux4,proto3" json:"Aux4,omitempty"`
	Aux5            *Indicator `protobuf:"bytes,13,opt,name=Aux5,proto3" json:"Aux5,omitempty"`
	Aux6            *Indicator `protobuf:"bytes,14,opt,name=Aux6,proto3" json:"Aux6,omitempty"`
	Valve4          *Indicator `protobuf:"bytes,15,opt,name=Valve4,proto3" json:"Valve4,omitempty"`
	Spillover       *Indicator `protobuf:"bytes,16,opt,name=Spillover,proto3" json:"Spillover,omitempty"`
	SystemOff       *Indicator `protobuf:"bytes,17,opt,name=SystemOff,proto3" json:"SystemOff,omitempty"`
	Aux7            *Indicator `protobuf:"bytes,18,opt,name=Aux7,proto3" json:"Aux7,omitempty"`
	Aux8            *Indicator `protobuf:"bytes,19,opt,name=Aux8,proto3" json:"Aux8,omitempty"`
	Aux9            *Indicator `protobuf:"bytes,20,opt,name=Aux9,proto3" json:"Aux9,omitempty"`
	Aux10           *Indicator `protobuf:"bytes,21,opt,name=Aux10,proto3" json:"Aux10,omitempty"`
	Aux11           *Indicator `protobuf:"bytes,22,opt,name=Aux11,proto3" json:"Aux11,omitempty"`
	Aux12           *Indicator `protobuf:"bytes,23,opt,name=Aux12,proto3" json:"Aux12,omitempty"`
	Aux13           *Indicator `protobuf:"bytes,24,opt,name=Aux13,proto3" json:"Aux13,omitempty"`
	Auz14           *Indicator `protobuf:"bytes,25,opt,name=Auz14,proto3" json:"Auz14,omitempty"`
	SuperChlorinate *Indicator `protobuf:"bytes,26,opt,name=SuperChlorinate,proto3" json:"SuperChlorinate,omitempty"`
}

func (x *MonitorStateResponse) Reset() {
	*x = MonitorStateResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_poolpi_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MonitorStateResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MonitorStateResponse) ProtoMessage() {}

func (x *MonitorStateResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_poolpi_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MonitorStateResponse.ProtoReflect.Descriptor instead.
func (*MonitorStateResponse) Descriptor() ([]byte, []int) {
	return file_pb_poolpi_proto_rawDescGZIP(), []int{1}
}

func (x *MonitorStateResponse) GetHeater1() *Indicator {
	if x != nil {
		return x.Heater1
	}
	return nil
}

func (x *MonitorStateResponse) GetValve3() *Indicator {
	if x != nil {
		return x.Valve3
	}
	return nil
}

func (x *MonitorStateResponse) GetCheckSystem() *Indicator {
	if x != nil {
		return x.CheckSystem
	}
	return nil
}

func (x *MonitorStateResponse) GetPool() *Indicator {
	if x != nil {
		return x.Pool
	}
	return nil
}

func (x *MonitorStateResponse) GetSpa() *Indicator {
	if x != nil {
		return x.Spa
	}
	return nil
}

func (x *MonitorStateResponse) GetFilter() *Indicator {
	if x != nil {
		return x.Filter
	}
	return nil
}

func (x *MonitorStateResponse) GetLights() *Indicator {
	if x != nil {
		return x.Lights
	}
	return nil
}

func (x *MonitorStateResponse) GetAux1() *Indicator {
	if x != nil {
		return x.Aux1
	}
	return nil
}

func (x *MonitorStateResponse) GetAux2() *Indicator {
	if x != nil {
		return x.Aux2
	}
	return nil
}

func (x *MonitorStateResponse) GetService() *Indicator {
	if x != nil {
		return x.Service
	}
	return nil
}

func (x *MonitorStateResponse) GetAux3() *Indicator {
	if x != nil {
		return x.Aux3
	}
	return nil
}

func (x *MonitorStateResponse) GetAux4() *Indicator {
	if x != nil {
		return x.Aux4
	}
	return nil
}

func (x *MonitorStateResponse) GetAux5() *Indicator {
	if x != nil {
		return x.Aux5
	}
	return nil
}

func (x *MonitorStateResponse) GetAux6() *Indicator {
	if x != nil {
		return x.Aux6
	}
	return nil
}

func (x *MonitorStateResponse) GetValve4() *Indicator {
	if x != nil {
		return x.Valve4
	}
	return nil
}

func (x *MonitorStateResponse) GetSpillover() *Indicator {
	if x != nil {
		return x.Spillover
	}
	return nil
}

func (x *MonitorStateResponse) GetSystemOff() *Indicator {
	if x != nil {
		return x.SystemOff
	}
	return nil
}

func (x *MonitorStateResponse) GetAux7() *Indicator {
	if x != nil {
		return x.Aux7
	}
	return nil
}

func (x *MonitorStateResponse) GetAux8() *Indicator {
	if x != nil {
		return x.Aux8
	}
	return nil
}

func (x *MonitorStateResponse) GetAux9() *Indicator {
	if x != nil {
		return x.Aux9
	}
	return nil
}

func (x *MonitorStateResponse) GetAux10() *Indicator {
	if x != nil {
		return x.Aux10
	}
	return nil
}

func (x *MonitorStateResponse) GetAux11() *Indicator {
	if x != nil {
		return x.Aux11
	}
	return nil
}

func (x *MonitorStateResponse) GetAux12() *Indicator {
	if x != nil {
		return x.Aux12
	}
	return nil
}

func (x *MonitorStateResponse) GetAux13() *Indicator {
	if x != nil {
		return x.Aux13
	}
	return nil
}

func (x *MonitorStateResponse) GetAuz14() *Indicator {
	if x != nil {
		return x.Auz14
	}
	return nil
}

func (x *MonitorStateResponse) GetSuperChlorinate() *Indicator {
	if x != nil {
		return x.SuperChlorinate
	}
	return nil
}

type Indicator struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Active  bool `protobuf:"varint,1,opt,name=Active,proto3" json:"Active,omitempty"`
	Caution bool `protobuf:"varint,2,opt,name=Caution,proto3" json:"Caution,omitempty"` // Light Blinking
}

func (x *Indicator) Reset() {
	*x = Indicator{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_poolpi_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Indicator) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Indicator) ProtoMessage() {}

func (x *Indicator) ProtoReflect() protoreflect.Message {
	mi := &file_pb_poolpi_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Indicator.ProtoReflect.Descriptor instead.
func (*Indicator) Descriptor() ([]byte, []int) {
	return file_pb_poolpi_proto_rawDescGZIP(), []int{2}
}

func (x *Indicator) GetActive() bool {
	if x != nil {
		return x.Active
	}
	return false
}

func (x *Indicator) GetCaution() bool {
	if x != nil {
		return x.Caution
	}
	return false
}

type MessagesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MessagesRequest) Reset() {
	*x = MessagesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_poolpi_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessagesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessagesRequest) ProtoMessage() {}

func (x *MessagesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_poolpi_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessagesRequest.ProtoReflect.Descriptor instead.
func (*MessagesRequest) Descriptor() ([]byte, []int) {
	return file_pb_poolpi_proto_rawDescGZIP(), []int{3}
}

type MessagesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message []byte `protobuf:"bytes,1,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *MessagesResponse) Reset() {
	*x = MessagesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_poolpi_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MessagesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessagesResponse) ProtoMessage() {}

func (x *MessagesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_poolpi_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MessagesResponse.ProtoReflect.Descriptor instead.
func (*MessagesResponse) Descriptor() ([]byte, []int) {
	return file_pb_poolpi_proto_rawDescGZIP(), []int{4}
}

func (x *MessagesResponse) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

type EventRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EventRequest) Reset() {
	*x = EventRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_poolpi_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventRequest) ProtoMessage() {}

func (x *EventRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pb_poolpi_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventRequest.ProtoReflect.Descriptor instead.
func (*EventRequest) Descriptor() ([]byte, []int) {
	return file_pb_poolpi_proto_rawDescGZIP(), []int{5}
}

type EventResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *EventResponse) Reset() {
	*x = EventResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_poolpi_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EventResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EventResponse) ProtoMessage() {}

func (x *EventResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_poolpi_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EventResponse.ProtoReflect.Descriptor instead.
func (*EventResponse) Descriptor() ([]byte, []int) {
	return file_pb_poolpi_proto_rawDescGZIP(), []int{6}
}

var File_pb_poolpi_proto protoreflect.FileDescriptor

var file_pb_poolpi_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x70, 0x62, 0x2f, 0x70, 0x6f, 0x6f, 0x6c, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x15, 0x0a, 0x13, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0xb2, 0x07, 0x0a, 0x14, 0x4d, 0x6f, 0x6e,
	0x69, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x24, 0x0a, 0x07, 0x48, 0x65, 0x61, 0x74, 0x65, 0x72, 0x31, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x07,
	0x48, 0x65, 0x61, 0x74, 0x65, 0x72, 0x31, 0x12, 0x22, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x76, 0x65,
	0x33, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61,
	0x74, 0x6f, 0x72, 0x52, 0x06, 0x56, 0x61, 0x6c, 0x76, 0x65, 0x33, 0x12, 0x2c, 0x0a, 0x0b, 0x43,
	0x68, 0x65, 0x63, 0x6b, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x0b, 0x43, 0x68,
	0x65, 0x63, 0x6b, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x12, 0x1e, 0x0a, 0x04, 0x50, 0x6f, 0x6f,
	0x6c, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61,
	0x74, 0x6f, 0x72, 0x52, 0x04, 0x50, 0x6f, 0x6f, 0x6c, 0x12, 0x1c, 0x0a, 0x03, 0x53, 0x70, 0x61,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x52, 0x03, 0x53, 0x70, 0x61, 0x12, 0x22, 0x0a, 0x06, 0x46, 0x69, 0x6c, 0x74, 0x65,
	0x72, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61,
	0x74, 0x6f, 0x72, 0x52, 0x06, 0x46, 0x69, 0x6c, 0x74, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x06, 0x4c,
	0x69, 0x67, 0x68, 0x74, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e,
	0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x06, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x12,
	0x1e, 0x0a, 0x04, 0x41, 0x75, 0x78, 0x31, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x04, 0x41, 0x75, 0x78, 0x31, 0x12,
	0x1e, 0x0a, 0x04, 0x41, 0x75, 0x78, 0x32, 0x18, 0x09, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e,
	0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x04, 0x41, 0x75, 0x78, 0x32, 0x12,
	0x24, 0x0a, 0x07, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x07, 0x53, 0x65,
	0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x1e, 0x0a, 0x04, 0x41, 0x75, 0x78, 0x33, 0x18, 0x0b, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52,
	0x04, 0x41, 0x75, 0x78, 0x33, 0x12, 0x1e, 0x0a, 0x04, 0x41, 0x75, 0x78, 0x34, 0x18, 0x0c, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52,
	0x04, 0x41, 0x75, 0x78, 0x34, 0x12, 0x1e, 0x0a, 0x04, 0x41, 0x75, 0x78, 0x35, 0x18, 0x0d, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52,
	0x04, 0x41, 0x75, 0x78, 0x35, 0x12, 0x1e, 0x0a, 0x04, 0x41, 0x75, 0x78, 0x36, 0x18, 0x0e, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52,
	0x04, 0x41, 0x75, 0x78, 0x36, 0x12, 0x22, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x76, 0x65, 0x34, 0x18,
	0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f,
	0x72, 0x52, 0x06, 0x56, 0x61, 0x6c, 0x76, 0x65, 0x34, 0x12, 0x28, 0x0a, 0x09, 0x53, 0x70, 0x69,
	0x6c, 0x6c, 0x6f, 0x76, 0x65, 0x72, 0x18, 0x10, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49,
	0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x09, 0x53, 0x70, 0x69, 0x6c, 0x6c, 0x6f,
	0x76, 0x65, 0x72, 0x12, 0x28, 0x0a, 0x09, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4f, 0x66, 0x66,
	0x18, 0x11, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74,
	0x6f, 0x72, 0x52, 0x09, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x4f, 0x66, 0x66, 0x12, 0x1e, 0x0a,
	0x04, 0x41, 0x75, 0x78, 0x37, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e,
	0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x04, 0x41, 0x75, 0x78, 0x37, 0x12, 0x1e, 0x0a,
	0x04, 0x41, 0x75, 0x78, 0x38, 0x18, 0x13, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e,
	0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x04, 0x41, 0x75, 0x78, 0x38, 0x12, 0x1e, 0x0a,
	0x04, 0x41, 0x75, 0x78, 0x39, 0x18, 0x14, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e,
	0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x04, 0x41, 0x75, 0x78, 0x39, 0x12, 0x20, 0x0a,
	0x05, 0x41, 0x75, 0x78, 0x31, 0x30, 0x18, 0x15, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49,
	0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x05, 0x41, 0x75, 0x78, 0x31, 0x30, 0x12,
	0x20, 0x0a, 0x05, 0x41, 0x75, 0x78, 0x31, 0x31, 0x18, 0x16, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a,
	0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x05, 0x41, 0x75, 0x78, 0x31,
	0x31, 0x12, 0x20, 0x0a, 0x05, 0x41, 0x75, 0x78, 0x31, 0x32, 0x18, 0x17, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x05, 0x41, 0x75,
	0x78, 0x31, 0x32, 0x12, 0x20, 0x0a, 0x05, 0x41, 0x75, 0x78, 0x31, 0x33, 0x18, 0x18, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x05,
	0x41, 0x75, 0x78, 0x31, 0x33, 0x12, 0x20, 0x0a, 0x05, 0x41, 0x75, 0x7a, 0x31, 0x34, 0x18, 0x19,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72,
	0x52, 0x05, 0x41, 0x75, 0x7a, 0x31, 0x34, 0x12, 0x34, 0x0a, 0x0f, 0x53, 0x75, 0x70, 0x65, 0x72,
	0x43, 0x68, 0x6c, 0x6f, 0x72, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x18, 0x1a, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x0a, 0x2e, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x52, 0x0f, 0x53, 0x75,
	0x70, 0x65, 0x72, 0x43, 0x68, 0x6c, 0x6f, 0x72, 0x69, 0x6e, 0x61, 0x74, 0x65, 0x22, 0x3d, 0x0a,
	0x09, 0x49, 0x6e, 0x64, 0x69, 0x63, 0x61, 0x74, 0x6f, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x41, 0x63,
	0x74, 0x69, 0x76, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x41, 0x63, 0x74, 0x69,
	0x76, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x43, 0x61, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x07, 0x43, 0x61, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x11, 0x0a, 0x0f,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22,
	0x2c, 0x0a, 0x10, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0xae, 0x02,
	0x0a, 0x0c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x9d,
	0x02, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x4e,
	0x6f, 0x6e, 0x65, 0x10, 0x00, 0x12, 0x0c, 0x0a, 0x08, 0x4b, 0x65, 0x79, 0x52, 0x69, 0x67, 0x68,
	0x74, 0x10, 0x01, 0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x4d, 0x65, 0x6e, 0x75, 0x10, 0x02,
	0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x4c, 0x65, 0x66, 0x74, 0x10, 0x03, 0x12, 0x0e, 0x0a,
	0x0a, 0x4b, 0x65, 0x79, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x10, 0x04, 0x12, 0x0c, 0x0a,
	0x08, 0x4b, 0x65, 0x79, 0x4d, 0x69, 0x6e, 0x75, 0x73, 0x10, 0x05, 0x12, 0x0b, 0x0a, 0x07, 0x4b,
	0x65, 0x79, 0x50, 0x6c, 0x75, 0x73, 0x10, 0x06, 0x12, 0x0e, 0x0a, 0x0a, 0x4b, 0x65, 0x79, 0x50,
	0x6f, 0x6f, 0x6c, 0x53, 0x70, 0x61, 0x10, 0x07, 0x12, 0x0d, 0x0a, 0x09, 0x4b, 0x65, 0x79, 0x46,
	0x69, 0x6c, 0x74, 0x65, 0x72, 0x10, 0x08, 0x12, 0x0d, 0x0a, 0x09, 0x4b, 0x65, 0x79, 0x4c, 0x69,
	0x67, 0x68, 0x74, 0x73, 0x10, 0x09, 0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x78,
	0x31, 0x10, 0x0a, 0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x78, 0x32, 0x10, 0x0b,
	0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x78, 0x33, 0x10, 0x0c, 0x12, 0x0b, 0x0a,
	0x07, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x78, 0x34, 0x10, 0x0d, 0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65,
	0x79, 0x41, 0x75, 0x78, 0x35, 0x10, 0x0e, 0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x41, 0x75,
	0x78, 0x36, 0x10, 0x0f, 0x12, 0x0b, 0x0a, 0x07, 0x4b, 0x65, 0x79, 0x41, 0x75, 0x78, 0x37, 0x10,
	0x10, 0x12, 0x0d, 0x0a, 0x09, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x76, 0x65, 0x33, 0x10, 0x11,
	0x12, 0x0d, 0x0a, 0x09, 0x4b, 0x65, 0x79, 0x56, 0x61, 0x6c, 0x76, 0x65, 0x34, 0x10, 0x12, 0x12,
	0x0d, 0x0a, 0x09, 0x4b, 0x65, 0x79, 0x48, 0x65, 0x61, 0x74, 0x65, 0x72, 0x10, 0x13, 0x22, 0x0f,
	0x0a, 0x0d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32,
	0xac, 0x01, 0x0a, 0x0a, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x3f,
	0x0a, 0x0c, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x14,
	0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x4d, 0x6f, 0x6e, 0x69, 0x74, 0x6f, 0x72, 0x53, 0x74,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x30, 0x01, 0x12,
	0x33, 0x0a, 0x08, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x12, 0x10, 0x2e, 0x4d, 0x65,
	0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e,
	0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x22, 0x00, 0x30, 0x01, 0x12, 0x28, 0x0a, 0x05, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0d, 0x2e,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x1c,
	0x5a, 0x1a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x63, 0x6f, 0x72,
	0x79, 0x62, 0x2f, 0x70, 0x6f, 0x6f, 0x6c, 0x70, 0x69, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_poolpi_proto_rawDescOnce sync.Once
	file_pb_poolpi_proto_rawDescData = file_pb_poolpi_proto_rawDesc
)

func file_pb_poolpi_proto_rawDescGZIP() []byte {
	file_pb_poolpi_proto_rawDescOnce.Do(func() {
		file_pb_poolpi_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_poolpi_proto_rawDescData)
	})
	return file_pb_poolpi_proto_rawDescData
}

var file_pb_poolpi_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pb_poolpi_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_pb_poolpi_proto_goTypes = []interface{}{
	(EventRequest_Event)(0),      // 0: EventRequest.Event
	(*MonitorStateRequest)(nil),  // 1: MonitorStateRequest
	(*MonitorStateResponse)(nil), // 2: MonitorStateResponse
	(*Indicator)(nil),            // 3: Indicator
	(*MessagesRequest)(nil),      // 4: MessagesRequest
	(*MessagesResponse)(nil),     // 5: MessagesResponse
	(*EventRequest)(nil),         // 6: EventRequest
	(*EventResponse)(nil),        // 7: EventResponse
}
var file_pb_poolpi_proto_depIdxs = []int32{
	3,  // 0: MonitorStateResponse.Heater1:type_name -> Indicator
	3,  // 1: MonitorStateResponse.Valve3:type_name -> Indicator
	3,  // 2: MonitorStateResponse.CheckSystem:type_name -> Indicator
	3,  // 3: MonitorStateResponse.Pool:type_name -> Indicator
	3,  // 4: MonitorStateResponse.Spa:type_name -> Indicator
	3,  // 5: MonitorStateResponse.Filter:type_name -> Indicator
	3,  // 6: MonitorStateResponse.Lights:type_name -> Indicator
	3,  // 7: MonitorStateResponse.Aux1:type_name -> Indicator
	3,  // 8: MonitorStateResponse.Aux2:type_name -> Indicator
	3,  // 9: MonitorStateResponse.Service:type_name -> Indicator
	3,  // 10: MonitorStateResponse.Aux3:type_name -> Indicator
	3,  // 11: MonitorStateResponse.Aux4:type_name -> Indicator
	3,  // 12: MonitorStateResponse.Aux5:type_name -> Indicator
	3,  // 13: MonitorStateResponse.Aux6:type_name -> Indicator
	3,  // 14: MonitorStateResponse.Valve4:type_name -> Indicator
	3,  // 15: MonitorStateResponse.Spillover:type_name -> Indicator
	3,  // 16: MonitorStateResponse.SystemOff:type_name -> Indicator
	3,  // 17: MonitorStateResponse.Aux7:type_name -> Indicator
	3,  // 18: MonitorStateResponse.Aux8:type_name -> Indicator
	3,  // 19: MonitorStateResponse.Aux9:type_name -> Indicator
	3,  // 20: MonitorStateResponse.Aux10:type_name -> Indicator
	3,  // 21: MonitorStateResponse.Aux11:type_name -> Indicator
	3,  // 22: MonitorStateResponse.Aux12:type_name -> Indicator
	3,  // 23: MonitorStateResponse.Aux13:type_name -> Indicator
	3,  // 24: MonitorStateResponse.Auz14:type_name -> Indicator
	3,  // 25: MonitorStateResponse.SuperChlorinate:type_name -> Indicator
	1,  // 26: Controller.MonitorState:input_type -> MonitorStateRequest
	4,  // 27: Controller.Messages:input_type -> MessagesRequest
	6,  // 28: Controller.Event:input_type -> EventRequest
	2,  // 29: Controller.MonitorState:output_type -> MonitorStateResponse
	5,  // 30: Controller.Messages:output_type -> MessagesResponse
	7,  // 31: Controller.Event:output_type -> EventResponse
	29, // [29:32] is the sub-list for method output_type
	26, // [26:29] is the sub-list for method input_type
	26, // [26:26] is the sub-list for extension type_name
	26, // [26:26] is the sub-list for extension extendee
	0,  // [0:26] is the sub-list for field type_name
}

func init() { file_pb_poolpi_proto_init() }
func file_pb_poolpi_proto_init() {
	if File_pb_poolpi_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_poolpi_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MonitorStateRequest); i {
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
		file_pb_poolpi_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MonitorStateResponse); i {
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
		file_pb_poolpi_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Indicator); i {
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
		file_pb_poolpi_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessagesRequest); i {
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
		file_pb_poolpi_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MessagesResponse); i {
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
		file_pb_poolpi_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventRequest); i {
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
		file_pb_poolpi_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EventResponse); i {
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
			RawDescriptor: file_pb_poolpi_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_poolpi_proto_goTypes,
		DependencyIndexes: file_pb_poolpi_proto_depIdxs,
		EnumInfos:         file_pb_poolpi_proto_enumTypes,
		MessageInfos:      file_pb_poolpi_proto_msgTypes,
	}.Build()
	File_pb_poolpi_proto = out.File
	file_pb_poolpi_proto_rawDesc = nil
	file_pb_poolpi_proto_goTypes = nil
	file_pb_poolpi_proto_depIdxs = nil
}