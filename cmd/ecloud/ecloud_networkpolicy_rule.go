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

func ecloudNetworkPolicyNetworkRuleRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "sub-commands relating to network policy network rules",
	}

	// Child commands
	cmd.AddCommand(ecloudNetworkPolicyNetworkRuleListCmd(f))

	return cmd
}

func ecloudNetworkPolicyNetworkRuleListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists network rules for network policy",
		Long:    "This command lists network rules for network policy",
		Example: "ans ecloud networkpolicy networkrule list np-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing network policy")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNetworkPolicyNetworkRuleList),
	}

	cmd.Flags().String("name", "", "Network rule name for filtering")

	return cmd
}

func ecloudNetworkPolicyNetworkRuleList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	rules, err := service.GetNetworkPolicyNetworkRules(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving network policy network rules: %s", err)
	}

	return output.CommandOutput(cmd, NetworkRuleCollection(rules))
}
