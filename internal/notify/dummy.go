package notify

import (
	"context"

	"github.com/slok/alertgram/internal/model"
)

type dummy int

// Dummy is a dummy notifier.
const Dummy = dummy(0)

func (dummy) Notify(ctx context.Context, alertGroup *model.AlertGroup) error { return nil }
func (dummy) Type() string                                                   { return "dummy" }
