package telegram_test

import (
	"context"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	notifymock "github.com/slok/alertgram/internal/mocks/notify"
	telegrammock "github.com/slok/alertgram/internal/mocks/notify/telegram"
	"github.com/slok/alertgram/internal/model"
	"github.com/slok/alertgram/internal/notify/telegram"
)

func GetBaseAlertGroup() *model.AlertGroup {
	return &model.AlertGroup{
		ID: "test-alert",
		Alerts: []model.Alert{
			{
				Labels: map[string]string{
					"alertname": "ServicePodIsRestarting",
				},
				Annotations: map[string]string{
					"message": "There has been restarting more than 5 times over 20 minutes",
				},
			},
			{
				Labels: map[string]string{
					"alertname": "ServicePodIsRestarting",
					"chatid":    "-1001234567890",
				},
				Annotations: map[string]string{
					"message": "There has been restarting more than 5 times over 20 minutes",
					"graph":   "https://prometheus.test/my-graph",
				},
			},
		},
	}
}

var errTest = errors.New("whatever")

func TestNotify(t *testing.T) {
	tests := map[string]struct {
		cfg        telegram.Config
		mocks      func(t *testing.T, mcli *telegrammock.Client, mr *notifymock.TemplateRenderer)
		alertGroup *model.AlertGroup
		expErr     error
	}{
		"A alertGroup should be rendered and send the message to telegram.": {
			cfg: telegram.Config{
				DefaultTelegramChatID: 1234,
			},
			mocks: func(t *testing.T, mcli *telegrammock.Client, mr *notifymock.TemplateRenderer) {
				expMsgData := "rendered template"
				expAlertGroup := GetBaseAlertGroup()
				mr.On("Render", mock.Anything, expAlertGroup).Once().Return(expMsgData, nil)

				expMsg := tgbotapi.MessageConfig{
					BaseChat:  tgbotapi.BaseChat{ChatID: 1234},
					ParseMode: "HTML",
					Text:      expMsgData,
				}
				mcli.On("Send", expMsg).Once().Return(tgbotapi.Message{}, nil)
			},
			alertGroup: GetBaseAlertGroup(),
		},

		"A error in the template rendering process should be processed.": {
			cfg: telegram.Config{
				DefaultTelegramChatID: 1234,
			},
			mocks: func(t *testing.T, mcli *telegrammock.Client, mr *notifymock.TemplateRenderer) {
				expAlertGroup := GetBaseAlertGroup()
				mr.On("Render", mock.Anything, expAlertGroup).Once().Return("", errTest)
			},
			alertGroup: GetBaseAlertGroup(),
			expErr:     errTest,
		},

		"A error in the notification send process should be processed with communication error.": {
			cfg: telegram.Config{
				DefaultTelegramChatID: 1234,
			},
			mocks: func(t *testing.T, mcli *telegrammock.Client, mr *notifymock.TemplateRenderer) {
				expMsgData := "rendered template"
				expAlertGroup := GetBaseAlertGroup()
				mr.On("Render", mock.Anything, expAlertGroup).Once().Return(expMsgData, nil)

				expMsg := tgbotapi.MessageConfig{
					BaseChat:  tgbotapi.BaseChat{ChatID: 1234},
					ParseMode: "HTML",
					Text:      expMsgData,
				}
				mcli.On("Send", expMsg).Once().Return(tgbotapi.Message{}, errTest)
			},
			alertGroup: GetBaseAlertGroup(),
			expErr:     telegram.ErrComm,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Mocks.
			mcli := &telegrammock.Client{}
			mr := &notifymock.TemplateRenderer{}
			test.mocks(t, mcli, mr)
			test.cfg.Client = mcli
			test.cfg.TemplateRenderer = mr

			// Execute.
			n, err := telegram.NewNotifier(test.cfg)
			require.NoError(err)
			err = n.Notify(context.TODO(), test.alertGroup)

			// Check.
			if test.expErr != nil && assert.Error(err) {
				assert.True(errors.Is(err, test.expErr))
			} else if assert.NoError(err) {
				mcli.AssertExpectations(t)
				mr.AssertExpectations(t)
			}
		})
	}
}
