package cmd

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/cli/test"
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

func TestOutputHandler_Handle(t *testing.T) {
	t.Run("JSONFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "json")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "{\"TestProperty1\":\"testvalue1\",\"TestProperty2\":\"testvalue2\"}", output)
	})

	t.Run("TemplateFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "template")
		handler.Template = "{{ .TestProperty1 }}"

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "testvalue1\n", output)
	})

	t.Run("ValueFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "value")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "fields1 test value 1 fields1 test value 2\nfields2 test value 1 fields2 test value 2\n", output)
	})

	t.Run("ValueFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := testOutputDataProvider{
			GetFieldDataError: errors.New("test error 1"),
		}

		handler := NewOutputHandler(prov, "value")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("CSVFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "csv")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "test_property_1,test_property_2\nfields1 test value 1,fields1 test value 2\nfields2 test value 1,fields2 test value 2\n", output)
	})

	t.Run("CSVFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := testOutputDataProvider{
			GetFieldDataError: errors.New("test error 1"),
		}

		handler := NewOutputHandler(prov, "csv")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("TableFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "table")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", output)
	})

	t.Run("TableFormat_GetFieldDataError_ReturnsError", func(t *testing.T) {
		prov := testOutputDataProvider{
			GetFieldDataError: errors.New("test error 1"),
		}

		handler := NewOutputHandler(prov, "table")

		err := handler.Handle()

		assert.NotNil(t, err)
		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("EmptyFormat_ExpectedTableOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "")

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", output)
	})

	t.Run("InvalidFormat_ExpectedTableOutputWithStdErrError", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "invalidformat")

		stdOut, stdErr := test.CatchStdOutStdErr(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "+----------------------+----------------------+\n|   TEST PROPERTY 1    |   TEST PROPERTY 2    |\n+----------------------+----------------------+\n| fields1 test value 1 | fields1 test value 2 |\n| fields2 test value 1 | fields2 test value 2 |\n+----------------------+----------------------+\n", stdOut)
		assert.Equal(t, "Invalid output format [invalidformat], defaulting to 'table'\n", stdErr)
	})

	t.Run("UnsupportedFormat_NoUnsupportedHandler_ExpectedError", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "table")
		handler.SupportedFormats = []string{"json"}

		err := handler.Handle()

		assert.Equal(t, "Unsupported output format [table], supported formats: json", err.Error())
	})

	t.Run("UnsupportedFormat_NoUnsupportedHandler_ExpectedError", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "table")
		handler.SupportedFormats = []string{"json"}

		err := handler.Handle()

		assert.Equal(t, "Unsupported output format [table], supported formats: json", err.Error())
	})

	t.Run("UnsupportedFormat_WithUnsupportedHandler_CallsHandler", func(t *testing.T) {
		called := false

		handler := NewOutputHandler(testOutputDataProvider{}, "table")
		handler.SupportedFormats = []string{"json"}
		handler.UnsupportedFormatHandler = func() error {
			called = true
			return nil
		}

		handler.Handle()

		assert.True(t, called)
	})

	t.Run("UnsupportedFormat_Error_ReturnsError", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "table")
		handler.SupportedFormats = []string{"json"}
		handler.UnsupportedFormatHandler = func() error {
			return errors.New("test error 1")
		}

		err := handler.Handle()

		assert.Equal(t, "test error 1", err.Error())
	})

	t.Run("SupportedFormat_ExpectedOutput", func(t *testing.T) {
		handler := NewOutputHandler(testOutputDataProvider{}, "value")
		handler.SupportedFormats = []string{"value"}

		output := test.CatchStdOut(t, func() {
			handler.Handle()
		})

		assert.Equal(t, "fields1 test value 1 fields1 test value 2\nfields2 test value 1 fields2 test value 2\n", output)
	})
}

func testResetViper() {
	viper.SetDefault("command_wait_timeout_seconds", 1200)
	viper.SetDefault("command_wait_sleep_seconds", 5)
}

func TestWaitForCommand(t *testing.T) {

	t.Run("SuccessAfter3Attempts", func(t *testing.T) {
		testResetViper()
		defer testResetViper()
		viper.SetDefault("command_wait_sleep_seconds", 1)

		attempt := 1
		f := func() (bool, error) {
			if attempt == 3 {
				return true, nil
			}

			attempt++
			return false, nil
		}

		r := WaitForCommand(f)

		assert.Nil(t, r)
		assert.Equal(t, 3, attempt)
	})

	t.Run("ErrorAfter3Attempts", func(t *testing.T) {
		testResetViper()
		defer testResetViper()
		viper.SetDefault("command_wait_sleep_seconds", 1)

		attempt := 1
		f := func() (bool, error) {
			if attempt == 3 {
				return false, errors.New("test error")
			}

			attempt++
			return false, nil
		}

		r := WaitForCommand(f)

		assert.NotNil(t, r)
		assert.Equal(t, 3, attempt)
	})

	t.Run("InvalidWaitTimeout", func(t *testing.T) {
		testResetViper()
		defer testResetViper()
		viper.SetDefault("command_wait_timeout_seconds", 0)

		attempt := 1
		f := func() (bool, error) {
			attempt++
			return false, nil
		}

		r := WaitForCommand(f)

		assert.NotNil(t, r)
		assert.Equal(t, 1, attempt)
	})

	t.Run("InvalidWaitSleep", func(t *testing.T) {
		testResetViper()
		defer testResetViper()
		viper.SetDefault("command_wait_sleep_seconds", 0)

		attempt := 1
		f := func() (bool, error) {
			attempt++
			return false, nil
		}

		r := WaitForCommand(f)

		assert.NotNil(t, r)
		assert.Equal(t, 1, attempt)
	})
}
