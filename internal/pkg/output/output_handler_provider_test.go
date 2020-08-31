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

func TestSerializedOutputHandlerProvider_convertField(t *testing.T) {
	t.Run("StringType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := "somestring"

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "somestring", output.Get("test_field").Value)
	})

	t.Run("BoolType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := true

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "true", output.Get("test_field").Value)
	})

	t.Run("IntType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := 123

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Int8Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int8(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Int16Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int16(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Int32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int32(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Int64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := int64(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("UintType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Uint8Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint8(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Uint16Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint16(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Uint32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint32(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Uint64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := uint64(123)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123", output.Get("test_field").Value)
	})

	t.Run("Float32Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := 123.4

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123.400000", output.Get("test_field").Value)
	})

	t.Run("Float64Type_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
		v := float64(123.4)

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "123.400000", output.Get("test_field").Value)
	})

	t.Run("StructType_ReturnsExpectedString", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil)
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
		o := NewSerializedOutputHandlerProvider(nil)
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

	t.Run("FieldHandler_Expected", func(t *testing.T) {
		o := NewSerializedOutputHandlerProvider(nil).WithFieldHandler("test_field", func(v *OrderedFields, fieldName string, reflectedValue reflect.Value) *OrderedFields {
			v.Set("test_field", NewFieldValue("testvaluefromhandler", false))
			return v
		})
		v := "somestring"

		output := o.convertField(NewOrderedFields(), "test_field", reflect.ValueOf(v))

		assert.Equal(t, "testvaluefromhandler", output.Get("test_field").Value)
	})

	t.Run("WithJsonTags_ReturnsExpectedFieldNames", func(t *testing.T) {
		type testType struct {
			Property1 string `json:"new_property_1"`
		}

		o := NewSerializedOutputHandlerProvider(nil)
		v := testType{}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Len(t, output.Keys(), 1)
		assert.Equal(t, "new_property_1", output.Keys()[0])
	})

	t.Run("WithoutJsonTags_ReturnsExpectedFieldNames", func(t *testing.T) {
		type testType struct {
			Property1 string
		}

		o := NewSerializedOutputHandlerProvider(nil)
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

		o := NewSerializedOutputHandlerProvider(nil)
		v := testType{}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Len(t, output.Keys(), 1)
	})

	t.Run("ArrayProperty_ReturnsExpectedFields", func(t *testing.T) {
		type testType struct {
			Property1 []string
		}

		o := NewSerializedOutputHandlerProvider(nil)
		v := testType{Property1: []string{"some", "value"}}

		output := o.convertField(NewOrderedFields(), "", reflect.ValueOf(v))

		assert.Len(t, output.Keys(), 1)
		assert.Equal(t, "[some value]", output.Get("property_1").Value)
	})
}
