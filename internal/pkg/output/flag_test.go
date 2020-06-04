package output

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseOutputFlag_EmptyString_ReturnsEmptyStrings(t *testing.T) {
	flag := ""

	name, arg := ParseOutputFlag(flag)

	assert.Equal(t, "", name)
	assert.Equal(t, "", arg)
}

func TestParseOutputFlag_NameWithNoArg_ReturnsName(t *testing.T) {
	flag := "something"

	name, arg := ParseOutputFlag(flag)

	assert.Equal(t, "something", name)
	assert.Equal(t, "", arg)
}

func TestParseOutputFlag_NameWithArg_ReturnsNameAndArg(t *testing.T) {
	flag := "something=somethingelse"

	name, arg := ParseOutputFlag(flag)

	assert.Equal(t, "something", name)
	assert.Equal(t, "somethingelse", arg)
}

func TestParseOutputFlag_NameWithEmptyArg_ReturnsNameAndEmptyArg(t *testing.T) {
	flag := "something="

	name, arg := ParseOutputFlag(flag)

	assert.Equal(t, "something", name)
	assert.Equal(t, "", arg)
}

func TestParseOutputFlag_NameWithArgContainingSplitChar_ReturnsNameAndArg(t *testing.T) {
	flag := "something=somethingelse=123"

	name, arg := ParseOutputFlag(flag)

	assert.Equal(t, "something", name)
	assert.Equal(t, "somethingelse=123", arg)
}

func TestParseOutputFlag_NameWithSingleQuotedArg_ReturnsNameAndArg(t *testing.T) {
	flag := "something='somethingelse'"

	name, arg := ParseOutputFlag(flag)

	assert.Equal(t, "something", name)
	assert.Equal(t, "somethingelse", arg)
}

func TestParseOutputFlag_NameWithDoubleQuotedArg_ReturnsNameAndArg(t *testing.T) {
	flag := "something=\"somethingelse\""

	name, arg := ParseOutputFlag(flag)

	assert.Equal(t, "something", name)
	assert.Equal(t, "somethingelse", arg)
}
