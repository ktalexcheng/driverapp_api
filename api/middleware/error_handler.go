package middleware

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

// Wrapper function to accept error returns of service handlers
func ErrorHandler(h Handler) http.HandlerFunc {
	logErr := func(e error) { log.Error().Stack().Err(errors.Wrap(e, "error")).Msg("") }

	return func(w http.ResponseWriter, r *http.Request) {
		// Execute handler
		err := h(w, r)

		if err != nil {
			// Log error
			logErr(err)

			response := map[string]interface{}{
				"error": err.Error(),
				// "stack": err.StackTrace(),
			}
			jsonResponse, _err := json.Marshal(response)
			if _err != nil {
				logErr(err)
				return
			}

			// Write response
			w.WriteHeader(http.StatusInternalServerError) // Call WriteHeader() first to prevent superfluous call
			w.Header().Set("Content-Type", "application/json")
			_, _err = w.Write(jsonResponse)
			if _err != nil {
				logErr(err)
				return
			}
		}
	}
}
