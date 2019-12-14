package forward

import (
	"context"
	"errors"
	"fmt"

	"github.com/slok/alertgram/internal/log"
	"github.com/slok/alertgram/internal/model"
)

// Properties are the  properties an AlertGroup can have
// when the forwarding process is done.
type Properties struct {
	// CustomChatID can be used when the forward should be done
	// to a different target (chat, group, channel, user...)
	// instead of using the default one.
	CustomChatID string
}

// Service is the domain service that forwards alerts
type Service interface {
	// Forward knows how to forward alerts from an input to an output.
	Forward(ctx context.Context, props Properties, alertGroup *model.AlertGroup) error
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

func (s service) Forward(ctx context.Context, props Properties, alertGroup *model.AlertGroup) error {
	// TODO(slok): Add better validation.
	if alertGroup == nil {
		return fmt.Errorf("alertgroup can't be empty: %w", ErrInvalidAlertGroup)
	}

	// TODO(slok): Add concurrency using workers.
	notification := Notification{
		AlertGroup: *alertGroup,
		ChatID:     props.CustomChatID,
	}
	for _, not := range s.notifiers {
		err := not.Notify(ctx, notification)
		if err != nil {
			s.logger.WithValues(log.KV{"notifier": not.Type(), "alertGroupID": alertGroup.ID}).
				Errorf("could not notify alert group: %s", err)
		}
	}

	return nil
}
