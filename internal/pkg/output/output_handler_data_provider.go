package output

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/iancoleman/strcase"
)

type OutputHandlerDataProvider interface {
	GetData() interface{}
	GetFieldData() ([]*OrderedFields, error)
}

type OutputHandlerDataProviderOption func(p *GenericOutputHandlerDataProvider)

type GenericOutputHandlerDataProvider struct {
	data          interface{}
	fieldDataFunc func() ([]*OrderedFields, error)
}

func NewGenericOutputHandlerDataProvider(opts ...OutputHandlerDataProviderOption) *GenericOutputHandlerDataProvider {
	p := &GenericOutputHandlerDataProvider{}
	for _, opt := range opts {
		opt(p)
	}
	return p
}

func WithData(data interface{}) OutputHandlerDataProviderOption {
	return func(p *GenericOutputHandlerDataProvider) {
		p.data = data
	}
}

func WithFieldDataFunc(fieldDataFunc func() ([]*OrderedFields, error)) OutputHandlerDataProviderOption {
	return func(p *GenericOutputHandlerDataProvider) {
		p.fieldDataFunc = fieldDataFunc
	}
}

func (p *GenericOutputHandlerDataProvider) GetData() interface{} {
	return p.data
}

func (p *GenericOutputHandlerDataProvider) GetFieldData() ([]*OrderedFields, error) {
	return p.fieldDataFunc()
}

type SerializedOutputHandlerDataProvider struct {
	*GenericOutputHandlerDataProvider
	defaultFields      []string
	ignoredFields      []string
	monetaryFields     []string
	fieldHandlers      map[string]FieldHandlerFunc
	fieldValueHandlers map[string]FieldValueHandlerFunc
}

type FieldHandlerFunc func(v *OrderedFields, fieldName string, reflectedValue reflect.Value) *OrderedFields

type FieldValueHandlerFunc func(reflectedValue reflect.Value) string

var MonetaryFieldValueHandler = func(reflectedValue reflect.Value) string {
	return fmt.Sprintf("%.2f", reflectedValue.Float())
}

func NewSerializedOutputHandlerDataProvider(items interface{}) *SerializedOutputHandlerDataProvider {
	return &SerializedOutputHandlerDataProvider{
		GenericOutputHandlerDataProvider: NewGenericOutputHandlerDataProvider(
			WithData(items),
		),
	}
}

func (o *SerializedOutputHandlerDataProvider) WithDefaultFields(fields []string) *SerializedOutputHandlerDataProvider {
	o.defaultFields = append(o.defaultFields, fields...)
	return o
}

func (o *SerializedOutputHandlerDataProvider) WithIgnoredFields(fields []string) *SerializedOutputHandlerDataProvider {
	o.ignoredFields = append(o.ignoredFields, fields...)
	return o
}

func (o *SerializedOutputHandlerDataProvider) WithMonetaryFields(fields []string) *SerializedOutputHandlerDataProvider {
	o.WithFieldValueHandler(MonetaryFieldValueHandler, fields...)
	return o
}

func (o *SerializedOutputHandlerDataProvider) WithFieldHandler(f FieldHandlerFunc, fieldNames ...string) *SerializedOutputHandlerDataProvider {
	if o.fieldHandlers == nil {
		o.fieldHandlers = make(map[string]FieldHandlerFunc)
	}

	for _, fieldName := range fieldNames {
		o.fieldHandlers[fieldName] = f
	}
	return o
}

func (o *SerializedOutputHandlerDataProvider) WithFieldValueHandler(f FieldValueHandlerFunc, fieldNames ...string) *SerializedOutputHandlerDataProvider {
	if o.fieldValueHandlers == nil {
		o.fieldValueHandlers = make(map[string]FieldValueHandlerFunc)
	}

	for _, fieldName := range fieldNames {
		o.fieldValueHandlers[fieldName] = f
	}
	return o
}

func (o *SerializedOutputHandlerDataProvider) GetFieldData() ([]*OrderedFields, error) {
	return o.convert(reflect.ValueOf(o.GetData())), nil
}

func (o *SerializedOutputHandlerDataProvider) convert(reflectedValue reflect.Value) []*OrderedFields {
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

func (o *SerializedOutputHandlerDataProvider) convertField(v *OrderedFields, fieldName string, reflectedValue reflect.Value) *OrderedFields {
	if o.fieldHandlers != nil && o.fieldHandlers[fieldName] != nil {
		return o.fieldHandlers[fieldName](v, fieldName, reflectedValue)
	}

	if o.fieldValueHandlers != nil && o.fieldValueHandlers[fieldName] != nil {
		return o.hydrateField(v, fieldName, o.fieldValueHandlers[fieldName](reflectedValue))
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
		return o.hydrateField(v, fieldName, fmt.Sprintf("%f", reflectedValue.Float()))
	case reflect.Ptr:
		if reflectedValue.IsNil() {
			return o.hydrateField(v, fieldName, fmt.Sprintf("%v", reflectedValue.Interface()))
		}
		return o.convertField(v, fieldName, reflectedValue.Elem())
	case reflect.Invalid:
		return nil
	}

	return o.hydrateField(v, fieldName, fmt.Sprintf("%v", reflectedValue.Interface()))
}

func (o *SerializedOutputHandlerDataProvider) hydrateField(v *OrderedFields, fieldName string, fieldValue string) *OrderedFields {
	if !o.isIgnoredField(fieldName) {
		v.Set(fieldName, NewFieldValue(fieldValue, o.isDefaultField(fieldName)))
	}

	return v
}

func (o *SerializedOutputHandlerDataProvider) isDefaultField(name string) bool {
	return o.fieldInFields(name, o.defaultFields)
}

func (o *SerializedOutputHandlerDataProvider) isIgnoredField(name string) bool {
	return o.fieldInFields(name, o.ignoredFields)
}

func (o *SerializedOutputHandlerDataProvider) fieldInFields(name string, fields []string) bool {
	for _, field := range fields {
		if field == name {
			return true
		}
	}

	return false
}
