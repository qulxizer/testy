package handler

import (
	"encoding/json"
	"net/http"

	"testy/models"
	"testy/runner"
)

func SubmitHandler(w http.ResponseWriter, r *http.Request) {
	var req models.SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, passed, err := runner.RunSubmission(r.Context(), req)
	if err != nil {
		http.Error(w, "Runner failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(models.SubmissionResponse{
		Output: output,
		Passed: passed,
	})
}
