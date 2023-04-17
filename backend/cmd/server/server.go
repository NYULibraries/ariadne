package server

import (
	"ariadne/api"
	"ariadne/log"
	"fmt"
	"github.com/spf13/cobra"
	"net/http"
	"strings"
)

const defaultPort = "8080"

var loggingLevel string
var port string

var ServerCmd = &cobra.Command{
	Use:     "server",
	Short:   "Start API server",
	Example: "ariadne server",
	Run: func(cmd *cobra.Command, args []string) {
		start()
	},
}

func init() {
	ServerCmd.Flags().StringVarP(&loggingLevel, "logging-level", "l",
		log.DefaultLevelStringOption,
		"Sets logging level: "+strings.Join(log.GetValidLevelOptionStrings(), ", ")+"")
	ServerCmd.Flags().StringVarP(&port, "port", "p", defaultPort, "Port to run server on")
}

func start() {
	router := api.NewRouter()

	normalizedLogLevel := strings.ToLower(loggingLevel)
	err := log.SetLevelByString(normalizedLogLevel)
	if err != nil {
		log.Fatal(api.MessageKey, err)
	}

	log.Info(api.MessageKey, fmt.Sprintf("Logging level set to \"%s\"", normalizedLogLevel))

	log.Info(api.MessageKey, "Listening on port "+port)

	log.Fatal(api.MessageKey, http.ListenAndServe(":"+port, router))
}
