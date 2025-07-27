package main

import (
	"cmp"
	"maps"
	"slices"
)

type contextStore interface {
	IDable[context]
	All() ([]context, error)
	Add(c context) error
	EditName(id id, name string) (context, error)
	SetNewCurrentTask(ctxId id, tId *id) error
	Delete(id id) error
}

type contextStoreImpl map[id]context

func (csi contextStoreImpl) All() ([]context, error) {
	values := maps.Values(csi)
	s := slices.Collect(values)
	slices.SortFunc(s, func(a, b context) int {
		return cmp.Compare(a.Name, b.Name)
	})
	return s, nil
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
func (csi contextStoreImpl) EditName(id id, newName string) (context, error) {
	ctx, found := csi[id]
	if !found {
		return context{}, idNotFound(id)
	}
	ctx.Name = newName
	csi[id] = ctx
	return ctx, nil
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
