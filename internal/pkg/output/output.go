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

// Value will format specified rows using given includeProperties by extracting field values,
// and output them to stdout
func Value(includeProperties []string, rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	properties := getPropertiesOrDefault(includeProperties, rows[0])
	for _, row := range rows {
		data := getPropertyData(properties, row)
		fmt.Println(strings.Join(data, " "))
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
// Included properties can be overriden by populating includeProperties parameter
func Table(includeProperties []string, rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	table := tablewriter.NewWriter(os.Stdout)

	// properties will hold our header values, and will be used to determine required fields
	// when iterating over rows to add data to table
	headers := getPropertiesOrDefault(includeProperties, rows[0])

	table.SetHeader(headers)

	// Loop through each row, adding required fields specified in headers to table
	for _, r := range rows {
		rowData := getPropertyData(headers, r)
		table.Append(rowData)
	}

	table.Render()

	return nil
}

// CSV outputs provided rows as CSV to stdout
func CSV(includeProperties []string, rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	w := csv.NewWriter(os.Stdout)

	// First retrieve properties and write to CSV buffer
	properties := getPropertiesOrDefault(includeProperties, rows[0])
	err := w.Write(properties)
	if err != nil {
		return err
	}

	for _, row := range rows {
		// For each row, obtain property data and and write to CSV buffer
		data := getPropertyData(properties, row)
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

// List will format specified rows using given includeProperties by extracting fields,
// and output them to stdout
func List(includeProperties []string, rows []*OrderedFields) error {
	if len(rows) < 1 {
		return nil
	}

	f := bufio.NewWriter(os.Stdout)
	defer f.Flush()

	properties := getPropertiesOrDefault(includeProperties, rows[0])
	maxPropertyLength := getMaxPropertyLength(properties)
	for i, row := range rows {
		if i > 0 {
			f.WriteString("\n")
		}

		for _, property := range properties {
			if row.Exists(strings.ToLower(property)) {
				f.WriteString(fmt.Sprintf("%s : %s\n", padProperty(property, maxPropertyLength), row.Get(strings.ToLower(property)).Value))
			}
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

func getPropertiesOrDefault(includeProperties []string, fields *OrderedFields) []string {
	var properties []string

	if len(includeProperties) > 0 {
		// For each property in includeProperties, add property to properties array
		// if that property exists in fields
		for _, prop := range includeProperties {
			for _, property := range fields.Keys() {
				if glob.Glob(strings.ToLower(prop), property) {
					properties = append(properties, property)
				}
			}
		}

	} else {
		// Use default fields
		for _, property := range fields.Keys() {
			if fields.Get(property).Default {
				properties = append(properties, property)
			}
		}
	}

	return properties
}

func getPropertyData(properties []string, fields *OrderedFields) []string {
	var propertyData []string
	for _, property := range properties {
		propertyData = append(propertyData, fields.Get(property).Value)
	}

	return propertyData
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
