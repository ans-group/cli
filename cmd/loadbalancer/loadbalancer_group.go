package loadbalancer

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func loadbalancerGroupRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "group",
		Short: "sub-commands relating to groups",
	}

	// Child commands
	cmd.AddCommand(loadbalancerGroupListCmd(f))
	cmd.AddCommand(loadbalancerGroupShowCmd(f))

	return cmd
}

func loadbalancerGroupListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists groups",
		Long:    "This command lists groups",
		Example: "ukfast loadbalancer group list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerGroupList(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Group name for filtering")

	return cmd
}

func loadbalancerGroupList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

	groups, err := service.GetGroups(params)
	if err != nil {
		return fmt.Errorf("Error retrieving groups: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerGroupsProvider(groups))
}

func loadbalancerGroupShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <group: id>...",
		Short:   "Shows a group",
		Long:    "This command shows one or more groups",
		Example: "ukfast loadbalancer group show rtr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing group")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerGroupShow(c.LoadBalancerService(), cmd, args)
		},
	}
}

func loadbalancerGroupShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var groups []loadbalancer.Group
	for _, arg := range args {
		group, err := service.GetGroup(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving group [%s]: %s", arg, err)
			continue
		}

		groups = append(groups, group)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerGroupsProvider(groups))
}
