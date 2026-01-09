package e2e

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

func TestPathRegexUUID(t *testing.T) {
	cfg, err := config.ParseFile("path_regex_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("any UUID accepts valid UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/uuid/550e8400-e29b-41d4-a716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":"550e8400-e29b-41d4-a716-446655440000","type":"any-uuid"}`, rec.Body.String())
	})

	t.Run("any UUID rejects invalid format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/uuid/not-a-uuid", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("any UUID rejects uppercase", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/uuid/550E8400-E29B-41D4-A716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("any UUID rejects short format", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/uuid/550e8400-e29b-41d4-a716", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPathRegexUUIDv4(t *testing.T) {
	cfg, err := config.ParseFile("path_regex_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("UUIDv4 accepts valid v4 UUID", func(t *testing.T) {
		// Valid UUIDv4: version=4, variant=1 (8, 9, a, b)
		req := httptest.NewRequest(http.MethodGet, "/uuid-v4/550e8400-e29b-41d4-a716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":"550e8400-e29b-41d4-a716-446655440000","type":"uuid-v4"}`, rec.Body.String())
	})

	t.Run("UUIDv4 accepts variant 8", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/uuid-v4/550e8400-e29b-41d4-8716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("UUIDv4 accepts variant 9", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/uuid-v4/550e8400-e29b-41d4-9716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("UUIDv4 accepts variant b", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/uuid-v4/550e8400-e29b-41d4-b716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("UUIDv4 rejects v7 UUID", func(t *testing.T) {
		// UUIDv7: version=7
		req := httptest.NewRequest(http.MethodGet, "/uuid-v4/018f6b2e-7d3c-7000-8000-000000000000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("UUIDv4 rejects v1 UUID", func(t *testing.T) {
		// UUIDv1: version=1
		req := httptest.NewRequest(http.MethodGet, "/uuid-v4/550e8400-e29b-11d4-a716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("UUIDv4 rejects invalid variant", func(t *testing.T) {
		// variant=c is not valid for RFC 4122
		req := httptest.NewRequest(http.MethodGet, "/uuid-v4/550e8400-e29b-41d4-c716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPathRegexUUIDv7(t *testing.T) {
	cfg, err := config.ParseFile("path_regex_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("UUIDv7 accepts valid v7 UUID", func(t *testing.T) {
		// Valid UUIDv7: version=7, variant=1
		req := httptest.NewRequest(http.MethodGet, "/uuid-v7/018f6b2e-7d3c-7000-8000-000000000000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":"018f6b2e-7d3c-7000-8000-000000000000","type":"uuid-v7"}`, rec.Body.String())
	})

	t.Run("UUIDv7 accepts variant a", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/uuid-v7/018f6b2e-7d3c-7000-a000-000000000000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("UUIDv7 rejects v4 UUID", func(t *testing.T) {
		// UUIDv4: version=4
		req := httptest.NewRequest(http.MethodGet, "/uuid-v7/550e8400-e29b-41d4-a716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("UUIDv7 rejects invalid variant", func(t *testing.T) {
		// variant=0 is not valid for RFC 4122
		req := httptest.NewRequest(http.MethodGet, "/uuid-v7/018f6b2e-7d3c-7000-0000-000000000000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPathRegexNumeric(t *testing.T) {
	cfg, err := config.ParseFile("path_regex_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("numeric accepts digits", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/numeric/12345", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":12345,"type":"numeric"}`, rec.Body.String())
	})

	t.Run("numeric accepts zero", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/numeric/0", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"id":0,"type":"numeric"}`, rec.Body.String())
	})

	t.Run("numeric rejects letters", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/numeric/abc", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("numeric rejects mixed", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/numeric/123abc", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("numeric rejects negative", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/numeric/-123", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("numeric rejects UUID", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/numeric/550e8400-e29b-41d4-a716-446655440000", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}

func TestPathRegexSlug(t *testing.T) {
	cfg, err := config.ParseFile("path_regex_test.yaml")
	require.NoError(t, err)

	mux, err := server.NewServeMux(nil, cfg, ".")
	require.NoError(t, err)

	t.Run("slug accepts lowercase", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slug/hello-world", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"slug":"hello-world","type":"slug"}`, rec.Body.String())
	})

	t.Run("slug accepts numbers", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slug/post-123", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"slug":"post-123","type":"slug"}`, rec.Body.String())
	})

	t.Run("slug accepts single word", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slug/hello", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.JSONEq(t, `{"slug":"hello","type":"slug"}`, rec.Body.String())
	})

	t.Run("slug rejects uppercase", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slug/Hello-World", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("slug rejects underscores", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slug/hello_world", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})

	t.Run("slug rejects spaces", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/slug/hello%20world", nil)
		rec := httptest.NewRecorder()

		mux.ServeHTTP(rec, req)
		require.Equal(t, http.StatusNotFound, rec.Code)
	})
}
