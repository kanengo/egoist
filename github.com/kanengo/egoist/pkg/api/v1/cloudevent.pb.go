// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.14.0
// source: api/v1/cloudevent.proto

package api

import (
	_struct "github.com/golang/protobuf/ptypes/struct"
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

type CloudEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              string                    `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Data            []byte                    `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Source          string                    `protobuf:"bytes,3,opt,name=source,proto3" json:"source,omitempty"`
	SpecVersion     string                    `protobuf:"bytes,4,opt,name=spec_version,json=specVersion,proto3" json:"spec_version,omitempty"`
	Type            string                    `protobuf:"bytes,5,opt,name=type,proto3" json:"type,omitempty"`
	DataContentType string                    `protobuf:"bytes,6,opt,name=data_content_type,json=dataContentType,proto3" json:"data_content_type,omitempty"`
	Timestamp       int64                     `protobuf:"varint,7,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	Extensions      map[string]*_struct.Value `protobuf:"bytes,9,rep,name=extensions,proto3" json:"extensions,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *CloudEvent) Reset() {
	*x = CloudEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_cloudevent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CloudEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CloudEvent) ProtoMessage() {}

func (x *CloudEvent) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_cloudevent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CloudEvent.ProtoReflect.Descriptor instead.
func (*CloudEvent) Descriptor() ([]byte, []int) {
	return file_api_v1_cloudevent_proto_rawDescGZIP(), []int{0}
}

func (x *CloudEvent) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *CloudEvent) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *CloudEvent) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *CloudEvent) GetSpecVersion() string {
	if x != nil {
		return x.SpecVersion
	}
	return ""
}

func (x *CloudEvent) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *CloudEvent) GetDataContentType() string {
	if x != nil {
		return x.DataContentType
	}
	return ""
}

func (x *CloudEvent) GetTimestamp() int64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

func (x *CloudEvent) GetExtensions() map[string]*_struct.Value {
	if x != nil {
		return x.Extensions
	}
	return nil
}

var File_api_v1_cloudevent_proto protoreflect.FileDescriptor

var file_api_v1_cloudevent_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6c, 0x6f, 0x75, 0x64, 0x65, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x1a, 0x1c, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62,
	0x75, 0x66, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0xe4, 0x02, 0x0a, 0x0a, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x0e,
	0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x12,
	0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x73, 0x70,
	0x65, 0x63, 0x5f, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x0b, 0x73, 0x70, 0x65, 0x63, 0x56, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x2a, 0x0a, 0x11, 0x64, 0x61, 0x74, 0x61, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e,
	0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x64, 0x61,
	0x74, 0x61, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x07, 0x20, 0x01, 0x28, 0x03,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x12, 0x42, 0x0a, 0x0a, 0x65,
	0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x09, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x22, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x43, 0x6c, 0x6f, 0x75, 0x64, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x1a,
	0x55, 0x0a, 0x0f, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x73, 0x69, 0x6f, 0x6e, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x2c, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x4a, 0x5a, 0x28, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6b, 0x61, 0x6e, 0x65, 0x6e, 0x67, 0x6f, 0x2f, 0x65, 0x67, 0x6f,
	0x69, 0x73, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x3b, 0x61,
	0x70, 0x69, 0xaa, 0x02, 0x1d, 0x45, 0x67, 0x6f, 0x69, 0x73, 0x74, 0x2e, 0x43, 0x6c, 0x69, 0x65,
	0x6e, 0x74, 0x2e, 0x41, 0x75, 0x74, 0x6f, 0x67, 0x65, 0x6e, 0x2e, 0x47, 0x72, 0x70, 0x63, 0x2e,
	0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_cloudevent_proto_rawDescOnce sync.Once
	file_api_v1_cloudevent_proto_rawDescData = file_api_v1_cloudevent_proto_rawDesc
)

func file_api_v1_cloudevent_proto_rawDescGZIP() []byte {
	file_api_v1_cloudevent_proto_rawDescOnce.Do(func() {
		file_api_v1_cloudevent_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_cloudevent_proto_rawDescData)
	})
	return file_api_v1_cloudevent_proto_rawDescData
}

var file_api_v1_cloudevent_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_api_v1_cloudevent_proto_goTypes = []interface{}{
	(*CloudEvent)(nil),    // 0: api.v1.CloudEvent
	nil,                   // 1: api.v1.CloudEvent.ExtensionsEntry
	(*_struct.Value)(nil), // 2: google.protobuf.Value
}
var file_api_v1_cloudevent_proto_depIdxs = []int32{
	1, // 0: api.v1.CloudEvent.extensions:type_name -> api.v1.CloudEvent.ExtensionsEntry
	2, // 1: api.v1.CloudEvent.ExtensionsEntry.value:type_name -> google.protobuf.Value
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_api_v1_cloudevent_proto_init() }
func file_api_v1_cloudevent_proto_init() {
	if File_api_v1_cloudevent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_api_v1_cloudevent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CloudEvent); i {
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
			RawDescriptor: file_api_v1_cloudevent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_api_v1_cloudevent_proto_goTypes,
		DependencyIndexes: file_api_v1_cloudevent_proto_depIdxs,
		MessageInfos:      file_api_v1_cloudevent_proto_msgTypes,
	}.Build()
	File_api_v1_cloudevent_proto = out.File
	file_api_v1_cloudevent_proto_rawDesc = nil
	file_api_v1_cloudevent_proto_goTypes = nil
	file_api_v1_cloudevent_proto_depIdxs = nil
}
