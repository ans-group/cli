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
	defaultFields     []string
	ignoredFields     []string
	monetaryFields    []string
	fieldHandlerFuncs map[string]FieldHandlerFunc
}

type FieldHandlerFunc func(v *OrderedFields, fieldName string, reflectedValue reflect.Value) *OrderedFields

func NewSerializedOutputHandlerDataProvider(items interface{}) *SerializedOutputHandlerDataProvider {
	return &SerializedOutputHandlerDataProvider{
		GenericOutputHandlerDataProvider: NewGenericOutputHandlerDataProvider(
			WithData(items),
		),
		fieldHandlerFuncs: make(map[string]FieldHandlerFunc),
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
	o.monetaryFields = append(o.monetaryFields, fields...)
	return o
}

func (o *SerializedOutputHandlerDataProvider) WithFieldHandler(fieldName string, f FieldHandlerFunc) *SerializedOutputHandlerDataProvider {
	if o.fieldHandlerFuncs == nil {
		o.fieldHandlerFuncs = make(map[string]FieldHandlerFunc)
	}
	o.fieldHandlerFuncs[fieldName] = f
	return o
}

func (o *SerializedOutputHandlerDataProvider) WithMultipleFieldHandler(fieldNames []string, f FieldHandlerFunc) *SerializedOutputHandlerDataProvider {
	for _, fieldName := range fieldNames {
		o.WithFieldHandler(fieldName, f)
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
	if o.fieldHandlerFuncs[fieldName] != nil {
		return o.fieldHandlerFuncs[fieldName](v, fieldName, reflectedValue)
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
		if o.isMonetaryField(fieldName) {
			return o.hydrateField(v, fieldName, fmt.Sprintf("%.2f", reflectedValue.Float()))
		}
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

func (o *SerializedOutputHandlerDataProvider) isMonetaryField(name string) bool {
	return o.fieldInFields(name, o.monetaryFields)
}

func (o *SerializedOutputHandlerDataProvider) fieldInFields(name string, fields []string) bool {
	for _, field := range fields {
		if field == name {
			return true
		}
	}

	return false
}
