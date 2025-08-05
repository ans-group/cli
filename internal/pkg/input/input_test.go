package input

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/ans-group/cli/test/test_input"
	"github.com/stretchr/testify/assert"
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
		assert.Equal(t, "error reading test from stdin input: test error", err.Error())
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
			return io.NopCloser(bytes.NewReader([]byte("test text")))
		}

		text, _ := ReadInput("test")

		assert.Equal(t, "test text", text)
	})
}
