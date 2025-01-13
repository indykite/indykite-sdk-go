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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        (unknown)
// source: indykite/entitymatching/v1beta1/model.proto

package entitymatchingv1beta1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	_ "github.com/envoyproxy/protoc-gen-validate/validate"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type PipelineStatus int32

const (
	PipelineStatus_PIPELINE_STATUS_STATUS_INVALID     PipelineStatus = 0
	PipelineStatus_PIPELINE_STATUS_STATUS_PENDING     PipelineStatus = 1
	PipelineStatus_PIPELINE_STATUS_STATUS_IN_PROGRESS PipelineStatus = 2
	PipelineStatus_PIPELINE_STATUS_STATUS_SUCCESS     PipelineStatus = 3
	PipelineStatus_PIPELINE_STATUS_STATUS_ERROR       PipelineStatus = 4
)

// Enum value maps for PipelineStatus.
var (
	PipelineStatus_name = map[int32]string{
		0: "PIPELINE_STATUS_STATUS_INVALID",
		1: "PIPELINE_STATUS_STATUS_PENDING",
		2: "PIPELINE_STATUS_STATUS_IN_PROGRESS",
		3: "PIPELINE_STATUS_STATUS_SUCCESS",
		4: "PIPELINE_STATUS_STATUS_ERROR",
	}
	PipelineStatus_value = map[string]int32{
		"PIPELINE_STATUS_STATUS_INVALID":     0,
		"PIPELINE_STATUS_STATUS_PENDING":     1,
		"PIPELINE_STATUS_STATUS_IN_PROGRESS": 2,
		"PIPELINE_STATUS_STATUS_SUCCESS":     3,
		"PIPELINE_STATUS_STATUS_ERROR":       4,
	}
)

func (x PipelineStatus) Enum() *PipelineStatus {
	p := new(PipelineStatus)
	*p = x
	return p
}

func (x PipelineStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PipelineStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_indykite_entitymatching_v1beta1_model_proto_enumTypes[0].Descriptor()
}

func (PipelineStatus) Type() protoreflect.EnumType {
	return &file_indykite_entitymatching_v1beta1_model_proto_enumTypes[0]
}

func (x PipelineStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PipelineStatus.Descriptor instead.
func (PipelineStatus) EnumDescriptor() ([]byte, []int) {
	return file_indykite_entitymatching_v1beta1_model_proto_rawDescGZIP(), []int{0}
}

type PropertyMapping struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// SourceNodeType is the type of the node that will be compared to nodes of TargetNodeType.
	SourceNodeType string `protobuf:"bytes,1,opt,name=source_node_type,json=sourceNodeType,proto3" json:"source_node_type,omitempty"`
	// SourceNodeProperty is a property of the source node that will be compared to TargetNodeProperty.
	SourceNodeProperty string `protobuf:"bytes,2,opt,name=source_node_property,json=sourceNodeProperty,proto3" json:"source_node_property,omitempty"`
	// TargetNodeType is the type of the node that will be compared to nodes of SourceNodeType.
	TargetNodeType string `protobuf:"bytes,3,opt,name=target_node_type,json=targetNodeType,proto3" json:"target_node_type,omitempty"`
	// TargetNodeProperty is a property of the source node that will be compared to SourceNodeProperty.
	TargetNodeProperty string `protobuf:"bytes,4,opt,name=target_node_property,json=targetNodeProperty,proto3" json:"target_node_property,omitempty"`
	// SimilarityScoreCutoff defines the threshold (in range [0,1]), above which entities will be automatically matched.
	SimilarityScoreCutoff float32 `protobuf:"fixed32,5,opt,name=similarity_score_cutoff,json=similarityScoreCutoff,proto3" json:"similarity_score_cutoff,omitempty"`
	unknownFields         protoimpl.UnknownFields
	sizeCache             protoimpl.SizeCache
}

func (x *PropertyMapping) Reset() {
	*x = PropertyMapping{}
	mi := &file_indykite_entitymatching_v1beta1_model_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PropertyMapping) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PropertyMapping) ProtoMessage() {}

func (x *PropertyMapping) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_entitymatching_v1beta1_model_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PropertyMapping.ProtoReflect.Descriptor instead.
func (*PropertyMapping) Descriptor() ([]byte, []int) {
	return file_indykite_entitymatching_v1beta1_model_proto_rawDescGZIP(), []int{0}
}

func (x *PropertyMapping) GetSourceNodeType() string {
	if x != nil {
		return x.SourceNodeType
	}
	return ""
}

func (x *PropertyMapping) GetSourceNodeProperty() string {
	if x != nil {
		return x.SourceNodeProperty
	}
	return ""
}

func (x *PropertyMapping) GetTargetNodeType() string {
	if x != nil {
		return x.TargetNodeType
	}
	return ""
}

func (x *PropertyMapping) GetTargetNodeProperty() string {
	if x != nil {
		return x.TargetNodeProperty
	}
	return ""
}

func (x *PropertyMapping) GetSimilarityScoreCutoff() float32 {
	if x != nil {
		return x.SimilarityScoreCutoff
	}
	return 0
}

type CustomPropertyMappings struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// SourceNodeProperty is a property of the source node that will be compared to TargetNodeProperty.
	SourceNodeProperty string `protobuf:"bytes,2,opt,name=source_node_property,json=sourceNodeProperty,proto3" json:"source_node_property,omitempty"`
	// TargetNodeProperty is a property of the source node that will be compared to SourceNodeProperty.
	TargetNodeProperty string `protobuf:"bytes,4,opt,name=target_node_property,json=targetNodeProperty,proto3" json:"target_node_property,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *CustomPropertyMappings) Reset() {
	*x = CustomPropertyMappings{}
	mi := &file_indykite_entitymatching_v1beta1_model_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *CustomPropertyMappings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CustomPropertyMappings) ProtoMessage() {}

func (x *CustomPropertyMappings) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_entitymatching_v1beta1_model_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CustomPropertyMappings.ProtoReflect.Descriptor instead.
func (*CustomPropertyMappings) Descriptor() ([]byte, []int) {
	return file_indykite_entitymatching_v1beta1_model_proto_rawDescGZIP(), []int{1}
}

func (x *CustomPropertyMappings) GetSourceNodeProperty() string {
	if x != nil {
		return x.SourceNodeProperty
	}
	return ""
}

func (x *CustomPropertyMappings) GetTargetNodeProperty() string {
	if x != nil {
		return x.TargetNodeProperty
	}
	return ""
}

var File_indykite_entitymatching_v1beta1_model_proto protoreflect.FileDescriptor

var file_indykite_entitymatching_v1beta1_model_proto_rawDesc = []byte{
	0x0a, 0x2b, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2f, 0x6d, 0x6f, 0x64, 0x65, 0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1f, 0x69,
	0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x6d, 0x61,
	0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x17,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x94, 0x02, 0x0a, 0x0f, 0x50, 0x72, 0x6f, 0x70,
	0x65, 0x72, 0x74, 0x79, 0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x12, 0x28, 0x0a, 0x10, 0x73,
	0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0e, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x6f, 0x64,
	0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x30, 0x0a, 0x14, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f,
	0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x12, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x6f, 0x64, 0x65, 0x50,
	0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x12, 0x28, 0x0a, 0x10, 0x74, 0x61, 0x72, 0x67, 0x65,
	0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0e, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x30, 0x0a, 0x14, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65,
	0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x12, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4e, 0x6f, 0x64, 0x65, 0x50, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x79, 0x12, 0x49, 0x0a, 0x17, 0x73, 0x69, 0x6d, 0x69, 0x6c, 0x61, 0x72, 0x69, 0x74,
	0x79, 0x5f, 0x73, 0x63, 0x6f, 0x72, 0x65, 0x5f, 0x63, 0x75, 0x74, 0x6f, 0x66, 0x66, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x02, 0x42, 0x11, 0xfa, 0x42, 0x0e, 0x0a, 0x0c, 0x1d, 0x00, 0x00, 0x80, 0x3f,
	0x2d, 0x00, 0x00, 0x00, 0x00, 0x40, 0x01, 0x52, 0x15, 0x73, 0x69, 0x6d, 0x69, 0x6c, 0x61, 0x72,
	0x69, 0x74, 0x79, 0x53, 0x63, 0x6f, 0x72, 0x65, 0x43, 0x75, 0x74, 0x6f, 0x66, 0x66, 0x22, 0x7c,
	0x0a, 0x16, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79,
	0x4d, 0x61, 0x70, 0x70, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x30, 0x0a, 0x14, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4e, 0x6f,
	0x64, 0x65, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x12, 0x30, 0x0a, 0x14, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x5f, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72,
	0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x12, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74,
	0x4e, 0x6f, 0x64, 0x65, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x2a, 0xc6, 0x01, 0x0a,
	0x0e, 0x50, 0x69, 0x70, 0x65, 0x6c, 0x69, 0x6e, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12,
	0x22, 0x0a, 0x1e, 0x50, 0x49, 0x50, 0x45, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49,
	0x44, 0x10, 0x00, 0x12, 0x22, 0x0a, 0x1e, 0x50, 0x49, 0x50, 0x45, 0x4c, 0x49, 0x4e, 0x45, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x50, 0x45,
	0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x26, 0x0a, 0x22, 0x50, 0x49, 0x50, 0x45, 0x4c,
	0x49, 0x4e, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55,
	0x53, 0x5f, 0x49, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x47, 0x52, 0x45, 0x53, 0x53, 0x10, 0x02, 0x12,
	0x22, 0x0a, 0x1e, 0x50, 0x49, 0x50, 0x45, 0x4c, 0x49, 0x4e, 0x45, 0x5f, 0x53, 0x54, 0x41, 0x54,
	0x55, 0x53, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53,
	0x53, 0x10, 0x03, 0x12, 0x20, 0x0a, 0x1c, 0x50, 0x49, 0x50, 0x45, 0x4c, 0x49, 0x4e, 0x45, 0x5f,
	0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x10, 0x04, 0x42, 0xae, 0x02, 0x0a, 0x23, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e,
	0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x6d, 0x61, 0x74,
	0x63, 0x68, 0x69, 0x6e, 0x67, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x42, 0x0a, 0x4d,
	0x6f, 0x64, 0x65, 0x6c, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x5d, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65,
	0x2f, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2d, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f,
	0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x2f, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x3b, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x6d, 0x61, 0x74, 0x63, 0x68,
	0x69, 0x6e, 0x67, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x45, 0x58,
	0xaa, 0x02, 0x1f, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x45, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0xca, 0x02, 0x1f, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x5c, 0x45, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x5c, 0x56, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0xe2, 0x02, 0x2b, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x5c,
	0x45, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x5c, 0x56,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0xea, 0x02, 0x21, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x3a, 0x3a, 0x45,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x69, 0x6e, 0x67, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_indykite_entitymatching_v1beta1_model_proto_rawDescOnce sync.Once
	file_indykite_entitymatching_v1beta1_model_proto_rawDescData = file_indykite_entitymatching_v1beta1_model_proto_rawDesc
)

func file_indykite_entitymatching_v1beta1_model_proto_rawDescGZIP() []byte {
	file_indykite_entitymatching_v1beta1_model_proto_rawDescOnce.Do(func() {
		file_indykite_entitymatching_v1beta1_model_proto_rawDescData = protoimpl.X.CompressGZIP(file_indykite_entitymatching_v1beta1_model_proto_rawDescData)
	})
	return file_indykite_entitymatching_v1beta1_model_proto_rawDescData
}

var file_indykite_entitymatching_v1beta1_model_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_indykite_entitymatching_v1beta1_model_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_indykite_entitymatching_v1beta1_model_proto_goTypes = []any{
	(PipelineStatus)(0),            // 0: indykite.entitymatching.v1beta1.PipelineStatus
	(*PropertyMapping)(nil),        // 1: indykite.entitymatching.v1beta1.PropertyMapping
	(*CustomPropertyMappings)(nil), // 2: indykite.entitymatching.v1beta1.CustomPropertyMappings
}
var file_indykite_entitymatching_v1beta1_model_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_indykite_entitymatching_v1beta1_model_proto_init() }
func file_indykite_entitymatching_v1beta1_model_proto_init() {
	if File_indykite_entitymatching_v1beta1_model_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_indykite_entitymatching_v1beta1_model_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_indykite_entitymatching_v1beta1_model_proto_goTypes,
		DependencyIndexes: file_indykite_entitymatching_v1beta1_model_proto_depIdxs,
		EnumInfos:         file_indykite_entitymatching_v1beta1_model_proto_enumTypes,
		MessageInfos:      file_indykite_entitymatching_v1beta1_model_proto_msgTypes,
	}.Build()
	File_indykite_entitymatching_v1beta1_model_proto = out.File
	file_indykite_entitymatching_v1beta1_model_proto_rawDesc = nil
	file_indykite_entitymatching_v1beta1_model_proto_goTypes = nil
	file_indykite_entitymatching_v1beta1_model_proto_depIdxs = nil
}
