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
	// get the current task for the context
	handleRoute(mux, "GET /tasks/current/{id}", handleGetCurrentTask(tasks, contexts))
	// get a task by id
	handleRoute(mux, "GET /tasks/{id}", handleGetTaskById(tasks))
	// push a task to the top of the context
	handleRoute(mux, "POST /tasks/current/{id}", handlePostCurrentTask(tasks, contexts))
	// push a task to the context to be sorted
	handleRoute(mux, "POST /tasks/{id}", handlePostTask(tasks, contexts))
	// edit the current task
	handleRoute(mux, "PUT /tasks/current/{id}", handlePutCurrentTask(tasks, contexts))
	// edit a task by id
	handleRoute(mux, "PUT /tasks/{id}", handlePutTaskById(tasks, contexts))
	// complete a task
	handleRoute(mux, "PUT /complete/{id}", handleCompleteCurrentTask(tasks, contexts))

	// get all contexts
	handleRoute(mux, "GET /context", handleGetContexts(contexts))
	// get context by id
	handleRoute(mux, "GET /context/{id}", handleGetContextById(contexts))
	// create a new context
	handleRoute(mux, "POST /context", handlePostContext(contexts))
	// edit a context by id
	handleRoute(mux, "PUT /context/{id}", handlePutContext(contexts))
	// delete a context and associated tasks
	handleRoute(mux, "DELETE /context/{id}", handleDeleteContext(contexts, tasks))

}
