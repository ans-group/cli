package helper

import (
	"errors"
	"testing"

	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestWaitForCommand(t *testing.T) {

	t.Run("SuccessAfter3Attempts", func(t *testing.T) {
		config.Reset()
		config.Set("test", "command_wait_sleep_seconds", 1)
		config.SwitchCurrentContext("test")
		defer config.Reset()

		attempt := 1
		f := func() (bool, error) {
			if attempt == 3 {
				return true, nil
			}

			attempt++
			return false, nil
		}

		r := WaitForCommand(f)

		assert.Nil(t, r)
		assert.Equal(t, 3, attempt)
	})

	t.Run("ErrorAfter3Attempts", func(t *testing.T) {
		config.Reset()
		config.Set("test", "command_wait_sleep_seconds", 1)
		config.SwitchCurrentContext("test")
		defer config.Reset()

		attempt := 1
		f := func() (bool, error) {
			if attempt == 3 {
				return false, errors.New("test error")
			}

			attempt++
			return false, nil
		}

		r := WaitForCommand(f)

		assert.NotNil(t, r)
		assert.Equal(t, 3, attempt)
	})
}
