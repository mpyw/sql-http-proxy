package main

import (
	"context"
	"log"
	"os"

	"github.com/mpyw/sql-http-proxy/internal/cli/commands"
)

func main() {
	if err := commands.App.Run(context.Background(), os.Args); err != nil {
		log.Fatalln(err)
	}
}
