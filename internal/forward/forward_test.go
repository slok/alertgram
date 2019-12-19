package forward_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/alertgram/internal/forward"
	forwardmock "github.com/slok/alertgram/internal/mocks/forward"
	"github.com/slok/alertgram/internal/model"
)

var errTest = errors.New("whatever")

func TestServiceForward(t *testing.T) {
	tests := map[string]struct {
		cfg        forward.ServiceConfig
		props      forward.Properties
		alertGroup *model.AlertGroup
		mock       func(ns []*forwardmock.Notifier)
		expErr     error
	}{
		"A nil alert group should fail.": {
			mock:   func(ns []*forwardmock.Notifier) {},
			expErr: forward.ErrInvalidAlertGroup,
		},

		"A forwarded alerts should be send to all notifiers.": {
			props: forward.Properties{
				CustomChatID: "-1001234567890",
			},
			alertGroup: &model.AlertGroup{
				ID:     "test-group",
				Alerts: []model.Alert{model.Alert{Name: "test"}},
			},
			mock: func(ns []*forwardmock.Notifier) {
				expNotification := forward.Notification{
					ChatID: "-1001234567890",
					AlertGroup: model.AlertGroup{
						ID:     "test-group",
						Alerts: []model.Alert{model.Alert{Name: "test"}},
					},
				}
				for _, n := range ns {
					n.On("Notify", mock.Anything, expNotification).Once().Return(nil)
				}
			},
		},

		"Errors from notifiers should be ignored to the callers and all notifiers should be called.": {
			alertGroup: &model.AlertGroup{
				ID:     "test-group",
				Alerts: []model.Alert{model.Alert{Name: "test"}},
			},
			mock: func(ns []*forwardmock.Notifier) {
				expNotification := forward.Notification{
					AlertGroup: model.AlertGroup{
						ID:     "test-group",
						Alerts: []model.Alert{model.Alert{Name: "test"}},
					},
				}
				for i, n := range ns {
					err := errTest
					// Set error in the first one.
					if i != 0 {
						err = nil
					}
					n.On("Notify", mock.Anything, expNotification).Once().Return(err)
					n.On("Type").Maybe().Return("")
				}
			},
		},

		"Alerts that have the label for custom chat ids should be grouped together.": {
			cfg: forward.ServiceConfig{
				AlertLabelChatID: "test_chat_id",
			},
			props: forward.Properties{
				CustomChatID: "-1001234567890",
			},
			alertGroup: &model.AlertGroup{
				ID: "test-group",
				Alerts: []model.Alert{
					{Name: "test-1", Labels: map[string]string{"test_chat_id": ""}},
					{Name: "test-2", Labels: map[string]string{"test_chat_id": "chat2"}},
					{Name: "test-3", Labels: map[string]string{"test_chat_id": "chat1"}},
					{Name: "test-4", Labels: map[string]string{"test_chat_id": "chat2"}},
					{Name: "test-3", Labels: map[string]string{"test_chat_id": "chat1"}},
					{Name: "test-6", Labels: map[string]string{"test_chat_id": "chat3"}},
					{Name: "test-7", Labels: map[string]string{"test_chat_id": ""}},
					{Name: "test-8", Labels: map[string]string{"test_chat_id": ""}},
					{Name: "test-9", Labels: map[string]string{"test_chat_id": "chat1"}},
				},
			},
			mock: func(ns []*forwardmock.Notifier) {
				expNotChatDef := forward.Notification{
					ChatID: "-1001234567890",
					AlertGroup: model.AlertGroup{ID: "test-group",
						Alerts: []model.Alert{
							{Name: "test-1", Labels: map[string]string{"test_chat_id": ""}},
							{Name: "test-7", Labels: map[string]string{"test_chat_id": ""}},
							{Name: "test-8", Labels: map[string]string{"test_chat_id": ""}},
						},
					},
				}
				expNotChat1 := forward.Notification{
					ChatID: "chat1",
					AlertGroup: model.AlertGroup{ID: "test-group-chat1",
						Alerts: []model.Alert{
							{Name: "test-3", Labels: map[string]string{"test_chat_id": "chat1"}},
							{Name: "test-3", Labels: map[string]string{"test_chat_id": "chat1"}},
							{Name: "test-9", Labels: map[string]string{"test_chat_id": "chat1"}},
						},
					},
				}
				expNotChat2 := forward.Notification{
					ChatID: "chat2",
					AlertGroup: model.AlertGroup{ID: "test-group-chat2",
						Alerts: []model.Alert{
							{Name: "test-2", Labels: map[string]string{"test_chat_id": "chat2"}},
							{Name: "test-4", Labels: map[string]string{"test_chat_id": "chat2"}},
						},
					},
				}
				expNotChat3 := forward.Notification{
					ChatID: "chat3",
					AlertGroup: model.AlertGroup{ID: "test-group-chat3",
						Alerts: []model.Alert{
							{Name: "test-6", Labels: map[string]string{"test_chat_id": "chat3"}},
						},
					},
				}
				for _, n := range ns {
					n.On("Notify", mock.Anything, expNotChatDef).Once().Return(nil)
					n.On("Notify", mock.Anything, expNotChat1).Once().Return(nil)
					n.On("Notify", mock.Anything, expNotChat2).Once().Return(nil)
					n.On("Notify", mock.Anything, expNotChat3).Once().Return(nil)
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
			svc, err := forward.NewService(test.cfg)
			require.NoError(err)

			err = svc.Forward(context.TODO(), test.props, test.alertGroup)

			if test.expErr != nil && assert.Error(err) {
				assert.True(errors.Is(err, test.expErr))
			} else if assert.NoError(err) {
				mn1.AssertExpectations(t)
				mn2.AssertExpectations(t)
			}
		})
	}
}
