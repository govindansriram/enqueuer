package main

import (
	"bytes"
	"encoding/json"
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

func Test_notAuthorizedHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		notAuthorizedHandler(w)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != 401 {
		t.Error("matching keys failed")
	}
}

func Test_tooLargeHandler(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tooLargeHandler(w)
	})

	handler.ServeHTTP(rr, req)

	if rr.Code != 413 {
		t.Error("matching keys failed")
	}
}

func Test_getLengthByteMap(t *testing.T) {
	byts, err := getLengthByteMap(10)

	if err != nil {
		t.Fatal(err)
	}

	dataMap := &map[string]int{}

	err = json.Unmarshal(byts, dataMap)

	if err != nil {
		t.Fatal(err)
	}

	if val, ok := (*dataMap)["length"]; !ok || val != 10 {
		t.Fatal("map is incorrect")
	}
}

func Test_ping(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ping", nil)
	req.Header.Set("x-api-key", "test-key")
	rr := httptest.NewRecorder()

	handler := enqueueHandler{
		key: "test-key",
	}

	handler.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Error("ping failed")
	}
}

type TestLenDriver struct {
	errs bool
	val  int
}

func (t TestLenDriver) Len() (int, error) {
	if t.errs {
		return -1, errors.New("test error")
	} else {
		return t.val, nil
	}
}

func (t TestLenDriver) Close() {
	return
}

func (t TestLenDriver) Push(mess []byte) (int, error) {
	return 100, nil
}

func Test_length(t *testing.T) {

	t.Run("test valid len", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/len", nil)
		req.Header.Set("x-api-key", "test-key")
		rr := httptest.NewRecorder()

		handler := enqueueHandler{
			key:    "test-key",
			driver: TestLenDriver{errs: false, val: 1000},
		}

		handler.ServeHTTP(rr, req)

		if rr.Code != 200 {
			t.Error("len failed")
		}

		body, err := io.ReadAll(rr.Body)

		if err != nil {
			t.Fatal(err)
		}

		var lenMap map[string]int

		err = json.Unmarshal(body, &lenMap)

		if err != nil {
			t.Fatal(err)
		}

		if val, ok := (lenMap)["length"]; !ok || val != 1000 {
			t.Fatal("map is incorrect")
		}
	})

	t.Run("test invalid len", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/len", nil)
		req.Header.Set("x-api-key", "test-key")
		rr := httptest.NewRecorder()

		handler := enqueueHandler{
			key:    "test-key",
			driver: TestLenDriver{errs: true},
		}

		handler.ServeHTTP(rr, req)

		if rr.Code != 500 {
			t.Error("received invalid code")
		}
	})
}

type TestEnqueueDriver struct {
	errs bool
	val  int
}

func (t TestEnqueueDriver) Len() (int, error) {
	return 1000, nil
}

func (t TestEnqueueDriver) Close() {
	return
}

func (t TestEnqueueDriver) Push(mess []byte) (int, error) {
	if t.errs {
		return -1, errors.New("message exceeds acceptable message size")
	} else {
		return t.val, nil
	}
}

func Test_enqueue(t *testing.T) {

	t.Run("test valid enqueue", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/enqueue", nil)
		req.Header.Set("x-api-key", "test-key")
		rr := httptest.NewRecorder()

		handler := enqueueHandler{
			key:    "test-key",
			driver: TestEnqueueDriver{errs: false, val: 1000},
		}

		handler.ServeHTTP(rr, req)

		if rr.Code != 200 {
			t.Error("enqueue failed")
		}

		body, err := io.ReadAll(rr.Body)

		if err != nil {
			t.Fatal(err)
		}

		var lenMap map[string]int

		err = json.Unmarshal(body, &lenMap)

		if err != nil {
			t.Fatal(err)
		}

		if val, ok := (lenMap)["length"]; !ok || val != 1000 {
			t.Fatal("map is incorrect")
		}
	})

	t.Run("test invalid len", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/enqueue", nil)
		req.Header.Set("x-api-key", "test-key")
		rr := httptest.NewRecorder()

		handler := enqueueHandler{
			key:    "test-key",
			driver: TestEnqueueDriver{errs: true},
		}

		handler.ServeHTTP(rr, req)

		if rr.Code != 413 {
			t.Error("received invalid code")
		}
	})
}
