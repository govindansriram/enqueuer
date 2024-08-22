package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type enqueueHandler struct {
	key string
}

func checkEquality(fullPath, value string) bool {
	splitList := strings.Split(fullPath, "/")
	if strings.ToLower(splitList[len(splitList)-1]) == value {
		return true
	}

	return false
}

func successHandler(w http.ResponseWriter, response ...[]byte) {
	w.WriteHeader(http.StatusOK)

	var err error

	if len(response) > 0 {
		_, err = w.Write(response[0])
	} else {
		_, err = w.Write([]byte("200 ok"))
	}

	if err != nil {
		log.Println(err)
	}
}

func serverErrorHandler(w http.ResponseWriter, err ...error) {
	w.WriteHeader(http.StatusInternalServerError)

	if len(err) == 0 {
		_, er := w.Write([]byte(err[0].Error()))

		if er != nil {
			log.Println(err)
		}
	}

}

func NotAuthorizedHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte("401 status unauthorized"))

	if err != nil {
		log.Println(err)
	}
}

func TooLargeHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusRequestEntityTooLarge)
	_, err := w.Write([]byte("413 queue cannot store a message this large"))

	if err != nil {
		log.Println(err)
	}
}

func getLengthByteMap(lg int) ([]byte, error) {
	m := make(map[string]int, 1)
	m["length"] = lg

	return json.Marshal(m)
}

func length(w http.ResponseWriter) {
	d := getDriver()
	lg, err := d.Len()

	if err != nil {
		serverErrorHandler(w, err)
	}

	byteMap, err := getLengthByteMap(lg)

	if err != nil {
		serverErrorHandler(w, err)
	}

	successHandler(w, byteMap)
}

func ping(w http.ResponseWriter) {
	successHandler(w)
}

func enqueue(w http.ResponseWriter, r *http.Request) {
	d := getDriver()

	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		serverErrorHandler(w, err)
	}

	lg, err := d.Push(bodyBytes)

	if err != nil && err.Error() == "message exceeds acceptable message size" {
		TooLargeHandler(w)
		return
	}

	if err != nil {
		serverErrorHandler(w, err)
	}

	byteMap, err := getLengthByteMap(lg)

	if err != nil {
		serverErrorHandler(w, err)
	}

	successHandler(w, byteMap)
}

func (h enqueueHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !authenticate(w, r, h.key) {
		return
	}

	comparePaths := func(desired string) bool {
		return checkEquality(r.URL.Path, desired)
	}

	switch {
	case r.Method == http.MethodGet && comparePaths("ping"):
		ping(w)
		return
	case r.Method == http.MethodGet && comparePaths("len"):
		length(w)
		return
	case r.Method == http.MethodPost && comparePaths("enqueue"):
		enqueue(w, r)
		return
	default:
		return
	}
}
