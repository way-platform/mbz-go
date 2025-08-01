// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: wayplatform/mbz/v1/push_message.proto

package mbzv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// A push message from the Mercedes-Benz Kafka push API.
type PushMessage struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Unique message identifier.
	MessageId string `protobuf:"bytes,1,opt,name=message_id,json=messageId,proto3" json:"message_id,omitempty"`
	// Vehicle identification number (VIN).
	Vin string `protobuf:"bytes,2,opt,name=vin,proto3" json:"vin,omitempty"`
	// Time when the message was created (in microseconds since Unix epoch).
	Time int64 `protobuf:"varint,3,opt,name=time,proto3" json:"time,omitempty"`
	// Message type.
	MessageType MessageType `protobuf:"varint,4,opt,name=message_type,json=messageType,proto3,enum=wayplatform.mbz.v1.MessageType" json:"message_type,omitempty"`
	// Version tag for the message.
	Version string `protobuf:"bytes,5,opt,name=version,proto3" json:"version,omitempty"`
	// Service associated with the message.
	ServiceId string `protobuf:"bytes,6,opt,name=service_id,json=serviceId,proto3" json:"service_id,omitempty"`
	// Sending behavior.
	SendingBehavior SendingBehavior `protobuf:"varint,7,opt,name=sending_behavior,json=sendingBehavior,proto3,enum=wayplatform.mbz.v1.SendingBehavior" json:"sending_behavior,omitempty"`
	// Signals (valid for SIGNALS message type).
	Signals       []*Signal `protobuf:"bytes,8,rep,name=signals,proto3" json:"signals,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *PushMessage) Reset() {
	*x = PushMessage{}
	mi := &file_wayplatform_mbz_v1_push_message_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PushMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMessage) ProtoMessage() {}

func (x *PushMessage) ProtoReflect() protoreflect.Message {
	mi := &file_wayplatform_mbz_v1_push_message_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMessage.ProtoReflect.Descriptor instead.
func (*PushMessage) Descriptor() ([]byte, []int) {
	return file_wayplatform_mbz_v1_push_message_proto_rawDescGZIP(), []int{0}
}

func (x *PushMessage) GetMessageId() string {
	if x != nil {
		return x.MessageId
	}
	return ""
}

func (x *PushMessage) GetVin() string {
	if x != nil {
		return x.Vin
	}
	return ""
}

func (x *PushMessage) GetTime() int64 {
	if x != nil {
		return x.Time
	}
	return 0
}

func (x *PushMessage) GetMessageType() MessageType {
	if x != nil {
		return x.MessageType
	}
	return MessageType_MESSAGE_TYPE_UNSPECIFIED
}

func (x *PushMessage) GetVersion() string {
	if x != nil {
		return x.Version
	}
	return ""
}

func (x *PushMessage) GetServiceId() string {
	if x != nil {
		return x.ServiceId
	}
	return ""
}

func (x *PushMessage) GetSendingBehavior() SendingBehavior {
	if x != nil {
		return x.SendingBehavior
	}
	return SendingBehavior_SENDING_BEHAVIOR_UNSPECIFIED
}

func (x *PushMessage) GetSignals() []*Signal {
	if x != nil {
		return x.Signals
	}
	return nil
}

var File_wayplatform_mbz_v1_push_message_proto protoreflect.FileDescriptor

const file_wayplatform_mbz_v1_push_message_proto_rawDesc = "" +
	"\n" +
	"%wayplatform/mbz/v1/push_message.proto\x12\x12wayplatform.mbz.v1\x1a%wayplatform/mbz/v1/message_type.proto\x1a)wayplatform/mbz/v1/sending_behavior.proto\x1a\x1fwayplatform/mbz/v1/signal.proto\"\xd5\x02\n" +
	"\vPushMessage\x12\x1d\n" +
	"\n" +
	"message_id\x18\x01 \x01(\tR\tmessageId\x12\x10\n" +
	"\x03vin\x18\x02 \x01(\tR\x03vin\x12\x12\n" +
	"\x04time\x18\x03 \x01(\x03R\x04time\x12B\n" +
	"\fmessage_type\x18\x04 \x01(\x0e2\x1f.wayplatform.mbz.v1.MessageTypeR\vmessageType\x12\x18\n" +
	"\aversion\x18\x05 \x01(\tR\aversion\x12\x1d\n" +
	"\n" +
	"service_id\x18\x06 \x01(\tR\tserviceId\x12N\n" +
	"\x10sending_behavior\x18\a \x01(\x0e2#.wayplatform.mbz.v1.SendingBehaviorR\x0fsendingBehavior\x124\n" +
	"\asignals\x18\b \x03(\v2\x1a.wayplatform.mbz.v1.SignalR\asignalsB\xda\x01\n" +
	"\x16com.wayplatform.mbz.v1B\x10PushMessageProtoP\x01ZDgithub.com/way-platform/mbz-go/proto/gen/go/wayplatform/mbz/v1;mbzv1\xa2\x02\x03WMX\xaa\x02\x12Wayplatform.Mbz.V1\xca\x02\x12Wayplatform\\Mbz\\V1\xe2\x02\x1eWayplatform\\Mbz\\V1\\GPBMetadata\xea\x02\x14Wayplatform::Mbz::V1b\x06proto3"

var (
	file_wayplatform_mbz_v1_push_message_proto_rawDescOnce sync.Once
	file_wayplatform_mbz_v1_push_message_proto_rawDescData []byte
)

func file_wayplatform_mbz_v1_push_message_proto_rawDescGZIP() []byte {
	file_wayplatform_mbz_v1_push_message_proto_rawDescOnce.Do(func() {
		file_wayplatform_mbz_v1_push_message_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_wayplatform_mbz_v1_push_message_proto_rawDesc), len(file_wayplatform_mbz_v1_push_message_proto_rawDesc)))
	})
	return file_wayplatform_mbz_v1_push_message_proto_rawDescData
}

var file_wayplatform_mbz_v1_push_message_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_wayplatform_mbz_v1_push_message_proto_goTypes = []any{
	(*PushMessage)(nil),  // 0: wayplatform.mbz.v1.PushMessage
	(MessageType)(0),     // 1: wayplatform.mbz.v1.MessageType
	(SendingBehavior)(0), // 2: wayplatform.mbz.v1.SendingBehavior
	(*Signal)(nil),       // 3: wayplatform.mbz.v1.Signal
}
var file_wayplatform_mbz_v1_push_message_proto_depIdxs = []int32{
	1, // 0: wayplatform.mbz.v1.PushMessage.message_type:type_name -> wayplatform.mbz.v1.MessageType
	2, // 1: wayplatform.mbz.v1.PushMessage.sending_behavior:type_name -> wayplatform.mbz.v1.SendingBehavior
	3, // 2: wayplatform.mbz.v1.PushMessage.signals:type_name -> wayplatform.mbz.v1.Signal
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_wayplatform_mbz_v1_push_message_proto_init() }
func file_wayplatform_mbz_v1_push_message_proto_init() {
	if File_wayplatform_mbz_v1_push_message_proto != nil {
		return
	}
	file_wayplatform_mbz_v1_message_type_proto_init()
	file_wayplatform_mbz_v1_sending_behavior_proto_init()
	file_wayplatform_mbz_v1_signal_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_wayplatform_mbz_v1_push_message_proto_rawDesc), len(file_wayplatform_mbz_v1_push_message_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_wayplatform_mbz_v1_push_message_proto_goTypes,
		DependencyIndexes: file_wayplatform_mbz_v1_push_message_proto_depIdxs,
		MessageInfos:      file_wayplatform_mbz_v1_push_message_proto_msgTypes,
	}.Build()
	File_wayplatform_mbz_v1_push_message_proto = out.File
	file_wayplatform_mbz_v1_push_message_proto_goTypes = nil
	file_wayplatform_mbz_v1_push_message_proto_depIdxs = nil
}
