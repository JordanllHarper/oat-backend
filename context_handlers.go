package main

import (
	"net/http"

	"github.com/google/uuid"
)

func handleGetContexts(contexts contextStore) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		contexts, err := contexts.All()
		if err != nil {
			return nil, err
		}
		return statusOk{bodyMultipleValues[context]{contexts}}, nil
	}
}
func handleGetContextById(contexts contextStore) HttpResponseHandler {
	return func(r *http.Request) (HttpResponse, error) {
		qId := r.PathValue(idKey)
		id, err := uuid.Parse(qId)
		if err != nil {
			return nil, malformedId{qId, err}
		}
		context, err := contexts.GetById(id)
		if err != nil {
			return nil, err
		}
		return statusOk{context}, nil
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
		return statusCreated{context}, nil
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
		context, err := contexts.EditName(id, putCtx.NewName)
		if err != nil {
			return nil, err
		}
		return statusOk{context}, nil
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
		return statusNoContent{}, nil
	}
}
