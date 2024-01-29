// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: proto/providerpb/flag_provider.proto

package providerpb

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

// GetFlagRequest is the request body for getting the current status of the
// flag.
type GetFlagRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the flag.
	FlagName string `protobuf:"bytes,2,opt,name=flag_name,json=flagName,proto3" json:"flag_name,omitempty"`
	// The environment to fetch the flag from.
	Environment string `protobuf:"bytes,1,opt,name=environment,proto3" json:"environment,omitempty"`
}

func (x *GetFlagRequest) Reset() {
	*x = GetFlagRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_providerpb_flag_provider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFlagRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFlagRequest) ProtoMessage() {}

func (x *GetFlagRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_providerpb_flag_provider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFlagRequest.ProtoReflect.Descriptor instead.
func (*GetFlagRequest) Descriptor() ([]byte, []int) {
	return file_proto_providerpb_flag_provider_proto_rawDescGZIP(), []int{0}
}

func (x *GetFlagRequest) GetFlagName() string {
	if x != nil {
		return x.FlagName
	}
	return ""
}

func (x *GetFlagRequest) GetEnvironment() string {
	if x != nil {
		return x.Environment
	}
	return ""
}

// GetFlagResponse is the response of getting the current staus of the flag.
type GetFlagResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The current status of the flag.
	Status bool `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
}

func (x *GetFlagResponse) Reset() {
	*x = GetFlagResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_providerpb_flag_provider_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetFlagResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetFlagResponse) ProtoMessage() {}

func (x *GetFlagResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_providerpb_flag_provider_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetFlagResponse.ProtoReflect.Descriptor instead.
func (*GetFlagResponse) Descriptor() ([]byte, []int) {
	return file_proto_providerpb_flag_provider_proto_rawDescGZIP(), []int{1}
}

func (x *GetFlagResponse) GetStatus() bool {
	if x != nil {
		return x.Status
	}
	return false
}

var File_proto_providerpb_flag_provider_proto protoreflect.FileDescriptor

var file_proto_providerpb_flag_provider_proto_rawDesc = []byte{
	0x0a, 0x24, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x70, 0x62, 0x2f, 0x66, 0x6c, 0x61, 0x67, 0x5f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0a, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x70, 0x62, 0x22, 0x4f, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x6c, 0x61, 0x67, 0x5f, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x6c, 0x61, 0x67, 0x4e, 0x61, 0x6d,
	0x65, 0x12, 0x20, 0x0a, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d,
	0x65, 0x6e, 0x74, 0x22, 0x29, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x32, 0x52,
	0x0a, 0x0c, 0x46, 0x6c, 0x61, 0x67, 0x50, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x12, 0x42,
	0x0a, 0x07, 0x47, 0x65, 0x74, 0x46, 0x6c, 0x61, 0x67, 0x12, 0x1a, 0x2e, 0x70, 0x72, 0x6f, 0x76,
	0x69, 0x64, 0x65, 0x72, 0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1b, 0x2e, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72,
	0x70, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x42, 0x2d, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d,
	0x2f, 0x77, 0x61, 0x64, 0x75, 0x68, 0x65, 0x6b, 0x2f, 0x66, 0x6c, 0x61, 0x67, 0x67, 0x65, 0x72,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x70,
	0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_providerpb_flag_provider_proto_rawDescOnce sync.Once
	file_proto_providerpb_flag_provider_proto_rawDescData = file_proto_providerpb_flag_provider_proto_rawDesc
)

func file_proto_providerpb_flag_provider_proto_rawDescGZIP() []byte {
	file_proto_providerpb_flag_provider_proto_rawDescOnce.Do(func() {
		file_proto_providerpb_flag_provider_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_providerpb_flag_provider_proto_rawDescData)
	})
	return file_proto_providerpb_flag_provider_proto_rawDescData
}

var file_proto_providerpb_flag_provider_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_providerpb_flag_provider_proto_goTypes = []interface{}{
	(*GetFlagRequest)(nil),  // 0: providerpb.GetFlagRequest
	(*GetFlagResponse)(nil), // 1: providerpb.GetFlagResponse
}
var file_proto_providerpb_flag_provider_proto_depIdxs = []int32{
	0, // 0: providerpb.FlagProvider.GetFlag:input_type -> providerpb.GetFlagRequest
	1, // 1: providerpb.FlagProvider.GetFlag:output_type -> providerpb.GetFlagResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_providerpb_flag_provider_proto_init() }
func file_proto_providerpb_flag_provider_proto_init() {
	if File_proto_providerpb_flag_provider_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_providerpb_flag_provider_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFlagRequest); i {
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
		file_proto_providerpb_flag_provider_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetFlagResponse); i {
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
			RawDescriptor: file_proto_providerpb_flag_provider_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_providerpb_flag_provider_proto_goTypes,
		DependencyIndexes: file_proto_providerpb_flag_provider_proto_depIdxs,
		MessageInfos:      file_proto_providerpb_flag_provider_proto_msgTypes,
	}.Build()
	File_proto_providerpb_flag_provider_proto = out.File
	file_proto_providerpb_flag_provider_proto_rawDesc = nil
	file_proto_providerpb_flag_provider_proto_goTypes = nil
	file_proto_providerpb_flag_provider_proto_depIdxs = nil
}
