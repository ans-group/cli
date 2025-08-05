package config

import (
	"testing"

	"github.com/ans-group/cli/test/test_output"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_configContextUpdate(t *testing.T) {
	t.Run("UpdatesSpecifiedContext", func(t *testing.T) {
		defer config.Reset()
		config.Init("")

		fs := afero.NewMemMapFs()
		cmd := configContextUpdateCmd(fs)
		cmd.Flags().Set("api-key", "someapikey")

		err := configContextUpdate(fs, cmd, []string{"somecontext"})

		config.SwitchCurrentContext("somecontext")

		apiKey := config.GetString("api_key")

		assert.Nil(t, err)
		assert.Equal(t, "someapikey", apiKey)
	})

	t.Run("UpdatesCurrentContext", func(t *testing.T) {
		defer config.Reset()
		fs := afero.NewMemMapFs()
		config.SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(`contexts:
  testcontext1:
    api_key: oldvalue
current_context: testcontext1
`), 0644)
		config.Init("/tmp/testconfig.yml")

		cmd := configContextUpdateCmd(fs)
		cmd.Flags().Set("api-key", "someapikey")
		cmd.Flags().Set("current", "true")

		config.SwitchCurrentContext("testcontext1")

		err := configContextUpdate(fs, cmd, []string{""})

		apiKey := config.GetString("api_key")

		assert.Nil(t, err)
		assert.Equal(t, "someapikey", apiKey)
	})
}

func Test_configContextList(t *testing.T) {
	t.Run("ListsContexts", func(t *testing.T) {
		defer config.Reset()
		fs := afero.NewMemMapFs()
		config.SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(`contexts:
  testcontext1:
    api_key: testkey
  testcontext2:
    api_key: testkey2
current_context: testcontext1
`), 0644)
		config.Init("/tmp/testconfig.yml")

		cmd := configContextListCmd()
		cmd.Flags().String("output", "value", "")

		test_output.AssertCombinedOutput(t, "testcontext1 true\ntestcontext2 false\n", "", func() {
			configContextList(cmd)
		})
	})
}

func Test_configContextShow(t *testing.T) {
	t.Run("ShowsCurrentContext", func(t *testing.T) {
		defer config.Reset()
		fs := afero.NewMemMapFs()
		config.SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(`contexts:
  testcontext1:
    api_key: testkey
  testcontext2:
    api_key: testkey2
current_context: testcontext1
`), 0644)
		config.Init("/tmp/testconfig.yml")

		cmd := configContextShowCmd()
		cmd.Flags().String("output", "value", "")

		test_output.AssertCombinedOutput(t, "testcontext1 true\n", "", func() {
			configContextShow(cmd)
		})
	})
}

func Test_configContextSwitchCmd_Args(t *testing.T) {
	t.Run("ValidArgs_NoError", func(t *testing.T) {
		err := configContextSwitchCmd(nil).Args(nil, []string{"testcontext"})

		assert.Nil(t, err)
	})

	t.Run("InvalidArgs_Error", func(t *testing.T) {
		err := configContextSwitchCmd(nil).Args(nil, []string{})

		assert.NotNil(t, err)
		assert.Equal(t, "missing context", err.Error())
	})
}

func Test_configContextSwitch(t *testing.T) {
	t.Run("SwitchesContext", func(t *testing.T) {
		defer config.Reset()
		fs := afero.NewMemMapFs()
		config.SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(`contexts:
  testcontext1:
    api_key: testkey
  testcontext2:
    api_key: testkey2
current_context: testcontext1
`), 0644)
		config.Init("/tmp/testconfig.yml")

		err := configContextSwitch(fs, configContextSwitchCmd(fs), []string{"testcontext2"})

		currentContext := config.GetCurrentContextName()

		assert.Nil(t, err)
		assert.Equal(t, "testcontext2", currentContext)
	})

	t.Run("NonExistentContextReturnsError", func(t *testing.T) {
		defer config.Reset()
		fs := afero.NewMemMapFs()
		config.SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(`contexts:
  testcontext1:
    api_key: testkey
  testcontext2:
    api_key: testkey2
current_context: testcontext1
`), 0644)
		config.Init("/tmp/testconfig.yml")

		err := configContextSwitch(fs, configContextSwitchCmd(fs), []string{"nonexistent"})

		assert.NotNil(t, err)
	})
}
