package logger

import "github.com/sirupsen/logrus"

type logrusadapter struct {
	entry *logrus.Entry
}

func (l *logrusadapter) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logrusadapter) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusadapter) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusadapter) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusadapter) Panicf(format string, args ...interface{}) {
	l.entry.Panicf(format, args...)
}

func (l *logrusadapter) WithField(key string, value interface{}) Logger {
	return &logrusadapter{l.entry.WithField(key, value)}
}

func (l *logrusadapter) WithFields(fields map[string]interface{}) Logger {
	return &logrusadapter{l.entry.WithFields(fields)}
}
