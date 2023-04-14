package util

import (
	"ariadne/util/diff"
	"os"
)

func Diff(path1 string, path2 string) (string, error) {
	bytes1, err := os.ReadFile(path1)
	if err != nil {
		return "", err
	}

	bytes2, err := os.ReadFile(path2)
	if err != nil {
		return "", err
	}

	diffBytes := diff.Diff(path1, bytes1, path2, bytes2)

	return string(diffBytes), nil
}
