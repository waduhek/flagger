// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.32.0
// 	protoc        v4.25.1
// source: proto/flagpb/flag.proto

package flagpb

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

// CreateFlagRequest is the request body for creating a new flag.
type CreateFlagRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// A unique name in the project for the flag.
	FlagName string `protobuf:"bytes,1,opt,name=flag_name,json=flagName,proto3" json:"flag_name,omitempty"`
	// The name of the project under which the flag is to be created.
	ProjectName string `protobuf:"bytes,2,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty"`
}

func (x *CreateFlagRequest) Reset() {
	*x = CreateFlagRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_flagpb_flag_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFlagRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFlagRequest) ProtoMessage() {}

func (x *CreateFlagRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_flagpb_flag_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFlagRequest.ProtoReflect.Descriptor instead.
func (*CreateFlagRequest) Descriptor() ([]byte, []int) {
	return file_proto_flagpb_flag_proto_rawDescGZIP(), []int{0}
}

func (x *CreateFlagRequest) GetFlagName() string {
	if x != nil {
		return x.FlagName
	}
	return ""
}

func (x *CreateFlagRequest) GetProjectName() string {
	if x != nil {
		return x.ProjectName
	}
	return ""
}

// CreateFlagResponse is the response for creating a new flag.
type CreateFlagResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *CreateFlagResponse) Reset() {
	*x = CreateFlagResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_flagpb_flag_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateFlagResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateFlagResponse) ProtoMessage() {}

func (x *CreateFlagResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_flagpb_flag_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateFlagResponse.ProtoReflect.Descriptor instead.
func (*CreateFlagResponse) Descriptor() ([]byte, []int) {
	return file_proto_flagpb_flag_proto_rawDescGZIP(), []int{1}
}

// UpdateFlagStatusRequest is the request body to update the status of a flag.
type UpdateFlagStatusRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The name of the project where the flag is created.
	ProjectName string `protobuf:"bytes,1,opt,name=project_name,json=projectName,proto3" json:"project_name,omitempty"`
	// The name of the environment in which the flag is to be updated.
	EnvironmentName string `protobuf:"bytes,2,opt,name=environment_name,json=environmentName,proto3" json:"environment_name,omitempty"`
	// The name of the flag to be updated.
	FlagName string `protobuf:"bytes,3,opt,name=flag_name,json=flagName,proto3" json:"flag_name,omitempty"`
	// The update to be made to the flag.
	IsActive bool `protobuf:"varint,4,opt,name=is_active,json=isActive,proto3" json:"is_active,omitempty"`
}

func (x *UpdateFlagStatusRequest) Reset() {
	*x = UpdateFlagStatusRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_flagpb_flag_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateFlagStatusRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateFlagStatusRequest) ProtoMessage() {}

func (x *UpdateFlagStatusRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_flagpb_flag_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateFlagStatusRequest.ProtoReflect.Descriptor instead.
func (*UpdateFlagStatusRequest) Descriptor() ([]byte, []int) {
	return file_proto_flagpb_flag_proto_rawDescGZIP(), []int{2}
}

func (x *UpdateFlagStatusRequest) GetProjectName() string {
	if x != nil {
		return x.ProjectName
	}
	return ""
}

func (x *UpdateFlagStatusRequest) GetEnvironmentName() string {
	if x != nil {
		return x.EnvironmentName
	}
	return ""
}

func (x *UpdateFlagStatusRequest) GetFlagName() string {
	if x != nil {
		return x.FlagName
	}
	return ""
}

func (x *UpdateFlagStatusRequest) GetIsActive() bool {
	if x != nil {
		return x.IsActive
	}
	return false
}

// UpdateFlagStatusResponse is the response for updating a flag.
type UpdateFlagStatusResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *UpdateFlagStatusResponse) Reset() {
	*x = UpdateFlagStatusResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_flagpb_flag_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateFlagStatusResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateFlagStatusResponse) ProtoMessage() {}

func (x *UpdateFlagStatusResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_flagpb_flag_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateFlagStatusResponse.ProtoReflect.Descriptor instead.
func (*UpdateFlagStatusResponse) Descriptor() ([]byte, []int) {
	return file_proto_flagpb_flag_proto_rawDescGZIP(), []int{3}
}

var File_proto_flagpb_flag_proto protoreflect.FileDescriptor

var file_proto_flagpb_flag_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x6c, 0x61, 0x67, 0x70, 0x62, 0x2f, 0x66,
	0x6c, 0x61, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x66, 0x6c, 0x61, 0x67, 0x70,
	0x62, 0x22, 0x53, 0x0a, 0x11, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x6c, 0x61, 0x67, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x6c, 0x61, 0x67, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x6c, 0x61, 0x67, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x72, 0x6f, 0x6a, 0x65,
	0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x14, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x46, 0x6c, 0x61, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0xa1, 0x01, 0x0a,
	0x17, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x6c, 0x61, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x21, 0x0a, 0x0c, 0x70, 0x72, 0x6f, 0x6a,
	0x65, 0x63, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x29, 0x0a, 0x10, 0x65,
	0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65, 0x6e, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0f, 0x65, 0x6e, 0x76, 0x69, 0x72, 0x6f, 0x6e, 0x6d, 0x65,
	0x6e, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x66, 0x6c, 0x61, 0x67, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x6c, 0x61, 0x67, 0x4e,
	0x61, 0x6d, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x41, 0x63, 0x74, 0x69, 0x76, 0x65,
	0x22, 0x1a, 0x0a, 0x18, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x6c, 0x61, 0x67, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x32, 0xa2, 0x01, 0x0a,
	0x04, 0x46, 0x6c, 0x61, 0x67, 0x12, 0x43, 0x0a, 0x0a, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46,
	0x6c, 0x61, 0x67, 0x12, 0x19, 0x2e, 0x66, 0x6c, 0x61, 0x67, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x46, 0x6c, 0x61, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a,
	0x2e, 0x66, 0x6c, 0x61, 0x67, 0x70, 0x62, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x46, 0x6c,
	0x61, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x55, 0x0a, 0x10, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x46, 0x6c, 0x61, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f,
	0x2e, 0x66, 0x6c, 0x61, 0x67, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46, 0x6c,
	0x61, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x20, 0x2e, 0x66, 0x6c, 0x61, 0x67, 0x70, 0x62, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x46,
	0x6c, 0x61, 0x67, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x77, 0x61, 0x64, 0x75, 0x68, 0x65, 0x6b, 0x2f, 0x66, 0x6c, 0x61, 0x67, 0x67, 0x65, 0x72, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x66, 0x6c, 0x61, 0x67, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_flagpb_flag_proto_rawDescOnce sync.Once
	file_proto_flagpb_flag_proto_rawDescData = file_proto_flagpb_flag_proto_rawDesc
)

func file_proto_flagpb_flag_proto_rawDescGZIP() []byte {
	file_proto_flagpb_flag_proto_rawDescOnce.Do(func() {
		file_proto_flagpb_flag_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_flagpb_flag_proto_rawDescData)
	})
	return file_proto_flagpb_flag_proto_rawDescData
}

var file_proto_flagpb_flag_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_flagpb_flag_proto_goTypes = []interface{}{
	(*CreateFlagRequest)(nil),        // 0: flagpb.CreateFlagRequest
	(*CreateFlagResponse)(nil),       // 1: flagpb.CreateFlagResponse
	(*UpdateFlagStatusRequest)(nil),  // 2: flagpb.UpdateFlagStatusRequest
	(*UpdateFlagStatusResponse)(nil), // 3: flagpb.UpdateFlagStatusResponse
}
var file_proto_flagpb_flag_proto_depIdxs = []int32{
	0, // 0: flagpb.Flag.CreateFlag:input_type -> flagpb.CreateFlagRequest
	2, // 1: flagpb.Flag.UpdateFlagStatus:input_type -> flagpb.UpdateFlagStatusRequest
	1, // 2: flagpb.Flag.CreateFlag:output_type -> flagpb.CreateFlagResponse
	3, // 3: flagpb.Flag.UpdateFlagStatus:output_type -> flagpb.UpdateFlagStatusResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_proto_flagpb_flag_proto_init() }
func file_proto_flagpb_flag_proto_init() {
	if File_proto_flagpb_flag_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_flagpb_flag_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateFlagRequest); i {
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
		file_proto_flagpb_flag_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateFlagResponse); i {
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
		file_proto_flagpb_flag_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateFlagStatusRequest); i {
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
		file_proto_flagpb_flag_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateFlagStatusResponse); i {
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
			RawDescriptor: file_proto_flagpb_flag_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_flagpb_flag_proto_goTypes,
		DependencyIndexes: file_proto_flagpb_flag_proto_depIdxs,
		MessageInfos:      file_proto_flagpb_flag_proto_msgTypes,
	}.Build()
	File_proto_flagpb_flag_proto = out.File
	file_proto_flagpb_flag_proto_rawDesc = nil
	file_proto_flagpb_flag_proto_goTypes = nil
	file_proto_flagpb_flag_proto_depIdxs = nil
}
