package tracer

import (
	"context"
	"fmt"
	"net/http"
	"time"

	newrelic "github.com/newrelic/go-agent"
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

// Start comment
func Start(w http.ResponseWriter, r *http.Request) {

	t := &tracer{
		newrelicTransaction: newrelicTransaction(w, r),
		attributes:          make(map[string]interface{}),
	}

	ctx := r.Context()
	ctx = context.WithValue(ctx, TRACER_CTX_KEY, t)
	*r = *r.WithContext(ctx)
}

// StartSegment comment
func StartSegment(ctx context.Context, name string) *Segment {
	s := &Segment{
		name:    name,
		startAt: time.Now(),
	}

	tracerI := ctx.Value(TRACER_CTX_KEY)
	if tracer, ok := tracerI.(*tracer); ok {
		s.tracer = tracer
		s.segment = newrelicStartSegment(tracer.newrelicTransaction, name)
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

		config.logger.WithFields(tracer.attributes).Info()
	}
}
