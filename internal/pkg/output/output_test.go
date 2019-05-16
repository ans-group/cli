package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func TestOutputExit_SetsHandler(t *testing.T) {
	code := 0
	oldOutputExit := SetOutputExit(func(c int) {
		code = c
	})
	defer func() { SetOutputExit(oldOutputExit) }()

	outputExit(1)

	assert.Equal(t, 1, code)
}

func TestOutputExit_ReturnsOldExitCode(t *testing.T) {
	code := 0
	oldOutputExit := SetOutputExit(func(c int) { code = c })
	defer func() { SetOutputExit(oldOutputExit) }()

	SetOutputExit(func(c int) {})(1)

	assert.Equal(t, 1, code)
}

func TestDebugLogger_Error_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Error("testerror")
	})

	assert.Equal(t, "testerror\n", output)
}

func TestDebugLogger_Warn_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Warn("testwarn")
	})

	assert.Equal(t, "testwarn\n", output)
}

func TestDebugLogger_Info_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Info("testinfo")
	})

	assert.Equal(t, "testinfo\n", output)
}

func TestDebugLogger_Debug_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Debug("testdebug")
	})

	assert.Equal(t, "testdebug\n", output)
}

func TestDebugLogger_Trace_OutputsError(t *testing.T) {
	l := DebugLogger{}

	output := test.CatchStdErr(t, func() {
		l.Trace("testtrace")
	})

	assert.Equal(t, "testtrace\n", output)
}

func TestError_ExpectedStderr(t *testing.T) {
	output := test.CatchStdErr(t, func() {
		Error("test")
	})

	assert.Equal(t, "test\n", output)
}

func TestErrorf_ExpectedStderr(t *testing.T) {
	output := test.CatchStdErr(t, func() {
		Errorf("test %s", "fmt")
	})

	assert.Equal(t, "test fmt\n", output)
}

func TestFatal_ExpectedExitCode(t *testing.T) {
	code := 0
	oldOutputExit := SetOutputExit(func(c int) {
		code = c
	})
	defer func() { outputExit = oldOutputExit }()

	output := test.CatchStdErr(t, func() {
		Fatal("test")
	})

	assert.Equal(t, "test\n", output)
	assert.Equal(t, 1, code)
}

func TestFatalf_ExpectedExitCode(t *testing.T) {
	code := 0
	oldOutputExit := SetOutputExit(func(c int) {
		code = c
	})
	defer func() { outputExit = oldOutputExit }()

	output := test.CatchStdErr(t, func() {
		Fatalf("test %s", "fmt")
	})

	assert.Equal(t, "test fmt\n", output)
	assert.Equal(t, 1, code)
}

func TestValue_ExpectedStdout(t *testing.T) {
	t.Run("SingleRowDefaultFields", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Default: true,
			Value:   "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Default: false,
			Value:   "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Default: true,
			Value:   "TestValue3",
		})
		rows = append(rows, fields)

		output := test.CatchStdOut(t, func() {
			Value([]string{}, rows)
		})

		assert.Equal(t, "TestValue1 TestValue3\n", output)
	})

	t.Run("MultipleRowsDefaultFields", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Default: true,
			Value:   "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Default: false,
			Value:   "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Default: true,
			Value:   "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		output := test.CatchStdOut(t, func() {
			Value([]string{}, rows)
		})

		assert.Equal(t, "Row1TestValue1 Row1TestValue3\nRow2TestValue1 Row2TestValue3\n", output)
	})

	t.Run("SingleRowIncludeSingleColumn", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Value: "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Value: "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Value: "TestValue3",
		})
		rows = append(rows, fields)

		includeColumns := []string{"testproperty1"}

		output := test.CatchStdOut(t, func() {
			Value(includeColumns, rows)
		})

		assert.Equal(t, "TestValue1\n", output)
	})

	t.Run("SingleRowIncludeMultipleColumns", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Value: "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Value: "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Value: "TestValue3",
		})
		rows = append(rows, fields)

		includeColumns := []string{"testproperty1", "testproperty3"}

		output := test.CatchStdOut(t, func() {
			Value(includeColumns, rows)
		})

		assert.Equal(t, "TestValue1 TestValue3\n", output)
	})

	t.Run("MultipleRowsIncludeSingleColumn", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Value: "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Value: "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Value: "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		includeColumns := []string{"testproperty1"}

		output := test.CatchStdOut(t, func() {
			Value(includeColumns, rows)
		})

		assert.Equal(t, "Row1TestValue1\nRow2TestValue1\n", output)
	})

	t.Run("MultipleRowsIncludeMultipleColumns", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Value: "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Value: "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Value: "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		includeColumns := []string{"testproperty1", "testproperty3"}

		output := test.CatchStdOut(t, func() {
			Value(includeColumns, rows)
		})

		assert.Equal(t, "Row1TestValue1 Row1TestValue3\nRow2TestValue1 Row2TestValue3\n", output)
	})

	t.Run("SingleRowIncludeColumnGlob", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Value: "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Value: "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Value: "TestValue3",
		})
		rows = append(rows, fields)

		includeColumns := []string{"*"}

		output := test.CatchStdOut(t, func() {
			Value(includeColumns, rows)
		})

		assert.Equal(t, "TestValue1 TestValue2 TestValue3\n", output)
	})

	t.Run("MultipleRowsIncludeColumnGlob", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Value: "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Value: "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Value: "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		includeColumns := []string{"*"}

		output := test.CatchStdOut(t, func() {
			Value(includeColumns, rows)
		})

		assert.Equal(t, "Row1TestValue1 Row1TestValue2 Row1TestValue3\nRow2TestValue1 Row2TestValue2 Row2TestValue3\n", output)
	})
}

func TestOutput_JSON_ExpectedStdout(t *testing.T) {
	t.Run("ExpectedStdOut", func(t *testing.T) {
		zone := safedns.Zone{Name: "testzone.com", Description: "testdescription"}

		output := test.CatchStdOut(t, func() {
			JSON(zone)
		})

		assert.Equal(t, "{\"name\":\"testzone.com\",\"description\":\"testdescription\"}", output)
	})

	t.Run("MarshalError_ReturnsError", func(t *testing.T) {
		type teststruct struct {
			Invalid chan int
		}

		err := JSON(teststruct{})

		assert.NotNil(t, err)
		assert.Equal(t, "failed to marshal json: json: unsupported type: chan int", err.Error())
	})
}

func TestTable_ExpectedStdout(t *testing.T) {
	t.Run("SingleRowDefaultFields", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Default: true,
			Value:   "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Default: false,
			Value:   "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Default: true,
			Value:   "TestValue3",
		})
		rows = append(rows, fields)

		output := test.CatchStdOut(t, func() {
			Table([]string{}, rows)
		})

		assert.Equal(t, "+---------------+---------------+\n| TESTPROPERTY1 | TESTPROPERTY3 |\n+---------------+---------------+\n| TestValue1    | TestValue3    |\n+---------------+---------------+\n", output)
	})

	t.Run("MultipleRowsDefaultFields", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Default: true,
			Value:   "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Default: false,
			Value:   "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Default: true,
			Value:   "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		output := test.CatchStdOut(t, func() {
			Table([]string{}, rows)
		})

		assert.Equal(t, "+----------------+----------------+\n| TESTPROPERTY1  | TESTPROPERTY3  |\n+----------------+----------------+\n| Row1TestValue1 | Row1TestValue3 |\n| Row2TestValue1 | Row2TestValue3 |\n+----------------+----------------+\n", output)
	})

	t.Run("SingleRowIncludeColumns", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Value: "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Value: "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Value: "TestValue3",
		})
		rows = append(rows, fields)

		includeColumns := []string{"testproperty1", "testproperty3"}

		output := test.CatchStdOut(t, func() {
			Table(includeColumns, rows)
		})

		assert.Equal(t, "+---------------+---------------+\n| TESTPROPERTY1 | TESTPROPERTY3 |\n+---------------+---------------+\n| TestValue1    | TestValue3    |\n+---------------+---------------+\n", output)
	})

	t.Run("MultipleRowsIncludeColumns", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Value: "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Value: "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Value: "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		includeColumns := []string{"testproperty1", "testproperty3"}

		output := test.CatchStdOut(t, func() {
			Table(includeColumns, rows)
		})

		assert.Equal(t, "+----------------+----------------+\n| TESTPROPERTY1  | TESTPROPERTY3  |\n+----------------+----------------+\n| Row1TestValue1 | Row1TestValue3 |\n| Row2TestValue1 | Row2TestValue3 |\n+----------------+----------------+\n", output)
	})

	t.Run("NoRows", func(t *testing.T) {
		var rows []*OrderedFields

		output := test.CatchStdOut(t, func() {
			err := Table([]string{}, rows)
			assert.Nil(t, err)
		})

		assert.Equal(t, "", output)
	})

	t.Run("SingleRowIncludeColumnGlob", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Value: "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Value: "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Value: "TestValue3",
		})
		rows = append(rows, fields)

		includeColumns := []string{"*"}

		output := test.CatchStdOut(t, func() {
			Table(includeColumns, rows)
		})

		assert.Equal(t, "+---------------+---------------+---------------+\n| TESTPROPERTY1 | TESTPROPERTY2 | TESTPROPERTY3 |\n+---------------+---------------+---------------+\n| TestValue1    | TestValue2    | TestValue3    |\n+---------------+---------------+---------------+\n", output)
	})

	t.Run("MultipleRowsIncludeColumnGlob", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Value: "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Value: "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Value: "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		includeColumns := []string{"*"}

		output := test.CatchStdOut(t, func() {
			Table(includeColumns, rows)
		})

		assert.Equal(t, "+----------------+----------------+----------------+\n| TESTPROPERTY1  | TESTPROPERTY2  | TESTPROPERTY3  |\n+----------------+----------------+----------------+\n| Row1TestValue1 | Row1TestValue2 | Row1TestValue3 |\n| Row2TestValue1 | Row2TestValue2 | Row2TestValue3 |\n+----------------+----------------+----------------+\n", output)
	})
}

func TestCSV_ExpectedStdout(t *testing.T) {
	t.Run("SingleRowDefaultFields", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Default: true,
			Value:   "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Default: false,
			Value:   "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Default: true,
			Value:   "TestValue3",
		})
		rows = append(rows, fields)

		output := test.CatchStdOut(t, func() {
			CSV([]string{}, rows)
		})

		assert.Equal(t, "testproperty1,testproperty3\nTestValue1,TestValue3\n", output)
	})

	t.Run("MultipleRowsDefaultFields", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Default: true,
			Value:   "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Default: false,
			Value:   "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Default: true,
			Value:   "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		output := test.CatchStdOut(t, func() {
			CSV([]string{}, rows)
		})

		assert.Equal(t, "testproperty1,testproperty3\nRow1TestValue1,Row1TestValue3\nRow2TestValue1,Row2TestValue3\n", output)
	})

	t.Run("SingleRowIncludeColumns", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Value: "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Value: "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Value: "TestValue3",
		})
		rows = append(rows, fields)

		includeColumns := []string{"testproperty1", "testproperty3"}

		output := test.CatchStdOut(t, func() {
			CSV(includeColumns, rows)
		})

		assert.Equal(t, "testproperty1,testproperty3\nTestValue1,TestValue3\n", output)
	})

	t.Run("MultipleRowsIncludeColumns", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Value: "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Value: "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Value: "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		includeColumns := []string{"testproperty1", "testproperty3"}

		output := test.CatchStdOut(t, func() {
			CSV(includeColumns, rows)
		})

		assert.Equal(t, "testproperty1,testproperty3\nRow1TestValue1,Row1TestValue3\nRow2TestValue1,Row2TestValue3\n", output)
	})

	t.Run("SingleRowIncludeColumnGlob", func(t *testing.T) {
		var rows []*OrderedFields
		fields := NewOrderedFields()
		fields.Set("testproperty1", FieldValue{
			Value: "TestValue1",
		})
		fields.Set("testproperty2", FieldValue{
			Value: "TestValue2",
		})
		fields.Set("testproperty3", FieldValue{
			Value: "TestValue3",
		})
		rows = append(rows, fields)

		includeColumns := []string{"*"}

		output := test.CatchStdOut(t, func() {
			CSV(includeColumns, rows)
		})

		assert.Equal(t, "testproperty1,testproperty2,testproperty3\nTestValue1,TestValue2,TestValue3\n", output)
	})

	t.Run("MultipleRowsIncludeColumnGlob", func(t *testing.T) {
		var rows []*OrderedFields

		row1fields := NewOrderedFields()
		row1fields.Set("testproperty1", FieldValue{
			Value: "Row1TestValue1",
		})
		row1fields.Set("testproperty2", FieldValue{
			Value: "Row1TestValue2",
		})
		row1fields.Set("testproperty3", FieldValue{
			Value: "Row1TestValue3",
		})
		rows = append(rows, row1fields)

		row2fields := NewOrderedFields()
		row2fields.Set("testproperty1", FieldValue{
			Value: "Row2TestValue1",
		})
		row2fields.Set("testproperty2", FieldValue{
			Value: "Row2TestValue2",
		})
		row2fields.Set("testproperty3", FieldValue{
			Value: "Row2TestValue3",
		})
		rows = append(rows, row2fields)

		includeColumns := []string{"*"}

		output := test.CatchStdOut(t, func() {
			CSV(includeColumns, rows)
		})

		assert.Equal(t, "testproperty1,testproperty2,testproperty3\nRow1TestValue1,Row1TestValue2,Row1TestValue3\nRow2TestValue1,Row2TestValue2,Row2TestValue3\n", output)
	})
}

func TestTemplate_ExpectedStdout(t *testing.T) {
	t.Run("InvalidTemplateError", func(t *testing.T) {
		template := "{{if invalid}}"
		type model struct {
			TestProperty1 string
		}

		err := Template(template, model{TestProperty1: "testvalue1"})

		assert.NotNil(t, err)
	})

	t.Run("ItemMissingPropertyTemplateError", func(t *testing.T) {
		template := "{{ .TestProperty2 }}"
		type model struct {
			TestProperty1 string
		}

		err := Template(template, model{TestProperty1: "testvalue1"})

		assert.NotNil(t, err)
	})

	t.Run("ItemSliceMissingPropertyTemplateError", func(t *testing.T) {
		template := "{{ .TestProperty2 }}"
		type model struct {
			TestProperty1 string
		}

		var models []model
		models = append(models, model{TestProperty1: "testvalue1"})
		models = append(models, model{TestProperty1: "testvalue1"})

		err := Template(template, models)

		assert.NotNil(t, err)
	})

	t.Run("Item", func(t *testing.T) {
		template := "{{ .TestProperty1 }}"
		type model struct {
			TestProperty1 string
		}

		output := test.CatchStdOut(t, func() {
			Template(template, model{TestProperty1: "testvalue1"})
		})

		assert.Equal(t, "testvalue1\n", output)
	})

	t.Run("ItemSlice", func(t *testing.T) {
		template := "{{ .TestProperty1 }}"
		type model struct {
			TestProperty1 string
		}

		var models []model
		models = append(models, model{TestProperty1: "testvalue1"})
		models = append(models, model{TestProperty1: "testvalue2"})

		output := test.CatchStdOut(t, func() {
			Template(template, models)
		})

		assert.Equal(t, "testvalue1\ntestvalue2\n", output)
	})
}

func TestOrderedFields_Set(t *testing.T) {
	t.Run("SetValue", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey", NewFieldValue("testvalue", true))

		v := f.Get("testkey")

		assert.Equal(t, "testvalue", v.Value)
		assert.Equal(t, true, v.Default)
	})

	t.Run("SetExistingValueOverwrite", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey", NewFieldValue("testvalue1", false))
		f.Set("testkey", NewFieldValue("testvalue2", true))

		v := f.Get("testkey")

		assert.Equal(t, "testvalue2", v.Value)
		assert.Equal(t, true, v.Default)
	})

	t.Run("KeysPopulated", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey", NewFieldValue("testvalue", true))

		keys := f.Keys()

		assert.Contains(t, keys, "testkey")
	})

	t.Run("ExistsTrue", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey", NewFieldValue("testvalue", true))

		exists := f.Exists("testkey")

		assert.True(t, exists)
	})

	t.Run("ExistsFalse", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey", NewFieldValue("testvalue", true))

		exists := f.Exists("testkey2")

		assert.False(t, exists)
	})

	t.Run("NonExistentReturnsEmptyValue", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey", NewFieldValue("testvalue", true))

		v := f.Get("testkey2")

		assert.Equal(t, "", v.Value)
	})
}

func Test_getColumnsOrDefault_Expected(t *testing.T) {
	t.Run("NoIncludeFields_DefaultColumns", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey1", NewFieldValue("testvalue1", false))
		f.Set("testkey2", NewFieldValue("testvalue2", true))
		f.Set("testkey3", NewFieldValue("testvalue3", false))

		columns := getColumnsOrDefault([]string{}, f)

		assert.Len(t, columns, 1)
		assert.Equal(t, "testkey2", columns[0])
	})

	t.Run("IncludeField_ExpectedColumn", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey1", NewFieldValue("testvalue1", true))
		f.Set("testkey2", NewFieldValue("testvalue2", false))
		f.Set("testkey3", NewFieldValue("testvalue3", false))

		columns := getColumnsOrDefault([]string{"testkey2"}, f)

		assert.Len(t, columns, 1)
		assert.Equal(t, "testkey2", columns[0])
	})

	t.Run("IncludeFieldCaseInsensitive_ExpectedColumn", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey1", NewFieldValue("testvalue1", true))
		f.Set("testkey2", NewFieldValue("testvalue2", false))
		f.Set("testkey3", NewFieldValue("testvalue3", false))

		columns := getColumnsOrDefault([]string{"TeStKeY2"}, f)

		assert.Len(t, columns, 1)
		assert.Equal(t, "testkey2", columns[0])
	})

	t.Run("IncludeMultipleField_ExpectedColumns", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey1", NewFieldValue("testvalue1", true))
		f.Set("testkey2", NewFieldValue("testvalue2", false))
		f.Set("testkey3", NewFieldValue("testvalue3", false))

		columns := getColumnsOrDefault([]string{"testkey1", "testkey2"}, f)

		assert.Len(t, columns, 2)
		assert.Equal(t, "testkey1", columns[0])
		assert.Equal(t, "testkey2", columns[1])
	})

	t.Run("GlobAll_ExpectedColumns", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey1", NewFieldValue("testvalue1", true))
		f.Set("testkey2", NewFieldValue("testvalue2", false))
		f.Set("testkey3", NewFieldValue("testvalue3", false))

		columns := getColumnsOrDefault([]string{"*"}, f)

		assert.Len(t, columns, 3)
		assert.Equal(t, "testkey1", columns[0])
		assert.Equal(t, "testkey2", columns[1])
		assert.Equal(t, "testkey3", columns[2])
	})

	t.Run("GlobStart_ExpectedColumns", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey1", NewFieldValue("testvalue1", true))
		f.Set("testkey2", NewFieldValue("testvalue2", false))
		f.Set("testkey3", NewFieldValue("testvalue3", false))
		f.Set("otherkey", NewFieldValue("othervalue", false))

		columns := getColumnsOrDefault([]string{"*key1"}, f)

		assert.Len(t, columns, 1)
		assert.Equal(t, "testkey1", columns[0])
	})

	t.Run("GlobEnd_ExpectedColumns", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey1", NewFieldValue("testvalue1", true))
		f.Set("testkey2", NewFieldValue("testvalue2", false))
		f.Set("testkey3", NewFieldValue("testvalue3", false))
		f.Set("otherkey", NewFieldValue("othervalue", false))

		columns := getColumnsOrDefault([]string{"testkey*"}, f)

		assert.Len(t, columns, 3)
		assert.Equal(t, "testkey1", columns[0])
		assert.Equal(t, "testkey2", columns[1])
		assert.Equal(t, "testkey3", columns[2])
		f.Set("otherkey", NewFieldValue("othervalue", false))
	})

	t.Run("GlobStartEnd_ExpectedColumns", func(t *testing.T) {
		f := NewOrderedFields()
		f.Set("testkey1", NewFieldValue("testvalue1", true))
		f.Set("testkey2", NewFieldValue("testvalue2", false))
		f.Set("testkey3", NewFieldValue("testvalue3", false))
		f.Set("otherkey", NewFieldValue("othervalue", false))

		columns := getColumnsOrDefault([]string{"*test*"}, f)

		assert.Len(t, columns, 3)
		assert.Equal(t, "testkey1", columns[0])
		assert.Equal(t, "testkey2", columns[1])
		assert.Equal(t, "testkey3", columns[2])
	})
}
