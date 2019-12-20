package main

import (
	"errors"
	"os"
	"time"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	// Version will be populated in compilation time.
	Version = "dev"
)

// flag descriptions.
const (
	descAMListenAddr       = "The listen address where the server will be listening to alertmanager's webhook request."
	descAMWebhookPath      = "The path where the server will be handling the alertmanager webhook alert requests."
	descAMChatIDQS         = "The optional query string key used to customize the chat id of the notification. Does not depend on the notifier type."
	descAMDMSPath          = "The path for the dead man switch alerts from the Alertmanger."
	descTelegramAPIToken   = "The token that will be used to use the telegram API to send the alerts."
	descTelegramDefChatID  = "The default ID of the chat (group/channel) in telegram where the alerts will be sent."
	descMetricsListenAddr  = "The listen address where the metrics will be being served."
	descMetricsPath        = "The path where the metrics will be being served."
	descMetricsHCPath      = "The path where the healthcheck will be being served, it uses the same port as the metrics."
	descDMSEnable          = "Enables the dead man switch, that will send an alert if no alert is received at regular intervals."
	descDMSInterval        = "The interval the dead mans switch needs to receive an alert to not activate and send a notification alert (in Go time duration)."
	descDMSChatID          = "The chat ID (group/channel/room) the dead man's witch will sent the alerts. Does not depend on the notifier type and if not set it will be used notifier default chat ID."
	descDebug              = "Run the application in debug mode."
	descNotifyDryRun       = "Dry run the notification and show in the terminal instead of sending."
	descNotifyTemplatePath = "The path to set a custom template for the notification messages."
	descAlertLabelChatID   = "The label of the alert that will carry the chat id to forward the alert."
)

const (
	defAMListenAddr      = ":8080"
	defAMWebhookPath     = "/alerts"
	defAMChatIDQS        = "chat-id"
	defAMDMSPath         = "/alerts/dms"
	defMetricsListenAddr = ":8081"
	defMetricsPath       = "/metrics"
	defMetricsHCPath     = "/status"
	defDMSInterval       = "15m"
	defAlertLabelChatID  = "chat_id"
)

// Config has the configuration of the application.
type Config struct {
	AlertmanagerListenAddr         string
	AlertmanagerWebhookPath        string
	AlertmanagerChatIDQQueryString string
	AlertmanagerDMSPath            string
	TeletramAPIToken               string
	TelegramChatID                 int64
	MetricsListenAddr              string
	MetricsPath                    string
	MetricsHCPath                  string
	DMSInterval                    time.Duration
	DMSEnable                      bool
	DMSChatID                      string
	NotifyTemplate                 *os.File
	DebugMode                      bool
	NotifyDryRun                   bool
	AlertLabelChatID               string

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
	c.app.Flag("alertmanager.chat-id-query-string", descAMChatIDQS).Default(defAMChatIDQS).StringVar(&c.AlertmanagerChatIDQQueryString)
	c.app.Flag("alertmanager.dead-mans-switch-path", descAMDMSPath).Default(defAMDMSPath).StringVar(&c.AlertmanagerDMSPath)
	c.app.Flag("telegram.api-token", descTelegramAPIToken).StringVar(&c.TeletramAPIToken)
	c.app.Flag("telegram.chat-id", descTelegramDefChatID).Int64Var(&c.TelegramChatID)
	c.app.Flag("metrics.listen-address", descMetricsListenAddr).Default(defMetricsListenAddr).StringVar(&c.MetricsListenAddr)
	c.app.Flag("metrics.path", descMetricsPath).Default(defMetricsPath).StringVar(&c.MetricsPath)
	c.app.Flag("metrics.health-path", descMetricsHCPath).Default(defMetricsHCPath).StringVar(&c.MetricsHCPath)
	c.app.Flag("dead-mans-switch.enable", descDMSEnable).BoolVar(&c.DMSEnable)
	c.app.Flag("dead-mans-switch.interval", descDMSInterval).Default(defDMSInterval).DurationVar(&c.DMSInterval)
	c.app.Flag("dead-mans-switch.chat-id", descDMSChatID).StringVar(&c.DMSChatID)
	c.app.Flag("notify.dry-run", descNotifyDryRun).BoolVar(&c.NotifyDryRun)
	c.app.Flag("notify.template-path", descNotifyTemplatePath).FileVar(&c.NotifyTemplate)
	c.app.Flag("alert.label-chat-id", descAlertLabelChatID).Default(defAlertLabelChatID).StringVar(&c.AlertLabelChatID)
	c.app.Flag("debug", descDebug).BoolVar(&c.DebugMode)
}

func (c *Config) validate() error {
	if !c.NotifyDryRun {
		if c.TeletramAPIToken == "" {
			return errors.New("telegram api token is required")
		}

		if c.TelegramChatID == 0 {
			return errors.New("telegram default chat ID is required")
		}
	}
	return nil
}
