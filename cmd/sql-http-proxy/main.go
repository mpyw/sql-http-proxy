package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/mpyw/sql-http-proxy/internal/cli/commands"
)

func main() {
	if err := commands.App.Run(context.Background(), os.Args); err != nil {
		slog.Error("Application error", "error", err)
		os.Exit(1)
	}
}
