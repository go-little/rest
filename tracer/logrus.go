package tracer

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/luizalabs/burzumlogs-sdk/go-burzumlogs-sdk/logzum"
)

type LoggerConfig struct {
	Stdout       bool
	BurzumToken  string
	BurzumConfig logzum.Config
}

// NewLogger comment
func NewLogger(config LoggerConfig) *logrus.Logger {

	logger := logrus.New()

	logger.Formatter = &logrus.JSONFormatter{
		TimestampFormat: time.RFC3339,
	}

	if !config.Stdout {
		logger.Out = ioutil.Discard
	}

	if config.BurzumToken != "" {
		burzumHook, err := logzum.New(config.BurzumToken)
		if err != nil {
			fmt.Printf("Error on start Burzum Hook: %v", err)
		} else {
			logger.AddHook(burzumHook)
			fmt.Printf("burzum logs with token: %s\n", config.BurzumToken)
		}
	}

	return logger
}
