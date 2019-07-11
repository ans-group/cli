package helper_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/helper"
)

func TestGetDestinationFilePath_DestinationNotDirectory_ReturnsDestination(t *testing.T) {
	fs := afero.NewMemMapFs()
	path, err := helper.GetDestinationFilePath(fs, "/path/to/file.txt", "/dest/dir")

	assert.Nil(t, err)
	assert.Equal(t, "/dest/dir", filepath.ToSlash(path))
}

func TestGetDestinationFilePath_DestinationDirectory_ReturnsDestinationWithFilename(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.Mkdir("/dest/dir", 655)
	path, err := helper.GetDestinationFilePath(fs, "/path/to/file.txt", "/dest/dir")

	assert.Nil(t, err)
	assert.Equal(t, "/dest/dir/file.txt", filepath.ToSlash(path))
}

func TestGetDestinationFilePath_DestinationOmitted_ReturnsPwdWithFilename(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.Mkdir("/dest/dir", 655)
	path, err := helper.GetDestinationFilePath(fs, "/path/to/file.txt", "")

	cwd, _ := os.Getwd()

	assert.Nil(t, err)
	assert.Equal(t, filepath.ToSlash(cwd+"/file.txt"), filepath.ToSlash(path))
}
