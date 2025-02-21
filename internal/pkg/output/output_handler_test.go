package output

import (
	"testing"

	"github.com/ans-group/cli/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type testModel struct {
	TestProperty1 string `json:"testproperty1"`
	TestProperty2 string `json:"testproperty2"`
	TestProperty3 string `json:"testproperty3"`
}

type testModelCollection []testModel

var collectionSingleRow = testModelCollection([]testModel{{"Row1TestValue1", "Row1TestValue2", "Row1TestValue3"}})
var collectionMultipleRows = testModelCollection([]testModel{{"Row1TestValue1", "Row1TestValue2", "Row1TestValue3"}, {"Row2TestValue1", "Row2TestValue2", "Row2TestValue3"}})

func TestOutputHandler_JSON(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		o := NewOutputHandler()

		output := test.CatchStdOut(t, func() {
			err := o.JSON("test")
			assert.NoError(t, err)
		})

		assert.Equal(t, "\"test\"", output)
	})

	t.Run("MarshalError", func(t *testing.T) {
		o := NewOutputHandler()

		err := o.JSON(func() {})

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to marshal json")
	})
}

func TestOutputHandler_Value(t *testing.T) {
	t.Run("SingleRowDefaultFields_ExpectedStdout", func(t *testing.T) {
		o := NewOutputHandler()

		output := test.CatchStdOut(t, func() {
			err := o.Value(&cobra.Command{}, collectionSingleRow)
			assert.NoError(t, err)
		})

		assert.Equal(t, "Row1TestValue1 Row1TestValue2 Row1TestValue3\n", output)
	})

	t.Run("MultipleRowsDefaultFields_ExpectedStdout", func(t *testing.T) {
		o := NewOutputHandler()
		output := test.CatchStdOut(t, func() {
			err := o.Value(&cobra.Command{}, collectionMultipleRows)
			assert.NoError(t, err)
		})

		assert.Equal(t, "Row1TestValue1 Row1TestValue2 Row1TestValue3\nRow2TestValue1 Row2TestValue2 Row2TestValue3\n", output)
	})
}

func TestOutputHandler_JSONPath(t *testing.T) {
	o := NewOutputHandler()

	output := test.CatchStdOut(t, func() {
		err := o.JSONPath("{[].TestProperty1}", collectionSingleRow)
		assert.NoError(t, err)
	})

	assert.Equal(t, "Row1TestValue1", output)
}

func TestOutputHandler_CSV(t *testing.T) {
	t.Run("SingleRowDefaultFields_ExpectedStdout", func(t *testing.T) {
		o := NewOutputHandler()

		output := test.CatchStdOut(t, func() {
			o.CSV(&cobra.Command{}, collectionSingleRow)
		})

		assert.Equal(t, "testproperty1,testproperty2,testproperty3\nRow1TestValue1,Row1TestValue2,Row1TestValue3\n", output)
	})

	t.Run("MultipleRowsDefaultFields_ExpectedStdout", func(t *testing.T) {
		o := NewOutputHandler()

		output := test.CatchStdOut(t, func() {
			o.CSV(&cobra.Command{}, collectionMultipleRows)
		})

		assert.Equal(t, "testproperty1,testproperty2,testproperty3\nRow1TestValue1,Row1TestValue2,Row1TestValue3\nRow2TestValue1,Row2TestValue2,Row2TestValue3\n", output)
	})
}

func TestOutputHandler_YAML(t *testing.T) {
	o := NewOutputHandler()

	output := test.CatchStdOut(t, func() {
		err := o.YAML("test")
		assert.NoError(t, err)
	})

	assert.Equal(t, "test\n", output)
}

func TestOutputHandler_Table(t *testing.T) {
	o := NewOutputHandler()

	output := test.CatchStdOut(t, func() {
		err := o.Table(&cobra.Command{}, collectionSingleRow)
		assert.NoError(t, err)
	})

	assert.Equal(t, "+----------------+----------------+----------------+\n| TESTPROPERTY1  | TESTPROPERTY2  | TESTPROPERTY3  |\n+----------------+----------------+----------------+\n| Row1TestValue1 | Row1TestValue2 | Row1TestValue3 |\n+----------------+----------------+----------------+\n", output)
}

func TestOutputHandler_List(t *testing.T) {
	o := NewOutputHandler()

	output := test.CatchStdOut(t, func() {
		err := o.List(&cobra.Command{}, collectionSingleRow)
		assert.NoError(t, err)
	})

	assert.Equal(t, "testproperty1 : Row1TestValue1\ntestproperty2 : Row1TestValue2\ntestproperty3 : Row1TestValue3\n", output)
}

func TestOutputHandler_Template(t *testing.T) {
	o := NewOutputHandler()

	output := test.CatchStdOut(t, func() {
		err := o.Template("{{.TestProperty1}}", collectionSingleRow)
		assert.NoError(t, err)
	})

	assert.Equal(t, "Row1TestValue1\n", output)
}
