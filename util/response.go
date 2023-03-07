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
func HTTPWriteJSONResponse(w http.ResponseWriter, status int, respBody *JSONResponse) error {
	response, err := json.Marshal(respBody)
	if err != nil {
		return err
	}

	w.WriteHeader(status)
	w.Write(response)
	return nil
}
