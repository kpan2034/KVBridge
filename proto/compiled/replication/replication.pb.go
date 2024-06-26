// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.20.3
// source: proto/replication.proto

package replication

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

type ReplicateWriteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// bool deferred = 1; // can the write be deferred? used when ack is not
	// needed AFTER write is made durable
	// In this case the client can get away with writing this value only to the
	// WAL and only syncing periodically
	Id    uint32 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Key   []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ReplicateWriteRequest) Reset() {
	*x = ReplicateWriteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplicateWriteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicateWriteRequest) ProtoMessage() {}

func (x *ReplicateWriteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicateWriteRequest.ProtoReflect.Descriptor instead.
func (*ReplicateWriteRequest) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{0}
}

func (x *ReplicateWriteRequest) GetId() uint32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *ReplicateWriteRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *ReplicateWriteRequest) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type ReplicateWriteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Resp:
	//
	//	*ReplicateWriteResponse_Ok
	//	*ReplicateWriteResponse_Value
	Resp isReplicateWriteResponse_Resp `protobuf_oneof:"resp"`
}

func (x *ReplicateWriteResponse) Reset() {
	*x = ReplicateWriteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ReplicateWriteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReplicateWriteResponse) ProtoMessage() {}

func (x *ReplicateWriteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReplicateWriteResponse.ProtoReflect.Descriptor instead.
func (*ReplicateWriteResponse) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{1}
}

func (m *ReplicateWriteResponse) GetResp() isReplicateWriteResponse_Resp {
	if m != nil {
		return m.Resp
	}
	return nil
}

func (x *ReplicateWriteResponse) GetOk() bool {
	if x, ok := x.GetResp().(*ReplicateWriteResponse_Ok); ok {
		return x.Ok
	}
	return false
}

func (x *ReplicateWriteResponse) GetValue() []byte {
	if x, ok := x.GetResp().(*ReplicateWriteResponse_Value); ok {
		return x.Value
	}
	return nil
}

type isReplicateWriteResponse_Resp interface {
	isReplicateWriteResponse_Resp()
}

type ReplicateWriteResponse_Ok struct {
	Ok bool `protobuf:"varint,1,opt,name=ok,proto3,oneof"`
}

type ReplicateWriteResponse_Value struct {
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3,oneof"`
}

func (*ReplicateWriteResponse_Ok) isReplicateWriteResponse_Resp() {}

func (*ReplicateWriteResponse_Value) isReplicateWriteResponse_Resp() {}

type GetKeyRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key []byte `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetKeyRequest) Reset() {
	*x = GetKeyRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKeyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeyRequest) ProtoMessage() {}

func (x *GetKeyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeyRequest.ProtoReflect.Descriptor instead.
func (*GetKeyRequest) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{2}
}

func (x *GetKeyRequest) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

type GetKeyResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok    bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Value []byte `protobuf:"bytes,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GetKeyResponse) Reset() {
	*x = GetKeyResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKeyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeyResponse) ProtoMessage() {}

func (x *GetKeyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeyResponse.ProtoReflect.Descriptor instead.
func (*GetKeyResponse) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{3}
}

func (x *GetKeyResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *GetKeyResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// Merkle tree
type MerkleTreeRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeyRangeLowerBound uint32 `protobuf:"varint,1,opt,name=keyRangeLowerBound,proto3" json:"keyRangeLowerBound,omitempty"`
	KeyRangeUpperBound uint32 `protobuf:"varint,2,opt,name=keyRangeUpperBound,proto3" json:"keyRangeUpperBound,omitempty"`
}

func (x *MerkleTreeRequest) Reset() {
	*x = MerkleTreeRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MerkleTreeRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MerkleTreeRequest) ProtoMessage() {}

func (x *MerkleTreeRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MerkleTreeRequest.ProtoReflect.Descriptor instead.
func (*MerkleTreeRequest) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{4}
}

func (x *MerkleTreeRequest) GetKeyRangeLowerBound() uint32 {
	if x != nil {
		return x.KeyRangeLowerBound
	}
	return 0
}

func (x *MerkleTreeRequest) GetKeyRangeUpperBound() uint32 {
	if x != nil {
		return x.KeyRangeUpperBound
	}
	return 0
}

type MerkleTreeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []uint32 `protobuf:"varint,1,rep,packed,name=data,proto3" json:"data,omitempty"`
}

func (x *MerkleTreeResponse) Reset() {
	*x = MerkleTreeResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MerkleTreeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MerkleTreeResponse) ProtoMessage() {}

func (x *MerkleTreeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MerkleTreeResponse.ProtoReflect.Descriptor instead.
func (*MerkleTreeResponse) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{5}
}

func (x *MerkleTreeResponse) GetData() []uint32 {
	if x != nil {
		return x.Data
	}
	return nil
}

type GetKeysInRangesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeyRangeList []*GetKeysInRangesRequest_KeyRange `protobuf:"bytes,1,rep,name=keyRangeList,proto3" json:"keyRangeList,omitempty"`
}

func (x *GetKeysInRangesRequest) Reset() {
	*x = GetKeysInRangesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKeysInRangesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeysInRangesRequest) ProtoMessage() {}

func (x *GetKeysInRangesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeysInRangesRequest.ProtoReflect.Descriptor instead.
func (*GetKeysInRangesRequest) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{6}
}

func (x *GetKeysInRangesRequest) GetKeyRangeList() []*GetKeysInRangesRequest_KeyRange {
	if x != nil {
		return x.KeyRangeList
	}
	return nil
}

type GetKeysInRangesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ok    bool   `protobuf:"varint,1,opt,name=ok,proto3" json:"ok,omitempty"`
	Key   []byte `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
	Value []byte `protobuf:"bytes,3,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *GetKeysInRangesResponse) Reset() {
	*x = GetKeysInRangesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKeysInRangesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeysInRangesResponse) ProtoMessage() {}

func (x *GetKeysInRangesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeysInRangesResponse.ProtoReflect.Descriptor instead.
func (*GetKeysInRangesResponse) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{7}
}

func (x *GetKeysInRangesResponse) GetOk() bool {
	if x != nil {
		return x.Ok
	}
	return false
}

func (x *GetKeysInRangesResponse) GetKey() []byte {
	if x != nil {
		return x.Key
	}
	return nil
}

func (x *GetKeysInRangesResponse) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

type GetKeysInRangesRequest_KeyRange struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	KeyRangeLowerBound uint32 `protobuf:"varint,1,opt,name=keyRangeLowerBound,proto3" json:"keyRangeLowerBound,omitempty"`
	KeyRangeUpperBound uint32 `protobuf:"varint,2,opt,name=keyRangeUpperBound,proto3" json:"keyRangeUpperBound,omitempty"`
}

func (x *GetKeysInRangesRequest_KeyRange) Reset() {
	*x = GetKeysInRangesRequest_KeyRange{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_replication_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetKeysInRangesRequest_KeyRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetKeysInRangesRequest_KeyRange) ProtoMessage() {}

func (x *GetKeysInRangesRequest_KeyRange) ProtoReflect() protoreflect.Message {
	mi := &file_proto_replication_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetKeysInRangesRequest_KeyRange.ProtoReflect.Descriptor instead.
func (*GetKeysInRangesRequest_KeyRange) Descriptor() ([]byte, []int) {
	return file_proto_replication_proto_rawDescGZIP(), []int{6, 0}
}

func (x *GetKeysInRangesRequest_KeyRange) GetKeyRangeLowerBound() uint32 {
	if x != nil {
		return x.KeyRangeLowerBound
	}
	return 0
}

func (x *GetKeysInRangesRequest_KeyRange) GetKeyRangeUpperBound() uint32 {
	if x != nil {
		return x.KeyRangeUpperBound
	}
	return 0
}

var File_proto_replication_proto protoreflect.FileDescriptor

var file_proto_replication_proto_rawDesc = []byte{
	0x0a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x08, 0x6b, 0x76, 0x62, 0x72, 0x69,
	0x64, 0x67, 0x65, 0x22, 0x4f, 0x0a, 0x15, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x02, 0x69, 0x64, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x22, 0x4a, 0x0a, 0x16, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x10,
	0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x02, 0x6f, 0x6b,
	0x12, 0x16, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x48,
	0x00, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x72, 0x65, 0x73, 0x70,
	0x22, 0x21, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03,
	0x6b, 0x65, 0x79, 0x22, 0x36, 0x0a, 0x0e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x73, 0x0a, 0x11, 0x4d,
	0x65, 0x72, 0x6b, 0x6c, 0x65, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x2e, 0x0a, 0x12, 0x6b, 0x65, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x77, 0x65,
	0x72, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x12, 0x6b, 0x65,
	0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x77, 0x65, 0x72, 0x42, 0x6f, 0x75, 0x6e, 0x64,
	0x12, 0x2e, 0x0a, 0x12, 0x6b, 0x65, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x55, 0x70, 0x70, 0x65,
	0x72, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x12, 0x6b, 0x65,
	0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x55, 0x70, 0x70, 0x65, 0x72, 0x42, 0x6f, 0x75, 0x6e, 0x64,
	0x22, 0x28, 0x0a, 0x12, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x54, 0x72, 0x65, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0d, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xd3, 0x01, 0x0a, 0x16, 0x47,
	0x65, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x49, 0x6e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x4d, 0x0a, 0x0c, 0x6b, 0x65, 0x79, 0x52, 0x61, 0x6e, 0x67,
	0x65, 0x4c, 0x69, 0x73, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x6b, 0x76,
	0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x49, 0x6e,
	0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4b, 0x65,
	0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x52, 0x0c, 0x6b, 0x65, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65,
	0x4c, 0x69, 0x73, 0x74, 0x1a, 0x6a, 0x0a, 0x08, 0x4b, 0x65, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65,
	0x12, 0x2e, 0x0a, 0x12, 0x6b, 0x65, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x77, 0x65,
	0x72, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x12, 0x6b, 0x65,
	0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x4c, 0x6f, 0x77, 0x65, 0x72, 0x42, 0x6f, 0x75, 0x6e, 0x64,
	0x12, 0x2e, 0x0a, 0x12, 0x6b, 0x65, 0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x55, 0x70, 0x70, 0x65,
	0x72, 0x42, 0x6f, 0x75, 0x6e, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x12, 0x6b, 0x65,
	0x79, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x55, 0x70, 0x70, 0x65, 0x72, 0x42, 0x6f, 0x75, 0x6e, 0x64,
	0x22, 0x51, 0x0a, 0x17, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x49, 0x6e, 0x52, 0x61, 0x6e,
	0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x6f,
	0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6b, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x32, 0xcc, 0x02, 0x0a, 0x12, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x53, 0x0a, 0x0e, 0x52, 0x65,
	0x70, 0x6c, 0x69, 0x63, 0x61, 0x74, 0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x12, 0x1f, 0x2e, 0x6b,
	0x76, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e,
	0x6b, 0x76, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x3b, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x12, 0x17, 0x2e, 0x6b, 0x76, 0x62, 0x72,
	0x69, 0x64, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x1a, 0x18, 0x2e, 0x6b, 0x76, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x47, 0x65,
	0x74, 0x4b, 0x65, 0x79, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x4a, 0x0a, 0x0d,
	0x47, 0x65, 0x74, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x54, 0x72, 0x65, 0x65, 0x12, 0x1b, 0x2e,
	0x6b, 0x76, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x54,
	0x72, 0x65, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1c, 0x2e, 0x6b, 0x76, 0x62,
	0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x72, 0x6b, 0x6c, 0x65, 0x54, 0x72, 0x65, 0x65,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x58, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4b,
	0x65, 0x79, 0x73, 0x49, 0x6e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x12, 0x20, 0x2e, 0x6b, 0x76,
	0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73, 0x49, 0x6e,
	0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e,
	0x6b, 0x76, 0x62, 0x72, 0x69, 0x64, 0x67, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4b, 0x65, 0x79, 0x73,
	0x49, 0x6e, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x30, 0x01, 0x42, 0x0f, 0x5a, 0x0d, 0x2e, 0x2f, 0x72, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_replication_proto_rawDescOnce sync.Once
	file_proto_replication_proto_rawDescData = file_proto_replication_proto_rawDesc
)

func file_proto_replication_proto_rawDescGZIP() []byte {
	file_proto_replication_proto_rawDescOnce.Do(func() {
		file_proto_replication_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_replication_proto_rawDescData)
	})
	return file_proto_replication_proto_rawDescData
}

var file_proto_replication_proto_msgTypes = make([]protoimpl.MessageInfo, 9)
var file_proto_replication_proto_goTypes = []interface{}{
	(*ReplicateWriteRequest)(nil),           // 0: kvbridge.ReplicateWriteRequest
	(*ReplicateWriteResponse)(nil),          // 1: kvbridge.ReplicateWriteResponse
	(*GetKeyRequest)(nil),                   // 2: kvbridge.GetKeyRequest
	(*GetKeyResponse)(nil),                  // 3: kvbridge.GetKeyResponse
	(*MerkleTreeRequest)(nil),               // 4: kvbridge.MerkleTreeRequest
	(*MerkleTreeResponse)(nil),              // 5: kvbridge.MerkleTreeResponse
	(*GetKeysInRangesRequest)(nil),          // 6: kvbridge.GetKeysInRangesRequest
	(*GetKeysInRangesResponse)(nil),         // 7: kvbridge.GetKeysInRangesResponse
	(*GetKeysInRangesRequest_KeyRange)(nil), // 8: kvbridge.GetKeysInRangesRequest.KeyRange
}
var file_proto_replication_proto_depIdxs = []int32{
	8, // 0: kvbridge.GetKeysInRangesRequest.keyRangeList:type_name -> kvbridge.GetKeysInRangesRequest.KeyRange
	0, // 1: kvbridge.ReplicationService.ReplicateWrite:input_type -> kvbridge.ReplicateWriteRequest
	2, // 2: kvbridge.ReplicationService.GetKey:input_type -> kvbridge.GetKeyRequest
	4, // 3: kvbridge.ReplicationService.GetMerkleTree:input_type -> kvbridge.MerkleTreeRequest
	6, // 4: kvbridge.ReplicationService.GetKeysInRanges:input_type -> kvbridge.GetKeysInRangesRequest
	1, // 5: kvbridge.ReplicationService.ReplicateWrite:output_type -> kvbridge.ReplicateWriteResponse
	3, // 6: kvbridge.ReplicationService.GetKey:output_type -> kvbridge.GetKeyResponse
	5, // 7: kvbridge.ReplicationService.GetMerkleTree:output_type -> kvbridge.MerkleTreeResponse
	7, // 8: kvbridge.ReplicationService.GetKeysInRanges:output_type -> kvbridge.GetKeysInRangesResponse
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_replication_proto_init() }
func file_proto_replication_proto_init() {
	if File_proto_replication_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_replication_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplicateWriteRequest); i {
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
		file_proto_replication_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ReplicateWriteResponse); i {
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
		file_proto_replication_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKeyRequest); i {
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
		file_proto_replication_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKeyResponse); i {
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
		file_proto_replication_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MerkleTreeRequest); i {
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
		file_proto_replication_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MerkleTreeResponse); i {
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
		file_proto_replication_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKeysInRangesRequest); i {
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
		file_proto_replication_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKeysInRangesResponse); i {
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
		file_proto_replication_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetKeysInRangesRequest_KeyRange); i {
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
	file_proto_replication_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*ReplicateWriteResponse_Ok)(nil),
		(*ReplicateWriteResponse_Value)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_replication_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   9,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_replication_proto_goTypes,
		DependencyIndexes: file_proto_replication_proto_depIdxs,
		MessageInfos:      file_proto_replication_proto_msgTypes,
	}.Build()
	File_proto_replication_proto = out.File
	file_proto_replication_proto_rawDesc = nil
	file_proto_replication_proto_goTypes = nil
	file_proto_replication_proto_depIdxs = nil
}
