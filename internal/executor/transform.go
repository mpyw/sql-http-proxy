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
	PostEach *js.Transformer
	PostAll  *js.Transformer
}

// CompileTransformOptions contains options for compiling transforms.
type CompileTransformOptions struct {
	ConfigDir   string
	Helpers     *js.CompiledHelpers
	ValueParser *mock.ValueParser
}

// CompileTransforms compiles all transforms from a config.Transform.
// Note: mock is now compiled separately via CompileMock.
func CompileTransforms(t *config.Transform, opts CompileTransformOptions) (*Transforms, error) {
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
			pre.SetHelpers(opts.Helpers)
			transforms.Pre = pre
		}
	}

	// Compile post.each transform
	if t.Post.Each != "" {
		postEach, err := js.CompilePost(t.Post.Each)
		if err != nil {
			errs = append(errs, fmt.Errorf("post.each: %w", err))
		} else {
			postEach.SetHelpers(opts.Helpers)
			transforms.PostEach = postEach
		}
	}

	// Compile post.all transform
	if t.Post.All != "" {
		postAll, err := js.CompilePost(t.Post.All)
		if err != nil {
			errs = append(errs, fmt.Errorf("post.all: %w", err))
		} else {
			postAll.SetHelpers(opts.Helpers)
			transforms.PostAll = postAll
		}
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return transforms, nil
}

// CompileMock compiles a mock source from config.Mock.
func CompileMock(m *config.Mock, opts CompileTransformOptions) (mock.Source, error) {
	if m == nil || m.IsEmpty() {
		return nil, nil
	}

	mockOpts := mock.CompileOptions{
		ConfigDir:   opts.ConfigDir,
		Helpers:     opts.Helpers,
		ValueParser: opts.ValueParser,
	}

	return mock.Compile(m, mockOpts)
}
