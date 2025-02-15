package model

type Submission struct {
	ID          *string             `json:"id"`
	Language    *Language           `json:"language"`
	Code        *string             `json:"code"`
	TimeLimit   int                 `json:"timeLimit"`
	MemoryLimit int                 `json:"memoryLimit"`
	TestCases   []TestCase          `json:"testCases"`
	Settings    *SubmissionSettings `json:"settings"`
}
