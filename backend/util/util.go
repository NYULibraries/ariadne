package util

import (
	"github.com/getsentry/sentry-go"
	"log"
	"os"
	"os/exec"
	"time"
)

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

func InitSentry() {
	// Get the value of the SENTRY_DSN environment variable
	dsn := os.Getenv("SENTRY_DSN")
	if dsn == "" {
		// SENTRY_DSN environment variable is not set
		log.Println("SENTRY_DSN environment variable is not set")
	} else {
		// SENTRY_DSN environment variable is set
		log.Println("SENTRY_DSN environment variable is set")
		// Initialize the Sentry Go SDK and set up error reporting
		err := sentry.Init(sentry.ClientOptions{
			Dsn:   dsn,
			Debug: true, // Enable debugging mode
		})
		if err != nil {
			// Failed to initialize the Sentry SDK
			log.Fatalf("sentry.Init: %s", err)
		}
		// Flush buffered events before the program terminates
		defer sentry.Flush(2 * time.Second)
	}
	sentry.CaptureMessage("It works!")
}
