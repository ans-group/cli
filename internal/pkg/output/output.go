package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/ryanuber/go-glob"
)

var outputExit func(code int) = os.Exit

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

// Value will format specified rows using given includeColumns by extracting field values,
// and output them to stdout
func Value(includeColumns []string, rows []*OrderedFields) error {
	for _, row := range rows {
		var out []string
		columns := getColumnsOrDefault(includeColumns, rows[0])
		for _, column := range columns {
			if row.Exists(strings.ToLower(column)) {
				out = append(out, row.Get(strings.ToLower(column)).Value)
			}
		}

		if len(out) > 0 {
			fmt.Println(strings.Join(out, " "))
		}
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

// Table takes an array of mapped fields (key being lowercased name), and outputs a table
// Included columns can be overriden by populating includeColumns parameter
func Table(includeColumns []string, rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)

	// columns will hold our header values, and will be used to determine required fields
	// when iterating over rows to add data to table
	columns := getColumnsOrDefault(includeColumns, rows[0])

	table.SetHeader(columns)

	// Loop through each row, adding required fields specified in columns to table
	for _, r := range rows {
		rowData := getColumnData(columns, r)
		table.Append(rowData)
	}

	table.Render()

	return nil
}

// CSV outputs provided rows as CSV to stdout
func CSV(includeColumns []string, rows []*OrderedFields) error {
	w := csv.NewWriter(os.Stdout)

	// First retrieve columns and write to CSV buffer
	columns := getColumnsOrDefault(includeColumns, rows[0])
	err := w.Write(columns)
	if err != nil {
		return err
	}

	for _, row := range rows {
		// For each row, obtain column data and and write to CSV buffer
		data := getColumnData(columns, row)
		err := w.Write(data)
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

func getColumnsOrDefault(includeColumns []string, fields *OrderedFields) []string {
	var columns []string

	if len(includeColumns) > 0 {
		// For each column in includeColumns, add column to columns array
		// if that column exists in fields
		for _, prop := range includeColumns {
			for _, column := range fields.Keys() {
				if glob.Glob(strings.ToLower(prop), column) {
					columns = append(columns, column)
				}
			}
		}

	} else {
		// Use default fields
		for _, column := range fields.Keys() {
			if fields.Get(column).Default {
				columns = append(columns, column)
			}
		}
	}

	return columns
}

func getColumnData(columns []string, fields *OrderedFields) []string {
	var columnData []string
	for _, column := range columns {
		columnData = append(columnData, fields.Get(column).Value)
	}

	return columnData
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

// OrderedFields holds a string map with field values, and a slice of keys for
// maintaining order
type OrderedFields struct {
	m    map[string]FieldValue
	keys []string
}

// NewOrderedFields returns a pointer to an initialized OrderedFields struct
func NewOrderedFields() *OrderedFields {
	return &OrderedFields{
		m: make(map[string]FieldValue),
	}
}

// Set adds/updates given key k with FieldValue v
func (o *OrderedFields) Set(k string, v FieldValue) {
	exists := o.Exists(k)
	o.m[k] = v
	if !exists {
		o.keys = append(o.keys, k)
	}
}

// Get retrieves FieldValue for given key k
func (o *OrderedFields) Get(k string) FieldValue {
	return o.m[k]
}

// Exists returns true if given key k exists, otherwise false
func (o *OrderedFields) Exists(k string) bool {
	_, exists := o.m[k]
	return exists
}

// Keys returns a list of ordered keys
func (o *OrderedFields) Keys() []string {
	return o.keys
}

// FieldValue holds the value for a table field
type FieldValue struct {
	Value   string
	Default bool
}

// NewFieldValue returns a new, initialized FieldValue struct
func NewFieldValue(value string, def bool) FieldValue {
	return FieldValue{
		Value:   value,
		Default: def,
	}
}
