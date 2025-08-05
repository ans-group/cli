package helper_test

import (
	"testing"

	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestInferTypeFromStringFlagValue_InfersInteger(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		flag := "123"

		value := helper.InferTypeFromStringFlagValue(flag)

		assert.Equal(t, 123, value)
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		flag := "-123"

		value := helper.InferTypeFromStringFlagValue(flag)

		assert.Equal(t, -123, value)
	})
}

func TestInferTypeFromStringFlagValue_InfersBool(t *testing.T) {
	t.Run("Lowercase", func(t *testing.T) {
		flag := "true"

		value := helper.InferTypeFromStringFlagValue(flag)

		assert.Equal(t, true, value)
	})

	t.Run("Uppercase", func(t *testing.T) {
		flag := "TRUE"

		value := helper.InferTypeFromStringFlagValue(flag)

		assert.Equal(t, true, value)
	})
}

func TestInferTypeFromStringFlagValue_InfersString(t *testing.T) {
	flag := "somestring"

	value := helper.InferTypeFromStringFlagValue(flag)

	assert.Equal(t, "somestring", value)
}

func TestGetFilteringInferOperator_Expected(t *testing.T) {
	t.Run("ContainsAsterisk_ReturnsLKOperator", func(t *testing.T) {
		f := helper.GetFilteringInferOperator("testproperty", "testvalue*")

		assert.Equal(t, "testproperty", f.Property)
		assert.Equal(t, connection.LKOperator, f.Operator)
		assert.Equal(t, "testvalue*", f.Value[0])
	})

	t.Run("ContainsComma_ReturnsINOperator", func(t *testing.T) {
		f := helper.GetFilteringInferOperator("testproperty", "testvalue1,testvalue2")

		assert.Equal(t, "testproperty", f.Property)
		assert.Equal(t, connection.INOperator, f.Operator)
		assert.Equal(t, "testvalue1,testvalue2", f.Value[0])
	})

	t.Run("Default_ReturnsEQOperator", func(t *testing.T) {
		f := helper.GetFilteringInferOperator("testproperty", "testvalue1")

		assert.Equal(t, "testproperty", f.Property)
		assert.Equal(t, connection.EQOperator, f.Operator)
		assert.Equal(t, "testvalue1", f.Value[0])
	})
}

func TestGetFilteringFromStringFlagValue_Expected(t *testing.T) {
	t.Run("SingleValue", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlagValue("testproperty:eq=value")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.EQOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value", filtering.Value[0])
	})

	t.Run("CommaSeparated", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlagValue("testproperty:in=value1,value2")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.INOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 2)
		assert.Equal(t, "value1", filtering.Value[0])
		assert.Equal(t, "value2", filtering.Value[1])
	})

	t.Run("EmptyFilter_ReturnsEmptyFiltering", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlagValue("")

		assert.Nil(t, err)
		assert.Equal(t, connection.APIRequestFiltering{}, filtering)
	})

	t.Run("MissingOperator_ReturnsEQFilter", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlagValue("testproperty=value")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.EQOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value", filtering.Value[0])
	})

	t.Run("MissingOperatorWithGlob_ReturnsLKFilter", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlagValue("testproperty=value*")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.LKOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value*", filtering.Value[0])
	})

	t.Run("MissingProperty_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlagValue(":eq=value")

		assert.NotNil(t, err)
		assert.Equal(t, "missing property for filtering", err.Error())
	})

	t.Run("MissingOperator_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlagValue("testproperty:=value")

		assert.NotNil(t, err)
		assert.Equal(t, "missing operator for filtering", err.Error())
	})

	t.Run("EmptyValue_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlagValue("testproperty:invalid=")

		assert.NotNil(t, err)
		assert.Equal(t, "missing value for filtering", err.Error())
	})

	t.Run("InvalidOperator_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlagValue("testproperty:invalid=value")

		assert.NotNil(t, err)
		assert.Equal(t, "invalid filtering operator", err.Error())
	})

	t.Run("MissingValue_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlagValue("testproperty:eq")

		assert.NotNil(t, err)
		assert.Equal(t, "missing value for filtering", err.Error())
	})
}

func TestGetFilteringArrayFromStringArrayFlagValue_Expected(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		filtering, err := helper.GetFilteringArrayFromStringArrayFlagValue([]string{"testproperty:eq=value"})

		assert.Nil(t, err)

		assert.Len(t, filtering, 1)
		assert.Equal(t, "testproperty", filtering[0].Property)
		assert.Equal(t, connection.EQOperator, filtering[0].Operator)
		assert.Equal(t, "value", filtering[0].Value[0])
	})

	t.Run("Multiple", func(t *testing.T) {
		filtering, err := helper.GetFilteringArrayFromStringArrayFlagValue([]string{"testproperty1:eq=value1", "testproperty2:lt=value2"})

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
		s := helper.GetSortingFromStringFlagValue("test")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})

	t.Run("WithDesc_Descending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlagValue("test:desc")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, true, s.Descending)
	})

	t.Run("WithDescMixedCase_Descending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlagValue("test:DeSc")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, true, s.Descending)
	})

	t.Run("InvalidOrder_Ascending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlagValue("test:invalid")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})

	t.Run("EmptyOrder_Ascending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlagValue("test:")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})
}

func TestGetAPIRequestParametersFromFlags(t *testing.T) {
	t.Run("SortFlagSpecified_SortHydrated", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("sort", "", "")
		cmd.ParseFlags([]string{"--sort=testproperty"})

		params, err := helper.GetAPIRequestParametersFromFlags(cmd)

		assert.Nil(t, err)
		assert.Equal(t, "testproperty", params.Sorting.Property)
	})

	t.Run("SortFlagSpecifiedWithDescending_SortDescendingTrue", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("sort", "", "")
		cmd.ParseFlags([]string{"--sort=testproperty:desc"})

		params, err := helper.GetAPIRequestParametersFromFlags(cmd)

		assert.Nil(t, err)
		assert.Equal(t, "testproperty", params.Sorting.Property)
		assert.True(t, params.Sorting.Descending)
	})

	t.Run("PageFlagSpecified_PaginationPageHydrated", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Int("page", 0, "")
		cmd.ParseFlags([]string{"--page=3"})

		params, err := helper.GetAPIRequestParametersFromFlags(cmd)

		assert.Nil(t, err)
		assert.Equal(t, 3, params.Pagination.Page)
	})

	t.Run("FilterFlagSpecified_HydratesFilter", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().StringArray("filter", []string{}, "")
		cmd.ParseFlags([]string{"--filter=testproperty=testvalue"})

		params, err := helper.GetAPIRequestParametersFromFlags(cmd)

		assert.Nil(t, err)
		assert.Equal(t, "testproperty", params.Filtering[0].Property)
		assert.Equal(t, "testvalue", params.Filtering[0].Value[0])
		assert.Equal(t, connection.EQOperator, params.Filtering[0].Operator)
	})

	t.Run("OptSpecified_OptApplied", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("name", "", "")
		cmd.ParseFlags([]string{"--name=test"})
		params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))

		assert.Nil(t, err)
		assert.Equal(t, "name", params.Filtering[0].Property)
		assert.Equal(t, "test", params.Filtering[0].Value[0])
		assert.Equal(t, connection.EQOperator, params.Filtering[0].Operator)
	})
}

func TestGetStringPtrFlagIfChanged(t *testing.T) {
	t.Run("FlagChanged_ReturnsValue", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("test", "", "")
		cmd.ParseFlags([]string{"--test=abcd"})

		value := helper.GetStringPtrFlagIfChanged(cmd, "test")

		expected := "abcd"
		assert.NotNil(t, value)
		assert.Equal(t, &expected, value)
	})

	t.Run("FlagUnchanged_ReturnsNil", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().String("test", "", "")

		value := helper.GetStringPtrFlagIfChanged(cmd, "test")

		assert.Nil(t, value)
	})
}

func TestGetBoolPtrFlagIfChanged(t *testing.T) {
	t.Run("FlagChanged_ReturnsValue", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("test", false, "")
		cmd.ParseFlags([]string{"--test"})

		value := helper.GetBoolPtrFlagIfChanged(cmd, "test")

		assert.NotNil(t, value)
		assert.True(t, *value)
	})

	t.Run("FlagUnchanged_ReturnsNil", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Bool("test", false, "")

		value := helper.GetBoolPtrFlagIfChanged(cmd, "test")

		assert.Nil(t, value)
	})
}

func TestGetIntPtrFlagIfChanged(t *testing.T) {
	t.Run("FlagChanged_ReturnsValue", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Int("test", 0, "")
		cmd.ParseFlags([]string{"--test=123"})

		value := helper.GetIntPtrFlagIfChanged(cmd, "test")

		expected := 123
		assert.NotNil(t, value)
		assert.Equal(t, &expected, value)
	})

	t.Run("FlagUnchanged_ReturnsNil", func(t *testing.T) {
		cmd := &cobra.Command{}
		cmd.Flags().Int("test", 0, "")

		value := helper.GetIntPtrFlagIfChanged(cmd, "test")

		assert.Nil(t, value)
	})
}
