package output

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"

	"github.com/iancoleman/strcase"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

type DefaultColumnable interface {
	DefaultColumns() []string
}

type Sortable interface {
	SortableColumns() []string
	DefaultSortColumn() string
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

type OutputHandler struct{}

func NewOutputHandler() *OutputHandler {
	return &OutputHandler{}
}

// func Output(cmd *cobra.Command, d interface{}) error {
// 	return NewOutputHandler().Output(cmd, d)
// }

// // Handle calls the relevant OutputProvider data retrieval methods for given value
// // in struct property 'Format'
// func (o *OutputHandler) Handle() error {
// 	if !o.supportedFormat() {
// 		return fmt.Errorf("Unsupported output format [%s], supported formats: %s", o.Format, strings.Join(o.SupportedFormats, ", "))
// 	}

// 	switch o.Format {
// 	case "json":
// 		return JSON(o.DataProvider.GetData())
// 	case "yaml":
// 		return YAML(o.DataProvider.GetData())
// 	case "jsonpath":
// 		return JSONPath(o.FormatArg, o.DataProvider.GetData())
// 	case "template":
// 		return Template(o.FormatArg, o.DataProvider.GetData())
// 	case "value":
// 		d, err := o.getProcessedFieldData()
// 		if err != nil {
// 			return err
// 		}
// 		return Value(d)
// 	case "csv":
// 		d, err := o.getProcessedFieldData()
// 		if err != nil {
// 			return err
// 		}
// 		return CSV(d)
// }

func (o *OutputHandler) Output(cmd *cobra.Command, d interface{}) error {
	format, _ := cmd.Flags().GetString("output")
	switch format {
	case "json":
		return o.JSON(d)
	default:
		Errorf("Invalid output format [%s], defaulting to 'table'", format)
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

	table := tablewriter.NewWriter(os.Stdout)

	table.SetHeader(columns)
	table.AppendBulk(rows)
	table.Render()

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
	defer f.Flush()

	maxPropertyLength := getMaxPropertyLength(columns)
	for i, row := range rows {
		if i > 0 {
			f.WriteString("\n")
		}

		for columnIndex, column := range columns {
			f.WriteString(fmt.Sprintf("%s : %s\n", padProperty(column, maxPropertyLength), row[columnIndex]))
		}
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

	for _, column := range rows[0].Keys() {
		if len(filteredColumnNames) > 0 {
			for _, filteredColumnName := range filteredColumnNames {
				if column == filteredColumnName {
					filteredColumns = append(filteredColumns, filteredColumnName)
				}
			}
		} else {
			filteredColumns = append(filteredColumns, column)
		}
	}

	for _, row := range rows {
		filteredRow := make([]string, len(filteredColumns))
		for filteredColumnIndex, filteredColumn := range filteredColumns {
			filteredRow[filteredColumnIndex] = row.Get(filteredColumn)
		}
		filteredRows = append(filteredRows, filteredRow)
	}

	if d, ok := d.(Sortable); ok && len(d.SortableColumns()) > 0 {
		var sortColumn string
		if cmd.Flags().Changed("sort") {
			sortColumn, _ = cmd.Flags().GetString("sort")
		} else {
			sortColumn = d.DefaultSortColumn()
		}
		sortIndex := indexOf(filteredColumns, sortColumn)
		if sortIndex != -1 {
			sort.Slice(filteredRows, func(i, j int) bool {
				return filteredRows[i][sortIndex] < filteredRows[j][sortIndex]
			})
		}
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

func indexOf(slice []string, s string) int {
	for i, v := range slice {
		if v == s {
			return i
		}
	}
	return -1
}
