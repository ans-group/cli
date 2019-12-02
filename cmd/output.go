package cmd

import (
	"fmt"
	"os"

	"github.com/ukfast/cli/internal/pkg/output"
)

var errorLevel int

// OutputWithCustomErrorLevel is a wrapper for OutputError, which sets global
// var errorLevel with provided level
func OutputWithCustomErrorLevel(level int, str string) {
	output.Error(str)
	errorLevel = level
}

// OutputWithCustomErrorLevelf is a wrapper for OutputWithCustomErrorLevel, which sets global
// var errorLevel with provided level
func OutputWithCustomErrorLevelf(level int, format string, a ...interface{}) {
	OutputWithCustomErrorLevel(level, fmt.Sprintf(format, a...))
}

// OutputWithErrorLevelf is a wrapper for OutputWithCustomErrorLevelf, which sets global
// var errorLevel to 1
func OutputWithErrorLevelf(format string, a ...interface{}) {
	OutputWithCustomErrorLevelf(1, format, a...)
}

// OutputWithErrorLevel is a wrapper for OutputWithCustomErrorLevel, which sets global
// var errorLevel to 1
func OutputWithErrorLevel(str string) {
	OutputWithCustomErrorLevel(1, str)
}

func Exit() {
	os.Exit(errorLevel)
}
