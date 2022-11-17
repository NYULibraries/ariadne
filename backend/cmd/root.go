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
	Use:   "ariadne",
	Short: "[TODO: short description]",
	Long: `[TODO: long multiline description]:

[Line 1]
[Line 2]
[Line 3]`,
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
