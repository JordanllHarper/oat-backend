package main

import (
	"net/http"

	"github.com/google/uuid"
)

func handleGetTaskById(tasks taskStore) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		task, err := getById(tasks, r.PathValue(idKey))
		if err != nil {
			return nil, err
		}
		return statusOk{task}, nil
	}
}

func handleGetCurrentTask(
	tasks taskStore,
	contexts contextStore,
) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		ctx, err := getCtxFromRq(contexts, r)
		if err != nil {
			return nil, err
		}
		hasTask, task, err := getCurrentTask(ctx, tasks)
		if err != nil {
			return nil, err
		}
		if !hasTask {
			return statusNoContent{}, nil
		}
		return statusOk{task}, nil
	}
}

func handleCompleteCurrentTask(
	tasks taskStore,
	contexts contextStore,
) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		ctxId, err := getCtxFromRq(contexts, r)
		if err != nil {
			return nil, err
		}
		hasTask, next, err := completeAndGetNextTask(ctxId, tasks, contexts)
		if err != nil {
			return nil, err
		}
		if !hasTask {
			return statusNoContent{}, nil
		}
		return statusOk{next}, nil
	}
}
func handlePostTask(
	tasks taskStore,
	contexts contextStore,
) HttpResponseHandler {
	type postTask struct {
		Title    string   `json:"title"`
		Notes    string   `json:"notes"`
		Priority priority `json:"priority"`
	}
	return func(r *http.Request) (HttpResponse, error) {
		ctx, err := getCtxFromRq(contexts, r)
		if err != nil {
			return nil, err
		}
		postTask, err := jsonDecode[postTask](r.Body)
		if err != nil {
			return nil, malformedBody{}
		}
		t, err := addTask(
			tasks,
			uuid.New,
			ctx.Id,
			postTask.Title,
			postTask.Notes,
			postTask.Priority,
		)
		if err != nil {
			return nil, err
		}
		if ctx.CurrentTaskId == nil {
			err = contexts.SetNewCurrentTask(ctx.Id, &t.Id)
			if err != nil {
				return nil, err
			}
		}
		return statusCreated{t}, nil
	}
}
func handlePostCurrentTask(
	tasks taskStore,
	contexts contextStore,
) HttpResponseHandler {
	type postTask struct {
		Title    string   `json:"title"`
		Notes    string   `json:"notes"`
		Priority priority `json:"priority"`
	}
	return func(r *http.Request) (HttpResponse, error) {
		ctx, err := getCtxFromRq(contexts, r)
		if err != nil {
			return nil, err
		}
		postTask, err := jsonDecode[postTask](r.Body)
		if err != nil {
			return nil, malformedBody{err}
		}
		t, err := addTask(
			tasks,
			uuid.New,
			ctx.Id,
			postTask.Title,
			postTask.Notes,
			postTask.Priority,
		)
		if err != nil {
			return nil, err
		}
		if err = contexts.SetNewCurrentTask(ctx.Id, &t.Id); err != nil {
			return nil, err
		}
		return statusCreated{t}, nil
	}
}
func handlePutTaskById(tasks taskStore, contexts contextStore) HttpResponseHandler {
	type putTask struct {
		ContextId *id       `json:"contextId"`
		Title     *string   `json:"title"`
		Notes     *string   `json:"notes"`
		Priority  *priority `json:"priority"`
	}

	// naive - caller responsible for validation
	updateTaskFields := func(pt putTask, t task) task {
		if pt.ContextId != nil {
			t.ContextId = *pt.ContextId
		}
		if pt.Title != nil {
			t.Title = *pt.Title
		}
		if pt.Notes != nil {
			t.Notes = *pt.Notes
		}
		if pt.Priority != nil {
			t.Priority = *pt.Priority
		}
		return t
	}

	return func(r *http.Request) (HttpResponse, error) {
		task, err := getTaskFromRq(tasks, r)
		if err != nil {
			return nil, err
		}
		putTask, err := jsonDecode[putTask](r.Body)
		if err != nil {
			return nil, malformedBody{err}
		}
		currentCtx, err := contexts.GetById(task.ContextId)
		if err != nil {
			return nil, err
		}

		modifiedTask := updateTaskFields(putTask, task)
		if !priorityValid(modifiedTask.Priority) {
			return nil, invalidPriority(*putTask.Priority)
		}
		targetCtx, err := contexts.GetById(modifiedTask.ContextId)
		if err != nil {
			return nil, err
		}

		// set the task as current for target - might want to change this in the future
		if currentCtx.Id != targetCtx.Id {
			if err = contexts.SetNewCurrentTask(targetCtx.Id, &modifiedTask.Id); err != nil {
				return nil, err
			}
		}

		t, err := tasks.ModifyTask(modifiedTask)
		if err != nil {
			return nil, err
		}
		// update the current task of the current ctx
		if *currentCtx.CurrentTaskId == t.Id {
			if _, _, err := setNextTask(currentCtx, tasks, contexts); err != nil {
				return nil, err
			}
		}
		return statusNoContent{}, nil
	}
}
