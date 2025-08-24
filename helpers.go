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

func setNextTask(
	ctx context,
	tasks taskGetter,
	contexts contextCurrentTaskSetter,
) (hasNextTask bool, t task, err error) {
	allTasks, err := tasks.AllForContext(ctx.Id)
	if err != nil {
		return false, task{}, err
	}
	if len(allTasks) == 0 {
		// no more tasks, set current task to nil
		return false, task{}, contexts.SetNewCurrentTask(ctx.Id, nil)
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
	tasks taskGetterDeleter,
	contexts contextCurrentTaskSetter,
) (hasNextTask bool, t task, err error) {
	if ctx.CurrentTaskId == nil {
		return false, task{}, nil
	}
	if err = tasks.RemoveTask(*ctx.CurrentTaskId); err != nil {
		return false, task{}, err
	}
	if err = contexts.SetNewCurrentTask(ctx.Id, nil); err != nil {
		return false, task{}, err
	}
	return setNextTask(ctx, tasks, contexts)
}

func priorityValid(priority priority) bool { return 1 <= priority && priority < 4 }

func addTask(
	tasks taskStore,
	idGenFunc func() id,
	ctxId id,
	title string,
	notes string,
	p priority,
) (task, error) {
	if !priorityValid(p) {
		return task{}, invalidPriority(p)
	}
	if strings.TrimSpace(title) == "" {
		return task{}, noTitle{}
	}
	notes = strings.TrimSpace(notes)
	t := newTask(
		idGenFunc,
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

func getAllContexts(contexts contextGetter) (HttpResponse, error) {
	all, err := contexts.All()
	if err != nil {
		return nil, err
	}
	return statusOk{bodyMultipleValues[context]{all}}, nil
}
