package middleware

import (
	"net/http"
	"time"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

type responseRecorder struct {
	http.ResponseWriter
	status int
}

func (r *responseRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func Logging(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			recorder := &responseRecorder{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			start := time.Now()

			next.ServeHTTP(recorder, r)

			duration := time.Since(start)

			requestID := RequestIDFromContext(r.Context())

			log.Info(
				"%s %s %d %s req=%s",
				r.Method,
				r.URL.Path,
				recorder.status,
				duration.Round(time.Microsecond),
				requestID,
			)
		})
	}
}
