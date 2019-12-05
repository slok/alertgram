// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"

	model "github.com/slok/alertgram/internal/model"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// Forward provides a mock function with given fields: ctx, alertGroup
func (_m *Service) Forward(ctx context.Context, alertGroup *model.AlertGroup) error {
	ret := _m.Called(ctx, alertGroup)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *model.AlertGroup) error); ok {
		r0 = rf(ctx, alertGroup)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}