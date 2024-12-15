package model

import "github.com/google/uuid"

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
	newFileName := uuid.NewString() + ".input"
	t.InputFile = &newFileName
	return newFileName
}

func (t *TestCase) GetExpectOutputFileName() string {
	if t.OutputFile != nil {
		return *t.OutputFile
	}
	newFileName := uuid.NewString() + ".output"
	t.OutputFile = &newFileName
	return newFileName
}
