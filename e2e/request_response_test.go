package e2e

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestRequestObject(t *testing.T) {
	cfg, err := config.ParseFile("request_response_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("request.method", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/method", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"method":"GET"}`, rec.Body.String())
	})

	t.Run("request.url", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/url?foo=bar", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"url":"/test/request/url?foo=bar"}`, rec.Body.String())
	})

	t.Run("request.headers.get()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/headers/get", nil)
		req.Header.Set("X-Custom-Header", "custom-value")
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{
			"customHeader": "custom-value",
			"missingHeader": null,
			"contentType": "application/json"
		}`, rec.Body.String())
	})

	t.Run("request.headers.has()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/headers/has", nil)
		req.Header.Set("X-Custom-Header", "exists")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"hasCustom": true, "hasMissing": false}`, rec.Body.String())
	})

	t.Run("request.headers.keys()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/headers/keys", nil)
		req.Header.Set("X-Header-A", "a")
		req.Header.Set("X-Header-B", "b")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		// Keys should include our custom headers (lowercase)
		assert.Contains(t, rec.Body.String(), "x-header-a")
		assert.Contains(t, rec.Body.String(), "x-header-b")
	})

	t.Run("request.headers.values()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/headers/values", nil)
		req.Header.Set("X-Test-Value", "test-value-123")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "test-value-123")
	})

	t.Run("request.headers.entries()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/headers/entries", nil)
		req.Header.Set("X-Entry-Test", "entry-value")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		// Entries should be [name, value] pairs
		assert.Contains(t, rec.Body.String(), "x-entry-test")
		assert.Contains(t, rec.Body.String(), "entry-value")
	})

	t.Run("request.headers.forEach()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/headers/foreach", nil)
		req.Header.Set("X-Foreach-Test", "foreach-value")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Contains(t, rec.Body.String(), "x-foreach-test")
		assert.Contains(t, rec.Body.String(), "foreach-value")
	})

	t.Run("request.headers is read-only", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/headers/readonly", nil)
		req.Header.Set("X-Custom-Header", "original")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"before": "original", "after": "original", "unchanged": true}`, rec.Body.String())
	})
}

func TestResponseObject(t *testing.T) {
	cfg, err := config.ParseFile("request_response_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("response.status get/set", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/status", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code) // 201
		require.JSONEq(t, `{"defaultStatus": 200, "newStatus": 201}`, rec.Body.String())
	})

	t.Run("response.statusText get/set", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/statustext", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code)
		require.JSONEq(t, `{"autoText": "Created", "customText": "Custom Created"}`, rec.Body.String())
	})

	t.Run("response.ok", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/ok", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{
			"okAt200": true,
			"okAt201": true,
			"okAt404": false,
			"okAt500": false
		}`, rec.Body.String())
	})

	t.Run("response.headers.set()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/headers/set", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "test-value", rec.Header().Get("X-Custom-Response"))
		assert.Equal(t, "no-cache", rec.Header().Get("Cache-Control"))
	})

	t.Run("response.headers.append()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/headers/append", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		// Multiple values for the same header
		values := rec.Header().Values("X-Multi")
		assert.Contains(t, values, "value1")
		assert.Contains(t, values, "value2")
	})

	t.Run("response.headers.delete()", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/headers/delete", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		assert.Empty(t, rec.Header().Get("X-To-Delete"))
		assert.Equal(t, "will-remain", rec.Header().Get("X-Kept"))
	})

	t.Run("response.headers iteration", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/headers/iteration", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		body := rec.Body.String()
		// Check that keys, values, entries contain our headers
		assert.Contains(t, body, "x-header-a")
		assert.Contains(t, body, "x-header-b")
		assert.Contains(t, body, "value-a")
		assert.Contains(t, body, "value-b")
	})

	t.Run("response.status invalid values ignored", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/status/invalid", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusAccepted, rec.Code) // 202 is the last valid status
		require.JSONEq(t, `{
			"before": 200,
			"after99": 200,
			"after600": 200,
			"after202": 202
		}`, rec.Body.String())
	})
}

func TestRequestInDifferentContexts(t *testing.T) {
	cfg, err := config.ParseFile("request_response_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("request in pre-transform", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/request/in-pre", nil)
		req.Header.Set("X-Test-Header", "exists")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"requestMethod": "GET", "hasHeader": true}`, rec.Body.String())
	})

	t.Run("response not available in pre-transform", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/test/response/not-in-pre", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"responseUndefined": true}`, rec.Body.String())
	})
}

func TestMutationRequestResponse(t *testing.T) {
	cfg, err := config.ParseFile("request_response_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("request/response in mutation", func(t *testing.T) {
		body := strings.NewReader(`{"name": "test"}`)
		req := httptest.NewRequest(http.MethodPost, "/test/mutation/request-response", body)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusCreated, rec.Code)
		require.JSONEq(t, `{
			"method": "POST",
			"contentType": "application/json"
		}`, rec.Body.String())
		assert.Equal(t, "/test/mutation/request-response/1", rec.Header().Get("Location"))
	})
}
