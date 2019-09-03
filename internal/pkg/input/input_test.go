package input

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test/test_input"
)

func TestReadInput(t *testing.T) {
	t.Run("StdinReadError_ReturnsError", func(t *testing.T) {
		oldReader := InputReader
		defer func() { InputReader = oldReader }()

		readErr := errors.New("test error")
		InputReader = func() io.Reader {
			return &test_input.TestReadCloser{
				ReadError: readErr,
			}
		}

		_, err := ReadInput("test")

		assert.NotNil(t, err)
		assert.Equal(t, "Error reading test from stdin input: test error", err.Error())
	})

	t.Run("BreakOnDot", func(t *testing.T) {
		oldReader := InputReader
		defer func() { InputReader = oldReader }()

		InputReader = func() io.Reader {
			return bytes.NewReader([]byte("test text\nmore text\n.\n"))
		}

		text, _ := ReadInput("test")

		assert.Equal(t, "test text\nmore text", text)
	})

	t.Run("BreakOnEOF", func(t *testing.T) {
		oldReader := InputReader
		defer func() { InputReader = oldReader }()

		InputReader = func() io.Reader {
			return ioutil.NopCloser(bytes.NewReader([]byte("test text")))
		}

		text, _ := ReadInput("test")

		assert.Equal(t, "test text", text)
	})
}
