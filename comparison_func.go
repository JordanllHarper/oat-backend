package main

func compareIdToTask(i1 id) func(i2 task) bool {
	return func(i2 task) bool { return i1 == i2.Id }
}

func compareTasksById(i1 task) func(i2 task) bool {
	return func(i2 task) bool { return i1.Id == i2.Id }
}
