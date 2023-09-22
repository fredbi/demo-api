package commands

import (
	"log"

	"github.com/fredbi/demo-api/app"
	"github.com/spf13/cobra"
)

// Execute the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

var (
	rootCmd = &cobra.Command{
		Use:   "images",
		Short: "images serves an API to create and manipulate images",
	}

	serveCmd = &cobra.Command{
		Use:   "serve",
		Short: "starts the API server",
		RunE: func(_ *cobra.Command, _ []string) error {
			config
			server := app.New()
			return server.Start()
		},
	}
)

func init() {
	// NOTE: in a production environment, this initialization typically
	// collects configuration from either configuration file, env vars or
	// flags and injects these into the app.
	rootCmd.AddCommand(serveCmd)
}
