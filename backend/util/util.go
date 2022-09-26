package util

import (
	"os/exec"
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
