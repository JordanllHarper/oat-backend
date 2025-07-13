package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
)

func newServer(
	taskStore taskStore,
	contextStore contextStore,
	logger *log.Logger,
) http.Handler {
	mux := http.NewServeMux()
	setupRoutes(mux, taskStore, contextStore)
	var handler http.Handler = mux
	handler = logResponses(handler, logger)
	handler = logStores(handler, logger, taskStore, contextStore)
	return handler
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Server exited with error: %v\n", err)
	}
}

func run() error {
	testContextId := uuid.MustParse("cdf053cb-d7c7-45e3-b1e6-18f291690caa")
	testTaskId := uuid.MustParse("cbec13d5-36f7-4bac-a326-3a75d555a995")
	cs := contextStoreImpl{
		// sample data
		testContextId: {
			Id:            testContextId,
			Name:          "Test context",
			CurrentTaskId: &testTaskId,
		},
	}
	ts := &taskStoreImpl{
		task{
			Id:        testTaskId,
			ContextId: testContextId,
			Title:     "Test task",
			Priority:  One,
			Notes:     nil,
		},
	}
	srv := newServer(
		ts,
		cs,
		log.Default(),
	)

	port := "8080"
	fmt.Println("Listening on port", port)
	return http.ListenAndServe(":"+port, srv)
}
