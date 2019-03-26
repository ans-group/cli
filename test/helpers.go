package test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/ukfast/cli/internal/pkg/output"
	"gopkg.in/go-playground/assert.v1"
)

func dieOn(err error, t *testing.T) {
	if err != nil {
		t.Fatal(err)
	}
}

// CatchStdErr returns output to `os.Stderr` from `f` as string
// https://groups.google.com/d/msg/golang-nuts/hVUtoeyNL7Y/HxEfVr70AwAJ
func CatchStdErr(t *testing.T, f func()) string {
	realStderr := os.Stderr
	defer func() { os.Stderr = realStderr }()
	r, fakeStderr, err := os.Pipe()
	dieOn(err, t)
	os.Stderr = fakeStderr
	f()
	dieOn(fakeStderr.Close(), t)
	newErrBytes, err := ioutil.ReadAll(r)
	dieOn(err, t)
	dieOn(r.Close(), t)

	return string(newErrBytes)
}

// CatchStdOut returns output to `os.Stdout` from `f` as string
// https://groups.google.com/d/msg/golang-nuts/hVUtoeyNL7Y/HxEfVr70AwAJ
func CatchStdOut(t *testing.T, f func()) string {
	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()
	r, fakeStdout, err := os.Pipe()
	dieOn(err, t)
	os.Stdout = fakeStdout
	f()
	// need to close here, otherwise ReadAll never gets "EOF".
	dieOn(fakeStdout.Close(), t)
	newOutBytes, err := ioutil.ReadAll(r)
	dieOn(err, t)
	dieOn(r.Close(), t)

	return string(newOutBytes)
}

// CatchStdOutStdErr returns output to `os.Stdout` and `os.Stderr` from `f` as strings
func CatchStdOutStdErr(t *testing.T, f func()) (stdOut string, stdErr string) {
	stdErr = CatchStdErr(t, func() {
		stdOut = CatchStdOut(t, f)
	})

	return stdOut, stdErr
}

// Output tests the output of func f, asserting stdout matches expected
func Output(t *testing.T, expected string, f func()) {
	output := CatchStdOut(t, f)
	assert.Equal(t, expected, output)
}

// ErrorOutput tests the output of func f, asserting stderr matches expected
func ErrorOutput(t *testing.T, expected string, f func()) {
	output := CatchStdErr(t, f)
	assert.Equal(t, expected, output)
}

// CombinedOutput tests the output of func f, asserting stdout/stderr matches
// expectedStdOut/expectedStdErr
func CombinedOutput(t *testing.T, expectedStdOut string, expectedStdErr string, f func()) {
	stdOut, stdErr := CatchStdOutStdErr(t, f)
	assert.Equal(t, expectedStdOut, stdOut)
	assert.Equal(t, expectedStdErr, stdErr)
}

// FatalOutput tests the output of func f, asserting stderr matches expected and
// exit code is equal to 1
func FatalOutput(t *testing.T, expected string, f func()) {
	code := 0
	oldOutputExit := output.SetOutputExit(func(c int) {
		code = c
	})
	defer func() { output.SetOutputExit(oldOutputExit) }()

	ErrorOutput(t, expected, f)
	assert.Equal(t, 1, code)
}
