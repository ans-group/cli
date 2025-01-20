package output_test

import (
	"testing"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/stretchr/testify/assert"
)

func TestOrderedFields_Set(t *testing.T) {
	t.Run("SetValue", func(t *testing.T) {
		f := output.NewOrderedFields()
		f.Set("testkey", output.NewFieldValue("testvalue", true))

		v := f.Get("testkey")

		assert.Equal(t, "testvalue", v.Value)
		assert.Equal(t, true, v.Default)
	})

	t.Run("SetExistingValueOverwrite", func(t *testing.T) {
		f := output.NewOrderedFields()
		f.Set("testkey", output.NewFieldValue("testvalue1", false))
		f.Set("testkey", output.NewFieldValue("testvalue2", true))

		v := f.Get("testkey")

		assert.Equal(t, "testvalue2", v.Value)
		assert.Equal(t, true, v.Default)
	})

	t.Run("KeysPopulated", func(t *testing.T) {
		f := output.NewOrderedFields()
		f.Set("testkey", output.NewFieldValue("testvalue", true))

		keys := f.Keys()

		assert.Contains(t, keys, "testkey")
	})

	t.Run("ExistsTrue", func(t *testing.T) {
		f := output.NewOrderedFields()
		f.Set("testkey", output.NewFieldValue("testvalue", true))

		exists := f.Exists("testkey")

		assert.True(t, exists)
	})

	t.Run("ExistsFalse", func(t *testing.T) {
		f := output.NewOrderedFields()
		f.Set("testkey", output.NewFieldValue("testvalue", true))

		exists := f.Exists("testkey2")

		assert.False(t, exists)
	})

	t.Run("NonExistentReturnsEmptyValue", func(t *testing.T) {
		f := output.NewOrderedFields()
		f.Set("testkey", output.NewFieldValue("testvalue", true))

		v := f.Get("testkey2")

		assert.Equal(t, "", v.Value)
	})
}
