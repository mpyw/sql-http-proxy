// Package executor provides query and mutation execution logic.
package executor

import "errors"

// ErrNotFound is returned when a "one" type query finds no rows.
var ErrNotFound = errors.New("not found")

// PhaseError wraps an error with the phase (pre/mock/post) it occurred in.
type PhaseError struct {
	Phase string // "pre", "mock", or "post"
	Err   error
}

func (e *PhaseError) Error() string {
	return e.Err.Error()
}

func (e *PhaseError) Unwrap() error {
	return e.Err
}

// WrapPreError wraps an error as a pre-transform error.
func WrapPreError(err error) error {
	if err == nil {
		return nil
	}
	return &PhaseError{Phase: "pre", Err: err}
}

// WrapMockError wraps an error as a mock-transform error.
func WrapMockError(err error) error {
	if err == nil {
		return nil
	}
	return &PhaseError{Phase: "mock", Err: err}
}

// WrapPostError wraps an error as a post-transform error.
func WrapPostError(err error) error {
	if err == nil {
		return nil
	}
	return &PhaseError{Phase: "post", Err: err}
}

// Record represents an executed query/mutation record for testing/logging.
type Record struct {
	Type   string         // "one", "many", or "none"
	IsMock bool           // true if mock execution
	SQL    string         // executed SQL (bound for DB, original for mock)
	Params map[string]any // named parameters (for mock)
	Args   []any          // positional arguments (for DB)
}

// Recorder is called when a query/mutation is executed.
type Recorder func(record Record)

// Options contains optional settings for executors.
type Options struct {
	Recorder Recorder
}
