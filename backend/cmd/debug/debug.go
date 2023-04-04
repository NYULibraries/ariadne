package debug

import (
	"ariadne/log"
	"github.com/spf13/cobra"
)

var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debugging utilities",
}

func init() {
	log.SetLevel(log.LevelError)
}
