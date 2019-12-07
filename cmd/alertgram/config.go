package main

import (
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Version will be populated in compilation time.
	Version = "dev"
)

// flag descriptions.
const (
	descAMListenAddr      = "The listen address where the server will be listening to alertmanager's webhook request."
	descAMWebhookPath     = "The path where the server will be handling the alertmanager webhook alert requests."
	descTelegramAPIToken  = "The token that will be used to use the telegram API to send the alerts."
	descTelegramDefChatID = "The default ID of the chat (group/channel) in telegram where the alerts will be sent."
	descDebug             = "Run the application in debug mode."
)

const (
	defAMListenAddr  = ":8080"
	defAMWebhookPath = "/alerts"
)

// Config has the configuration of the application.
type Config struct {
	AlertmanagerListenAddr  string
	AlertmanagerWebhookPath string
	TeletramAPIToken        string
	TelegramChatID          int64
	DebugMode               bool

	app *kingpin.Application
}

// NewConfig returns a new configuration for the apps.
func NewConfig() (*Config, error) {
	c := &Config{
		app: kingpin.New("alertgram", "Forward your alerts to telegram.").DefaultEnvars(),
	}
	c.app.Version(Version)

	c.registerFlags()

	if _, err := c.app.Parse(os.Args[1:]); err != nil {
		return nil, err
	}
	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) registerFlags() {
	c.app.Flag("alertmanager.listen-address", descAMListenAddr).Default(defAMListenAddr).StringVar(&c.AlertmanagerListenAddr)
	c.app.Flag("alertmanager.webhook-path", descAMWebhookPath).Default(defAMWebhookPath).StringVar(&c.AlertmanagerWebhookPath)
	c.app.Flag("telegram.api-token", descTelegramAPIToken).Required().StringVar(&c.TeletramAPIToken)
	c.app.Flag("telegram.chat-id", descTelegramDefChatID).Required().Int64Var(&c.TelegramChatID)
	c.app.Flag("debug", descDebug).BoolVar(&c.DebugMode)
}

func (c *Config) validate() error {
	return nil
}
