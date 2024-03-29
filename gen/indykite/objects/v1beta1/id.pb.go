// Copyright (c) 2020 IndyKite
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
// source: indykite/objects/v1beta1/id.proto

package objectsv1beta1

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

// Identifier is a universally unique identifier (UUID) a 128-bit number used to identify information in system.
type Identifier struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Id:
	//	*Identifier_IdString
	//	*Identifier_IdBytes
	Id isIdentifier_Id `protobuf_oneof:"id"`
}

func (x *Identifier) Reset() {
	*x = Identifier{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_objects_v1beta1_id_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Identifier) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Identifier) ProtoMessage() {}

func (x *Identifier) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_objects_v1beta1_id_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Identifier.ProtoReflect.Descriptor instead.
func (*Identifier) Descriptor() ([]byte, []int) {
	return file_indykite_objects_v1beta1_id_proto_rawDescGZIP(), []int{0}
}

func (m *Identifier) GetId() isIdentifier_Id {
	if m != nil {
		return m.Id
	}
	return nil
}

func (x *Identifier) GetIdString() string {
	if x, ok := x.GetId().(*Identifier_IdString); ok {
		return x.IdString
	}
	return ""
}

func (x *Identifier) GetIdBytes() []byte {
	if x, ok := x.GetId().(*Identifier_IdBytes); ok {
		return x.IdBytes
	}
	return nil
}

type isIdentifier_Id interface {
	isIdentifier_Id()
}

type Identifier_IdString struct {
	//String representation of an RFC4122 compliant UUID.
	IdString string `protobuf:"bytes,7,opt,name=id_string,json=idString,proto3,oneof"`
}

type Identifier_IdBytes struct {
	// Byte[16] array representation of an RFC4122 compliant UUID.
	IdBytes []byte `protobuf:"bytes,8,opt,name=id_bytes,json=idBytes,proto3,oneof"`
}

func (*Identifier_IdString) isIdentifier_Id() {}

func (*Identifier_IdBytes) isIdentifier_Id() {}

// ObjectReference ...
type ObjectReference struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	//UUID of the top level Customer.
	CustomerId *Identifier `protobuf:"bytes,1,opt,name=customer_id,json=customerId,proto3" json:"customer_id,omitempty"`
	//UUID of Application Space in Customer.
	AppSpaceId *Identifier `protobuf:"bytes,2,opt,name=app_space_id,json=appSpaceId,proto3" json:"app_space_id,omitempty"`
	//UUID of Application in Application Space.
	AppId *Identifier `protobuf:"bytes,3,opt,name=app_id,json=appId,proto3" json:"app_id,omitempty"`
	//UUID of Tenant in Application Space.
	TenantId *Identifier `protobuf:"bytes,4,opt,name=tenant_id,json=tenantId,proto3" json:"tenant_id,omitempty"`
	//Gives a hint about what the identifier refers to. Usually a URL to the schema of the target object.
	TypeHint string `protobuf:"bytes,6,opt,name=type_hint,json=typeHint,proto3" json:"type_hint,omitempty"`
	//UUID of Object to refer to.
	Id *Identifier `protobuf:"bytes,7,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *ObjectReference) Reset() {
	*x = ObjectReference{}
	if protoimpl.UnsafeEnabled {
		mi := &file_indykite_objects_v1beta1_id_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ObjectReference) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ObjectReference) ProtoMessage() {}

func (x *ObjectReference) ProtoReflect() protoreflect.Message {
	mi := &file_indykite_objects_v1beta1_id_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ObjectReference.ProtoReflect.Descriptor instead.
func (*ObjectReference) Descriptor() ([]byte, []int) {
	return file_indykite_objects_v1beta1_id_proto_rawDescGZIP(), []int{1}
}

func (x *ObjectReference) GetCustomerId() *Identifier {
	if x != nil {
		return x.CustomerId
	}
	return nil
}

func (x *ObjectReference) GetAppSpaceId() *Identifier {
	if x != nil {
		return x.AppSpaceId
	}
	return nil
}

func (x *ObjectReference) GetAppId() *Identifier {
	if x != nil {
		return x.AppId
	}
	return nil
}

func (x *ObjectReference) GetTenantId() *Identifier {
	if x != nil {
		return x.TenantId
	}
	return nil
}

func (x *ObjectReference) GetTypeHint() string {
	if x != nil {
		return x.TypeHint
	}
	return ""
}

func (x *ObjectReference) GetId() *Identifier {
	if x != nil {
		return x.Id
	}
	return nil
}

var File_indykite_objects_v1beta1_id_proto protoreflect.FileDescriptor

var file_indykite_objects_v1beta1_id_proto_rawDesc = []byte{
	0x0a, 0x21, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f, 0x6f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2f, 0x69, 0x64, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x18, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x6f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x22, 0x4e, 0x0a,
	0x0a, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x12, 0x1d, 0x0a, 0x09, 0x69,
	0x64, 0x5f, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x08, 0x69, 0x64, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x12, 0x1b, 0x0a, 0x08, 0x69, 0x64,
	0x5f, 0x62, 0x79, 0x74, 0x65, 0x73, 0x18, 0x08, 0x20, 0x01, 0x28, 0x0c, 0x48, 0x00, 0x52, 0x07,
	0x69, 0x64, 0x42, 0x79, 0x74, 0x65, 0x73, 0x42, 0x04, 0x0a, 0x02, 0x69, 0x64, 0x22, 0xf3, 0x02,
	0x0a, 0x0f, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x52, 0x65, 0x66, 0x65, 0x72, 0x65, 0x6e, 0x63,
	0x65, 0x12, 0x45, 0x0a, 0x0b, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74,
	0x65, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61,
	0x31, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x0a, 0x63, 0x75,
	0x73, 0x74, 0x6f, 0x6d, 0x65, 0x72, 0x49, 0x64, 0x12, 0x46, 0x0a, 0x0c, 0x61, 0x70, 0x70, 0x5f,
	0x73, 0x70, 0x61, 0x63, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24,
	0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74,
	0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69,
	0x66, 0x69, 0x65, 0x72, 0x52, 0x0a, 0x61, 0x70, 0x70, 0x53, 0x70, 0x61, 0x63, 0x65, 0x49, 0x64,
	0x12, 0x3b, 0x0a, 0x06, 0x61, 0x70, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x24, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x49, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x05, 0x61, 0x70, 0x70, 0x49, 0x64, 0x12, 0x41, 0x0a,
	0x09, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x24, 0x2e, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x2e, 0x49, 0x64, 0x65, 0x6e,
	0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52, 0x08, 0x74, 0x65, 0x6e, 0x61, 0x6e, 0x74, 0x49, 0x64,
	0x12, 0x1b, 0x0a, 0x09, 0x74, 0x79, 0x70, 0x65, 0x5f, 0x68, 0x69, 0x6e, 0x74, 0x18, 0x06, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x74, 0x79, 0x70, 0x65, 0x48, 0x69, 0x6e, 0x74, 0x12, 0x34, 0x0a,
	0x02, 0x69, 0x64, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x69, 0x6e, 0x64, 0x79,
	0x6b, 0x69, 0x74, 0x65, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x2e, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x69, 0x65, 0x72, 0x52,
	0x02, 0x69, 0x64, 0x42, 0xf8, 0x01, 0x0a, 0x1c, 0x63, 0x6f, 0x6d, 0x2e, 0x69, 0x6e, 0x64, 0x79,
	0x6b, 0x69, 0x74, 0x65, 0x2e, 0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x76, 0x31, 0x62,
	0x65, 0x74, 0x61, 0x31, 0x42, 0x07, 0x49, 0x64, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a,
	0x4d, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x64, 0x79,
	0x6b, 0x69, 0x74, 0x65, 0x2f, 0x6a, 0x61, 0x72, 0x76, 0x69, 0x73, 0x2d, 0x73, 0x64, 0x6b, 0x2d,
	0x67, 0x6f, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x69, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2f,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2f, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x3b,
	0x6f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x76, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xa2, 0x02,
	0x03, 0x49, 0x4f, 0x58, 0xaa, 0x02, 0x18, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x2e,
	0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x2e, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xca,
	0x02, 0x18, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x5c, 0x4f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x73, 0x5c, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0xe2, 0x02, 0x24, 0x49, 0x6e, 0x64,
	0x79, 0x6b, 0x69, 0x74, 0x65, 0x5c, 0x4f, 0x62, 0x6a, 0x65, 0x63, 0x74, 0x73, 0x5c, 0x56, 0x31,
	0x62, 0x65, 0x74, 0x61, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x1a, 0x49, 0x6e, 0x64, 0x79, 0x6b, 0x69, 0x74, 0x65, 0x3a, 0x3a, 0x4f, 0x62,
	0x6a, 0x65, 0x63, 0x74, 0x73, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x65, 0x74, 0x61, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_indykite_objects_v1beta1_id_proto_rawDescOnce sync.Once
	file_indykite_objects_v1beta1_id_proto_rawDescData = file_indykite_objects_v1beta1_id_proto_rawDesc
)

func file_indykite_objects_v1beta1_id_proto_rawDescGZIP() []byte {
	file_indykite_objects_v1beta1_id_proto_rawDescOnce.Do(func() {
		file_indykite_objects_v1beta1_id_proto_rawDescData = protoimpl.X.CompressGZIP(file_indykite_objects_v1beta1_id_proto_rawDescData)
	})
	return file_indykite_objects_v1beta1_id_proto_rawDescData
}

var file_indykite_objects_v1beta1_id_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_indykite_objects_v1beta1_id_proto_goTypes = []interface{}{
	(*Identifier)(nil),      // 0: indykite.objects.v1beta1.Identifier
	(*ObjectReference)(nil), // 1: indykite.objects.v1beta1.ObjectReference
}
var file_indykite_objects_v1beta1_id_proto_depIdxs = []int32{
	0, // 0: indykite.objects.v1beta1.ObjectReference.customer_id:type_name -> indykite.objects.v1beta1.Identifier
	0, // 1: indykite.objects.v1beta1.ObjectReference.app_space_id:type_name -> indykite.objects.v1beta1.Identifier
	0, // 2: indykite.objects.v1beta1.ObjectReference.app_id:type_name -> indykite.objects.v1beta1.Identifier
	0, // 3: indykite.objects.v1beta1.ObjectReference.tenant_id:type_name -> indykite.objects.v1beta1.Identifier
	0, // 4: indykite.objects.v1beta1.ObjectReference.id:type_name -> indykite.objects.v1beta1.Identifier
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_indykite_objects_v1beta1_id_proto_init() }
func file_indykite_objects_v1beta1_id_proto_init() {
	if File_indykite_objects_v1beta1_id_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_indykite_objects_v1beta1_id_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Identifier); i {
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
		file_indykite_objects_v1beta1_id_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ObjectReference); i {
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
	file_indykite_objects_v1beta1_id_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Identifier_IdString)(nil),
		(*Identifier_IdBytes)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_indykite_objects_v1beta1_id_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_indykite_objects_v1beta1_id_proto_goTypes,
		DependencyIndexes: file_indykite_objects_v1beta1_id_proto_depIdxs,
		MessageInfos:      file_indykite_objects_v1beta1_id_proto_msgTypes,
	}.Build()
	File_indykite_objects_v1beta1_id_proto = out.File
	file_indykite_objects_v1beta1_id_proto_rawDesc = nil
	file_indykite_objects_v1beta1_id_proto_goTypes = nil
	file_indykite_objects_v1beta1_id_proto_depIdxs = nil
}
