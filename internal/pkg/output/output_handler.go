package output

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/iancoleman/strcase"
	"github.com/olekukonko/tablewriter"
	"github.com/olekukonko/tablewriter/renderer"
	"github.com/olekukonko/tablewriter/tw"
	"github.com/ryanuber/go-glob"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"k8s.io/client-go/util/jsonpath"
)

type DefaultColumnable interface {
	DefaultColumns() []string
}

type FieldValueHandlerFunc func(reflectedValue reflect.Value) string

type HandlesFields interface {
	FieldValueHandlers() map[string]FieldValueHandlerFunc
}

var MonetaryFieldValueHandler = func(reflectedValue reflect.Value) string {
	return fmt.Sprintf("%.2f", reflectedValue.Float())
}

type ProvidesFields interface {
	Fields() []*OrderedFields
}

type OutputHandlerOpts map[string]interface{}

type OutputHandler struct {
	additionalColumns []string
}

func NewOutputHandler(opts ...OutputHandlerOption) *OutputHandler {
	h := &OutputHandler{}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

func (o *OutputHandler) Output(cmd *cobra.Command, d interface{}) error {
	var flag string

	if cmd.Flags().Changed("output") {
		flag, _ = cmd.Flags().GetString("output")
	}

	if len(flag) == 0 {
		flag = "table"
		outputDefault := config.GetString("output.default")
		if len(outputDefault) > 0 {
			flag = outputDefault
		}
	}

	format, arg := ParseOutputFlag(flag)

	switch format {
	case "json":
		return o.JSON(d)
	case "list":
		return o.List(cmd, d)
	case "value":
		return o.Value(cmd, d)
	case "csv":
		return o.CSV(cmd, d)
	case "yaml":
		return o.YAML(d)
	case "jsonpath":
		return o.JSONPath(arg, d)
	case "template":
		return o.Template(arg, d)
	default:
		Errorf("invalid output format [%s], defaulting to 'table'", format)
		fallthrough
	case "table":
		return o.Table(cmd, d)
	}
}

func (o *OutputHandler) JSON(d interface{}) error {
	out, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal json: %s", err)
	}

	_, err = fmt.Print(string(out[:]))

	return err
}

// Table takes an array of mapped fields (key being lowercased name), and outputs a table
func (o *OutputHandler) Table(cmd *cobra.Command, d interface{}) error {
	columns, rows := o.getData(cmd, d)
	if len(rows) < 1 {
		return nil
	}

	// Get output style from configuration, default to "ascii"
	var symbols tw.Symbols
	outputStyle := config.GetString("output.table.style")
	switch outputStyle {
	case "unicode":
		symbols = tw.NewSymbols(tw.StyleDefault)
	default:
		symbols = tw.NewSymbols(tw.StyleASCII)
	}

	table := tablewriter.NewTable(os.Stdout,
		tablewriter.WithHeaderAlignment(tw.AlignCenter),
		tablewriter.WithRowAlignment(tw.AlignLeft),
		tablewriter.WithRenderer(renderer.NewBlueprint(tw.Rendition{Symbols: symbols})),
		tablewriter.WithConfig(tablewriter.Config{
			Row: tw.CellConfig{
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNormal},
				ColMaxWidths: tw.CellWidth{Global: 30},
			},
			Header: tw.CellConfig{
				Formatting:   tw.CellFormatting{AutoWrap: tw.WrapNormal},
				ColMaxWidths: tw.CellWidth{Global: 30},
			},
		}))

	table.Header(columns)
	_ = table.Bulk(rows)
	_ = table.Render()

	return nil
}

// List will format specified rows using given includeProperties by extracting fields,
// and output them to stdout
func (o *OutputHandler) List(cmd *cobra.Command, d interface{}) error {
	columns, rows := o.getData(cmd, d)
	if len(rows) < 1 {
		return nil
	}

	f := bufio.NewWriter(os.Stdout)
	defer func() { _ = f.Flush() }()

	maxPropertyLength := getMaxPropertyLength(columns)
	for i, row := range rows {
		if i > 0 {
			_, _ = f.WriteString("\n")
		}

		for columnIndex, column := range columns {
			formatted := formatListPropertyValue(column, row[columnIndex], maxPropertyLength)
			_, _ = fmt.Fprint(f, formatted)
		}

	}

	return nil
}

// Value will format specified rows using given includeProperties by extracting field values,
// and output them to stdout
func (o *OutputHandler) Value(cmd *cobra.Command, d interface{}) error {
	columns, rows := o.getData(cmd, d)
	if len(rows) < 1 {
		return nil
	}

	for _, row := range rows {
		var rowData []string
		for columnIndex := range columns {
			rowData = append(rowData, row[columnIndex])
		}
		fmt.Println(strings.Join(rowData, " "))
	}

	return nil
}

// YAML marshals and outputs value v to stdout
func (o *OutputHandler) YAML(d interface{}) error {
	out, err := yaml.Marshal(d)
	if err != nil {
		return fmt.Errorf("failed to marshal yaml: %s", err)
	}

	_, err = fmt.Print(string(out[:]))

	return err
}

// JSONPath marshals and outputs value v to stdout
func (o *OutputHandler) JSONPath(query string, d interface{}) error {
	j := jsonpath.New("clioutput")
	err := j.Parse(query)
	if err != nil {
		return fmt.Errorf("failed to parse jsonpath template: %w", err)
	}

	err = j.Execute(os.Stdout, d)
	if err != nil {
		return fmt.Errorf("failed to execute jsonpath: %w", err)
	}

	return nil
}

// CSV outputs provided rows as CSV to stdout
func (o *OutputHandler) CSV(cmd *cobra.Command, d interface{}) error {
	columns, rows := o.getData(cmd, d)
	if len(rows) < 1 {
		return nil
	}

	w := csv.NewWriter(os.Stdout)

	// First retrieve properties and write to CSV buffer
	err := w.Write(columns)
	if err != nil {
		return err
	}

	for _, row := range rows {
		// For each row, obtain property data and and write to CSV buffer
		var rowData []string
		for columnIndex := range columns {
			rowData = append(rowData, row[columnIndex])
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

// Template will format i with given Golang template t, and output resulting string
// to stdout
func (o *OutputHandler) Template(t string, d interface{}) error {
	tmpl, err := template.New("output").Parse(t)
	if err != nil {
		return fmt.Errorf("failed to create template: %s", err.Error())
	}

	switch reflect.TypeOf(d).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(d)
		for i := 0; i < s.Len(); i++ {
			err = tmpl.Execute(os.Stdout, s.Index(i))
			if err != nil {
				return fmt.Errorf("failed to execute template on slice: %s", err.Error())
			}
			fmt.Print("\n")
		}
	default:
		err = tmpl.Execute(os.Stdout, d)
		if err != nil {
			return fmt.Errorf("failed to execute template: %s", err.Error())
		}
		fmt.Print("\n")
	}

	return nil
}

func (o *OutputHandler) getData(cmd *cobra.Command, d interface{}) (filteredColumns []string, filteredRows [][]string) {
	rows := o.convert(d, reflect.ValueOf(d))
	if len(rows) == 0 {
		return
	}

	var filteredColumnNames []string
	if cmd.Flags().Changed("property") {
		filteredColumnNames, _ = cmd.Flags().GetStringSlice("property")
	} else if d, ok := d.(DefaultColumnable); ok && len(d.DefaultColumns()) > 0 {
		filteredColumnNames = d.DefaultColumns()
	}

	// Always add additional columns from options
	if len(o.additionalColumns) > 0 {
		filteredColumnNames = append(filteredColumnNames, o.additionalColumns...)
	}

	if len(filteredColumnNames) > 0 {
		for _, column := range rows[0].Keys() {
			for _, filteredColumnName := range filteredColumnNames {
				if glob.Glob(strings.ToLower(filteredColumnName), column) {
					filteredColumns = append(filteredColumns, column)
				}
			}
		}
	} else {
		filteredColumns = rows[0].Keys()
	}

	for _, row := range rows {
		filteredRow := make([]string, len(filteredColumns))
		for filteredColumnIndex, filteredColumn := range filteredColumns {
			filteredRow[filteredColumnIndex] = row.Get(filteredColumn)
		}
		filteredRows = append(filteredRows, filteredRow)
	}

	return filteredColumns, filteredRows
}

func (o *OutputHandler) convert(d interface{}, reflectedValue reflect.Value) []*OrderedFields {
	if dProvidesFields, ok := d.(ProvidesFields); ok {
		return dProvidesFields.Fields()
	}

	fields := []*OrderedFields{}

	switch reflectedValue.Kind() {
	case reflect.Slice:
		for i := 0; i < reflectedValue.Len(); i++ {
			fields = append(fields, o.convert(d, reflectedValue.Index(i))...)
		}
	case reflect.Struct:
		fields = append(fields, o.convertField(d, NewOrderedFields(), "", reflectedValue))
	}

	return fields
}

func (o *OutputHandler) convertField(d interface{}, v *OrderedFields, fieldName string, reflectedValue reflect.Value) *OrderedFields {
	if dHandlesFields, ok := d.(HandlesFields); ok {
		fieldHandlers := dHandlesFields.FieldValueHandlers()
		if fieldHandlers != nil && fieldHandlers[fieldName] != nil {
			v.Set(fieldName, fieldHandlers[fieldName](reflectedValue))
			return v
		}
	}

	switch reflectedValue.Kind() {
	case reflect.Struct:
		reflectedValueType := reflectedValue.Type()

		for i := 0; i < reflectedValueType.NumField(); i++ {
			reflectedValueField := reflectedValue.Field(i)
			reflectedValueTypeField := reflectedValueType.Field(i)

			if !reflectedValueField.CanInterface() {
				// Skip unexported field
				continue
			}
			childFieldName := ""
			if !reflectedValueTypeField.Anonymous {
				jsonTag := reflectedValueTypeField.Tag.Get("json")
				if jsonTag != "" {
					childFieldName = jsonTag
				} else {
					childFieldName = strcase.ToSnake(reflectedValueTypeField.Name)
				}
			}

			if len(fieldName) > 0 {
				childFieldName = fieldName + "." + childFieldName
			}

			o.convertField(d, v, childFieldName, reflectedValueField)
		}

		return v
	case reflect.String:
		v.Set(fieldName, reflectedValue.String())
		return v
	case reflect.Bool:
		v.Set(fieldName, strconv.FormatBool(reflectedValue.Bool()))
		return v
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		v.Set(fieldName, strconv.FormatInt(reflectedValue.Int(), 10))
		return v
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		v.Set(fieldName, strconv.FormatUint(reflectedValue.Uint(), 10))
		return v
	case reflect.Float32, reflect.Float64:
		v.Set(fieldName, fmt.Sprintf("%f", reflectedValue.Float()))
		return v
	case reflect.Ptr:
		if !reflectedValue.IsNil() {
			return o.convertField(d, v, fieldName, reflectedValue.Elem())
		}
	case reflect.Invalid:
		return nil
	}

	v.Set(fieldName, fmt.Sprintf("%v", reflectedValue.Interface()))
	return v
}

func formatListPropertyValue(property string, value string, maxPropertyLength int) string {
	paddedProperty := property + strings.Repeat(" ", maxPropertyLength-len(property))
	valueLines := strings.Split(value, "\n")

	var result strings.Builder
	for i, line := range valueLines {
		if i == 0 {
			result.WriteString(fmt.Sprintf("%s : %s\n", paddedProperty, line))
		} else {
			result.WriteString(fmt.Sprintf("%s   %s\n", strings.Repeat(" ", maxPropertyLength), line))
		}
	}
	return result.String()
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
