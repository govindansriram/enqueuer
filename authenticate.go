package main

import (
	"net/http"
)

func authenticate(w http.ResponseWriter, r *http.Request, validKey string) bool {
	apiKey := r.Header.Get("x-api-key")

	if apiKey == validKey {
		return true
	}

	notAuthorizedHandler(w)
	return false
}
