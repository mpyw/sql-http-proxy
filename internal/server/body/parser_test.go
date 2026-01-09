package body

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
)

func TestNewParser(t *testing.T) {
	p := NewParser([]config.AcceptType{config.AcceptJSON, config.AcceptForm})
	assert.NotNil(t, p)
	assert.Equal(t, []config.AcceptType{config.AcceptJSON, config.AcceptForm}, p.accepts)
}

func TestParser_Parse_JSON(t *testing.T) {
	t.Run("valid JSON", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptJSON})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"Alice","age":30}`))
		req.Header.Set("Content-Type", "application/json")

		result, err := p.Parse(req)
		require.NoError(t, err)
		assert.Equal(t, "Alice", result["name"])
		assert.Equal(t, float64(30), result["age"])
	})

	t.Run("no Content-Type defaults to JSON", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptJSON})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key":"value"}`))

		result, err := p.Parse(req)
		require.NoError(t, err)
		assert.Equal(t, "value", result["key"])
	})

	t.Run("no Content-Type without JSON accept", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptForm})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key":"value"}`))

		_, err := p.Parse(req)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnsupportedMediaType)
	})

	t.Run("JSON not accepted", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptForm})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"key":"value"}`))
		req.Header.Set("Content-Type", "application/json")

		_, err := p.Parse(req)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnsupportedMediaType)
	})

	t.Run("invalid JSON", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptJSON})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{invalid}`))
		req.Header.Set("Content-Type", "application/json")

		_, err := p.Parse(req)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrBadRequest)
	})
}

func TestParser_Parse_FormURLEncoded(t *testing.T) {
	t.Run("valid form data", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptForm})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("name=Alice&age=30"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		result, err := p.Parse(req)
		require.NoError(t, err)
		assert.Equal(t, "Alice", result["name"])
		assert.Equal(t, "30", result["age"])
	})

	t.Run("form not accepted", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptJSON})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("name=Alice"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		_, err := p.Parse(req)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnsupportedMediaType)
	})

	t.Run("multiple values", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptForm})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("tags=a&tags=b&tags=c"))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		result, err := p.Parse(req)
		require.NoError(t, err)
		tags, ok := result["tags"].([]any)
		require.True(t, ok)
		assert.Len(t, tags, 3)
	})
}

func TestParser_Parse_Multipart(t *testing.T) {
	t.Run("valid multipart form", func(t *testing.T) {
		body := "--boundary\r\n" +
			"Content-Disposition: form-data; name=\"field1\"\r\n\r\n" +
			"value1\r\n" +
			"--boundary\r\n" +
			"Content-Disposition: form-data; name=\"field2\"\r\n\r\n" +
			"value2\r\n" +
			"--boundary--\r\n"

		p := NewParser([]config.AcceptType{config.AcceptForm})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

		result, err := p.Parse(req)
		require.NoError(t, err)
		assert.Equal(t, "value1", result["field1"])
		assert.Equal(t, "value2", result["field2"])
	})

	t.Run("multipart not accepted", func(t *testing.T) {
		body := "--boundary\r\n" +
			"Content-Disposition: form-data; name=\"field1\"\r\n\r\n" +
			"value1\r\n" +
			"--boundary--\r\n"

		p := NewParser([]config.AcceptType{config.AcceptJSON})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

		_, err := p.Parse(req)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrUnsupportedMediaType)
	})

	t.Run("missing boundary", func(t *testing.T) {
		p := NewParser([]config.AcceptType{config.AcceptForm})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
		req.Header.Set("Content-Type", "multipart/form-data")

		_, err := p.Parse(req)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrBadRequest)
	})

	t.Run("skips file uploads", func(t *testing.T) {
		body := "--boundary\r\n" +
			"Content-Disposition: form-data; name=\"file\"; filename=\"test.txt\"\r\n" +
			"Content-Type: text/plain\r\n\r\n" +
			"file content\r\n" +
			"--boundary\r\n" +
			"Content-Disposition: form-data; name=\"field\"\r\n\r\n" +
			"value\r\n" +
			"--boundary--\r\n"

		p := NewParser([]config.AcceptType{config.AcceptForm})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

		result, err := p.Parse(req)
		require.NoError(t, err)
		assert.Equal(t, "value", result["field"])
		assert.Nil(t, result["file"])
	})

	t.Run("multiple values same name", func(t *testing.T) {
		body := "--boundary\r\n" +
			"Content-Disposition: form-data; name=\"items\"\r\n\r\n" +
			"item1\r\n" +
			"--boundary\r\n" +
			"Content-Disposition: form-data; name=\"items\"\r\n\r\n" +
			"item2\r\n" +
			"--boundary--\r\n"

		p := NewParser([]config.AcceptType{config.AcceptForm})
		req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
		req.Header.Set("Content-Type", "multipart/form-data; boundary=boundary")

		result, err := p.Parse(req)
		require.NoError(t, err)
		items, ok := result["items"].([]any)
		require.True(t, ok)
		assert.Len(t, items, 2)
	})
}

func TestParser_Parse_UnsupportedMediaType(t *testing.T) {
	p := NewParser([]config.AcceptType{config.AcceptJSON, config.AcceptForm})
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set("Content-Type", "text/plain")

	_, err := p.Parse(req)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrUnsupportedMediaType)
}

func TestParser_Parse_InvalidContentType(t *testing.T) {
	p := NewParser([]config.AcceptType{config.AcceptJSON})
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set("Content-Type", "invalid content type;;;")

	_, err := p.Parse(req)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrBadRequest)
}

func TestReadBody(t *testing.T) {
	t.Run("normal read", func(t *testing.T) {
		data, err := readBody(strings.NewReader("hello"), "")
		require.NoError(t, err)
		assert.Equal(t, []byte("hello"), data)
	})

	t.Run("body too large", func(t *testing.T) {
		largeData := make([]byte, MaxBodySize+1)
		_, err := readBody(bytes.NewReader(largeData), "")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrBodyTooLarge)
	})

	t.Run("read error", func(t *testing.T) {
		_, err := readBody(&errorReader{}, "")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrBadRequest)
	})

	t.Run("charset conversion", func(t *testing.T) {
		// UTF-8 encoded data with charset hint
		data, err := readBody(strings.NewReader("hello"), "utf-8")
		require.NoError(t, err)
		assert.Equal(t, []byte("hello"), data)
	})

	t.Run("invalid charset", func(t *testing.T) {
		_, err := readBody(strings.NewReader("hello"), "invalid-charset-xyz")
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrBadRequest)
	})
}

// errorReader is a reader that always returns an error
type errorReader struct{}

func (e *errorReader) Read(p []byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}
