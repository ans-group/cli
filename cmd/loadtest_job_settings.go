package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestJobSettingsRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "sub-commands relating to job settings",
	}

	// Child commands
	cmd.AddCommand(loadtestJobSettingsShowCmd())

	return cmd
}

func loadtestJobSettingsShowCmd() *cobra.Command {
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
			return loadtestJobSettingsShow(getClient().LTaaSService(), cmd, args)
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

	return outputLoadTestJobSettings(allSettings)
}
