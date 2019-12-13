package telegram

import (
	"context"
	"errors"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/internalerrors"
	"github.com/slok/alertgram/internal/log"
	"github.com/slok/alertgram/internal/model"
	"github.com/slok/alertgram/internal/notify"
)

var (
	// ErrComm will be used when the communication to telegram fails.
	ErrComm = errors.New("error communicating with telegram")
)

// Config is the configuration of the Notifier.
type Config struct {
	// This is the ID of the channel or group where the alerts
	// will be sent by default.
	// Got from here https://github.com/GabrielRF/telegram-id#web-channel-id
	// You ca get the the ID like this:
	// - Enter the telegram web app and there to the channel/group.
	// - Check the URL, it has this schema: https://web.telegram.org/#/im?p=c1234567891_12345678912345678912
	// - Get the `c1234567891_`, get this part: `1234567891`.
	// - Add `-100` (until you have 13 characters), this should be the chat ID: `-1001234567891`
	DefaultTelegramChatID int64
	// TemplateRenderer is the renderer that will be used to render the
	// notifications before sending to Telegram.
	TemplateRenderer notify.TemplateRenderer
	// Client is the telegram client is compatible with "github.com/go-telegram-bot-api/telegram-bot-api"
	// library client API.
	Client Client
	// Logger is the logger.
	Logger log.Logger
}

func (c *Config) defaults() error {
	if c.Client == nil {
		return fmt.Errorf("telegram client is required")
	}

	if c.TemplateRenderer == nil {
		c.TemplateRenderer = notify.DefaultTemplateRenderer
	}

	if c.Logger == nil {
		c.Logger = log.Dummy
	}

	return nil
}

type notifier struct {
	tplRenderer notify.TemplateRenderer
	cfg         Config
	client      Client
	logger      log.Logger
}

// NewNotifier returns a notifier is a Telegram notifier
// that knows how to send alerts to telegram.
func NewNotifier(cfg Config) (forward.Notifier, error) {
	err := cfg.defaults()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, internalerrors.ErrInvalidConfiguration)
	}

	return &notifier{
		cfg:         cfg,
		tplRenderer: cfg.TemplateRenderer,
		client:      cfg.Client,
		logger:      cfg.Logger.WithValues(log.KV{"notifier": "telegram"}),
	}, nil
}

func (n notifier) Notify(ctx context.Context, alertGroup *model.AlertGroup) error {
	logger := n.logger.WithValues(log.KV{"alertGroup": alertGroup.ID, "alertsNumber": len(alertGroup.Alerts)})
	select {
	case <-ctx.Done():
		logger.Infof("context cancelled, not notifying alerts")
		return nil
	default:
	}

	msg, err := n.alertGroupToMessage(ctx, alertGroup)
	if err != nil {
		return fmt.Errorf("could not format the alerts to message: %w", err)
	}

	res, err := n.client.Send(msg)
	if err != nil {
		err = fmt.Errorf("%w:  %s", ErrComm, err)
		return fmt.Errorf("error sending telegram message: %w", err)
	}
	logger.Infof("telegram message sent")
	logger.Debugf("telegram response: %+v", res)

	return nil
}

func (n notifier) alertGroupToMessage(ctx context.Context, a *model.AlertGroup) (tgbotapi.Chattable, error) {
	data, err := n.tplRenderer.Render(ctx, a)
	if err != nil {
		return nil, fmt.Errorf("error rendering alerts to template: %w", err)
	}

	msg := tgbotapi.NewMessage(n.cfg.DefaultTelegramChatID, data)
	msg.ParseMode = "HTML"
	msg.DisableWebPagePreview = true // TODO(slok): Make it configurable?
	return msg, nil
}

func (n notifier) Type() string { return "telegram" }

// Client is an small abstraction for the telegram-bot-api client.
// the client of this lib should satisfy directly.
// More info here: https://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api.
type Client interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
