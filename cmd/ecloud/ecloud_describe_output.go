package ecloud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// describeMaxDerivedColumns is the maximum number of columns derived for a child
// collection which declares no default columns of its own.
const describeMaxDerivedColumns = 5

// DescribeResult is the assembled description of a single resource. It is the
// structure marshalled directly for JSON and YAML output.
type DescribeResult struct {
	Type     string             `json:"type"`
	ID       string             `json:"id"`
	Resource any                `json:"resource"`
	Related  []DescribeRelated  `json:"related,omitempty"`
	Children []DescribeChildren `json:"children,omitempty"`
}

// DescribeRelated is a resolved parent reference.
type DescribeRelated struct {
	Field string `json:"field"`
	Type  string `json:"type"`
	ID    string `json:"id"`
	Name  string `json:"name,omitempty"`
	Error string `json:"error,omitempty"`
}

// DescribeChildren is one resolved child collection.
type DescribeChildren struct {
	Title string `json:"title"`
	Items any    `json:"items"`
	// Total is the number of items available, which exceeds the number of items
	// retrieved where the collection has been truncated by --child-limit.
	Total     int    `json:"total"`
	Truncated bool   `json:"truncated"`
	Error     string `json:"error,omitempty"`
}

// describeField is a single flattened field of a resource.
type describeField struct {
	// Name is the dotted json path of the field, e.g. "sync.status".
	Name string
	// Value is the rendered value of the field.
	Value string
	// Boolean indicates that the field holds a boolean value, which remains
	// meaningful when zero-valued.
	Boolean bool
	// Zero indicates that the field holds the zero value for its type.
	Zero bool
}

// ecloudDescribeOutput renders describe results in the format given by the --output
// flag. The shared output handler flattens data into a single level of string fields,
// so cannot represent the nested documents produced by describe, which therefore
// handles its own rendering.
func ecloudDescribeOutput(cmd *cobra.Command, results []DescribeResult) error {
	format, _ := output.ParseOutputFlag(describeOutputFlag(cmd))

	switch format {
	case "json":
		return describeMarshal(cmd.OutOrStdout(), results, false)
	case "json-pretty":
		return describeMarshal(cmd.OutOrStdout(), results, true)
	case "yaml":
		out, err := yaml.Marshal(results)
		if err != nil {
			return fmt.Errorf("ecloud: failed to marshal yaml: %w", err)
		}

		_, err = cmd.OutOrStdout().Write(out)
		return err
	case "table":
	default:
		output.Errorf("describe does not support output format [%s], using default", format)
	}

	describeRender(cmd.OutOrStdout(), results)

	return nil
}

// describeOutputFlag returns the configured output format, following the same
// precedence as the shared output handler.
func describeOutputFlag(cmd *cobra.Command) string {
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

	return flag
}

func describeMarshal(w io.Writer, results []DescribeResult, pretty bool) error {
	var out []byte
	var err error
	if pretty {
		out, err = json.MarshalIndent(results, "", "  ")
	} else {
		out, err = json.Marshal(results)
	}
	if err != nil {
		return fmt.Errorf("ecloud: failed to marshal json: %w", err)
	}

	_, err = fmt.Fprintln(w, string(out))

	return err
}

// describeRender writes the human-readable representation of describe results.
func describeRender(w io.Writer, results []DescribeResult) {
	for i, result := range results {
		if i > 0 {
			_, _ = fmt.Fprintln(w)
		}

		describeRenderResult(w, result)
	}
}

func describeRenderResult(w io.Writer, result DescribeResult) {
	_, _ = fmt.Fprintf(w, "%s  %s\n", result.Type, result.ID)

	// Parent references are displayed in the related section, so are omitted from
	// the resource block. Any reference which was not resolved into that section,
	// such as a non-string identifier, remains here.
	relatedFields := make([]string, len(result.Related))
	for i, reference := range result.Related {
		relatedFields[i] = reference.Field
	}

	describeRenderFields(w, describeResourceFields(result.Resource, relatedFields))

	if len(result.Related) > 0 {
		_, _ = fmt.Fprintf(w, "\nRelated\n")
		describeRenderRelated(w, result.Related)
	}

	for _, children := range result.Children {
		_, _ = fmt.Fprintln(w)
		describeRenderChildren(w, children)
	}
}

// describeResourceFields returns the fields of a resource for display, with the sync
// and task sub-objects flattened, and zero-valued fields omitted. Booleans are
// retained when false, as they remain meaningful. Fields named in relatedFields are
// omitted, as they are displayed in the related section instead.
func describeResourceFields(resource any, relatedFields []string) []describeField {
	var fields []describeField
	for _, field := range describeFlattenFields(resource) {
		switch {
		case field.Name == "id":
			continue
		case field.Name == "sync.status":
			field.Name = "sync"
		case field.Name == "task.in_progress":
			field.Name = "task_in_progress"
		case strings.HasPrefix(field.Name, "sync."), strings.HasPrefix(field.Name, "task."):
			continue
		case slices.Contains(relatedFields, field.Name):
			continue
		}

		if field.Zero && !field.Boolean {
			continue
		}

		fields = append(fields, field)
	}

	return fields
}

func describeRenderFields(w io.Writer, fields []describeField) {
	buf, tw := describeBlockWriter()
	for _, field := range fields {
		_, _ = fmt.Fprintf(tw, "  %s\t%s\n", field.Name, field.Value)
	}
	describeFlushBlock(w, buf, tw)
}

func describeRenderRelated(w io.Writer, related []DescribeRelated) {
	buf, tw := describeBlockWriter()
	for _, reference := range related {
		name := reference.Name
		switch {
		case reference.Error != "":
			name = fmt.Sprintf("(%s)", reference.Error)
		case name != "":
			name = fmt.Sprintf("(%s)", name)
		}

		_, _ = fmt.Fprintf(tw, "  %s\t%s\t%s\n", strings.TrimSuffix(reference.Field, "_id"), reference.ID, name)
	}
	describeFlushBlock(w, buf, tw)
}

// describeBlockWriter returns a buffered tabwriter for rendering an aligned block.
func describeBlockWriter() (*bytes.Buffer, *tabwriter.Writer) {
	buf := &bytes.Buffer{}

	return buf, tabwriter.NewWriter(buf, 0, 8, 2, ' ', 0)
}

// describeFlushBlock writes an aligned block to w, stripping the trailing padding
// left by the tabwriter on rows with empty trailing cells.
func describeFlushBlock(w io.Writer, buf *bytes.Buffer, tw *tabwriter.Writer) {
	_ = tw.Flush()
	if buf.Len() < 1 {
		return
	}

	for line := range strings.SplitSeq(strings.TrimSuffix(buf.String(), "\n"), "\n") {
		_, _ = fmt.Fprintln(w, strings.TrimRight(line, " "))
	}
}

func describeRenderChildren(w io.Writer, children DescribeChildren) {
	if children.Error != "" {
		_, _ = fmt.Fprintf(w, "%s\n  %s\n", children.Title, children.Error)
		return
	}

	items := describeItems(children.Items)
	count := strconv.Itoa(len(items))
	if children.Truncated {
		count = fmt.Sprintf("%d of %d, truncated", len(items), children.Total)
	}

	_, _ = fmt.Fprintf(w, "%s (%s)\n", children.Title, count)
	if len(items) < 1 {
		return
	}

	columns := describeChildColumns(children.Items, items[0])

	buf, tw := describeBlockWriter()
	_, _ = fmt.Fprintf(tw, "  %s\n", strings.ToUpper(strings.Join(columns, "\t")))
	for _, item := range items {
		fields := describeFlattenFields(item.Interface())

		values := make([]string, len(columns))
		for i, column := range columns {
			values[i] = describeFieldValue(fields, column)
		}

		_, _ = fmt.Fprintf(tw, "  %s\n", strings.Join(values, "\t"))
	}
	describeFlushBlock(w, buf, tw)
}

// describeItems returns the individual items of a child collection.
func describeItems(collection any) []reflect.Value {
	reflectedValue := reflect.ValueOf(collection)
	if reflectedValue.Kind() != reflect.Slice {
		return nil
	}

	items := make([]reflect.Value, reflectedValue.Len())
	for i := range items {
		items[i] = reflectedValue.Index(i)
	}

	return items
}

// describeChildColumns returns the columns to display for a child collection,
// preferring the default columns declared by the collection type where available.
func describeChildColumns(collection any, item reflect.Value) []string {
	if columnable, ok := collection.(output.DefaultColumnable); ok {
		if columns := columnable.DefaultColumns(); len(columns) > 0 {
			return columns
		}
	}

	var columns []string
	for _, field := range describeFlattenFields(item.Interface()) {
		if len(columns) >= describeMaxDerivedColumns {
			break
		}

		columns = append(columns, field.Name)
	}

	return columns
}

// describeFieldValue returns the value of the named field. Nested fields may be
// addressed either by their dotted path (e.g. "sync.status") or by the equivalent
// underscored name (e.g. "sync_status"), as used by the default column definitions.
func describeFieldValue(fields []describeField, name string) string {
	for _, field := range fields {
		if field.Name == name || strings.ReplaceAll(field.Name, ".", "_") == name {
			return field.Value
		}
	}

	return ""
}

// describeFlattenFields returns the json-tagged fields of a struct in declaration
// order, recursing into nested structs to produce dotted field names.
func describeFlattenFields(resource any) []describeField {
	var fields []describeField
	describeFlattenValue("", reflect.ValueOf(resource), &fields)

	return fields
}

func describeFlattenValue(name string, value reflect.Value, fields *[]describeField) {
	if value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return
		}

		value = value.Elem()
	}

	// Structs are recursed into rather than emitted as a field of their own, which
	// includes the unnamed top-level value being flattened.
	if value.Kind() == reflect.Struct {
		reflectedType := value.Type()
		for i := range reflectedType.NumField() {
			field := reflectedType.Field(i)
			if !value.Field(i).CanInterface() {
				continue
			}

			childName := describeJSONFieldName(field)
			if name != "" {
				childName = name + "." + childName
			}

			describeFlattenValue(childName, value.Field(i), fields)
		}

		return
	}

	if name == "" {
		return
	}

	formatted, ok := describeFormatValue(value)
	if !ok {
		return
	}

	*fields = append(*fields, describeField{
		Name:    name,
		Value:   formatted,
		Boolean: value.Kind() == reflect.Bool,
		Zero:    value.IsZero(),
	})
}

// describeFormatValue renders a scalar field value, reporting whether the value is
// of a displayable kind.
func describeFormatValue(value reflect.Value) (string, bool) {
	switch value.Kind() {
	case reflect.String:
		return value.String(), true
	case reflect.Bool:
		return strconv.FormatBool(value.Bool()), true
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(value.Int(), 10), true
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(value.Uint(), 10), true
	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(value.Float(), 'f', -1, 64), true
	case reflect.Slice:
		// Resource tags are the only collection held inline by an eCloud resource.
		// They are rendered as 'scope:name' to match the tags column displayed by
		// the list and show commands.
		if tags, ok := value.Interface().([]ecloud.ResourceTag); ok {
			names := make([]string, len(tags))
			for i, tag := range tags {
				names[i] = tag.Scope + ":" + tag.Name
			}

			return strings.Join(names, ", "), true
		}
	}

	return "", false
}
