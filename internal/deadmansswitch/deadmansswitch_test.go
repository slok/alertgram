package deadmansswitch_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/alertgram/internal/deadmansswitch"
	"github.com/slok/alertgram/internal/forward"
	forwardmock "github.com/slok/alertgram/internal/mocks/forward"
	"github.com/slok/alertgram/internal/model"
)

func TestServiceDeadMansSwitch(t *testing.T) {
	tests := map[string]struct {
		cfg    deadmansswitch.Config
		exec   func(svc deadmansswitch.Service) error
		mock   func(ns []*forwardmock.Notifier)
		expErr error
	}{
		"If the alert is not received in the interval it should notify.": {
			cfg: deadmansswitch.Config{
				Interval: 10 * time.Millisecond,
			},
			exec: func(svc deadmansswitch.Service) error {
				// Give time to interval to act.
				time.Sleep(15 * time.Millisecond)
				return nil
			},
			mock: func(ns []*forwardmock.Notifier) {
				for _, n := range ns {
					n.On("Notify", mock.Anything, mock.Anything).Once().Return(nil)
					n.On("Type").Maybe().Return("")
				}
			},
		},

		"If the alert is received in the interval it should not notify.": {
			cfg: deadmansswitch.Config{
				Interval: 10 * time.Millisecond,
			},
			exec: func(svc deadmansswitch.Service) error {
				// Give time to interval to act.
				time.Sleep(6 * time.Millisecond)
				err := svc.PushSwitch(context.TODO(), &model.AlertGroup{})
				time.Sleep(6 * time.Millisecond)
				return err
			},
			mock: func(ns []*forwardmock.Notifier) {},
		},

		"If the alert is received and then stops being received in the interval it should not notify.": {
			cfg: deadmansswitch.Config{
				Interval: 10 * time.Millisecond,
			},
			exec: func(svc deadmansswitch.Service) error {
				time.Sleep(6 * time.Millisecond)
				err := svc.PushSwitch(context.TODO(), &model.AlertGroup{})
				if err != nil {
					return err
				}
				time.Sleep(6 * time.Millisecond)
				err = svc.PushSwitch(context.TODO(), &model.AlertGroup{})
				time.Sleep(15 * time.Millisecond)
				return err
			},
			mock: func(ns []*forwardmock.Notifier) {
				for _, n := range ns {
					n.On("Notify", mock.Anything, mock.Anything).Once().Return(nil)
					n.On("Type").Maybe().Return("")
				}
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			mn1 := &forwardmock.Notifier{}
			mn2 := &forwardmock.Notifier{}
			test.mock([]*forwardmock.Notifier{mn1, mn2})

			test.cfg.Notifiers = []forward.Notifier{mn1, mn2}
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			svc, err := deadmansswitch.NewService(ctx, test.cfg)
			require.NoError(err)
			err = test.exec(svc)

			if test.expErr != nil && assert.Error(err) {
				assert.True(errors.Is(err, test.expErr))
			} else if assert.NoError(err) {
				mn1.AssertExpectations(t)
				mn2.AssertExpectations(t)
			}
		})
	}
}
