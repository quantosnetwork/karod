// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.19.4
// source: p2p/proto/messages.proto

package p2p

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

type DataType int32

const (
	DataType_UNKNOWN               DataType = 0
	DataType_AUTH                  DataType = 1
	DataType_BLOCKCHAIN_MISC_ASSET DataType = 2
	DataType_TRANSACTION           DataType = 3
	DataType_WALLET                DataType = 4
	DataType_BLOCK                 DataType = 5
	DataType_CONTRACT              DataType = 6
	DataType_DNS                   DataType = 7
	DataType_ERROR                 DataType = 8
	DataType_INFO                  DataType = 9
	DataType_META                  DataType = 10
	DataType_NFT                   DataType = 11
	DataType_CONTENT               DataType = 12
	DataType_STATE                 DataType = 13
	DataType_INDEX                 DataType = 14
	DataType_RECEIPT               DataType = 15
	DataType_HEADER                DataType = 16
	DataType_REDIRECT              DataType = 17
	DataType_PLUGIN                DataType = 18
	DataType_ADDRESS               DataType = 19
	DataType_KEY                   DataType = 20
	DataType_ROUTER                DataType = 21
	DataType_PROXY                 DataType = 22
	DataType_QUANTIX               DataType = 23
	DataType_VERSION               DataType = 24
	DataType_FILE                  DataType = 25
	DataType_ORACLE                DataType = 26
	DataType_GENESIS               DataType = 27
	DataType_BRIDGE                DataType = 28
	DataType_CHILD_CHAIN_LAYERX    DataType = 29
	DataType_COMMAND               DataType = 998
	DataType_INTERNAL              DataType = 999
	DataType_EGG                   DataType = 99887
)

// Enum value maps for DataType.
var (
	DataType_name = map[int32]string{
		0:     "UNKNOWN",
		1:     "AUTH",
		2:     "BLOCKCHAIN_MISC_ASSET",
		3:     "TRANSACTION",
		4:     "WALLET",
		5:     "BLOCK",
		6:     "CONTRACT",
		7:     "DNS",
		8:     "ERROR",
		9:     "INFO",
		10:    "META",
		11:    "NFT",
		12:    "CONTENT",
		13:    "STATE",
		14:    "INDEX",
		15:    "RECEIPT",
		16:    "HEADER",
		17:    "REDIRECT",
		18:    "PLUGIN",
		19:    "ADDRESS",
		20:    "KEY",
		21:    "ROUTER",
		22:    "PROXY",
		23:    "QUANTIX",
		24:    "VERSION",
		25:    "FILE",
		26:    "ORACLE",
		27:    "GENESIS",
		28:    "BRIDGE",
		29:    "CHILD_CHAIN_LAYERX",
		998:   "COMMAND",
		999:   "INTERNAL",
		99887: "EGG",
	}
	DataType_value = map[string]int32{
		"UNKNOWN":               0,
		"AUTH":                  1,
		"BLOCKCHAIN_MISC_ASSET": 2,
		"TRANSACTION":           3,
		"WALLET":                4,
		"BLOCK":                 5,
		"CONTRACT":              6,
		"DNS":                   7,
		"ERROR":                 8,
		"INFO":                  9,
		"META":                  10,
		"NFT":                   11,
		"CONTENT":               12,
		"STATE":                 13,
		"INDEX":                 14,
		"RECEIPT":               15,
		"HEADER":                16,
		"REDIRECT":              17,
		"PLUGIN":                18,
		"ADDRESS":               19,
		"KEY":                   20,
		"ROUTER":                21,
		"PROXY":                 22,
		"QUANTIX":               23,
		"VERSION":               24,
		"FILE":                  25,
		"ORACLE":                26,
		"GENESIS":               27,
		"BRIDGE":                28,
		"CHILD_CHAIN_LAYERX":    29,
		"COMMAND":               998,
		"INTERNAL":              999,
		"EGG":                   99887,
	}
)

func (x DataType) Enum() *DataType {
	p := new(DataType)
	*p = x
	return p
}

func (x DataType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (DataType) Descriptor() protoreflect.EnumDescriptor {
	return file_p2p_proto_messages_proto_enumTypes[0].Descriptor()
}

func (DataType) Type() protoreflect.EnumType {
	return &file_p2p_proto_messages_proto_enumTypes[0]
}

func (x DataType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use DataType.Descriptor instead.
func (DataType) EnumDescriptor() ([]byte, []int) {
	return file_p2p_proto_messages_proto_rawDescGZIP(), []int{0}
}

type ApiRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AuthToken string `protobuf:"bytes,1,opt,name=authToken,proto3" json:"authToken,omitempty"`
	ReqID     string `protobuf:"bytes,2,opt,name=ReqID,proto3" json:"ReqID,omitempty"`
	ApiPath   string `protobuf:"bytes,3,opt,name=ApiPath,proto3" json:"ApiPath,omitempty"`
	ReqData   []byte `protobuf:"bytes,4,opt,name=ReqData,proto3" json:"ReqData,omitempty"`
}

func (x *ApiRequest) Reset() {
	*x = ApiRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2p_proto_messages_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiRequest) ProtoMessage() {}

func (x *ApiRequest) ProtoReflect() protoreflect.Message {
	mi := &file_p2p_proto_messages_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiRequest.ProtoReflect.Descriptor instead.
func (*ApiRequest) Descriptor() ([]byte, []int) {
	return file_p2p_proto_messages_proto_rawDescGZIP(), []int{0}
}

func (x *ApiRequest) GetAuthToken() string {
	if x != nil {
		return x.AuthToken
	}
	return ""
}

func (x *ApiRequest) GetReqID() string {
	if x != nil {
		return x.ReqID
	}
	return ""
}

func (x *ApiRequest) GetApiPath() string {
	if x != nil {
		return x.ApiPath
	}
	return ""
}

func (x *ApiRequest) GetReqData() []byte {
	if x != nil {
		return x.ReqData
	}
	return nil
}

type ApiResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReqID       string   `protobuf:"bytes,1,opt,name=ReqID,proto3" json:"ReqID,omitempty"`
	ResID       string   `protobuf:"bytes,2,opt,name=ResID,proto3" json:"ResID,omitempty"`
	ResData     []byte   `protobuf:"bytes,3,opt,name=ResData,proto3" json:"ResData,omitempty"`
	ResDataType DataType `protobuf:"varint,4,opt,name=ResDataType,proto3,enum=p2ppb.DataType" json:"ResDataType,omitempty"`
	StatusCode  int32    `protobuf:"varint,5,opt,name=StatusCode,proto3" json:"StatusCode,omitempty"`
}

func (x *ApiResponse) Reset() {
	*x = ApiResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2p_proto_messages_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ApiResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ApiResponse) ProtoMessage() {}

func (x *ApiResponse) ProtoReflect() protoreflect.Message {
	mi := &file_p2p_proto_messages_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ApiResponse.ProtoReflect.Descriptor instead.
func (*ApiResponse) Descriptor() ([]byte, []int) {
	return file_p2p_proto_messages_proto_rawDescGZIP(), []int{1}
}

func (x *ApiResponse) GetReqID() string {
	if x != nil {
		return x.ReqID
	}
	return ""
}

func (x *ApiResponse) GetResID() string {
	if x != nil {
		return x.ResID
	}
	return ""
}

func (x *ApiResponse) GetResData() []byte {
	if x != nil {
		return x.ResData
	}
	return nil
}

func (x *ApiResponse) GetResDataType() DataType {
	if x != nil {
		return x.ResDataType
	}
	return DataType_UNKNOWN
}

func (x *ApiResponse) GetStatusCode() int32 {
	if x != nil {
		return x.StatusCode
	}
	return 0
}

type Error struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReqID   string `protobuf:"bytes,1,opt,name=ReqID,proto3" json:"ReqID,omitempty"`
	Code    int32  `protobuf:"varint,2,opt,name=Code,proto3" json:"Code,omitempty"`
	Message string `protobuf:"bytes,3,opt,name=Message,proto3" json:"Message,omitempty"`
}

func (x *Error) Reset() {
	*x = Error{}
	if protoimpl.UnsafeEnabled {
		mi := &file_p2p_proto_messages_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Error) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Error) ProtoMessage() {}

func (x *Error) ProtoReflect() protoreflect.Message {
	mi := &file_p2p_proto_messages_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Error.ProtoReflect.Descriptor instead.
func (*Error) Descriptor() ([]byte, []int) {
	return file_p2p_proto_messages_proto_rawDescGZIP(), []int{2}
}

func (x *Error) GetReqID() string {
	if x != nil {
		return x.ReqID
	}
	return ""
}

func (x *Error) GetCode() int32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Error) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_p2p_proto_messages_proto protoreflect.FileDescriptor

var file_p2p_proto_messages_proto_rawDesc = []byte{
	0x0a, 0x18, 0x70, 0x32, 0x70, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x6d, 0x65, 0x73, 0x73,
	0x61, 0x67, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x32, 0x70, 0x70,
	0x62, 0x22, 0x74, 0x0a, 0x0a, 0x41, 0x70, 0x69, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1c, 0x0a, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x09, 0x61, 0x75, 0x74, 0x68, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x14, 0x0a,
	0x05, 0x52, 0x65, 0x71, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x52, 0x65,
	0x71, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x41, 0x70, 0x69, 0x50, 0x61, 0x74, 0x68, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x41, 0x70, 0x69, 0x50, 0x61, 0x74, 0x68, 0x12, 0x18, 0x0a,
	0x07, 0x52, 0x65, 0x71, 0x44, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x07,
	0x52, 0x65, 0x71, 0x44, 0x61, 0x74, 0x61, 0x22, 0xa6, 0x01, 0x0a, 0x0b, 0x41, 0x70, 0x69, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x65, 0x71, 0x49, 0x44,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x52, 0x65, 0x71, 0x49, 0x44, 0x12, 0x14, 0x0a,
	0x05, 0x52, 0x65, 0x73, 0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x52, 0x65,
	0x73, 0x49, 0x44, 0x12, 0x18, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x07, 0x52, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x12, 0x31, 0x0a,
	0x0b, 0x52, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x0f, 0x2e, 0x70, 0x32, 0x70, 0x70, 0x62, 0x2e, 0x44, 0x61, 0x74, 0x61, 0x54,
	0x79, 0x70, 0x65, 0x52, 0x0b, 0x52, 0x65, 0x73, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65,
	0x12, 0x1e, 0x0a, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x05,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x43, 0x6f, 0x64, 0x65,
	0x22, 0x4b, 0x0a, 0x05, 0x45, 0x72, 0x72, 0x6f, 0x72, 0x12, 0x14, 0x0a, 0x05, 0x52, 0x65, 0x71,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x52, 0x65, 0x71, 0x49, 0x44, 0x12,
	0x12, 0x0a, 0x04, 0x43, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x43,
	0x6f, 0x64, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x03,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2a, 0xaf, 0x03,
	0x0a, 0x08, 0x44, 0x61, 0x74, 0x61, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x4e,
	0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x41, 0x55, 0x54, 0x48, 0x10,
	0x01, 0x12, 0x19, 0x0a, 0x15, 0x42, 0x4c, 0x4f, 0x43, 0x4b, 0x43, 0x48, 0x41, 0x49, 0x4e, 0x5f,
	0x4d, 0x49, 0x53, 0x43, 0x5f, 0x41, 0x53, 0x53, 0x45, 0x54, 0x10, 0x02, 0x12, 0x0f, 0x0a, 0x0b,
	0x54, 0x52, 0x41, 0x4e, 0x53, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x10, 0x03, 0x12, 0x0a, 0x0a,
	0x06, 0x57, 0x41, 0x4c, 0x4c, 0x45, 0x54, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x42, 0x4c, 0x4f,
	0x43, 0x4b, 0x10, 0x05, 0x12, 0x0c, 0x0a, 0x08, 0x43, 0x4f, 0x4e, 0x54, 0x52, 0x41, 0x43, 0x54,
	0x10, 0x06, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x4e, 0x53, 0x10, 0x07, 0x12, 0x09, 0x0a, 0x05, 0x45,
	0x52, 0x52, 0x4f, 0x52, 0x10, 0x08, 0x12, 0x08, 0x0a, 0x04, 0x49, 0x4e, 0x46, 0x4f, 0x10, 0x09,
	0x12, 0x08, 0x0a, 0x04, 0x4d, 0x45, 0x54, 0x41, 0x10, 0x0a, 0x12, 0x07, 0x0a, 0x03, 0x4e, 0x46,
	0x54, 0x10, 0x0b, 0x12, 0x0b, 0x0a, 0x07, 0x43, 0x4f, 0x4e, 0x54, 0x45, 0x4e, 0x54, 0x10, 0x0c,
	0x12, 0x09, 0x0a, 0x05, 0x53, 0x54, 0x41, 0x54, 0x45, 0x10, 0x0d, 0x12, 0x09, 0x0a, 0x05, 0x49,
	0x4e, 0x44, 0x45, 0x58, 0x10, 0x0e, 0x12, 0x0b, 0x0a, 0x07, 0x52, 0x45, 0x43, 0x45, 0x49, 0x50,
	0x54, 0x10, 0x0f, 0x12, 0x0a, 0x0a, 0x06, 0x48, 0x45, 0x41, 0x44, 0x45, 0x52, 0x10, 0x10, 0x12,
	0x0c, 0x0a, 0x08, 0x52, 0x45, 0x44, 0x49, 0x52, 0x45, 0x43, 0x54, 0x10, 0x11, 0x12, 0x0a, 0x0a,
	0x06, 0x50, 0x4c, 0x55, 0x47, 0x49, 0x4e, 0x10, 0x12, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x44, 0x44,
	0x52, 0x45, 0x53, 0x53, 0x10, 0x13, 0x12, 0x07, 0x0a, 0x03, 0x4b, 0x45, 0x59, 0x10, 0x14, 0x12,
	0x0a, 0x0a, 0x06, 0x52, 0x4f, 0x55, 0x54, 0x45, 0x52, 0x10, 0x15, 0x12, 0x09, 0x0a, 0x05, 0x50,
	0x52, 0x4f, 0x58, 0x59, 0x10, 0x16, 0x12, 0x0b, 0x0a, 0x07, 0x51, 0x55, 0x41, 0x4e, 0x54, 0x49,
	0x58, 0x10, 0x17, 0x12, 0x0b, 0x0a, 0x07, 0x56, 0x45, 0x52, 0x53, 0x49, 0x4f, 0x4e, 0x10, 0x18,
	0x12, 0x08, 0x0a, 0x04, 0x46, 0x49, 0x4c, 0x45, 0x10, 0x19, 0x12, 0x0a, 0x0a, 0x06, 0x4f, 0x52,
	0x41, 0x43, 0x4c, 0x45, 0x10, 0x1a, 0x12, 0x0b, 0x0a, 0x07, 0x47, 0x45, 0x4e, 0x45, 0x53, 0x49,
	0x53, 0x10, 0x1b, 0x12, 0x0a, 0x0a, 0x06, 0x42, 0x52, 0x49, 0x44, 0x47, 0x45, 0x10, 0x1c, 0x12,
	0x16, 0x0a, 0x12, 0x43, 0x48, 0x49, 0x4c, 0x44, 0x5f, 0x43, 0x48, 0x41, 0x49, 0x4e, 0x5f, 0x4c,
	0x41, 0x59, 0x45, 0x52, 0x58, 0x10, 0x1d, 0x12, 0x0c, 0x0a, 0x07, 0x43, 0x4f, 0x4d, 0x4d, 0x41,
	0x4e, 0x44, 0x10, 0xe6, 0x07, 0x12, 0x0d, 0x0a, 0x08, 0x49, 0x4e, 0x54, 0x45, 0x52, 0x4e, 0x41,
	0x4c, 0x10, 0xe7, 0x07, 0x12, 0x09, 0x0a, 0x03, 0x45, 0x47, 0x47, 0x10, 0xaf, 0x8c, 0x06, 0x42,
	0x08, 0x5a, 0x06, 0x70, 0x62, 0x2f, 0x70, 0x32, 0x70, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_p2p_proto_messages_proto_rawDescOnce sync.Once
	file_p2p_proto_messages_proto_rawDescData = file_p2p_proto_messages_proto_rawDesc
)

func file_p2p_proto_messages_proto_rawDescGZIP() []byte {
	file_p2p_proto_messages_proto_rawDescOnce.Do(func() {
		file_p2p_proto_messages_proto_rawDescData = protoimpl.X.CompressGZIP(file_p2p_proto_messages_proto_rawDescData)
	})
	return file_p2p_proto_messages_proto_rawDescData
}

var file_p2p_proto_messages_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_p2p_proto_messages_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_p2p_proto_messages_proto_goTypes = []interface{}{
	(DataType)(0),       // 0: p2ppb.DataType
	(*ApiRequest)(nil),  // 1: p2ppb.ApiRequest
	(*ApiResponse)(nil), // 2: p2ppb.ApiResponse
	(*Error)(nil),       // 3: p2ppb.Error
}
var file_p2p_proto_messages_proto_depIdxs = []int32{
	0, // 0: p2ppb.ApiResponse.ResDataType:type_name -> p2ppb.DataType
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_p2p_proto_messages_proto_init() }
func file_p2p_proto_messages_proto_init() {
	if File_p2p_proto_messages_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_p2p_proto_messages_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiRequest); i {
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
		file_p2p_proto_messages_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ApiResponse); i {
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
		file_p2p_proto_messages_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Error); i {
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
			RawDescriptor: file_p2p_proto_messages_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_p2p_proto_messages_proto_goTypes,
		DependencyIndexes: file_p2p_proto_messages_proto_depIdxs,
		EnumInfos:         file_p2p_proto_messages_proto_enumTypes,
		MessageInfos:      file_p2p_proto_messages_proto_msgTypes,
	}.Build()
	File_p2p_proto_messages_proto = out.File
	file_p2p_proto_messages_proto_rawDesc = nil
	file_p2p_proto_messages_proto_goTypes = nil
	file_p2p_proto_messages_proto_depIdxs = nil
}