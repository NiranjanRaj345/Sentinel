package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/NiranjanRaj345/sentinel/services/node-agent/internal/logger"
)

func Recovery(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			defer func() {
				if err := recover(); err != nil {

					log.Error(
						"Panic recovered: %v\n%s",
						err,
						debug.Stack(),
					)

					http.Error(
						w,
						http.StatusText(http.StatusInternalServerError),
						http.StatusInternalServerError,
					)
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
