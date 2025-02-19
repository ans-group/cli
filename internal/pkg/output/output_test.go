package output

import (
	"testing"

	"github.com/ans-group/cli/test"
	"github.com/stretchr/testify/assert"
)

func TestOutputExit_SetsHandler(t *testing.T) {
	code := 0
	oldOutputExit := SetOutputExit(func(c int) {
		code = c
	})
	defer func() { SetOutputExit(oldOutputExit) }()

	outputExit(1)

	assert.Equal(t, 1, code)
}

func TestOutputExit_ReturnsOldExitCode(t *testing.T) {
	code := 0
	oldOutputExit := SetOutputExit(func(c int) { code = c })
	defer func() { SetOutputExit(oldOutputExit) }()

	SetOutputExit(func(c int) {})(1)

	assert.Equal(t, 1, code)
}

func TestDebugLogger_Error_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Error("testerror")
	})

	assert.Equal(t, "testerror\n", output)
}

func TestDebugLogger_Warn_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Warn("testwarn")
	})

	assert.Equal(t, "testwarn\n", output)
}

func TestDebugLogger_Info_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Info("testinfo")
	})

	assert.Equal(t, "testinfo\n", output)
}

func TestDebugLogger_Debug_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Debug("testdebug")
	})

	assert.Equal(t, "testdebug\n", output)
}

func TestDebugLogger_Trace_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Trace("testtrace")
	})

	assert.Equal(t, "testtrace\n", output)
}

func TestError_ExpectedStderr(t *testing.T) {
	output := test.CatchStdErr(t, func() {
		Error("test")
	})

	assert.Equal(t, "test\n", output)
}

func TestErrorf_ExpectedStderr(t *testing.T) {
	output := test.CatchStdErr(t, func() {
		Errorf("test %s", "fmt")
	})

	assert.Equal(t, "test fmt\n", output)
}

func TestFatal_ExpectedExitCode(t *testing.T) {
	code := 0
	oldOutputExit := SetOutputExit(func(c int) {
		code = c
	})
	defer func() { outputExit = oldOutputExit }()

	output := test.CatchStdErr(t, func() {
		Fatal("test")
	})

	assert.Equal(t, "test\n", output)
	assert.Equal(t, 1, code)
}

func TestFatalf_ExpectedExitCode(t *testing.T) {
	code := 0
	oldOutputExit := SetOutputExit(func(c int) {
		code = c
	})
	defer func() { outputExit = oldOutputExit }()

	output := test.CatchStdErr(t, func() {
		Fatalf("test %s", "fmt")
	})

	assert.Equal(t, "test fmt\n", output)
	assert.Equal(t, 1, code)
}

func TestOutputWithCustomErrorLevel_ExpectedExitCode(t *testing.T) {
	output := test.CatchStdErr(t, func() {
		OutputWithCustomErrorLevel(5, "test")
	})

	assert.Equal(t, "test\n", output)
	assert.Equal(t, 5, errorLevel)
}

func TestOutputWithCustomErrorLevelf_ExpectedExitCode(t *testing.T) {
	output := test.CatchStdErr(t, func() {
		OutputWithCustomErrorLevelf(5, "test %s", "fmt")
	})

	assert.Equal(t, "test fmt\n", output)
	assert.Equal(t, 5, errorLevel)
}

func TestOutputWithErrorLevelf_ExpectedExitCode(t *testing.T) {
	output := test.CatchStdErr(t, func() {
		OutputWithErrorLevelf("test %s", "fmt")
	})

	assert.Equal(t, "test fmt\n", output)
	assert.Equal(t, 1, errorLevel)
}

func TestOutputWithErrorLevel_ExpectedExitCode(t *testing.T) {
	output := test.CatchStdErr(t, func() {
		OutputWithErrorLevel("test")
	})

	assert.Equal(t, "test\n", output)
	assert.Equal(t, 1, errorLevel)
}

func TestExitWithErrorLevel_ExpectedExitCode(t *testing.T) {
	errorLevel = 5
	code := 0
	oldOutputExit := SetOutputExit(func(c int) {
		code = c
	})
	defer func() { outputExit = oldOutputExit }()
	ExitWithErrorLevel()

	assert.Equal(t, 5, code)
}
