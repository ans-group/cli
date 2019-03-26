package test_output

import (
	"testing"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"gopkg.in/go-playground/assert.v1"
)

// AssertOutput tests the output of func f, asserting stdout matches expected
func AssertOutput(t *testing.T, expected string, f func()) {
	output := test.CatchStdOut(t, f)
	assert.Equal(t, expected, output)
}

// AssertErrorOutput tests the output of func f, asserting stderr matches expected
func AssertErrorOutput(t *testing.T, expected string, f func()) {
	output := test.CatchStdErr(t, f)
	assert.Equal(t, expected, output)
}

// AssertCombinedOutput tests the output of func f, asserting stdout/stderr matches
// expectedStdOut/expectedStdErr
func AssertCombinedOutput(t *testing.T, expectedStdOut string, expectedStdErr string, f func()) {
	stdOut, stdErr := test.CatchStdOutStdErr(t, f)
	assert.Equal(t, expectedStdOut, stdOut)
	assert.Equal(t, expectedStdErr, stdErr)
}

// AssertFatalOutput tests the output of func f, asserting stderr matches expected and
// exit code is equal to 1
func AssertFatalOutput(t *testing.T, expected string, f func()) {
	code := 0
	oldOutputExit := output.SetOutputExit(func(c int) {
		code = c
	})
	defer func() { output.SetOutputExit(oldOutputExit) }()

	AssertErrorOutput(t, expected, f)
	assert.Equal(t, 1, code)
}
