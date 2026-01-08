// Package commands provides the command-line interface for sql-http-proxy.
package commands

import (
	"context"
	"fmt"
	"log"
	"net/http"

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

	log.Println("Parsing configuration")
	cfg, err := config.ParseFile(filename)
	if err != nil {
		return err
	}

	log.Println("Validating transforms")
	if err := cfg.ValidateTransforms(); err != nil {
		return err
	}

	if validateOnly {
		log.Println("Validation successful")
		return nil
	}

	log.Println("Connecting to database")
	conn, err := db.Connect(cfg)
	if err != nil {
		return err
	}
	if conn != nil {
		defer func() {
			_ = conn.Close()
		}()
	} else {
		log.Println("All endpoints use mock - skipping database connection")
	}

	log.Println("Creating HTTP handlers")
	mux, err := server.NewServeMux(conn, cfg)
	if err != nil {
		return err
	}

	log.Printf("Launching HTTP server on %s\n", listen)
	if err := http.ListenAndServe(listen, mux); err != nil {
		return fmt.Errorf("server error: %w", err)
	}
	return nil
}

// App is the main CLI application.
var App = MakeApp()
