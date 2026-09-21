package middleware

import (
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

func validCredentials(users map[string]string, username, password string) bool {
	expected, exists := users[username]
	if !exists {
		return false
	}
	return subtle.ConstantTimeCompare([]byte(expected), []byte(password)) == 1
}

func unauthorized(w http.ResponseWriter) {
	appErr := models.ErrUnauthorized("invalid credentials")
	w.Header().Set("WWW-Authenticate", `Basic realm="employee-management"`)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.Code)
	_ = json.NewEncoder(w).Encode(appErr)
}
