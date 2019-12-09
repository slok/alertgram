package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/oklog/run"

	"github.com/slok/alertgram/internal/forward"
	internalhttp "github.com/slok/alertgram/internal/http"
	"github.com/slok/alertgram/internal/http/alertmanager"
	"github.com/slok/alertgram/internal/log"
	"github.com/slok/alertgram/internal/log/logrus"
	"github.com/slok/alertgram/internal/notify"
	"github.com/slok/alertgram/internal/notify/telegram"
)

// Main is the main application.
type Main struct {
	cfg    *Config
	logger log.Logger
}

// Run runs the main application.
func (m *Main) Run() error {
	// Initialization and setup.
	var err error
	m.cfg, err = NewConfig()
	if err != nil {
		return err
	}
	m.logger = logrus.New(m.cfg.DebugMode)

	// Dependencies.
	tgCli, err := tgbotapi.NewBotAPI(m.cfg.TeletramAPIToken)
	if err != nil {
		return err
	}
	tplRenderer := notify.DefaultTemplateRenderer

	var notifier forward.Notifier
	if m.cfg.NotifyDryRun {
		notifier = notify.NewLogger(tplRenderer, m.logger)
	} else {
		notifier, err = telegram.NewNotifier(telegram.Config{
			TemplateRenderer:      tplRenderer,
			Client:                tgCli,
			DefaultTelegramChatID: m.cfg.TelegramChatID,
			Logger:                m.logger,
		})
		if err != nil {
			return err
		}
	}

	// Domain services.
	forwardSvc := forward.NewService([]forward.Notifier{notifier}, m.logger)

	var g run.Group

	// Alertmanager webhook server.
	{
		logger := m.logger.WithValues(log.KV{"server": "alertmanager-handler"})
		h, err := alertmanager.NewHandler(alertmanager.Config{
			WebhookPath: m.cfg.AlertmanagerWebhookPath,
			Forwarder:   forwardSvc,
			Logger:      logger,
		})
		if err != nil {
			return err
		}
		server, err := internalhttp.NewServer(internalhttp.Config{
			Handler:       h,
			ListenAddress: m.cfg.AlertmanagerListenAddr,
			Logger:        logger,
		})
		if err != nil {
			return err
		}

		g.Add(
			func() error {
				return server.ListenAndServe()
			},
			func(_ error) {
				if err := server.DrainAndShutdown(); err != nil {
					logger.Errorf("error while draining connections")
				}
			})
	}

	// Capture signals.
	{
		logger := m.logger.WithValues(log.KV{"service": "main"})
		sigC := make(chan os.Signal, 1)
		exitC := make(chan struct{})
		signal.Notify(sigC, syscall.SIGTERM, syscall.SIGINT)
		g.Add(
			func() error {
				select {
				case <-sigC:
					logger.Infof("signal captured")
				case <-exitC:
				}
				return nil
			},
			func(_ error) {
				close(exitC)
			})
	}

	return g.Run()
}

func main() {
	m := Main{}
	if err := m.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error running the app: %s", err)
		os.Exit(1)
	}

	fmt.Println("goodbye!")
}
