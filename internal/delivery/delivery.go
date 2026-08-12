package delivery

import "net/http"

var (
	Forbidden = http.StatusForbidden
	Unauthorized = http.StatusUnauthorized
	InternalServerError = http.StatusInternalServerError
)