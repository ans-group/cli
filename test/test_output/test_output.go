package test_output

import (
	"testing"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/cli/test"
	"gopkg.in/go-playground/assert.v1"
)

// AssertOutputFunc tests the output of func f, asserting using func assertFunc
func AssertOutputFunc(t *testing.T, assertFunc func(stdErr string), f func()) {
	output := test.CatchStdOut(t, f)
	assertFunc(output)
}

// AssertOutput tests the output of func f, asserting stdout matches expected
func AssertOutput(t *testing.T, expected string, f func()) {
	AssertOutputFunc(t, func(stdOut string) {
		assert.Equal(t, expected, stdOut)
	}, f)
}

// AssertErrorOutputFunc tests the output of func f, asserting using func assertFunc
func AssertErrorOutputFunc(t *testing.T, assertFunc func(stdErr string), f func()) {
	output := test.CatchStdErr(t, f)
	assertFunc(output)
}

// AssertErrorOutput tests the output of func f, asserting stderr matches expected
func AssertErrorOutput(t *testing.T, expected string, f func()) {
	AssertErrorOutputFunc(t, func(stdErr string) {
		assert.Equal(t, expected, stdErr)
	}, f)
}

// AssertCombinedOutputFunc tests the output of func f, asserting using func assertFunc
func AssertCombinedOutputFunc(t *testing.T, assertFunc func(stdOut, stdErr string), f func()) {
	stdOut, stdErr := test.CatchStdOutStdErr(t, f)
	assertFunc(stdOut, stdErr)
}

// AssertCombinedOutput tests the output of func f, asserting stdout/stderr matches
// expectedStdOut/expectedStdErr
func AssertCombinedOutput(t *testing.T, expectedStdOut string, expectedStdErr string, f func()) {
	AssertCombinedOutputFunc(t, func(stdOut, stdErr string) {
		assert.Equal(t, expectedStdOut, stdOut)
		assert.Equal(t, expectedStdErr, stdErr)
	}, f)
}

// AssertErrorLevelFunc tests the output of func f, asserting using func assertFunc and
// exit code is equal to level
func AssertErrorLevelFunc(t *testing.T, level int, assertFunc func(stdErr string), f func()) {
	code := 0
	oldOutputExit := output.SetOutputExit(func(c int) {
		code = c
	})
	defer func() { output.SetOutputExit(oldOutputExit) }()

	AssertErrorOutputFunc(t, assertFunc, f)
	assert.Equal(t, level, code)
}

// AssertErrorLevelOutput tests the output of func f, asserting stderr matches expected and
// exit code is equal to 1
func AssertErrorLevelOutput(t *testing.T, level int, expected string, f func()) {
	AssertErrorLevelFunc(t, level, func(stdErr string) {
		assert.Equal(t, expected, stdErr)
	}, f)
}

// AssertErrorLevel tests that exit code is equal to 1
func AssertErrorLevel(t *testing.T, f func()) {
	AssertErrorLevelFunc(t, 1, func(stdErr string) {}, f)
}
