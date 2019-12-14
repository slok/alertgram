package forward

import (
	"context"
	"errors"
	"fmt"

	"github.com/slok/alertgram/internal/log"
	"github.com/slok/alertgram/internal/model"
)

// Service is the domain service that forwards alerts
type Service interface {
	// Forward knows how to forward alerts from an input to an output.
	Forward(ctx context.Context, alertGroup *model.AlertGroup) error
}

type service struct {
	notifiers []Notifier
	logger    log.Logger
}

// NewService returns a new forward.Service.
func NewService(notifiers []Notifier, l log.Logger) Service {
	return &service{
		notifiers: notifiers,
		logger:    l.WithValues(log.KV{"service": "forward.Service"}),
	}
}

var (
	// ErrInvalidAlertGroup will be used when the alertgroup is not valid.
	ErrInvalidAlertGroup = errors.New("invalid alert group")
)

func (s service) Forward(ctx context.Context, alertGroup *model.AlertGroup) error {
	// TODO(slok): Add better validation.
	if alertGroup == nil {
		return fmt.Errorf("alertgroup can't be empty: %w", ErrInvalidAlertGroup)
	}

	// TODO(slok): Add concurrency using workers.
	for _, not := range s.notifiers {
		err := not.Notify(ctx, Notification{AlertGroup: *alertGroup})
		if err != nil {
			s.logger.WithValues(log.KV{"notifier": not.Type(), "alertGroupID": alertGroup.ID}).
				Errorf("could not notify alert group: %s", err)
		}
	}

	return nil
}
