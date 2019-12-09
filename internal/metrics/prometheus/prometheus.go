package prometheus

import (
	"context"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	httpmetrics "github.com/slok/go-http-metrics/metrics"
	httpmetricsprometheus "github.com/slok/go-http-metrics/metrics/prometheus"

	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/notify"
)

const prefix = "alertgram"

// Recorder knows how to measure the different metrics
// interfaces of the application.
type Recorder struct {
	httpmetrics.Recorder

	forwardServiceOpDurHistogram   *prometheus.HistogramVec
	forwardNotifierOpDurHistogram  *prometheus.HistogramVec
	templateRendererOpDurHistogram *prometheus.HistogramVec
}

// New returns a new Prometheus recorder for the app.
func New(reg prometheus.Registerer) *Recorder {
	r := &Recorder{
		Recorder: httpmetricsprometheus.NewRecorder(httpmetricsprometheus.Config{
			Registry: reg,
		}),

		forwardServiceOpDurHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "forward",
			Name:      "operation_duration_seconds",
			Help:      "The duration of the operation in forward service.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"operation", "success"}),

		forwardNotifierOpDurHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "notifier",
			Name:      "operation_duration_seconds",
			Help:      "The duration of the operation in notifier.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"type", "operation", "success"}),

		templateRendererOpDurHistogram: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: prefix,
			Subsystem: "template_renderer",
			Name:      "operation_duration_seconds",
			Help:      "The duration of the operation in template renderer.",
			Buckets:   prometheus.DefBuckets,
		}, []string{"type", "operation", "success"}),
	}

	// Register all the metrics.
	reg.MustRegister(
		r.forwardServiceOpDurHistogram,
		r.forwardNotifierOpDurHistogram,
		r.templateRendererOpDurHistogram,
	)

	return r
}

// ObserveForwardNotifierOpDuration satifies forward.NotifierMetricsRecorder interface.
func (r Recorder) ObserveForwardNotifierOpDuration(ctx context.Context, notType string, op string, success bool, t time.Duration) {
	r.forwardNotifierOpDurHistogram.WithLabelValues(notType, op, strconv.FormatBool(success)).Observe(t.Seconds())
}

// ObserveForwardServiceOpDuration satisfies forward.ServiceMetricsRecorder interface.
func (r Recorder) ObserveForwardServiceOpDuration(ctx context.Context, op string, success bool, t time.Duration) {
	r.forwardServiceOpDurHistogram.WithLabelValues(op, strconv.FormatBool(success)).Observe(t.Seconds())
}

// ObserveTemplateRendererOpDuration satisfies notify.TemplateRendererMetricsRecorder interface.
func (r Recorder) ObserveTemplateRendererOpDuration(ctx context.Context, rendererType string, op string, success bool, t time.Duration) {
	r.templateRendererOpDurHistogram.WithLabelValues(rendererType, op, strconv.FormatBool(success)).Observe(t.Seconds())
}

// Ensure that the recorder implements the different interfaces of the app.
var _ forward.NotifierMetricsRecorder = &Recorder{}
var _ forward.ServiceMetricsRecorder = &Recorder{}
var _ notify.TemplateRendererMetricsRecorder = &Recorder{}
var _ httpmetrics.Recorder = &Recorder{}
