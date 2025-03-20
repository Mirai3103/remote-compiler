package model

type SubmissionResult struct {
	SubmissionID    *string `json:"submissionId"`
	TestCaseID      *string `json:"testCaseId"`
	Status          *string `json:"status"`
	Stdout          *string `json:"stdout"`
	MemoryUsageInKb float64 `json:"memoryUsageInKb"`
	TimeUsageInMs   float64 `json:"timeUsageInMs"`
}
