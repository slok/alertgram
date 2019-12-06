package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/slok/alertgram/internal/log"
)

const (
	drainTimeoutDef  = 2 * time.Second
	listenAddressDef = ":8080"
)

// Config is the server configuration.
type Config struct {
	// ListenAddress is where the server will be listening.
	// By default will listen on :8080.
	ListenAddress string
	// DrainTimeout is the draining timeout, by default is 2 seconds.
	DrainTimeout time.Duration
	// Handler is the handler that will serve the server.
	Handler http.Handler
	// Logger is the logger used by the server.
	Logger log.Logger
}

func (c *Config) defaults() error {
	if c.Handler == nil {
		return fmt.Errorf("handler is required")
	}

	if c.Logger == nil {
		c.Logger = log.Dummy
	}

	if c.ListenAddress == "" {
		c.ListenAddress = listenAddressDef
	}

	if c.DrainTimeout == 0 {
		c.DrainTimeout = drainTimeoutDef
	}

	return nil
}

// Server is a Server that serves a handler.
type Server struct {
	server        *http.Server
	listenAddress string
	drainTimeout  time.Duration
	logger        log.Logger
}

// NewServer returns a new HTTP server.
func NewServer(cfg Config) (*Server, error) {
	// fulfill with default configuration if needed.
	err := cfg.defaults()
	if err != nil {
		return nil, err
	}

	// Create the handler mux and the internal http server.
	httpServer := &http.Server{
		Handler: cfg.Handler,
		Addr:    cfg.ListenAddress,
	}

	// Create our HTTP Server.
	return &Server{
		server:        httpServer,
		listenAddress: cfg.ListenAddress,
		drainTimeout:  cfg.DrainTimeout,
		logger: cfg.Logger.WithValues(log.KV{
			"service": "http-server",
			"addr":    cfg.ListenAddress,
		}),
	}, nil
}

// ListenAndServe runs the server.
func (s *Server) ListenAndServe() error {
	s.logger.Infof("server listening on %s...", s.listenAddress)
	return s.server.ListenAndServe()
}

// DrainAndShutdown will drain the connections and shutdown the server.
func (s *Server) DrainAndShutdown() error {
	s.logger.Infof("start draining connections...")

	ctx, cancel := context.WithTimeout(context.Background(), s.drainTimeout)
	defer cancel()

	err := s.server.Shutdown(ctx)
	if err != nil {
		return fmt.Errorf("something happened while draining connections: %w", err)
	}

	s.logger.Infof("connections drained")
	return nil
}
