package middleware

import (
	"net/http"

	"io/ioutil"

	"bytes"

	"github.com/gorilla/mux"

	"github.com/go-little/rest/response"

	"github.com/go-little/rest/tracer"

	newrelic "github.com/newrelic/go-agent"
)

type TracerMiddlewareConfig struct {
	LoggerConfig   tracer.LoggerConfig
	NewrelicConfig newrelic.Config
}

// TracerMiddleware comment
func TracerMiddleware(config TracerMiddlewareConfig) func(next http.Handler) http.Handler {

	tracer.Config(config.LoggerConfig, config.NewrelicConfig)

	return func(next http.Handler) http.Handler {

		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tracer.Start(w, r)
			defer tracer.End(r.Context())

			ctx := r.Context()

			segment := tracer.StartSegment(ctx, "http")
			defer segment.End()

			putSegmentRequestAttr(segment, r)

			rw := response.NewResponseWriterWrapper(w)
			next.ServeHTTP(rw, r)

			putSegmentResponseAttr(segment, rw)

		})

	}
}

func putSegmentRequestAttr(segment *tracer.Segment, r *http.Request) {
	segment.Attr("request_header", r.Header)

	segment.Attr("request_method", r.Method)

	route := mux.CurrentRoute(r)
	pathPattern, _ := route.GetPathTemplate()
	segment.Attr("request_path_pattern", pathPattern)

	segment.Attr("request_path", r.URL.Path)

	segment.Attr("request_querystring", r.URL.RawQuery)

	if r.Body != http.NoBody {
		buf, _ := ioutil.ReadAll(r.Body)
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf))
		segment.Attr("request_body", string(buf))
	}
}

func putSegmentResponseAttr(segment *tracer.Segment, rw *response.ResponseWriterWrapper) {
	segment.Attr("response_header", rw.Header())
	segment.Attr("response_status_code", rw.StatusCode)
	segment.Attr("response_body", string(rw.Body))
}
