package main

import (
	"maps"
	"slices"
)

type (
	contextGetter interface {
		All() ([]context, error)
	}
	contextAdder interface {
		Add(c context) error
	}
	contextEditer interface {
		EditName(id id, name string) error
	}
	contextCurrentTaskSetter interface {
		SetNewCurrentTask(ctxId id, tId *id) error
	}
	contextDeleter interface {
		Delete(id id) error
	}
	contextStore interface {
		IDable[context]
		contextGetter
		contextAdder
		contextEditer
		contextCurrentTaskSetter
		contextDeleter
	}

	contextStoreImpl map[id]context
)

func (csi contextStoreImpl) All() ([]context, error) {
	values := slices.Collect(maps.Values(csi))
	if values == nil {
		// return an empty context list rather than nil
		return []context{}, nil
	}
	return values, nil
}

func (csi contextStoreImpl) GetById(id id) (context, error) {
	ctx, found := csi[id]
	if !found {
		return context{}, idNotFound(id)
	}
	return ctx, nil
}
func (csi contextStoreImpl) Add(c context) error {
	_, found := csi[c.Id]
	if found {
		return idAlreadyExists(c.Id)
	}
	csi[c.Id] = c
	return nil
}
func (csi contextStoreImpl) EditName(id id, newName string) error {
	ctx, found := csi[id]
	if !found {
		return idNotFound(id)
	}
	ctx.Name = newName
	csi[id] = ctx
	return nil
}
func (csi contextStoreImpl) Delete(id id) error {
	_, found := csi[id]
	if !found {
		return idNotFound(id)
	}
	delete(csi, id)
	return nil
}

func (csi contextStoreImpl) SetNewCurrentTask(ctxId id, tId *id) error {
	ctx, err := csi.GetById(ctxId)
	if err != nil {
		return err
	}
	ctx.CurrentTaskId = tId
	csi[ctx.Id] = ctx
	return nil
}
