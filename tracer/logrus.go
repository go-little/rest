package tracer

import (
	"time"

	"github.com/sirupsen/logrus"
)

// RootLogger comment
func RootLogger() *logrus.Logger {

	rootLogger := logrus.New()

	rootLogger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	return rootLogger
}
