package internal

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

func (l *Language) GenerateSourceFileName() string {
	newFileName := uuid.NewString() + *l.SourceFileExt
	l.SourceFileName = &newFileName
	return newFileName
}

func (l *Language) GenerateBinaryFileName() string {
	newFileName := uuid.NewString() + *l.BinaryFileExt
	l.BinaryFileName = &newFileName
	return newFileName
}
