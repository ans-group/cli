package helper

import (
	"errors"
	"path/filepath"

	"github.com/spf13/afero"
)

// GetDestinationFilePath returns the destination path, given source file 'source' and optional destination path 'destination'
func GetDestinationFilePath(fs afero.Fs, source string, destination string) (string, error) {
	if len(source) < 1 {
		return "", errors.New("missing source")
	}

	var targetFilePath string
	if len(destination) > 0 {
		targetFilePath = destination
		if ok, _ := afero.IsDir(fs, targetFilePath); ok {
			targetFilePath = filepath.Join(targetFilePath, filepath.Base(source))
		}
	} else {
		targetFilePath = filepath.Base(source)
	}

	return filepath.Clean(targetFilePath), nil
}
