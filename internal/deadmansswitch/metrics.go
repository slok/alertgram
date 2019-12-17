package deadmansswitch

import (
	"context"
	"time"

	"github.com/slok/alertgram/internal/model"
)

// ServiceMetricsRecorder knows how to record metrics on deadmansswitch.Service.
type ServiceMetricsRecorder interface {
	ObserveDMSServiceOpDuration(ctx context.Context, op string, success bool, t time.Duration)
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

func (m measureService) PushSwitch(ctx context.Context, ag *model.AlertGroup) (err error) {
	defer func(t0 time.Time) {
		m.rec.ObserveDMSServiceOpDuration(ctx, "PushSwitch", err == nil, time.Since(t0))
	}(time.Now())
	return m.next.PushSwitch(ctx, ag)
}
