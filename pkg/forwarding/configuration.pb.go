// Code generated by protoc-gen-go. DO NOT EDIT.
// source: forwarding/configuration.proto

package forwarding

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// Configuration encodes session configuration parameters. It is used for create
// commands to specify configuration options, for loading global configuration
// options, and for storing a merged configuration inside sessions.
type Configuration struct {
	// SocketOverwriteMode specifies whether or not existing Unix domain sockets
	// should be overwritten when creating new listener sockets.
	SocketOverwriteMode SocketOverwriteMode `protobuf:"varint,41,opt,name=socketOverwriteMode,proto3,enum=forwarding.SocketOverwriteMode" json:"socketOverwriteMode,omitempty"`
	// SocketOwner specifies the owner identifier to use for Unix domain
	// listener sockets.
	SocketOwner string `protobuf:"bytes,42,opt,name=socketOwner,proto3" json:"socketOwner,omitempty"`
	// SocketGroup specifies the group identifier to use for Unix domain
	// listener sockets.
	SocketGroup string `protobuf:"bytes,43,opt,name=socketGroup,proto3" json:"socketGroup,omitempty"`
	// SocketPermissionMode specifies the permission mode to use for Unix domain
	// listener sockets.
	SocketPermissionMode uint32   `protobuf:"varint,44,opt,name=socketPermissionMode,proto3" json:"socketPermissionMode,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Configuration) Reset()         { *m = Configuration{} }
func (m *Configuration) String() string { return proto.CompactTextString(m) }
func (*Configuration) ProtoMessage()    {}
func (*Configuration) Descriptor() ([]byte, []int) {
	return fileDescriptor_5e51e4766fb5528c, []int{0}
}

func (m *Configuration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Configuration.Unmarshal(m, b)
}
func (m *Configuration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Configuration.Marshal(b, m, deterministic)
}
func (m *Configuration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Configuration.Merge(m, src)
}
func (m *Configuration) XXX_Size() int {
	return xxx_messageInfo_Configuration.Size(m)
}
func (m *Configuration) XXX_DiscardUnknown() {
	xxx_messageInfo_Configuration.DiscardUnknown(m)
}

var xxx_messageInfo_Configuration proto.InternalMessageInfo

func (m *Configuration) GetSocketOverwriteMode() SocketOverwriteMode {
	if m != nil {
		return m.SocketOverwriteMode
	}
	return SocketOverwriteMode_SocketOverwriteModeDefault
}

func (m *Configuration) GetSocketOwner() string {
	if m != nil {
		return m.SocketOwner
	}
	return ""
}

func (m *Configuration) GetSocketGroup() string {
	if m != nil {
		return m.SocketGroup
	}
	return ""
}

func (m *Configuration) GetSocketPermissionMode() uint32 {
	if m != nil {
		return m.SocketPermissionMode
	}
	return 0
}

func init() {
	proto.RegisterType((*Configuration)(nil), "forwarding.Configuration")
}

func init() { proto.RegisterFile("forwarding/configuration.proto", fileDescriptor_5e51e4766fb5528c) }

var fileDescriptor_5e51e4766fb5528c = []byte{
	// 229 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0x4b, 0xcb, 0x2f, 0x2a,
	0x4f, 0x2c, 0x4a, 0xc9, 0xcc, 0x4b, 0xd7, 0x4f, 0xce, 0xcf, 0x4b, 0xcb, 0x4c, 0x2f, 0x2d, 0x4a,
	0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x42, 0xc8, 0x4b, 0xa9,
	0x21, 0xa9, 0x2d, 0xce, 0x4f, 0xce, 0x4e, 0x2d, 0x89, 0xcf, 0x2f, 0x4b, 0x2d, 0x2a, 0x2f, 0xca,
	0x2c, 0x49, 0x8d, 0xcf, 0xcd, 0x4f, 0x49, 0x85, 0xe8, 0x51, 0xba, 0xc5, 0xc8, 0xc5, 0xeb, 0x8c,
	0x6c, 0x96, 0x50, 0x20, 0x97, 0x30, 0x44, 0x83, 0x3f, 0x4c, 0xbd, 0x6f, 0x7e, 0x4a, 0xaa, 0x84,
	0xa6, 0x02, 0xa3, 0x06, 0x9f, 0x91, 0xbc, 0x1e, 0xc2, 0x5c, 0xbd, 0x60, 0x4c, 0x65, 0x41, 0xd8,
	0xf4, 0x0a, 0x29, 0x70, 0x71, 0x43, 0x85, 0xcb, 0xf3, 0x52, 0x8b, 0x24, 0xb4, 0x14, 0x18, 0x35,
	0x38, 0x83, 0x90, 0x85, 0x10, 0x2a, 0xdc, 0x8b, 0xf2, 0x4b, 0x0b, 0x24, 0xb4, 0x91, 0x55, 0x80,
	0x85, 0x84, 0x8c, 0xb8, 0x44, 0x20, 0xdc, 0x80, 0xd4, 0xa2, 0xdc, 0xcc, 0xe2, 0xe2, 0xcc, 0xfc,
	0x3c, 0xb0, 0xbb, 0x74, 0x14, 0x18, 0x35, 0x78, 0x83, 0xb0, 0xca, 0x39, 0xe9, 0x44, 0x69, 0xa5,
	0x67, 0x96, 0x64, 0x94, 0x26, 0xe9, 0x25, 0xe7, 0xe7, 0xea, 0x67, 0x24, 0x96, 0xe5, 0x27, 0xeb,
	0x66, 0xe6, 0xeb, 0xe7, 0x96, 0x96, 0x24, 0xa6, 0xa7, 0xe6, 0xe9, 0x17, 0x64, 0xa7, 0xeb, 0x23,
	0xbc, 0x93, 0xc4, 0x06, 0x0e, 0x11, 0x63, 0x40, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb3, 0xf4, 0xbc,
	0x18, 0x67, 0x01, 0x00, 0x00,
}