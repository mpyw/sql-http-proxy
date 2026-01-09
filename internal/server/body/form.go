package body

import (
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/url"

	"github.com/samber/lo"

	"github.com/mpyw/sql-http-proxy/internal/charset"
)

// parseURLEncoded parses URL-encoded form body with optional charset conversion.
func parseURLEncoded(body io.Reader, charsetName string) (map[string]any, error) {
	data, err := readBody(body, charsetName)
	if err != nil {
		return nil, err
	}

	values, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, fmt.Errorf("%w: invalid form data: %v", ErrBadRequest, err)
	}

	return formValuesToMap(values), nil
}

// parseMultipart parses multipart form body with optional charset conversion.
func parseMultipart(body io.Reader, boundary, charsetName string) (map[string]any, error) {
	reader := multipart.NewReader(body, boundary)
	result := make(map[string]any)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("%w: failed to read multipart: %v", ErrBadRequest, err)
		}

		// Skip file uploads (only process form fields)
		if part.FileName() != "" {
			if err := part.Close(); err != nil {
				slog.Warn("Failed to close multipart part", "error", err)
			}
			continue
		}

		name := part.FormName()
		if name == "" {
			if err := part.Close(); err != nil {
				slog.Warn("Failed to close multipart part", "error", err)
			}
			continue
		}

		data, err := io.ReadAll(part)
		if closeErr := part.Close(); closeErr != nil {
			slog.Warn("Failed to close multipart part", "error", closeErr)
		}
		if err != nil {
			return nil, fmt.Errorf("%w: failed to read multipart field: %v", ErrBadRequest, err)
		}

		// Apply charset conversion if specified
		if charsetName != "" {
			data, err = charset.ToUTF8(data, charsetName)
			if err != nil {
				return nil, fmt.Errorf("%w: charset conversion failed: %v", ErrBadRequest, err)
			}
		}

		// Handle multiple values with same name
		if existing, ok := result[name]; ok {
			switch v := existing.(type) {
			case []any:
				result[name] = append(v, string(data))
			default:
				result[name] = []any{v, string(data)}
			}
		} else {
			result[name] = string(data)
		}
	}

	return result, nil
}

// formValuesToMap converts url.Values to map[string]any.
// Single values are stored as strings, multiple values as []any.
func formValuesToMap(values url.Values) map[string]any {
	result := make(map[string]any)
	for key, vals := range values {
		if len(vals) == 1 {
			result[key] = vals[0]
		} else {
			result[key] = lo.ToAnySlice(vals)
		}
	}
	return result
}
