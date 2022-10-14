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

func TestSerializedOutputHandlerDataProvider_GetData_ReturnsData(t *testing.T) {
	data := testOutputData{}
	o := NewSerializedOutputHandlerDataProvider(data)

	output := o.GetData()

	assert.Equal(t, data, output)
}

func TestSerializedOutputHandlerDataProvider_GetFieldData(t *testing.T) {
	t.Run("SingleStruct_ReturnsExpectedFields", func(t *testing.T) {
		data := testOutputData{}
		o := NewSerializedOutputHandlerDataProvider(data)

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output, 1)
		assert.Len(t, output[0].Keys(), 2)
	})

	t.Run("MultipleStructs_ReturnsExpectedFields", func(t *testing.T) {
		data1 := testOutputData{}
		data2 := testOutputData{}
		o := NewSerializedOutputHandlerDataProvider([]testOutputData{data1, data2})

		output, err := o.GetFieldData()

		assert.Nil(t, err)
		assert.Len(t, output, 2)
		assert.Len(t, output[0].Keys(), 2)
	})
}

func TestSerializedOutputHandlerDataProvider_isIgnoredField(t *testing.T) {
	t.Run("ItemInIgnoredFields_ReturnsTrue", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil).WithIgnoredFields([]string{"field1", "field2", "field3"})
		defaultField := o.isIgnoredField("field2")

		assert.True(t, defaultField)
	})

	t.Run("ItemNotInIgnoredFields_ReturnsFalse", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil).WithIgnoredFields([]string{"field1", "field2", "field3"})
		defaultField := o.isIgnoredField("somefield")

		assert.False(t, defaultField)
	})
}

func TestSerializedOutputHandlerDataProvider_convertField(t *testing.T) {
	t.Run("StringType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := "somestring"

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "somestring", output.Get("test_field").Value)
	})

	t.Run("BoolType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := true

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "true", output.Get("test_field").Value)
	})

	t.Run("IntType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := 123

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Int8Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := int8(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Int16Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := int16(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Int32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := int32(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Int64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := int64(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("UintType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := uint(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Uint8Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := uint8(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Uint16Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := uint16(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Uint32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := uint32(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Uint64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := uint64(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Float32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := 123.4

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123.400000", output.Get("test_field").Value)
	})

	t.Run("Float64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		v := float64(123.4)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123.400000", output.Get("test_field").Value)
	})

	t.Run("StructType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		type somestruct struct {
			Field1 string `json:"field1"`
			Field2 string `json:"field2"`
		}
		v := somestruct{
			Field1: "field1value",
			Field2: "field2value",
		}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Equal(t, "field1value", output.Get("field1").Value)
		assert.Equal(t, "field2value", output.Get("field2").Value)
	})

	t.Run("NestedStructType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		type somestruct struct {
			Struct1 struct {
				Field1 string `json:"field1"`
			} `json:"struct1"`
		}
		v := somestruct{
			Struct1: struct {
				Field1 string `json:"field1"`
			}{Field1: "field1value"},
		}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Equal(t, "field1value", output.Get("struct1_field1").Value)
	})

	t.Run("StructTypeWithNil_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerDataProvider(nil)
		type somestruct struct {
			Field1 *string `json:"field1"`
		}
		v := somestruct{}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))
		assert.Equal(t, "<nil>", output.Get("field1").Value)
	})

	t.Run("WithJsonTags_ReturnsExpectedFieldNames", func(t *testing.T) {
		type testType struct {
			Property1 string `json:"new_property_1"`
		}

		o := NewSerializedOutputHandlerDataProvider(nil)
		v := testType{}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Len(t, output.Keys(), 1)
		assert.Equal(t, "new_property_1", output.Keys()[0])
	})

	t.Run("WithoutJsonTags_ReturnsExpectedFieldNames", func(t *testing.T) {
		type testType struct {
			Property1 string
		}

		o := NewSerializedOutputHandlerDataProvider(nil)
		v := testType{}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Len(t, output.Keys(), 1)
		assert.Equal(t, "property_1", output.Keys()[0])
	})

	t.Run("UnexportedProperty_ReturnsExpectedFields", func(t *testing.T) {
		type testType struct {
			Property1 string
			property2 string
		}

		o := NewSerializedOutputHandlerDataProvider(nil)
		v := testType{}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Len(t, output.Keys(), 1)
	})

	t.Run("ArrayProperty_ReturnsExpectedFields", func(t *testing.T) {
		type testType struct {
			Property1 []string
		}

		o := NewSerializedOutputHandlerDataProvider(nil)
		v := testType{Property1: []string{"some", "value"}}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Len(t, output.Keys(), 1)
		assert.Equal(t, "[some value]", output.Get("property_1").Value)
	})
}
