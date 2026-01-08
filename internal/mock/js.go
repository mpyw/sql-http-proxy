package mock

import (
	"fmt"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// jsSource wraps js.Transformer for mock data generation.
type jsSource struct {
	transformer *js.Transformer
}

// SetHelpers sets the compiled helpers for this mock source.
func (s *jsSource) SetHelpers(h *js.CompiledHelpers) {
	s.transformer.SetHelpers(h)
}

// compileJS compiles JavaScript code into a jsSource.
func compileJS(code string) (*jsSource, error) {
	t, err := js.CompileMockJS(code)
	if err != nil {
		return nil, fmt.Errorf("failed to compile mock: %w", err)
	}
	return &jsSource{transformer: t}, nil
}

// Data executes the JavaScript function and returns the mock data.
func (s *jsSource) Data(ctx map[string]any, sql string, input map[string]any, tc *js.TransformContext) (any, map[string]any, error) {
	result, err := s.transformer.ApplyMock(ctx, sql, input, tc)
	if err != nil {
		return nil, nil, err
	}
	return result.Output, result.Ctx, nil
}
