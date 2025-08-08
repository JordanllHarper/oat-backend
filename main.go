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
	handler = logRequests(handler, logger)
	handler = logStores(handler, logger, taskStore, contextStore)
	return handler
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Server exited with error: %v\n", err)
	}
}

func run() error {
	homeContextId := uuid.MustParse("cdf053cb-d7c7-45e3-b1e6-18f291690caa")
	workContextId := uuid.MustParse("039c80f0-42be-41d2-812b-bd323356a892")
	testTaskId := uuid.MustParse("cbec13d5-36f7-4bac-a326-3a75d555a995")
	// sample data
	cs := contextStoreImpl{
		homeContextId: {
			Id:            homeContextId,
			Name:          "Home",
			CurrentTaskId: &testTaskId,
		},

		workContextId: {
			Id:            workContextId,
			Name:          "Work",
			CurrentTaskId: nil,
		},
	}
	ts := &taskStoreImpl{
		task{
			Id:        testTaskId,
			ContextId: homeContextId,
			Title:     "Test task",
			Priority:  One,
			Notes:     "",
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
