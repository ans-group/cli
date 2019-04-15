package test_output

import (
	"testing"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
	"gopkg.in/go-playground/assert.v1"
)

// AssertOutputFunc tests the output of func f, asserting using func assertFunc
func AssertOutputFunc(t *testing.T, f func(), assertFunc func(stdErr string)) {
	output := test.CatchStdOut(t, f)
	assertFunc(output)
}

// AssertOutput tests the output of func f, asserting stdout matches expected
func AssertOutput(t *testing.T, expected string, f func()) {
	AssertOutputFunc(t, f, func(stdOut string) {
		assert.Equal(t, expected, stdOut)
	})
}

// AssertErrorOutputFunc tests the output of func f, asserting using func assertFunc
func AssertErrorOutputFunc(t *testing.T, f func(), assertFunc func(stdErr string)) {
	output := test.CatchStdErr(t, f)
	assertFunc(output)
}

// AssertErrorOutput tests the output of func f, asserting stderr matches expected
func AssertErrorOutput(t *testing.T, expected string, f func()) {
	AssertErrorOutputFunc(t, f, func(stdErr string) {
		assert.Equal(t, expected, stdErr)
	})
}

// AssertCombinedOutputFunc tests the output of func f, asserting using func assertFunc
func AssertCombinedOutputFunc(t *testing.T, f func(), assertFunc func(stdOut, stdErr string)) {
	stdOut, stdErr := test.CatchStdOutStdErr(t, f)
	assertFunc(stdOut, stdErr)
}

// AssertCombinedOutput tests the output of func f, asserting stdout/stderr matches
// expectedStdOut/expectedStdErr
func AssertCombinedOutput(t *testing.T, expectedStdOut string, expectedStdErr string, f func()) {
	AssertCombinedOutputFunc(t, f, func(stdOut, stdErr string) {
		assert.Equal(t, expectedStdOut, stdOut)
		assert.Equal(t, expectedStdErr, stdErr)
	})
}

// AssertFatalOutputFunc tests the output of func f, asserting using func assertFunc and
// exit code is equal to 1
func AssertFatalOutputFunc(t *testing.T, f func(), assertFunc func(stdErr string)) {
	code := 0
	oldOutputExit := output.SetOutputExit(func(c int) {
		code = c
	})
	defer func() { output.SetOutputExit(oldOutputExit) }()

	AssertErrorOutputFunc(t, f, assertFunc)
	assert.Equal(t, 1, code)
}

// AssertFatalOutput tests the output of func f, asserting stderr matches expected and
// exit code is equal to 1
func AssertFatalOutput(t *testing.T, expected string, f func()) {
	AssertFatalOutputFunc(t, f, func(stdErr string) {
		assert.Equal(t, expected, stdErr)
	})
}
