package forward_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/log"
	forwardmock "github.com/slok/alertgram/internal/mocks/forward"
	"github.com/slok/alertgram/internal/model"
)

var errTest = errors.New("whatever")

func TestServiceForward(t *testing.T) {
	tests := map[string]struct {
		alert  model.Alert
		mock   func(ns []*forwardmock.Notifier)
		expErr error
	}{
		"A forwarded alert should be send to all notifiers.": {
			alert: model.Alert{Name: "test"},
			mock: func(ns []*forwardmock.Notifier) {
				expAlert := model.Alert{Name: "test"}
				for _, n := range ns {
					n.On("Notify", mock.Anything, expAlert).Once().Return(nil)
				}
			},
		},

		"Errors from notifiers should be ignored to the callers and all notifiers should be called.": {
			alert: model.Alert{Name: "test"},
			mock: func(ns []*forwardmock.Notifier) {
				expAlert := model.Alert{Name: "test"}
				for i, n := range ns {
					err := errTest
					// Set error in the first one.
					if i != 0 {
						err = nil
					}
					n.On("Notify", mock.Anything, expAlert).Once().Return(err)
					n.On("Type").Maybe().Return("")
				}
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			mn1 := &forwardmock.Notifier{}
			mn2 := &forwardmock.Notifier{}
			test.mock([]*forwardmock.Notifier{mn1, mn2})

			svc := forward.NewService([]forward.Notifier{mn1, mn2}, log.Dummy)
			err := svc.Forward(context.TODO(), test.alert)

			if test.expErr != nil && assert.Error(err) {
				assert.True(errors.Is(err, test.expErr))
			} else if assert.NoError(err) {
				mn1.AssertExpectations(t)
				mn2.AssertExpectations(t)
			}
		})
	}
}
