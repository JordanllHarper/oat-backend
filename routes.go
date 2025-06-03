package main

import (
	"net/http"
)

const idKey = "id"

func setupRoutes(
	mux *http.ServeMux,
	tasks taskStore,
	contexts contextStore,
) {
	handleRoute(mux, "GET /tasks/current", handleGetCurrentTask(tasks, contexts))
	handleRoute(mux, "POST /tasks/current", handlePostCurrentTask(tasks, contexts))
	handleRoute(mux, "PUT /tasks/current", handlePutCurrentTask(tasks))
	handleRoute(mux, "GET /tasks/{id}", handleGetTaskById(tasks))
	handleRoute(mux, "POST /tasks", handlePostTask(tasks, contexts))
	handleRoute(mux, "PUT /tasks/{id}", handlePostTask(tasks, contexts))
	handleRoute(mux, "PUT /complete", handleCompleteCurrentTask(tasks, contexts))

	handleRoute(mux, "GET /context", handleGetContexts(contexts))
	handleRoute(mux, "GET /context/{id}", handleGetContextById(contexts))
	handleRoute(mux, "POST /context", handlePostContext(contexts))
	handleRoute(mux, "PUT /context/{id}", handlePutContext(contexts))
	handleRoute(mux, "DELETE /context/{id}", handleDeleteContext(contexts))

}
