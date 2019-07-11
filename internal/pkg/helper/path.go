package helper

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/afero"
)

// GetDestinationFilePath returns the destination path, given source file 'source' and optional destination path 'destination'
func GetDestinationFilePath(fs afero.Fs, source string, destination string) (string, error) {
	var targetFilePath string
	if len(destination) > 0 {
		targetFilePath = destination
		if ok, _ := afero.IsDir(fs, targetFilePath); ok {
			targetFilePath = filepath.Join(targetFilePath, filepath.Base(source))
		}
	} else {
		dir, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("Error determining current directory: %s", err)
		}
		targetFilePath = filepath.Join(dir, filepath.Base(source))
	}

	return filepath.Clean(targetFilePath), nil
}
