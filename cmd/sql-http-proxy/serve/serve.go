package serve

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/mpyw/sql-http-proxy/config"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Launch a server",
		RunE:  Run,
	}
	cmd.Flags().StringP("config", "f", "sql-http-proxy.json", "config json file")
	cmd.Flags().StringP("listen", "l", ":8080", "HTTP host:port")
	return cmd
}

func Run(cmd *cobra.Command, _ []string) error {
	log.Println("Parsing configuration")
	filename, err := cmd.Flags().GetString("config")
	if err != nil {
		return fmt.Errorf("missing or wrong value or required flag \"--config\" \"-f\": %w", err)
	}
	listen, err := cmd.Flags().GetString("listen")
	if err != nil {
		return fmt.Errorf("missing or wrong value or required flag \"--listen\" \"-l\": %w", err)
	}
	cfg, err := config.ParseFile(filename)
	if err != nil {
		return err
	}
	driverName, err := cfg.Driver()
	if err != nil {
		return err
	}

	log.Println("Connecting to database")
	db, err := sqlx.Open(driverName, cfg.DSN)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	log.Printf("Launching HTTP server on %s\n", listen)
	mux := http.NewServeMux()
	for _, query := range cfg.Queries {
		mux.Handle(query.Path, CreateHandler(db, query))
	}
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		res := responder{w, r}
		res.Error(404, errors.New("not found"))
	}))
	return http.ListenAndServe(listen, mux)
}

func CreateHandler(db *sqlx.DB, query config.Query) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Accepting request: %s\n", r.RequestURI)
		res := responder{w, r}
		params := extractValuesFromUrl(r.URL.Query(), query.Argc)

		switch query.Type {
		case config.QueryTypeOne:
			log.Printf("Type: One, Query: %s, Params: %+v\n", query.SQL, params)
			row := db.QueryRowxContext(r.Context(), query.SQL, params...)
			if row.Err() != nil {
				res.Error(http.StatusInternalServerError, row.Err())
				return
			}
			entry := map[string]any{}
			if err := row.MapScan(entry); err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					res.Respond(http.StatusNotFound, nil)
				} else {
					res.Error(http.StatusInternalServerError, err)
				}
				return
			}
			res.Respond(http.StatusOK, entry)

		case config.QueryTypeMany:
			log.Printf("Type: Many, Query: %s, Params: %+v\n", query.SQL, params)
			rows, err := db.QueryxContext(r.Context(), query.SQL, params...)
			if err != nil {
				res.Error(http.StatusInternalServerError, err)
				return
			}
			entries := make([]map[string]any, 0)
			for rows.Next() {
				entry := map[string]any{}
				if err := rows.MapScan(entry); err != nil {
					res.Error(http.StatusInternalServerError, err)
					return
				}
				entries = append(entries, entry)
			}
			res.Respond(http.StatusOK, entries)

		default:
			res.Error(http.StatusInternalServerError, errors.New("unsupported query type"))
		}
	})
}
func extractValuesFromUrl(values url.Values, argc int) []any {
	entries := make([]any, 0, len(values))
	for i := 1; i <= argc; i++ {
		key := fmt.Sprintf("$%d", i)
		if value, ok := values[key]; ok && len(value) > 0 {
			entries = append(entries, value[0])
		}
	}
	return entries
}
