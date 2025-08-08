package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	HttpResponse interface {
		Body() any
		StatusCode() int
	}

	bodyMultipleValues[T any] struct {
		Values []T `json:"values"`
	}
	statusOk        struct{ body any }
	statusCreated   struct{ body any }
	statusNoContent struct{}
)

func (stat statusOk) Body() any       { return stat.body }
func (stat statusOk) StatusCode() int { return http.StatusOK }

func (stat statusCreated) Body() any       { return stat.body }
func (stat statusCreated) StatusCode() int { return http.StatusCreated }

func (stat statusNoContent) Body() any       { return nil }
func (stat statusNoContent) StatusCode() int { return http.StatusNoContent }

type (
	HttpError interface {
		StatusCode() int
		error
	}
	internalServerError struct{ error }
	HttpResponseHandler func(r *http.Request) (HttpResponse, error)
)

func (err internalServerError) StatusCode() int { return http.StatusInternalServerError }

func mapErrToHttpError(err error) HttpError {
	var httpErr HttpError
	if errors.As(err, &httpErr) {
		return httpErr
	}
	return internalServerError{err}
}

func handleRoute(
	mux *http.ServeMux,
	pat string,
	handler HttpResponseHandler,
) {
	mux.HandleFunc(
		pat,
		func(w http.ResponseWriter, r *http.Request) {
			response, err := handler(r)
			if err != nil {
				httpErr := mapErrToHttpError(err)
				http.Error(w, httpErr.Error(), httpErr.StatusCode())
				return
			}
			code := response.StatusCode()
			w.WriteHeader(code)
			if response.Body() == nil {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			if err = jsonEncode(w, response.Body()); err != nil {
				log.SetPrefix("json")
				log.Println("An Error occurred while jsonEncoding", err)
				http.Error(w, fmt.Sprintf("An error occurred while encoding json: %v", err), http.StatusInternalServerError)
				return
			}
		},
	)
}

func jsonEncode(w http.ResponseWriter, body any) error { return json.NewEncoder(w).Encode(body) }

func jsonDecode[T any](body io.ReadCloser) (T, error) {
	var t T
	err := json.NewDecoder(body).Decode(&t)
	return t, err
}
