package internal

type Language struct {
	Version        string `yaml:"version" json:"version"`
	Name           string `yaml:"name" json:"name"`
	SourceFileExt  string `yaml:"source_file" json:"sourceFileExt"`
	BinaryFileExt  string `yaml:"binary_file" json:"binaryFileExt"`
	CompileCommand string `yaml:"compile_command" json:"compileCommand"`
	RunCommand     string `yaml:"run_command" json:"runCommand"`
}
