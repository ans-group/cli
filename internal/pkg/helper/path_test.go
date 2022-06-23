package helper_test

import (
	"path/filepath"
	"testing"

	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetDestinationFilePath_MissingSource_ReturnsError(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.Mkdir("/dest/dir", 655)
	_, err := helper.GetDestinationFilePath(fs, "", "")

	assert.NotNil(t, err)
}

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

func TestGetDestinationFilePath_DestinationOmitted_ReturnsFilename(t *testing.T) {
	fs := afero.NewMemMapFs()
	fs.Mkdir("/dest/dir", 655)
	path, err := helper.GetDestinationFilePath(fs, "/path/to/file.txt", "")

	assert.Nil(t, err)
	assert.Equal(t, "file.txt", path)
}
