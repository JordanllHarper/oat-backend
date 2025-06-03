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

func handleGetCurrentTask(tasks taskStore, contexts contextStore) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		ctx, err := getCtx(contexts, r)
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
		ctxId, err := getCtx(contexts, r)
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
func handlePostTask(tasks taskStore, contexts contextStore) HttpResponseHandler {
	type postTask struct {
		Title    string   `json:"title"`
		Notes    string   `json:"notes"`
		Priority priority `json:"priority"`
	}
	return func(r *http.Request) (HttpResponse, error) {
		ctx, err := getCtx(contexts, r)
		if err != nil {
			return nil, err
		}
		postTask, err := jsonDecode[postTask](r.Body)
		if err != nil {
			return nil, malformedBody{}
		}
		newTask := newTask(
			ctx.Id,
			postTask.Title,
			postTask.Notes,
			postTask.Priority,
		)
		if err := tasks.InsertTask(
			newTask,
		); err != nil {
			return nil, err
		}
		return statusCreated{newTask}, nil
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
		ctx, err := getCtx(contexts, r)
		if err != nil {
			return nil, err
		}
		postTask, err := jsonDecode[postTask](r.Body)
		if err != nil {
			return nil, malformedBody{err}
		}
		newTask := newTask(
			ctx.Id,
			postTask.Title,
			postTask.Notes,
			postTask.Priority,
		)
		if err = tasks.InsertTask(newTask); err != nil {
			return nil, err
		}
		if err = contexts.SetNewCurrentTask(ctx.Id, newTask.Id); err != nil {
			return nil, err
		}
		return statusCreated{newTask}, nil
	}
}
func handlePutCurrentTask(tasks taskStore) HttpResponseHandler {
	type putTask struct {
		ContextId *id       `json:"contextId"`
		Title     *string   `json:"title"`
		Notes     *string   `json:"notes"`
		Priority  *priority `json:"priority"`
	}
	return func(r *http.Request) (HttpResponse, error) {
		currentTaskToModify, err := getById(tasks, r.PathValue(idKey))
		if err != nil {
			return nil, err
		}
		putTask, err := jsonDecode[putTask](r.Body)
		if err != nil {
			return nil, malformedBody{err}
		}
		if putTask.ContextId != nil {
			currentTaskToModify.ContextId = *putTask.ContextId
		}
		if putTask.Title != nil {
			currentTaskToModify.Title = *putTask.Title
		}
		if putTask.Notes != nil {
			currentTaskToModify.Notes = *putTask.Notes
		}
		if putTask.Priority != nil {
			currentTaskToModify.Priority = *putTask.Priority
		}
		newTask, err := tasks.ModifyTask(currentTaskToModify)
		if err != nil {
			return nil, err
		}
		return statusOk{newTask}, nil
	}
}
