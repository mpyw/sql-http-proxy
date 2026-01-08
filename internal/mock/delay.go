package mock

import (
	"time"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// DelayedSource wraps a Source and adds artificial latency.
type DelayedSource struct {
	source Source
	delay  time.Duration
}

// NewDelayedSource creates a DelayedSource that adds delay before returning data.
func NewDelayedSource(source Source, delay time.Duration) *DelayedSource {
	return &DelayedSource{
		source: source,
		delay:  delay,
	}
}

// Data returns the data from the wrapped source after the configured delay.
func (d *DelayedSource) Data(ctx map[string]any, sql string, input map[string]any, tc *js.TransformContext) (any, map[string]any, error) {
	time.Sleep(d.delay)
	return d.source.Data(ctx, sql, input, tc)
}
