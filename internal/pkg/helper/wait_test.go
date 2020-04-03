package helper

import (
	"errors"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/ukfast/cli/test"
)

func TestWaitForCommand(t *testing.T) {

	t.Run("SuccessAfter3Attempts", func(t *testing.T) {
		test.TestResetViper()
		defer test.TestResetViper()
		viper.SetDefault("command_wait_sleep_seconds", 1)

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
		test.TestResetViper()
		defer test.TestResetViper()
		viper.SetDefault("command_wait_sleep_seconds", 1)

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

	t.Run("InvalidWaitTimeout", func(t *testing.T) {
		test.TestResetViper()
		defer test.TestResetViper()
		viper.SetDefault("command_wait_timeout_seconds", 0)

		attempt := 1
		f := func() (bool, error) {
			attempt++
			return false, nil
		}

		r := WaitForCommand(f)

		assert.NotNil(t, r)
		assert.Equal(t, 1, attempt)
	})

	t.Run("InvalidWaitSleep", func(t *testing.T) {
		test.TestResetViper()
		defer test.TestResetViper()
		viper.SetDefault("command_wait_sleep_seconds", 0)

		attempt := 1
		f := func() (bool, error) {
			attempt++
			return false, nil
		}

		r := WaitForCommand(f)

		assert.NotNil(t, r)
		assert.Equal(t, 1, attempt)
	})
}
