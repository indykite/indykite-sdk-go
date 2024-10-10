// Copyright (c) 2023 IndyKite
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Ingest Service Description.

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: indykite/ingest/v1beta3/ingest_api.proto

package ingestv1beta3

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.62.0 or later.
const _ = grpc.SupportPackageIsVersion8

const (
	IngestAPI_StreamRecords_FullMethodName                     = "/indykite.ingest.v1beta3.IngestAPI/StreamRecords"
	IngestAPI_IngestRecord_FullMethodName                      = "/indykite.ingest.v1beta3.IngestAPI/IngestRecord"
	IngestAPI_BatchUpsertNodes_FullMethodName                  = "/indykite.ingest.v1beta3.IngestAPI/BatchUpsertNodes"
	IngestAPI_BatchUpsertRelationships_FullMethodName          = "/indykite.ingest.v1beta3.IngestAPI/BatchUpsertRelationships"
	IngestAPI_BatchDeleteNodes_FullMethodName                  = "/indykite.ingest.v1beta3.IngestAPI/BatchDeleteNodes"
	IngestAPI_BatchDeleteRelationships_FullMethodName          = "/indykite.ingest.v1beta3.IngestAPI/BatchDeleteRelationships"
	IngestAPI_BatchDeleteNodeProperties_FullMethodName         = "/indykite.ingest.v1beta3.IngestAPI/BatchDeleteNodeProperties"
	IngestAPI_BatchDeleteRelationshipProperties_FullMethodName = "/indykite.ingest.v1beta3.IngestAPI/BatchDeleteRelationshipProperties"
	IngestAPI_BatchDeleteNodeTags_FullMethodName               = "/indykite.ingest.v1beta3.IngestAPI/BatchDeleteNodeTags"
)

// IngestAPIClient is the client API for IngestAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// IngestAPI represents the service interface for data ingestion.
type IngestAPIClient interface {
	StreamRecords(ctx context.Context, opts ...grpc.CallOption) (IngestAPI_StreamRecordsClient, error)
	// Deprecated: Do not use.
	IngestRecord(ctx context.Context, in *IngestRecordRequest, opts ...grpc.CallOption) (*IngestRecordResponse, error)
	BatchUpsertNodes(ctx context.Context, in *BatchUpsertNodesRequest, opts ...grpc.CallOption) (*BatchUpsertNodesResponse, error)
	BatchUpsertRelationships(ctx context.Context, in *BatchUpsertRelationshipsRequest, opts ...grpc.CallOption) (*BatchUpsertRelationshipsResponse, error)
	BatchDeleteNodes(ctx context.Context, in *BatchDeleteNodesRequest, opts ...grpc.CallOption) (*BatchDeleteNodesResponse, error)
	BatchDeleteRelationships(ctx context.Context, in *BatchDeleteRelationshipsRequest, opts ...grpc.CallOption) (*BatchDeleteRelationshipsResponse, error)
	BatchDeleteNodeProperties(ctx context.Context, in *BatchDeleteNodePropertiesRequest, opts ...grpc.CallOption) (*BatchDeleteNodePropertiesResponse, error)
	BatchDeleteRelationshipProperties(ctx context.Context, in *BatchDeleteRelationshipPropertiesRequest, opts ...grpc.CallOption) (*BatchDeleteRelationshipPropertiesResponse, error)
	BatchDeleteNodeTags(ctx context.Context, in *BatchDeleteNodeTagsRequest, opts ...grpc.CallOption) (*BatchDeleteNodeTagsResponse, error)
}

type ingestAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewIngestAPIClient(cc grpc.ClientConnInterface) IngestAPIClient {
	return &ingestAPIClient{cc}
}

func (c *ingestAPIClient) StreamRecords(ctx context.Context, opts ...grpc.CallOption) (IngestAPI_StreamRecordsClient, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &IngestAPI_ServiceDesc.Streams[0], IngestAPI_StreamRecords_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &ingestAPIStreamRecordsClient{ClientStream: stream}
	return x, nil
}

type IngestAPI_StreamRecordsClient interface {
	Send(*StreamRecordsRequest) error
	Recv() (*StreamRecordsResponse, error)
	grpc.ClientStream
}

type ingestAPIStreamRecordsClient struct {
	grpc.ClientStream
}

func (x *ingestAPIStreamRecordsClient) Send(m *StreamRecordsRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *ingestAPIStreamRecordsClient) Recv() (*StreamRecordsResponse, error) {
	m := new(StreamRecordsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// Deprecated: Do not use.
func (c *ingestAPIClient) IngestRecord(ctx context.Context, in *IngestRecordRequest, opts ...grpc.CallOption) (*IngestRecordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(IngestRecordResponse)
	err := c.cc.Invoke(ctx, IngestAPI_IngestRecord_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestAPIClient) BatchUpsertNodes(ctx context.Context, in *BatchUpsertNodesRequest, opts ...grpc.CallOption) (*BatchUpsertNodesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchUpsertNodesResponse)
	err := c.cc.Invoke(ctx, IngestAPI_BatchUpsertNodes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestAPIClient) BatchUpsertRelationships(ctx context.Context, in *BatchUpsertRelationshipsRequest, opts ...grpc.CallOption) (*BatchUpsertRelationshipsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchUpsertRelationshipsResponse)
	err := c.cc.Invoke(ctx, IngestAPI_BatchUpsertRelationships_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestAPIClient) BatchDeleteNodes(ctx context.Context, in *BatchDeleteNodesRequest, opts ...grpc.CallOption) (*BatchDeleteNodesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchDeleteNodesResponse)
	err := c.cc.Invoke(ctx, IngestAPI_BatchDeleteNodes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestAPIClient) BatchDeleteRelationships(ctx context.Context, in *BatchDeleteRelationshipsRequest, opts ...grpc.CallOption) (*BatchDeleteRelationshipsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchDeleteRelationshipsResponse)
	err := c.cc.Invoke(ctx, IngestAPI_BatchDeleteRelationships_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestAPIClient) BatchDeleteNodeProperties(ctx context.Context, in *BatchDeleteNodePropertiesRequest, opts ...grpc.CallOption) (*BatchDeleteNodePropertiesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchDeleteNodePropertiesResponse)
	err := c.cc.Invoke(ctx, IngestAPI_BatchDeleteNodeProperties_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestAPIClient) BatchDeleteRelationshipProperties(ctx context.Context, in *BatchDeleteRelationshipPropertiesRequest, opts ...grpc.CallOption) (*BatchDeleteRelationshipPropertiesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchDeleteRelationshipPropertiesResponse)
	err := c.cc.Invoke(ctx, IngestAPI_BatchDeleteRelationshipProperties_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ingestAPIClient) BatchDeleteNodeTags(ctx context.Context, in *BatchDeleteNodeTagsRequest, opts ...grpc.CallOption) (*BatchDeleteNodeTagsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BatchDeleteNodeTagsResponse)
	err := c.cc.Invoke(ctx, IngestAPI_BatchDeleteNodeTags_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IngestAPIServer is the server API for IngestAPI service.
// All implementations should embed UnimplementedIngestAPIServer
// for forward compatibility
//
// IngestAPI represents the service interface for data ingestion.
type IngestAPIServer interface {
	StreamRecords(IngestAPI_StreamRecordsServer) error
	// Deprecated: Do not use.
	IngestRecord(context.Context, *IngestRecordRequest) (*IngestRecordResponse, error)
	BatchUpsertNodes(context.Context, *BatchUpsertNodesRequest) (*BatchUpsertNodesResponse, error)
	BatchUpsertRelationships(context.Context, *BatchUpsertRelationshipsRequest) (*BatchUpsertRelationshipsResponse, error)
	BatchDeleteNodes(context.Context, *BatchDeleteNodesRequest) (*BatchDeleteNodesResponse, error)
	BatchDeleteRelationships(context.Context, *BatchDeleteRelationshipsRequest) (*BatchDeleteRelationshipsResponse, error)
	BatchDeleteNodeProperties(context.Context, *BatchDeleteNodePropertiesRequest) (*BatchDeleteNodePropertiesResponse, error)
	BatchDeleteRelationshipProperties(context.Context, *BatchDeleteRelationshipPropertiesRequest) (*BatchDeleteRelationshipPropertiesResponse, error)
	BatchDeleteNodeTags(context.Context, *BatchDeleteNodeTagsRequest) (*BatchDeleteNodeTagsResponse, error)
}

// UnimplementedIngestAPIServer should be embedded to have forward compatible implementations.
type UnimplementedIngestAPIServer struct {
}

func (UnimplementedIngestAPIServer) StreamRecords(IngestAPI_StreamRecordsServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRecords not implemented")
}
func (UnimplementedIngestAPIServer) IngestRecord(context.Context, *IngestRecordRequest) (*IngestRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IngestRecord not implemented")
}
func (UnimplementedIngestAPIServer) BatchUpsertNodes(context.Context, *BatchUpsertNodesRequest) (*BatchUpsertNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchUpsertNodes not implemented")
}
func (UnimplementedIngestAPIServer) BatchUpsertRelationships(context.Context, *BatchUpsertRelationshipsRequest) (*BatchUpsertRelationshipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchUpsertRelationships not implemented")
}
func (UnimplementedIngestAPIServer) BatchDeleteNodes(context.Context, *BatchDeleteNodesRequest) (*BatchDeleteNodesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDeleteNodes not implemented")
}
func (UnimplementedIngestAPIServer) BatchDeleteRelationships(context.Context, *BatchDeleteRelationshipsRequest) (*BatchDeleteRelationshipsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDeleteRelationships not implemented")
}
func (UnimplementedIngestAPIServer) BatchDeleteNodeProperties(context.Context, *BatchDeleteNodePropertiesRequest) (*BatchDeleteNodePropertiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDeleteNodeProperties not implemented")
}
func (UnimplementedIngestAPIServer) BatchDeleteRelationshipProperties(context.Context, *BatchDeleteRelationshipPropertiesRequest) (*BatchDeleteRelationshipPropertiesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDeleteRelationshipProperties not implemented")
}
func (UnimplementedIngestAPIServer) BatchDeleteNodeTags(context.Context, *BatchDeleteNodeTagsRequest) (*BatchDeleteNodeTagsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDeleteNodeTags not implemented")
}

// UnsafeIngestAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IngestAPIServer will
// result in compilation errors.
type UnsafeIngestAPIServer interface {
	mustEmbedUnimplementedIngestAPIServer()
}

func RegisterIngestAPIServer(s grpc.ServiceRegistrar, srv IngestAPIServer) {
	s.RegisterService(&IngestAPI_ServiceDesc, srv)
}

func _IngestAPI_StreamRecords_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(IngestAPIServer).StreamRecords(&ingestAPIStreamRecordsServer{ServerStream: stream})
}

type IngestAPI_StreamRecordsServer interface {
	Send(*StreamRecordsResponse) error
	Recv() (*StreamRecordsRequest, error)
	grpc.ServerStream
}

type ingestAPIStreamRecordsServer struct {
	grpc.ServerStream
}

func (x *ingestAPIStreamRecordsServer) Send(m *StreamRecordsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *ingestAPIStreamRecordsServer) Recv() (*StreamRecordsRequest, error) {
	m := new(StreamRecordsRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _IngestAPI_IngestRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IngestRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestAPIServer).IngestRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngestAPI_IngestRecord_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestAPIServer).IngestRecord(ctx, req.(*IngestRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngestAPI_BatchUpsertNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchUpsertNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestAPIServer).BatchUpsertNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngestAPI_BatchUpsertNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestAPIServer).BatchUpsertNodes(ctx, req.(*BatchUpsertNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngestAPI_BatchUpsertRelationships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchUpsertRelationshipsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestAPIServer).BatchUpsertRelationships(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngestAPI_BatchUpsertRelationships_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestAPIServer).BatchUpsertRelationships(ctx, req.(*BatchUpsertRelationshipsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngestAPI_BatchDeleteNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchDeleteNodesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestAPIServer).BatchDeleteNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngestAPI_BatchDeleteNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestAPIServer).BatchDeleteNodes(ctx, req.(*BatchDeleteNodesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngestAPI_BatchDeleteRelationships_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchDeleteRelationshipsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestAPIServer).BatchDeleteRelationships(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngestAPI_BatchDeleteRelationships_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestAPIServer).BatchDeleteRelationships(ctx, req.(*BatchDeleteRelationshipsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngestAPI_BatchDeleteNodeProperties_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchDeleteNodePropertiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestAPIServer).BatchDeleteNodeProperties(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngestAPI_BatchDeleteNodeProperties_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestAPIServer).BatchDeleteNodeProperties(ctx, req.(*BatchDeleteNodePropertiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngestAPI_BatchDeleteRelationshipProperties_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchDeleteRelationshipPropertiesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestAPIServer).BatchDeleteRelationshipProperties(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngestAPI_BatchDeleteRelationshipProperties_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestAPIServer).BatchDeleteRelationshipProperties(ctx, req.(*BatchDeleteRelationshipPropertiesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IngestAPI_BatchDeleteNodeTags_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchDeleteNodeTagsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IngestAPIServer).BatchDeleteNodeTags(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IngestAPI_BatchDeleteNodeTags_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IngestAPIServer).BatchDeleteNodeTags(ctx, req.(*BatchDeleteNodeTagsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IngestAPI_ServiceDesc is the grpc.ServiceDesc for IngestAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IngestAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "indykite.ingest.v1beta3.IngestAPI",
	HandlerType: (*IngestAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IngestRecord",
			Handler:    _IngestAPI_IngestRecord_Handler,
		},
		{
			MethodName: "BatchUpsertNodes",
			Handler:    _IngestAPI_BatchUpsertNodes_Handler,
		},
		{
			MethodName: "BatchUpsertRelationships",
			Handler:    _IngestAPI_BatchUpsertRelationships_Handler,
		},
		{
			MethodName: "BatchDeleteNodes",
			Handler:    _IngestAPI_BatchDeleteNodes_Handler,
		},
		{
			MethodName: "BatchDeleteRelationships",
			Handler:    _IngestAPI_BatchDeleteRelationships_Handler,
		},
		{
			MethodName: "BatchDeleteNodeProperties",
			Handler:    _IngestAPI_BatchDeleteNodeProperties_Handler,
		},
		{
			MethodName: "BatchDeleteRelationshipProperties",
			Handler:    _IngestAPI_BatchDeleteRelationshipProperties_Handler,
		},
		{
			MethodName: "BatchDeleteNodeTags",
			Handler:    _IngestAPI_BatchDeleteNodeTags_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamRecords",
			Handler:       _IngestAPI_StreamRecords_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "indykite/ingest/v1beta3/ingest_api.proto",
}
