package middleware

import (
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

// PanicRecoveryMiddleware recovers from panics and returns a JSON error response.
func PanicRecoveryMiddleware(next http.Handler) http.Handler {
	logger := zap.NewExample().Sugar()
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Recovered from panic: %v", err)
				var errMessage = "Internal Server Error"
				if err, ok := err.(error); ok {
					errMessage = err.Error()
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				logger.Fatal(json.NewEncoder(w).Encode(ErrorResponse{
					Error: errMessage,
				}))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
