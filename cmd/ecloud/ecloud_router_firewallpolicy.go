package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudRouterFirewallPolicyRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "firewallpolicy",
		Short: "sub-commands relating to router firewall policies",
	}

	// Child commands
	cmd.AddCommand(ecloudRouterFirewallPolicyListCmd(f))

	return cmd
}

func ecloudRouterFirewallPolicyListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists router firewall policies",
		Long:    "This command lists router firewall policies",
		Example: "ukfast ecloud router firewallpolicy list rtr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing router")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudRouterFirewallPolicyList),
	}

	cmd.Flags().String("name", "", "Firewall policy name for filtering")

	return cmd
}

func ecloudRouterFirewallPolicyList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	policies, err := service.GetRouterFirewallPolicies(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving router firewall policies: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudFirewallPoliciesProvider(policies))
}
