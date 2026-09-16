package server

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"testy/middlewares"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
}

func (s *Server) getJWTFromReboot(username, password string) (string, error) {
	rr, err := http.NewRequest(http.MethodPost, "https://learn.reboot01.com/api/auth/signin", nil)
	rr.SetBasicAuth(username, password)
	resp, err := s.client.Do(rr)
	if err != nil {
		return "", err
	}
	if resp.StatusCode == http.StatusUnauthorized {
		return "", err
	} else if resp.StatusCode != http.StatusOK {
		return "", err
	}

	scanner := bufio.NewScanner(resp.Body)
	if !scanner.Scan() {
		return "", fmt.Errorf("failed to read jwt")
	}

	token := scanner.Text()
	return token, nil
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1048576)
	username, password, ok := r.BasicAuth()
	if !ok || username == "" || password == "" {
		http.Error(w, "username and password are required", http.StatusBadRequest)
		return
	}

	token, err := s.getJWTFromReboot(username, password)
	if err != nil || token == "" {
		http.Error(w, "wrong username or password", http.StatusUnauthorized)
		return
	}

	user, err := s.db.CreateUser(r.Context(), username)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	rawToken, err := middlewares.GenerateToken()
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	tokenHash := middlewares.HashToken(rawToken)
	expiresAt := time.Now().Add(middlewares.SessionDuration)

	err = s.db.CreateSession(r.Context(), tokenHash, user.ID, expiresAt)
	if err != nil {
		log.Fatalln(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	middlewares.SetSessionCookie(w, rawToken, expiresAt)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(LoginResponse{
		ID:       user.ID,
		Username: user.Username,
	})
}
