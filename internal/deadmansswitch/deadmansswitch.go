package deadmansswitch

import (
	"context"
	"fmt"
	"time"

	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/log"
	"github.com/slok/alertgram/internal/model"
)

// Service is a Dead man's switch
//
// A dead man's switch is a process where at regular intervals if some kind of signal is
// not received it will be activated. This usually is used to check that some kind
// of system is working, in this case if we don't receive an alert we assume that something
// is not working and we should notify.
type Service interface {
	// PushSwitch will disable the dead man's switch when it's pushed and reset
	// the interval for activation.
	PushSwitch(ctx context.Context, alertGroup *model.AlertGroup) error
}

// Config is the Service configuration.
type Config struct {
	CustomChatID string
	Interval     time.Duration
	Notifiers    []forward.Notifier
	Logger       log.Logger
}

func (c *Config) defaults() error {
	if c.Logger == nil {
		c.Logger = log.Dummy
	}
	return nil
}

type service struct {
	cfg       Config
	dmsSwitch chan *model.AlertGroup
	notifiers []forward.Notifier
	logger    log.Logger
}

// NewService returns a Dead mans's switch service.
// When creating a new instance it will start the dead man's switch interval
// it can only stop once and it's done when the received context is done.
func NewService(ctx context.Context, cfg Config) (Service, error) {
	err := cfg.defaults()
	if err != nil {
		return nil, fmt.Errorf("invalid dead man's switch service configuration: %w", err)
	}
	s := &service{
		cfg:       cfg,
		dmsSwitch: make(chan *model.AlertGroup, 1),
		notifiers: cfg.Notifiers,
		logger:    cfg.Logger.WithValues(log.KV{"service": "deadMansSwitch"}),
	}
	go s.startDMS(ctx)

	return s, nil
}

func (s *service) PushSwitch(_ context.Context, alertGroup *model.AlertGroup) error {
	if alertGroup != nil {
		s.dmsSwitch <- alertGroup
	}
	return nil
}

func (s *service) activate(ctx context.Context) error {
	dmsNotification := forward.Notification{
		ChatID: s.cfg.CustomChatID,
		AlertGroup: model.AlertGroup{
			ID: "DeadMansSwitchActive",
			Alerts: []model.Alert{
				model.Alert{
					ID:       "DeadMansSwitchActive",
					Name:     "DeadMansSwitchActive",
					StartsAt: time.Now(),
					Status:   model.AlertStatusFiring,
					Labels: map[string]string{
						"alertname": "DeadMansSwitchActive",
						"severity":  "critical",
						"origin":    "alertgram",
					},
					Annotations: map[string]string{
						"message": "The Dead man's switch has been activated! This usually means that your monitoring/alerting system is not working",
					},
				},
			},
		},
	}

	// TODO(slok): Add concurrency using workers.
	for _, not := range s.notifiers {
		err := not.Notify(ctx, dmsNotification)
		if err != nil {
			s.logger.WithValues(log.KV{"notifier": not.Type(), "alertGroupID": dmsNotification.AlertGroup.ID}).
				Errorf("could not notify alert group: %s", err)
		}
	}
	return nil
}

// startDMS will start the DeadMansSwitch process.
// It will be listening to the signals to know
// that we are alive, if not received in the interval the
// Dead mans switch should assume we are dead and will activate
// this means executing the received function.
func (s *service) startDMS(ctx context.Context) {
	logger := s.logger.WithValues(log.KV{"interval": s.cfg.Interval})
	logger.Infof("dead man's switch started with an interval of %s", s.cfg.Interval)

	for {
		select {
		case <-ctx.Done():
			logger.Infof("context done, stopping dead man's switch")
			return
		case <-time.After(s.cfg.Interval):
			logger.Infof("no switch pushed during interval wait, dead mans switch activated!")
			err := s.activate(ctx)
			if err != nil {
				logger.Errorf("something happened when activating the dead man's switch")
			}
		case <-s.dmsSwitch:
			logger.Debugf("dead mans switch pushed, deactivated")
		}
	}
}

// DisabledService is a Dead man switch service that doesn't do anything.
const DisabledService = dummyService(0)

type dummyService int

func (dummyService) PushSwitch(ctx context.Context, alertGroup *model.AlertGroup) error { return nil }
