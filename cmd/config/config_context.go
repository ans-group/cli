package config

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/config"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func configContextRootCmd(fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "sub-commands relating to CLI config",
	}

	// Child commands
	cmd.AddCommand(configContextUpdateCmd(fs))
	cmd.AddCommand(configContextListCmd())
	cmd.AddCommand(configContextShowCmd())
	cmd.AddCommand(configContextSwitchCmd(fs))

	return cmd
}

func configContextUpdateCmd(fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update",
		Short:   "Updates context configuration",
		Long:    "This command updates context configuration",
		Example: "ans config context update mycontext --api-key \"secretkey\"",
		RunE: func(cmd *cobra.Command, args []string) error {
			return configContextUpdate(fs, cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().Bool("current", false, "Specifies that current context should be updated")
	cmd.Flags().String("api-key", "", "Specifies API key")
	cmd.Flags().Int("api-timeout-seconds", 0, "Specifies API timeout in seconds")
	cmd.Flags().String("api-uri", "", "Specifies API URI")
	cmd.Flags().Bool("api-insecure", false, "Specifies API TLS validation should be disabled")
	cmd.Flags().Bool("api-debug", false, "Specifies API debug logging should be enabled")
	cmd.Flags().Int("api-pagination-perpage", 0, "Specifies how many items should be retrieved per-page for paginated API requests")
	cmd.Flags().Int("command-wait-timeout-seconds", 0, "Specifies how long commands supporting 'wait' parameter should wait")
	cmd.Flags().Int("command-wait-sleep-seconds", 0, "Specifies how often commands supporting 'wait' parameter should poll")

	return cmd
}

func configContextUpdate(fs afero.Fs, cmd *cobra.Command, args []string) error {
	updated := false

	updateCurrentContext, _ := cmd.Flags().GetBool("current")

	set := func(name string, flagName string, value interface{}) {
		if cmd.Flags().Changed(flagName) {
			if updateCurrentContext {
				err := config.SetCurrentContext(name, value)
				if err != nil {
					output.Fatalf("Failed to update current context: %s", err)
				}
			} else {
				for _, context := range args {
					config.Set(context, name, value)
				}
			}
			updated = true
		}
	}

	apiKey, _ := cmd.Flags().GetString("api-key")
	set("api_key", "api-key", apiKey)
	apiTimeoutSeconds, _ := cmd.Flags().GetInt("api-timeout-seconds")
	set("api_timeout_seconds", "api-timeout-seconds", apiTimeoutSeconds)
	apiURI, _ := cmd.Flags().GetString("api-uri")
	set("api_uri", "api-uri", apiURI)
	apiInsecure, _ := cmd.Flags().GetBool("api-insecure")
	set("api_insecure", "api-insecure", apiInsecure)
	apiDebug, _ := cmd.Flags().GetBool("api-debug")
	set("api_debug", "api-debug", apiDebug)
	apiPaginationPerPage, _ := cmd.Flags().GetInt("api-pagination-perpage")
	set("api_pagination_perpage", "api-pagination-perpage", apiPaginationPerPage)
	commandWaitTimeoutSeconds, _ := cmd.Flags().GetInt("command-wait-timeout-seconds")
	set("command_wait_timeout_seconds", "command-wait-timeout-seconds", commandWaitTimeoutSeconds)
	commandWaitSleepSeconds, _ := cmd.Flags().GetInt("command-wait-sleep-seconds")
	set("command_wait_sleep_seconds", "command-wait-sleep-seconds", commandWaitSleepSeconds)

	if updated {
		config.SetFs(fs)
		return config.Save()
	}

	return nil
}

func configContextListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists contexts",
		Long:    "This command lists contexts",
		Example: "ans config context list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return configContextList(cmd)
		},
	}
}

func configContextList(cmd *cobra.Command) error {
	contextNames := config.GetContextNames()
	currentContextName := config.GetCurrentContextName()

	var contexts []Context

	for _, contextName := range contextNames {
		contexts = append(contexts, Context{
			Name:   contextName,
			Active: contextName == currentContextName,
		})
	}

	return output.CommandOutput(cmd, ContextCollection(contexts))
}

func configContextShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows current context",
		Long:    "This command shows the current context",
		Example: "ans config context show",
		RunE: func(cmd *cobra.Command, args []string) error {
			return configContextShow(cmd)
		},
	}
}

func configContextShow(cmd *cobra.Command) error {
	currentContextName := config.GetCurrentContextName()

	if len(currentContextName) < 1 {
		return errors.New("no context set")
	}

	context := Context{
		Name:   currentContextName,
		Active: true,
	}

	return output.CommandOutput(cmd, ContextCollection([]Context{context}))
}

func configContextSwitchCmd(fs afero.Fs) *cobra.Command {
	return &cobra.Command{
		Use:     "switch",
		Short:   "Switches current context",
		Long:    "This command switches the current context",
		Example: "ans config context switch mycontext",
		Aliases: []string{"use", "select"},
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing context")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return configContextSwitch(fs, cmd, args)
		},
	}
}

func configContextSwitch(fs afero.Fs, cmd *cobra.Command, args []string) error {
	config.SetFs(fs)

	err := config.SwitchCurrentContext(args[0])
	if err != nil {
		return fmt.Errorf("failed to switch context: %s", err)
	}

	err = config.Save()
	if err != nil {
		return fmt.Errorf("failed to write new context: %s", err)
	}

	fmt.Printf("Switched to context \"%s\"\n", args[0])
	return nil
}
