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
	DefaultFields     []string
	IgnoredFields     []string
	MonetaryFields    []string
	FieldHandlerFuncs map[string]FieldHandlerFunc
}

type FieldHandlerFunc func(v *OrderedFields, fieldName string, reflectedValue reflect.Value) *OrderedFields

func NewSerializedOutputHandlerProvider(items interface{}) *SerializedOutputHandlerProvider {
	return &SerializedOutputHandlerProvider{
		GenericOutputHandlerProvider: NewGenericOutputHandlerProvider(
			WithData(items),
		),
	}
}

func (o *SerializedOutputHandlerProvider) WithDefaultFields(fields []string) *SerializedOutputHandlerProvider {
	o.DefaultFields = append(o.DefaultFields, fields...)
	return o
}

func (o *SerializedOutputHandlerProvider) WithIgnoredFields(fields []string) *SerializedOutputHandlerProvider {
	o.IgnoredFields = append(o.IgnoredFields, fields...)
	return o
}

func (o *SerializedOutputHandlerProvider) WithMonetaryFields(fields []string) *SerializedOutputHandlerProvider {
	o.MonetaryFields = append(o.MonetaryFields, fields...)
	return o
}

func (o *SerializedOutputHandlerProvider) WithFieldHandler(fieldName string, f FieldHandlerFunc) *SerializedOutputHandlerProvider {
	if o.FieldHandlerFuncs == nil {
		o.FieldHandlerFuncs = make(map[string]FieldHandlerFunc)
	}
	o.FieldHandlerFuncs[fieldName] = f
	return o
}

func (o *SerializedOutputHandlerProvider) WithMultipleFieldHandler(fieldNames []string, f FieldHandlerFunc) *SerializedOutputHandlerProvider {
	for _, fieldName := range fieldNames {
		o.WithFieldHandler(fieldName, f)
	}
	return o
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
		fields = append(fields, o.convertField(NewOrderedFields(), "", reflectedValue))
	}

	return fields
}

func (o *SerializedOutputHandlerProvider) convertField(v *OrderedFields, fieldName string, reflectedValue reflect.Value) *OrderedFields {
	if o.FieldHandlerFuncs[fieldName] != nil {
		return o.FieldHandlerFuncs[fieldName](v, fieldName, reflectedValue)
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
			jsonTag := reflectedValueTypeField.Tag.Get("json")
			if jsonTag != "" {
				childFieldName = jsonTag
			} else {
				childFieldName = strcase.ToSnake(reflectedValueTypeField.Name)
			}

			if len(fieldName) > 0 {
				childFieldName = fieldName + "_" + childFieldName
			}

			o.convertField(v, childFieldName, reflectedValueField)
		}

		return v
	case reflect.String:
		return o.hydrateField(v, fieldName, reflectedValue.String())
	case reflect.Bool:
		return o.hydrateField(v, fieldName, strconv.FormatBool(reflectedValue.Bool()))
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return o.hydrateField(v, fieldName, strconv.FormatInt(reflectedValue.Int(), 10))
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return o.hydrateField(v, fieldName, strconv.FormatUint(reflectedValue.Uint(), 10))
	case reflect.Float32, reflect.Float64:
		if o.isMonetaryField(fieldName) {
			return o.hydrateField(v, fieldName, fmt.Sprintf("%.2f", reflectedValue.Float()))
		}
		return o.hydrateField(v, fieldName, fmt.Sprintf("%f", reflectedValue.Float()))
	case reflect.Ptr:
		return o.convertField(v, fieldName, reflectedValue.Elem())
	}

	return o.hydrateField(v, fieldName, fmt.Sprintf("%v", reflectedValue.Interface()))
}

func (o *SerializedOutputHandlerProvider) hydrateField(v *OrderedFields, fieldName string, fieldValue string) *OrderedFields {
	if !o.isIgnoredField(fieldName) {
		v.Set(fieldName, NewFieldValue(fieldValue, o.isDefaultField(fieldName)))
	}

	return v
}

func (o *SerializedOutputHandlerProvider) isDefaultField(name string) bool {
	return o.fieldInFields(name, o.DefaultFields)
}

func (o *SerializedOutputHandlerProvider) isIgnoredField(name string) bool {
	return o.fieldInFields(name, o.IgnoredFields)
}

func (o *SerializedOutputHandlerProvider) isMonetaryField(name string) bool {
	return o.fieldInFields(name, o.MonetaryFields)
}

func (o *SerializedOutputHandlerProvider) fieldInFields(name string, fields []string) bool {
	for _, field := range fields {
		if strings.ToLower(field) == strings.ToLower(name) {
			return true
		}
	}

	return false
}
