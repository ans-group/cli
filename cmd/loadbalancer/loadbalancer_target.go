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

func loadbalancerTargetRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "target",
		Short: "sub-commands relating to targets",
	}

	// Child commands
	cmd.AddCommand(loadbalancerTargetListCmd(f))
	cmd.AddCommand(loadbalancerTargetShowCmd(f))

	return cmd
}

func loadbalancerTargetListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists targets",
		Long:    "This command lists targets",
		Example: "ukfast loadbalancer target list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerTargetList(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Target name for filtering")

	return cmd
}

func loadbalancerTargetList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("name", "name"))

	targets, err := service.GetTargets(params)
	if err != nil {
		return fmt.Errorf("Error retrieving targets: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider(targets))
}

func loadbalancerTargetShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <target: id>...",
		Short:   "Shows a target",
		Long:    "This command shows one or more targets",
		Example: "ukfast loadbalancer target show rtr-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing target")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerTargetShow(c.LoadBalancerService(), cmd, args)
		},
	}
}

func loadbalancerTargetShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var targets []loadbalancer.Target
	for _, arg := range args {
		target, err := service.GetTarget(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving target [%s]: %s", arg, err)
			continue
		}

		targets = append(targets, target)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerTargetsProvider(targets))
}
