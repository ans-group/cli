package safedns

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func safednsSettingsRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "settings",
		Short: "sub-commands relating to settings",
	}

	// Child commands
	cmd.AddCommand(safednsSettingsShowCmd(f))

	return cmd
}

func safednsSettingsShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows settings for account",
		Long:    "This command shows SafeDNS settings for account",
		Example: "ukfast safedns settings show",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsSettingsShow(c.SafeDNSService(), cmd, args)
		},
	}
}

func safednsSettingsShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	settings, err := service.GetSettings()
	if err != nil {
		return fmt.Errorf("Error retrieving settings: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSSettingsProvider([]safedns.Settings{settings}))
}
