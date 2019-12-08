package telegram

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html/template"

	"github.com/Masterminds/sprig/v3"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/slok/alertgram/internal/forward"
	"github.com/slok/alertgram/internal/internalerrors"
	"github.com/slok/alertgram/internal/log"
	"github.com/slok/alertgram/internal/model"
)

var (
	// ErrComm will be used when the communication to telegram fails.
	ErrComm = errors.New("error communicating with telegram")
	// ErrRenderTemplate will be used when there is an error rendering the alerts
	// to a template.
	ErrRenderTemplate = errors.New("error rendering template")
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
	// AlertMessageTemplate is the template that will be used to render the final telegram
	// message.
	// The templates use https://github.com/Masterminds/sprig to render.
	// The templates will be rendered and send to telegram in HTML parse mode.
	// If not set it will use the default simple one.
	AlertMessageTemplate string
	// DryRun will not send the notification and just print in the terminal.
	DryRun bool
	// Client is the telegram client is compatible with "github.com/go-telegram-bot-api/telegram-bot-api"
	// library client API.
	Client Client
	// Logger is the logger.
	Logger log.Logger
}

func (c *Config) defaults() error {
	if c.AlertMessageTemplate == "" {
		c.AlertMessageTemplate = defTemplate
	}

	if c.Client == nil {
		return fmt.Errorf("telegram client is required")
	}

	if c.Logger == nil {
		c.Logger = log.Dummy
	}

	return nil
}

type notifier struct {
	msgTpl *template.Template
	cfg    Config
	client Client
	logger log.Logger
}

// NewNotifier returns a notifier is a Telegram notifier
// that knows how to send alerts to telegram.
func NewNotifier(cfg Config) (forward.Notifier, error) {
	err := cfg.defaults()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", err, internalerrors.ErrInvalidConfiguration)
	}

	// Load the telegram message template.
	tpl, err := template.New("alerts").Funcs(sprig.FuncMap()).Parse(cfg.AlertMessageTemplate)
	if err != nil {
		return nil, fmt.Errorf("error parsing alert template: %s: %w", err, internalerrors.ErrInvalidConfiguration)
	}

	return &notifier{
		cfg:    cfg,
		msgTpl: tpl,
		client: cfg.Client,
		logger: cfg.Logger.WithValues(log.KV{"notifier": "telegram"}),
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

	msg, err := n.alertGroupToMessage(alertGroup)
	if err != nil {
		return fmt.Errorf("could not format the alerts to message: %w", err)
	}

	if n.cfg.DryRun {
		logger.Infof("dry run message to telegram:	\n%+v", msg)
		return nil
	}

	res, err := n.client.Send(msg)
	if err != nil {
		err = fmt.Errorf("%w:  %s", ErrComm, err)
		return fmt.Errorf("error sending telegram message: %w", err)
	}
	logger.Debugf("telegram response: %+v", res)

	return nil
}

func (n notifier) alertGroupToMessage(a *model.AlertGroup) (tgbotapi.Chattable, error) {
	var b bytes.Buffer
	err := n.msgTpl.Execute(&b, a)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrRenderTemplate, err)
		return nil, fmt.Errorf("error rendering alerts to template: %w", err)
	}

	msg := tgbotapi.NewMessage(n.cfg.DefaultTelegramChatID, b.String())
	msg.ParseMode = "HTML"
	return msg, nil
}

func (n notifier) Type() string { return "telegram" }

// Client is an small abstraction for the telegram-bot-api client.
// the client of this lib should satisfy directly.
// More info here: https://godoc.org/github.com/go-telegram-bot-api/telegram-bot-api.
type Client interface {
	Send(c tgbotapi.Chattable) (tgbotapi.Message, error)
}
