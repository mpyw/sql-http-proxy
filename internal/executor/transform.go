package executor

import (
	"errors"
	"fmt"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/js"
	"github.com/mpyw/sql-http-proxy/internal/mock"
)

// Transforms holds pre-compiled transforms.
type Transforms struct {
	Pre      *js.Transformer
	Mock     mock.Source
	PostEach *js.Transformer
	PostAll  *js.Transformer
}

// IsMock returns true if mock is configured.
func (t *Transforms) IsMock() bool {
	return t.Mock != nil
}

// CompileTransforms compiles all transforms from a config.Transform.
// configDir is the directory of the YAML config file (for resolving relative paths in mock).
func CompileTransforms(t *config.Transform, configDir string) (*Transforms, error) {
	if t == nil {
		return &Transforms{}, nil
	}

	transforms := &Transforms{}
	var errs []error

	// Compile pre-transform
	if t.Pre != "" {
		pre, err := js.CompilePre(t.Pre)
		if err != nil {
			errs = append(errs, fmt.Errorf("pre-transform: %w", err))
		} else {
			transforms.Pre = pre
		}
	}

	// Compile mock source
	if t.Mock != nil && !t.Mock.IsEmpty() {
		mockSource, err := mock.Compile(t.Mock, configDir)
		if err != nil {
			errs = append(errs, fmt.Errorf("mock: %w", err))
		} else {
			transforms.Mock = mockSource
		}
	}

	// Compile post.each transform
	if t.Post.Each != "" {
		postEach, err := compilePost(t.Post.Each)
		if err != nil {
			errs = append(errs, fmt.Errorf("post.each: %w", err))
		} else {
			transforms.PostEach = postEach
		}
	}

	// Compile post.all transform
	if t.Post.All != "" {
		postAll, err := compilePost(t.Post.All)
		if err != nil {
			errs = append(errs, fmt.Errorf("post.all: %w", err))
		} else {
			transforms.PostAll = postAll
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return transforms, nil
}

func compilePost(code string) (*js.Transformer, error) {
	return js.CompilePost(code)
}
