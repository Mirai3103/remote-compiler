package model

type Submission struct {
	ID              *string             `json:"id"`
	Language        *Language           `json:"language"`
	Code            *string             `json:"code"`
	TimeLimitInMs   int                 `json:"timeLimitInMs"`
	MemoryLimitInKb int                 `json:"memoryLimitInKb"`
	TestCases       []TestCase          `json:"testCases"`
	Settings        *SubmissionSettings `json:"settings"`
}
