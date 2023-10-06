package cmd

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "feta",
	Short: "Feta API is the API for the testing website Feta",
	Run: func(cmd *cobra.Command, args []string) {
		// nothing
	},
}

func Execute(ctx context.Context) {
	rootCmd.AddCommand(APICmd(ctx))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
