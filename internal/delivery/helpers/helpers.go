package helpers

import (
	"encoding/json"
	"errors"
	"net/http"
)

var (
	Forbidden = http.StatusForbidden
	Unauthorized = http.StatusUnauthorized
	InternalServerError = http.StatusInternalServerError
	BadRequest = http.StatusBadRequest
	OK = http.StatusOK
	NotFound = http.StatusNotFound

	ErrNotFound = errors.New("User ID not found")
	ErrQueryEmpty = errors.New("Err: Query Empty")
	ErrCampaignNotActive = errors.New("campaign is not active")
	ErrAuth = errors.New("refresh token missing")
)

func RespondWithError(w http.ResponseWriter, code int, err error) {
	errMsg := "unknown error"
	if err != nil {
		errMsg = err.Error()
	}
	Respond(w, code, map[string]string{"error": errMsg})
}

func Respond(w http.ResponseWriter, code int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if message != nil {
		json.NewEncoder(w).Encode(message)
	}
}