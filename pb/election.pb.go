// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.1
// 	protoc        v5.27.1
// source: proto/election.proto

package pb

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

type VoteResponse_VoteType int32

const (
	VoteResponse_VOTE_REQUESTED VoteResponse_VoteType = 0
	VoteResponse_VOTE_GIVEN     VoteResponse_VoteType = 1
	VoteResponse_VOTE_REFUSED   VoteResponse_VoteType = 2
)

// Enum value maps for VoteResponse_VoteType.
var (
	VoteResponse_VoteType_name = map[int32]string{
		0: "VOTE_REQUESTED",
		1: "VOTE_GIVEN",
		2: "VOTE_REFUSED",
	}
	VoteResponse_VoteType_value = map[string]int32{
		"VOTE_REQUESTED": 0,
		"VOTE_GIVEN":     1,
		"VOTE_REFUSED":   2,
	}
)

func (x VoteResponse_VoteType) Enum() *VoteResponse_VoteType {
	p := new(VoteResponse_VoteType)
	*p = x
	return p
}

func (x VoteResponse_VoteType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (VoteResponse_VoteType) Descriptor() protoreflect.EnumDescriptor {
	return file_proto_election_proto_enumTypes[0].Descriptor()
}

func (VoteResponse_VoteType) Type() protoreflect.EnumType {
	return &file_proto_election_proto_enumTypes[0]
}

func (x VoteResponse_VoteType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use VoteResponse_VoteType.Descriptor instead.
func (VoteResponse_VoteType) EnumDescriptor() ([]byte, []int) {
	return file_proto_election_proto_rawDescGZIP(), []int{1, 0}
}

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	LogfileIndex uint64 `protobuf:"varint,1,opt,name=logfile_index,json=logfileIndex,proto3" json:"logfile_index,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_election_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_election_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_proto_election_proto_rawDescGZIP(), []int{0}
}

func (x *VoteRequest) GetLogfileIndex() uint64 {
	if x != nil {
		return x.LogfileIndex
	}
	return 0
}

type VoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	VoteType VoteResponse_VoteType `protobuf:"varint,1,opt,name=vote_type,json=voteType,proto3,enum=proto.VoteResponse_VoteType" json:"vote_type,omitempty"`
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_election_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_election_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_proto_election_proto_rawDescGZIP(), []int{1}
}

func (x *VoteResponse) GetVoteType() VoteResponse_VoteType {
	if x != nil {
		return x.VoteType
	}
	return VoteResponse_VOTE_REQUESTED
}

var File_proto_election_proto protoreflect.FileDescriptor

var file_proto_election_proto_rawDesc = []byte{
	0x0a, 0x14, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x32, 0x0a,
	0x0b, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x0d,
	0x6c, 0x6f, 0x67, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x69, 0x6e, 0x64, 0x65, 0x78, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0c, 0x6c, 0x6f, 0x67, 0x66, 0x69, 0x6c, 0x65, 0x49, 0x6e, 0x64, 0x65,
	0x78, 0x22, 0x8b, 0x01, 0x0a, 0x0c, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x39, 0x0a, 0x09, 0x76, 0x6f, 0x74, 0x65, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x6f,
	0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x08, 0x76, 0x6f, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x22, 0x40, 0x0a,
	0x08, 0x56, 0x6f, 0x74, 0x65, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x56, 0x4f, 0x54,
	0x45, 0x5f, 0x52, 0x45, 0x51, 0x55, 0x45, 0x53, 0x54, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a,
	0x0a, 0x56, 0x4f, 0x54, 0x45, 0x5f, 0x47, 0x49, 0x56, 0x45, 0x4e, 0x10, 0x01, 0x12, 0x10, 0x0a,
	0x0c, 0x56, 0x4f, 0x54, 0x45, 0x5f, 0x52, 0x45, 0x46, 0x55, 0x53, 0x45, 0x44, 0x10, 0x02, 0x32,
	0x46, 0x0a, 0x0f, 0x45, 0x6c, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x33, 0x0a, 0x06, 0x56, 0x6f, 0x74, 0x69, 0x6e, 0x67, 0x12, 0x12, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x06, 0x5a, 0x04, 0x2e, 0x2f, 0x70, 0x62, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_election_proto_rawDescOnce sync.Once
	file_proto_election_proto_rawDescData = file_proto_election_proto_rawDesc
)

func file_proto_election_proto_rawDescGZIP() []byte {
	file_proto_election_proto_rawDescOnce.Do(func() {
		file_proto_election_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_election_proto_rawDescData)
	})
	return file_proto_election_proto_rawDescData
}

var file_proto_election_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_proto_election_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_proto_election_proto_goTypes = []interface{}{
	(VoteResponse_VoteType)(0), // 0: proto.VoteResponse.VoteType
	(*VoteRequest)(nil),        // 1: proto.VoteRequest
	(*VoteResponse)(nil),       // 2: proto.VoteResponse
}
var file_proto_election_proto_depIdxs = []int32{
	0, // 0: proto.VoteResponse.vote_type:type_name -> proto.VoteResponse.VoteType
	1, // 1: proto.ElectionService.Voting:input_type -> proto.VoteRequest
	2, // 2: proto.ElectionService.Voting:output_type -> proto.VoteResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_proto_election_proto_init() }
func file_proto_election_proto_init() {
	if File_proto_election_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_proto_election_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
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
		file_proto_election_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResponse); i {
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
			RawDescriptor: file_proto_election_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_election_proto_goTypes,
		DependencyIndexes: file_proto_election_proto_depIdxs,
		EnumInfos:         file_proto_election_proto_enumTypes,
		MessageInfos:      file_proto_election_proto_msgTypes,
	}.Build()
	File_proto_election_proto = out.File
	file_proto_election_proto_rawDesc = nil
	file_proto_election_proto_goTypes = nil
	file_proto_election_proto_depIdxs = nil
}
