package httphandlers

import (
	"encoding/json"
	"net/http"
)

var errSomethingWentWrongMsg = "an error occured, please try again later"

func ErrorResponse(rw http.ResponseWriter, status, message string, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(rw).Encode(map[string]string{
		"status":  status,
		"message": message,
	})
}
