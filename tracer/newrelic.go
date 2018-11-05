package tracer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	newrelic "github.com/newrelic/go-agent"
)

type NewrelicConfig struct {
	AppName string
	License string
}

type NewrelicWrapper struct {
	newrelic.Application
}

func NewNewrelicWrapper(config newrelic.Config) *NewrelicWrapper {
	var app newrelic.Application

	if config.AppName != "" && config.License != "" {
		var err error
		newRelicConfig := newrelic.NewConfig(config.AppName, config.License)
		app, err = newrelic.NewApplication(newRelicConfig)

		if err != nil {
			fmt.Printf("Error on start NewRelic App: %v", err)
		} else {
			fmt.Printf("New Relic app: %s\n", config.AppName)
		}

	}

	return &NewrelicWrapper{
		Application: app,
	}
}

func (n *NewrelicWrapper) newrelicTransaction(w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	if n.Application != nil {
		route := mux.CurrentRoute(r)

		method, _ := route.GetMethods()
		pathPattern, _ := route.GetPathTemplate()

		txnName := fmt.Sprintf("%s (%s)", pathPattern, method)
		txn := n.Application.StartTransaction(txnName, w, r)

		return txn
	}
	return nil
}

func (n *NewrelicWrapper) newrelicStartSegment(txn newrelic.Transaction, name string) *newrelic.Segment {
	var segment *newrelic.Segment
	if txn != nil {
		segment = newrelic.StartSegment(txn, name)
	}
	return segment
}

func (n *NewrelicWrapper) newrelicStartExternalSegment(txn newrelic.Transaction, req *http.Request) *newrelic.ExternalSegment {
	var segment *newrelic.ExternalSegment
	if txn != nil {
		segment = newrelic.StartExternalSegment(txn, req)
	}
	return segment
}
