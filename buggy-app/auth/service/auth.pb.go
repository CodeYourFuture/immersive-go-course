// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.7
// source: auth/service/auth.proto

package service

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

type Input struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *Input) Reset() {
	*x = Input{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_service_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Input) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Input) ProtoMessage() {}

func (x *Input) ProtoReflect() protoreflect.Message {
	mi := &file_auth_service_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Input.ProtoReflect.Descriptor instead.
func (*Input) Descriptor() ([]byte, []int) {
	return file_auth_service_auth_proto_rawDescGZIP(), []int{0}
}

func (x *Input) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Input) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type Result struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Allow bool `protobuf:"varint,1,opt,name=allow,proto3" json:"allow,omitempty"`
}

func (x *Result) Reset() {
	*x = Result{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_service_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Result) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Result) ProtoMessage() {}

func (x *Result) ProtoReflect() protoreflect.Message {
	mi := &file_auth_service_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Result.ProtoReflect.Descriptor instead.
func (*Result) Descriptor() ([]byte, []int) {
	return file_auth_service_auth_proto_rawDescGZIP(), []int{1}
}

func (x *Result) GetAllow() bool {
	if x != nil {
		return x.Allow
	}
	return false
}

var File_auth_service_auth_proto protoreflect.FileDescriptor

var file_auth_service_auth_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2f, 0x61,
	0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x73, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x22, 0x33, 0x0a, 0x05, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70,
	0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x1e, 0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x05, 0x61, 0x6c, 0x6c, 0x6f, 0x77, 0x32, 0x33, 0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12,
	0x2b, 0x0a, 0x06, 0x56, 0x65, 0x72, 0x69, 0x66, 0x79, 0x12, 0x0e, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x0f, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x42, 0x46, 0x5a, 0x44,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x43, 0x6f, 0x64, 0x65, 0x59,
	0x6f, 0x75, 0x72, 0x46, 0x75, 0x74, 0x75, 0x72, 0x65, 0x2f, 0x69, 0x6d, 0x6d, 0x65, 0x72, 0x73,
	0x69, 0x76, 0x65, 0x2d, 0x67, 0x6f, 0x2d, 0x63, 0x6f, 0x75, 0x72, 0x73, 0x65, 0x2f, 0x62, 0x75,
	0x67, 0x67, 0x79, 0x2d, 0x61, 0x70, 0x70, 0x2f, 0x61, 0x75, 0x74, 0x68, 0x2f, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_service_auth_proto_rawDescOnce sync.Once
	file_auth_service_auth_proto_rawDescData = file_auth_service_auth_proto_rawDesc
)

func file_auth_service_auth_proto_rawDescGZIP() []byte {
	file_auth_service_auth_proto_rawDescOnce.Do(func() {
		file_auth_service_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_service_auth_proto_rawDescData)
	})
	return file_auth_service_auth_proto_rawDescData
}

var file_auth_service_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_auth_service_auth_proto_goTypes = []interface{}{
	(*Input)(nil),  // 0: service.Input
	(*Result)(nil), // 1: service.Result
}
var file_auth_service_auth_proto_depIdxs = []int32{
	0, // 0: service.Auth.Verify:input_type -> service.Input
	1, // 1: service.Auth.Verify:output_type -> service.Result
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_auth_service_auth_proto_init() }
func file_auth_service_auth_proto_init() {
	if File_auth_service_auth_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_auth_service_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Input); i {
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
		file_auth_service_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Result); i {
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
			RawDescriptor: file_auth_service_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_service_auth_proto_goTypes,
		DependencyIndexes: file_auth_service_auth_proto_depIdxs,
		MessageInfos:      file_auth_service_auth_proto_msgTypes,
	}.Build()
	File_auth_service_auth_proto = out.File
	file_auth_service_auth_proto_rawDesc = nil
	file_auth_service_auth_proto_goTypes = nil
	file_auth_service_auth_proto_depIdxs = nil
}
