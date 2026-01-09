// Package commands provides the command-line interface for sql-http-proxy.
package commands

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/urfave/cli/v3"

	"github.com/mpyw/sql-http-proxy/internal/config"
	"github.com/mpyw/sql-http-proxy/internal/db"
	"github.com/mpyw/sql-http-proxy/internal/server"
)

// ShutdownTimeout is the maximum time to wait for graceful shutdown.
const ShutdownTimeout = 30 * time.Second

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

	configDir := filepath.Dir(filename)

	if validateOnly {
		slog.Info("Validation successful")
		return nil
	}

	slog.Info("Connecting to database")
	conn, err := db.Connect(cfg, configDir)
	if err != nil {
		return err
	}
	if conn != nil {
		defer func() {
			if err := conn.Close(); err != nil {
				slog.Warn("Failed to close database connection", "error", err)
			}
		}()
	} else {
		slog.Info("All endpoints use mock - skipping database connection")
	}

	slog.Info("Creating HTTP handlers")
	mux, err := server.NewServeMux(conn, cfg, configDir)
	if err != nil {
		return err
	}

	srv := &http.Server{
		Addr:    listen,
		Handler: mux,
	}

	// Channel to receive server errors
	serverErr := make(chan error, 1)

	// Start server in goroutine
	go func() {
		slog.Info("Launching HTTP server", "address", listen)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("server error: %w", err)
		}
		close(serverErr)
	}()

	// Wait for interrupt signal or server error
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-serverErr:
		return err
	case sig := <-quit:
		slog.Info("Received shutdown signal", "signal", sig)
	}

	// Graceful shutdown with timeout
	shutdownCtx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()

	slog.Info("Shutting down server...")
	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("server shutdown error: %w", err)
	}

	slog.Info("Server stopped gracefully")
	return nil
}

// App is the main CLI application.
var App = MakeApp()
