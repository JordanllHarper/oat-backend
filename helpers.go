package main

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/google/uuid"
)

type IDable[T any] interface {
	GetById(id id) (T, error)
}

func getById[T any](g IDable[T], idstr string) (T, error) {
	var t T
	id, err := uuid.Parse(idstr)
	if err != nil {
		return t, malformedId{idstr, err}
	}
	return g.GetById(id)
}

func getCtxFromRq(contexts contextStore, r *http.Request) (context, error) {
	qId := r.PathValue(idKey)
	if strings.TrimSpace(qId) == "" {
		return context{}, noContextProvided{}
	}
	id, err := uuid.Parse(qId)
	if err != nil {
		return context{}, malformedId{qId, err}
	}
	return contexts.GetById(id)
}

func getTaskFromRq(tasks taskStore, r *http.Request) (task, error) {
	qId := r.PathValue(idKey)
	if strings.TrimSpace(qId) == "" {
		return task{}, noContextProvided{}
	}
	id, err := uuid.Parse(qId)
	if err != nil {
		return task{}, malformedId{qId, err}
	}
	return tasks.GetById(id)
}

func getCurrentTask(
	ctx context,
	tasks taskStore,
) (hasTask bool, t task, err error) {
	if ctx.CurrentTaskId == nil {
		return false, task{}, nil
	}
	t, err = tasks.GetById(*ctx.CurrentTaskId)
	if err != nil {
		return false, task{}, fmt.Errorf("Getting current task for context id %s: %w", ctx.Id, err)
	}
	return true, t, nil
}

func getNextTask(
	ctx context,
	tasks taskStore,
	contexts contextStore,
) (hasNextTask bool, t task, err error) {
	allTasks, err := tasks.All()
	if err != nil {
		return false, task{}, err
	}
	if len(allTasks) == 0 {
		return false, task{}, nil
	}
	slices.SortFunc(allTasks, func(a, b task) int {
		return cmp.Compare(a.Priority, b.Priority)
	})
	highPriority := allTasks[0]
	if err = contexts.SetNewCurrentTask(ctx.Id, &highPriority.Id); err != nil {
		return false, task{}, err
	}
	return true, highPriority, nil
}

func completeAndGetNextTask(
	ctx context,
	tasks taskStore,
	contexts contextStore,
) (hasNextTask bool, t task, err error) {
	if err = tasks.RemoveTask(*ctx.CurrentTaskId); err != nil {
		return false, task{}, err
	}
	if err = contexts.SetNewCurrentTask(ctx.Id, nil); err != nil {
		return false, task{}, err
	}
	return getNextTask(ctx, tasks, contexts)
}

func priorityValid(priority priority) bool { return 1 <= priority && priority < 4 }

func addTask(
	tasks taskStore,
	ctxId id,
	title string,
	notes *string,
	p priority,
) (task, error) {
	if !priorityValid(p) {
		return task{}, invalidPriority(p)
	}
	if strings.TrimSpace(title) == "" {
		return task{}, noTitle{}
	}
	if notes != nil {
		// if a user adds newlines intentionally we don't want to remove them
		*notes = trimPreserveNewline(*notes)
		if *notes == "" {
			notes = nil
		}
	}
	t := newTask(
		ctxId,
		title,
		notes,
		p,
	)
	if err := tasks.InsertTask(t); err != nil {
		return task{}, err
	}
	return t, nil
}
