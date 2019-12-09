package notify

import (
	"context"

	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/log"
	"github.com/slok/alertgram/internal/model"
)

type dummy int

// Dummy is a dummy notifier.
const Dummy = dummy(0)

func (dummy) Notify(ctx context.Context, alertGroup *model.AlertGroup) error { return nil }
func (dummy) Type() string                                                   { return "dummy" }

type logger struct {
	renderer TemplateRenderer
	logger   log.Logger
}

// NewLogger returns a notifier that only logs the renderer alerts,
// normally used to develop or dry/run.
func NewLogger(r TemplateRenderer, l log.Logger) forward.Notifier {
	return &logger{
		renderer: r,
		logger:   l.WithValues(log.KV{"notifier": "logger"}),
	}
}

func (l logger) Notify(ctx context.Context, alertGroup *model.AlertGroup) error {
	logger := l.logger.WithValues(log.KV{"alertGroup": alertGroup.ID, "alertsNumber": len(alertGroup.Alerts)})

	alertText, err := l.renderer.Render(ctx, alertGroup)
	if err != nil {
		return err
	}
	logger.Infof("alert: %s", alertText)

	return nil
}
func (logger) Type() string { return "logger" }
