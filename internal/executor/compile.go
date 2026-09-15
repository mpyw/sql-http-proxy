package executor

import (
	"errors"
	"fmt"

	"github.com/samber/lo"

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
	var err error

	for _, entry := range []lo.Tuple4[string, func(string) (*js.Transformer, error), string, **js.Transformer]{
		lo.T4("pre-transform", js.CompilePre, t.Pre, &transforms.Pre),
		lo.T4("post.each", js.CompilePost, t.Post.Each, &transforms.PostEach),
		lo.T4("post.all", js.CompilePost, t.Post.All, &transforms.PostAll),
	} {
		name, compile, script, target := entry.Unpack()
		if script == "" {
			continue
		}
		compiled, compileErr := compile(script)
		if compileErr != nil {
			err = errors.Join(err, fmt.Errorf("%s: %w", name, compileErr))
		} else {
			compiled.SetHelpers(opts.Helpers)
			*target = compiled
		}
	}

	if err != nil {
		return nil, err
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
