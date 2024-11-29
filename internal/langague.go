package internal

type Language struct {
	Version        string `yaml:"version" json:"version"`
	Name           string `yaml:"name" json:"name"`
	SourceFile     string `yaml:"source_file" json:"source_file"`
	BinaryFile     string `yaml:"binary_file" json:"binary_file"`
	InputFile      string `yaml:"input_file" json:"input_file"`
	OutputFile     string `yaml:"output_file" json:"output_file"`
	CompileCommand string `yaml:"compile_command" json:"compile_command"`
	RunCommand     string `yaml:"run_command" json:"run_command"`
}

var CPP = Language{
	Version:        "17",
	Name:           "C++ 17",
	SourceFile:     "main.cpp",
	BinaryFile:     "main",
	InputFile:      "input.txt",
	OutputFile:     "output.txt",
	CompileCommand: "g++ -std=c++17 -O2 -o $BINARY_FILE $SOURCE_FILE",
	RunCommand:     "$BINARY_FILE $INPUT_FILE $OUTPUT_FILE",
}

var Shell = Language{
	Version:    "5",
	Name:       "Shell",
	SourceFile: "main.sh",
	InputFile:  "input.txt",
	OutputFile: "output.txt",
	RunCommand: "sh $SOURCE_FILE $INPUT_FILE $OUTPUT_FILE",
}

var TypeScript = Language{
	Version:        "3.7",
	Name:           "TypeScript 3.7",
	SourceFile:     "main.ts",
	BinaryFile:     "main.js",
	InputFile:      "input.txt",
	OutputFile:     "output.txt",
	CompileCommand: "tsc $SOURCE_FILE",
	RunCommand:     "node $BINARY_FILE $INPUT_FILE $OUTPUT_FILE",
}

var Python = Language{
	Version:    "3.8",
	Name:       "Python 3.8",
	SourceFile: "main.py",
	InputFile:  "input.txt",
	OutputFile: "output.txt",
	RunCommand: "python3 $SOURCE_FILE $INPUT_FILE $OUTPUT_FILE",
}

var Java = Language{
	Version:        "11",
	Name:           "Java 11",
	SourceFile:     "Main.java",
	BinaryFile:     "Main.class",
	InputFile:      "input.txt",
	OutputFile:     "output.txt",
	CompileCommand: "javac $SOURCE_FILE",
	RunCommand:     "java Main $INPUT_FILE $OUTPUT_FILE",
}

var Go = Language{
	Version:        "1.13",
	Name:           "Go 1.13",
	SourceFile:     "main.go",
	BinaryFile:     "main",
	InputFile:      "input.txt",
	OutputFile:     "output.txt",
	CompileCommand: "go build -o $BINARY_FILE $SOURCE_FILE",
	RunCommand:     "$BINARY_FILE $INPUT_FILE $OUTPUT_FILE",
}
