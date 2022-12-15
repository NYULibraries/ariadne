package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"ariadne/cmd/debug"
	"ariadne/cmd/server"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "ariadne",
	Long: "`ariadne`" + ` is the backend application for the NYU Libraries OpenURL link resolver.
Use ariadne to start the API server and to debug backend requests and responses.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(debug.DebugCmd)
	rootCmd.AddCommand(server.ServerCmd)
}
