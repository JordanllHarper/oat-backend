package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func logResponses(h http.Handler, logger *log.Logger) http.Handler {
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

func newServer(
	taskStore taskStore,
	contextStore contextStore,
	logger *log.Logger,
) http.Handler {
	mux := http.NewServeMux()
	setupRoutes(mux, taskStore, contextStore)
	var handler http.Handler = mux
	handler = logResponses(handler, logger)
	// handler = logStores(handler, logger, taskStore, contextStore)
	return handler
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Server exited with error: %v\n", err)
	}
}

func run() error {
	cs := contextStoreImpl{
		uuid.MustParse("cdf053cb-d7c7-45e3-b1e6-18f291690caa"): {
			Id:            uuid.MustParse("cdf053cb-d7c7-45e3-b1e6-18f291690caa"),
			Name:          "Test context",
			CurrentTaskId: nil,
		},
	}
	ts := &taskStoreImpl{}
	srv := newServer(
		ts,
		cs,
		log.Default(),
	)

	fmt.Println("Listening on port 8080")
	return http.ListenAndServe(":8080", srv)
}
