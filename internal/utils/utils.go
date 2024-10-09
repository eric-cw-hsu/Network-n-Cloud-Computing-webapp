package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetProjectRootPath() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for currentDir != "/" {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}
		currentDir = filepath.Dir(currentDir)
	}
	return "", fmt.Errorf("project root not found")
}
