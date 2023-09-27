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

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: indykite/auditsink/v1beta1/ingest.proto

package auditsinkv1beta1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type UpsertData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*UpsertData_Node
	//	*UpsertData_Relation
	Data isUpsertData_Data `protobuf_oneof:"data"`
}

func (x *UpsertData) Reset() {
	*x = UpsertData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpsertData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertData) ProtoMessage() {}

func (x *UpsertData) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpsertData.ProtoReflect.Descriptor instead.
func (*UpsertData) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{0}
}

func (m *UpsertData) GetData() isUpsertData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *UpsertData) GetNode() *Node {
	if x, ok := x.GetData().(*UpsertData_Node); ok {
		return x.Node
	}
	return nil
}

func (x *UpsertData) GetRelation() *Relation {
	if x, ok := x.GetData().(*UpsertData_Relation); ok {
		return x.Relation
	}
	return nil
}

type isUpsertData_Data interface {
	isUpsertData_Data()
}

type UpsertData_Node struct {
	Node *Node `protobuf:"bytes,1,opt,name=node,proto3,oneof"`
}

type UpsertData_Relation struct {
	Relation *Relation `protobuf:"bytes,2,opt,name=relation,proto3,oneof"`
}

func (*UpsertData_Node) isUpsertData_Data() {}

func (*UpsertData_Relation) isUpsertData_Data() {}

type DeleteData struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Data:
	//
	//	*DeleteData_Node
	//	*DeleteData_Relation
	//	*DeleteData_NodeProperty
	//	*DeleteData_RelationProperty
	Data isDeleteData_Data `protobuf_oneof:"data"`
}

func (x *DeleteData) Reset() {
	*x = DeleteData{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteData) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteData) ProtoMessage() {}

func (x *DeleteData) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteData.ProtoReflect.Descriptor instead.
func (*DeleteData) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{1}
}

func (m *DeleteData) GetData() isDeleteData_Data {
	if m != nil {
		return m.Data
	}
	return nil
}

func (x *DeleteData) GetNode() *NodeMatch {
	if x, ok := x.GetData().(*DeleteData_Node); ok {
		return x.Node
	}
	return nil
}

func (x *DeleteData) GetRelation() *RelationMatch {
	if x, ok := x.GetData().(*DeleteData_Relation); ok {
		return x.Relation
	}
	return nil
}

func (x *DeleteData) GetNodeProperty() *DeleteData_NodePropertyMatch {
	if x, ok := x.GetData().(*DeleteData_NodeProperty); ok {
		return x.NodeProperty
	}
	return nil
}

func (x *DeleteData) GetRelationProperty() *DeleteData_RelationPropertyMatch {
	if x, ok := x.GetData().(*DeleteData_RelationProperty); ok {
		return x.RelationProperty
	}
	return nil
}

type isDeleteData_Data interface {
	isDeleteData_Data()
}

type DeleteData_Node struct {
	Node *NodeMatch `protobuf:"bytes,1,opt,name=node,proto3,oneof"`
}

type DeleteData_Relation struct {
	Relation *RelationMatch `protobuf:"bytes,2,opt,name=relation,proto3,oneof"`
}

type DeleteData_NodeProperty struct {
	NodeProperty *DeleteData_NodePropertyMatch `protobuf:"bytes,3,opt,name=node_property,json=nodeProperty,proto3,oneof"`
}

type DeleteData_RelationProperty struct {
	RelationProperty *DeleteData_RelationPropertyMatch `protobuf:"bytes,4,opt,name=relation_property,json=relationProperty,proto3,oneof"`
}

func (*DeleteData_Node) isDeleteData_Data() {}

func (*DeleteData_Relation) isDeleteData_Data() {}

func (*DeleteData_NodeProperty) isDeleteData_Data() {}

func (*DeleteData_RelationProperty) isDeleteData_Data() {}

type DigitalTwin struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExternalId string `protobuf:"bytes,1,opt,name=external_id,json=externalId,proto3" json:"external_id,omitempty"`
	Type       string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
	TenantId   string `protobuf:"bytes,3,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
}

func (x *DigitalTwin) Reset() {
	*x = DigitalTwin{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DigitalTwin) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DigitalTwin) ProtoMessage() {}

func (x *DigitalTwin) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DigitalTwin.ProtoReflect.Descriptor instead.
func (*DigitalTwin) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{2}
}

func (x *DigitalTwin) GetExternalId() string {
	if x != nil {
		return x.ExternalId
	}
	return ""
}

func (x *DigitalTwin) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *DigitalTwin) GetTenantId() string {
	if x != nil {
		return x.TenantId
	}
	return ""
}

type Resource struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExternalId string `protobuf:"bytes,1,opt,name=external_id,json=externalId,proto3" json:"external_id,omitempty"`
	Type       string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *Resource) Reset() {
	*x = Resource{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Resource) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Resource) ProtoMessage() {}

func (x *Resource) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Resource.ProtoReflect.Descriptor instead.
func (*Resource) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{3}
}

func (x *Resource) GetExternalId() string {
	if x != nil {
		return x.ExternalId
	}
	return ""
}

func (x *Resource) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type Node struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Type:
	//
	//	*Node_DigitalTwin
	//	*Node_Resource
	Type isNode_Type `protobuf_oneof:"type"`
}

func (x *Node) Reset() {
	*x = Node{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Node) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Node) ProtoMessage() {}

func (x *Node) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Node.ProtoReflect.Descriptor instead.
func (*Node) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{4}
}

func (m *Node) GetType() isNode_Type {
	if m != nil {
		return m.Type
	}
	return nil
}

func (x *Node) GetDigitalTwin() *DigitalTwin {
	if x, ok := x.GetType().(*Node_DigitalTwin); ok {
		return x.DigitalTwin
	}
	return nil
}

func (x *Node) GetResource() *Resource {
	if x, ok := x.GetType().(*Node_Resource); ok {
		return x.Resource
	}
	return nil
}

type isNode_Type interface {
	isNode_Type()
}

type Node_DigitalTwin struct {
	DigitalTwin *DigitalTwin `protobuf:"bytes,1,opt,name=digital_twin,json=digitalTwin,proto3,oneof"`
}

type Node_Resource struct {
	Resource *Resource `protobuf:"bytes,2,opt,name=resource,proto3,oneof"`
}

func (*Node_DigitalTwin) isNode_Type() {}

func (*Node_Resource) isNode_Type() {}

type Relation struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Match *RelationMatch `protobuf:"bytes,1,opt,name=match,proto3" json:"match,omitempty"`
}

func (x *Relation) Reset() {
	*x = Relation{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Relation) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Relation) ProtoMessage() {}

func (x *Relation) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Relation.ProtoReflect.Descriptor instead.
func (*Relation) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{5}
}

func (x *Relation) GetMatch() *RelationMatch {
	if x != nil {
		return x.Match
	}
	return nil
}

type NodeMatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ExternalId string `protobuf:"bytes,1,opt,name=external_id,json=externalId,proto3" json:"external_id,omitempty"`
	Type       string `protobuf:"bytes,2,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *NodeMatch) Reset() {
	*x = NodeMatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NodeMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NodeMatch) ProtoMessage() {}

func (x *NodeMatch) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NodeMatch.ProtoReflect.Descriptor instead.
func (*NodeMatch) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{6}
}

func (x *NodeMatch) GetExternalId() string {
	if x != nil {
		return x.ExternalId
	}
	return ""
}

func (x *NodeMatch) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type RelationMatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SourceMatch *NodeMatch `protobuf:"bytes,1,opt,name=source_match,json=sourceMatch,proto3" json:"source_match,omitempty"`
	TargetMatch *NodeMatch `protobuf:"bytes,2,opt,name=target_match,json=targetMatch,proto3" json:"target_match,omitempty"`
	Type        string     `protobuf:"bytes,3,opt,name=type,proto3" json:"type,omitempty"`
}

func (x *RelationMatch) Reset() {
	*x = RelationMatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RelationMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RelationMatch) ProtoMessage() {}

func (x *RelationMatch) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RelationMatch.ProtoReflect.Descriptor instead.
func (*RelationMatch) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{7}
}

func (x *RelationMatch) GetSourceMatch() *NodeMatch {
	if x != nil {
		return x.SourceMatch
	}
	return nil
}

func (x *RelationMatch) GetTargetMatch() *NodeMatch {
	if x != nil {
		return x.TargetMatch
	}
	return nil
}

func (x *RelationMatch) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

type DeleteData_NodePropertyMatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Match *NodeMatch `protobuf:"bytes,1,opt,name=match,proto3" json:"match,omitempty"`
	Key   string     `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeleteData_NodePropertyMatch) Reset() {
	*x = DeleteData_NodePropertyMatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteData_NodePropertyMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteData_NodePropertyMatch) ProtoMessage() {}

func (x *DeleteData_NodePropertyMatch) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteData_NodePropertyMatch.ProtoReflect.Descriptor instead.
func (*DeleteData_NodePropertyMatch) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{1, 0}
}

func (x *DeleteData_NodePropertyMatch) GetMatch() *NodeMatch {
	if x != nil {
		return x.Match
	}
	return nil
}

func (x *DeleteData_NodePropertyMatch) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type DeleteData_RelationPropertyMatch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Match *RelationMatch `protobuf:"bytes,1,opt,name=match,proto3" json:"match,omitempty"`
	Key   string         `protobuf:"bytes,2,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *DeleteData_RelationPropertyMatch) Reset() {
	*x = DeleteData_RelationPropertyMatch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteData_RelationPropertyMatch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteData_RelationPropertyMatch) ProtoMessage() {}

func (x *DeleteData_RelationPropertyMatch) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteData_RelationPropertyMatch.ProtoReflect.Descriptor instead.
func (*DeleteData_RelationPropertyMatch) Descriptor() ([]byte, []int) {
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP(), []int{1, 1}
}

func (x *DeleteData_RelationPropertyMatch) GetMatch() *RelationMatch {
	if x != nil {
		return x.Match
	}
	return nil
}

func (x *DeleteData_RelationPropertyMatch) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

var File_indykite_auditsink_v1beta1_ingest_proto protoreflect.FileDescriptor

var file_indykite_auditsink_v1beta1_ingest_proto_rawDesc = []byte{
	0x0a, 0x27, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74,
	0x73, 0x69, 0x6e, 0x6b, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x69, 0x6e, 0x67,
	0x65, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x1a, 0x69, 0x6e, 0x64, 0x79, 0x6b,
	0x69, 0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x22, 0x90, 0x01, 0x0a, 0x0a, 0x55, 0x70, 0x73, 0x65, 0x72, 0x74,
	0x44, 0x61, 0x74, 0x61, 0x12, 0x36, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x20, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x75,
	0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e,
	0x4e, 0x6f, 0x64, 0x65, 0x48, 0x00, 0x52, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x12, 0x42, 0x0a, 0x08,
	0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73,
	0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x42, 0x06, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22, 0xb8, 0x04, 0x0a, 0x0a, 0x44, 0x65, 0x6c,
	0x65, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x12, 0x3b, 0x0a, 0x04, 0x6e, 0x6f, 0x64, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65,
	0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x48, 0x00, 0x52, 0x04,
	0x6e, 0x6f, 0x64, 0x65, 0x12, 0x47, 0x0a, 0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74,
	0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x74, 0x63,
	0x68, 0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x5f, 0x0a,
	0x0d, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x38, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e,
	0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74, 0x61, 0x2e, 0x4e, 0x6f, 0x64,
	0x65, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x48, 0x00,
	0x52, 0x0c, 0x6e, 0x6f, 0x64, 0x65, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x12, 0x6b,
	0x0a, 0x11, 0x72, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x70, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x79, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x3c, 0x2e, 0x69, 0x6e, 0x64, 0x79,
	0x6b, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x44, 0x61, 0x74,
	0x61, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72,
	0x74, 0x79, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x48, 0x00, 0x52, 0x10, 0x72, 0x65, 0x6c, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x1a, 0x62, 0x0a, 0x11, 0x4e,
	0x6f, 0x64, 0x65, 0x50, 0x72, 0x6f, 0x70, 0x65, 0x72, 0x74, 0x79, 0x4d, 0x61, 0x74, 0x63, 0x68,
	0x12, 0x3b, 0x0a, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74,
	0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4e, 0x6f, 0x64,
	0x65, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x1a,
	0x6a, 0x0a, 0x15, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x50, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x79, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x3f, 0x0a, 0x05, 0x6d, 0x61, 0x74, 0x63,
	0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69,
	0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x52, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x42, 0x06, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x22, 0x5f, 0x0a, 0x0b, 0x44, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x54, 0x77,
	0x69, 0x6e, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x1b, 0x0a, 0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61,
	0x6e, 0x74, 0x49, 0x64, 0x22, 0x3f, 0x0a, 0x08, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x49,
	0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xa0, 0x01, 0x0a, 0x04, 0x4e, 0x6f, 0x64, 0x65, 0x12, 0x4c,
	0x0a, 0x0c, 0x64, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x5f, 0x74, 0x77, 0x69, 0x6e, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e,
	0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x44, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x54, 0x77, 0x69, 0x6e, 0x48, 0x00, 0x52,
	0x0b, 0x64, 0x69, 0x67, 0x69, 0x74, 0x61, 0x6c, 0x54, 0x77, 0x69, 0x6e, 0x12, 0x42, 0x0a, 0x08,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73,
	0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f,
	0x75, 0x72, 0x63, 0x65, 0x48, 0x00, 0x52, 0x08, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65,
	0x42, 0x06, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0x4b, 0x0a, 0x08, 0x52, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x3f, 0x0a, 0x05, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x61,
	0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2e, 0x52, 0x65, 0x6c, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x05,
	0x6d, 0x61, 0x74, 0x63, 0x68, 0x22, 0x40, 0x0a, 0x09, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x74,
	0x63, 0x68, 0x12, 0x1f, 0x0a, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61,
	0x6c, 0x49, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x22, 0xb7, 0x01, 0x0a, 0x0d, 0x52, 0x65, 0x6c, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x48, 0x0a, 0x0c, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x6d, 0x61, 0x74, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x25, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74,
	0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4e, 0x6f, 0x64,
	0x65, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x52, 0x0b, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x4d, 0x61,
	0x74, 0x63, 0x68, 0x12, 0x48, 0x0a, 0x0c, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x5f, 0x6d, 0x61,
	0x74, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x25, 0x2e, 0x69, 0x6e, 0x64, 0x79,
	0x6b, 0x69, 0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x4d, 0x61, 0x74, 0x63, 0x68,
	0x52, 0x0b, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x4d, 0x61, 0x74, 0x63, 0x68, 0x12, 0x12, 0x0a,
	0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x42, 0x8c, 0x02, 0x0a, 0x1e, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69,
	0x74, 0x65, 0x2e, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x42, 0x0b, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x50, 0x01, 0x5a, 0x53, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f,
	0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74,
	0x65, 0x2d, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x69, 0x6e, 0x64,
	0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x2f,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x61, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e,
	0x6b, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x41, 0x58, 0xaa, 0x02,
	0x1a, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x41, 0x75, 0x64, 0x69, 0x74, 0x73,
	0x69, 0x6e, 0x6b, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xca, 0x02, 0x1a, 0x49, 0x6e,
	0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x5c, 0x41, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b,
	0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xe2, 0x02, 0x26, 0x49, 0x6e, 0x64, 0x79, 0x6b,
	0x69, 0x74, 0x65, 0x5c, 0x41, 0x75, 0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x5c, 0x56, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x1c, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x3a, 0x3a, 0x41, 0x75,
	0x64, 0x69, 0x74, 0x73, 0x69, 0x6e, 0x6b, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_indykite_auditsink_v1beta1_ingest_proto_rawDescOnce sync.Once
	file_indykite_auditsink_v1beta1_ingest_proto_rawDescData = file_indykite_auditsink_v1beta1_ingest_proto_rawDesc
)

func file_indykite_auditsink_v1beta1_ingest_proto_rawDescGZIP() []byte {
	file_indykite_auditsink_v1beta1_ingest_proto_rawDescOnce.Do(func() {
		file_indykite_auditsink_v1beta1_ingest_proto_rawDescData = protoimpl.X.CompressGZIP(file_indykite_auditsink_v1beta1_ingest_proto_rawDescData)
	})
	return file_indykite_auditsink_v1beta1_ingest_proto_rawDescData
}

var file_indykite_auditsink_v1beta1_ingest_proto_msgTypes = make([]protoimpl.MessageInfo, 10)
var file_indykite_auditsink_v1beta1_ingest_proto_goTypes = []interface{}{
	(*UpsertData)(nil),                       // 0: indykite.auditsink.v1beta1.UpsertData
	(*DeleteData)(nil),                       // 1: indykite.auditsink.v1beta1.DeleteData
	(*DigitalTwin)(nil),                      // 2: indykite.auditsink.v1beta1.DigitalTwin
	(*Resource)(nil),                         // 3: indykite.auditsink.v1beta1.Resource
	(*Node)(nil),                             // 4: indykite.auditsink.v1beta1.Node
	(*Relation)(nil),                         // 5: indykite.auditsink.v1beta1.Relation
	(*NodeMatch)(nil),                        // 6: indykite.auditsink.v1beta1.NodeMatch
	(*RelationMatch)(nil),                    // 7: indykite.auditsink.v1beta1.RelationMatch
	(*DeleteData_NodePropertyMatch)(nil),     // 8: indykite.auditsink.v1beta1.DeleteData.NodePropertyMatch
	(*DeleteData_RelationPropertyMatch)(nil), // 9: indykite.auditsink.v1beta1.DeleteData.RelationPropertyMatch
}
var file_indykite_auditsink_v1beta1_ingest_proto_depIdxs = []int32{
	4,  // 0: indykite.auditsink.v1beta1.UpsertData.node:type_name -> indykite.auditsink.v1beta1.Node
	5,  // 1: indykite.auditsink.v1beta1.UpsertData.relation:type_name -> indykite.auditsink.v1beta1.Relation
	6,  // 2: indykite.auditsink.v1beta1.DeleteData.node:type_name -> indykite.auditsink.v1beta1.NodeMatch
	7,  // 3: indykite.auditsink.v1beta1.DeleteData.relation:type_name -> indykite.auditsink.v1beta1.RelationMatch
	8,  // 4: indykite.auditsink.v1beta1.DeleteData.node_property:type_name -> indykite.auditsink.v1beta1.DeleteData.NodePropertyMatch
	9,  // 5: indykite.auditsink.v1beta1.DeleteData.relation_property:type_name -> indykite.auditsink.v1beta1.DeleteData.RelationPropertyMatch
	2,  // 6: indykite.auditsink.v1beta1.Node.digital_twin:type_name -> indykite.auditsink.v1beta1.DigitalTwin
	3,  // 7: indykite.auditsink.v1beta1.Node.resource:type_name -> indykite.auditsink.v1beta1.Resource
	7,  // 8: indykite.auditsink.v1beta1.Relation.match:type_name -> indykite.auditsink.v1beta1.RelationMatch
	6,  // 9: indykite.auditsink.v1beta1.RelationMatch.source_match:type_name -> indykite.auditsink.v1beta1.NodeMatch
	6,  // 10: indykite.auditsink.v1beta1.RelationMatch.target_match:type_name -> indykite.auditsink.v1beta1.NodeMatch
	6,  // 11: indykite.auditsink.v1beta1.DeleteData.NodePropertyMatch.match:type_name -> indykite.auditsink.v1beta1.NodeMatch
	7,  // 12: indykite.auditsink.v1beta1.DeleteData.RelationPropertyMatch.match:type_name -> indykite.auditsink.v1beta1.RelationMatch
	13, // [13:13] is the sub-list for method output_type
	13, // [13:13] is the sub-list for method input_type
	13, // [13:13] is the sub-list for extension type_name
	13, // [13:13] is the sub-list for extension extendee
	0,  // [0:13] is the sub-list for field type_name
}

func init() { file_indykite_auditsink_v1beta1_ingest_proto_init() }
func file_indykite_auditsink_v1beta1_ingest_proto_init() {
	if File_indykite_auditsink_v1beta1_ingest_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpsertData); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteData); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DigitalTwin); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Resource); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Node); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Relation); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NodeMatch); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RelationMatch); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteData_NodePropertyMatch); i {
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
		file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteData_RelationPropertyMatch); i {
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
	file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*UpsertData_Node)(nil),
		(*UpsertData_Relation)(nil),
	}
	file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*DeleteData_Node)(nil),
		(*DeleteData_Relation)(nil),
		(*DeleteData_NodeProperty)(nil),
		(*DeleteData_RelationProperty)(nil),
	}
	file_indykite_auditsink_v1beta1_ingest_proto_msgTypes[4].OneofWrappers = []interface{}{
		(*Node_DigitalTwin)(nil),
		(*Node_Resource)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_indykite_auditsink_v1beta1_ingest_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   10,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_indykite_auditsink_v1beta1_ingest_proto_goTypes,
		DependencyIndexes: file_indykite_auditsink_v1beta1_ingest_proto_depIdxs,
		MessageInfos:      file_indykite_auditsink_v1beta1_ingest_proto_msgTypes,
	}.Build()
	File_indykite_auditsink_v1beta1_ingest_proto = out.File
	file_indykite_auditsink_v1beta1_ingest_proto_rawDesc = nil
	file_indykite_auditsink_v1beta1_ingest_proto_goTypes = nil
	file_indykite_auditsink_v1beta1_ingest_proto_depIdxs = nil
}
