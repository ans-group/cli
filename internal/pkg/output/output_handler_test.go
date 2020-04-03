package output

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test"
)

var testOutputHandlerProvider = NewGenericOutputHandlerProvider(
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
		handler := NewOutputHandler(testOutputHandlerProvider, "json")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "{\"TestProperty1\":\"testvalue1\",\"TestProperty2\":\"testvalue2\"}", output)
	})

	t.Run("TemplateFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerProvider, "template")
		handler.Template = "{{ .TestProperty1 }}"

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "testvalue1\n", output)
	})

	t.Run("ValueFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerProvider, "value")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "fields1 test value 1 fields1 test value 2\nfields2 test value 1 fields2 test value 2\n", output)
	})

	t.Run("ValueFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := NewGenericOutputHandlerProvider(
			WithFieldDataFunc(func() ([]*OrderedFields, error) {
				return nil, errors.New("test error 1")
			}),
		)

		handler := NewOutputHandler(prov, "value")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("CSVFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerProvider, "csv")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "test_property_1,test_property_2\nfields1 test value 1,fields1 test value 2\nfields2 test value 1,fields2 test value 2\n", output)
	})

	t.Run("CSVFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := NewGenericOutputHandlerProvider(
			WithFieldDataFunc(func() ([]*OrderedFields, error) {
				return nil, errors.New("test error 1")
			}),
		)

		handler := NewOutputHandler(prov, "csv")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("TableFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerProvider, "table")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", output)
	})

	t.Run("TableFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := NewGenericOutputHandlerProvider(
			WithFieldDataFunc(func() ([]*OrderedFields, error) {
				return nil, errors.New("test error 1")
			}),
		)

		handler := NewOutputHandler(prov, "table")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("EmptyFormat_ExpectedTableOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerProvider, "")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", output)
	})

	t.Run("InvalidFormat_ExpectedTableOutputWithStdErrError", func(t *testing.T) {
		handler := NewOutputHandler(testOutputHandlerProvider, "invalidformat")

		stdOut, stdErr := test.CatchStdOutStdErr(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", stdOut)
		assert.Equal(t, "Invalid output format [invalidformat], defaulting to 'table'\n", stdErr)
	})

	t.Run("UnsupportedFormat_NoUnsupportedHandler_ExpectedError", func(t *testing.T) {
		handler := NewOutputHandler(NewGenericOutputHandlerProvider(WithSupportedFormats([]string{"json"})), "table")

		err := handler.Handle()

		assert.Equal(t, "Unsupported output format [table], supported formats: json", err.Error())
	})

	t.Run("UnsupportedFormat_NoUnsupportedHandler_ExpectedError", func(t *testing.T) {
		handler := NewOutputHandler(NewGenericOutputHandlerProvider(WithSupportedFormats([]string{"json"})), "table")

		err := handler.Handle()

		assert.Equal(t, "Unsupported output format [table], supported formats: json", err.Error())
	})

	t.Run("UnsupportedFormat_WithUnsupportedHandler_CallsHandler", func(t *testing.T) {
		called := false

		handler := NewOutputHandler(NewGenericOutputHandlerProvider(WithSupportedFormats([]string{"json"})), "table")
		handler.UnsupportedFormatHandler = func() error {
			called = true
			return nil
		}

		handler.Handle()

		assert.True(t, called)
	})

	t.Run("UnsupportedFormat_Error_ReturnsError", func(t *testing.T) {
		handler := NewOutputHandler(NewGenericOutputHandlerProvider(WithSupportedFormats([]string{"json"})), "table")
		handler.UnsupportedFormatHandler = func() error {
			return errors.New("test error 1")
		}

		err := handler.Handle()

		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("SupportedFormat_ExpectedOutput", func(t *testing.T) {
		prov := testOutputHandlerProvider
		prov.supportedFormats = []string{"value"}

		handler := NewOutputHandler(prov, "value")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "fields1 test value 1 fields1 test value 2\nfields2 test value 1 fields2 test value 2\n", output)
	})
}
