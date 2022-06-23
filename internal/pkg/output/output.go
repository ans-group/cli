package output

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/util/jsonpath"

	"github.com/ans-group/sdk-go/pkg/connection"
)

var outputExit func(code int) = os.Exit
var errorLevel int

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
	os.Stderr.WriteString(str + "\n")
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

// Value will format specified rows using given includeProperties by extracting field values,
// and output them to stdout
func Value(rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	for _, row := range rows {
		var rowData []string
		for _, fieldKey := range row.Keys() {
			rowData = append(rowData, row.Get(fieldKey).Value)
		}
		fmt.Println(strings.Join(rowData, " "))
	}

	return nil
}

// JSON marshals and outputs value v to stdout
func JSON(v interface{}) error {
	out, err := json.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	_, err = fmt.Print(string(out[:]))

	return err
}

// YAML marshals and outputs value v to stdout
func YAML(v interface{}) error {
	out, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %s", err)
	}

	_, err = fmt.Print(string(out[:]))

	return err
}

// JSONPath marshals and outputs value v to stdout
func JSONPath(query string, v interface{}) error {
	j := jsonpath.New("clioutput")
	err := j.Parse(query)
	if err != nil {
		return fmt.Errorf("Failed to parse jsonpath template: %w", err)
	}

	err = j.Execute(os.Stdout, v)
	if err != nil {
		return fmt.Errorf("Failed to execute jsonpath: %w", err)
	}

	return nil
}

// Table takes an array of mapped fields (key being lowercased name), and outputs a table
func Table(rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)

	// properties will hold our header values, and will be used to determine required fields
	// when iterating over rows to add data to table
	headers := rows[0].Keys()

	table.SetHeader(headers)

	// Loop through each row, adding required fields specified in headers to table
	for _, row := range rows {
		var rowData []string
		for _, header := range headers {
			rowData = append(rowData, row.Get(header).Value)
		}
		table.Append(rowData)
	}

	table.Render()

	return nil
}

// CSV outputs provided rows as CSV to stdout
func CSV(rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	w := csv.NewWriter(os.Stdout)

	// First retrieve properties and write to CSV buffer
	headers := rows[0].Keys()
	err := w.Write(headers)
	if err != nil {
		return err
	}

	for _, row := range rows {
		// For each row, obtain property data and and write to CSV buffer
		var rowData []string
		for _, header := range headers {
			rowData = append(rowData, row.Get(header).Value)
		}
		err := w.Write(rowData)
		if err != nil {
			return err
		}

		// Finally flush CSV buffer to stdout
		w.Flush()
		err = w.Error()
		if err != nil {
			return err
		}
	}

	return nil
}

// List will format specified rows using given includeProperties by extracting fields,
// and output them to stdout
func List(rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()

	maxPropertyLength := getMaxPropertyLength(rows[0].Keys())
	for i, row := range rows {
		if i > 0 {
			f.WriteString("\n")
		}

		for _, fieldKey := range row.Keys() {
			f.WriteString(fmt.Sprintf("%s : %s\n", padProperty(fieldKey, maxPropertyLength), row.Get(fieldKey).Value))
		}
	}

	return nil
}

func padProperty(property string, maxLength int) string {
	diff := maxLength - len(property)
	if diff > 0 {
		return property + strings.Repeat(" ", diff)
	}
	return property
}

func getMaxPropertyLength(properties []string) int {
	maxLength := 0
	for _, property := range properties {
		length := len(property)
		if length > maxLength {
			maxLength = length
		}
	}
	return maxLength
}

// Template will format i with given Golang template t, and output resulting string
// to stdout
func Template(t string, i interface{}) error {
	tmpl, err := template.New("output").Parse(t)
	if err != nil {
		return fmt.Errorf("failed to create template: %s", err.Error())
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(i)
		for i := 0; i < s.Len(); i++ {
			err = tmpl.Execute(os.Stdout, s.Index(i))
			if err != nil {
				return fmt.Errorf("failed to execute template on slice: %s", err.Error())
			}
			fmt.Print("\n")
		}

		break
	default:
		err = tmpl.Execute(os.Stdout, i)
		if err != nil {
			return fmt.Errorf("failed to execute template: %s", err.Error())
		}
		fmt.Print("\n")
	}

	return nil
}

func CommandOutputPaginated[T any](cmd *cobra.Command, out OutputHandlerDataProvider, paginated *connection.Paginated[T]) error {
	err := CommandOutput(cmd, out)
	if err != nil {
		return err
	}

	Errorf("Page %d/%d", paginated.CurrentPage(), paginated.TotalPages())
	return nil
}

func CommandOutput(cmd *cobra.Command, out OutputHandlerDataProvider) error {
	// Format flag deprecated, however we'll check to see whether populated first and use it
	var flag string
	if cmd.Flags().Changed("format") {
		flag, _ = cmd.Flags().GetString("format")
	} else {
		flag, _ = cmd.Flags().GetString("output")
	}

	name, arg := ParseOutputFlag(flag)

	// outputtemplate flag deprecated, however we'll check to see whether populated first and use it
	if name == "template" && cmd.Flags().Changed("outputtemplate") {
		arg, _ = cmd.Flags().GetString("outputtemplate")
	}

	handler := NewOutputHandler(out, name, arg)
	handler.Properties, _ = cmd.Flags().GetStringSlice("property")

	return handler.Handle()
}
