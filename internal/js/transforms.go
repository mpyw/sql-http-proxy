package js

import (
	"errors"
	"fmt"
)

// Transforms holds pre-compiled JavaScript transformers.
type Transforms struct {
	Pre      *Transformer
	Mock     *Transformer
	PostEach *Transformer
	PostAll  *Transformer
}

// IsMock returns true if mock transform is configured.
func (t *Transforms) IsMock() bool {
	return t.Mock != nil
}

// CompileTransforms compiles all transforms from source strings.
// Returns all compilation errors joined, not just the first one.
func CompileTransforms(pre, mock, postEach, postAll string) (*Transforms, error) {
	transforms := &Transforms{}
	var errs []error

	compilers := []struct {
		field   string
		source  string
		target  **Transformer
		compile func(string) (*Transformer, error)
	}{
		{"pre-transform", pre, &transforms.Pre, CompilePre},
		{"mock", mock, &transforms.Mock, CompilePre},
		{"post.each", postEach, &transforms.PostEach, CompilePost},
		{"post.all", postAll, &transforms.PostAll, CompilePost},
	}

	for _, c := range compilers {
		if c.source == "" {
			continue
		}
		compiled, err := c.compile(c.source)
		if err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", c.field, err))
			continue
		}
		*c.target = compiled
	}

	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	return transforms, nil
}
