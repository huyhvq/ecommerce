package cmd

import (
	"github.com/huyhvq/ecommerce/cmd/api"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pd",
	Short: "Product discovery service",
	Long:  "Product discovery service",
	// Run: func(cmd *cobra.Command, args []string) { },
}

func init() {
	rootCmd.AddCommand(api.ServeCmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
