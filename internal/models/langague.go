package model

import snowflakeid "github.com/Mirai3103/remote-compiler/pkg/snowflake_id"

type Language struct {
	SourceFileExt  *string `yaml:"source_file" json:"sourceFileExt"`
	BinaryFileExt  *string `yaml:"binary_file" json:"binaryFileExt"`
	CompileCommand *string `yaml:"compile_command" json:"compileCommand"`
	RunCommand     *string `yaml:"run_command" json:"runCommand"`
	sourceFileName *string `yaml:"-" json:"-"`
	binaryFileName *string `yaml:"-" json:"-"`
}

func (l *Language) GetSourceFileName() string {
	if l.SourceFileExt == nil {
		return ""
	}
	if l.sourceFileName != nil {
		return *l.sourceFileName
	}
	newFileName := snowflakeid.NewString() + *l.SourceFileExt
	l.sourceFileName = &newFileName
	return newFileName
}

func (l *Language) GetBinaryFileName() string {
	if l.BinaryFileExt == nil {
		return ""
	}
	if l.binaryFileName != nil {
		return *l.binaryFileName
	}
	newFileName := snowflakeid.NewString() + *l.BinaryFileExt
	l.binaryFileName = &newFileName
	return newFileName
}
