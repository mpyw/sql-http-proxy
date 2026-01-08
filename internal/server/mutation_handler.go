package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/js"
)

// CreateMutationHandler creates an HTTP handler for a mutation.
// db can be nil if mock is configured.
func CreateMutationHandler(db *sqlx.DB, mutation config.Mutation) (http.Handler, error) {
	return CreateMutationHandlerWithOptions(db, mutation, HandlerOptions{})
}

// CreateMutationHandlerWithOptions creates an HTTP handler for a mutation with options.
// db can be nil if mock is configured.
func CreateMutationHandlerWithOptions(db *sqlx.DB, mutation config.Mutation, opts HandlerOptions) (http.Handler, error) {
	// Pre-compile transformers
	var preTransform, mockTransform *js.Transformer
	var postTransformEach, postTransformAll *js.Transformer
	var err error

	if mutation.Transform != nil {
		if mutation.Transform.Pre != "" {
			preTransform, err = js.CompilePre(mutation.Transform.Pre)
			if err != nil {
				return nil, err
			}
		}
		if mutation.Transform.Mock != "" {
			mockTransform, err = js.CompilePre(mutation.Transform.Mock)
			if err != nil {
				return nil, err
			}
		}
		if mutation.Transform.Post.Each != "" {
			postTransformEach, err = js.Compile(mutation.Transform.Post.Each, []string{"ctx", "input", "output"})
			if err != nil {
				return nil, err
			}
		}
		if mutation.Transform.Post.All != "" {
			postTransformAll, err = js.Compile(mutation.Transform.Post.All, []string{"ctx", "input", "output"})
			if err != nil {
				return nil, err
			}
		}
	}

	expectedMethod := mutation.GetMethod()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check HTTP method
		if r.Method != expectedMethod {
			w.Header().Set("Allow", expectedMethod)
			res := &responder{w: w}
			res.Error(http.StatusMethodNotAllowed, errors.New("method not allowed"))
			return
		}

		log.Printf("Accepting request: %s %s\n", r.Method, r.RequestURI)
		res := &responder{w: w}

		// Read params from request body (JSON)
		params, err := extractBodyParams(r)
		if err != nil {
			res.Error(http.StatusBadRequest, err)
			return
		}

		// Shared context for pre/mock/post transforms (mutable)
		ctx := make(map[string]any)

		// Keep original params for post-transform
		originalParams := make(map[string]any, len(params))
		for k, v := range params {
			originalParams[k] = v
		}

		// SQL can be modified by pre-transform
		currentSQL := mutation.SQL

		// Apply pre-transform to params
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

		// Use mock if configured, otherwise execute mutation
		if mockTransform != nil {
			handleMockMutation(res, mockTransform, ctx, currentSQL, params, originalParams, mutation, postTransformEach, postTransformAll, opts.Recorder)
		} else {
			handleDBMutation(res, db, r.Context(), ctx, currentSQL, params, originalParams, mutation, postTransformEach, postTransformAll, opts.Recorder)
		}
	}), nil
}

// extractBodyParams reads JSON body and returns as map.
func extractBodyParams(r *http.Request) (map[string]any, error) {
	if r.Body == nil {
		return make(map[string]any), nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	if len(body) == 0 {
		return make(map[string]any), nil
	}

	var params map[string]any
	if err := json.Unmarshal(body, &params); err != nil {
		return nil, err
	}
	return params, nil
}

// handleMockMutation handles a mutation using mock data.
func handleMockMutation(res *responder, mockTransform *js.Transformer, ctx map[string]any, sqlStr string, params, originalParams map[string]any, mutation config.Mutation, postTransformEach, postTransformAll *js.Transformer, recorder QueryRecorder) {
	log.Printf("Type: %s (Mock Mutation), SQL: %s, Params: %+v\n", mutation.Type, sqlStr, params)

	if recorder != nil {
		recorder(QueryRecord{
			Type:   string(mutation.Type),
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

	switch mutation.Type {
	case config.MutationTypeNone:
		// No return value - warn if mock returns something
		if mockResult.Output != nil {
			log.Printf("Warning: mock for 'none' mutation returned a value, but it will be ignored")
		}
		res.Respond(http.StatusNoContent, nil)

	case config.MutationTypeOne:
		// Mock should return object or null for "one"
		var entry map[string]any
		if mockResult.Output == nil {
			entry = nil
		} else if m, ok := mockResult.Output.(map[string]any); ok {
			entry = m
		} else {
			res.Error(http.StatusInternalServerError, errors.New("mock for 'one' mutation must return object or null"))
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
		if result == nil {
			res.Respond(http.StatusNoContent, nil)
		} else {
			res.Respond(http.StatusOK, result)
		}

	case config.MutationTypeMany:
		// Mock should return array for "many"
		var entries []map[string]any
		if arr, ok := mockResult.Output.([]any); ok {
			entries = make([]map[string]any, len(arr))
			for i, item := range arr {
				if m, ok := item.(map[string]any); ok {
					entries[i] = m
				} else {
					res.Error(http.StatusInternalServerError, errors.New("mock for 'many' mutation must return array of objects"))
					return
				}
			}
		} else if mockResult.Output == nil {
			entries = []map[string]any{}
		} else {
			res.Error(http.StatusInternalServerError, errors.New("mock for 'many' mutation must return array"))
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
		res.Error(http.StatusInternalServerError, errors.New("unsupported mutation type"))
	}
}

// handleDBMutation handles a mutation using the database.
func handleDBMutation(res *responder, db *sqlx.DB, reqCtx context.Context, ctx map[string]any, currentSQL string, params, originalParams map[string]any, mutation config.Mutation, postTransformEach, postTransformAll *js.Transformer, recorder QueryRecorder) {
	// Convert named parameters to positional
	boundSQL, args, err := sqlx.Named(currentSQL, params)
	if err != nil {
		res.Error(http.StatusBadRequest, err)
		return
	}
	boundSQL = db.Rebind(boundSQL)

	if recorder != nil {
		recorder(QueryRecord{
			Type:   string(mutation.Type),
			IsMock: false,
			SQL:    boundSQL,
			Args:   args,
		})
	}

	switch mutation.Type {
	case config.MutationTypeNone:
		log.Printf("Type: None (Mutation), Query: %s, Params: %+v\n", boundSQL, args)
		_, err := db.ExecContext(reqCtx, boundSQL, args...)
		if err != nil {
			res.Error(http.StatusInternalServerError, err)
			return
		}
		res.Respond(http.StatusNoContent, nil)

	case config.MutationTypeOne:
		log.Printf("Type: One (Mutation), Query: %s, Params: %+v\n", boundSQL, args)
		row := db.QueryRowxContext(reqCtx, boundSQL, args...)
		if row.Err() != nil {
			res.Error(http.StatusInternalServerError, row.Err())
			return
		}
		var entry map[string]any
		scanned := map[string]any{}
		if err := row.MapScan(scanned); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
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
		if result == nil {
			res.Respond(http.StatusNoContent, nil)
		} else {
			res.Respond(http.StatusOK, result)
		}

	case config.MutationTypeMany:
		log.Printf("Type: Many (Mutation), Query: %s, Params: %+v\n", boundSQL, args)
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
		res.Error(http.StatusInternalServerError, errors.New("unsupported mutation type"))
	}
}
