package main

import (
	"github.com/mpyw/sql-http-proxy/cmd/sql-http-proxy/serve"
	"github.com/spf13/cobra"
	"log"
)

func main() {
	root := &cobra.Command{
		Use:     "sql-http-proxy",
		Short:   "sql-http-proxy is a JSON configuration-based HTTP to SQL proxy server",
		Version: "0.0.1",
	}
	root.AddCommand(serve.NewCommand())
	if err := root.Execute(); err != nil {
		log.Fatalln(err)
		return
	}
}
