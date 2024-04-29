// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: proto/replication.proto

package replication

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ReplicationServiceClient is the client API for ReplicationService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ReplicationServiceClient interface {
	// replicate a single op, ideally it should be a single
	// function that takes in the op + args the client recieved but for now
	// just a simple replicate write function that always replicate
	ReplicateWrite(ctx context.Context, in *ReplicateWriteRequest, opts ...grpc.CallOption) (*ReplicateWriteResponse, error)
	GetKey(ctx context.Context, in *GetKeyRequest, opts ...grpc.CallOption) (*GetKeyResponse, error)
	GetMerkleTree(ctx context.Context, in *MerkleTreeRequest, opts ...grpc.CallOption) (*MerkleTreeResponse, error)
}

type replicationServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewReplicationServiceClient(cc grpc.ClientConnInterface) ReplicationServiceClient {
	return &replicationServiceClient{cc}
}

func (c *replicationServiceClient) ReplicateWrite(ctx context.Context, in *ReplicateWriteRequest, opts ...grpc.CallOption) (*ReplicateWriteResponse, error) {
	out := new(ReplicateWriteResponse)
	err := c.cc.Invoke(ctx, "/kvbridge.ReplicationService/ReplicateWrite", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicationServiceClient) GetKey(ctx context.Context, in *GetKeyRequest, opts ...grpc.CallOption) (*GetKeyResponse, error) {
	out := new(GetKeyResponse)
	err := c.cc.Invoke(ctx, "/kvbridge.ReplicationService/GetKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *replicationServiceClient) GetMerkleTree(ctx context.Context, in *MerkleTreeRequest, opts ...grpc.CallOption) (*MerkleTreeResponse, error) {
	out := new(MerkleTreeResponse)
	err := c.cc.Invoke(ctx, "/kvbridge.ReplicationService/GetMerkleTree", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ReplicationServiceServer is the server API for ReplicationService service.
// All implementations must embed UnimplementedReplicationServiceServer
// for forward compatibility
type ReplicationServiceServer interface {
	// replicate a single op, ideally it should be a single
	// function that takes in the op + args the client recieved but for now
	// just a simple replicate write function that always replicate
	ReplicateWrite(context.Context, *ReplicateWriteRequest) (*ReplicateWriteResponse, error)
	GetKey(context.Context, *GetKeyRequest) (*GetKeyResponse, error)
	GetMerkleTree(context.Context, *MerkleTreeRequest) (*MerkleTreeResponse, error)
	mustEmbedUnimplementedReplicationServiceServer()
}

// UnimplementedReplicationServiceServer must be embedded to have forward compatible implementations.
type UnimplementedReplicationServiceServer struct {
}

func (UnimplementedReplicationServiceServer) ReplicateWrite(context.Context, *ReplicateWriteRequest) (*ReplicateWriteResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReplicateWrite not implemented")
}
func (UnimplementedReplicationServiceServer) GetKey(context.Context, *GetKeyRequest) (*GetKeyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetKey not implemented")
}
func (UnimplementedReplicationServiceServer) GetMerkleTree(context.Context, *MerkleTreeRequest) (*MerkleTreeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetMerkleTree not implemented")
}
func (UnimplementedReplicationServiceServer) mustEmbedUnimplementedReplicationServiceServer() {}

// UnsafeReplicationServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ReplicationServiceServer will
// result in compilation errors.
type UnsafeReplicationServiceServer interface {
	mustEmbedUnimplementedReplicationServiceServer()
}

func RegisterReplicationServiceServer(s grpc.ServiceRegistrar, srv ReplicationServiceServer) {
	s.RegisterService(&ReplicationService_ServiceDesc, srv)
}

func _ReplicationService_ReplicateWrite_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReplicateWriteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServiceServer).ReplicateWrite(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvbridge.ReplicationService/ReplicateWrite",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServiceServer).ReplicateWrite(ctx, req.(*ReplicateWriteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReplicationService_GetKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetKeyRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServiceServer).GetKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvbridge.ReplicationService/GetKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServiceServer).GetKey(ctx, req.(*GetKeyRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReplicationService_GetMerkleTree_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MerkleTreeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReplicationServiceServer).GetMerkleTree(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/kvbridge.ReplicationService/GetMerkleTree",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReplicationServiceServer).GetMerkleTree(ctx, req.(*MerkleTreeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ReplicationService_ServiceDesc is the grpc.ServiceDesc for ReplicationService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ReplicationService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "kvbridge.ReplicationService",
	HandlerType: (*ReplicationServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ReplicateWrite",
			Handler:    _ReplicationService_ReplicateWrite_Handler,
		},
		{
			MethodName: "GetKey",
			Handler:    _ReplicationService_GetKey_Handler,
		},
		{
			MethodName: "GetMerkleTree",
			Handler:    _ReplicationService_GetMerkleTree_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/replication.proto",
}
