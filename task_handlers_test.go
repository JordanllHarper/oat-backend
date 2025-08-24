package main

import "testing"

func Test_handleGetTaskById(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tasks taskStore
		want  HttpResponseHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleGetTaskById(tt.tasks)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("handleGetTaskById() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleGetCurrentTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tasks    taskStore
		contexts contextStore
		want     HttpResponseHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleGetCurrentTask(tt.tasks, tt.contexts)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("handleGetCurrentTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handleCompleteCurrentTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tasks    taskStore
		contexts contextStore
		want     HttpResponseHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handleCompleteCurrentTask(tt.tasks, tt.contexts)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("handleCompleteCurrentTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlePostTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tasks    taskStore
		contexts contextStore
		want     HttpResponseHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlePostTask(tt.tasks, tt.contexts)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("handlePostTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlePostCurrentTask(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tasks    taskStore
		contexts contextStore
		want     HttpResponseHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlePostCurrentTask(tt.tasks, tt.contexts)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("handlePostCurrentTask() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_handlePutTaskById(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		tasks    taskStore
		contexts contextStore
		want     HttpResponseHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := handlePutTaskById(tt.tasks, tt.contexts)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("handlePutTaskById() = %v, want %v", got, tt.want)
			}
		})
	}
}
