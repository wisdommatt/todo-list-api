// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
	tasks "github.com/wisdommatt/creativeadvtech-assessment/components/tasks"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// CreateTask provides a mock function with given fields: ctx, task
func (_m *Service) CreateTask(ctx context.Context, task tasks.Task) (*tasks.Task, error) {
	ret := _m.Called(ctx, task)

	var r0 *tasks.Task
	if rf, ok := ret.Get(0).(func(context.Context, tasks.Task) *tasks.Task); ok {
		r0 = rf(ctx, task)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tasks.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, tasks.Task) error); ok {
		r1 = rf(ctx, task)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteTask provides a mock function with given fields: ctx, taskID
func (_m *Service) DeleteTask(ctx context.Context, taskID string) (*tasks.Task, error) {
	ret := _m.Called(ctx, taskID)

	var r0 *tasks.Task
	if rf, ok := ret.Get(0).(func(context.Context, string) *tasks.Task); ok {
		r0 = rf(ctx, taskID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tasks.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, taskID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTask provides a mock function with given fields: ctx, taskID
func (_m *Service) GetTask(ctx context.Context, taskID string) (*tasks.Task, error) {
	ret := _m.Called(ctx, taskID)

	var r0 *tasks.Task
	if rf, ok := ret.Get(0).(func(context.Context, string) *tasks.Task); ok {
		r0 = rf(ctx, taskID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*tasks.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, taskID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetTasks provides a mock function with given fields: ctx, userID, lastID, limit
func (_m *Service) GetTasks(ctx context.Context, userID string, lastID string, limit int) ([]tasks.Task, error) {
	ret := _m.Called(ctx, userID, lastID, limit)

	var r0 []tasks.Task
	if rf, ok := ret.Get(0).(func(context.Context, string, string, int) []tasks.Task); ok {
		r0 = rf(ctx, userID, lastID, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]tasks.Task)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, string, int) error); ok {
		r1 = rf(ctx, userID, lastID, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
