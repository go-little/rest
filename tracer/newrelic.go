package tracer

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	newrelic "github.com/newrelic/go-agent"
)

func newrelicTransaction(w http.ResponseWriter, r *http.Request) newrelic.Transaction {
	if config.newrelicApp != nil {
		route := mux.CurrentRoute(r)

		method, _ := route.GetMethods()
		pathPattern, _ := route.GetPathTemplate()

		txnName := fmt.Sprintf("%s %s", method, pathPattern)
		txn := config.newrelicApp.StartTransaction(txnName, w, r)

		return txn
	}
	return nil
}

func newrelicStartSegment(txn newrelic.Transaction, name string) *newrelic.Segment {
	var segment *newrelic.Segment
	if txn != nil {
		segment = newrelic.StartSegment(txn, name)
	}
	return segment
}
