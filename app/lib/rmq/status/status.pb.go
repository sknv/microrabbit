// Code generated by protoc-gen-go. DO NOT EDIT.
// source: status.proto

package status

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

type Status struct {
	Code                 uint32            `protobuf:"varint,1,opt,name=code,proto3" json:"code,omitempty"`
	Message              string            `protobuf:"bytes,2,opt,name=message,proto3" json:"message,omitempty"`
	Meta                 map[string]string `protobuf:"bytes,3,rep,name=meta,proto3" json:"meta,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *Status) Reset()         { *m = Status{} }
func (m *Status) String() string { return proto.CompactTextString(m) }
func (*Status) ProtoMessage()    {}
func (*Status) Descriptor() ([]byte, []int) {
	return fileDescriptor_status_09744f5b70125540, []int{0}
}
func (m *Status) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Status.Unmarshal(m, b)
}
func (m *Status) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Status.Marshal(b, m, deterministic)
}
func (dst *Status) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Status.Merge(dst, src)
}
func (m *Status) XXX_Size() int {
	return xxx_messageInfo_Status.Size(m)
}
func (m *Status) XXX_DiscardUnknown() {
	xxx_messageInfo_Status.DiscardUnknown(m)
}

var xxx_messageInfo_Status proto.InternalMessageInfo

func (m *Status) GetCode() uint32 {
	if m != nil {
		return m.Code
	}
	return 0
}

func (m *Status) GetMessage() string {
	if m != nil {
		return m.Message
	}
	return ""
}

func (m *Status) GetMeta() map[string]string {
	if m != nil {
		return m.Meta
	}
	return nil
}

func init() {
	proto.RegisterType((*Status)(nil), "rmq.status.Status")
	proto.RegisterMapType((map[string]string)(nil), "rmq.status.Status.MetaEntry")
}

func init() { proto.RegisterFile("status.proto", fileDescriptor_status_09744f5b70125540) }

var fileDescriptor_status_09744f5b70125540 = []byte{
	// 172 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x29, 0x2e, 0x49, 0x2c,
	0x29, 0x2d, 0xd6, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x2a, 0xca, 0x2d, 0xd4, 0x83, 0x88,
	0x28, 0x2d, 0x64, 0xe4, 0x62, 0x0b, 0x06, 0x33, 0x85, 0x84, 0xb8, 0x58, 0x92, 0xf3, 0x53, 0x52,
	0x25, 0x18, 0x15, 0x18, 0x35, 0x78, 0x83, 0xc0, 0x6c, 0x21, 0x09, 0x2e, 0xf6, 0xdc, 0xd4, 0xe2,
	0xe2, 0xc4, 0xf4, 0x54, 0x09, 0x26, 0x05, 0x46, 0x0d, 0xce, 0x20, 0x18, 0x57, 0xc8, 0x80, 0x8b,
	0x25, 0x37, 0xb5, 0x24, 0x51, 0x82, 0x59, 0x81, 0x59, 0x83, 0xdb, 0x48, 0x46, 0x0f, 0x61, 0xa6,
	0x1e, 0xc4, 0x3c, 0x3d, 0xdf, 0xd4, 0x92, 0x44, 0xd7, 0xbc, 0x92, 0xa2, 0xca, 0x20, 0xb0, 0x4a,
	0x29, 0x73, 0x2e, 0x4e, 0xb8, 0x90, 0x90, 0x00, 0x17, 0x73, 0x76, 0x6a, 0x25, 0xd8, 0x2e, 0xce,
	0x20, 0x10, 0x53, 0x48, 0x84, 0x8b, 0xb5, 0x2c, 0x31, 0xa7, 0x14, 0x66, 0x11, 0x84, 0x63, 0xc5,
	0x64, 0xc1, 0xe8, 0xc4, 0x11, 0xc5, 0x06, 0x31, 0x39, 0x89, 0x0d, 0xec, 0x01, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x81, 0x53, 0xa3, 0xfb, 0xd0, 0x00, 0x00, 0x00,
}
