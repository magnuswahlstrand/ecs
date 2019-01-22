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
	// as an example, we prepend a shamrock to all log messages
	// but you can do whatever you want here.

	var systemName interface{} = "unknown"
	if val, ok := entry.Data["system"]; ok {
		systemName = val
		// delete(entry.Data, "system_name")
	}
	entry.Message = fmt.Sprintf("%-10s: ", systemName) + entry.Message
	return f.TextFormatter.Format(entry)
}

func NewSystemLogger(level ...logrus.Level) *logrus.Logger {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)
	if len(level) > 0 {
		logger.SetLevel(level[0])
	}
	logger.SetFormatter(&SystemFormatter{})
	return logger
}
