package logging

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Logger interface {
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	WithField(key string, value interface{}) *logrus.Entry
}

type SystemFormatter struct {
	logrus.TextFormatter
}

func (f *SystemFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	entry.Message = fmt.Sprintf("%-80s", fmt.Sprintf("%10s: ", entry.Data["s"])+entry.Message)
	return f.TextFormatter.Format(entry)
}

func NewLogger(level ...logrus.Level) Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	if len(level) > 0 {
		logger.SetLevel(level[0])
	}
	logger.SetFormatter(&SystemFormatter{})
	return logger
}
