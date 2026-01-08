// Package commands provides the command-line interface for sql-http-proxy.
package commands

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path/filepath"

	"github.com/urfave/cli/v3"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/db"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

// Version is set by goreleaser via ldflags.
var Version = "dev"

// MakeApp creates a new CLI application instance.
func MakeApp() *cli.Command {
	return &cli.Command{
		Name:    "sql-http-proxy",
		Usage:   "YAML configuration-based HTTP to SQL proxy server",
		Version: Version,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   ".sql-http-proxy.yaml",
				Usage:   "Config YAML file",
			},
			&cli.StringFlag{
				Name:    "listen",
				Aliases: []string{"l"},
				Value:   ":8080",
				Usage:   "HTTP host:port",
			},
			&cli.BoolFlag{
				Name:  "validate",
				Usage: "Validate config and transforms without starting server",
			},
		},
		Action: action,
	}
}

func action(ctx context.Context, cmd *cli.Command) error {
	filename := cmd.String("config")
	listen := cmd.String("listen")
	validateOnly := cmd.Bool("validate")

	slog.Info("Parsing configuration")
	cfg, err := config.ParseFile(filename)
	if err != nil {
		return err
	}

	slog.Info("Validating transforms")
	if err := cfg.ValidateTransforms(); err != nil {
		return err
	}

	if validateOnly {
		slog.Info("Validation successful")
		return nil
	}

	slog.Info("Connecting to database")
	conn, err := db.Connect(cfg)
	if err != nil {
		return err
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
	} else {
		slog.Info("All endpoints use mock - skipping database connection")
	}

	slog.Info("Creating HTTP handlers")
	configDir := filepath.Dir(filename)
	mux, err := server.NewServeMux(conn, cfg, configDir)
	if err != nil {
		return err
	}

	slog.Info("Launching HTTP server", "address", listen)
	if err := http.ListenAndServe(listen, mux); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// App is the main CLI application.
var App = MakeApp()
