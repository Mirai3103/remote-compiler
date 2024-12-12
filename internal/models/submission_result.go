package internal

type SubmissionResult struct {
	SubmissionID string `json:"submissionId"`
	TestCaseID   string `json:"testCaseId"`
	Status       string `json:"status"`
	Stdout       string `json:"stdout"`
	Stderr       string `json:"stderr"`
}
