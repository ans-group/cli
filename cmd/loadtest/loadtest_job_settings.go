package loadtest

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestJobSettingsRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "sub-commands relating to job settings",
	}

	// Child commands
	cmd.AddCommand(loadtestJobSettingsShowCmd(f))

	return cmd
}

func loadtestJobSettingsShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows job settings",
		Long:    "This command shows job settings",
		Example: "ukfast loadtest job settings show",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing job")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestJobSettingsShow(f.NewClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestJobSettingsShow(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	var allSettings []ltaas.JobSettings
	for _, arg := range args {
		settings, err := service.GetJobSettings(arg)
		if err != nil {
			return fmt.Errorf("Error retrieving job settings: %s", err)
		}

		allSettings = append(allSettings, settings)
	}

	return output.CommandOutput(cmd, OutputLoadTestJobSettingsProvider(allSettings))
}
