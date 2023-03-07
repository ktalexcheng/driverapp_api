package middleware

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Logs start and end of endpoint handlers
func LogHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO: Log request parameters
		log.Info().Msg(fmt.Sprintf("Processing %s %s", r.Method, r.URL))

		next.ServeHTTP(w, r)

		log.Info().Msg(fmt.Sprintf("Completed %s %s", r.Method, r.URL))
	})
}
