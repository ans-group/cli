package helper_test

import (
	"testing"

	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/stretchr/testify/assert"
)

func TestJoinInt_EmptyArray_ReturnsEmptyString(t *testing.T) {
	var arr []int

	out := helper.JoinInt(arr, ",")

	assert.Equal(t, "", out)
}

func TestJoinInt_NoneEmptyArray_ReturnsEmptyString(t *testing.T) {
	arr := []int{1, 2, 3}

	out := helper.JoinInt(arr, ",")

	assert.Equal(t, "1,2,3", out)
}
