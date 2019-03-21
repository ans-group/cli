package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/sdk-go/pkg/connection"
)

func TestInferTypeFromStringFlag_InfersInteger(t *testing.T) {
	t.Run("PositiveNumber", func(t *testing.T) {
		flag := "123"

		value := helper.InferTypeFromStringFlag(flag)

		assert.Equal(t, 123, value)
	})

	t.Run("NegativeNumber", func(t *testing.T) {
		flag := "-123"

		value := helper.InferTypeFromStringFlag(flag)

		assert.Equal(t, -123, value)
	})
}

func TestInferTypeFromStringFlag_InfersBool(t *testing.T) {
	t.Run("Lowercase", func(t *testing.T) {
		flag := "true"

		value := helper.InferTypeFromStringFlag(flag)

		assert.Equal(t, true, value)
	})

	t.Run("Uppercase", func(t *testing.T) {
		flag := "TRUE"

		value := helper.InferTypeFromStringFlag(flag)

		assert.Equal(t, true, value)
	})
}

func TestInferTypeFromStringFlag_InfersString(t *testing.T) {
	flag := "somestring"

	value := helper.InferTypeFromStringFlag(flag)

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

func TestGetFilteringFromStringFlag_Expected(t *testing.T) {
	t.Run("SingleValue", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlag("testproperty:eq=value")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.EQOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value", filtering.Value[0])
	})

	t.Run("CommaSeparated", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlag("testproperty:in=value1,value2")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.INOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 2)
		assert.Equal(t, "value1", filtering.Value[0])
		assert.Equal(t, "value2", filtering.Value[1])
	})

	t.Run("EmptyFilter_ReturnsEmptyFiltering", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlag("")

		assert.Nil(t, err)
		assert.Equal(t, connection.APIRequestFiltering{}, filtering)
	})

	t.Run("MissingOperator_ReturnsEQFilter", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlag("testproperty=value")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.EQOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value", filtering.Value[0])
	})

	t.Run("MissingOperatorWithGlob_ReturnsLKFilter", func(t *testing.T) {
		filtering, err := helper.GetFilteringFromStringFlag("testproperty=value*")

		assert.Nil(t, err)

		assert.Equal(t, "testproperty", filtering.Property)
		assert.Equal(t, connection.LKOperator, filtering.Operator)
		assert.Len(t, filtering.Value, 1)
		assert.Equal(t, "value*", filtering.Value[0])
	})

	t.Run("MissingProperty_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlag(":eq=value")

		assert.NotNil(t, err)
		assert.Equal(t, "Missing property for filtering", err.Error())
	})

	t.Run("MissingOperator_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlag("testproperty:=value")

		assert.NotNil(t, err)
		assert.Equal(t, "Missing operator for filtering", err.Error())
	})

	t.Run("EmptyValue_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlag("testproperty:invalid=")

		assert.NotNil(t, err)
		assert.Equal(t, "Missing value for filtering", err.Error())
	})

	t.Run("InvalidOperator_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlag("testproperty:invalid=value")

		assert.NotNil(t, err)
		assert.Equal(t, "Invalid filtering operator", err.Error())
	})

	t.Run("MissingValue_ReturnsError", func(t *testing.T) {
		_, err := helper.GetFilteringFromStringFlag("testproperty:eq")

		assert.NotNil(t, err)
		assert.Equal(t, "Missing value for filtering", err.Error())
	})
}

func TestGetFilteringFromStringArrayFlag_Expected(t *testing.T) {
	t.Run("Single", func(t *testing.T) {
		filtering, err := helper.GetFilteringArrayFromStringArrayFlag([]string{"testproperty:eq=value"})

		assert.Nil(t, err)

		assert.Len(t, filtering, 1)
		assert.Equal(t, "testproperty", filtering[0].Property)
		assert.Equal(t, connection.EQOperator, filtering[0].Operator)
		assert.Equal(t, "value", filtering[0].Value[0])
	})

	t.Run("Multiple", func(t *testing.T) {
		filtering, err := helper.GetFilteringArrayFromStringArrayFlag([]string{"testproperty1:eq=value1", "testproperty2:lt=value2"})

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

func TestGetSortingFromStringFlag_Expected(t *testing.T) {
	t.Run("Default_Ascending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlag("test")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})

	t.Run("WithDesc_Descending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlag("test:desc")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, true, s.Descending)
	})

	t.Run("WithDescMixedCase_Descending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlag("test:DeSc")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, true, s.Descending)
	})

	t.Run("InvalidOrder_Ascending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlag("test:invalid")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})

	t.Run("EmptyOrder_Ascending", func(t *testing.T) {
		s := helper.GetSortingFromStringFlag("test:")

		assert.Equal(t, "test", s.Property)
		assert.Equal(t, false, s.Descending)
	})
}
