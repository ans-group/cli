package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudFirewallRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firewall",
		Short: "sub-commands relating to firewalls",
	}

	// Child commands
	cmd.AddCommand(ecloudFirewallListCmd())
	cmd.AddCommand(ecloudFirewallShowCmd())

	return cmd
}

func ecloudFirewallListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists firewalls",
		Long:    "This command lists firewalls",
		Example: "ukfast ecloud firewall list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudFirewallList(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudFirewallList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	firewalls, err := service.GetFirewalls(params)
	if err != nil {
		output.Fatalf("Error retrieving firewalls: %s", err)
		return
	}

	outputECloudFirewalls(firewalls)
}

func ecloudFirewallShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <firewall: id>...",
		Short:   "Shows a firewall",
		Long:    "This command shows one or more firewalls",
		Example: "ukfast ecloud vm firewall 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing firewall")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudFirewallShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudFirewallShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
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

	outputECloudFirewalls(firewalls)
}
