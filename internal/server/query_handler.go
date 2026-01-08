package server

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"net/url"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/js"
)

// CreateHandler creates an HTTP handler for a query.
// db can be nil if mock is configured.
func CreateHandler(db *sqlx.DB, query config.Query) (http.Handler, error) {
	return CreateHandlerWithOptions(db, query, HandlerOptions{})
}

// CreateHandlerWithOptions creates an HTTP handler with options.
// db can be nil if mock is configured.
func CreateHandlerWithOptions(db *sqlx.DB, query config.Query, opts HandlerOptions) (http.Handler, error) {
	// Pre-compile transformers
	var preTransform, mockTransform *js.Transformer
	var postTransformEach, postTransformAll *js.Transformer
	var err error

	if query.Transform != nil {
		if query.Transform.Pre != "" {
			preTransform, err = js.CompilePre(query.Transform.Pre)
			if err != nil {
				return nil, err
			}
		}
		if query.Transform.Mock != "" {
			mockTransform, err = js.CompilePre(query.Transform.Mock)
			if err != nil {
				return nil, err
			}
		}
		if query.Transform.Post.Each != "" {
			postTransformEach, err = js.Compile(query.Transform.Post.Each, []string{"ctx", "input", "output"})
			if err != nil {
				return nil, err
			}
		}
		if query.Transform.Post.All != "" {
			postTransformAll, err = js.Compile(query.Transform.Post.All, []string{"ctx", "input", "output"})
			if err != nil {
				return nil, err
			}
		}
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Accepting request: %s\n", r.RequestURI)
		res := &responder{w: w}
		params := extractNamedParams(r.URL.Query())

		// Shared context for pre/mock/post transforms (mutable)
		ctx := make(map[string]any)

		// Keep original params for post-transform
		originalParams := make(map[string]any, len(params))
		for k, v := range params {
			originalParams[k] = v
		}

		// SQL can be modified by pre-transform
		currentSQL := query.SQL

		// Apply pre-transform to params (with access to ctx and sql as free variables)
		if preTransform != nil {
			result, err := preTransform.ApplyPre(ctx, currentSQL, params)
			if err != nil {
				handleTransformError(res, err, http.StatusBadRequest)
				return
			}
			params = result.Output
			currentSQL = result.SQL
			ctx = result.Ctx
		}

		// Use mock if configured, otherwise query DB
		if mockTransform != nil {
			handleMockQuery(res, mockTransform, ctx, currentSQL, params, originalParams, query, postTransformEach, postTransformAll, opts.Recorder)
		} else {
			handleDBQuery(res, db, r.Context(), ctx, currentSQL, params, originalParams, query, postTransformEach, postTransformAll, opts.Recorder)
		}
	}), nil
}

// handleMockQuery handles a query using mock data.
func handleMockQuery(res *responder, mockTransform *js.Transformer, ctx map[string]any, sqlStr string, params, originalParams map[string]any, query config.Query, postTransformEach, postTransformAll *js.Transformer, recorder QueryRecorder) {
	log.Printf("Type: %s (Mock), SQL: %s, Params: %+v\n", query.Type, sqlStr, params)

	if recorder != nil {
		recorder(QueryRecord{
			Type:   string(query.Type),
			IsMock: true,
			SQL:    sqlStr,
			Params: params,
		})
	}

	mockResult, err := mockTransform.ApplyMock(ctx, sqlStr, params)
	if err != nil {
		handleTransformError(res, err, http.StatusInternalServerError)
		return
	}
	ctx = mockResult.Ctx

	switch query.Type {
	case config.QueryTypeOne:
		// Mock should return object or null for "one"
		var entry map[string]any
		if mockResult.Output == nil {
			if !query.HandleNotFound {
				res.Error(http.StatusNotFound, errors.New("not found"))
				return
			}
			entry = nil
		} else if m, ok := mockResult.Output.(map[string]any); ok {
			entry = m
		} else {
			res.Error(http.StatusInternalServerError, errors.New("mock for 'one' query must return object or null"))
			return
		}

		var result any = entry
		if postTransformAll != nil {
			result, err = postTransformAll.Apply(ctx, originalParams, entry)
			if err != nil {
				handleTransformError(res, err, http.StatusInternalServerError)
				return
			}
		}
		res.Respond(http.StatusOK, result)

	case config.QueryTypeMany:
		// Mock should return array for "many"
		var entries []map[string]any
		if arr, ok := mockResult.Output.([]any); ok {
			entries = make([]map[string]any, len(arr))
			for i, item := range arr {
				if m, ok := item.(map[string]any); ok {
					entries[i] = m
				} else {
					res.Error(http.StatusInternalServerError, errors.New("mock for 'many' query must return array of objects"))
					return
				}
			}
		} else if mockResult.Output == nil {
			entries = []map[string]any{}
		} else {
			res.Error(http.StatusInternalServerError, errors.New("mock for 'many' query must return array"))
			return
		}

		var result any = entries
		if postTransformEach != nil {
			eachResult, err := postTransformEach.ApplyToEachRow(ctx, originalParams, entries)
			if err != nil {
				handleTransformError(res, err, http.StatusInternalServerError)
				return
			}
			result = eachResult

			if postTransformAll != nil {
				entriesFromEach := make([]map[string]any, len(eachResult))
				for i, r := range eachResult {
					if m, ok := r.(map[string]any); ok {
						entriesFromEach[i] = m
					} else {
						res.Error(http.StatusInternalServerError, errors.New("post.each transform must return object"))
						return
					}
				}
				result, err = postTransformAll.ApplyToAllRows(ctx, originalParams, entriesFromEach)
				if err != nil {
					handleTransformError(res, err, http.StatusInternalServerError)
					return
				}
			}
		} else if postTransformAll != nil {
			result, err = postTransformAll.ApplyToAllRows(ctx, originalParams, entries)
			if err != nil {
				handleTransformError(res, err, http.StatusInternalServerError)
				return
			}
		}
		res.Respond(http.StatusOK, result)

	default:
		res.Error(http.StatusInternalServerError, errors.New("unsupported query type"))
	}
}

// handleDBQuery handles a query using the database.
func handleDBQuery(res *responder, db *sqlx.DB, reqCtx context.Context, ctx map[string]any, currentSQL string, params, originalParams map[string]any, query config.Query, postTransformEach, postTransformAll *js.Transformer, recorder QueryRecorder) {
	// Convert named parameters to positional
	boundSQL, args, err := sqlx.Named(currentSQL, params)
	if err != nil {
		res.Error(http.StatusBadRequest, err)
		return
	}
	boundSQL = db.Rebind(boundSQL)

	if recorder != nil {
		recorder(QueryRecord{
			Type:   string(query.Type),
			IsMock: false,
			SQL:    boundSQL,
			Args:   args,
		})
	}

	switch query.Type {
	case config.QueryTypeOne:
		log.Printf("Type: One, Query: %s, Params: %+v\n", boundSQL, args)
		row := db.QueryRowxContext(reqCtx, boundSQL, args...)
		if row.Err() != nil {
			res.Error(http.StatusInternalServerError, row.Err())
			return
		}
		var entry map[string]any
		scanned := map[string]any{}
		if err := row.MapScan(scanned); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				if !query.HandleNotFound {
					res.Error(http.StatusNotFound, errors.New("not found"))
					return
				}
				entry = nil
			} else {
				res.Error(http.StatusInternalServerError, err)
				return
			}
		} else {
			entry = scanned
		}

		var result any = entry
		if postTransformAll != nil {
			result, err = postTransformAll.Apply(ctx, originalParams, entry)
			if err != nil {
				handleTransformError(res, err, http.StatusInternalServerError)
				return
			}
		}
		res.Respond(http.StatusOK, result)

	case config.QueryTypeMany:
		log.Printf("Type: Many, Query: %s, Params: %+v\n", boundSQL, args)
		rows, err := db.QueryxContext(reqCtx, boundSQL, args...)
		if err != nil {
			res.Error(http.StatusInternalServerError, err)
			return
		}
		defer func() { _ = rows.Close() }()

		entries := make([]map[string]any, 0)
		for rows.Next() {
			entry := map[string]any{}
			if err := rows.MapScan(entry); err != nil {
				res.Error(http.StatusInternalServerError, err)
				return
			}
			entries = append(entries, entry)
		}
		if err := rows.Err(); err != nil {
			res.Error(http.StatusInternalServerError, err)
			return
		}

		var result any = entries
		if postTransformEach != nil {
			eachResult, err := postTransformEach.ApplyToEachRow(ctx, originalParams, entries)
			if err != nil {
				handleTransformError(res, err, http.StatusInternalServerError)
				return
			}
			result = eachResult

			if postTransformAll != nil {
				entriesFromEach := make([]map[string]any, len(eachResult))
				for i, r := range eachResult {
					if m, ok := r.(map[string]any); ok {
						entriesFromEach[i] = m
					} else {
						res.Error(http.StatusInternalServerError, errors.New("post.each transform must return object"))
						return
					}
				}
				result, err = postTransformAll.ApplyToAllRows(ctx, originalParams, entriesFromEach)
				if err != nil {
					handleTransformError(res, err, http.StatusInternalServerError)
					return
				}
			}
		} else if postTransformAll != nil {
			result, err = postTransformAll.ApplyToAllRows(ctx, originalParams, entries)
			if err != nil {
				handleTransformError(res, err, http.StatusInternalServerError)
				return
			}
		}
		res.Respond(http.StatusOK, result)

	default:
		res.Error(http.StatusInternalServerError, errors.New("unsupported query type"))
	}
}

func extractNamedParams(values url.Values) map[string]any {
	params := make(map[string]any, len(values))
	for key, vals := range values {
		if len(vals) > 0 {
			params[key] = vals[0]
		}
	}
	return params
}
