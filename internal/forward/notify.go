package forward

import (
	"context"

	"github.com/slok/alertgram/internal/model"
)

// Notifier knows how to notify alerts to different backends.
type Notifier interface {
	Notify(ctx context.Context, alert model.Alert) error
	Type() string
}
