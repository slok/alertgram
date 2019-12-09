package notify_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/alertgram/internal/model"
	"github.com/slok/alertgram/internal/notify"
)

func TestTemplateRenderer(t *testing.T) {
	tests := map[string]struct {
		alertGroup *model.AlertGroup
		renderer   func() notify.TemplateRenderer
		expData    string
		expErr     error
	}{
		"Invalid template should return an error.": {
			renderer: func() notify.TemplateRenderer {
				r, _ := notify.NewHTMLTemplateRenderer("{{ .ID }}")
				return r
			},
			expErr: notify.ErrRenderTemplate,
		},

		"Custom template should render the alerts correctly.": {
			alertGroup: &model.AlertGroup{
				ID:     "test-alert",
				Alerts: []model.Alert{{}, {}, {}},
			},
			renderer: func() notify.TemplateRenderer {
				r, _ := notify.NewHTMLTemplateRenderer("{{ .ID }} has {{ .Alerts | len }} alerts.")
				return r
			},
			expData: "test-alert has 3 alerts.",
		},

		"Default template should render the alerts correctly.": {
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
			expData: `
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
			renderer: func() notify.TemplateRenderer { return notify.DefaultTemplateRenderer },
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)

			r := test.renderer()
			gotData, err := r.Render(context.TODO(), test.alertGroup)

			if test.expErr != nil && assert.Error(err) {
				assert.True(errors.Is(err, test.expErr))
			} else if assert.NoError(err) {
				assert.Equal(test.expData, gotData)
			}
		})
	}
}
