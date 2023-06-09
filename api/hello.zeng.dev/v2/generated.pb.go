/*
boilerplate text in generated file header
*/

// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2/generated.proto

package v2

import (
	fmt "fmt"

	io "io"

	proto "github.com/gogo/protobuf/proto"
	k8s_io_apimachinery_pkg_apis_meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	math "math"
	math_bits "math/bits"
	reflect "reflect"
	strings "strings"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

func (m *Foo) Reset()      { *m = Foo{} }
func (*Foo) ProtoMessage() {}
func (*Foo) Descriptor() ([]byte, []int) {
	return fileDescriptor_18c52dbbfb632f7c, []int{0}
}
func (m *Foo) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Foo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *Foo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Foo.Merge(m, src)
}
func (m *Foo) XXX_Size() int {
	return m.Size()
}
func (m *Foo) XXX_DiscardUnknown() {
	xxx_messageInfo_Foo.DiscardUnknown(m)
}

var xxx_messageInfo_Foo proto.InternalMessageInfo

func (m *FooCondition) Reset()      { *m = FooCondition{} }
func (*FooCondition) ProtoMessage() {}
func (*FooCondition) Descriptor() ([]byte, []int) {
	return fileDescriptor_18c52dbbfb632f7c, []int{1}
}
func (m *FooCondition) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FooCondition) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *FooCondition) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FooCondition.Merge(m, src)
}
func (m *FooCondition) XXX_Size() int {
	return m.Size()
}
func (m *FooCondition) XXX_DiscardUnknown() {
	xxx_messageInfo_FooCondition.DiscardUnknown(m)
}

var xxx_messageInfo_FooCondition proto.InternalMessageInfo

func (m *FooConfig) Reset()      { *m = FooConfig{} }
func (*FooConfig) ProtoMessage() {}
func (*FooConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_18c52dbbfb632f7c, []int{2}
}
func (m *FooConfig) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FooConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *FooConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FooConfig.Merge(m, src)
}
func (m *FooConfig) XXX_Size() int {
	return m.Size()
}
func (m *FooConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_FooConfig.DiscardUnknown(m)
}

var xxx_messageInfo_FooConfig proto.InternalMessageInfo

func (m *FooList) Reset()      { *m = FooList{} }
func (*FooList) ProtoMessage() {}
func (*FooList) Descriptor() ([]byte, []int) {
	return fileDescriptor_18c52dbbfb632f7c, []int{3}
}
func (m *FooList) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FooList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *FooList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FooList.Merge(m, src)
}
func (m *FooList) XXX_Size() int {
	return m.Size()
}
func (m *FooList) XXX_DiscardUnknown() {
	xxx_messageInfo_FooList.DiscardUnknown(m)
}

var xxx_messageInfo_FooList proto.InternalMessageInfo

func (m *FooSpec) Reset()      { *m = FooSpec{} }
func (*FooSpec) ProtoMessage() {}
func (*FooSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_18c52dbbfb632f7c, []int{4}
}
func (m *FooSpec) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FooSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *FooSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FooSpec.Merge(m, src)
}
func (m *FooSpec) XXX_Size() int {
	return m.Size()
}
func (m *FooSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_FooSpec.DiscardUnknown(m)
}

var xxx_messageInfo_FooSpec proto.InternalMessageInfo

func (m *FooStatus) Reset()      { *m = FooStatus{} }
func (*FooStatus) ProtoMessage() {}
func (*FooStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_18c52dbbfb632f7c, []int{5}
}
func (m *FooStatus) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *FooStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	b = b[:cap(b)]
	n, err := m.MarshalToSizedBuffer(b)
	if err != nil {
		return nil, err
	}
	return b[:n], nil
}
func (m *FooStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FooStatus.Merge(m, src)
}
func (m *FooStatus) XXX_Size() int {
	return m.Size()
}
func (m *FooStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_FooStatus.DiscardUnknown(m)
}

var xxx_messageInfo_FooStatus proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Foo)(nil), "github.com.phosae.x_kubernetes.api.hello.zeng.dev.v2.Foo")
	proto.RegisterType((*FooCondition)(nil), "github.com.phosae.x_kubernetes.api.hello.zeng.dev.v2.FooCondition")
	proto.RegisterType((*FooConfig)(nil), "github.com.phosae.x_kubernetes.api.hello.zeng.dev.v2.FooConfig")
	proto.RegisterType((*FooList)(nil), "github.com.phosae.x_kubernetes.api.hello.zeng.dev.v2.FooList")
	proto.RegisterType((*FooSpec)(nil), "github.com.phosae.x_kubernetes.api.hello.zeng.dev.v2.FooSpec")
	proto.RegisterType((*FooStatus)(nil), "github.com.phosae.x_kubernetes.api.hello.zeng.dev.v2.FooStatus")
}

func init() {
	proto.RegisterFile("github.com/phosae/x-kubernetes/api/hello.zeng.dev/v2/generated.proto", fileDescriptor_18c52dbbfb632f7c)
}

var fileDescriptor_18c52dbbfb632f7c = []byte{
	// 623 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0x41, 0x6e, 0xd3, 0x40,
	0x14, 0x8d, 0x9b, 0xa4, 0xb4, 0xd3, 0x82, 0xaa, 0x59, 0x95, 0x4a, 0x38, 0x55, 0xd8, 0x74, 0xd3,
	0x19, 0x12, 0x45, 0x08, 0x16, 0x08, 0xc9, 0x20, 0x4b, 0x48, 0xad, 0x40, 0x86, 0x55, 0x85, 0x5a,
	0x26, 0xce, 0x74, 0x3c, 0xa4, 0xf6, 0x58, 0x99, 0x89, 0xd5, 0xb2, 0xe2, 0x08, 0xec, 0xb8, 0x00,
	0x17, 0x40, 0xe2, 0x0e, 0x74, 0xd9, 0x65, 0x57, 0x11, 0x35, 0xb7, 0xe8, 0x0a, 0xcd, 0x8c, 0x1d,
	0xa7, 0xad, 0x10, 0x55, 0xd8, 0xf9, 0x3f, 0xff, 0xff, 0xde, 0xcb, 0xff, 0x2f, 0x06, 0x2f, 0x19,
	0x57, 0xd1, 0xb8, 0x8f, 0x42, 0x11, 0xe3, 0x34, 0x12, 0x92, 0x50, 0x7c, 0xbc, 0x3d, 0x1c, 0xf7,
	0xe9, 0x28, 0xa1, 0x8a, 0x4a, 0x4c, 0x52, 0x8e, 0x23, 0x7a, 0x74, 0x24, 0xd0, 0x27, 0x9a, 0x30,
	0x34, 0xa0, 0x19, 0xce, 0xba, 0x98, 0xd1, 0x84, 0x8e, 0x88, 0xa2, 0x03, 0x94, 0x8e, 0x84, 0x12,
	0xb0, 0x57, 0xb1, 0x20, 0xcb, 0x82, 0x8e, 0x0f, 0x2a, 0x16, 0x44, 0x52, 0x8e, 0xae, 0xb2, 0xa0,
	0xac, 0xbb, 0xb1, 0x3d, 0xa3, 0xcd, 0x04, 0x13, 0xd8, 0x90, 0xf5, 0xc7, 0x87, 0xa6, 0x32, 0x85,
	0x79, 0xb2, 0x22, 0x1b, 0xbd, 0xe1, 0x13, 0x89, 0xb8, 0xd0, 0x96, 0x62, 0x12, 0x46, 0x3c, 0xa1,
	0xa3, 0x13, 0x9c, 0x0e, 0x99, 0x06, 0x24, 0x8e, 0xa9, 0x22, 0x38, 0xeb, 0x5c, 0xb7, 0xb6, 0x81,
	0xff, 0x36, 0x35, 0x1a, 0x27, 0x8a, 0xc7, 0xf4, 0xc6, 0xc0, 0xe3, 0x7f, 0x0d, 0xc8, 0x30, 0xa2,
	0x31, 0xb9, 0x3e, 0xd7, 0xfe, 0xbe, 0x00, 0xea, 0xbe, 0x10, 0xf0, 0x03, 0x58, 0xd2, 0x5e, 0x06,
	0x44, 0x91, 0x75, 0x67, 0xd3, 0xd9, 0x5a, 0xe9, 0x3e, 0x42, 0x96, 0x12, 0xcd, 0x52, 0xa2, 0x74,
	0xc8, 0x34, 0x20, 0x91, 0xee, 0x46, 0x59, 0x07, 0xbd, 0xee, 0x7f, 0xa4, 0xa1, 0xda, 0xa5, 0x8a,
	0x78, 0xf0, 0x74, 0xd2, 0xaa, 0xe5, 0x93, 0x16, 0xa8, 0xb0, 0x60, 0xca, 0x0a, 0x0f, 0x40, 0x43,
	0xa6, 0x34, 0x5c, 0x5f, 0x30, 0xec, 0xcf, 0xd0, 0x3c, 0xcb, 0x47, 0xbe, 0x10, 0x6f, 0x53, 0x1a,
	0x7a, 0xab, 0x85, 0x54, 0x43, 0x57, 0x81, 0x21, 0x86, 0x0c, 0x2c, 0x4a, 0x45, 0xd4, 0x58, 0xae,
	0xd7, 0x8d, 0xc4, 0xf3, 0xf9, 0x25, 0x0c, 0x8d, 0x77, 0xaf, 0x10, 0x59, 0xb4, 0x75, 0x50, 0xd0,
	0xb7, 0xbf, 0x39, 0x60, 0xd5, 0x17, 0xe2, 0x85, 0x48, 0x06, 0x5c, 0x71, 0x91, 0xc0, 0x1e, 0x68,
	0xa8, 0x93, 0x94, 0x9a, 0xc5, 0x2d, 0x7b, 0x9b, 0xa5, 0xb7, 0x77, 0x27, 0x29, 0xbd, 0x9c, 0xb4,
	0xd6, 0x66, 0x7b, 0x35, 0x16, 0x98, 0x6e, 0xb8, 0x3f, 0xf5, 0xbb, 0x60, 0xe6, 0xfc, 0xab, 0x72,
	0x97, 0x93, 0xd6, 0xad, 0xb2, 0x83, 0xa6, 0xdc, 0xd7, 0x6c, 0xee, 0x80, 0x65, 0xab, 0x7c, 0xc8,
	0x19, 0x7c, 0x00, 0xea, 0xb1, 0x64, 0x85, 0xc3, 0x95, 0x42, 0xa9, 0xbe, 0x2b, 0x59, 0xa0, 0x71,
	0xb8, 0x09, 0x1a, 0xb1, 0x64, 0x9d, 0xc2, 0xc9, 0x74, 0xbb, 0xbb, 0x92, 0x75, 0x02, 0xf3, 0xa6,
	0xfd, 0xd3, 0x01, 0x77, 0x7c, 0x21, 0x76, 0xb8, 0x54, 0xf0, 0xfd, 0x8d, 0xb0, 0xa0, 0xdb, 0x85,
	0x45, 0x4f, 0x9b, 0xa8, 0xac, 0x15, 0x0a, 0x4b, 0x25, 0x32, 0x13, 0x94, 0x7d, 0xd0, 0xe4, 0x8a,
	0xc6, 0x7a, 0x2d, 0xf5, 0xad, 0x95, 0xee, 0xd3, 0xb9, 0xcf, 0xe8, 0xdd, 0x2d, 0x54, 0x9a, 0xaf,
	0x34, 0x5f, 0x60, 0x69, 0xdb, 0x5f, 0xed, 0x2f, 0xd1, 0xc9, 0x81, 0x0f, 0x41, 0x93, 0xc7, 0x84,
	0x95, 0xa7, 0xab, 0x06, 0x34, 0x18, 0xd8, 0x77, 0x3a, 0x58, 0xa1, 0xd9, 0x62, 0x91, 0xdd, 0xf9,
	0x83, 0x65, 0x8f, 0x51, 0x05, 0xcb, 0xd6, 0x41, 0x41, 0xdf, 0xfe, 0xe1, 0x98, 0x93, 0xd9, 0x3b,
	0x42, 0x0c, 0x9a, 0x69, 0x44, 0x64, 0xe9, 0xed, 0x7e, 0xe9, 0xed, 0x8d, 0x06, 0x2f, 0x27, 0xad,
	0x25, 0x5f, 0x08, 0xf3, 0x1c, 0xd8, 0x3e, 0x98, 0x01, 0x10, 0x96, 0x59, 0x28, 0xb7, 0xe7, 0xfd,
	0x8f, 0x57, 0x4b, 0x55, 0xfd, 0xaf, 0xa7, 0x90, 0x0c, 0x66, 0x94, 0xbc, 0xbd, 0xd3, 0x0b, 0xb7,
	0x76, 0x76, 0xe1, 0xd6, 0xce, 0x2f, 0xdc, 0xda, 0xe7, 0xdc, 0x75, 0x4e, 0x73, 0xd7, 0x39, 0xcb,
	0x5d, 0xe7, 0x3c, 0x77, 0x9d, 0x5f, 0xb9, 0xeb, 0x7c, 0xf9, 0xed, 0xd6, 0xf6, 0x7a, 0xf3, 0x7c,
	0xb3, 0xff, 0x04, 0x00, 0x00, 0xff, 0xff, 0x39, 0x70, 0x7c, 0x92, 0xea, 0x05, 0x00, 0x00,
}

func (m *Foo) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Foo) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Foo) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Status.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size, err := m.Spec.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	{
		size, err := m.ObjectMeta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *FooCondition) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FooCondition) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FooCondition) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	i -= len(m.Status)
	copy(dAtA[i:], m.Status)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Status)))
	i--
	dAtA[i] = 0x12
	i -= len(m.Type)
	copy(dAtA[i:], m.Type)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Type)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *FooConfig) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FooConfig) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FooConfig) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	i -= len(m.Msg1)
	copy(dAtA[i:], m.Msg1)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Msg1)))
	i--
	dAtA[i] = 0x12
	i -= len(m.Msg)
	copy(dAtA[i:], m.Msg)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Msg)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *FooList) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FooList) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FooList) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Items) > 0 {
		for iNdEx := len(m.Items) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Items[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenerated(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	{
		size, err := m.ListMeta.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *FooSpec) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FooSpec) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FooSpec) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Config.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintGenerated(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	i -= len(m.Image)
	copy(dAtA[i:], m.Image)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Image)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func (m *FooStatus) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *FooStatus) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *FooStatus) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Conditions) > 0 {
		for iNdEx := len(m.Conditions) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Conditions[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintGenerated(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x12
		}
	}
	i -= len(m.Phase)
	copy(dAtA[i:], m.Phase)
	i = encodeVarintGenerated(dAtA, i, uint64(len(m.Phase)))
	i--
	dAtA[i] = 0xa
	return len(dAtA) - i, nil
}

func encodeVarintGenerated(dAtA []byte, offset int, v uint64) int {
	offset -= sovGenerated(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Foo) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ObjectMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = m.Spec.Size()
	n += 1 + l + sovGenerated(uint64(l))
	l = m.Status.Size()
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *FooCondition) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Type)
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.Status)
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *FooConfig) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Msg)
	n += 1 + l + sovGenerated(uint64(l))
	l = len(m.Msg1)
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *FooList) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = m.ListMeta.Size()
	n += 1 + l + sovGenerated(uint64(l))
	if len(m.Items) > 0 {
		for _, e := range m.Items {
			l = e.Size()
			n += 1 + l + sovGenerated(uint64(l))
		}
	}
	return n
}

func (m *FooSpec) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Image)
	n += 1 + l + sovGenerated(uint64(l))
	l = m.Config.Size()
	n += 1 + l + sovGenerated(uint64(l))
	return n
}

func (m *FooStatus) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Phase)
	n += 1 + l + sovGenerated(uint64(l))
	if len(m.Conditions) > 0 {
		for _, e := range m.Conditions {
			l = e.Size()
			n += 1 + l + sovGenerated(uint64(l))
		}
	}
	return n
}

func sovGenerated(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozGenerated(x uint64) (n int) {
	return sovGenerated(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (this *Foo) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&Foo{`,
		`ObjectMeta:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.ObjectMeta), "ObjectMeta", "v1.ObjectMeta", 1), `&`, ``, 1) + `,`,
		`Spec:` + strings.Replace(strings.Replace(this.Spec.String(), "FooSpec", "FooSpec", 1), `&`, ``, 1) + `,`,
		`Status:` + strings.Replace(strings.Replace(this.Status.String(), "FooStatus", "FooStatus", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *FooCondition) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&FooCondition{`,
		`Type:` + fmt.Sprintf("%v", this.Type) + `,`,
		`Status:` + fmt.Sprintf("%v", this.Status) + `,`,
		`}`,
	}, "")
	return s
}
func (this *FooConfig) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&FooConfig{`,
		`Msg:` + fmt.Sprintf("%v", this.Msg) + `,`,
		`Msg1:` + fmt.Sprintf("%v", this.Msg1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *FooList) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForItems := "[]Foo{"
	for _, f := range this.Items {
		repeatedStringForItems += strings.Replace(strings.Replace(f.String(), "Foo", "Foo", 1), `&`, ``, 1) + ","
	}
	repeatedStringForItems += "}"
	s := strings.Join([]string{`&FooList{`,
		`ListMeta:` + strings.Replace(strings.Replace(fmt.Sprintf("%v", this.ListMeta), "ListMeta", "v1.ListMeta", 1), `&`, ``, 1) + `,`,
		`Items:` + repeatedStringForItems + `,`,
		`}`,
	}, "")
	return s
}
func (this *FooSpec) String() string {
	if this == nil {
		return "nil"
	}
	s := strings.Join([]string{`&FooSpec{`,
		`Image:` + fmt.Sprintf("%v", this.Image) + `,`,
		`Config:` + strings.Replace(strings.Replace(this.Config.String(), "FooConfig", "FooConfig", 1), `&`, ``, 1) + `,`,
		`}`,
	}, "")
	return s
}
func (this *FooStatus) String() string {
	if this == nil {
		return "nil"
	}
	repeatedStringForConditions := "[]FooCondition{"
	for _, f := range this.Conditions {
		repeatedStringForConditions += strings.Replace(strings.Replace(f.String(), "FooCondition", "FooCondition", 1), `&`, ``, 1) + ","
	}
	repeatedStringForConditions += "}"
	s := strings.Join([]string{`&FooStatus{`,
		`Phase:` + fmt.Sprintf("%v", this.Phase) + `,`,
		`Conditions:` + repeatedStringForConditions + `,`,
		`}`,
	}, "")
	return s
}
func valueToStringGenerated(v interface{}) string {
	rv := reflect.ValueOf(v)
	if rv.IsNil() {
		return "nil"
	}
	pv := reflect.Indirect(rv).Interface()
	return fmt.Sprintf("*%v", pv)
}
func (m *Foo) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Foo: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Foo: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ObjectMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ObjectMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Spec", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Spec.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Status.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FooCondition) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FooCondition: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FooCondition: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Type", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Type = FooConditionType(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Status", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Status = k8s_io_apimachinery_pkg_apis_meta_v1.ConditionStatus(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FooConfig) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FooConfig: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FooConfig: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg1", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg1 = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FooList) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FooList: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FooList: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ListMeta", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ListMeta.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Items", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Items = append(m.Items, Foo{})
			if err := m.Items[len(m.Items)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FooSpec) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FooSpec: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FooSpec: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Image", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Image = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Config", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.Config.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *FooStatus) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: FooStatus: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: FooStatus: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Phase", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Phase = FooPhase(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Conditions", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthGenerated
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthGenerated
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Conditions = append(m.Conditions, FooCondition{})
			if err := m.Conditions[len(m.Conditions)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipGenerated(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthGenerated
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipGenerated(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowGenerated
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowGenerated
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthGenerated
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupGenerated
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthGenerated
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthGenerated        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowGenerated          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupGenerated = fmt.Errorf("proto: unexpected end of group")
)
