package telegram_test

import (
	"context"
	"errors"
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	telegrammock "github.com/slok/alertgram/internal/mocks/notify/telegram"
	"github.com/slok/alertgram/internal/model"
	"github.com/slok/alertgram/internal/notify/telegram"
)

func TestNotify(t *testing.T) {
	tests := map[string]struct {
		cfg        telegram.Config
		mocks      func(t *testing.T, mcli *telegrammock.Client)
		alertGroup *model.AlertGroup
		expErr     error
	}{
		"An alertGroup should be rendered and notify to telegram (Default template).": {
			cfg: telegram.Config{
				DefaultTelegramChatID: 1234,
			},
			mocks: func(t *testing.T, mcli *telegrammock.Client) {
				expMsg := tgbotapi.MessageConfig{
					BaseChat:  tgbotapi.BaseChat{ChatID: 1234},
					ParseMode: "HTML",
					Text: `
🚨🚨 FIRING 2 🚨🚨
💥💥💥 <b>ServicePodIsRestarting</b> 💥💥💥
  There has been restarting more than 5 times over 20 minutes
	🔹 chatid: -1001234567890
	🔹 job: kubernetes-metrics
	🔹 owner: team1
	🔹 pod: ns1/pod-service1-f76c976c4-9hlgv
	🔹 severity: telegram
	🔸 <a href="https://prometheus.test/my-graph">graph</a>
	🔸 <a href="https://github.test/runbooks/pod-restarting.md">runbook</a>
💥💥💥 <b>ServicePodIsRestarting</b> 💥💥💥
  There has been restarting more than 5 times over 20 minutes
	🔹 chatid: -1001234567890
	🔹 job: kubernetes-metrics
	🔹 owner: team1
	🔹 pod: ns1/pod-service64-f5c7dd9cfc5-8scht
	🔹 severity: telegram
	🔸 <a href="https://prometheus.test/my-graph">graph</a>
	🔸 <a href="https://github.test/runbooks/pod-restarting.md">runbook</a>
`,
				}
				mcli.On("Send", expMsg).Once().Return(tgbotapi.Message{}, nil)
			},
			alertGroup: &model.AlertGroup{
				ID: "test-alert",
				Alerts: []model.Alert{
					{
						Labels: map[string]string{
							"alertname": "ServicePodIsRestarting",
							"chatid":    "-1001234567890",
							"job":       "kubernetes-metrics",
							"owner":     "team1",
							"pod":       "ns1/pod-service1-f76c976c4-9hlgv",
							"severity":  "telegram",
						},
						Annotations: map[string]string{
							"message": "There has been restarting more than 5 times over 20 minutes",
							"graph":   "https://prometheus.test/my-graph",
							"runbook": "https://github.test/runbooks/pod-restarting.md",
						},
					},
					{
						Labels: map[string]string{
							"alertname": "ServicePodIsRestarting",
							"chatid":    "-1001234567890",
							"job":       "kubernetes-metrics",
							"owner":     "team1",
							"pod":       "ns1/pod-service64-f5c7dd9cfc5-8scht",
							"severity":  "telegram",
						},
						Annotations: map[string]string{
							"message": "There has been restarting more than 5 times over 20 minutes",
							"graph":   "https://prometheus.test/my-graph",
							"runbook": "https://github.test/runbooks/pod-restarting.md",
						},
					},
				},
			},
		},

		"An alertGroup should be rendered and notify to telegram (Custom template).": {
			cfg: telegram.Config{
				DefaultTelegramChatID: 1234,
				AlertMessageTemplate:  "{{ .ID }} has {{ .Alerts | len }} alerts.",
			},
			mocks: func(t *testing.T, mcli *telegrammock.Client) {
				expMsg := tgbotapi.MessageConfig{
					BaseChat:  tgbotapi.BaseChat{ChatID: 1234},
					Text:      "test-alert has 3 alerts.",
					ParseMode: "HTML",
				}
				mcli.On("Send", expMsg).Once().Return(tgbotapi.Message{}, nil)
			},
			alertGroup: &model.AlertGroup{
				ID:     "test-alert",
				Alerts: []model.Alert{{}, {}, {}},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			mcli := &telegrammock.Client{}
			test.mocks(t, mcli)

			test.cfg.Client = mcli
			n, err := telegram.NewNotifier(test.cfg)
			require.NoError(err)
			err = n.Notify(context.TODO(), test.alertGroup)

			if test.expErr != nil && assert.Error(err) {
				assert.True(errors.Is(err, test.expErr))
			} else if assert.NoError(err) {
				mcli.AssertExpectations(t)
			}
		})
	}
}
