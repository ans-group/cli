package test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/spf13/viper"
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

func TestResetViper() {
	viper.SetDefault("command_wait_timeout_seconds", 1200)
	viper.SetDefault("command_wait_sleep_seconds", 5)
}
