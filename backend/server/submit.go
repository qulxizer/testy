package server

import (
	"encoding/json"
	"log"
	"net/http"

	"testy/middlewares"
	"testy/models"
	"testy/runner"
)

func (s *Server) SubmitHandler(w http.ResponseWriter, r *http.Request) {
	userID, ok := middlewares.GetUserID(r.Context())
	if !ok {
		log.Fatalln("unauthorized")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	log.Printf("userID submit: %v\n", userID)
	var req models.SubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatalln(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	output, passed, err := runner.RunSubmission(r.Context(), req)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "Runner failed: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(models.SubmissionResponse{
		Output: output,
		Passed: passed,
	})
}
