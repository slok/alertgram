package forward

import (
	"context"

	"github.com/slok/alertgram/internal/model"
)

// Notification is the notification that wants to be sed
// via a notifier.
type Notification struct {
	// ChatID is an ID to send the notification. In
	// Telegram could be a channel/group ID, in Slack
	// a room or a user.
	ChatID     string
	AlertGroup model.AlertGroup
}

// Notifier knows how to notify alerts to different backends.
type Notifier interface {
	Notify(ctx context.Context, notification Notification) error
	Type() string
}
