package main

import (
	"net/http"

	"github.com/google/uuid"
)

func handleGetContexts(contexts contextStore) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		return getAllContexts(contexts)
	}
}
func handleGetContextById(contexts contextStore) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		ctx, err := getById(contexts, r.PathValue(idKey))
		if err != nil {
			return nil, err
		}
		return statusOk{ctx}, nil
	}
}
func handlePostContext(contexts contextStore) HttpResponseHandler {
	type postCtx struct {
		Name string `json:"name"`
	}
	return func(r *http.Request) (HttpResponse, error) {
		postCtx, err := jsonDecode[postCtx](r.Body)
		if err != nil {
			return nil, err
		}
		context := newContext(postCtx.Name)
		if err := contexts.Add(context); err != nil {
			return nil, err
		}
		return getAllContexts(contexts)
	}
}
func handlePutContext(contexts contextStore) HttpResponseHandler {
	type putCtx struct {
		NewName string `json:"newName"`
	}
	return func(r *http.Request) (HttpResponse, error) {
		putCtx, err := jsonDecode[putCtx](r.Body)
		if err != nil {
			return nil, err
		}
		qId := r.PathValue(idKey)
		id, err := uuid.Parse(qId)
		if err != nil {
			return nil, malformedId{qId, err}
		}
		if err = contexts.EditName(id, putCtx.NewName); err != nil {
			return nil, err
		}
		return getAllContexts(contexts)
	}
}

func handleDeleteContext(contexts contextStore, tasks taskStore) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		qId := r.PathValue(idKey)
		id, err := uuid.Parse(qId)
		if err != nil {
			return nil, malformedId{qId, err}
		}
		if err = contexts.Delete(id); err != nil {
			return nil, err
		}
		if err = tasks.DeleteByContext(id); err != nil {
			return nil, err
		}
		return getAllContexts(contexts)
	}
}
