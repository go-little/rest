package tracer

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	newrelic "github.com/newrelic/go-agent"
	"github.com/sirupsen/logrus"
)

const TRACER_CTX_KEY = "TRACER_CTX_KEY"

type tracer struct {
	newrelicTransaction newrelic.Transaction
	attributes          map[string]interface{}
}

type Segment struct {
	name    string
	startAt time.Time
	endAt   time.Time
	tracer  *tracer
	segment *newrelic.Segment
}

type ExternalSegment struct {
	name            string
	startAt         time.Time
	endAt           time.Time
	tracer          *tracer
	externalSegment *newrelic.ExternalSegment
}

var logger *logrus.Logger
var newrelicWrapper *NewrelicWrapper

func Config(loggerConfig LoggerConfig, newrelicConfig newrelic.Config) {
	logger = NewLogger(loggerConfig)
	newrelicWrapper = NewNewrelicWrapper(newrelicConfig)
}

// Start comment
func Start(w http.ResponseWriter, r *http.Request) {

	t := &tracer{
		newrelicTransaction: newrelicWrapper.newrelicTransaction(w, r),
		attributes:          make(map[string]interface{}),
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, TRACER_CTX_KEY, t)
	*r = *r.WithContext(ctx)
}

func normalizeSegmentName(name string) string {
	name = strings.Replace(name, " ", "_", -1)
	name = strings.ToLower(name)
	return name
}

// StartSegment comment
func StartSegment(ctx context.Context, name string) *Segment {
	s := &Segment{
		name:    normalizeSegmentName(name),
		startAt: time.Now(),
	}

	tracerI := ctx.Value(TRACER_CTX_KEY)
	if tracer, ok := tracerI.(*tracer); ok {
		s.tracer = tracer
		s.segment = newrelicWrapper.newrelicStartSegment(tracer.newrelicTransaction, name)
	}

	return s
}

func (s *Segment) Attr(key string, value interface{}) *Segment {
	if s.tracer != nil {
		s.tracer.attributes[fmt.Sprintf("%s.%s", s.name, key)] = value
	}
	return s
}

func (s *Segment) End() {
	s.endAt = time.Now()
	s.Attr("elapsed_milliseconds", s.endAt.Sub(s.startAt)/time.Millisecond)
	if s.segment != nil {
		s.segment.End()
	}
}

func StartExternalSegment(ctx context.Context, name string, req *http.Request) *ExternalSegment {
	s := &ExternalSegment{
		name:    normalizeSegmentName(name),
		startAt: time.Now(),
	}

	tracerI := ctx.Value(TRACER_CTX_KEY)
	if tracer, ok := tracerI.(*tracer); ok {
		s.tracer = tracer
		s.externalSegment = newrelicWrapper.newrelicStartExternalSegment(tracer.newrelicTransaction, req)
	}

	return s
}

func (s *ExternalSegment) Attr(key string, value interface{}) *ExternalSegment {
	if s.tracer != nil {
		s.tracer.attributes[fmt.Sprintf("%s.%s", s.name, key)] = value
	}
	return s
}

func (s *ExternalSegment) End(response *http.Response) {
	s.endAt = time.Now()
	s.Attr("elapsed_milliseconds", s.endAt.Sub(s.startAt)/time.Millisecond)
	if s.externalSegment != nil {
		s.externalSegment.Response = response
		s.externalSegment.End()
	}
}

// End comment
func End(ctx context.Context) {
	tracerI := ctx.Value(TRACER_CTX_KEY)
	if tracer, ok := tracerI.(*tracer); ok {
		if tracer.newrelicTransaction != nil {
			for key, value := range tracer.attributes {
				tracer.newrelicTransaction.AddAttribute(key, value)
			}
			tracer.newrelicTransaction.End()
		}

		logger.WithFields(tracer.attributes).Info()
	}
}
