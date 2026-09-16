package middlewares

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"testy/db"
	"time"
)

type contextKey string

const SessionCookieName = "session_token"
const SessionDuration = 24 * time.Hour

const UserIDKey contextKey = "user_id"

func GetUserID(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(UserIDKey).(int64)
	return id, ok
}

func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(b), nil
}

func HashToken(token string) []byte {
	h := sha256.Sum256([]byte(token))
	return h[:]
}

func SetSessionCookie(w http.ResponseWriter, token string, expires time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     SessionCookieName,
		Value:    token,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode, // CSRF protection
	})
}

func StartSessionCleaner(db *db.DB, interval time.Duration) {
	log.Println("session cleaner up")
	go func() {
		ticker := time.NewTicker(interval)
		for range ticker.C {
			_, _ = db.Pool.Exec(context.TODO(), "DELETE FROM sessions WHERE expires_at < NOW()")
		}
	}()
}

func AuthMiddleware(db *db.DB) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			tokenHash := HashToken(cookie.Value)

			userID, err := db.GetUserIDBySession(r.Context(), tokenHash)
			if err != nil {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
