package cmd

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
)

type testOutputDataProvider struct {
	GetFieldDataError error
}
type testOutputData struct {
	TestProperty1 string
	TestProperty2 string
}

func (o testOutputDataProvider) GetData() interface{} {
	return testOutputData{TestProperty1: "testvalue1", TestProperty2: "testvalue2"}
}

func (o testOutputDataProvider) GetFieldData() ([]*output.OrderedFields, error) {
	var data []*output.OrderedFields
	fields1 := output.NewOrderedFields()
	fields1.Set("test_property_1", output.NewFieldValue("fields1 test value 1", true))
	fields1.Set("test_property_2", output.NewFieldValue("fields1 test value 2", true))
	fields2 := output.NewOrderedFields()
	fields2.Set("test_property_1", output.NewFieldValue("fields2 test value 1", true))
	fields2.Set("test_property_2", output.NewFieldValue("fields2 test value 2", true))

	data = append(data, fields1, fields2)
	return data, o.GetFieldDataError
}

func TestGenericOutputHandlerProvider_GetData_ReturnsData(t *testing.T) {
	data := testOutputData{}
	o := NewGenericOutputHandlerProvider(data, nil, nil)

	output := o.GetData()

	assert.Equal(t, data, output)
}

func TestGenericOutputHandlerProvider_GetFieldData(t *testing.T) {
	t.Run("SingleStruct_ReturnsExpectedFields", func(t *testing.T) {
		data := testOutputData{}
		o := NewGenericOutputHandlerProvider(data, nil, nil)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output, 1)
		assert.Len(t, output[0].Keys(), 2)
	})

	t.Run("MultipleStructs_ReturnsExpectedFields", func(t *testing.T) {
		data1 := testOutputData{}
		data2 := testOutputData{}
		o := NewGenericOutputHandlerProvider([]testOutputData{data1, data2}, nil, nil)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output, 2)
		assert.Len(t, output[0].Keys(), 2)
	})

	t.Run("WithJsonTags_ReturnsExpectedFieldNames", func(t *testing.T) {
		type testType struct {
			Property1 string `json:"new_property_1"`
		}

		data := testType{}
		o := NewGenericOutputHandlerProvider(data, nil, nil)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output[0].Keys(), 1)
		assert.Equal(t, "new_property_1", output[0].Keys()[0])
	})

	t.Run("WithoutJsonTags_ReturnsExpectedFieldNames", func(t *testing.T) {
		type testType struct {
			Property1 string
		}

		data := testType{}
		o := NewGenericOutputHandlerProvider(data, nil, nil)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output[0].Keys(), 1)
		assert.Equal(t, "property_1", output[0].Keys()[0])
	})

	t.Run("UnexportedProperty_ReturnsExpectedFields", func(t *testing.T) {
		type testType struct {
			Property1 string
			property2 string
		}

		data := testType{}
		o := NewGenericOutputHandlerProvider(data, nil, nil)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output, 1)
		assert.Len(t, output[0].Keys(), 1)
	})

	t.Run("ArrayProperty_ReturnsExpectedFields", func(t *testing.T) {
		type testType struct {
			Property1 []string
		}

		data := testType{Property1: []string{"some", "value"}}
		o := NewGenericOutputHandlerProvider(data, nil, nil)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output[0].Keys(), 1)
		assert.Equal(t, "some, value", output[0].Get("property_1").Value)
	})
}

func TestGenericOutputHandlerProvider_isDefaultField(t *testing.T) {
	t.Run("ItemInDefaultFields_ReturnsTrue", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, []string{"field1", "field2", "field3"}, nil)
		defaultField := o.isDefaultField("field2")

		assert.True(t, defaultField)
	})

	t.Run("ItemNotInDefaultFields_ReturnsFalse", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, []string{"field1", "field2", "field3"}, nil)
		defaultField := o.isDefaultField("somefield")

		assert.False(t, defaultField)
	})
}

func TestGenericOutputHandlerProvider_isIgnoredField(t *testing.T) {
	t.Run("ItemInIgnoredFields_ReturnsTrue", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, []string{"field1", "field2", "field3"})
		defaultField := o.isIgnoredField("field2")

		assert.True(t, defaultField)
	})

	t.Run("ItemNotInIgnoredFields_ReturnsFalse", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, []string{"field1", "field2", "field3"})
		defaultField := o.isIgnoredField("somefield")

		assert.False(t, defaultField)
	})
}

func TestGenericOutputHandlerProvider_fieldToString(t *testing.T) {
	t.Run("StringType_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := "somestring"

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "somestring", output)
	})

	t.Run("BoolType_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := true

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "true", output)
	})

	t.Run("IntType_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := int(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Int8Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := int8(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Int16Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := int16(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Int32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := int32(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Int64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := int64(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("UintType_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := uint(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Uint8Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := uint8(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Uint16Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := uint16(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Uint32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := uint32(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Uint64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := uint64(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Float32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := 123.4

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123.400000", output)
	})

	t.Run("Float64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		v := float64(123.4)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123.400000", output)
	})

	t.Run("UnknownType_ReturnsExpectedString", func(t *testing.T) {
		o := NewGenericOutputHandlerProvider(nil, nil, nil)
		type somestruct struct{}
		v := somestruct{}

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "{}", output)
	})
}
