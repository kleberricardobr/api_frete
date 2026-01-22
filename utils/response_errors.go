package utils

import (
	"encoding/json"
	"net/http"
)

func AddError(w http.ResponseWriter, error string, statusCode uint) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(int(statusCode))
	json.NewEncoder(w).Encode(map[string]string{
		"error": error,
	})
}
