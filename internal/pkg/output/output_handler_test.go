package output

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test"
)

var testOutputHandlerDataProvider = NewGenericOutputHandlerDataProvider(
	WithData(testOutputData{TestProperty1: "testvalue1", TestProperty2: "testvalue2"}),
	WithFieldDataFunc(func() ([]*OrderedFields, error) {
		var data []*OrderedFields
		fields1 := NewOrderedFields()
		fields1.Set("test_property_1", NewFieldValue("fields1 test value 1", true))
		fields1.Set("test_property_2", NewFieldValue("fields1 test value 2", true))
		fields2 := NewOrderedFields()
		fields2.Set("test_property_1", NewFieldValue("fields2 test value 1", true))
		fields2.Set("test_property_2", NewFieldValue("fields2 test value 2", true))

		data = append(data, fields1, fields2)
		return data, nil
	}),
)

func TestOutputHandler_Handle(t *testing.T) {
	t.Run("JSONFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "json", "")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "{\"TestProperty1\":\"testvalue1\",\"TestProperty2\":\"testvalue2\"}", output)
	})

	t.Run("TemplateFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "template", "{{ .TestProperty1 }}")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "testvalue1\n", output)
	})

	t.Run("ValueFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "value", "")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "fields1 test value 1 fields1 test value 2\nfields2 test value 1 fields2 test value 2\n", output)
	})

	t.Run("ValueFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := NewGenericOutputHandlerDataProvider(
			WithFieldDataFunc(func() ([]*OrderedFields, error) {
				return nil, errors.New("test error 1")
			}),
		)

		handler := NewOutputHandler(prov, "value", "")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("CSVFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "csv", "")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "test_property_1,test_property_2\nfields1 test value 1,fields1 test value 2\nfields2 test value 1,fields2 test value 2\n", output)
	})

	t.Run("CSVFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := NewGenericOutputHandlerDataProvider(
			WithFieldDataFunc(func() ([]*OrderedFields, error) {
				return nil, errors.New("test error 1")
			}),
		)

		handler := NewOutputHandler(prov, "csv", "")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("ListFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "list", "")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "test_property_1 : fields1 test value 1\ntest_property_2 : fields1 test value 2\n\ntest_property_1 : fields2 test value 1\ntest_property_2 : fields2 test value 2\n", output)
	})

	t.Run("ListFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := NewGenericOutputHandlerDataProvider(
			WithFieldDataFunc(func() ([]*OrderedFields, error) {
				return nil, errors.New("test error 1")
			}),
		)

		handler := NewOutputHandler(prov, "list", "")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("JSONPathFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "jsonpath", "{}")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "{testvalue1 testvalue2}", output)
	})

	t.Run("TableFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "table", "")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", output)
	})

	t.Run("TableFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := NewGenericOutputHandlerDataProvider(
			WithFieldDataFunc(func() ([]*OrderedFields, error) {
				return nil, errors.New("test error 1")
			}),
		)

		handler := NewOutputHandler(prov, "table", "")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("EmptyFormat_ExpectedTableOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "", "")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", output)
	})

	t.Run("InvalidFormat_ExpectedTableOutputWithStdErrError", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerDataProvider, "invalidformat", "")

		stdOut, stdErr := test.CatchStdOutStdErr(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", stdOut)
		assert.Equal(t, "Invalid output format [invalidformat], defaulting to 'table'\n", stdErr)
	})

	t.Run("UnsupportedFormat_NoUnsupportedHandler_ExpectedError", func(t *testing.T) {
		handler := NewOutputHandler(NewGenericOutputHandlerDataProvider(), "table", "").WithSupportedFormats([]string{"json"})

		err := handler.Handle()

		assert.Equal(t, "Unsupported output format [table], supported formats: json", err.Error())
	})

	t.Run("UnsupportedFormat_NoUnsupportedHandler_ExpectedError", func(t *testing.T) {
		handler := NewOutputHandler(NewGenericOutputHandlerDataProvider(), "table", "").WithSupportedFormats([]string{"json"})

		err := handler.Handle()

		assert.Equal(t, "Unsupported output format [table], supported formats: json", err.Error())
	})

	t.Run("UnsupportedFormat_Error_ReturnsError", func(t *testing.T) {
		handler := NewOutputHandler(NewGenericOutputHandlerDataProvider(), "table", "").WithSupportedFormats([]string{"json"})

		err := handler.Handle()

		assert.NotNil(t, err)
	})

	t.Run("SupportedFormat_ExpectedOutput", func(t *testing.T) {
		prov := testOutputHandlerDataProvider

		handler := NewOutputHandler(prov, "value", "").WithSupportedFormats([]string{"value"})

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "fields1 test value 1 fields1 test value 2\nfields2 test value 1 fields2 test value 2\n", output)
	})
}
