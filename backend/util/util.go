package util

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
)

var validEnvironments = map[string]bool{
	"dev":  true,
	"prod": true,
}

func Diff(path1 string, path2 string) (string, error) {
	diffCmd := "diff"

	outputBytes, err := exec.Command(diffCmd, "-r", path1, path2).CombinedOutput()
	if err != nil {
		switch err.(type) {
		case *exec.ExitError:
			// `diff` ran successfully with non-zero exit code.  Report the
			// differences.
		default:
			// `diff` command failed to run.
			return "", err
		}
	}

	return string(outputBytes), nil
}

// InitSentry initializes the Sentry SDK and sets up error reporting.
// It uses the SENTRY_DSN and SENTRY_ENVIRONMENT environment variables to configure the SDK.
func InitSentry() {
	// Get the values of the SENTRY_DSN and SENTRY_ENVIRONMENT environment variables
	dsn := os.Getenv("SENTRY_DSN")
	environment := strings.ToLower(os.Getenv("SENTRY_ENVIRONMENT"))

	if dsn == "" || environment == "" {
		// SENTRY_DSN or SENTRY_ENVIRONMENT environment variable is not set
		log.Println("Both SENTRY_DSN and SENTRY_ENVIRONMENT environment variables must be set with valid values for Sentry to be initialized")
	} else if !validEnvironments[environment] {
		log.Printf("Invalid value for SENTRY_ENVIRONMENT: %s\n", environment)
	} else {
		// Initialize the Sentry Go SDK and set up error reporting
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              dsn,
			Environment:      environment,
			AttachStacktrace: true,
		})
		if err != nil {
			// Failed to initialize the Sentry SDK
			log.Printf("sentry.Init: %s", err)
		}
	}
}

// RecoverWithSentry is a utility function that recovers from panics and sends the panic information to Sentry.
// It should be called in a defer statement in your main function or in the top-level function of a goroutine.
// https://docs.sentry.io/platforms/go/panics/
func RecoverWithSentry() {
	if r := recover(); r != nil {
		sentry.CurrentHub().Recover(r)
		sentry.Flush(2 * time.Second)
		// Optionally, re-panic after sending the panic information to Sentry
		panic(r)
	}
}
