package main

import (
	"encoding/json"
	"os"

	model "github.com/Mirai3103/remote-compiler/internal/models"
	"gopkg.in/yaml.v3"
)

func ptr(s string) *string {
	return &s
}

func main() {
	cpp := model.Language{
		SourceFileExt:  ptr(".cpp"),
		BinaryFileExt:  ptr(".out"),
		CompileCommand: ptr("g++ $SourceFileName -o $BinaryFileName"),
		RunCommand:     ptr("$BinaryFileName"),
	}
	c := model.Language{
		SourceFileExt:  ptr(".c"),
		BinaryFileExt:  ptr(".out"),
		CompileCommand: ptr("gcc $SourceFileName -o $BinaryFileName"),
		RunCommand:     ptr("$BinaryFileName"),
	}
	python := model.Language{
		SourceFileExt:  ptr(".py"),
		BinaryFileExt:  nil,
		CompileCommand: nil,
		RunCommand:     ptr("python3 $SourceFileName"),
	}
	goLang := model.Language{
		SourceFileExt:  ptr(".go"),
		BinaryFileExt:  ptr(".out"),
		CompileCommand: ptr("go build -o $BinaryFileName $SourceFileName"),
		RunCommand:     ptr("$BinaryFileName"),
	}
	node := model.Language{
		SourceFileExt:  ptr(".js"),
		BinaryFileExt:  nil,
		CompileCommand: nil,
		RunCommand:     ptr("node $SourceFileName"),
	}
	shell := model.Language{
		SourceFileExt:  ptr(".sh"),
		BinaryFileExt:  nil,
		CompileCommand: nil,
		RunCommand:     ptr("/bin/bash $SourceFileName"),
	}

	listLang := []model.Language{cpp, c, python, goLang, node, shell}
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
