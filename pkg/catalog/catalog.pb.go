// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.31.0
// 	protoc        (unknown)
// source: catalog/catalog.proto

package catalog

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Entry_AddressType int32

const (
	// Deprecated.
	// Unknown address type (should only exist for old commits)
	// is resolved (to Relative or Full) by the prefix of the address.
	Entry_BY_PREFIX_DEPRECATED Entry_AddressType = 0
	Entry_RELATIVE             Entry_AddressType = 1
	Entry_FULL                 Entry_AddressType = 2
)

// Enum value maps for Entry_AddressType.
var (
	Entry_AddressType_name = map[int32]string{
		0: "BY_PREFIX_DEPRECATED",
		1: "RELATIVE",
		2: "FULL",
	}
	Entry_AddressType_value = map[string]int32{
		"BY_PREFIX_DEPRECATED": 0,
		"RELATIVE":             1,
		"FULL":                 2,
	}
)

func (x Entry_AddressType) Enum() *Entry_AddressType {
	p := new(Entry_AddressType)
	*p = x
	return p
}

func (x Entry_AddressType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Entry_AddressType) Descriptor() protoreflect.EnumDescriptor {
	return file_catalog_catalog_proto_enumTypes[0].Descriptor()
}

func (Entry_AddressType) Type() protoreflect.EnumType {
	return &file_catalog_catalog_proto_enumTypes[0]
}

func (x Entry_AddressType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Entry_AddressType.Descriptor instead.
func (Entry_AddressType) EnumDescriptor() ([]byte, []int) {
	return file_catalog_catalog_proto_rawDescGZIP(), []int{0, 0}
}

type Entry struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address      string                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	LastModified *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_modified,json=lastModified,proto3" json:"last_modified,omitempty"`
	Size         int64                  `protobuf:"varint,3,opt,name=size,proto3" json:"size,omitempty"`
	ETag         string                 `protobuf:"bytes,4,opt,name=e_tag,json=eTag,proto3" json:"e_tag,omitempty"`
	Metadata     map[string]string      `protobuf:"bytes,5,rep,name=metadata,proto3" json:"metadata,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	AddressType  Entry_AddressType      `protobuf:"varint,6,opt,name=address_type,json=addressType,proto3,enum=catalog.Entry_AddressType" json:"address_type,omitempty"`
	ContentType  string                 `protobuf:"bytes,7,opt,name=content_type,json=contentType,proto3" json:"content_type,omitempty"`
}

func (x *Entry) Reset() {
	*x = Entry{}
	if protoimpl.UnsafeEnabled {
		mi := &file_catalog_catalog_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Entry) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Entry) ProtoMessage() {}

func (x *Entry) ProtoReflect() protoreflect.Message {
	mi := &file_catalog_catalog_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Entry.ProtoReflect.Descriptor instead.
func (*Entry) Descriptor() ([]byte, []int) {
	return file_catalog_catalog_proto_rawDescGZIP(), []int{0}
}

func (x *Entry) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *Entry) GetLastModified() *timestamppb.Timestamp {
	if x != nil {
		return x.LastModified
	}
	return nil
}

func (x *Entry) GetSize() int64 {
	if x != nil {
		return x.Size
	}
	return 0
}

func (x *Entry) GetETag() string {
	if x != nil {
		return x.ETag
	}
	return ""
}

func (x *Entry) GetMetadata() map[string]string {
	if x != nil {
		return x.Metadata
	}
	return nil
}

func (x *Entry) GetAddressType() Entry_AddressType {
	if x != nil {
		return x.AddressType
	}
	return Entry_BY_PREFIX_DEPRECATED
}

func (x *Entry) GetContentType() string {
	if x != nil {
		return x.ContentType
	}
	return ""
}

var File_catalog_catalog_proto protoreflect.FileDescriptor

var file_catalog_catalog_proto_rawDesc = []byte{
	0x0a, 0x15, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2f, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f,
	0x67, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67,
	0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0xa5, 0x03, 0x0a, 0x05, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x61,
	0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64,
	0x64, 0x72, 0x65, 0x73, 0x73, 0x12, 0x3f, 0x0a, 0x0d, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x6d, 0x6f,
	0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x6c, 0x61, 0x73, 0x74, 0x4d, 0x6f,
	0x64, 0x69, 0x66, 0x69, 0x65, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x73, 0x69, 0x7a, 0x65, 0x12, 0x13, 0x0a, 0x05, 0x65, 0x5f,
	0x74, 0x61, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x65, 0x54, 0x61, 0x67, 0x12,
	0x38, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x18, 0x05, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x1c, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x45, 0x6e, 0x74, 0x72,
	0x79, 0x2e, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52,
	0x08, 0x6d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x3d, 0x0a, 0x0c, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x1a, 0x2e, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x2e, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x2e,
	0x41, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x52, 0x0b, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x21, 0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x74,
	0x65, 0x6e, 0x74, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x54, 0x79, 0x70, 0x65, 0x1a, 0x3b, 0x0a, 0x0d, 0x4d,
	0x65, 0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x3f, 0x0a, 0x0b, 0x41, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x54, 0x79, 0x70, 0x65, 0x12, 0x18, 0x0a, 0x14, 0x42, 0x59, 0x5f, 0x50, 0x52,
	0x45, 0x46, 0x49, 0x58, 0x5f, 0x44, 0x45, 0x50, 0x52, 0x45, 0x43, 0x41, 0x54, 0x45, 0x44, 0x10,
	0x00, 0x12, 0x0c, 0x0a, 0x08, 0x52, 0x45, 0x4c, 0x41, 0x54, 0x49, 0x56, 0x45, 0x10, 0x01, 0x12,
	0x08, 0x0a, 0x04, 0x46, 0x55, 0x4c, 0x4c, 0x10, 0x02, 0x42, 0x24, 0x5a, 0x22, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72, 0x65, 0x65, 0x76, 0x65, 0x73, 0x65,
	0x2f, 0x6c, 0x61, 0x6b, 0x65, 0x66, 0x73, 0x2f, 0x63, 0x61, 0x74, 0x61, 0x6c, 0x6f, 0x67, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_catalog_catalog_proto_rawDescOnce sync.Once
	file_catalog_catalog_proto_rawDescData = file_catalog_catalog_proto_rawDesc
)

func file_catalog_catalog_proto_rawDescGZIP() []byte {
	file_catalog_catalog_proto_rawDescOnce.Do(func() {
		file_catalog_catalog_proto_rawDescData = protoimpl.X.CompressGZIP(file_catalog_catalog_proto_rawDescData)
	})
	return file_catalog_catalog_proto_rawDescData
}

var file_catalog_catalog_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_catalog_catalog_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_catalog_catalog_proto_goTypes = []interface{}{
	(Entry_AddressType)(0),        // 0: catalog.Entry.AddressType
	(*Entry)(nil),                 // 1: catalog.Entry
	nil,                           // 2: catalog.Entry.MetadataEntry
	(*timestamppb.Timestamp)(nil), // 3: google.protobuf.Timestamp
}
var file_catalog_catalog_proto_depIdxs = []int32{
	3, // 0: catalog.Entry.last_modified:type_name -> google.protobuf.Timestamp
	2, // 1: catalog.Entry.metadata:type_name -> catalog.Entry.MetadataEntry
	0, // 2: catalog.Entry.address_type:type_name -> catalog.Entry.AddressType
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_catalog_catalog_proto_init() }
func file_catalog_catalog_proto_init() {
	if File_catalog_catalog_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_catalog_catalog_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Entry); i {
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
			RawDescriptor: file_catalog_catalog_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_catalog_catalog_proto_goTypes,
		DependencyIndexes: file_catalog_catalog_proto_depIdxs,
		EnumInfos:         file_catalog_catalog_proto_enumTypes,
		MessageInfos:      file_catalog_catalog_proto_msgTypes,
	}.Build()
	File_catalog_catalog_proto = out.File
	file_catalog_catalog_proto_rawDesc = nil
	file_catalog_catalog_proto_goTypes = nil
	file_catalog_catalog_proto_depIdxs = nil
}
