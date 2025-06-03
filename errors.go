package main

import (
	"fmt"
	"net/http"
)

type (
	HttpError interface {
		StatusCode() int
		error
	}

	internalServerError struct{ error }
	idNotFound          id
	idAlreadyExists     id
	malformedId         struct {
		id string
		error
	}
	malformedBody          struct{ error }
	noMoreTasks            struct{}
	unsupportedRoute       string
	noContextProvided      struct{}
	couldntFindCurrentTask struct {
		taskId, contextId id
	}
	noCurrentTask id
)

func (err internalServerError) StatusCode() int { return http.StatusInternalServerError }

func (err idNotFound) StatusCode() int { return http.StatusBadRequest }
func (err idNotFound) Error() string   { return fmt.Sprintf("Id %s not found", id(err)) }

func (err idAlreadyExists) StatusCode() int { return http.StatusBadRequest }
func (err idAlreadyExists) Error() string   { return fmt.Sprintf("Id %s already exists", id(err)) }

func (err malformedId) StatusCode() int { return http.StatusBadRequest }
func (err malformedId) Error() string {
	return fmt.Sprintf("Id %s is malformed due to: %s ", err.id, err.error)
}

func (err noMoreTasks) StatusCode() int { return http.StatusBadRequest }
func (err noMoreTasks) Error() string   { return fmt.Sprintf("No more tasks") }

func (err malformedBody) StatusCode() int { return http.StatusBadRequest }
func (err malformedBody) Error() string   { return fmt.Sprintf("Malformed body: %s", err.error.Error()) }

func (err unsupportedRoute) StatusCode() int { return http.StatusBadRequest }
func (err unsupportedRoute) Error() string   { return fmt.Sprintf("Unsupported route %s", string(err)) }

func (err noContextProvided) StatusCode() int { return http.StatusBadRequest }
func (err noContextProvided) Error() string   { return fmt.Sprintf("No context id provided") }

func (err couldntFindCurrentTask) StatusCode() int { return http.StatusInternalServerError }
func (err couldntFindCurrentTask) Error() string {
	return fmt.Sprintf("Registered current task %s in context %s couldn't be found", id(err.taskId), id(err.contextId))
}

func (err noCurrentTask) StatusCode() int { return http.StatusBadRequest }
func (err noCurrentTask) Error() string   { return fmt.Sprintf("No current task for: %s", id(err)) }
