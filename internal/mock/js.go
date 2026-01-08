package mock

import (
	"fmt"

	"github.com/mpyw/sql-http-proxy/internal/js"
)

// JSSource wraps js.Transformer for mock data generation.
type JSSource struct {
	transformer *js.Transformer
}

// SetHelpers sets the compiled helpers for this mock source.
func (s *JSSource) SetHelpers(h *js.CompiledHelpers) {
	s.transformer.SetHelpers(h)
}

// CompileJS compiles JavaScript code into a JSSource.
func CompileJS(code string) (*JSSource, error) {
	t, err := js.CompileMockJS(code)
	if err != nil {
		return nil, fmt.Errorf("failed to compile mock: %w", err)
	}
	return &JSSource{transformer: t}, nil
}

// Data executes the JavaScript function and returns the mock data.
func (s *JSSource) Data(ctx map[string]any, sql string, input map[string]any, tc *js.TransformContext) (any, map[string]any, error) {
	result, err := s.transformer.ApplyMock(ctx, sql, input, tc)
	if err != nil {
		return nil, nil, err
	}
	return result.Output, result.Ctx, nil
}
