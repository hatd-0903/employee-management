package middleware

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"net/http"

	"employee-management/internal/models"
)

// BasicAuth enforces HTTP Basic Authentication against a fixed set of
// username/password pairs. /healthz is left open for container health checks.
func BasicAuth(users map[string]string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/healthz" {
				next.ServeHTTP(w, r)
				return
			}

			username, password, ok := r.BasicAuth()
			if !ok || !validCredentials(users, username, password) {
				unauthorized(w)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// validCredentials always hashes and compares, even for an unknown username,
// so response time doesn't leak which usernames exist or how long a
// password is (see review comment on this function).
func validCredentials(users map[string]string, username, password string) bool {
	expected, exists := users[username]

	got := sha256.Sum256([]byte(password))
	want := sha256.Sum256([]byte(expected))
	match := subtle.ConstantTimeCompare(got[:], want[:]) == 1

	return match && exists
}

func unauthorized(w http.ResponseWriter) {
	appErr := models.ErrUnauthorized("invalid credentials")
	w.Header().Set("WWW-Authenticate", `Basic realm="employee-management"`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	_ = json.NewEncoder(w).Encode(appErr)
}
