package logrus

import (
	"github.com/sirupsen/logrus"

	"github.com/slok/alertgram/internal/log"
)

type logger struct {
	*logrus.Entry
}

// New returns a new implementation of a logrus logger.
// If not debug mode it will use JSON logging.
func New(debug bool) log.Logger {
	l := logrus.New()
	if debug {
		l.SetLevel(logrus.DebugLevel)
	} else {
		l.SetFormatter(&logrus.JSONFormatter{})
	}

	return &logger{Entry: logrus.NewEntry(l)}
}

func (l logger) WithValues(vals map[string]interface{}) log.Logger {
	return logger{l.Entry.WithFields(vals)}
}
