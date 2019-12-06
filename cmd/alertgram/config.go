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
	descAMListenAddr  = "The listen address where the server will be listening to alertmanager's webhook request."
	descAMWebhookPath = "The path where the server will be handling the alertmanager webhook alert requests."
	descDebug         = "Run the application in debug mode."
)

const (
	defAMListenAddr  = ":8080"
	defAMWebhookPath = "/alerts"
)

// Config has the configuration of the application.
type Config struct {
	AlertmanagerListenAddr  string
	AlertmanagerWebhookPath string
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

	c.app.Parse(os.Args[1:])
	if err := c.validate(); err != nil {
		return nil, err
	}

	return c, nil
}

func (c *Config) registerFlags() {
	c.app.Flag("alertmanager.listen-address", descAMListenAddr).Default(defAMListenAddr).StringVar(&c.AlertmanagerListenAddr)
	c.app.Flag("alertmanager.webhook-path", descAMWebhookPath).Default(defAMWebhookPath).StringVar(&c.AlertmanagerWebhookPath)
	c.app.Flag("debug", descDebug).BoolVar(&c.DebugMode)
}

func (c *Config) validate() error {
	return nil
}
