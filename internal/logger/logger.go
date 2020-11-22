package logger

import (
	"fmt"
	"io/ioutil"

	"github.com/edusalguero/roteiro.git/internal/config"
	"github.com/sirupsen/logrus"
)

//go:generate mockgen -source=./logger.go -destination=./mock/logger.go
type Logger interface {
	Info(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Fatalf(string, ...interface{})
	Panicf(string, ...interface{})

	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger
}

func New(cnf config.Logger) (Logger, error) {
	lvl, err := logrus.ParseLevel(cnf.Level)
	if err != nil {
		return nil, fmt.Errorf("could not parse level")
	}

	internalLogger := logrus.New()
	internalLogger.SetLevel(lvl)
	internalLogger.SetFormatter(new(logrus.JSONFormatter))

	return &logrusadapter{
		entry: logrus.NewEntry(internalLogger),
	}, nil
}

// NewNopLogger returns a logger that discards all log messages.
func NewNopLogger() Logger {
	l := logrus.New()
	l.SetOutput(ioutil.Discard)
	return &logrusadapter{
		entry: logrus.NewEntry(l),
	}
}
