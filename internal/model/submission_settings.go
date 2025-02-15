package model

// message SubmissionSettings {
// 	bool with_trim = 1;
// 	bool with_case_sensitive = 2;
// 	bool with_whitespace = 3;
//   }

type SubmissionSettings struct {
	WithTrim          bool `json:"with_trim"`
	WithCaseSensitive bool `json:"with_case_sensitive"`
	WithWhitespace    bool `json:"with_whitespace"`
}
