package model

import "github.com/google/uuid"

type TestCase struct {
	ID           *string `json:"id"`
	Input        *string `json:"input"`
	ExpectOutput *string `json:"expectOutput"`
	InputFile    *string `json:"inputFile"`
	OutputFile   *string `json:"outputFile"`
}

func (t *TestCase) GenerateInputFileName() string {
	newFileName := uuid.NewString() + ".in"
	t.InputFile = &newFileName
	return newFileName
}

func (t *TestCase) GenerateExpectOutputFileName() string {
	newFileName := uuid.NewString() + ".out"
	t.OutputFile = &newFileName
	return newFileName
}
