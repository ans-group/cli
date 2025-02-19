package output

import (
	"testing"

	"github.com/ans-group/cli/test"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

type testCollection struct {
	fields []*OrderedFields
}

func (t *testCollection) Fields() []*OrderedFields {
	return t.fields
}

func newTestCollection(fields []*OrderedFields) *testCollection {
	return &testCollection{
		fields: fields,
	}
}

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
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", "TestValue1")
		fields.Set("testproperty2", "TestValue2")
		fields.Set("testproperty3", "TestValue3")
		rows = append(rows, fields)

		collection := newTestCollection(rows)
		o := NewOutputHandler()

		output := test.CatchStdOut(t, func() {
			err := o.Value(&cobra.Command{}, collection)
			assert.NoError(t, err)
		})

		assert.Equal(t, "TestValue1 TestValue2 TestValue3\n", output)
	})

	t.Run("MultipleRowsDefaultFields_ExpectedStdout", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", "Row1TestValue1")
		row1fields.Set("testproperty2", "Row1TestValue2")
		row1fields.Set("testproperty3", "Row1TestValue3")
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", "Row2TestValue1")
		row2fields.Set("testproperty2", "Row2TestValue2")
		row2fields.Set("testproperty3", "Row2TestValue3")
		rows = append(rows, row2fields)

		collection := newTestCollection(rows)
		o := NewOutputHandler()
		output := test.CatchStdOut(t, func() {
			err := o.Value(&cobra.Command{}, collection)
			assert.NoError(t, err)
		})

		assert.Equal(t, "Row1TestValue1 Row1TestValue2 Row1TestValue3\nRow2TestValue1 Row2TestValue2 Row2TestValue3\n", output)
	})
}

func TestOutputHandler_CSV(t *testing.T) {
	t.Run("SingleRowDefaultFields_ExpectedStdout", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", "TestValue1")
		fields.Set("testproperty2", "TestValue2")
		fields.Set("testproperty3", "TestValue3")
		rows = append(rows, fields)

		collection := newTestCollection(rows)
		o := NewOutputHandler()

		output := test.CatchStdOut(t, func() {
			o.CSV(&cobra.Command{}, collection)
		})

		assert.Equal(t, "testproperty1,testproperty2,testproperty3\nTestValue1,TestValue2,TestValue3\n", output)
	})

	t.Run("MultipleRowsDefaultFields_ExpectedStdout", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", "Row1TestValue1")
		row1fields.Set("testproperty2", "Row1TestValue2")
		row1fields.Set("testproperty3", "Row1TestValue3")
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", "Row2TestValue1")
		row2fields.Set("testproperty2", "Row2TestValue2")
		row2fields.Set("testproperty3", "Row2TestValue3")
		rows = append(rows, row2fields)

		collection := newTestCollection(rows)
		o := NewOutputHandler()

		output := test.CatchStdOut(t, func() {
			o.CSV(&cobra.Command{}, collection)
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
	var rows []*OrderedFields
	fields := NewOrderedFields()
	fields.Set("testproperty1", "TestValue1")
	fields.Set("testproperty2", "TestValue2")
	fields.Set("testproperty3", "TestValue3")
	rows = append(rows, fields)

	collection := newTestCollection(rows)
	o := NewOutputHandler()

	output := test.CatchStdOut(t, func() {
		err := o.Table(&cobra.Command{}, collection)
		assert.NoError(t, err)
	})

	assert.Equal(t, "+---------------+---------------+---------------+\n| TESTPROPERTY1 | TESTPROPERTY2 | TESTPROPERTY3 |\n+---------------+---------------+---------------+\n| TestValue1    | TestValue2    | TestValue3    |\n+---------------+---------------+---------------+\n", output)
}

func TestOutputHandler_List(t *testing.T) {
	var rows []*OrderedFields
	fields := NewOrderedFields()
	fields.Set("testproperty1", "TestValue1")
	fields.Set("testproperty2", "TestValue2")
	fields.Set("testproperty3", "TestValue3")
	rows = append(rows, fields)

	collection := newTestCollection(rows)
	o := NewOutputHandler()

	output := test.CatchStdOut(t, func() {
		err := o.List(&cobra.Command{}, collection)
		assert.NoError(t, err)
	})

	assert.Equal(t, "testproperty1 : TestValue1\ntestproperty2 : TestValue2\ntestproperty3 : TestValue3\n", output)
}
