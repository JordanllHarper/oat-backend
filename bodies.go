package main

import "net/http"

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
