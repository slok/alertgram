package notify

import (
	"context"
	"time"

	"github.com/slok/alertgram/internal/model"
)

// TemplateRendererMetricsRecorder knows how to record the metrics in the TemplateRenderer.
type TemplateRendererMetricsRecorder interface {
	ObserveTemplateRendererOpDuration(ctx context.Context, rendererType string, op string, success bool, t time.Duration)
}

type measureTemplateRenderer struct {
	rendererType string
	rec          TemplateRendererMetricsRecorder
	next         TemplateRenderer
}

// NewMeasureTemplateRenderer wraps a template renderer and measures using metrics.
func NewMeasureTemplateRenderer(rendererType string, rec TemplateRendererMetricsRecorder, next TemplateRenderer) TemplateRenderer {
	return &measureTemplateRenderer{
		rendererType: rendererType,
		rec:          rec,
		next:         next,
	}
}

func (m measureTemplateRenderer) Render(ctx context.Context, ag *model.AlertGroup) (_ string, err error) {
	defer func(t0 time.Time) {
		m.rec.ObserveTemplateRendererOpDuration(ctx, m.rendererType, "Render", err == nil, time.Since(t0))
	}(time.Now())
	return m.next.Render(ctx, ag)
}
