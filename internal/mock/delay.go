package mock

import (
	"time"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// delayedSource wraps a Source and adds artificial latency.
type delayedSource struct {
	source Source
	delay  time.Duration
}

// newDelayedSource creates a delayedSource that adds delay before returning data.
func newDelayedSource(source Source, delay time.Duration) *delayedSource {
	return &delayedSource{
		source: source,
		delay:  delay,
	}
}

// Data returns the data from the wrapped source after the configured delay.
func (d *delayedSource) Data(ctx map[string]any, sql string, input map[string]any, tc *js.TransformContext) (any, map[string]any, error) {
	time.Sleep(d.delay)
	return d.source.Data(ctx, sql, input, tc)
}
