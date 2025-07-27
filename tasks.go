package main

import (
	"slices"
)

type (
	taskStore interface {
		IDable[task]
		All() ([]task, error)
		AllForContext(id id) ([]task, error)
		InsertTask(t task) error
		ModifyTask(t task) (task, error)
		RemoveTask(tId id) error
		DeleteByContext(cId id) error
	}

	taskStoreImpl []task
)

func (tsi taskStoreImpl) All() ([]task, error) { return tsi, nil }
func (tsi taskStoreImpl) AllForContext(id id) ([]task, error) {
	newTasks := []task{}
	// NOTE: Might not matter as this is for testing purposes but we are throwing branching out the window here, consider storing lists of task lists for each context
	for _, task := range tsi {
		if task.ContextId == id {
			newTasks = append(newTasks, task)
		}
	}
	return newTasks, nil
}

func (tsi *taskStoreImpl) ModifyTask(modified task) (task, error) {
	index := slices.IndexFunc(*tsi, compareTasksById(modified))
	if index == -1 {
		return task{}, idNotFound(modified.Id)
	}
	(*tsi)[index] = modified
	return modified, nil
}

func (tsi *taskStoreImpl) GetById(id id) (task, error) {
	for _, t := range *tsi {
		if t.Id == id {
			return t, nil
		}
	}
	return task{}, idNotFound(id)
}

func (tsi *taskStoreImpl) InsertTask(t task) error {
	idExists := slices.ContainsFunc(*tsi, compareTasksById(t))
	if idExists {
		return idAlreadyExists(t.Id)
	}
	*tsi = append(*tsi, t)
	return nil
}

func (tsi *taskStoreImpl) RemoveTask(tId id) error {
	idx := slices.IndexFunc(*tsi, compareIdToTask(tId))
	if idx == -1 {
		return idNotFound(tId)
	}
	*tsi = sliceUnorderedRemove((*tsi), idx)
	return nil

}
func (tsi *taskStoreImpl) DeleteByContext(cId id) error {
	*tsi = slices.DeleteFunc(*tsi, func(t task) bool {
		return t.ContextId == cId
	})
	return nil
}
