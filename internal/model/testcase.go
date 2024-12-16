package model

import snowflakeid "github.com/Mirai3103/remote-compiler/pkg/snowflake_id"

type TestCase struct {
	ID           *string `json:"id"`
	Input        *string `json:"input"`
	ExpectOutput *string `json:"expectOutput"`
	InputFile    *string `json:"inputFile"`
	OutputFile   *string `json:"outputFile"`
}

func (t *TestCase) GetInputFileName() string {
	if t.InputFile != nil {
		return *t.InputFile
	}
	newFileName := snowflakeid.NewString() + ".input"
	t.InputFile = &newFileName
	return newFileName
}

func (t *TestCase) GetExpectOutputFileName() string {
	if t.OutputFile != nil {
		return *t.OutputFile
	}
	newFileName := snowflakeid.NewString() + ".output"
	t.OutputFile = &newFileName
	return newFileName
}
