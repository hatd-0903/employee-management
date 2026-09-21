package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"

	"employee-management/internal/models"
)

// Recover catches panics from downstream handlers and returns the standard
// JSON error envelope instead of letting the server crash the connection.
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Printf("panic recovered: %v\n%s", rec, debug.Stack())

				appErr := models.ErrInternal(fmt.Errorf("%v", rec))
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(appErr.Code)
				_ = json.NewEncoder(w).Encode(appErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
