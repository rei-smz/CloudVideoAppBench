// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.2
// source: recording.proto

package RecordingControl

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

type StartRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FileName string `protobuf:"bytes,1,opt,name=fileName,proto3" json:"fileName,omitempty"`
}

func (x *StartRequest) Reset() {
	*x = StartRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recording_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StartRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StartRequest) ProtoMessage() {}

func (x *StartRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recording_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StartRequest.ProtoReflect.Descriptor instead.
func (*StartRequest) Descriptor() ([]byte, []int) {
	return file_recording_proto_rawDescGZIP(), []int{0}
}

func (x *StartRequest) GetFileName() string {
	if x != nil {
		return x.FileName
	}
	return ""
}

type StopRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *StopRequest) Reset() {
	*x = StopRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recording_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *StopRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StopRequest) ProtoMessage() {}

func (x *StopRequest) ProtoReflect() protoreflect.Message {
	mi := &file_recording_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StopRequest.ProtoReflect.Descriptor instead.
func (*StopRequest) Descriptor() ([]byte, []int) {
	return file_recording_proto_rawDescGZIP(), []int{1}
}

type ServerReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success bool `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
}

func (x *ServerReply) Reset() {
	*x = ServerReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_recording_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ServerReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServerReply) ProtoMessage() {}

func (x *ServerReply) ProtoReflect() protoreflect.Message {
	mi := &file_recording_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ServerReply.ProtoReflect.Descriptor instead.
func (*ServerReply) Descriptor() ([]byte, []int) {
	return file_recording_proto_rawDescGZIP(), []int{2}
}

func (x *ServerReply) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

var File_recording_proto protoreflect.FileDescriptor

var file_recording_proto_rawDesc = []byte{
	0x0a, 0x0f, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x09, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x22, 0x2a, 0x0a, 0x0c,
	0x53, 0x74, 0x61, 0x72, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1a, 0x0a, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x66, 0x69, 0x6c, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x0d, 0x0a, 0x0b, 0x53, 0x74, 0x6f, 0x70,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x22, 0x27, 0x0a, 0x0b, 0x53, 0x65, 0x72, 0x76, 0x65,
	0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73,
	0x32, 0x87, 0x01, 0x0a, 0x0f, 0x52, 0x65, 0x63, 0x6f, 0x64, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e,
	0x74, 0x72, 0x6f, 0x6c, 0x12, 0x3a, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12, 0x17, 0x2e,
	0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x74, 0x61, 0x72, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69,
	0x6e, 0x67, 0x2e, 0x53, 0x65, 0x72, 0x76, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00,
	0x12, 0x38, 0x0a, 0x04, 0x73, 0x74, 0x6f, 0x70, 0x12, 0x16, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72,
	0x64, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x74, 0x6f, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x16, 0x2e, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x2e, 0x53, 0x65, 0x72,
	0x76, 0x65, 0x72, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x42, 0x14, 0x5a, 0x12, 0x2e, 0x2f,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x69, 0x6e, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_recording_proto_rawDescOnce sync.Once
	file_recording_proto_rawDescData = file_recording_proto_rawDesc
)

func file_recording_proto_rawDescGZIP() []byte {
	file_recording_proto_rawDescOnce.Do(func() {
		file_recording_proto_rawDescData = protoimpl.X.CompressGZIP(file_recording_proto_rawDescData)
	})
	return file_recording_proto_rawDescData
}

var file_recording_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_recording_proto_goTypes = []any{
	(*StartRequest)(nil), // 0: recording.StartRequest
	(*StopRequest)(nil),  // 1: recording.StopRequest
	(*ServerReply)(nil),  // 2: recording.ServerReply
}
var file_recording_proto_depIdxs = []int32{
	0, // 0: recording.RecodingControl.start:input_type -> recording.StartRequest
	1, // 1: recording.RecodingControl.stop:input_type -> recording.StopRequest
	2, // 2: recording.RecodingControl.start:output_type -> recording.ServerReply
	2, // 3: recording.RecodingControl.stop:output_type -> recording.ServerReply
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_recording_proto_init() }
func file_recording_proto_init() {
	if File_recording_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_recording_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*StartRequest); i {
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
		file_recording_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*StopRequest); i {
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
		file_recording_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*ServerReply); i {
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
			RawDescriptor: file_recording_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_recording_proto_goTypes,
		DependencyIndexes: file_recording_proto_depIdxs,
		MessageInfos:      file_recording_proto_msgTypes,
	}.Build()
	File_recording_proto = out.File
	file_recording_proto_rawDesc = nil
	file_recording_proto_goTypes = nil
	file_recording_proto_depIdxs = nil
}
