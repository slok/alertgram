package alertmanager_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/prometheus/alertmanager/notify/webhook"
	"github.com/prometheus/alertmanager/template"
	prommodel "github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/alertgram/internal/http/alertmanager"
	forwardmock "github.com/slok/alertgram/internal/mocks/forward"
	"github.com/slok/alertgram/internal/model"
)

var t0 = time.Now().UTC()

func getBaseAlert() model.Alert {
	return model.Alert{
		ID:           "test-alert",
		Name:         "test-alert-name",
		StartsAt:     t0.Add(-10 * time.Minute),
		EndsAt:       t0.Add(-3 * time.Minute),
		Status:       model.AlertStatusFiring,
		Labels:       map[string]string{prommodel.AlertNameLabel: "test-alert-name", "lK1": "lV1", "lK2": "lV2"},
		Annotations:  map[string]string{"aK1": "aV1", "aK2": "aV2"},
		GeneratorURL: "http://test.com",
	}
}

func getBaseAlerts() *model.AlertGroup {
	al1 := getBaseAlert()
	al1.ID += "-1"
	al2 := getBaseAlert()
	al2.ID += "-2"

	return &model.AlertGroup{
		ID:     "test-group",
		Labels: map[string]string{"glK1": "glV1", "glK2": "glV2"},
		Alerts: []model.Alert{al1, al2},
	}
}

func getBaseAlertmanagerAlert() template.Alert {
	return template.Alert{
		Fingerprint:  "test-alert",
		Status:       string(prommodel.AlertFiring),
		Labels:       map[string]string{prommodel.AlertNameLabel: "test-alert-name", "lK1": "lV1", "lK2": "lV2"},
		Annotations:  map[string]string{"aK1": "aV1", "aK2": "aV2"},
		StartsAt:     t0.Add(-10 * time.Minute),
		EndsAt:       t0.Add(-3 * time.Minute),
		GeneratorURL: "http://test.com",
	}
}

func getBaseAlertmanagerAlerts() webhook.Message {
	al1 := getBaseAlertmanagerAlert()
	al1.Fingerprint += "-1"
	al2 := getBaseAlertmanagerAlert()
	al2.Fingerprint += "-2"

	return webhook.Message{
		Data: &template.Data{
			Receiver:          "test-recv",
			Status:            string(prommodel.AlertFiring),
			Alerts:            template.Alerts{al1, al2},
			GroupLabels:       map[string]string{"glK1": "glV1", "glK2": "glV2"},
			CommonLabels:      map[string]string{"gclK1": "gclV1", "gclK2": "gclV2"},
			CommonAnnotations: map[string]string{"gcaK1": "gcaV1", "gcaK2": "gcaV2"},
			ExternalURL:       "http://test.com",
		},
		Version:  "4",
		GroupKey: "test-group",
	}
}

func TestHandleAlerts(t *testing.T) {
	tests := map[string]struct {
		webhookAlertJSON func(t *testing.T) string
		mock             func(t *testing.T, msvc *forwardmock.Service)
		expCode          int
	}{
		"Alertmanager webhook alerts request should be handled correctly.": {
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *forwardmock.Service) {
				expAlerts := getBaseAlerts()
				msvc.On("Forward", mock.Anything, expAlerts).Once().Return(nil)
			},
			expCode: http.StatusOK,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Mocks.
			msvc := &forwardmock.Service{}
			test.mock(t, msvc)

			// Execute.
			h, _ := alertmanager.NewHandler(alertmanager.Config{
				WebhookPath: "/test-alerts",
				Forwarder:   msvc,
			})
			srv := httptest.NewServer(h)
			defer srv.Close()
			req, err := http.NewRequest(http.MethodPost, srv.URL+"/test-alerts", strings.NewReader(test.webhookAlertJSON(t)))
			require.NoError(err)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(err)

			// Check.
			assert.Equal(test.expCode, resp.StatusCode)
			msvc.AssertExpectations(t)
		})
	}
}
