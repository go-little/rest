package tracer

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/luizalabs/burzumlogs-sdk/go-burzumlogs-sdk/logzum"
	newrelic "github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
)

type tracerConfig struct {
	logger      *logrus.Logger
	newrelicApp newrelic.Application
}

var config tracerConfig

func init() {

	config = tracerConfig{}

	logger := RootLogger()

	logStdout := strings.ToLower(os.Getenv("LOG_STDOUT")) == "true"
	if !logStdout {
		logger.Out = ioutil.Discard
	}
	fmt.Printf("log stdout: %v\n", logStdout)

	burzumToken := os.Getenv("BURZUM_TOKEN")

	if burzumToken != "" {
		burzumHook, err := logzum.New(burzumToken)
		if err != nil {
			fmt.Printf("Error on start Burzum Hook: %v", err)
		} else {
			logger.AddHook(burzumHook)
			fmt.Printf("burzum logs with token: %s\n", burzumToken)
		}
	}

	config.logger = logger

	newrelicAppName := os.Getenv("NEWRELIC_APP_NAME")
	newrelicLicense := os.Getenv("NEWRELIC_LICENSE")
	if newrelicAppName != "" && newrelicLicense != "" {
		newRelicConfig := newrelic.NewConfig(newrelicAppName, newrelicLicense)
		app, err := newrelic.NewApplication(newRelicConfig)

		if err != nil {
			fmt.Printf("Error on start NewRelic App: %v", err)
		} else {
			config.newrelicApp = app
			fmt.Printf("New Relic app: %s\n", newrelicAppName)
		}

	}
}
