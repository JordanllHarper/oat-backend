package main

import (
	"log"
	"net/http"
)

func logRequests(h http.Handler, logger *log.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			logger.Println("Received request: ", r.Method, r.URL)
			h.ServeHTTP(w, r)
		},
	)
}

func logStores(
	h http.Handler,
	logger *log.Logger,
	tasks taskStore,
	contexts contextStore,
) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(w, r)
			logger.Println(tasks)
			logger.Println(contexts)
		},
	)
}
