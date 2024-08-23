package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type Publisher interface {
	Len() (int, error)
	Close()
	Push(mess []byte) (int, error)
}

type enqueueHandler struct {
	key    string
	driver Publisher
}

func checkEquality(fullPath, value string) bool {
	splitList := strings.Split(fullPath, "/")
	if strings.ToLower(splitList[len(splitList)-1]) == value {
		return true
	}

	return false
}

func successHandler(w http.ResponseWriter, response []byte) {
	w.WriteHeader(http.StatusOK)

	var err error

	if len(response) != 0 {
		_, err = w.Write(response)
	}

	if err != nil {
		log.Println(err)
	}
}

func serverErrorHandler(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)

	if err != nil {
		_, er := w.Write([]byte(err.Error()))

		if er != nil {
			log.Println(err)
		}
	}
}

func notAuthorizedHandler(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	_, err := w.Write([]byte("401 status unauthorized"))
	if err != nil {
		log.Println(err)
	}
}

func tooLargeHandler(w http.ResponseWriter) {
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

func length(w http.ResponseWriter, driver Publisher) {
	lg, err := driver.Len()

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
	successHandler(w, nil)
}

func enqueue(w http.ResponseWriter, r *http.Request, driver Publisher) {
	bodyBytes, err := io.ReadAll(r.Body)

	if err != nil {
		serverErrorHandler(w, err)
	}

	lg, err := driver.Push(bodyBytes)

	if err != nil && err.Error() == "message exceeds acceptable message size" {
		tooLargeHandler(w)
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
		length(w, h.driver)
		return
	case r.Method == http.MethodPost && comparePaths("enqueue"):
		enqueue(w, r, h.driver)
		return
	default:
		w.WriteHeader(404)
		return
	}
}
