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

	assert.Equal(t, "┌────────────────┬────────────────┬────────────────┐\n│ TESTPROPERTY 1 │ TESTPROPERTY 2 │ TESTPROPERTY 3 │\n├────────────────┼────────────────┼────────────────┤\n│ Row1TestValue1 │ Row1TestValue2 │ Row1TestValue3 │\n└────────────────┴────────────────┴────────────────┘\n", output)
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

func TestOutputHandler_WithAdditionalColumns(t *testing.T) {
	t.Run("NewOutputHandler_WithAdditionalColumns", func(t *testing.T) {
		o := NewOutputHandler(WithAdditionalColumns("testproperty1", "testproperty2"))
		assert.Equal(t, []string{"testproperty1", "testproperty2"}, o.additionalColumns)
	})

	t.Run("NewOutputHandler_WithAdditionalColumnsMultipleCalls", func(t *testing.T) {
		o := NewOutputHandler(
			WithAdditionalColumns("testproperty1"),
			WithAdditionalColumns("testproperty2"),
		)
		assert.Equal(t, []string{"testproperty1", "testproperty2"}, o.additionalColumns)
	})
}

func TestOutputHandler_AdditionalColumnsInTable(t *testing.T) {
	t.Run("AdditionalColumns_Table", func(t *testing.T) {
		o := NewOutputHandler(WithAdditionalColumns("testproperty1", "testproperty3"))

		output := test.CatchStdOut(t, func() {
			err := o.Table(&cobra.Command{}, collectionSingleRow)
			assert.NoError(t, err)
		})

		expected := "┌────────────────┬────────────────┐\n│ TESTPROPERTY 1 │ TESTPROPERTY 3 │\n├────────────────┼────────────────┤\n│ Row1TestValue1 │ Row1TestValue3 │\n└────────────────┴────────────────┘\n"
		assert.Equal(t, expected, output)
	})

	t.Run("AdditionalColumnsWithGlobPattern_Table", func(t *testing.T) {
		o := NewOutputHandler(WithAdditionalColumns("testproperty*"))

		output := test.CatchStdOut(t, func() {
			err := o.Table(&cobra.Command{}, collectionSingleRow)
			assert.NoError(t, err)
		})

		expected := "┌────────────────┬────────────────┬────────────────┐\n│ TESTPROPERTY 1 │ TESTPROPERTY 2 │ TESTPROPERTY 3 │\n├────────────────┼────────────────┼────────────────┤\n│ Row1TestValue1 │ Row1TestValue2 │ Row1TestValue3 │\n└────────────────┴────────────────┴────────────────┘\n"
		assert.Equal(t, expected, output)
	})
}

func TestOutputHandler_AdditionalColumnsInCSV(t *testing.T) {
	t.Run("AdditionalColumns_CSV", func(t *testing.T) {
		o := NewOutputHandler(WithAdditionalColumns("testproperty1", "testproperty3"))

		output := test.CatchStdOut(t, func() {
			err := o.CSV(&cobra.Command{}, collectionSingleRow)
			assert.NoError(t, err)
		})

		assert.Equal(t, "testproperty1,testproperty3\nRow1TestValue1,Row1TestValue3\n", output)
	})

	t.Run("AdditionalColumns_CSV_MultipleRows", func(t *testing.T) {
		o := NewOutputHandler(WithAdditionalColumns("testproperty2"))

		output := test.CatchStdOut(t, func() {
			err := o.CSV(&cobra.Command{}, collectionMultipleRows)
			assert.NoError(t, err)
		})

		assert.Equal(t, "testproperty2\nRow1TestValue2\nRow2TestValue2\n", output)
	})
}

func TestOutputHandler_AdditionalColumnsInList(t *testing.T) {
	t.Run("AdditionalColumns_List", func(t *testing.T) {
		o := NewOutputHandler(WithAdditionalColumns("testproperty1", "testproperty3"))

		output := test.CatchStdOut(t, func() {
			err := o.List(&cobra.Command{}, collectionSingleRow)
			assert.NoError(t, err)
		})

		assert.Equal(t, "testproperty1 : Row1TestValue1\ntestproperty3 : Row1TestValue3\n", output)
	})
}

func TestOutputHandler_AdditionalColumnsInValue(t *testing.T) {
	t.Run("AdditionalColumns_Value", func(t *testing.T) {
		o := NewOutputHandler(WithAdditionalColumns("testproperty2"))

		output := test.CatchStdOut(t, func() {
			err := o.Value(&cobra.Command{}, collectionSingleRow)
			assert.NoError(t, err)
		})

		assert.Equal(t, "Row1TestValue2\n", output)
	})
}

type testModelWithDefaults struct {
	TestProperty1 string `json:"testproperty1"`
	TestProperty2 string `json:"testproperty2"`
	TestProperty3 string `json:"testproperty3"`
}

func (t testModelWithDefaults) DefaultColumns() []string {
	return []string{"testproperty1"}
}

func TestOutputHandler_AdditionalColumnsWithDefaultColumns(t *testing.T) {
	t.Run("DefaultColumns_OnlyDefaults", func(t *testing.T) {
		// Test with a single item that implements DefaultColumnable, not a slice
		data := testModelWithDefaults{"Row1TestValue1", "Row1TestValue2", "Row1TestValue3"}
		o := NewOutputHandler()

		output := test.CatchStdOut(t, func() {
			err := o.Table(&cobra.Command{}, data) // Pass the item directly, not as a slice
			assert.NoError(t, err)
		})

		// With DefaultColumns interface, should only show the default column (testproperty1)
		expected := "┌────────────────┐\n│ TESTPROPERTY 1 │\n├────────────────┤\n│ Row1TestValue1 │\n└────────────────┘\n"
		assert.Equal(t, expected, output)
	})

	t.Run("DefaultColumns_WithAdditionalColumns", func(t *testing.T) {
		data := testModelWithDefaults{"Row1TestValue1", "Row1TestValue2", "Row1TestValue3"}
		o := NewOutputHandler(WithAdditionalColumns("testproperty2", "testproperty3"))

		output := test.CatchStdOut(t, func() {
			err := o.Table(&cobra.Command{}, data) // Pass the item directly, not as a slice
			assert.NoError(t, err)
		})

		// With additional columns, should show default columns + additional columns
		expected := "┌────────────────┬────────────────┬────────────────┐\n│ TESTPROPERTY 1 │ TESTPROPERTY 2 │ TESTPROPERTY 3 │\n├────────────────┼────────────────┼────────────────┤\n│ Row1TestValue1 │ Row1TestValue2 │ Row1TestValue3 │\n└────────────────┴────────────────┴────────────────┘\n"
		assert.Equal(t, expected, output)
	})

	t.Run("SliceOfDefaultColumns_WithAdditionalColumns", func(t *testing.T) {
		// Test with a slice - DefaultColumns interface won't work, but additional columns should
		data := testModelWithDefaults{"Row1TestValue1", "Row1TestValue2", "Row1TestValue3"}
		o := NewOutputHandler(WithAdditionalColumns("testproperty2", "testproperty3"))

		output := test.CatchStdOut(t, func() {
			err := o.Table(&cobra.Command{}, []testModelWithDefaults{data})
			assert.NoError(t, err)
		})

		// When passing a slice, DefaultColumns interface doesn't work, only additional columns show
		expected := "┌────────────────┬────────────────┐\n│ TESTPROPERTY 2 │ TESTPROPERTY 3 │\n├────────────────┼────────────────┤\n│ Row1TestValue2 │ Row1TestValue3 │\n└────────────────┴────────────────┘\n"
		assert.Equal(t, expected, output)
	})
}

func TestOutputHandler_AdditionalColumnsWithPropertyFilter(t *testing.T) {
	t.Run("PropertyFilter_WithAdditionalColumns", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().StringSlice("property", nil, "")
		cmd.Flags().Set("property", "testproperty2")

		o := NewOutputHandler(WithAdditionalColumns("testproperty1", "testproperty3"))

		output := test.CatchStdOut(t, func() {
			err := o.Table(cmd, collectionSingleRow)
			assert.NoError(t, err)
		})

		// When property flag is set, additional columns are still added to the filter
		// The order follows the original column order in the struct: testproperty1, testproperty2, testproperty3
		expected := "┌────────────────┬────────────────┬────────────────┐\n│ TESTPROPERTY 1 │ TESTPROPERTY 2 │ TESTPROPERTY 3 │\n├────────────────┼────────────────┼────────────────┤\n│ Row1TestValue1 │ Row1TestValue2 │ Row1TestValue3 │\n└────────────────┴────────────────┴────────────────┘\n"
		assert.Equal(t, expected, output)
	})
}
