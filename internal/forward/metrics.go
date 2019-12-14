package forward

import (
	"context"
	"time"

	"github.com/slok/alertgram/internal/model"
)

// ServiceMetricsRecorder knows how to record metrics on forward.Service.
type ServiceMetricsRecorder interface {
	ObserveForwardServiceOpDuration(ctx context.Context, op string, success bool, t time.Duration)
}

type measureService struct {
	rec  ServiceMetricsRecorder
	next Service
}

// NewMeasureService wraps a service and measures using metrics.
func NewMeasureService(rec ServiceMetricsRecorder, next Service) Service {
	return &measureService{
		rec:  rec,
		next: next,
	}
}

func (m measureService) Forward(ctx context.Context, props Properties, ag *model.AlertGroup) (err error) {
	defer func(t0 time.Time) {
		m.rec.ObserveForwardServiceOpDuration(ctx, "Forward", err == nil, time.Since(t0))
	}(time.Now())
	return m.next.Forward(ctx, props, ag)
}

// NotifierMetricsRecorder knows how to record metrics on forward.Notifier.
type NotifierMetricsRecorder interface {
	ObserveForwardNotifierOpDuration(ctx context.Context, notifierType string, op string, success bool, t time.Duration)
}

type measureNotifier struct {
	notifierType string
	rec          NotifierMetricsRecorder
	next         Notifier
}

// NewMeasureNotifier wraps a notifier and measures using metrics.
func NewMeasureNotifier(rec NotifierMetricsRecorder, next Notifier) Notifier {
	return &measureNotifier{
		notifierType: next.Type(),
		rec:          rec,
		next:         next,
	}
}

func (m measureNotifier) Notify(ctx context.Context, n Notification) (err error) {
	defer func(t0 time.Time) {
		m.rec.ObserveForwardNotifierOpDuration(ctx, m.notifierType, "Notify", err == nil, time.Since(t0))
	}(time.Now())
	return m.next.Notify(ctx, n)
}

func (m measureNotifier) Type() string {
	defer func(t0 time.Time) {
		m.rec.ObserveForwardNotifierOpDuration(context.TODO(), m.notifierType, "Type", true, time.Since(t0))
	}(time.Now())
	return m.next.Type()
}
