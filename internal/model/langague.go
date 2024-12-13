package model

import "github.com/google/uuid"

type Language struct {
	Version        *string `yaml:"version" json:"version"`
	Name           *string `yaml:"name" json:"name"`
	SourceFileExt  *string `yaml:"source_file" json:"sourceFileExt"`
	BinaryFileExt  *string `yaml:"binary_file" json:"binaryFileExt"`
	CompileCommand *string `yaml:"compile_command" json:"compileCommand"`
	RunCommand     *string `yaml:"run_command" json:"runCommand"`
	sourceFileName *string `yaml:"-" json:"-"`
	binaryFileName *string `yaml:"-" json:"-"`
}

func (l *Language) GetSourceFileName() string {
	if l.sourceFileName != nil {
		return *l.sourceFileName
	}
	newFileName := uuid.NewString() + *l.SourceFileExt
	l.sourceFileName = &newFileName
	return newFileName
}

func (l *Language) GetBinaryFileName() string {
	if l.binaryFileName != nil {
		return *l.binaryFileName
	}
	newFileName := uuid.NewString() + *l.BinaryFileExt
	l.binaryFileName = &newFileName
	return newFileName
}
