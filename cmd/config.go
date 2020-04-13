package cmd

import (
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func ConfigRootCmd(fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "sub-commands relating to CLI config",
	}

	// Child commands
	cmd.AddCommand(configSetCommand(fs))

	// Child root commands

	return cmd
}

func configSetCommand(fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "set",
		Short:   "Sets configuration properties",
		Long:    "This command sets configuration properties",
		Example: "ukfast config set --api_key \"secretkey\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			return configSet(fs, cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("api_key", "", "Specifies API key")
	cmd.Flags().Int("api_timeout_seconds", 0, "Specifies API timeout in seconds")
	cmd.Flags().String("api_uri", "", "Specifies API URI")
	cmd.Flags().Bool("api_insecure", false, "Specifies API TLS validation should be disabled")
	cmd.Flags().Bool("api_debug", false, "Specifies API debug logging should be enabled")
	cmd.Flags().Int("api_pagination_perpage", 0, "Specifies how many items should be retrieved per-page for paginated API requests")
	cmd.Flags().Int("command_wait_timeout_seconds", 0, "Specifies how long commands supporting 'wait' parameter should wait")
	cmd.Flags().Int("command_wait_sleep_seconds", 0, "Specifies how often commands supporting 'wait' parameter should poll")

	return cmd
}

func configSet(fs afero.Fs, cmd *cobra.Command, args []string) error {
	updated := false

	set := func(name string, value interface{}) {
		if cmd.Flags().Changed(name) {
			viper.Set(name, value)
			updated = true
		}
	}

	apiKey, _ := cmd.Flags().GetString("api_key")
	set("api_key", apiKey)
	apiTimeoutSeconds, _ := cmd.Flags().GetInt("api_timeout_seconds")
	set("api_timeout_seconds", apiTimeoutSeconds)
	apiURI, _ := cmd.Flags().GetString("api_uri")
	set("api_uri", apiURI)
	apiInsecure, _ := cmd.Flags().GetBool("api_insecure")
	set("api_insecure", apiInsecure)
	apiDebug, _ := cmd.Flags().GetBool("api_debug")
	set("api_debug", apiDebug)
	apiPaginationPerPage, _ := cmd.Flags().GetInt("api_pagination_perpage")
	set("api_pagination_perpage", apiPaginationPerPage)
	commandWaitTimeoutSeconds, _ := cmd.Flags().GetInt("command_wait_timeout_seconds")
	set("command_wait_timeout_seconds", commandWaitTimeoutSeconds)
	commandWaitSleepSeconds, _ := cmd.Flags().GetInt("command_wait_sleep_seconds")
	set("command_wait_sleep_seconds", commandWaitSleepSeconds)

	configFile := viper.GetViper().ConfigFileUsed()
	if len(configFile) < 1 {
		configFile = defaultConfigFile
	}

	if updated {
		viper.SetFs(fs)
		return viper.WriteConfigAs(configFile)
	}

	return nil
}
