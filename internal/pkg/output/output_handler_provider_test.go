package output

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testOutputData struct {
	TestProperty1 string
	TestProperty2 string
}

func TestSerializedOutputHandlerProvider_GetData_ReturnsData(t *testing.T) {
	data := testOutputData{}
	o := NewSerializedOutputHandlerProvider(data)

	output := o.GetData()

	assert.Equal(t, data, output)
}

func TestSerializedOutputHandlerProvider_GetFieldData(t *testing.T) {
	t.Run("SingleStruct_ReturnsExpectedFields", func(t *testing.T) {
		data := testOutputData{}
		o := NewSerializedOutputHandlerProvider(data)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output, 1)
		assert.Len(t, output[0].Keys(), 2)
	})

	t.Run("MultipleStructs_ReturnsExpectedFields", func(t *testing.T) {
		data1 := testOutputData{}
		data2 := testOutputData{}
		o := NewSerializedOutputHandlerProvider([]testOutputData{data1, data2})

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
		o := NewSerializedOutputHandlerProvider(data)

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
		o := NewSerializedOutputHandlerProvider(data)

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
		o := NewSerializedOutputHandlerProvider(data)

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
		o := NewSerializedOutputHandlerProvider(data)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output[0].Keys(), 1)
		assert.Equal(t, "some, value", output[0].Get("property_1").Value)
	})
}

func TestSerializedOutputHandlerProvider_isDefaultField(t *testing.T) {
	t.Run("ItemInDefaultFields_ReturnsTrue", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil).WithDefaultFields([]string{"field1", "field2", "field3"})
		defaultField := o.isDefaultField("field2")

		assert.True(t, defaultField)
	})

	t.Run("ItemNotInDefaultFields_ReturnsFalse", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil).WithDefaultFields([]string{"field1", "field2", "field3"})
		defaultField := o.isDefaultField("somefield")

		assert.False(t, defaultField)
	})
}

func TestSerializedOutputHandlerProvider_isIgnoredField(t *testing.T) {
	t.Run("ItemInIgnoredFields_ReturnsTrue", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil).WithIgnoredFields([]string{"field1", "field2", "field3"})
		defaultField := o.isIgnoredField("field2")

		assert.True(t, defaultField)
	})

	t.Run("ItemNotInIgnoredFields_ReturnsFalse", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil).WithIgnoredFields([]string{"field1", "field2", "field3"})
		defaultField := o.isIgnoredField("somefield")

		assert.False(t, defaultField)
	})
}

func TestSerializedOutputHandlerProvider_fieldToString(t *testing.T) {
	t.Run("StringType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := "somestring"

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "somestring", output)
	})

	t.Run("BoolType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := true

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "true", output)
	})

	t.Run("IntType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Int8Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int8(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Int16Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int16(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Int32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int32(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Int64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int64(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("UintType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Uint8Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint8(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Uint16Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint16(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Uint32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint32(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Uint64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint64(123)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123", output)
	})

	t.Run("Float32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := 123.4

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123.400000", output)
	})

	t.Run("Float64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := float64(123.4)

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "123.400000", output)
	})

	t.Run("UnknownType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		type somestruct struct{}
		v := somestruct{}

		output := o.fieldToString(reflect.ValueOf(v))

		assert.Equal(t, "{}", output)
	})
}
