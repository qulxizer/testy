package models

type SubmissionRequest struct {
	Exercise  string            `json:"exercise"`
	Filename  string            `json:"filename"`
	Code      string            `json:"code"`
	Main      *string           `json:"main"`
	Files     map[string]string `json:"files"`
	TestImage string            `json:"testImage"`
}

type SubmissionResponse struct {
	Output string `json:"output"`
	Passed bool   `json:"passed"`
}
