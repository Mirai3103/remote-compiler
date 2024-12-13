package model

import "github.com/google/uuid"

type Language struct {
	Version        *string `yaml:"version" json:"version"`
	Name           *string `yaml:"name" json:"name"`
	SourceFileExt  *string `yaml:"source_file" json:"sourceFileExt"`
	BinaryFileExt  *string `yaml:"binary_file" json:"binaryFileExt"`
	CompileCommand *string `yaml:"compile_command" json:"compileCommand"`
	RunCommand     *string `yaml:"run_command" json:"runCommand"`
	SourceFileName *string `yaml:"-" json:"-"`
	BinaryFileName *string `yaml:"-" json:"-"`
}

func (l *Language) GetSourceFileName() string {
	if l.SourceFileName != nil {
		return *l.SourceFileName
	}
	newFileName := uuid.NewString() + *l.SourceFileExt
	l.SourceFileName = &newFileName
	return newFileName
}

func (l *Language) GetBinaryFileName() string {
	if l.BinaryFileName != nil {
		return *l.BinaryFileName
	}
	newFileName := uuid.NewString() + *l.BinaryFileExt
	l.BinaryFileName = &newFileName
	return newFileName
}
