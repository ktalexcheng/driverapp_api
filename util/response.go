package util

import (
	"encoding/json"
	"net/http"
)

// A basic JSON response with a simple message field
type JSONResponse struct {
	Message string `json:"message"`
}

// Utility to write status code to response header and JSONResponse to body
func HTTPWriteJSONBody(w http.ResponseWriter, status int, respBody interface{}) error {
	w.WriteHeader(status)

	if respBody != nil {
		response, err := json.Marshal(respBody)
		if err != nil {
			return err
		}

		_, err = w.Write(response)
		if err != nil {
			return err
		}
	}

	return nil
}
