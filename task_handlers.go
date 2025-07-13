package main

import (
	"net/http"
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
		Notes    *string  `json:"notes"`
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
			ctx.Id,
			postTask.Title,
			postTask.Notes,
			postTask.Priority,
		)
		if err != nil {
			return nil, err
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
		Notes    *string  `json:"notes"`
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
			ctx.Id,
			postTask.Title,
			postTask.Notes,
			postTask.Priority,
		)
		if err != nil {
			return nil, err
		}
		if err = contexts.SetNewCurrentTask(ctx.Id, t.Id); err != nil {
			return nil, err
		}
		return statusCreated{t}, nil
	}
}
func handlePutCurrentTask(tasks taskStore, contexts contextStore) HttpResponseHandler {
	type putTask struct {
		ContextId *id       `json:"contextId"`
		Title     *string   `json:"title"`
		Notes     *string   `json:"notes"`
		Priority  *priority `json:"priority"`
	}
	return func(r *http.Request) (HttpResponse, error) {
		ctx, err := getCtxFromRq(contexts, r)
		if err != nil {
			return nil, err
		}
		currentTaskId := ctx.CurrentTaskId
		if currentTaskId == nil {
			return nil, noCurrentTask(ctx.Id)
		}
		currentTask, err := tasks.GetById(*currentTaskId)
		if err != nil {
			return nil, err
		}
		putTask, err := jsonDecode[putTask](r.Body)
		if err != nil {
			return nil, malformedBody{err}
		}

		if putTask.ContextId != nil {
			_, err := contexts.GetById(*putTask.ContextId)
			if err != nil {
				return nil, err
			}
			currentTask.ContextId = *putTask.ContextId
		}
		if putTask.Title != nil {
			currentTask.Title = *putTask.Title
		}
		if putTask.Notes != nil {
			currentTask.Notes = putTask.Notes
		}
		if putTask.Priority != nil {
			if !priorityValid(*putTask.Priority) {
				return nil, invalidPriority(*putTask.Priority)
			}
			currentTask.Priority = *putTask.Priority
		}
		t, err := tasks.ModifyTask(currentTask)
		if err != nil {
			return nil, err
		}
		return statusOk{t}, nil
	}
}
