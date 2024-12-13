package main

import (
	"encoding/json"
	"os"

	"github.com/Mirai3103/remote-compiler/internal/model"
	"gopkg.in/yaml.v3"
)

func ptr(s string) *string {
	return &s
}

func main() {
	cpp := model.Language{
		Version:        ptr("g++ 12.2.0"),
		Name:           ptr("C++"),
		SourceFileExt:  ptr(".cpp"),
		BinaryFileExt:  ptr(".out"),
		CompileCommand: ptr("g++ $SourceFileName -o $BinaryFileName"),
		RunCommand:     ptr("$BinaryFileName"),
	}
	c := model.Language{
		Version:        ptr("gcc 12.2.0"),
		Name:           ptr("C"),
		SourceFileExt:  ptr(".c"),
		BinaryFileExt:  ptr(".out"),
		CompileCommand: ptr("gcc $SourceFileName -o $BinaryFileName"),
		RunCommand:     ptr("$BinaryFileName"),
	}
	python := model.Language{
		Version:        ptr("Python 3.11.2"),
		Name:           ptr("Python3"),
		SourceFileExt:  ptr(".py"),
		BinaryFileExt:  nil,
		CompileCommand: nil,
		RunCommand:     ptr("python3 $SourceFileName"),
	}
	goLang := model.Language{
		Version:        ptr("Go 1.23.4"),
		Name:           ptr("Go"),
		SourceFileExt:  ptr(".go"),
		BinaryFileExt:  ptr(".out"),
		CompileCommand: ptr("go build -o $BinaryFileName $SourceFileName"),
		RunCommand:     ptr("$BinaryFileName"),
	}
	node := model.Language{
		Version:        ptr("Node.js 16.13.0"),
		Name:           ptr("Node.js"),
		SourceFileExt:  ptr(".js"),
		BinaryFileExt:  nil,
		CompileCommand: nil,
		RunCommand:     ptr("node $SourceFileName"),
	}

	listLang := []model.Language{cpp, c, python, goLang, node}
	json, _ := json.MarshalIndent(listLang, "", "  ")
	err := os.WriteFile("languages.json", json, 0644)
	if err != nil {
		panic(err)
	}
	yaml, _ := yaml.Marshal(listLang)
	err = os.WriteFile("languages.yaml", yaml, 0644)
	if err != nil {
		panic(err)
	}

}
