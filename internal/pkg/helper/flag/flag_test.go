package flag_test

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/sdk-go/pkg/connection"
)

func TestInferTypeFromStringFlagValue_InfersInteger(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		flag := "123"

		value := flaghelper.InferTypeFromStringFlagValue(flag)

		assert.Equal(t, 123, value)
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		flag := "-123"

		value := flaghelper.InferTypeFromStringFlagValue(flag)

		assert.Equal(t, -123, value)
	})
}

func TestInferTypeFromStringFlagValue_InfersBool(t *testing.T) {
	t.Run("Lowercase", func(t *testing.T) {
		flag := "true"

		value := flaghelper.InferTypeFromStringFlagValue(flag)

		assert.Equal(t, true, value)
	})

	t.Run("Uppercase", func(t *testing.T) {
		flag := "TRUE"

		value := flaghelper.InferTypeFromStringFlagValue(flag)

		assert.Equal(t, true, value)
	})
}

func TestInferTypeFromStringFlagValue_InfersString(t *testing.T) {
	flag := "somestring"

	value := flaghelper.InferTypeFromStringFlagValue(flag)

	assert.Equal(t, "somestring", value)
}

func TestGetFilteringInferOperator_Expected(t *testing.T) {
	t.Run("ContainsAsterisk_ReturnsLKOperator", func(t *testing.T) {
		f := flaghelper.GetFilteringInferOperator("testproperty", "testvalue*")

		assert.Equal(t, "testproperty", f.Property)
		assert.Equal(t, connection.LKOperator, f.Operator)
		assert.Equal(t, "testvalue*", f.Value[0])
	})

	t.Run("ContainsComma_ReturnsINOperator", func(t *testing.T) {
		f := flaghelper.GetFilteringInferOperator("testproperty", "testvalue1,testvalue2")

		assert.Equal(t, "testproperty", f.Property)
		assert.Equal(t, connection.INOperator, f.Operator)
		assert.Equal(t, "testvalue1,testvalue2", f.Value[0])
	})

	t.Run("Default_ReturnsEQOperator", func(t *testing.T) {
		f := flaghelper.GetFilteringInferOperator("testproperty", "testvalue1")

		assert.Equal(t, "testproperty", f.Property)
		assert.Equal(t, connection.EQOperator, f.Operator)
		assert.Equal(t, "testvalue1", f.Value[0])
	})
}

func TestGetFilteringFromStringFlagValue_Expected(t *testing.T) {
	t.Run("SingleValue", func(t *testing.T) {
		filtering, err := flaghelper.GetFilteringFromStringFlagValue("testproperty:eq=value")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.EQOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value", filtering.Value[0])
	})

	t.Run("CommaSeparated", func(t *testing.T) {
		filtering, err := flaghelper.GetFilteringFromStringFlagValue("testproperty:in=value1,value2")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.INOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 2)
		assert.Equal(t, "value1", filtering.Value[0])
		assert.Equal(t, "value2", filtering.Value[1])
	})

	t.Run("EmptyFilter_ReturnsEmptyFiltering", func(t *testing.T) {
		filtering, err := flaghelper.GetFilteringFromStringFlagValue("")

		assert.Nil(t, err)
		assert.Equal(t, connection.APIRequestFiltering{}, filtering)
	})

	t.Run("MissingOperator_ReturnsEQFilter", func(t *testing.T) {
		filtering, err := flaghelper.GetFilteringFromStringFlagValue("testproperty=value")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.EQOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value", filtering.Value[0])
	})

	t.Run("MissingOperatorWithGlob_ReturnsLKFilter", func(t *testing.T) {
		filtering, err := flaghelper.GetFilteringFromStringFlagValue("testproperty=value*")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.LKOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value*", filtering.Value[0])
	})

	t.Run("MissingProperty_ReturnsError", func(t *testing.T) {
		_, err := flaghelper.GetFilteringFromStringFlagValue(":eq=value")

		assert.NotNil(t, err)
		assert.Equal(t, "Missing property for filtering", err.Error())
	})

	t.Run("MissingOperator_ReturnsError", func(t *testing.T) {
		_, err := flaghelper.GetFilteringFromStringFlagValue("testproperty:=value")

		assert.NotNil(t, err)
		assert.Equal(t, "Missing operator for filtering", err.Error())
	})

	t.Run("EmptyValue_ReturnsError", func(t *testing.T) {
		_, err := flaghelper.GetFilteringFromStringFlagValue("testproperty:invalid=")

		assert.NotNil(t, err)
		assert.Equal(t, "Missing value for filtering", err.Error())
	})

	t.Run("InvalidOperator_ReturnsError", func(t *testing.T) {
		_, err := flaghelper.GetFilteringFromStringFlagValue("testproperty:invalid=value")

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid filtering operator", err.Error())
	})

	t.Run("MissingValue_ReturnsError", func(t *testing.T) {
		_, err := flaghelper.GetFilteringFromStringFlagValue("testproperty:eq")

		assert.NotNil(t, err)
		assert.Equal(t, "Missing value for filtering", err.Error())
	})
}

func TestGetFilteringArrayFromStringArrayFlagValue_Expected(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		filtering, err := flaghelper.GetFilteringArrayFromStringArrayFlagValue([]string{"testproperty:eq=value"})

		assert.Nil(t, err)

		assert.Len(t, filtering, 1)
		assert.Equal(t, "testproperty", filtering[0].Property)
		assert.Equal(t, connection.EQOperator, filtering[0].Operator)
		assert.Equal(t, "value", filtering[0].Value[0])
	})

	t.Run("Multiple", func(t *testing.T) {
		filtering, err := flaghelper.GetFilteringArrayFromStringArrayFlagValue([]string{"testproperty1:eq=value1", "testproperty2:lt=value2"})

		assert.Nil(t, err)

		assert.Len(t, filtering, 2)
		assert.Equal(t, "testproperty1", filtering[0].Property)
		assert.Equal(t, connection.EQOperator, filtering[0].Operator)
		assert.Equal(t, "value1", filtering[0].Value[0])
		assert.Equal(t, "testproperty2", filtering[1].Property)
		assert.Equal(t, connection.LTOperator, filtering[1].Operator)
		assert.Equal(t, "value2", filtering[1].Value[0])
	})
}

func TestGetSortingFromStringFlagValue_Expected(t *testing.T) {
	t.Run("Default_Ascending", func(t *testing.T) {
		s := flaghelper.GetSortingFromStringFlagValue("test")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})

	t.Run("WithDesc_Descending", func(t *testing.T) {
		s := flaghelper.GetSortingFromStringFlagValue("test:desc")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, true, s.Descending)
	})

	t.Run("WithDescMixedCase_Descending", func(t *testing.T) {
		s := flaghelper.GetSortingFromStringFlagValue("test:DeSc")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, true, s.Descending)
	})

	t.Run("InvalidOrder_Ascending", func(t *testing.T) {
		s := flaghelper.GetSortingFromStringFlagValue("test:invalid")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})

	t.Run("EmptyOrder_Ascending", func(t *testing.T) {
		s := flaghelper.GetSortingFromStringFlagValue("test:")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})
}

func TestHydrateAPIRequestParametersWithStringFilterFlag(t *testing.T) {
	t.Run("FlagNotSpecified_NoFilterHydrated", func(t *testing.T) {
		cmd := &cobra.Command{}
		params := connection.APIRequestParameters{}

		flaghelper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, flaghelper.NewStringFilterFlag("noneexistent", "noneexistent"))

		assert.Len(t, params.Filtering, 0)
	})

	t.Run("FlagSpecified_FilterHydrated", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("name", "", "")
		cmd.ParseFlags([]string{"--name=test"})
		params := connection.APIRequestParameters{}

		flaghelper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, flaghelper.NewStringFilterFlag("name", "name"))

		assert.Len(t, params.Filtering, 1)
		assert.Equal(t, "name", params.Filtering[0].Property)
	})
}

func Test_GetChangedOrDefaultPtrString(t *testing.T) {
	t.Run("FlagChanged_ReturnsValue", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("test-flag", "", "")
		cmd.ParseFlags([]string{"--test-flag=test"})

		val := flaghelper.GetChangedOrDefaultPtrString(cmd, "test-flag")

		assert.Equal(t, "test", *val)
	})

	t.Run("FlagNotChanged_ReturnsNil", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("test-flag", "", "")

		val := flaghelper.GetChangedOrDefaultPtrString(cmd, "test-flag")

		assert.Nil(t, val)
	})
}

func Test_GetChangedOrDefaultPtrBool(t *testing.T) {
	t.Run("FlagChanged_ReturnsValue", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("test-flag", false, "")
		cmd.ParseFlags([]string{"--test-flag=true"})

		val := flaghelper.GetChangedOrDefaultPtrBool(cmd, "test-flag")

		assert.True(t, *val)
	})

	t.Run("FlagNotChanged_ReturnsNil", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("test-flag", false, "")

		val := flaghelper.GetChangedOrDefaultPtrBool(cmd, "test-flag")

		assert.Nil(t, val)
	})
}

func Test_GetChangedOrDefaultPtrInt(t *testing.T) {
	t.Run("FlagChanged_ReturnsValue", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Int("test-flag", 0, "")
		cmd.ParseFlags([]string{"--test-flag=123"})

		val := flaghelper.GetChangedOrDefaultPtrInt(cmd, "test-flag")

		assert.Equal(t, 123, *val)
	})

	t.Run("FlagNotChanged_ReturnsNil", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Int("test-flag", 0, "")

		val := flaghelper.GetChangedOrDefaultPtrInt(cmd, "test-flag")

		assert.Nil(t, val)
	})
}
