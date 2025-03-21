// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.5
// 	protoc        v5.29.3
// source: proto/execution_service.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type SubmissionSettings struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	WithTrim          bool                   `protobuf:"varint,1,opt,name=with_trim,json=withTrim,proto3" json:"with_trim,omitempty"`
	WithCaseSensitive bool                   `protobuf:"varint,2,opt,name=with_case_sensitive,json=withCaseSensitive,proto3" json:"with_case_sensitive,omitempty"`
	WithWhitespace    bool                   `protobuf:"varint,3,opt,name=with_whitespace,json=withWhitespace,proto3" json:"with_whitespace,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *SubmissionSettings) Reset() {
	*x = SubmissionSettings{}
	mi := &file_proto_execution_service_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmissionSettings) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmissionSettings) ProtoMessage() {}

func (x *SubmissionSettings) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execution_service_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmissionSettings.ProtoReflect.Descriptor instead.
func (*SubmissionSettings) Descriptor() ([]byte, []int) {
	return file_proto_execution_service_proto_rawDescGZIP(), []int{0}
}

func (x *SubmissionSettings) GetWithTrim() bool {
	if x != nil {
		return x.WithTrim
	}
	return false
}

func (x *SubmissionSettings) GetWithCaseSensitive() bool {
	if x != nil {
		return x.WithCaseSensitive
	}
	return false
}

func (x *SubmissionSettings) GetWithWhitespace() bool {
	if x != nil {
		return x.WithWhitespace
	}
	return false
}

type Language struct {
	state          protoimpl.MessageState `protogen:"open.v1"`
	SourceFileExt  string                 `protobuf:"bytes,3,opt,name=source_file_ext,json=sourceFileExt,proto3" json:"source_file_ext,omitempty"`
	BinaryFileExt  string                 `protobuf:"bytes,4,opt,name=binary_file_ext,json=binaryFileExt,proto3" json:"binary_file_ext,omitempty"`
	CompileCommand string                 `protobuf:"bytes,5,opt,name=compile_command,json=compileCommand,proto3" json:"compile_command,omitempty"`
	RunCommand     string                 `protobuf:"bytes,6,opt,name=run_command,json=runCommand,proto3" json:"run_command,omitempty"`
	unknownFields  protoimpl.UnknownFields
	sizeCache      protoimpl.SizeCache
}

func (x *Language) Reset() {
	*x = Language{}
	mi := &file_proto_execution_service_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Language) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Language) ProtoMessage() {}

func (x *Language) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execution_service_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Language.ProtoReflect.Descriptor instead.
func (*Language) Descriptor() ([]byte, []int) {
	return file_proto_execution_service_proto_rawDescGZIP(), []int{1}
}

func (x *Language) GetSourceFileExt() string {
	if x != nil {
		return x.SourceFileExt
	}
	return ""
}

func (x *Language) GetBinaryFileExt() string {
	if x != nil {
		return x.BinaryFileExt
	}
	return ""
}

func (x *Language) GetCompileCommand() string {
	if x != nil {
		return x.CompileCommand
	}
	return ""
}

func (x *Language) GetRunCommand() string {
	if x != nil {
		return x.RunCommand
	}
	return ""
}

type SubmissionResult struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	SubmissionId    string                 `protobuf:"bytes,1,opt,name=submission_id,json=submissionId,proto3" json:"submission_id,omitempty"`
	TestCaseId      string                 `protobuf:"bytes,2,opt,name=test_case_id,json=testCaseId,proto3" json:"test_case_id,omitempty"`
	Status          string                 `protobuf:"bytes,3,opt,name=status,proto3" json:"status,omitempty"`
	Stdout          string                 `protobuf:"bytes,4,opt,name=stdout,proto3" json:"stdout,omitempty"`
	MemoryUsageInKb float32                `protobuf:"fixed32,5,opt,name=memory_usage_in_kb,json=memoryUsageInKb,proto3" json:"memory_usage_in_kb,omitempty"`
	TimeUsageInMs   float32                `protobuf:"fixed32,6,opt,name=time_usage_in_ms,json=timeUsageInMs,proto3" json:"time_usage_in_ms,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *SubmissionResult) Reset() {
	*x = SubmissionResult{}
	mi := &file_proto_execution_service_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *SubmissionResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SubmissionResult) ProtoMessage() {}

func (x *SubmissionResult) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execution_service_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SubmissionResult.ProtoReflect.Descriptor instead.
func (*SubmissionResult) Descriptor() ([]byte, []int) {
	return file_proto_execution_service_proto_rawDescGZIP(), []int{2}
}

func (x *SubmissionResult) GetSubmissionId() string {
	if x != nil {
		return x.SubmissionId
	}
	return ""
}

func (x *SubmissionResult) GetTestCaseId() string {
	if x != nil {
		return x.TestCaseId
	}
	return ""
}

func (x *SubmissionResult) GetStatus() string {
	if x != nil {
		return x.Status
	}
	return ""
}

func (x *SubmissionResult) GetStdout() string {
	if x != nil {
		return x.Stdout
	}
	return ""
}

func (x *SubmissionResult) GetMemoryUsageInKb() float32 {
	if x != nil {
		return x.MemoryUsageInKb
	}
	return 0
}

func (x *SubmissionResult) GetTimeUsageInMs() float32 {
	if x != nil {
		return x.TimeUsageInMs
	}
	return 0
}

type TestCase struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Input         string                 `protobuf:"bytes,2,opt,name=input,proto3" json:"input,omitempty"`
	ExpectOutput  string                 `protobuf:"bytes,3,opt,name=expect_output,json=expectOutput,proto3" json:"expect_output,omitempty"`
	InputFile     string                 `protobuf:"bytes,4,opt,name=input_file,json=inputFile,proto3" json:"input_file,omitempty"`
	OutputFile    string                 `protobuf:"bytes,5,opt,name=output_file,json=outputFile,proto3" json:"output_file,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TestCase) Reset() {
	*x = TestCase{}
	mi := &file_proto_execution_service_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TestCase) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestCase) ProtoMessage() {}

func (x *TestCase) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execution_service_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestCase.ProtoReflect.Descriptor instead.
func (*TestCase) Descriptor() ([]byte, []int) {
	return file_proto_execution_service_proto_rawDescGZIP(), []int{3}
}

func (x *TestCase) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *TestCase) GetInput() string {
	if x != nil {
		return x.Input
	}
	return ""
}

func (x *TestCase) GetExpectOutput() string {
	if x != nil {
		return x.ExpectOutput
	}
	return ""
}

func (x *TestCase) GetInputFile() string {
	if x != nil {
		return x.InputFile
	}
	return ""
}

func (x *TestCase) GetOutputFile() string {
	if x != nil {
		return x.OutputFile
	}
	return ""
}

type Submission struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	Id              string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Language        *Language              `protobuf:"bytes,2,opt,name=language,proto3" json:"language,omitempty"`
	Code            string                 `protobuf:"bytes,3,opt,name=code,proto3" json:"code,omitempty"`
	TimeLimitInMs   int32                  `protobuf:"varint,4,opt,name=time_limit_in_ms,json=timeLimitInMs,proto3" json:"time_limit_in_ms,omitempty"`
	MemoryLimitInKb int32                  `protobuf:"varint,5,opt,name=memory_limit_in_kb,json=memoryLimitInKb,proto3" json:"memory_limit_in_kb,omitempty"`
	TestCases       []*TestCase            `protobuf:"bytes,6,rep,name=test_cases,json=testCases,proto3" json:"test_cases,omitempty"`
	Settings        *SubmissionSettings    `protobuf:"bytes,7,opt,name=settings,proto3" json:"settings,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *Submission) Reset() {
	*x = Submission{}
	mi := &file_proto_execution_service_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Submission) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Submission) ProtoMessage() {}

func (x *Submission) ProtoReflect() protoreflect.Message {
	mi := &file_proto_execution_service_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Submission.ProtoReflect.Descriptor instead.
func (*Submission) Descriptor() ([]byte, []int) {
	return file_proto_execution_service_proto_rawDescGZIP(), []int{4}
}

func (x *Submission) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Submission) GetLanguage() *Language {
	if x != nil {
		return x.Language
	}
	return nil
}

func (x *Submission) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Submission) GetTimeLimitInMs() int32 {
	if x != nil {
		return x.TimeLimitInMs
	}
	return 0
}

func (x *Submission) GetMemoryLimitInKb() int32 {
	if x != nil {
		return x.MemoryLimitInKb
	}
	return 0
}

func (x *Submission) GetTestCases() []*TestCase {
	if x != nil {
		return x.TestCases
	}
	return nil
}

func (x *Submission) GetSettings() *SubmissionSettings {
	if x != nil {
		return x.Settings
	}
	return nil
}

var File_proto_execution_service_proto protoreflect.FileDescriptor

var file_proto_execution_service_proto_rawDesc = string([]byte{
	0x0a, 0x1d, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f,
	0x6e, 0x5f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x8a, 0x01, 0x0a, 0x12, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65,
	0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x12, 0x1b, 0x0a, 0x09, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x74,
	0x72, 0x69, 0x6d, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x08, 0x77, 0x69, 0x74, 0x68, 0x54,
	0x72, 0x69, 0x6d, 0x12, 0x2e, 0x0a, 0x13, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x63, 0x61, 0x73, 0x65,
	0x5f, 0x73, 0x65, 0x6e, 0x73, 0x69, 0x74, 0x69, 0x76, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x11, 0x77, 0x69, 0x74, 0x68, 0x43, 0x61, 0x73, 0x65, 0x53, 0x65, 0x6e, 0x73, 0x69, 0x74,
	0x69, 0x76, 0x65, 0x12, 0x27, 0x0a, 0x0f, 0x77, 0x69, 0x74, 0x68, 0x5f, 0x77, 0x68, 0x69, 0x74,
	0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0e, 0x77, 0x69,
	0x74, 0x68, 0x57, 0x68, 0x69, 0x74, 0x65, 0x73, 0x70, 0x61, 0x63, 0x65, 0x22, 0xa4, 0x01, 0x0a,
	0x08, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x12, 0x26, 0x0a, 0x0f, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x5f, 0x65, 0x78, 0x74, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0d, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x46, 0x69, 0x6c, 0x65, 0x45, 0x78,
	0x74, 0x12, 0x26, 0x0a, 0x0f, 0x62, 0x69, 0x6e, 0x61, 0x72, 0x79, 0x5f, 0x66, 0x69, 0x6c, 0x65,
	0x5f, 0x65, 0x78, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0d, 0x62, 0x69, 0x6e, 0x61,
	0x72, 0x79, 0x46, 0x69, 0x6c, 0x65, 0x45, 0x78, 0x74, 0x12, 0x27, 0x0a, 0x0f, 0x63, 0x6f, 0x6d,
	0x70, 0x69, 0x6c, 0x65, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0e, 0x63, 0x6f, 0x6d, 0x70, 0x69, 0x6c, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x61,
	0x6e, 0x64, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x75, 0x6e, 0x5f, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x72, 0x75, 0x6e, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x22, 0xdf, 0x01, 0x0a, 0x10, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69,
	0x6f, 0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x73, 0x75, 0x62, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x73, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x49, 0x64, 0x12, 0x20, 0x0a,
	0x0c, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x63, 0x61, 0x73, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x74, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x49, 0x64, 0x12,
	0x16, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x74, 0x64, 0x6f, 0x75, 0x74, 0x12,
	0x2b, 0x0a, 0x12, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x5f,
	0x69, 0x6e, 0x5f, 0x6b, 0x62, 0x18, 0x05, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0f, 0x6d, 0x65, 0x6d,
	0x6f, 0x72, 0x79, 0x55, 0x73, 0x61, 0x67, 0x65, 0x49, 0x6e, 0x4b, 0x62, 0x12, 0x27, 0x0a, 0x10,
	0x74, 0x69, 0x6d, 0x65, 0x5f, 0x75, 0x73, 0x61, 0x67, 0x65, 0x5f, 0x69, 0x6e, 0x5f, 0x6d, 0x73,
	0x18, 0x06, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x55, 0x73, 0x61, 0x67,
	0x65, 0x49, 0x6e, 0x4d, 0x73, 0x22, 0x95, 0x01, 0x0a, 0x08, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61,
	0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x78, 0x70, 0x65,
	0x63, 0x74, 0x5f, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x65, 0x78, 0x70, 0x65, 0x63, 0x74, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x1d, 0x0a,
	0x0a, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x09, 0x69, 0x6e, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b,
	0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x5f, 0x66, 0x69, 0x6c, 0x65, 0x18, 0x05, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x46, 0x69, 0x6c, 0x65, 0x22, 0x88, 0x02,
	0x0a, 0x0a, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x25, 0x0a, 0x08,
	0x6c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x09,
	0x2e, 0x4c, 0x61, 0x6e, 0x67, 0x75, 0x61, 0x67, 0x65, 0x52, 0x08, 0x6c, 0x61, 0x6e, 0x67, 0x75,
	0x61, 0x67, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x27, 0x0a, 0x10, 0x74, 0x69, 0x6d, 0x65, 0x5f,
	0x6c, 0x69, 0x6d, 0x69, 0x74, 0x5f, 0x69, 0x6e, 0x5f, 0x6d, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x0d, 0x74, 0x69, 0x6d, 0x65, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x49, 0x6e, 0x4d, 0x73,
	0x12, 0x2b, 0x0a, 0x12, 0x6d, 0x65, 0x6d, 0x6f, 0x72, 0x79, 0x5f, 0x6c, 0x69, 0x6d, 0x69, 0x74,
	0x5f, 0x69, 0x6e, 0x5f, 0x6b, 0x62, 0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x0f, 0x6d, 0x65,
	0x6d, 0x6f, 0x72, 0x79, 0x4c, 0x69, 0x6d, 0x69, 0x74, 0x49, 0x6e, 0x4b, 0x62, 0x12, 0x28, 0x0a,
	0x0a, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x63, 0x61, 0x73, 0x65, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x09, 0x2e, 0x54, 0x65, 0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x52, 0x09, 0x74, 0x65,
	0x73, 0x74, 0x43, 0x61, 0x73, 0x65, 0x73, 0x12, 0x2f, 0x0a, 0x08, 0x73, 0x65, 0x74, 0x74, 0x69,
	0x6e, 0x67, 0x73, 0x18, 0x07, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x53, 0x75, 0x62, 0x6d,
	0x69, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x52, 0x08,
	0x73, 0x65, 0x74, 0x74, 0x69, 0x6e, 0x67, 0x73, 0x32, 0x41, 0x0a, 0x10, 0x45, 0x78, 0x65, 0x63,
	0x75, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x2d, 0x0a, 0x07,
	0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x65, 0x12, 0x0b, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73,
	0x73, 0x69, 0x6f, 0x6e, 0x1a, 0x11, 0x2e, 0x53, 0x75, 0x62, 0x6d, 0x69, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x00, 0x30, 0x01, 0x42, 0x2c, 0x5a, 0x2a, 0x67,
	0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x4d, 0x69, 0x72, 0x61, 0x69, 0x33,
	0x31, 0x30, 0x33, 0x2f, 0x72, 0x65, 0x6d, 0x6f, 0x74, 0x65, 0x2d, 0x63, 0x6f, 0x6d, 0x70, 0x69,
	0x6c, 0x65, 0x72, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
})

var (
	file_proto_execution_service_proto_rawDescOnce sync.Once
	file_proto_execution_service_proto_rawDescData []byte
)

func file_proto_execution_service_proto_rawDescGZIP() []byte {
	file_proto_execution_service_proto_rawDescOnce.Do(func() {
		file_proto_execution_service_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_execution_service_proto_rawDesc), len(file_proto_execution_service_proto_rawDesc)))
	})
	return file_proto_execution_service_proto_rawDescData
}

var file_proto_execution_service_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_proto_execution_service_proto_goTypes = []any{
	(*SubmissionSettings)(nil), // 0: SubmissionSettings
	(*Language)(nil),           // 1: Language
	(*SubmissionResult)(nil),   // 2: SubmissionResult
	(*TestCase)(nil),           // 3: TestCase
	(*Submission)(nil),         // 4: Submission
}
var file_proto_execution_service_proto_depIdxs = []int32{
	1, // 0: Submission.language:type_name -> Language
	3, // 1: Submission.test_cases:type_name -> TestCase
	0, // 2: Submission.settings:type_name -> SubmissionSettings
	4, // 3: ExecutionService.Execute:input_type -> Submission
	2, // 4: ExecutionService.Execute:output_type -> SubmissionResult
	4, // [4:5] is the sub-list for method output_type
	3, // [3:4] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_proto_execution_service_proto_init() }
func file_proto_execution_service_proto_init() {
	if File_proto_execution_service_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_execution_service_proto_rawDesc), len(file_proto_execution_service_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_execution_service_proto_goTypes,
		DependencyIndexes: file_proto_execution_service_proto_depIdxs,
		MessageInfos:      file_proto_execution_service_proto_msgTypes,
	}.Build()
	File_proto_execution_service_proto = out.File
	file_proto_execution_service_proto_goTypes = nil
	file_proto_execution_service_proto_depIdxs = nil
}
