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

// CompileTransformOptions contains options for compiling transforms.
type CompileTransformOptions struct {
	ConfigDir   string
	Helpers     *js.CompiledHelpers
	ValueParser *mock.ValueParser
}

// CompileTransforms compiles all transforms from a config.Transform.
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

	// Compile mock source
	if t.Mock != nil && !t.Mock.IsEmpty() {
		mockOpts := mock.CompileOptions{
			ConfigDir:   opts.ConfigDir,
			Helpers:     opts.Helpers,
			ValueParser: opts.ValueParser,
		}
		mockSource, err := mock.Compile(t.Mock, mockOpts)
		if err != nil {
			errs = append(errs, fmt.Errorf("mock: %w", err))
		} else {
			transforms.Mock = mockSource
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
