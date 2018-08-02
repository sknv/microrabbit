// Code generated by protoc-gen-go. DO NOT EDIT.
// source: errors.proto

package rmq

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type Error struct {
	Code                 uint32            `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string            `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Meta                 map[string]string `protobuf:"bytes,3,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Error) Reset()         { *m = Error{} }
func (m *Error) String() string { return proto.CompactTextString(m) }
func (*Error) ProtoMessage()    {}
func (*Error) Descriptor() ([]byte, []int) {
	return fileDescriptor_errors_f937a0284c00b6aa, []int{0}
}
func (m *Error) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Error.Unmarshal(m, b)
}
func (m *Error) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Error.Marshal(b, m, deterministic)
}
func (dst *Error) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Error.Merge(dst, src)
}
func (m *Error) XXX_Size() int {
	return xxx_messageInfo_Error.Size(m)
}
func (m *Error) XXX_DiscardUnknown() {
	xxx_messageInfo_Error.DiscardUnknown(m)
}

var xxx_messageInfo_Error proto.InternalMessageInfo

func (m *Error) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Error) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Error) GetMeta() map[string]string {
	if m != nil {
		return m.Meta
	}
	return nil
}

func init() {
	proto.RegisterType((*Error)(nil), "rmq.Error")
	proto.RegisterMapType((map[string]string)(nil), "rmq.Error.MetaEntry")
}

func init() { proto.RegisterFile("errors.proto", fileDescriptor_errors_f937a0284c00b6aa) }

var fileDescriptor_errors_f937a0284c00b6aa = []byte{
	// 169 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x49, 0x2d, 0x2a, 0xca,
	0x2f, 0x2a, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2e, 0xca, 0x2d, 0x54, 0x9a, 0xc1,
	0xc8, 0xc5, 0xea, 0x0a, 0x12, 0x15, 0x12, 0xe2, 0x62, 0x49, 0xce, 0x4f, 0x49, 0x95, 0x60, 0x54,
	0x60, 0xd4, 0xe0, 0x0d, 0x02, 0xb3, 0x85, 0x24, 0xb8, 0xd8, 0x73, 0x53, 0x8b, 0x8b, 0x13, 0xd3,
	0x53, 0x25, 0x98, 0x14, 0x18, 0x35, 0x38, 0x83, 0x60, 0x5c, 0x21, 0x0d, 0x2e, 0x96, 0xdc, 0xd4,
	0x92, 0x44, 0x09, 0x66, 0x05, 0x66, 0x0d, 0x6e, 0x23, 0x11, 0xbd, 0xa2, 0xdc, 0x42, 0x3d, 0xb0,
	0x39, 0x7a, 0xbe, 0xa9, 0x25, 0x89, 0xae, 0x79, 0x25, 0x45, 0x95, 0x41, 0x60, 0x15, 0x52, 0xe6,
	0x5c, 0x9c, 0x70, 0x21, 0x21, 0x01, 0x2e, 0xe6, 0xec, 0xd4, 0x4a, 0xb0, 0x1d, 0x9c, 0x41, 0x20,
	0xa6, 0x90, 0x08, 0x17, 0x6b, 0x59, 0x62, 0x4e, 0x29, 0xcc, 0x02, 0x08, 0xc7, 0x8a, 0xc9, 0x82,
	0xd1, 0x89, 0x35, 0x0a, 0xe4, 0xc2, 0x24, 0x36, 0xb0, 0x6b, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff,
	0xff, 0x1c, 0x15, 0xa7, 0xd8, 0xbd, 0x00, 0x00, 0x00,
}
