// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.23.0
// 	protoc        v3.14.0
// source: synchronization/core/archive.proto

package core

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// Archive is a wrapper around Entry that allows identification of non-existent
// roots when serializing. In-memory, a nil-Entry (that arrives without any
// error) represents an absence of content on the filesystem. Unfortunately,
// there is no way to represent that as an encoded message (an empty byte slice
// would successfully decode to an empty directory entry). By adding a level of
// indirection that allows for an unset root entry, we can encode Entry messages
// in a way that allows us to represent absence.
type Archive struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Root *Entry `protobuf:"bytes,1,opt,name=root,proto3" json:"root,omitempty"`
}

func (x *Archive) Reset() {
	*x = Archive{}
	if protoimpl.UnsafeEnabled {
		mi := &file_synchronization_core_archive_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Archive) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Archive) ProtoMessage() {}

func (x *Archive) ProtoReflect() protoreflect.Message {
	mi := &file_synchronization_core_archive_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Archive.ProtoReflect.Descriptor instead.
func (*Archive) Descriptor() ([]byte, []int) {
	return file_synchronization_core_archive_proto_rawDescGZIP(), []int{0}
}

func (x *Archive) GetRoot() *Entry {
	if x != nil {
		return x.Root
	}
	return nil
}

var File_synchronization_core_archive_proto protoreflect.FileDescriptor

var file_synchronization_core_archive_proto_rawDesc = []byte{
	0x0a, 0x22, 0x73, 0x79, 0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x2f, 0x63, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x04, 0x63, 0x6f, 0x72, 0x65, 0x1a, 0x20, 0x73, 0x79, 0x6e, 0x63,
	0x68, 0x72, 0x6f, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x63, 0x6f, 0x72, 0x65,
	0x2f, 0x65, 0x6e, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x2a, 0x0a, 0x07,
	0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x12, 0x1f, 0x0a, 0x04, 0x72, 0x6f, 0x6f, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x63, 0x6f, 0x72, 0x65, 0x2e, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x04, 0x72, 0x6f, 0x6f, 0x74, 0x42, 0x38, 0x5a, 0x36, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x6d, 0x75, 0x74, 0x61, 0x67, 0x65, 0x6e, 0x2d, 0x69,
	0x6f, 0x2f, 0x6d, 0x75, 0x74, 0x61, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x73, 0x79,
	0x6e, 0x63, 0x68, 0x72, 0x6f, 0x6e, 0x69, 0x7a, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x63, 0x6f,
	0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_synchronization_core_archive_proto_rawDescOnce sync.Once
	file_synchronization_core_archive_proto_rawDescData = file_synchronization_core_archive_proto_rawDesc
)

func file_synchronization_core_archive_proto_rawDescGZIP() []byte {
	file_synchronization_core_archive_proto_rawDescOnce.Do(func() {
		file_synchronization_core_archive_proto_rawDescData = protoimpl.X.CompressGZIP(file_synchronization_core_archive_proto_rawDescData)
	})
	return file_synchronization_core_archive_proto_rawDescData
}

var file_synchronization_core_archive_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_synchronization_core_archive_proto_goTypes = []interface{}{
	(*Archive)(nil), // 0: core.Archive
	(*Entry)(nil),   // 1: core.Entry
}
var file_synchronization_core_archive_proto_depIdxs = []int32{
	1, // 0: core.Archive.root:type_name -> core.Entry
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_synchronization_core_archive_proto_init() }
func file_synchronization_core_archive_proto_init() {
	if File_synchronization_core_archive_proto != nil {
		return
	}
	file_synchronization_core_entry_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_synchronization_core_archive_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Archive); i {
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
			RawDescriptor: file_synchronization_core_archive_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_synchronization_core_archive_proto_goTypes,
		DependencyIndexes: file_synchronization_core_archive_proto_depIdxs,
		MessageInfos:      file_synchronization_core_archive_proto_msgTypes,
	}.Build()
	File_synchronization_core_archive_proto = out.File
	file_synchronization_core_archive_proto_rawDesc = nil
	file_synchronization_core_archive_proto_goTypes = nil
	file_synchronization_core_archive_proto_depIdxs = nil
}
