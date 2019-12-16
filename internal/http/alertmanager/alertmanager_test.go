package alertmanager_test

import (
	"encoding/json"
	"errors"
	"fmt"
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

	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/http/alertmanager"
	"github.com/slok/alertgram/internal/internalerrors"
	deadmansswitchmock "github.com/slok/alertgram/internal/mocks/deadmansswitch"
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
		config           alertmanager.Config
		urlPath          string
		webhookAlertJSON func(t *testing.T) string
		mock             func(t *testing.T, msvc *forwardmock.Service)
		expCode          int
	}{
		"Alertmanager webhook alerts request should be handled correctly (with defaults).": {
			urlPath: "/alerts",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *forwardmock.Service) {
				expAlerts := getBaseAlerts()
				expProps := forward.Properties{}
				msvc.On("Forward", mock.Anything, expProps, expAlerts).Once().Return(nil)
			},
			expCode: http.StatusOK,
		},

		"Alertmanager webhook alerts request should be handled correctly (with custom params).": {
			config: alertmanager.Config{
				WebhookPath:       "/test-alerts",
				ChatIDQueryString: "custom-telegram-chat-id",
			},
			urlPath: "/test-alerts?custom-telegram-chat-id=-1009876543210",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *forwardmock.Service) {
				expAlerts := getBaseAlerts()
				expProps := forward.Properties{
					CustomChatID: "-1009876543210",
				}
				msvc.On("Forward", mock.Anything, expProps, expAlerts).Once().Return(nil)
			},
			expCode: http.StatusOK,
		},

		"Alertmanager webhook internal errors should be propagated to clients (forwarding).": {
			urlPath: "/alerts",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *forwardmock.Service) {
				expAlerts := getBaseAlerts()
				expProps := forward.Properties{}
				msvc.On("Forward", mock.Anything, expProps, expAlerts).Once().Return(errors.New("whatever"))
			},
			expCode: http.StatusInternalServerError,
		},

		"Alertmanager webhook configuration errors should be propagated to clients (forwarding).": {
			urlPath: "/alerts",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *forwardmock.Service) {
				expAlerts := getBaseAlerts()
				expProps := forward.Properties{}
				err := fmt.Errorf("custom error: %w", internalerrors.ErrInvalidConfiguration)
				msvc.On("Forward", mock.Anything, expProps, expAlerts).Once().Return(err)
			},
			expCode: http.StatusBadRequest,
		},

		"Alertmanager webhook configuration errors on notification should be propagated to clients (alert mapping).": {
			urlPath: "/alerts",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				wa.Version = "v3"
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock:    func(t *testing.T, msvc *forwardmock.Service) {},
			expCode: http.StatusBadRequest,
		},

		"Alertmanager webhook configuration errors on notification should be propagated to clients (JSON formatting).": {
			urlPath: "/alerts",
			webhookAlertJSON: func(t *testing.T) string {
				return "{"
			},
			mock:    func(t *testing.T, msvc *forwardmock.Service) {},
			expCode: http.StatusBadRequest,
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
			test.config.ForwardService = msvc
			h, _ := alertmanager.NewHandler(test.config)
			srv := httptest.NewServer(h)
			defer srv.Close()
			req, err := http.NewRequest(http.MethodPost, srv.URL+test.urlPath, strings.NewReader(test.webhookAlertJSON(t)))
			require.NoError(err)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(err)

			// Check.
			assert.Equal(test.expCode, resp.StatusCode)
			msvc.AssertExpectations(t)
		})
	}
}

func TestHandleDeadMansSwitch(t *testing.T) {
	tests := map[string]struct {
		config           alertmanager.Config
		urlPath          string
		webhookAlertJSON func(t *testing.T) string
		mock             func(t *testing.T, msvc *deadmansswitchmock.Service)
		expCode          int
	}{

		"Dead man's switch request should be handled correctly (with defaults).": {
			urlPath: "/alerts/dms",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *deadmansswitchmock.Service) {
				expAlerts := getBaseAlerts()
				msvc.On("PushSwitch", mock.Anything, expAlerts).Once().Return(nil)
			},
			expCode: http.StatusOK,
		},

		"Dead man's switch request should be handled correctly (with custom params).": {
			config: alertmanager.Config{
				DeadMansSwitchPath: "/dead-mans-switch",
			},
			urlPath: "/dead-mans-switch",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *deadmansswitchmock.Service) {
				expAlerts := getBaseAlerts()
				msvc.On("PushSwitch", mock.Anything, expAlerts).Once().Return(nil)
			},
			expCode: http.StatusOK,
		},

		"Dead man's switch webhook internal errors should be propagated to clients (PushSwitch).": {
			urlPath: "/alerts/dms",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *deadmansswitchmock.Service) {
				expAlerts := getBaseAlerts()
				msvc.On("PushSwitch", mock.Anything, expAlerts).Once().Return(errors.New("whatever"))
			},
			expCode: http.StatusInternalServerError,
		},

		"Dead man's switch webhook configuration errors should be propagated to clients (PushSwitch).": {
			urlPath: "/alerts/dms",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock: func(t *testing.T, msvc *deadmansswitchmock.Service) {
				expAlerts := getBaseAlerts()
				err := fmt.Errorf("custom error: %w", internalerrors.ErrInvalidConfiguration)
				msvc.On("PushSwitch", mock.Anything, expAlerts).Once().Return(err)
			},
			expCode: http.StatusBadRequest,
		},

		"Dead man's switch webhook configuration errors on notification should be propagated to clients (alert mapping).": {
			urlPath: "/alerts/dms",
			webhookAlertJSON: func(t *testing.T) string {
				wa := getBaseAlertmanagerAlerts()
				wa.Version = "v3"
				body, err := json.Marshal(wa)
				require.NoError(t, err)
				return string(body)
			},
			mock:    func(t *testing.T, msvc *deadmansswitchmock.Service) {},
			expCode: http.StatusBadRequest,
		},

		"Dead man's switch configuration errors on notification should be propagated to clients (JSON formatting).": {
			urlPath: "/alerts/dms",
			webhookAlertJSON: func(t *testing.T) string {
				return "{"
			},
			mock:    func(t *testing.T, msvc *deadmansswitchmock.Service) {},
			expCode: http.StatusBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Mocks.
			mfw := &forwardmock.Service{}
			mdms := &deadmansswitchmock.Service{}
			test.mock(t, mdms)

			// Execute.
			test.config.ForwardService = mfw
			test.config.DeadMansSwitchService = mdms
			h, err := alertmanager.NewHandler(test.config)
			require.NoError(err)
			srv := httptest.NewServer(h)
			defer srv.Close()
			req, err := http.NewRequest(http.MethodPost, srv.URL+test.urlPath, strings.NewReader(test.webhookAlertJSON(t)))
			require.NoError(err)
			resp, err := http.DefaultClient.Do(req)
			require.NoError(err)

			// Check.
			assert.Equal(test.expCode, resp.StatusCode)
			mfw.AssertExpectations(t)
			mdms.AssertExpectations(t)
		})
	}
}
