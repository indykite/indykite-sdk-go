// Copyright (c) 2021 IndyKite
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
// 	protoc-gen-go v1.28.1
// 	protoc        (unknown)
// source: indykite/identity/v1beta1/document.proto

package identityv1beta1

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	v1beta1 "github.com/indykite/jarvis-sdk-go/gen/indykite/objects/v1beta1"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Document is an IndyKite document object.
type Document struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Name of the document resource.
	// Format: `databases/{application_id}/documents/{document_path}`.
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	// Fields are the key/value pairs of the document.
	//
	// The map keys represent field names.
	Fields map[string]*v1beta1.Value `protobuf:"bytes,2,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	// CreateTime when the document was created.
	CreateTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=create_time,json=createTime,proto3" json:"create_time,omitempty"`
	// UpdateTime when the document was last changed.
	UpdateTime *timestamppb.Timestamp `protobuf:"bytes,4,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
}

func (x *Document) Reset() {
	*x = Document{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Document) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Document) ProtoMessage() {}

func (x *Document) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Document.ProtoReflect.Descriptor instead.
func (*Document) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta1_document_proto_rawDescGZIP(), []int{0}
}

func (x *Document) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Document) GetFields() map[string]*v1beta1.Value {
	if x != nil {
		return x.Fields
	}
	return nil
}

func (x *Document) GetCreateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.CreateTime
	}
	return nil
}

func (x *Document) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

// DocumentMask used to restrict a get or update operation on a document to a subset of its fields.
type DocumentMask struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// FieldPaths is a list of fields in the mask.
	FieldPaths []string `protobuf:"bytes,1,rep,name=field_paths,json=fieldPaths,proto3" json:"field_paths,omitempty"`
}

func (x *DocumentMask) Reset() {
	*x = DocumentMask{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocumentMask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentMask) ProtoMessage() {}

func (x *DocumentMask) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentMask.ProtoReflect.Descriptor instead.
func (*DocumentMask) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta1_document_proto_rawDescGZIP(), []int{1}
}

func (x *DocumentMask) GetFieldPaths() []string {
	if x != nil {
		return x.FieldPaths
	}
	return nil
}

// Precondition used for conditional operations on a Document.
type Precondition struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to ConditionType:
	//	*Precondition_Exists
	//	*Precondition_UpdateTime
	ConditionType isPrecondition_ConditionType `protobuf_oneof:"condition_type"`
}

func (x *Precondition) Reset() {
	*x = Precondition{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Precondition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Precondition) ProtoMessage() {}

func (x *Precondition) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Precondition.ProtoReflect.Descriptor instead.
func (*Precondition) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta1_document_proto_rawDescGZIP(), []int{2}
}

func (m *Precondition) GetConditionType() isPrecondition_ConditionType {
	if m != nil {
		return m.ConditionType
	}
	return nil
}

func (x *Precondition) GetExists() bool {
	if x, ok := x.GetConditionType().(*Precondition_Exists); ok {
		return x.Exists
	}
	return false
}

func (x *Precondition) GetUpdateTime() *timestamppb.Timestamp {
	if x, ok := x.GetConditionType().(*Precondition_UpdateTime); ok {
		return x.UpdateTime
	}
	return nil
}

type isPrecondition_ConditionType interface {
	isPrecondition_ConditionType()
}

type Precondition_Exists struct {
	// Exists set to `true` when the target document must exist else set to `false`.
	Exists bool `protobuf:"varint,1,opt,name=exists,proto3,oneof"`
}

type Precondition_UpdateTime struct {
	// UpdateTime when set, the target document must exist and have been last updated at that time.
	UpdateTime *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=update_time,json=updateTime,proto3,oneof"`
}

func (*Precondition_Exists) isPrecondition_ConditionType() {}

func (*Precondition_UpdateTime) isPrecondition_ConditionType() {}

// Write is a single operation on a document.
type Write struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Operation to execute.
	//
	// Types that are assignable to Operation:
	//	*Write_Update
	//	*Write_Delete
	//	*Write_Transform
	Operation isWrite_Operation `protobuf_oneof:"operation"`
	// UpdateMask is the fields to update in this write.
	//
	// This field can be set only when the operation is `update`.
	// If the mask is not set for an `update` and the document exists, any
	// existing data will be overwritten.
	// If the mask is set and the document on the server has fields not covered by
	// the mask, they are left unchanged.
	// Fields referenced in the mask, but not present in the input document, are
	// deleted from the document on the server.
	UpdateMask *DocumentMask `protobuf:"bytes,3,opt,name=update_mask,json=updateMask,proto3" json:"update_mask,omitempty"`
	// UpdateTransforms represents the transforms to perform after update.
	//
	// This field can be set only when the operation is `update`.
	UpdateTransforms []*DocumentTransform_FieldTransform `protobuf:"bytes,7,rep,name=update_transforms,json=updateTransforms,proto3" json:"update_transforms,omitempty"`
	// CurrentDocument is an optional precondition on the document.
	//
	// The write will fail if this is set and not met by the target document.
	CurrentDocument *Precondition `protobuf:"bytes,4,opt,name=current_document,json=currentDocument,proto3" json:"current_document,omitempty"`
}

func (x *Write) Reset() {
	*x = Write{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Write) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Write) ProtoMessage() {}

func (x *Write) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Write.ProtoReflect.Descriptor instead.
func (*Write) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta1_document_proto_rawDescGZIP(), []int{3}
}

func (m *Write) GetOperation() isWrite_Operation {
	if m != nil {
		return m.Operation
	}
	return nil
}

func (x *Write) GetUpdate() *Document {
	if x, ok := x.GetOperation().(*Write_Update); ok {
		return x.Update
	}
	return nil
}

func (x *Write) GetDelete() string {
	if x, ok := x.GetOperation().(*Write_Delete); ok {
		return x.Delete
	}
	return ""
}

func (x *Write) GetTransform() *DocumentTransform {
	if x, ok := x.GetOperation().(*Write_Transform); ok {
		return x.Transform
	}
	return nil
}

func (x *Write) GetUpdateMask() *DocumentMask {
	if x != nil {
		return x.UpdateMask
	}
	return nil
}

func (x *Write) GetUpdateTransforms() []*DocumentTransform_FieldTransform {
	if x != nil {
		return x.UpdateTransforms
	}
	return nil
}

func (x *Write) GetCurrentDocument() *Precondition {
	if x != nil {
		return x.CurrentDocument
	}
	return nil
}

type isWrite_Operation interface {
	isWrite_Operation()
}

type Write_Update struct {
	// Document to write.
	Update *Document `protobuf:"bytes,1,opt,name=update,proto3,oneof"`
}

type Write_Delete struct {
	// Delete is a document name to delete.
	//
	// Format: `databases/{application_id}/documents/{document_path}`.
	Delete string `protobuf:"bytes,2,opt,name=delete,proto3,oneof"`
}

type Write_Transform struct {
	// Transform represent a transformation to a document.
	Transform *DocumentTransform `protobuf:"bytes,6,opt,name=transform,proto3,oneof"`
}

func (*Write_Update) isWrite_Operation() {}

func (*Write_Delete) isWrite_Operation() {}

func (*Write_Transform) isWrite_Operation() {}

// WriteResult represents the result of applying a write.
type WriteResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// UpdateTime is the last update time of the document after applying the write.
	UpdateTime *timestamppb.Timestamp `protobuf:"bytes,1,opt,name=update_time,json=updateTime,proto3" json:"update_time,omitempty"`
}

func (x *WriteResult) Reset() {
	*x = WriteResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *WriteResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WriteResult) ProtoMessage() {}

func (x *WriteResult) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use WriteResult.ProtoReflect.Descriptor instead.
func (*WriteResult) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta1_document_proto_rawDescGZIP(), []int{4}
}

func (x *WriteResult) GetUpdateTime() *timestamppb.Timestamp {
	if x != nil {
		return x.UpdateTime
	}
	return nil
}

// DocumentTransform represents the transformation of a document.
type DocumentTransform struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Document is the name of the document to transform.
	Document string `protobuf:"bytes,1,opt,name=document,proto3" json:"document,omitempty"`
	// FieldTransforms is the list of transformations to apply to the fields of the document.
	FieldTransforms []*DocumentTransform_FieldTransform `protobuf:"bytes,2,rep,name=field_transforms,json=fieldTransforms,proto3" json:"field_transforms,omitempty"`
}

func (x *DocumentTransform) Reset() {
	*x = DocumentTransform{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocumentTransform) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentTransform) ProtoMessage() {}

func (x *DocumentTransform) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentTransform.ProtoReflect.Descriptor instead.
func (*DocumentTransform) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta1_document_proto_rawDescGZIP(), []int{5}
}

func (x *DocumentTransform) GetDocument() string {
	if x != nil {
		return x.Document
	}
	return ""
}

func (x *DocumentTransform) GetFieldTransforms() []*DocumentTransform_FieldTransform {
	if x != nil {
		return x.FieldTransforms
	}
	return nil
}

// FieldTransform represents the transformation of a field of the document.
type DocumentTransform_FieldTransform struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// FieldPath is the path of the field.
	FieldPath string `protobuf:"bytes,1,opt,name=field_path,json=fieldPath,proto3" json:"field_path,omitempty"`
	// The transformation to apply on the field.
	//
	// Types that are assignable to TransformType:
	//	*DocumentTransform_FieldTransform_AppendMissingElements
	//	*DocumentTransform_FieldTransform_RemoveAllFromArray
	TransformType isDocumentTransform_FieldTransform_TransformType `protobuf_oneof:"transform_type"`
}

func (x *DocumentTransform_FieldTransform) Reset() {
	*x = DocumentTransform_FieldTransform{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DocumentTransform_FieldTransform) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DocumentTransform_FieldTransform) ProtoMessage() {}

func (x *DocumentTransform_FieldTransform) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta1_document_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DocumentTransform_FieldTransform.ProtoReflect.Descriptor instead.
func (*DocumentTransform_FieldTransform) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta1_document_proto_rawDescGZIP(), []int{5, 0}
}

func (x *DocumentTransform_FieldTransform) GetFieldPath() string {
	if x != nil {
		return x.FieldPath
	}
	return ""
}

func (m *DocumentTransform_FieldTransform) GetTransformType() isDocumentTransform_FieldTransform_TransformType {
	if m != nil {
		return m.TransformType
	}
	return nil
}

func (x *DocumentTransform_FieldTransform) GetAppendMissingElements() *v1beta1.ArrayValue {
	if x, ok := x.GetTransformType().(*DocumentTransform_FieldTransform_AppendMissingElements); ok {
		return x.AppendMissingElements
	}
	return nil
}

func (x *DocumentTransform_FieldTransform) GetRemoveAllFromArray() *v1beta1.ArrayValue {
	if x, ok := x.GetTransformType().(*DocumentTransform_FieldTransform_RemoveAllFromArray); ok {
		return x.RemoveAllFromArray
	}
	return nil
}

type isDocumentTransform_FieldTransform_TransformType interface {
	isDocumentTransform_FieldTransform_TransformType()
}

type DocumentTransform_FieldTransform_AppendMissingElements struct {
	// AppendMissingElements transforms the field by appending the given elements in order
	// if they are not already present in the current field value.
	// If the field is not an array, or if the field does not yet exist, it is
	// first set to the empty array.
	//
	// Equivalent numbers of different types (e.g. 3L and 3.0) are
	// considered equal when checking if a value is missing.
	// NaN is equal to NaN, and Null is equal to Null.
	// If the input contains multiple equivalent values, only the first will
	// be considered.
	//
	// The corresponding transform_result will be the null value.
	AppendMissingElements *v1beta1.ArrayValue `protobuf:"bytes,6,opt,name=append_missing_elements,json=appendMissingElements,proto3,oneof"`
}

type DocumentTransform_FieldTransform_RemoveAllFromArray struct {
	// remove_all_from_array Removes all of the given elements from the array in the field.
	// If the field is not an array, or if the field does not yet exist, it is
	// set to the empty array.
	RemoveAllFromArray *v1beta1.ArrayValue `protobuf:"bytes,7,opt,name=remove_all_from_array,json=removeAllFromArray,proto3,oneof"`
}

func (*DocumentTransform_FieldTransform_AppendMissingElements) isDocumentTransform_FieldTransform_TransformType() {
}

func (*DocumentTransform_FieldTransform_RemoveAllFromArray) isDocumentTransform_FieldTransform_TransformType() {
}

var File_indykite_identity_v1beta1_document_proto protoreflect.FileDescriptor

var file_indykite_identity_v1beta1_document_proto_rawDesc = []byte{
	0x0a, 0x28, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x64, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x69, 0x6e, 0x64, 0x79,
	0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x25, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65,
	0x2f, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31,
	0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbd, 0x02,
	0x0a, 0x08, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x47,
	0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x2f,
	0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x12, 0x3b, 0x0a, 0x0b, 0x63, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x63, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x69, 0x6d, 0x65, 0x12, 0x3b, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74,
	0x69, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65,
	0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d,
	0x65, 0x1a, 0x5a, 0x0a, 0x0b, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79,
	0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b,
	0x65, 0x79, 0x12, 0x35, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x1f, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x6f, 0x62, 0x6a,
	0x65, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x2f, 0x0a,
	0x0c, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x61, 0x73, 0x6b, 0x12, 0x1f, 0x0a,
	0x0b, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x73, 0x18, 0x01, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x0a, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x50, 0x61, 0x74, 0x68, 0x73, 0x22, 0x79,
	0x0a, 0x0c, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x18,
	0x0a, 0x06, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00,
	0x52, 0x06, 0x65, 0x78, 0x69, 0x73, 0x74, 0x73, 0x12, 0x3d, 0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e,
	0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x48, 0x00, 0x52, 0x0a, 0x75, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x42, 0x10, 0x0a, 0x0e, 0x63, 0x6f, 0x6e, 0x64, 0x69,
	0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x22, 0xc3, 0x03, 0x0a, 0x05, 0x57, 0x72,
	0x69, 0x74, 0x65, 0x12, 0x3d, 0x0a, 0x06, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x23, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69,
	0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e,
	0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x06, 0x75, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x12, 0x18, 0x0a, 0x06, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x00, 0x52, 0x06, 0x64, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x4c, 0x0a, 0x09,
	0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2c, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x48, 0x00, 0x52,
	0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x48, 0x0a, 0x0b, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x5f, 0x6d, 0x61, 0x73, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x27, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x4d, 0x61, 0x73, 0x6b, 0x52, 0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x61, 0x73, 0x6b, 0x12, 0x68, 0x0a, 0x11, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x18, 0x07, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x3b, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x75,
	0x6d, 0x65, 0x6e, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x46, 0x69,
	0x65, 0x6c, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x52, 0x10, 0x75, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x73, 0x12, 0x52,
	0x0a, 0x10, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b,
	0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x50, 0x72, 0x65, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x52, 0x0f, 0x63, 0x75, 0x72, 0x72, 0x65, 0x6e, 0x74, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65,
	0x6e, 0x74, 0x42, 0x0b, 0x0a, 0x09, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x22,
	0x4a, 0x0a, 0x0b, 0x57, 0x72, 0x69, 0x74, 0x65, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x3b,
	0x0a, 0x0b, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52,
	0x0a, 0x75, 0x70, 0x64, 0x61, 0x74, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x22, 0x96, 0x03, 0x0a, 0x11,
	0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72,
	0x6d, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x66, 0x0a,
	0x10, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d,
	0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3b, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69,
	0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74, 0x54, 0x72, 0x61, 0x6e,
	0x73, 0x66, 0x6f, 0x72, 0x6d, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x6f, 0x72, 0x6d, 0x52, 0x0f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x54, 0x72, 0x61, 0x6e, 0x73,
	0x66, 0x6f, 0x72, 0x6d, 0x73, 0x1a, 0xfc, 0x01, 0x0a, 0x0e, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x54,
	0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x12, 0x1d, 0x0a, 0x0a, 0x66, 0x69, 0x65, 0x6c,
	0x64, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x66, 0x69,
	0x65, 0x6c, 0x64, 0x50, 0x61, 0x74, 0x68, 0x12, 0x5e, 0x0a, 0x17, 0x61, 0x70, 0x70, 0x65, 0x6e,
	0x64, 0x5f, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x5f, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x6e,
	0x74, 0x73, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b,
	0x69, 0x74, 0x65, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65,
	0x74, 0x61, 0x31, 0x2e, 0x41, 0x72, 0x72, 0x61, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00,
	0x52, 0x15, 0x61, 0x70, 0x70, 0x65, 0x6e, 0x64, 0x4d, 0x69, 0x73, 0x73, 0x69, 0x6e, 0x67, 0x45,
	0x6c, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x73, 0x12, 0x59, 0x0a, 0x15, 0x72, 0x65, 0x6d, 0x6f, 0x76,
	0x65, 0x5f, 0x61, 0x6c, 0x6c, 0x5f, 0x66, 0x72, 0x6f, 0x6d, 0x5f, 0x61, 0x72, 0x72, 0x61, 0x79,
	0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74,
	0x65, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x41, 0x72, 0x72, 0x61, 0x79, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x12,
	0x72, 0x65, 0x6d, 0x6f, 0x76, 0x65, 0x41, 0x6c, 0x6c, 0x46, 0x72, 0x6f, 0x6d, 0x41, 0x72, 0x72,
	0x61, 0x79, 0x42, 0x10, 0x0a, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x6f, 0x72, 0x6d, 0x5f,
	0x74, 0x79, 0x70, 0x65, 0x42, 0x85, 0x02, 0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e, 0x64,
	0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x42, 0x0d, 0x44, 0x6f, 0x63, 0x75, 0x6d, 0x65, 0x6e, 0x74,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x4f, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x6a, 0x61, 0x72,
	0x76, 0x69, 0x73, 0x2d, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x69,
	0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74,
	0x79, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xa2, 0x02, 0x03, 0x49, 0x49, 0x58, 0xaa, 0x02,
	0x19, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xca, 0x02, 0x19, 0x49, 0x6e, 0x64,
	0x79, 0x6b, 0x69, 0x74, 0x65, 0x5c, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5c, 0x56,
	0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xe2, 0x02, 0x25, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74,
	0x65, 0x5c, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5c, 0x56, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02,
	0x1b, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x3a, 0x3a, 0x49, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_indykite_identity_v1beta1_document_proto_rawDescOnce sync.Once
	file_indykite_identity_v1beta1_document_proto_rawDescData = file_indykite_identity_v1beta1_document_proto_rawDesc
)

func file_indykite_identity_v1beta1_document_proto_rawDescGZIP() []byte {
	file_indykite_identity_v1beta1_document_proto_rawDescOnce.Do(func() {
		file_indykite_identity_v1beta1_document_proto_rawDescData = protoimpl.X.CompressGZIP(file_indykite_identity_v1beta1_document_proto_rawDescData)
	})
	return file_indykite_identity_v1beta1_document_proto_rawDescData
}

var file_indykite_identity_v1beta1_document_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_indykite_identity_v1beta1_document_proto_goTypes = []interface{}{
	(*Document)(nil),                         // 0: indykite.identity.v1beta1.Document
	(*DocumentMask)(nil),                     // 1: indykite.identity.v1beta1.DocumentMask
	(*Precondition)(nil),                     // 2: indykite.identity.v1beta1.Precondition
	(*Write)(nil),                            // 3: indykite.identity.v1beta1.Write
	(*WriteResult)(nil),                      // 4: indykite.identity.v1beta1.WriteResult
	(*DocumentTransform)(nil),                // 5: indykite.identity.v1beta1.DocumentTransform
	nil,                                      // 6: indykite.identity.v1beta1.Document.FieldsEntry
	(*DocumentTransform_FieldTransform)(nil), // 7: indykite.identity.v1beta1.DocumentTransform.FieldTransform
	(*timestamppb.Timestamp)(nil),            // 8: google.protobuf.Timestamp
	(*v1beta1.Value)(nil),                    // 9: indykite.objects.v1beta1.Value
	(*v1beta1.ArrayValue)(nil),               // 10: indykite.objects.v1beta1.ArrayValue
}
var file_indykite_identity_v1beta1_document_proto_depIdxs = []int32{
	6,  // 0: indykite.identity.v1beta1.Document.fields:type_name -> indykite.identity.v1beta1.Document.FieldsEntry
	8,  // 1: indykite.identity.v1beta1.Document.create_time:type_name -> google.protobuf.Timestamp
	8,  // 2: indykite.identity.v1beta1.Document.update_time:type_name -> google.protobuf.Timestamp
	8,  // 3: indykite.identity.v1beta1.Precondition.update_time:type_name -> google.protobuf.Timestamp
	0,  // 4: indykite.identity.v1beta1.Write.update:type_name -> indykite.identity.v1beta1.Document
	5,  // 5: indykite.identity.v1beta1.Write.transform:type_name -> indykite.identity.v1beta1.DocumentTransform
	1,  // 6: indykite.identity.v1beta1.Write.update_mask:type_name -> indykite.identity.v1beta1.DocumentMask
	7,  // 7: indykite.identity.v1beta1.Write.update_transforms:type_name -> indykite.identity.v1beta1.DocumentTransform.FieldTransform
	2,  // 8: indykite.identity.v1beta1.Write.current_document:type_name -> indykite.identity.v1beta1.Precondition
	8,  // 9: indykite.identity.v1beta1.WriteResult.update_time:type_name -> google.protobuf.Timestamp
	7,  // 10: indykite.identity.v1beta1.DocumentTransform.field_transforms:type_name -> indykite.identity.v1beta1.DocumentTransform.FieldTransform
	9,  // 11: indykite.identity.v1beta1.Document.FieldsEntry.value:type_name -> indykite.objects.v1beta1.Value
	10, // 12: indykite.identity.v1beta1.DocumentTransform.FieldTransform.append_missing_elements:type_name -> indykite.objects.v1beta1.ArrayValue
	10, // 13: indykite.identity.v1beta1.DocumentTransform.FieldTransform.remove_all_from_array:type_name -> indykite.objects.v1beta1.ArrayValue
	14, // [14:14] is the sub-list for method output_type
	14, // [14:14] is the sub-list for method input_type
	14, // [14:14] is the sub-list for extension type_name
	14, // [14:14] is the sub-list for extension extendee
	0,  // [0:14] is the sub-list for field type_name
}

func init() { file_indykite_identity_v1beta1_document_proto_init() }
func file_indykite_identity_v1beta1_document_proto_init() {
	if File_indykite_identity_v1beta1_document_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_indykite_identity_v1beta1_document_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Document); i {
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
		file_indykite_identity_v1beta1_document_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DocumentMask); i {
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
		file_indykite_identity_v1beta1_document_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Precondition); i {
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
		file_indykite_identity_v1beta1_document_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Write); i {
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
		file_indykite_identity_v1beta1_document_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*WriteResult); i {
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
		file_indykite_identity_v1beta1_document_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DocumentTransform); i {
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
		file_indykite_identity_v1beta1_document_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DocumentTransform_FieldTransform); i {
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
	file_indykite_identity_v1beta1_document_proto_msgTypes[2].OneofWrappers = []interface{}{
		(*Precondition_Exists)(nil),
		(*Precondition_UpdateTime)(nil),
	}
	file_indykite_identity_v1beta1_document_proto_msgTypes[3].OneofWrappers = []interface{}{
		(*Write_Update)(nil),
		(*Write_Delete)(nil),
		(*Write_Transform)(nil),
	}
	file_indykite_identity_v1beta1_document_proto_msgTypes[7].OneofWrappers = []interface{}{
		(*DocumentTransform_FieldTransform_AppendMissingElements)(nil),
		(*DocumentTransform_FieldTransform_RemoveAllFromArray)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_indykite_identity_v1beta1_document_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_indykite_identity_v1beta1_document_proto_goTypes,
		DependencyIndexes: file_indykite_identity_v1beta1_document_proto_depIdxs,
		MessageInfos:      file_indykite_identity_v1beta1_document_proto_msgTypes,
	}.Build()
	File_indykite_identity_v1beta1_document_proto = out.File
	file_indykite_identity_v1beta1_document_proto_rawDesc = nil
	file_indykite_identity_v1beta1_document_proto_goTypes = nil
	file_indykite_identity_v1beta1_document_proto_depIdxs = nil
}
