package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudFirewallRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firewall",
		Short: "sub-commands relating to firewalls",
	}

	// Child commands
	cmd.AddCommand(ecloudFirewallListCmd(f))
	cmd.AddCommand(ecloudFirewallShowCmd(f))

	return cmd
}

func ecloudFirewallListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists firewalls",
		Long:    "This command lists firewalls",
		Example: "ans ecloud firewall list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudFirewallList(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudFirewallList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	firewalls, err := service.GetFirewalls(params)
	if err != nil {
		return fmt.Errorf("error retrieving firewalls: %s", err)
	}

	return output.CommandOutput(cmd, FirewallCollection(firewalls))
}

func ecloudFirewallShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <firewall: id>...",
		Short:   "Shows a firewall",
		Long:    "This command shows one or more firewalls",
		Example: "ans ecloud vm firewall 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing firewall")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudFirewallShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudFirewallShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var firewalls []ecloud.Firewall
	for _, arg := range args {
		firewallID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid firewall ID [%s]", arg)
			continue
		}

		firewall, err := service.GetFirewall(firewallID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving firewall [%s]: %s", arg, err)
			continue
		}

		firewalls = append(firewalls, firewall)
	}

	return output.CommandOutput(cmd, FirewallCollection(firewalls))
}
