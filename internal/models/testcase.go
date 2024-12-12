package internal

type TestCase struct {
	ID           string `json:"id"`
	Input        string `json:"input"`
	ExpectOutput string `json:"expectOutput"`
}
