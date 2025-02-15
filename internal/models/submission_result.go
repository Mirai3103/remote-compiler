package model

type SubmissionResult struct {
	SubmissionID *string `json:"submissionId"`
	TestCaseID   *string `json:"testCaseId"`
	Status       *string `json:"status"`
	Stdout       *string `json:"stdout"`
	MemoryUsage  float64 `json:"memoryUsage"`
	TimeUsage    float64 `json:"timeUsage"`
}
