// Copyright (c) 2022 IndyKite
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
// 	protoc-gen-go v1.34.2
// 	protoc        (unknown)
// source: indykite/identity/v1beta2/consent.proto

package identityv1beta2

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ConsentReceipt struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PiiPrincipalId string                 `protobuf:"bytes,1,opt,name=pii_principal_id,json=piiPrincipalId,proto3" json:"pii_principal_id,omitempty"`
	PiiProcessor   *PiiProcessor          `protobuf:"bytes,2,opt,name=pii_processor,json=piiProcessor,proto3" json:"pii_processor,omitempty"`
	Items          []*ConsentReceipt_Item `protobuf:"bytes,3,rep,name=items,proto3" json:"items,omitempty"`
}

func (x *ConsentReceipt) Reset() {
	*x = ConsentReceipt{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta2_consent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsentReceipt) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsentReceipt) ProtoMessage() {}

func (x *ConsentReceipt) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta2_consent_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsentReceipt.ProtoReflect.Descriptor instead.
func (*ConsentReceipt) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta2_consent_proto_rawDescGZIP(), []int{0}
}

func (x *ConsentReceipt) GetPiiPrincipalId() string {
	if x != nil {
		return x.PiiPrincipalId
	}
	return ""
}

func (x *ConsentReceipt) GetPiiProcessor() *PiiProcessor {
	if x != nil {
		return x.PiiProcessor
	}
	return nil
}

func (x *ConsentReceipt) GetItems() []*ConsentReceipt_Item {
	if x != nil {
		return x.Items
	}
	return nil
}

type PiiController struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PiiControllerId string `protobuf:"bytes,1,opt,name=pii_controller_id,json=piiControllerId,proto3" json:"pii_controller_id,omitempty"`
	DisplayName     string `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
}

func (x *PiiController) Reset() {
	*x = PiiController{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta2_consent_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PiiController) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PiiController) ProtoMessage() {}

func (x *PiiController) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta2_consent_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PiiController.ProtoReflect.Descriptor instead.
func (*PiiController) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta2_consent_proto_rawDescGZIP(), []int{1}
}

func (x *PiiController) GetPiiControllerId() string {
	if x != nil {
		return x.PiiControllerId
	}
	return ""
}

func (x *PiiController) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

type PiiProcessor struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PiiProcessorId          string   `protobuf:"bytes,1,opt,name=pii_processor_id,json=piiProcessorId,proto3" json:"pii_processor_id,omitempty"`
	DisplayName             string   `protobuf:"bytes,2,opt,name=display_name,json=displayName,proto3" json:"display_name,omitempty"`
	Description             string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Owner                   string   `protobuf:"bytes,4,opt,name=owner,proto3" json:"owner,omitempty"`
	PolicyUri               string   `protobuf:"bytes,5,opt,name=policy_uri,json=policyUri,proto3" json:"policy_uri,omitempty"`
	TermsOfServiceUri       string   `protobuf:"bytes,6,opt,name=terms_of_service_uri,json=termsOfServiceUri,proto3" json:"terms_of_service_uri,omitempty"`
	ClientUri               string   `protobuf:"bytes,7,opt,name=client_uri,json=clientUri,proto3" json:"client_uri,omitempty"`
	LogoUri                 string   `protobuf:"bytes,8,opt,name=logo_uri,json=logoUri,proto3" json:"logo_uri,omitempty"`
	UserSupportEmailAddress string   `protobuf:"bytes,9,opt,name=user_support_email_address,json=userSupportEmailAddress,proto3" json:"user_support_email_address,omitempty"`
	AdditionalContacts      []string `protobuf:"bytes,10,rep,name=additional_contacts,json=additionalContacts,proto3" json:"additional_contacts,omitempty"`
}

func (x *PiiProcessor) Reset() {
	*x = PiiProcessor{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta2_consent_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PiiProcessor) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PiiProcessor) ProtoMessage() {}

func (x *PiiProcessor) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta2_consent_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PiiProcessor.ProtoReflect.Descriptor instead.
func (*PiiProcessor) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta2_consent_proto_rawDescGZIP(), []int{2}
}

func (x *PiiProcessor) GetPiiProcessorId() string {
	if x != nil {
		return x.PiiProcessorId
	}
	return ""
}

func (x *PiiProcessor) GetDisplayName() string {
	if x != nil {
		return x.DisplayName
	}
	return ""
}

func (x *PiiProcessor) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *PiiProcessor) GetOwner() string {
	if x != nil {
		return x.Owner
	}
	return ""
}

func (x *PiiProcessor) GetPolicyUri() string {
	if x != nil {
		return x.PolicyUri
	}
	return ""
}

func (x *PiiProcessor) GetTermsOfServiceUri() string {
	if x != nil {
		return x.TermsOfServiceUri
	}
	return ""
}

func (x *PiiProcessor) GetClientUri() string {
	if x != nil {
		return x.ClientUri
	}
	return ""
}

func (x *PiiProcessor) GetLogoUri() string {
	if x != nil {
		return x.LogoUri
	}
	return ""
}

func (x *PiiProcessor) GetUserSupportEmailAddress() string {
	if x != nil {
		return x.UserSupportEmailAddress
	}
	return ""
}

func (x *PiiProcessor) GetAdditionalContacts() []string {
	if x != nil {
		return x.AdditionalContacts
	}
	return nil
}

type ConsentReceipt_Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ConsentId       string                 `protobuf:"bytes,1,opt,name=consent_id,json=consentId,proto3" json:"consent_id,omitempty"`
	PiiController   *PiiController         `protobuf:"bytes,2,opt,name=pii_controller,json=piiController,proto3" json:"pii_controller,omitempty"`
	ConsentedAtTime *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=consented_at_time,json=consentedAtTime,proto3" json:"consented_at_time,omitempty"`
	Properties      []string               `protobuf:"bytes,4,rep,name=properties,proto3" json:"properties,omitempty"`
}

func (x *ConsentReceipt_Item) Reset() {
	*x = ConsentReceipt_Item{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_identity_v1beta2_consent_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ConsentReceipt_Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ConsentReceipt_Item) ProtoMessage() {}

func (x *ConsentReceipt_Item) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_identity_v1beta2_consent_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ConsentReceipt_Item.ProtoReflect.Descriptor instead.
func (*ConsentReceipt_Item) Descriptor() ([]byte, []int) {
	return file_indykite_identity_v1beta2_consent_proto_rawDescGZIP(), []int{0, 0}
}

func (x *ConsentReceipt_Item) GetConsentId() string {
	if x != nil {
		return x.ConsentId
	}
	return ""
}

func (x *ConsentReceipt_Item) GetPiiController() *PiiController {
	if x != nil {
		return x.PiiController
	}
	return nil
}

func (x *ConsentReceipt_Item) GetConsentedAtTime() *timestamppb.Timestamp {
	if x != nil {
		return x.ConsentedAtTime
	}
	return nil
}

func (x *ConsentReceipt_Item) GetProperties() []string {
	if x != nil {
		return x.Properties
	}
	return nil
}

var File_indykite_identity_v1beta2_consent_proto protoreflect.FileDescriptor

var file_indykite_identity_v1beta2_consent_proto_rawDesc = []byte{
	0x0a, 0x27, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x2f, 0x63, 0x6f, 0x6e, 0x73,
	0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x19, 0x69, 0x6e, 0x64, 0x79, 0x6b,
	0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x32, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xaf, 0x03, 0x0a, 0x0e, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e,
	0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x12, 0x28, 0x0a, 0x10, 0x70, 0x69, 0x69, 0x5f,
	0x70, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x70, 0x69, 0x69, 0x50, 0x72, 0x69, 0x6e, 0x63, 0x69, 0x70, 0x61, 0x6c,
	0x49, 0x64, 0x12, 0x4c, 0x0a, 0x0d, 0x70, 0x69, 0x69, 0x5f, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x27, 0x2e, 0x69, 0x6e, 0x64, 0x79,
	0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x32, 0x2e, 0x50, 0x69, 0x69, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73,
	0x6f, 0x72, 0x52, 0x0c, 0x70, 0x69, 0x69, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x12, 0x44, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32,
	0x2e, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x2e, 0x43, 0x6f, 0x6e, 0x73,
	0x65, 0x6e, 0x74, 0x52, 0x65, 0x63, 0x65, 0x69, 0x70, 0x74, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52,
	0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x1a, 0xde, 0x01, 0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12,
	0x1d, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x4f,
	0x0a, 0x0e, 0x70, 0x69, 0x69, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74,
	0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74,
	0x61, 0x32, 0x2e, 0x50, 0x69, 0x69, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72,
	0x52, 0x0d, 0x70, 0x69, 0x69, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x12,
	0x46, 0x0a, 0x11, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74, 0x65, 0x64, 0x5f, 0x61, 0x74, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0f, 0x63, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74, 0x65,
	0x64, 0x41, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1e, 0x0a, 0x0a, 0x70, 0x72, 0x6f, 0x70, 0x65,
	0x72, 0x74, 0x69, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0a, 0x70, 0x72, 0x6f,
	0x70, 0x65, 0x72, 0x74, 0x69, 0x65, 0x73, 0x22, 0x5e, 0x0a, 0x0d, 0x50, 0x69, 0x69, 0x43, 0x6f,
	0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x12, 0x2a, 0x0a, 0x11, 0x70, 0x69, 0x69, 0x5f,
	0x63, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c, 0x65, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0f, 0x70, 0x69, 0x69, 0x43, 0x6f, 0x6e, 0x74, 0x72, 0x6f, 0x6c, 0x6c,
	0x65, 0x72, 0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f,
	0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70,
	0x6c, 0x61, 0x79, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x8b, 0x03, 0x0a, 0x0c, 0x50, 0x69, 0x69, 0x50,
	0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x12, 0x28, 0x0a, 0x10, 0x70, 0x69, 0x69, 0x5f,
	0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x70, 0x69, 0x69, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x6f, 0x72,
	0x49, 0x64, 0x12, 0x21, 0x0a, 0x0c, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61, 0x79, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x69, 0x73, 0x70, 0x6c, 0x61,
	0x79, 0x4e, 0x61, 0x6d, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63,
	0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72,
	0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x77, 0x6e, 0x65, 0x72, 0x12, 0x1d, 0x0a,
	0x0a, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x70, 0x6f, 0x6c, 0x69, 0x63, 0x79, 0x55, 0x72, 0x69, 0x12, 0x2f, 0x0a, 0x14,
	0x74, 0x65, 0x72, 0x6d, 0x73, 0x5f, 0x6f, 0x66, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x5f, 0x75, 0x72, 0x69, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x74, 0x65, 0x72, 0x6d,
	0x73, 0x4f, 0x66, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x55, 0x72, 0x69, 0x12, 0x1d, 0x0a,
	0x0a, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x63, 0x6c, 0x69, 0x65, 0x6e, 0x74, 0x55, 0x72, 0x69, 0x12, 0x19, 0x0a, 0x08,
	0x6c, 0x6f, 0x67, 0x6f, 0x5f, 0x75, 0x72, 0x69, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x6c, 0x6f, 0x67, 0x6f, 0x55, 0x72, 0x69, 0x12, 0x3b, 0x0a, 0x1a, 0x75, 0x73, 0x65, 0x72, 0x5f,
	0x73, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x5f, 0x65, 0x6d, 0x61, 0x69, 0x6c, 0x5f, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x17, 0x75, 0x73, 0x65,
	0x72, 0x53, 0x75, 0x70, 0x70, 0x6f, 0x72, 0x74, 0x45, 0x6d, 0x61, 0x69, 0x6c, 0x41, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x12, 0x2f, 0x0a, 0x13, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e,
	0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x6e, 0x74, 0x61, 0x63, 0x74, 0x73, 0x18, 0x0a, 0x20, 0x03, 0x28,
	0x09, 0x52, 0x12, 0x61, 0x64, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x61, 0x6c, 0x43, 0x6f, 0x6e,
	0x74, 0x61, 0x63, 0x74, 0x73, 0x42, 0x86, 0x02, 0x0a, 0x1d, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e,
	0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x2e,
	0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x42, 0x0c, 0x43, 0x6f, 0x6e, 0x73, 0x65, 0x6e, 0x74,
	0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x51, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x69, 0x6e, 0x64,
	0x79, 0x6b, 0x69, 0x74, 0x65, 0x2d, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x2f, 0x67, 0x65, 0x6e,
	0x2f, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x69, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x74, 0x79, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x3b, 0x69, 0x64, 0x65, 0x6e, 0x74,
	0x69, 0x74, 0x79, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0xa2, 0x02, 0x03, 0x49, 0x49, 0x58,
	0xaa, 0x02, 0x19, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x49, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x74, 0x79, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0xca, 0x02, 0x19, 0x49,
	0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x5c, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79,
	0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0xe2, 0x02, 0x25, 0x49, 0x6e, 0x64, 0x79, 0x6b,
	0x69, 0x74, 0x65, 0x5c, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x74, 0x79, 0x5c, 0x56, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x32, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61,
	0xea, 0x02, 0x1b, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x3a, 0x3a, 0x49, 0x64, 0x65,
	0x6e, 0x74, 0x69, 0x74, 0x79, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x32, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_indykite_identity_v1beta2_consent_proto_rawDescOnce sync.Once
	file_indykite_identity_v1beta2_consent_proto_rawDescData = file_indykite_identity_v1beta2_consent_proto_rawDesc
)

func file_indykite_identity_v1beta2_consent_proto_rawDescGZIP() []byte {
	file_indykite_identity_v1beta2_consent_proto_rawDescOnce.Do(func() {
		file_indykite_identity_v1beta2_consent_proto_rawDescData = protoimpl.X.CompressGZIP(file_indykite_identity_v1beta2_consent_proto_rawDescData)
	})
	return file_indykite_identity_v1beta2_consent_proto_rawDescData
}

var file_indykite_identity_v1beta2_consent_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_indykite_identity_v1beta2_consent_proto_goTypes = []any{
	(*ConsentReceipt)(nil),        // 0: indykite.identity.v1beta2.ConsentReceipt
	(*PiiController)(nil),         // 1: indykite.identity.v1beta2.PiiController
	(*PiiProcessor)(nil),          // 2: indykite.identity.v1beta2.PiiProcessor
	(*ConsentReceipt_Item)(nil),   // 3: indykite.identity.v1beta2.ConsentReceipt.Item
	(*timestamppb.Timestamp)(nil), // 4: google.protobuf.Timestamp
}
var file_indykite_identity_v1beta2_consent_proto_depIdxs = []int32{
	2, // 0: indykite.identity.v1beta2.ConsentReceipt.pii_processor:type_name -> indykite.identity.v1beta2.PiiProcessor
	3, // 1: indykite.identity.v1beta2.ConsentReceipt.items:type_name -> indykite.identity.v1beta2.ConsentReceipt.Item
	1, // 2: indykite.identity.v1beta2.ConsentReceipt.Item.pii_controller:type_name -> indykite.identity.v1beta2.PiiController
	4, // 3: indykite.identity.v1beta2.ConsentReceipt.Item.consented_at_time:type_name -> google.protobuf.Timestamp
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_indykite_identity_v1beta2_consent_proto_init() }
func file_indykite_identity_v1beta2_consent_proto_init() {
	if File_indykite_identity_v1beta2_consent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_indykite_identity_v1beta2_consent_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ConsentReceipt); i {
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
		file_indykite_identity_v1beta2_consent_proto_msgTypes[1].Exporter = func(v any, i int) any {
			switch v := v.(*PiiController); i {
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
		file_indykite_identity_v1beta2_consent_proto_msgTypes[2].Exporter = func(v any, i int) any {
			switch v := v.(*PiiProcessor); i {
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
		file_indykite_identity_v1beta2_consent_proto_msgTypes[3].Exporter = func(v any, i int) any {
			switch v := v.(*ConsentReceipt_Item); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_indykite_identity_v1beta2_consent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_indykite_identity_v1beta2_consent_proto_goTypes,
		DependencyIndexes: file_indykite_identity_v1beta2_consent_proto_depIdxs,
		MessageInfos:      file_indykite_identity_v1beta2_consent_proto_msgTypes,
	}.Build()
	File_indykite_identity_v1beta2_consent_proto = out.File
	file_indykite_identity_v1beta2_consent_proto_rawDesc = nil
	file_indykite_identity_v1beta2_consent_proto_goTypes = nil
	file_indykite_identity_v1beta2_consent_proto_depIdxs = nil
}
