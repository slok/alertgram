package alertmanager

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/slok/go-http-metrics/metrics"
	metricsmiddleware "github.com/slok/go-http-metrics/middleware"
	metricsmiddlewaregin "github.com/slok/go-http-metrics/middleware/gin"

	"github.com/slok/alertgram/internal/deadmansswitch"
	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/log"
)

// Config is the configuration of the WebhookHandler.
type Config struct {
	MetricsRecorder       metrics.Recorder
	WebhookPath           string
	ChatIDQueryString     string
	ForwardService        forward.Service
	DeadMansSwitchPath    string
	DeadMansSwitchService deadmansswitch.Service
	Debug                 bool
	Logger                log.Logger
}

func (c *Config) defaults() error {
	if c.WebhookPath == "" {
		c.WebhookPath = "/alerts"
	}

	if c.ForwardService == nil {
		return fmt.Errorf("forward can't be nil")
	}

	if c.ChatIDQueryString == "" {
		c.ChatIDQueryString = "chat-id"
	}

	if c.DeadMansSwitchService == nil {
		c.DeadMansSwitchService = deadmansswitch.DisabledService
	}

	if c.DeadMansSwitchPath == "" {
		c.DeadMansSwitchPath = "/alerts/dms"
	}

	if c.Logger == nil {
		c.Logger = log.Dummy
	}

	return nil
}

// More info here: https://prometheus.io/docs/alerting/configuration/#webhook_config.
type webhookHandler struct {
	cfg              Config
	engine           *gin.Engine
	forwarder        forward.Service
	deadmansswitcher deadmansswitch.Service
	logger           log.Logger
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
		cfg:              cfg,
		engine:           gin.New(),
		forwarder:        cfg.ForwardService,
		deadmansswitcher: cfg.DeadMansSwitchService,
		logger:           cfg.Logger,
	}

	// Metrics middleware.
	mdlw := metricsmiddleware.New(metricsmiddleware.Config{
		Service:  "alertmanager-api",
		Recorder: cfg.MetricsRecorder,
	})
	w.engine.Use(metricsmiddlewaregin.Handler("", mdlw))

	// Register routes.
	w.routes()

	return w.engine, nil
}

func (w webhookHandler) routes() {
	w.engine.POST(w.cfg.WebhookPath, w.HandleAlerts())

	// Only enable dead man's switch if required.
	if w.deadmansswitcher != nil && w.deadmansswitcher != deadmansswitch.DisabledService {
		w.engine.POST(w.cfg.DeadMansSwitchPath, w.HandleDeadMansSwitch())
	}
}
