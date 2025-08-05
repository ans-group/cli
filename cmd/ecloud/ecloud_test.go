package ecloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetKeyValueFromStringFlag(t *testing.T) {
	t.Run("Valid_NoError", func(t *testing.T) {
		flag := "testkey=testvalue"

		key, value, err := GetKeyValueFromStringFlag(flag)

		assert.Nil(t, err)
		assert.Equal(t, "testkey", key)
		assert.Equal(t, "testvalue", value)
	})

	t.Run("Empty_Error", func(t *testing.T) {
		flag := ""

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})

	t.Run("OnlyKey_Error", func(t *testing.T) {
		flag := "testkey"

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})

	t.Run("MissingValue_Error", func(t *testing.T) {
		flag := "testkey="

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})

	t.Run("MissingKey_Error", func(t *testing.T) {
		flag := "=testvalue"

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})

	t.Run("MultiValue_Error", func(t *testing.T) {
		flag := "testkey=testvalue1=testvalue2"

		_, _, err := GetKeyValueFromStringFlag(flag)

		assert.NotNil(t, err)
	})
}

func TestGetCreateTagRequestFromStringArrayFlag(t *testing.T) {
	t.Run("None_NoError", func(t *testing.T) {
		var tagFlags []string

		r, err := GetCreateTagRequestFromStringArrayFlag(tagFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 0)
	})

	t.Run("Single", func(t *testing.T) {
		var tagFlags []string
		tagFlags = append(tagFlags, "testkey1=testvalue1")

		r, err := GetCreateTagRequestFromStringArrayFlag(tagFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 1)
		assert.Equal(t, "testkey1", r[0].Key)
		assert.Equal(t, "testvalue1", r[0].Value)
	})

	t.Run("Multiple", func(t *testing.T) {
		var tagFlags []string
		tagFlags = append(tagFlags, "testkey1=testvalue1")
		tagFlags = append(tagFlags, "testkey2=testvalue2")

		r, err := GetCreateTagRequestFromStringArrayFlag(tagFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 2)
		assert.Equal(t, "testkey1", r[0].Key)
		assert.Equal(t, "testvalue1", r[0].Value)
		assert.Equal(t, "testkey2", r[1].Key)
		assert.Equal(t, "testvalue2", r[1].Value)
	})

	t.Run("Invalid_ReturnsError", func(t *testing.T) {
		tagFlags := []string{"invalid"}

		_, err := GetCreateTagRequestFromStringArrayFlag(tagFlags)

		assert.NotNil(t, err)
		assert.Equal(t, "invalid format, expecting: key=value", err.Error())
	})
}

func TestGetCreateVirtualMachineRequestParameterFromStringArrayFlag(t *testing.T) {
	t.Run("None_NoError", func(t *testing.T) {
		var parameterFlags []string

		r, err := GetCreateVirtualMachineRequestParameterFromStringArrayFlag(parameterFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 0)
	})

	t.Run("Single", func(t *testing.T) {
		var parameterFlags []string
		parameterFlags = append(parameterFlags, "testkey1=testvalue1")

		r, err := GetCreateVirtualMachineRequestParameterFromStringArrayFlag(parameterFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 1)
		assert.Equal(t, "testkey1", r[0].Key)
		assert.Equal(t, "testvalue1", r[0].Value)
	})

	t.Run("Multiple", func(t *testing.T) {
		var parameterFlags []string
		parameterFlags = append(parameterFlags, "testkey1=testvalue1")
		parameterFlags = append(parameterFlags, "testkey2=testvalue2")

		r, err := GetCreateVirtualMachineRequestParameterFromStringArrayFlag(parameterFlags)

		assert.Nil(t, err)
		assert.Len(t, r, 2)
		assert.Equal(t, "testkey1", r[0].Key)
		assert.Equal(t, "testvalue1", r[0].Value)
		assert.Equal(t, "testkey2", r[1].Key)
		assert.Equal(t, "testvalue2", r[1].Value)
	})

	t.Run("Invalid_ReturnsError", func(t *testing.T) {
		parameterFlags := []string{"invalid"}

		_, err := GetCreateVirtualMachineRequestParameterFromStringArrayFlag(parameterFlags)

		assert.NotNil(t, err)
		assert.Equal(t, "invalid format, expecting: key=value", err.Error())
	})
}
