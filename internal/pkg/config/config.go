package config

import (
	"errors"
	"fmt"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
)

var defaultConfigFile string
var configName = ".ukfast"
var initialised bool

// Init initialises the config package
func Init(configPath string) error {
	viper.SetEnvPrefix("ukf")
	viper.AutomaticEnv()

	if len(configPath) > 0 {
		viper.SetConfigFile(configPath)
	} else {
		// Find home directory
		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		// Search config in home directory with name ".ukfast" (without extension)
		viper.AddConfigPath(home)
		viper.SetConfigName(configName)
		defaultConfigFile = fmt.Sprintf("%s/%s.yml", home, configName)
	}

	// If a config file is found, read it in
	err := viper.ReadInConfig()
	if len(configPath) > 0 && err != nil {
		return fmt.Errorf("Failed to read config from file '%s': %s", configPath, err.Error())
	}

	initialised = true

	return nil
}

// Save saves the config to configured config file (or default)
func Save() error {
	if !initialised {
		return errors.New("Config not initialised")
	}

	configFile := viper.ConfigFileUsed()
	if len(configFile) < 1 {
		configFile = defaultConfigFile
	}

	return viper.WriteConfigAs(configFile)
}

// SetFs sets the filesystem instance to use
func SetFs(fs afero.Fs) {
	viper.SetFs(fs)
}

// GetCurrentContextName returns the name of the current context
func GetCurrentContextName() string {
	return viper.GetString("current_context")
}

func GetContextNames() []string {
	var contextNames []string
	contexts := viper.GetStringMap(getContextBaseKey())
	for contextName := range contexts {
		contextNames = append(contextNames, contextName)
	}

	return contextNames
}

func getCurrentContextKeyOrDefault(key string) string {
	return getContextKeyOrDefault(GetCurrentContextName(), key)
}

func getContextKeyOrDefault(contextName string, key string) string {
	if len(contextName) > 0 {
		return getContextSubKey(contextName, key)
	}

	return key
}

func getContextBaseKey() string {
	return "contexts"
}

func getContextKey(name string) string {
	return fmt.Sprintf("%s.%s", getContextBaseKey(), name)
}

func getContextSubKey(name string, key string) string {
	return fmt.Sprintf("%s.%s", getContextKey(name), key)
}

func SetCurrentContext(key string, value any) error {
	contextName := GetCurrentContextName()
	if len(contextName) < 1 {
		return errors.New("current context not set")
	}

	return Set(GetCurrentContextName(), key, value)
}

func Set(contextName string, key string, value any) error {
	if len(contextName) < 1 {
		return errors.New("context name not provided")
	}

	viper.Set(getContextSubKey(contextName, key), value)
	return nil
}

func SetDefault(contextName string, key string, value any) {
	viper.SetDefault(getContextSubKey(contextName, key), value)
}

func SwitchCurrentContext(contextName string) error {
	if !ContextExists(contextName) {
		return fmt.Errorf("context not defined with name '%s'", contextName)
	}

	viper.Set("current_context", contextName)
	return nil
}

func ContextExists(contextName string) bool {
	return viper.IsSet(getContextKey(contextName))
}

func Reset() {
	viper.Reset()
}

func GetString(key string) string {
	return viper.GetString(getCurrentContextKeyOrDefault(key))
}

func GetInt(key string) int {
	return viper.GetInt(getCurrentContextKeyOrDefault(key))
}

func GetBool(key string) bool {
	return viper.GetBool(getCurrentContextKeyOrDefault(key))
}

func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(getCurrentContextKeyOrDefault(key))
}
