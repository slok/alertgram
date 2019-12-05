package alertmanager

import (
	"errors"

	"github.com/prometheus/alertmanager/notify/webhook"
	prommodel "github.com/prometheus/common/model"

	"github.com/slok/alertgram/internal/model"
)

var (
	// ErrCantDeserialize will be used when for some reason the
	// received data can't be deserialized.
	ErrCantDeserialize = errors.New("cant deserialize the received alerts")
)

// alertGroupV4 are the alertGroup received by the webhook
// It uses the V4 version of the webhook format.
//
// https://github.com/prometheus/alertmanager/blob/5cb556e4b2247f2c5d8cebdef88e2a634a46863a/notify/webhook/webhook.go#L85
type alertGroupV4 webhook.Message

func (a alertGroupV4) toDomain() (*model.AlertGroup, error) {
	// Map alerts.
	alerts := make([]model.Alert, len(a.Alerts))
	for i := 0; i < len(a.Alerts); i++ {
		alert := a.Alerts[i]
		alerts[i] = model.Alert{
			ID:          alert.Fingerprint,
			Name:        alert.Labels[prommodel.AlertNameLabel],
			Start:       alert.StartsAt,
			End:         alert.EndsAt,
			Status:      alertStatusToDomain(alert.Status),
			Labels:      alert.Labels,
			Annotations: alert.Annotations,
		}
	}

	ag := &model.AlertGroup{
		ID:          a.GroupKey,
		Status:      alertStatusToDomain(a.Status),
		Labels:      a.CommonLabels,
		Annotations: a.CommonAnnotations,
		Alerts:      alerts,
	}

	return ag, nil
}

func alertStatusToDomain(st string) model.AlertStatus {
	switch prommodel.AlertStatus(st) {
	case prommodel.AlertFiring:
		return model.AlertStatusFiring
	case prommodel.AlertResolved:
		return model.AlertStatusResolved
	default:
		return model.AlertStatusUnknown
	}
}
