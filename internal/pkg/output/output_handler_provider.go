package output

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

type OutputHandlerProvider interface {
	GetData() interface{}
	GetFieldData() ([]*OrderedFields, error)
	SupportedFormats() []string
}

type ProviderOption func(p *GenericOutputHandlerProvider)

type GenericOutputHandlerProvider struct {
	data             interface{}
	fieldDataFunc    func() ([]*OrderedFields, error)
	supportedFormats []string
}

func NewGenericOutputHandlerProvider(opts ...ProviderOption) *GenericOutputHandlerProvider {
	p := &GenericOutputHandlerProvider{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func WithData(data interface{}) ProviderOption {
	return func(p *GenericOutputHandlerProvider) {
		p.data = data
	}
}

func WithSupportedFormats(supportedFormats []string) ProviderOption {
	return func(p *GenericOutputHandlerProvider) {
		p.supportedFormats = supportedFormats
	}
}

func WithFieldDataFunc(fieldDataFunc func() ([]*OrderedFields, error)) ProviderOption {
	return func(p *GenericOutputHandlerProvider) {
		p.fieldDataFunc = fieldDataFunc
	}
}

func (p *GenericOutputHandlerProvider) GetData() interface{} {
	return p.data
}

func (p *GenericOutputHandlerProvider) GetFieldData() ([]*OrderedFields, error) {
	return p.fieldDataFunc()
}

func (p *GenericOutputHandlerProvider) SupportedFormats() []string {
	return p.supportedFormats
}

type SerializedOutputHandlerProvider struct {
	*GenericOutputHandlerProvider
	DefaultFields []string
	IgnoredFields []string
}

func NewSerializedOutputHandlerProvider(items interface{}, defaultFields []string, ignoredFields []string) *SerializedOutputHandlerProvider {
	return &SerializedOutputHandlerProvider{
		GenericOutputHandlerProvider: NewGenericOutputHandlerProvider(
			WithData(items),
		),
		DefaultFields: defaultFields,
		IgnoredFields: ignoredFields,
	}
}

func (o *SerializedOutputHandlerProvider) GetFieldData() ([]*OrderedFields, error) {
	return o.convert(reflect.ValueOf(o.GetData())), nil
}

func (o *SerializedOutputHandlerProvider) convert(reflectedValue reflect.Value) []*OrderedFields {
	fields := []*OrderedFields{}

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

func (o *SerializedOutputHandlerProvider) convertStruct(reflectedValue reflect.Value) *OrderedFields {
	fields := NewOrderedFields()
	reflectedValueType := reflectedValue.Type()

	for i := 0; i < reflectedValueType.NumField(); i++ {
		reflectedValueField := reflectedValue.Field(i)
		reflectedValueTypeField := reflectedValueType.Field(i)

		if !reflectedValueField.CanInterface() {
			// Skip unexported field
			continue
		}
		fieldName := ""

		jsonTag := reflectedValueTypeField.Tag.Get("json")
		if jsonTag != "" {
			fieldName = jsonTag
		} else {
			fieldName = strcase.ToSnake(reflectedValueTypeField.Name)
		}

		if !o.isIgnoredField(fieldName) {
			fields.Set(fieldName, NewFieldValue(o.fieldToString(reflectedValueField), o.isDefaultField(fieldName)))
		}
	}

	return fields
}

func (o *SerializedOutputHandlerProvider) fieldToString(reflectedValue reflect.Value) string {
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

func (o *SerializedOutputHandlerProvider) isDefaultField(name string) bool {
	return o.fieldInFields(name, o.DefaultFields)
}

func (o *SerializedOutputHandlerProvider) isIgnoredField(name string) bool {
	return o.fieldInFields(name, o.IgnoredFields)
}

func (o *SerializedOutputHandlerProvider) fieldInFields(name string, fields []string) bool {
	for _, field := range fields {
		if strings.ToLower(field) == strings.ToLower(name) {
			return true
		}
	}

	return false
}
