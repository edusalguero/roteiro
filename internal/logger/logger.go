package logger

import (
	"fmt"
	"io/ioutil"

	"github.com/sirupsen/logrus"
)

var (
	DefaultLogger Logger = NewLogger()
)

type Config struct {
	Level string `default:"debug"`
}

//go:generate mockgen -source=./logger.go -destination=./mock/logger.go
type Logger interface {
	Info(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Panicf(string, ...interface{})
	Criticalf(string, ...interface{})

	WithField(string, interface{}) Logger
	WithFields(map[string]interface{}) Logger
}

func NewLogger() Logger {
	logger, _ := New(Config{Level: "info"})

	return logger
}
func New(cnf Config) (Logger, error) {
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
