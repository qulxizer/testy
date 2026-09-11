package models

type SubmissionRequest struct {
	Exercise string            `json:"exercise"`
	Filename string            `json:"filename"`
	Code     string            `json:"code"`
	Main     *string           `json:"main,omitempty"`
	Files    map[string]string `json:"files,omitempty"`
}

type SubmissionResponse struct {
	Output string `json:"output"`
	Passed bool   `json:"passed"`
}
