package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_Authenticate(t *testing.T) {

	const ak = "test-key"

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authenticate(w, r, ak) {
			w.WriteHeader(200)
			return
		}
	})

	t.Run("test same api-keys", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("x-api-key", "test-key")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != 200 {
			t.Error("matching keys failed")
		}
	})

	t.Run("test diff api-keys", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("x-api-key", "test-ke")

		rr := httptest.NewRecorder()

		handler.ServeHTTP(rr, req)

		if rr.Code != 401 {
			t.Error("diff keys passed")
		}
	})
}
