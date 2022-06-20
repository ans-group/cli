package config

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

var defaultConfig = `contexts:
  testcontext1:
    somekey: somevalue
  testcontext2:
    somekey: somevalue
current_context: testcontext1
`

func TestInit(t *testing.T) {
	t.Run("ReadsConfigFile", func(t *testing.T) {
		defer Reset()
		fs := afero.NewMemMapFs()
		SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(`testconfig:
  somekey: somevalue
`), 0644)

		Init("/tmp/testconfig.yml")
		someKey := GetString("testconfig.somekey")

		assert.Equal(t, "somevalue", someKey)
	})

	t.Run("ReturnsErrorConfigNotFound", func(t *testing.T) {
		defer Reset()
		fs := afero.NewMemMapFs()
		SetFs(fs)

		err := Init("/tmp/testconfig.yml")

		assert.NotNil(t, err)
	})
}

func TestSave(t *testing.T) {
	t.Run("SavesToDefinedConfig", func(t *testing.T) {
		defer Reset()
		fs := afero.NewMemMapFs()
		SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(`contexts:
  somecontext:
    somekey: somevalue
`), 0644)

		Init("/tmp/testconfig.yml")
		Set("somecontext", "somekey", "newvalue")
		Save()

		content, _ := afero.ReadFile(fs, "/tmp/testconfig.yml")

		expected := `contexts:
  somecontext:
    somekey: newvalue
`

		assert.Equal(t, expected, string(content))
	})

	t.Run("SavesToDefaultConfig", func(t *testing.T) {
		defer Reset()
		fs := afero.NewMemMapFs()
		SetFs(fs)

		Init("")
		defaultConfigFile = "/tmp/defaultconfig.yml"
		Set("somecontext", "somekey", "newvalue")
		Save()

		content, _ := afero.ReadFile(fs, defaultConfigFile)

		expected := `contexts:
  somecontext:
    somekey: newvalue
`

		assert.Equal(t, expected, string(content))
	})

	t.Run("ReturnsErrorWhenNotInitialised", func(t *testing.T) {
		defer Reset()
		initialised = false
		err := Save()
		assert.NotNil(t, err)
	})
}

func TestGetCurrentContextName(t *testing.T) {
	defer Reset()
	fs := afero.NewMemMapFs()
	SetFs(fs)

	afero.WriteFile(fs, "/tmp/testconfig.yml", []byte("current_context: testcontext"), 0644)

	Init("/tmp/testconfig.yml")

	currentContext := GetCurrentContextName()

	assert.Equal(t, "testcontext", currentContext)
}

func TestGetContextNames(t *testing.T) {
	defer Reset()
	fs := afero.NewMemMapFs()
	SetFs(fs)

	afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(`contexts:
  testcontext1:
    somekey: somevalue
  testcontext2:
    somekey: somevalue
`), 0644)

	Init("/tmp/testconfig.yml")

	contexts := GetContextNames()

	assert.Contains(t, contexts, "testcontext1")
	assert.Contains(t, contexts, "testcontext2")
}

func TestSetCurrentContext(t *testing.T) {
	t.Run("SetsCurrentContextValue", func(t *testing.T) {
		defer Reset()
		fs := afero.NewMemMapFs()
		SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(defaultConfig), 0644)

		Init("/tmp/testconfig.yml")

		SetCurrentContext("somekey", "somenewvalue")
		value := GetString("somekey")

		assert.Equal(t, "somenewvalue", value)
	})

	t.Run("CurrentContextNotSetReturnsError", func(t *testing.T) {
		defer Reset()
		err := SetCurrentContext("somekey", "somenewvalue")

		assert.NotNil(t, err)
	})
}

func TestSet(t *testing.T) {
	t.Run("SetsContextValue", func(t *testing.T) {
		defer Reset()
		fs := afero.NewMemMapFs()
		SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(defaultConfig), 0644)

		Init("/tmp/testconfig.yml")

		Set("testcontext1", "somekey", "somenewvalue")
		value := GetString("somekey")

		assert.Equal(t, "somenewvalue", value)
	})

	t.Run("SetsDefaultValue", func(t *testing.T) {
		defer Reset()
		fs := afero.NewMemMapFs()
		SetFs(fs)

		afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(defaultConfig), 0644)

		Init("/tmp/testconfig.yml")

		Set("", "somekey", "somenewvalue")
		value := GetString("somekey")

		Set("", "current_context", "")
		defaultValue := GetString("somekey")

		assert.Equal(t, "somevalue", value)
		assert.Equal(t, "somenewvalue", defaultValue)
	})
}

func TestSwitchCurrentContext(t *testing.T) {
	defer Reset()
	fs := afero.NewMemMapFs()
	SetFs(fs)

	afero.WriteFile(fs, "/tmp/testconfig.yml", []byte(defaultConfig), 0644)

	Init("/tmp/testconfig.yml")

	SwitchCurrentContext("testcontext2")
	currentContext := GetCurrentContextName()

	assert.Equal(t, "testcontext2", currentContext)
}
