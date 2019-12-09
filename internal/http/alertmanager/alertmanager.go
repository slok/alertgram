package alertmanager

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slok/go-http-metrics/metrics"
	metricsmiddleware "github.com/slok/go-http-metrics/middleware"
	metricsmiddlewaregin "github.com/slok/go-http-metrics/middleware/gin"

	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/log"
)

// Config is the configuration of the WebhookHandler.
type Config struct {
	MetricsRecorder metrics.Recorder
	WebhookPath     string
	Forwarder       forward.Service
	Debug           bool
	Logger          log.Logger
}

func (c *Config) defaults() error {
	if c.WebhookPath == "" {
		c.WebhookPath = "/alerts"
	}

	if c.Forwarder == nil {
		return fmt.Errorf("forward can't be nil")
	}

	if c.Logger == nil {
		c.Logger = log.Dummy
	}

	return nil
}

// More info here: https://prometheus.io/docs/alerting/configuration/#webhook_config.
type webhookHandler struct {
	cfg       Config
	engine    *gin.Engine
	forwarder forward.Service
	logger    log.Logger
}

// NewHandler is an HTTP handler that knows how to handle
// alertmanager webhook alerts.
func NewHandler(cfg Config) (http.Handler, error) {
	err := cfg.defaults()
	if err != nil {
		return nil, err
	}

	if !cfg.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	w := webhookHandler{
		cfg:       cfg,
		engine:    gin.New(),
		forwarder: cfg.Forwarder,
		logger:    cfg.Logger,
	}

	// Metrics middleware.
	mdlw := metricsmiddleware.New(metricsmiddleware.Config{
		Recorder: cfg.MetricsRecorder,
	})
	w.engine.Use(metricsmiddlewaregin.Handler("", mdlw))

	// Register routes.
	w.routes()

	return w.engine, nil
}

func (w webhookHandler) routes() {
	w.engine.POST(w.cfg.WebhookPath, w.HandleAlerts())
}
