// Code generated by protoc-gen-go. DO NOT EDIT.
// source: order.proto

package pb

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

type OrderCreateCommand struct {
	OrderId              string                          `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	CustomerId           string                          `protobuf:"bytes,2,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Status               string                          `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	CreatedOn            int64                           `protobuf:"varint,4,opt,name=created_on,json=createdOn,proto3" json:"created_on,omitempty"`
	RestaurantId         string                          `protobuf:"bytes,5,opt,name=restaurant_id,json=restaurantId,proto3" json:"restaurant_id,omitempty"`
	Amount               float32                         `protobuf:"fixed32,6,opt,name=amount,proto3" json:"amount,omitempty"`
	OrderItems           []*OrderCreateCommand_OrderItem `protobuf:"bytes,7,rep,name=order_items,json=orderItems,proto3" json:"order_items,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *OrderCreateCommand) Reset()         { *m = OrderCreateCommand{} }
func (m *OrderCreateCommand) String() string { return proto.CompactTextString(m) }
func (*OrderCreateCommand) ProtoMessage()    {}
func (*OrderCreateCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_719bea36e3298c99, []int{0}
}
func (m *OrderCreateCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderCreateCommand.Unmarshal(m, b)
}
func (m *OrderCreateCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderCreateCommand.Marshal(b, m, deterministic)
}
func (dst *OrderCreateCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderCreateCommand.Merge(dst, src)
}
func (m *OrderCreateCommand) XXX_Size() int {
	return xxx_messageInfo_OrderCreateCommand.Size(m)
}
func (m *OrderCreateCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderCreateCommand.DiscardUnknown(m)
}

var xxx_messageInfo_OrderCreateCommand proto.InternalMessageInfo

func (m *OrderCreateCommand) GetOrderId() string {
	if m != nil {
		return m.OrderId
	}
	return ""
}

func (m *OrderCreateCommand) GetCustomerId() string {
	if m != nil {
		return m.CustomerId
	}
	return ""
}

func (m *OrderCreateCommand) GetStatus() string {
	if m != nil {
		return m.Status
	}
	return ""
}

func (m *OrderCreateCommand) GetCreatedOn() int64 {
	if m != nil {
		return m.CreatedOn
	}
	return 0
}

func (m *OrderCreateCommand) GetRestaurantId() string {
	if m != nil {
		return m.RestaurantId
	}
	return ""
}

func (m *OrderCreateCommand) GetAmount() float32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func (m *OrderCreateCommand) GetOrderItems() []*OrderCreateCommand_OrderItem {
	if m != nil {
		return m.OrderItems
	}
	return nil
}

type OrderCreateCommand_OrderItem struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	UnitPrice            float32  `protobuf:"fixed32,3,opt,name=unit_price,json=unitPrice,proto3" json:"unit_price,omitempty"`
	Quantity             int32    `protobuf:"varint,4,opt,name=quantity,proto3" json:"quantity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderCreateCommand_OrderItem) Reset()         { *m = OrderCreateCommand_OrderItem{} }
func (m *OrderCreateCommand_OrderItem) String() string { return proto.CompactTextString(m) }
func (*OrderCreateCommand_OrderItem) ProtoMessage()    {}
func (*OrderCreateCommand_OrderItem) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_719bea36e3298c99, []int{0, 0}
}
func (m *OrderCreateCommand_OrderItem) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderCreateCommand_OrderItem.Unmarshal(m, b)
}
func (m *OrderCreateCommand_OrderItem) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderCreateCommand_OrderItem.Marshal(b, m, deterministic)
}
func (dst *OrderCreateCommand_OrderItem) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderCreateCommand_OrderItem.Merge(dst, src)
}
func (m *OrderCreateCommand_OrderItem) XXX_Size() int {
	return xxx_messageInfo_OrderCreateCommand_OrderItem.Size(m)
}
func (m *OrderCreateCommand_OrderItem) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderCreateCommand_OrderItem.DiscardUnknown(m)
}

var xxx_messageInfo_OrderCreateCommand_OrderItem proto.InternalMessageInfo

func (m *OrderCreateCommand_OrderItem) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *OrderCreateCommand_OrderItem) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *OrderCreateCommand_OrderItem) GetUnitPrice() float32 {
	if m != nil {
		return m.UnitPrice
	}
	return 0
}

func (m *OrderCreateCommand_OrderItem) GetQuantity() int32 {
	if m != nil {
		return m.Quantity
	}
	return 0
}

type OrderPaymentDebitedCommand struct {
	OrderId              string   `protobuf:"bytes,1,opt,name=order_id,json=orderId,proto3" json:"order_id,omitempty"`
	CustomerId           string   `protobuf:"bytes,2,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	Amount               float32  `protobuf:"fixed32,3,opt,name=amount,proto3" json:"amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderPaymentDebitedCommand) Reset()         { *m = OrderPaymentDebitedCommand{} }
func (m *OrderPaymentDebitedCommand) String() string { return proto.CompactTextString(m) }
func (*OrderPaymentDebitedCommand) ProtoMessage()    {}
func (*OrderPaymentDebitedCommand) Descriptor() ([]byte, []int) {
	return fileDescriptor_order_719bea36e3298c99, []int{1}
}
func (m *OrderPaymentDebitedCommand) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderPaymentDebitedCommand.Unmarshal(m, b)
}
func (m *OrderPaymentDebitedCommand) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderPaymentDebitedCommand.Marshal(b, m, deterministic)
}
func (dst *OrderPaymentDebitedCommand) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderPaymentDebitedCommand.Merge(dst, src)
}
func (m *OrderPaymentDebitedCommand) XXX_Size() int {
	return xxx_messageInfo_OrderPaymentDebitedCommand.Size(m)
}
func (m *OrderPaymentDebitedCommand) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderPaymentDebitedCommand.DiscardUnknown(m)
}

var xxx_messageInfo_OrderPaymentDebitedCommand proto.InternalMessageInfo

func (m *OrderPaymentDebitedCommand) GetOrderId() string {
	if m != nil {
		return m.OrderId
	}
	return ""
}

func (m *OrderPaymentDebitedCommand) GetCustomerId() string {
	if m != nil {
		return m.CustomerId
	}
	return ""
}

func (m *OrderPaymentDebitedCommand) GetAmount() float32 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func init() {
	proto.RegisterType((*OrderCreateCommand)(nil), "pb.OrderCreateCommand")
	proto.RegisterType((*OrderCreateCommand_OrderItem)(nil), "pb.OrderCreateCommand.OrderItem")
	proto.RegisterType((*OrderPaymentDebitedCommand)(nil), "pb.OrderPaymentDebitedCommand")
}

func init() { proto.RegisterFile("order.proto", fileDescriptor_order_719bea36e3298c99) }

var fileDescriptor_order_719bea36e3298c99 = []byte{
	// 290 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x91, 0xbf, 0x4e, 0xc3, 0x30,
	0x10, 0xc6, 0x95, 0xa4, 0xff, 0x72, 0x85, 0xc5, 0x03, 0x0a, 0x95, 0x10, 0x51, 0x59, 0x32, 0x75,
	0x80, 0x27, 0x40, 0x65, 0xe9, 0xd4, 0x2a, 0x2f, 0x10, 0x39, 0xf1, 0x0d, 0x19, 0x6c, 0x07, 0xfb,
	0x3c, 0xf4, 0xc9, 0x59, 0x91, 0xaf, 0x6e, 0x41, 0x62, 0x65, 0xbb, 0xef, 0xf7, 0xd9, 0x77, 0xfa,
	0xee, 0x60, 0x6d, 0x9d, 0x42, 0xb7, 0x9b, 0x9c, 0x25, 0x2b, 0xf2, 0xa9, 0xdf, 0x7e, 0xe5, 0x20,
	0x8e, 0x91, 0xed, 0x1d, 0x4a, 0xc2, 0xbd, 0xd5, 0x5a, 0x1a, 0x25, 0x1e, 0x61, 0xc5, 0x2f, 0xbb,
	0x51, 0x55, 0x59, 0x9d, 0x35, 0x65, 0xbb, 0x64, 0x7d, 0x50, 0xe2, 0x19, 0xd6, 0x43, 0xf0, 0x64,
	0xf5, 0xc5, 0xcd, 0xd9, 0x85, 0x2b, 0x3a, 0x28, 0xf1, 0x00, 0x0b, 0x4f, 0x92, 0x82, 0xaf, 0x0a,
	0xf6, 0x92, 0x12, 0x4f, 0x00, 0x03, 0x0f, 0x51, 0x9d, 0x35, 0xd5, 0xac, 0xce, 0x9a, 0xa2, 0x2d,
	0x13, 0x39, 0x1a, 0xf1, 0x02, 0xf7, 0x0e, 0x3d, 0xc9, 0xe0, 0xa4, 0xa1, 0xd8, 0x79, 0xce, 0xbf,
	0xef, 0x7e, 0xe0, 0xa5, 0xb7, 0xd4, 0x36, 0x18, 0xaa, 0x16, 0x75, 0xd6, 0xe4, 0x6d, 0x52, 0xe2,
	0x3d, 0x25, 0xeb, 0x46, 0x42, 0xed, 0xab, 0x65, 0x5d, 0x34, 0xeb, 0xd7, 0x7a, 0x37, 0xf5, 0xbb,
	0xbf, 0xe1, 0x2e, 0xe8, 0x40, 0xa8, 0x5b, 0xb0, 0xd7, 0xd2, 0x6f, 0x0c, 0x94, 0x37, 0x43, 0x08,
	0x98, 0x0d, 0x56, 0x61, 0xca, 0xce, 0x75, 0x64, 0x46, 0x6a, 0x4c, 0x89, 0xb9, 0x8e, 0x99, 0x82,
	0x19, 0xa9, 0x9b, 0xdc, 0x38, 0x20, 0xe7, 0xcd, 0xdb, 0x32, 0x92, 0x53, 0x04, 0x62, 0x03, 0xab,
	0xcf, 0x20, 0x0d, 0x8d, 0x74, 0xe6, 0xc0, 0xf3, 0xf6, 0xa6, 0xb7, 0x13, 0x6c, 0x78, 0xde, 0x49,
	0x9e, 0x35, 0x1a, 0xfa, 0xc0, 0x7e, 0x24, 0x54, 0xff, 0x74, 0x80, 0xb4, 0xa4, 0xe2, 0xf7, 0x92,
	0xfa, 0x05, 0x9f, 0xfd, 0xed, 0x3b, 0x00, 0x00, 0xff, 0xff, 0x35, 0x82, 0x87, 0x5b, 0x05, 0x02,
	0x00, 0x00,
}