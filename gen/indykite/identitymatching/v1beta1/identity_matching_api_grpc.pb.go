// Copyright (c) 2024 IndyKite
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

// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.4.0
// - protoc             (unknown)
// source: indykite/identitymatching/v1beta1/identity_matching_api.proto

package identitymatchingv1beta1

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
	IdentityMatchingAPI_RunIdentityMatchingPipeline_FullMethodName  = "/indykite.identitymatching.v1beta1.IdentityMatchingAPI/RunIdentityMatchingPipeline"
	IdentityMatchingAPI_ReadSuggestedPropertyMapping_FullMethodName = "/indykite.identitymatching.v1beta1.IdentityMatchingAPI/ReadSuggestedPropertyMapping"
	IdentityMatchingAPI_ReadEntityMatchingReport_FullMethodName     = "/indykite.identitymatching.v1beta1.IdentityMatchingAPI/ReadEntityMatchingReport"
)

// IdentityMatchingAPIClient is the client API for IdentityMatchingAPI service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type IdentityMatchingAPIClient interface {
	// RunIdentityMatchingPipeline by Pipeline ID and optional property mapping.
	RunIdentityMatchingPipeline(ctx context.Context, in *RunIdentityMatchingPipelineRequest, opts ...grpc.CallOption) (*RunIdentityMatchingPipelineResponse, error)
	// ReadSuggestedPropertyMapping by Pipeline Name or ID.
	ReadSuggestedPropertyMapping(ctx context.Context, in *ReadSuggestedPropertyMappingRequest, opts ...grpc.CallOption) (*ReadSuggestedPropertyMappingResponse, error)
	// ReadEntityMatchingReport by Pipeline Name or ID for a successful Pipeline.
	ReadEntityMatchingReport(ctx context.Context, in *ReadEntityMatchingReportRequest, opts ...grpc.CallOption) (*ReadEntityMatchingReportResponse, error)
}

type identityMatchingAPIClient struct {
	cc grpc.ClientConnInterface
}

func NewIdentityMatchingAPIClient(cc grpc.ClientConnInterface) IdentityMatchingAPIClient {
	return &identityMatchingAPIClient{cc}
}

func (c *identityMatchingAPIClient) RunIdentityMatchingPipeline(ctx context.Context, in *RunIdentityMatchingPipelineRequest, opts ...grpc.CallOption) (*RunIdentityMatchingPipelineResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RunIdentityMatchingPipelineResponse)
	err := c.cc.Invoke(ctx, IdentityMatchingAPI_RunIdentityMatchingPipeline_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityMatchingAPIClient) ReadSuggestedPropertyMapping(ctx context.Context, in *ReadSuggestedPropertyMappingRequest, opts ...grpc.CallOption) (*ReadSuggestedPropertyMappingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadSuggestedPropertyMappingResponse)
	err := c.cc.Invoke(ctx, IdentityMatchingAPI_ReadSuggestedPropertyMapping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *identityMatchingAPIClient) ReadEntityMatchingReport(ctx context.Context, in *ReadEntityMatchingReportRequest, opts ...grpc.CallOption) (*ReadEntityMatchingReportResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReadEntityMatchingReportResponse)
	err := c.cc.Invoke(ctx, IdentityMatchingAPI_ReadEntityMatchingReport_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IdentityMatchingAPIServer is the server API for IdentityMatchingAPI service.
// All implementations should embed UnimplementedIdentityMatchingAPIServer
// for forward compatibility
type IdentityMatchingAPIServer interface {
	// RunIdentityMatchingPipeline by Pipeline ID and optional property mapping.
	RunIdentityMatchingPipeline(context.Context, *RunIdentityMatchingPipelineRequest) (*RunIdentityMatchingPipelineResponse, error)
	// ReadSuggestedPropertyMapping by Pipeline Name or ID.
	ReadSuggestedPropertyMapping(context.Context, *ReadSuggestedPropertyMappingRequest) (*ReadSuggestedPropertyMappingResponse, error)
	// ReadEntityMatchingReport by Pipeline Name or ID for a successful Pipeline.
	ReadEntityMatchingReport(context.Context, *ReadEntityMatchingReportRequest) (*ReadEntityMatchingReportResponse, error)
}

// UnimplementedIdentityMatchingAPIServer should be embedded to have forward compatible implementations.
type UnimplementedIdentityMatchingAPIServer struct {
}

func (UnimplementedIdentityMatchingAPIServer) RunIdentityMatchingPipeline(context.Context, *RunIdentityMatchingPipelineRequest) (*RunIdentityMatchingPipelineResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RunIdentityMatchingPipeline not implemented")
}
func (UnimplementedIdentityMatchingAPIServer) ReadSuggestedPropertyMapping(context.Context, *ReadSuggestedPropertyMappingRequest) (*ReadSuggestedPropertyMappingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadSuggestedPropertyMapping not implemented")
}
func (UnimplementedIdentityMatchingAPIServer) ReadEntityMatchingReport(context.Context, *ReadEntityMatchingReportRequest) (*ReadEntityMatchingReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadEntityMatchingReport not implemented")
}

// UnsafeIdentityMatchingAPIServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to IdentityMatchingAPIServer will
// result in compilation errors.
type UnsafeIdentityMatchingAPIServer interface {
	mustEmbedUnimplementedIdentityMatchingAPIServer()
}

func RegisterIdentityMatchingAPIServer(s grpc.ServiceRegistrar, srv IdentityMatchingAPIServer) {
	s.RegisterService(&IdentityMatchingAPI_ServiceDesc, srv)
}

func _IdentityMatchingAPI_RunIdentityMatchingPipeline_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RunIdentityMatchingPipelineRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityMatchingAPIServer).RunIdentityMatchingPipeline(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityMatchingAPI_RunIdentityMatchingPipeline_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityMatchingAPIServer).RunIdentityMatchingPipeline(ctx, req.(*RunIdentityMatchingPipelineRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityMatchingAPI_ReadSuggestedPropertyMapping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadSuggestedPropertyMappingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityMatchingAPIServer).ReadSuggestedPropertyMapping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityMatchingAPI_ReadSuggestedPropertyMapping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityMatchingAPIServer).ReadSuggestedPropertyMapping(ctx, req.(*ReadSuggestedPropertyMappingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _IdentityMatchingAPI_ReadEntityMatchingReport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadEntityMatchingReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IdentityMatchingAPIServer).ReadEntityMatchingReport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: IdentityMatchingAPI_ReadEntityMatchingReport_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IdentityMatchingAPIServer).ReadEntityMatchingReport(ctx, req.(*ReadEntityMatchingReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// IdentityMatchingAPI_ServiceDesc is the grpc.ServiceDesc for IdentityMatchingAPI service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var IdentityMatchingAPI_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "indykite.identitymatching.v1beta1.IdentityMatchingAPI",
	HandlerType: (*IdentityMatchingAPIServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RunIdentityMatchingPipeline",
			Handler:    _IdentityMatchingAPI_RunIdentityMatchingPipeline_Handler,
		},
		{
			MethodName: "ReadSuggestedPropertyMapping",
			Handler:    _IdentityMatchingAPI_ReadSuggestedPropertyMapping_Handler,
		},
		{
			MethodName: "ReadEntityMatchingReport",
			Handler:    _IdentityMatchingAPI_ReadEntityMatchingReport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "indykite/identitymatching/v1beta1/identity_matching_api.proto",
}
