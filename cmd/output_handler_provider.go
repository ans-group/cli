package cmd

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/ukfast/cli/internal/pkg/output"
)

type OutputHandlerProvider interface {
	GetData() interface{}
	GetFieldData() ([]*output.OrderedFields, error)
}

type GenericOutputHandlerProvider struct {
	items         interface{}
	DefaultFields []string
}

func NewGenericOutputHandlerProvider(items interface{}, defaultFields []string) *GenericOutputHandlerProvider {
	return &GenericOutputHandlerProvider{items: items, DefaultFields: defaultFields}
}

func (o *GenericOutputHandlerProvider) GetData() interface{} {
	return o.items
}

func (o *GenericOutputHandlerProvider) GetFieldData() ([]*output.OrderedFields, error) {
	return o.convert(reflect.ValueOf(o.items)), nil
}

func (o *GenericOutputHandlerProvider) convert(reflectedValue reflect.Value) []*output.OrderedFields {
	fields := []*output.OrderedFields{}

	switch reflectedValue.Kind() {
	case reflect.Slice:
		for i := 0; i < reflectedValue.Len(); i++ {
			fields = append(fields, o.convert(reflectedValue.Index(i))...)
		}
	case reflect.Struct:
		fields = append(fields, o.convertStruct(reflectedValue))
	}

	return fields
}

func (o *GenericOutputHandlerProvider) convertStruct(reflectedValue reflect.Value) *output.OrderedFields {
	fields := output.NewOrderedFields()
	reflectedValueType := reflectedValue.Type()

	for i := 0; i < reflectedValueType.NumField(); i++ {
		reflectedValueField := reflectedValue.Field(i)
		reflectedValueTypeField := reflectedValueType.Field(i)

		if !reflectedValueField.CanInterface() {
			// Skip unexported field
			continue
		}

		fieldName := strcase.ToSnake(reflectedValueTypeField.Name)

		fields.Set(fieldName, output.NewFieldValue(o.fieldToString(reflectedValueField), o.isDefaultField(fieldName)))
	}

	return fields
}

func (o *GenericOutputHandlerProvider) fieldToString(reflectedValue reflect.Value) string {
	switch reflectedValue.Kind() {
	case reflect.Slice:
		var items []string
		for i := 0; i < reflectedValue.Len(); i++ {
			items = append(items, o.fieldToString(reflectedValue.Index(i)))
		}

		return strings.Join(items, ", ")
	case reflect.String:
		return reflectedValue.String()
	case reflect.Bool:
		return strconv.FormatBool(reflectedValue.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(reflectedValue.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(reflectedValue.Uint(), 10)
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%f", reflectedValue.Float())
	}

	return fmt.Sprintf("%v", reflectedValue.Interface())
}

func (o *GenericOutputHandlerProvider) isDefaultField(name string) bool {
	for _, field := range o.DefaultFields {
		if strings.ToLower(field) == strings.ToLower(name) {
			return true
		}
	}

	return false
}
