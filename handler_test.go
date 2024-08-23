package main

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_checkEquality(t *testing.T) {

	const url = "localhost:3000/test/my/api"

	t.Run("invalid method", func(t *testing.T) {
		value := "api"

		if !checkEquality(url, value) {
			t.Error("failed to detect valid method")
		}
	})

	t.Run("invalid method", func(t *testing.T) {
		value := "ppi"

		if checkEquality(url, value) {
			t.Error("failed to detect invalid method")
		}
	})
}

func Test_successHandler(t *testing.T) {

	t.Run("success no message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			successHandler(w, nil)
		})

		handler.ServeHTTP(rr, req)

		if rr.Code != 200 {
			t.Error("matching keys failed")
		}

		if data, err := io.ReadAll(rr.Body); err == nil {
			if len(data) != 0 {
				t.Error("got message body")
			}
		} else {
			t.Error(err)
		}
	})

	t.Run("success w/ message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			successHandler(w, []byte("ayyyyyyye bruh"))
		})

		handler.ServeHTTP(rr, req)

		if rr.Code != 200 {
			t.Error("matching keys failed")
		}

		if data, err := io.ReadAll(rr.Body); err == nil {
			if !bytes.Equal(data, []byte("ayyyyyyye bruh")) {
				t.Error("got invalid message")
			}
		} else {
			t.Error(err)
		}
	})
}

func Test_serverErrorHandler(t *testing.T) {

	t.Run("error no message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serverErrorHandler(w, nil)
		})

		handler.ServeHTTP(rr, req)

		if rr.Code != 500 {
			t.Error("matching keys failed")
		}

		if data, err := io.ReadAll(rr.Body); err == nil {
			if len(data) != 0 {
				t.Error("got message body")
			}
		} else {
			t.Error(err)
		}
	})

	t.Run("error no message", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rr := httptest.NewRecorder()

		handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serverErrorHandler(w, errors.New("test error"))
		})

		handler.ServeHTTP(rr, req)

		if rr.Code != 500 {
			t.Error("matching keys failed")
		}

		if data, err := io.ReadAll(rr.Body); err == nil {
			if !bytes.Equal(data, []byte("test error")) {
				t.Error("got different message body")
			}
		} else {
			t.Error(err)
		}
	})
}
