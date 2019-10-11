package helper_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/internal/pkg/helper"
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
