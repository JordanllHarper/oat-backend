package main

import (
	"github.com/google/uuid"
)

const (
	One   priority = 1
	Two   priority = 2
	Three priority = 3
)

type (
	id = uuid.UUID

	priority int
	task     struct {
		Id        id       `json:"id"`
		ContextId id       `json:"contextId"`
		Title     string   `json:"title"`
		Notes     string   `json:"notes"`
		Priority  priority `json:"priority"`
	}

	context struct {
		Id            id     `json:"id"`
		Name          string `json:"name"`
		CurrentTaskId *id    `json:"currentTaskId"`
	}
)

func newTask(
	contextId id,
	title string,
	notes string,
	priority priority,
) task {
	return task{
		Id:        uuid.New(),
		ContextId: contextId,
		Title:     title,
		Notes:     notes,
		Priority:  priority,
	}
}

func newContext(name string) context {
	return context{
		Id:   uuid.New(),
		Name: name,
	}
}
