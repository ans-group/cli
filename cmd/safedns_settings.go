package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func safednsSettingsRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "sub-commands relating to settings",
	}

	// Child commands
	cmd.AddCommand(safednsSettingsShowCmd())

	return cmd
}

func safednsSettingsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows settings for account",
		Long:    "This command shows SafeDNS settings for account",
		Example: "ukfast safedns settings show",
		RunE: func(cmd *cobra.Command, args []string) error {
			return safednsSettingsShow(getClient().SafeDNSService(), cmd, args)
		},
	}
}

func safednsSettingsShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	settings, err := service.GetSettings()
	if err != nil {
		return fmt.Errorf("Error retrieving settings: %s", err)
	}

	outputSafeDNSSettings([]safedns.Settings{settings})
	return nil
}
