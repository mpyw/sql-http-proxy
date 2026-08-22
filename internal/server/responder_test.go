package server

import (
	"errors"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResponder_Respond_V1Compatible pins the response bytes that
// encoding/json v1 produced. Responses are this project's public API, so each
// case here is a v2 default that marshalOptions deliberately opts out of; a
// failure means a v2 default leaked into the wire format.
func TestResponder_Respond_V1Compatible(t *testing.T) {
	ts := time.Date(2026, 8, 22, 12, 34, 56, 789000000, time.UTC)

	tests := []struct {
		name    string
		payload any
		want    string
	}{
		{
			// v2 preserves insertion order; v1 always sorted.
			name:    "map keys are sorted",
			payload: map[string]any{"z": 1, "m": 2, "a": 3},
			want:    `{"a":3,"m":2,"z":1}`,
		},
		{
			// v2 renders these as [] and {}.
			name:    "nil slice and nil map are null",
			payload: map[string]any{"rows": []map[string]any(nil), "obj": map[string]any(nil)},
			want:    `{"obj":null,"rows":null}`,
		},
		{
			name:    "empty non-nil slice stays an empty array",
			payload: []map[string]any{},
			want:    `[]`,
		},
		{
			// v2 rejects invalid UTF-8 outright, which would 500 a row that
			// serves fine today.
			name:    "invalid UTF-8 becomes U+FFFD",
			payload: map[string]any{"name": "a\xff\xfeb"},
			want:    "{\"name\":\"a��b\"}",
		},
		{
			name:    "byte slices are base64",
			payload: map[string]any{"blob": []byte{0x00, 0xff}},
			want:    `{"blob":"AP8="}`,
		},
		{
			name:    "time is RFC 3339",
			payload: map[string]any{"ts": ts},
			want:    `{"ts":"2026-08-22T12:34:56.789Z"}`,
		},
		{
			name:    "SQL NULL stays null",
			payload: map[string]any{"n": nil},
			want:    `{"n":null}`,
		},
		{
			// The one intentional difference from v1, which escaped these.
			// Responses are application/json, so escaping bought nothing.
			name:    "HTML metacharacters are not escaped",
			payload: map[string]any{"html": "<b>&</b>"},
			want:    `{"html":"<b>&</b>"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			(&responder{w: w}).Respond(http.StatusOK, tt.payload)

			assert.Equal(t, http.StatusOK, w.Code)
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
			assert.Equal(t, tt.want, w.Body.String())
		})
	}
}

func TestResponder_Respond_NoContent(t *testing.T) {
	w := httptest.NewRecorder()
	(&responder{w: w}).Respond(http.StatusNoContent, map[string]any{"ignored": true})

	assert.Equal(t, http.StatusNoContent, w.Code)
	assert.Empty(t, w.Body.String())
}

func TestResponder_Respond_UnmarshalableValue(t *testing.T) {
	// v2 rejects NaN and Inf just as v1 did. Marshaling happens before the
	// header is written, so this surfaces as a 500 rather than a truncated body.
	w := httptest.NewRecorder()
	(&responder{w: w}).Respond(http.StatusOK, map[string]any{"n": math.Inf(1)})

	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "error")
}

func TestResponder_Error(t *testing.T) {
	t.Run("wraps the message as JSON", func(t *testing.T) {
		w := httptest.NewRecorder()
		(&responder{w: w}).Error(http.StatusBadRequest, errors.New(`bad "input"`))

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
		assert.JSONEq(t, `{"error":"bad \"input\""}`, w.Body.String())
	})

	t.Run("survives a non-UTF-8 message", func(t *testing.T) {
		// responder.Error marshals with lo.Must, so a driver error carrying raw
		// bytes must not be able to panic the handler.
		w := httptest.NewRecorder()
		require.NotPanics(t, func() {
			(&responder{w: w}).Error(http.StatusInternalServerError, errors.New("driver\xff\xfe"))
		})

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Equal(t, "{\"error\":\"driver��\"}", w.Body.String())
	})
}

func TestValidateStatus(t *testing.T) {
	tests := []struct {
		in, want int
	}{
		{200, 200},
		{100, 100},
		{599, 599},
		{99, http.StatusInternalServerError},
		{600, http.StatusInternalServerError},
		{0, http.StatusInternalServerError},
		{-1, http.StatusInternalServerError},
	}

	for _, tt := range tests {
		assert.Equal(t, tt.want, validateStatus(tt.in), "validateStatus(%d)", tt.in)
	}
}
