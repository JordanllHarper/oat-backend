package main

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

type (
	fixture_dummyTaskStoreFailure    struct{}
	fixture_dummyTaskStoreSuccess    struct{ tasks []task }
	fixture_dummyContextStoreSuccess struct{ contexts []context }
	fixture_dummyContextStoreFailure struct{ contexts []context }
)

func (f fixture_dummyContextStoreFailure) SetNewCurrentTask(ctxId id, tid *id) error {
	return fixture_error()
}
func (f fixture_dummyTaskStoreFailure) InsertTask(t task) error                      { return fixture_error() }
func (f fixture_dummyTaskStoreFailure) All() ([]task, error)                         { return nil, fixture_error() }
func (f fixture_dummyContextStoreSuccess) SetNewCurrentTask(ctxId id, tid *id) error { return nil }
func (f fixture_dummyContextStoreFailure) Add(c context) error                       { return fixture_error() }
func (f fixture_dummyContextStoreSuccess) All() ([]context, error) {
	return []context{}, nil
}
func (f fixture_dummyContextStoreFailure) All() ([]context, error) {
	return nil, fixture_error()
}
func (f fixture_dummyTaskStoreFailure) AllForContext(id id) ([]task, error) {
	return nil, fixture_error()
}
func (f fixture_dummyTaskStoreFailure) ModifyTask(t task) (task, error) {
	return task{}, fixture_error()
}
func (f fixture_dummyTaskStoreFailure) RemoveTask(tId id) error             { return fixture_error() }
func (f fixture_dummyTaskStoreFailure) DeleteByContext(cId id) error        { return fixture_error() }
func (f fixture_dummyTaskStoreFailure) GetById(id id) (task, error)         { return task{}, fixture_error() }
func (f fixture_dummyTaskStoreSuccess) All() ([]task, error)                { return f.tasks, nil }
func (f fixture_dummyTaskStoreSuccess) AllForContext(id id) ([]task, error) { return f.tasks, nil }
func (f fixture_dummyTaskStoreSuccess) InsertTask(t task) error             { return nil }
func (f fixture_dummyTaskStoreSuccess) ModifyTask(t task) (task, error)     { return f.tasks[0], nil }
func (f fixture_dummyTaskStoreSuccess) RemoveTask(tId id) error             { return nil }
func (f fixture_dummyTaskStoreSuccess) DeleteByContext(cId id) error        { return nil }

func (f fixture_dummyTaskStoreSuccess) GetById(id id) (task, error) {
	return fixture_task(), nil
}

func fixture_error() error     { return fmt.Errorf("Test error") }
func fixture_dummyValidId() id { return uuid.MustParse("11111111-1111-1111-1111-111111111111") }

func fixture_task() task {
	return task{
		Id:        fixture_dummyValidId(),
		ContextId: fixture_dummyValidId(),
		Title:     "test",
		Notes:     "notes",
		Priority:  1,
	}
}
func fixture_contextWithTask() context {
	currentTaskId := fixture_dummyValidId()
	return context{
		Id:            fixture_dummyValidId(),
		Name:          "name",
		CurrentTaskId: &currentTaskId,
	}
}

func fixture_contextNilTask() context {
	return context{
		Id:            fixture_dummyValidId(),
		Name:          "name",
		CurrentTaskId: nil,
	}
}

func fixture_tasksStoreEmpty() taskStore { return &taskStoreImpl{} }
func fixture_tasksStorePopulatedMatching() taskStore {
	return &taskStoreImpl{
		fixture_task(),
	}
}

func Test_getCurrentTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ctx         context
		tasks       taskStore
		wantHasTask bool
		wantTask    task
		wantErr     bool
	}{
		{
			name:        "Ctx's current task id is nil",
			ctx:         fixture_contextNilTask(),
			tasks:       fixture_tasksStoreEmpty(),
			wantHasTask: false,
			wantTask:    task{},
			wantErr:     false,
		},
		{
			name:        "tasks.getById fails",
			ctx:         fixture_contextWithTask(),
			tasks:       fixture_dummyTaskStoreFailure{},
			wantHasTask: false,
			wantTask:    task{},
			wantErr:     true,
		},
		{
			name:        "success",
			ctx:         fixture_contextWithTask(),
			tasks:       fixture_tasksStorePopulatedMatching(),
			wantHasTask: true,
			wantTask:    fixture_task(),
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHasTask, gotTask, gotErr := getCurrentTask(tt.ctx, tt.tasks)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getCurrentTask() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getCurrentTask() succeeded unexpectedly")
			}
			if gotHasTask != tt.wantHasTask {
				t.Errorf("getCurrentTask() = %v, want %v", gotHasTask, tt.wantHasTask)
			}
			if gotTask != tt.wantTask {
				t.Errorf("getCurrentTask() = %v, want %v", gotTask, tt.wantTask)
			}
		})
	}
}

func Test_setNextTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ctx             context
		tasks           taskGetter
		contexts        contextCurrentTaskSetter
		wantHasNextTask bool
		wantTask        task
		wantErr         bool
	}{
		{
			name:            "Getting all tasks fails",
			ctx:             fixture_contextWithTask(),
			tasks:           fixture_dummyTaskStoreFailure{},
			contexts:        nil,
			wantHasNextTask: false,
			wantTask:        task{},
			wantErr:         true,
		},
		{
			name: "Length of all tasks is 0",
			ctx:  context{},
			tasks: fixture_dummyTaskStoreSuccess{
				tasks: []task{},
			},
			contexts: fixture_dummyContextStoreSuccess{
				contexts: []context{},
			},
			wantHasNextTask: false,
			wantTask:        task{},
			wantErr:         false,
		},
		{
			name:            "Err from contexts.SetNewCurrentTask",
			ctx:             context{},
			tasks:           fixture_dummyTaskStoreSuccess{},
			contexts:        fixture_dummyContextStoreFailure{},
			wantHasNextTask: false,
			wantTask:        task{},
			wantErr:         true,
		},
		{
			name: "success",
			ctx:  context{},
			tasks: fixture_dummyTaskStoreSuccess{
				tasks: []task{
					fixture_task(),
				},
			},
			contexts:        fixture_dummyContextStoreSuccess{},
			wantHasNextTask: true,
			wantTask:        fixture_task(),
			wantErr:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHasNextTask, gotTask, gotErr := setNextTask(tt.ctx, tt.tasks, tt.contexts)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("setNextTask() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("setNextTask() succeeded unexpectedly")
			}
			if gotHasNextTask != tt.wantHasNextTask {
				t.Errorf("setNextTask() = %v, want %v", gotHasNextTask, tt.wantHasNextTask)
			}
			if gotTask != tt.wantTask {
				t.Errorf("setNextTask() = %v, want %v", gotTask, tt.wantTask)
			}
		})
	}
}

func Test_completeAndGetNextTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		ctx             context
		tasks           taskGetterDeleter
		contexts        contextCurrentTaskSetter
		wantHasNextTask bool
		wantTask        task
		wantErr         bool
	}{
		{
			name:            "ctx.CurrentTaskId is nil",
			ctx:             fixture_contextNilTask(),
			tasks:           fixture_dummyTaskStoreFailure{},
			contexts:        nil,
			wantHasNextTask: false,
			wantTask:        task{},
			wantErr:         false,
		},
		{
			name:            "tasks.RemoveTask fails",
			ctx:             fixture_contextWithTask(),
			tasks:           fixture_dummyTaskStoreFailure{},
			contexts:        nil,
			wantHasNextTask: false,
			wantTask:        task{},
			wantErr:         true,
		},
		{
			name:            "contexts.SetNewCurrentTask fails",
			ctx:             fixture_contextWithTask(),
			tasks:           fixture_dummyTaskStoreSuccess{},
			contexts:        fixture_dummyContextStoreFailure{},
			wantHasNextTask: false,
			wantTask:        task{},
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotHasNextTask, gotTask, gotErr := completeAndGetNextTask(tt.ctx, tt.tasks, tt.contexts)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("completeAndGetNextTask() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("completeAndGetNextTask() succeeded unexpectedly")
			}
			if gotHasNextTask != tt.wantHasNextTask {
				t.Errorf("completeAndGetNextTask() = %v, want %v", gotHasNextTask, tt.wantHasNextTask)
			}
			if gotTask != tt.wantTask {
				t.Errorf("completeAndGetNextTask() = %v, want %v", gotTask, tt.wantTask)
			}
		})
	}
}

func Test_priorityValid(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		priority priority
		want     bool
	}{
		{
			name:     "Too low priority",
			priority: 0,
			want:     false,
		},
		{
			name:     "Valid priority",
			priority: 1,
			want:     true,
		},
		{
			name:     "Valid priority",
			priority: 2,
			want:     true,
		},
		{
			name:     "Valid priority",
			priority: 3,
			want:     true,
		},
		{
			name:     "Too high priority",
			priority: 5,
			want:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := priorityValid(tt.priority)
			if got != tt.want {
				t.Errorf("priorityValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_addTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tasks    taskStore
		ctxId    id
		title    string
		notes    string
		p        priority
		wantTask task
		wantErr  bool
	}{
		{
			name:     "Invalid priority",
			tasks:    nil,
			ctxId:    id{},
			title:    "",
			notes:    "",
			p:        0,
			wantTask: task{},
			wantErr:  true,
		},
		{
			name:     "Empty title",
			tasks:    nil,
			ctxId:    id{},
			title:    "",
			notes:    "",
			p:        1,
			wantTask: task{},
			wantErr:  true,
		},
		{
			name:     "tasks.InsertTask fails",
			tasks:    fixture_dummyTaskStoreFailure{},
			ctxId:    id{},
			title:    "title",
			notes:    "",
			p:        1,
			wantTask: task{},
			wantErr:  true,
		},
		{
			name:  "success",
			tasks: fixture_dummyTaskStoreSuccess{},
			ctxId: fixture_dummyValidId(),
			title: "title",
			notes: "notes",
			p:     1,
			wantTask: task{
				Title:    "title",
				Notes:    "notes",
				Priority: 1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotTask, gotErr := addTask(tt.tasks, fixture_dummyValidId, tt.ctxId, tt.title, tt.notes, tt.p)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("addTask() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("addTask() succeeded unexpectedly")
			}
			if gotTask == tt.wantTask {
				t.Errorf("addTask() = %v, want %v", gotTask, tt.wantTask)
			}
		})
	}
}

func Test_getAllContexts(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		contexts     contextGetter
		wantContexts HttpResponse
		wantErr      bool
	}{
		{
			name:         "contexts.All fails",
			contexts:     fixture_dummyContextStoreFailure{},
			wantContexts: nil,
			wantErr:      true,
		},
		{
			name:         "success",
			contexts:     fixture_dummyContextStoreSuccess{},
			wantContexts: statusOk{bodyMultipleValues[context]{}},
			wantErr:      false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := getAllContexts(tt.contexts)
			if gotErr != nil {
				if !tt.wantErr {
					t.Errorf("getAllContexts() failed: %v", gotErr)
				}
				return
			}
			if tt.wantErr {
				t.Fatal("getAllContexts() succeeded unexpectedly")
			}
		})
	}
}
