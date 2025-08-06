package output

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ans-group/sdk-go/pkg/connection"
)

var outputExit func(code int) = os.Exit
var errorLevel int

type OutputHandlerOption func(*OutputHandler)

func WithAdditionalColumns(columns ...string) OutputHandlerOption {
	return func(h *OutputHandler) {
		h.additionalColumns = append(h.additionalColumns, columns...)
	}
}

func SetOutputExit(e func(code int)) func(code int) {
	oldOutputExit := outputExit
	outputExit = e

	return oldOutputExit
}

type DebugLogger struct {
}

func (l *DebugLogger) Error(msg string) {
	Error(msg)
}

func (l *DebugLogger) Warn(msg string) {
	Error(msg)
}

func (l *DebugLogger) Info(msg string) {
	Error(msg)
}

func (l *DebugLogger) Debug(msg string) {
	Error(msg)
}

func (l *DebugLogger) Trace(msg string) {
	Error(msg)
}

// Error writes specified string to stderr
func Error(str string) {
	_, _ = os.Stderr.WriteString(str + "\n")
}

// Errorf writes specified string with formatting to stderr
func Errorf(format string, a ...interface{}) {
	Error(fmt.Sprintf(format, a...))
}

// Fatal writes specified string to stderr and calls outputExit to
// exit with 1
func Fatal(str string) {
	Error(str)
	outputExit(1)
}

// Fatalf writes specified string with formatting to stderr and calls
// outputExit to exit with 1
func Fatalf(format string, a ...interface{}) {
	Fatal(fmt.Sprintf(format, a...))
}

// OutputWithCustomErrorLevel is a wrapper for OutputError, which sets global
// var errorLevel with provided level
func OutputWithCustomErrorLevel(level int, str string) {
	Error(str)
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

func ExitWithErrorLevel() {
	outputExit(errorLevel)
}

func CommandOutputPaginated[T any](cmd *cobra.Command, d interface{}, paginated *connection.Paginated[T]) error {
	err := CommandOutput(cmd, d)
	if err != nil {
		return err
	}

	Errorf("page %d/%d", paginated.CurrentPage(), paginated.TotalPages())
	return nil
}

func CommandOutput(cmd *cobra.Command, d interface{}, opts ...OutputHandlerOption) error {
	return NewOutputHandler(opts...).Output(cmd, d)
}
