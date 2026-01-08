// Package body provides HTTP request body parsing with Content-Type handling.
package body

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/http"
	"slices"

	"github.com/mpyw/sql-http-proxy/internal/charset"
)

// AcceptType constants
const (
	AcceptJSON = "json"
	AcceptForm = "form"
)

// MaxBodySize is the maximum allowed request body size (10MB).
const MaxBodySize = 10 * 1024 * 1024

// Error types
var (
	ErrUnsupportedMediaType = errors.New("unsupported media type")
	ErrBadRequest           = errors.New("bad request")
	ErrBodyTooLarge         = errors.New("request body too large")
)

// Parser parses HTTP request bodies based on Content-Type.
type Parser struct {
	accepts []string
}

// NewParser creates a new Parser with the given accepted types.
func NewParser(accepts []string) *Parser {
	return &Parser{accepts: accepts}
}

// Parse parses the request body based on Content-Type.
// Returns ErrUnsupportedMediaType if Content-Type is not in accepts.
// Returns ErrBadRequest if body parsing fails.
// Returns ErrBodyTooLarge if body exceeds MaxBodySize.
func (p *Parser) Parse(r *http.Request) (map[string]any, error) {
	// Limit body size to prevent DoS
	body := io.LimitReader(r.Body, MaxBodySize+1)

	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		// Default to json if no Content-Type specified and json is accepted
		if slices.Contains(p.accepts, AcceptJSON) {
			return p.parseJSON(body, "")
		}
		return nil, fmt.Errorf("%w: missing Content-Type header", ErrUnsupportedMediaType)
	}

	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid Content-Type: %v", ErrBadRequest, err)
	}

	charsetName := params["charset"]

	switch mediaType {
	case "application/json":
		if !slices.Contains(p.accepts, AcceptJSON) {
			return nil, fmt.Errorf("%w: application/json not accepted", ErrUnsupportedMediaType)
		}
		return p.parseJSON(body, charsetName)

	case "application/x-www-form-urlencoded":
		if !slices.Contains(p.accepts, AcceptForm) {
			return nil, fmt.Errorf("%w: application/x-www-form-urlencoded not accepted", ErrUnsupportedMediaType)
		}
		return p.parseURLEncoded(body, charsetName)

	case "multipart/form-data":
		if !slices.Contains(p.accepts, AcceptForm) {
			return nil, fmt.Errorf("%w: multipart/form-data not accepted", ErrUnsupportedMediaType)
		}
		boundary := params["boundary"]
		if boundary == "" {
			return nil, fmt.Errorf("%w: missing boundary in multipart/form-data", ErrBadRequest)
		}
		return p.parseMultipart(body, boundary, charsetName)

	default:
		return nil, fmt.Errorf("%w: %s", ErrUnsupportedMediaType, mediaType)
	}
}

func (p *Parser) parseJSON(body io.Reader, charsetName string) (map[string]any, error) {
	data, err := io.ReadAll(body)
	if err != nil {
		return nil, fmt.Errorf("%w: failed to read body: %v", ErrBadRequest, err)
	}
	if len(data) > MaxBodySize {
		return nil, ErrBodyTooLarge
	}

	if charsetName != "" {
		data, err = charset.ToUTF8(data, charsetName)
		if err != nil {
			return nil, fmt.Errorf("%w: charset conversion failed: %v", ErrBadRequest, err)
		}
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("%w: invalid JSON: %v", ErrBadRequest, err)
	}
	return result, nil
}
