package forward

import (
	"context"

	"github.com/slok/alertgram/internal/log"
	"github.com/slok/alertgram/internal/model"
)

// Service is the domain service that forwards alerts
type Service interface {
	// Forward knows how to forward alerts from an input to an output.
	Forward(ctx context.Context, alert model.Alert) error
}

type service struct {
	notifiers []Notifier
	logger    log.Logger
}

// NewService returns a new forward.Service.
func NewService(notifiers []Notifier, l log.Logger) Service {
	return &service{
		notifiers: notifiers,
		logger:    l.WithData(log.KV{"service": "forward.Service"}),
	}
}

func (s service) Forward(ctx context.Context, alert model.Alert) error {
	// TODO(slok): Validate alert.
	// TODO(slok): Add concurrency using workers.
	for _, not := range s.notifiers {
		err := not.Notify(ctx, alert)
		if err != nil {
			s.logger.WithData(log.KV{"notifier": not.Type(), "alertID": alert.ID}).Errorf("could not notify alert")
		}
	}

	return nil
}
