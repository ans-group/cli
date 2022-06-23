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

func ecloudLoadBalancerSpecRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "loadbalancerspec",
		Short: "sub-commands relating to load balancer specs",
	}

	// Child commands
	cmd.AddCommand(ecloudLoadBalancerSpecListCmd(f))
	cmd.AddCommand(ecloudLoadBalancerSpecShowCmd(f))

	return cmd
}

func ecloudLoadBalancerSpecListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists load balancer specs",
		Long:    "This command lists load balancer specs",
		Example: "ukfast ecloud loadbalancer list",
		RunE:    ecloudCobraRunEFunc(f, ecloudLoadBalancerSpecList),
	}

	cmd.Flags().String("name", "", "Name for filtering")

	return cmd
}

func ecloudLoadBalancerSpecList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	lbs, err := service.GetLoadBalancerSpecs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving load balancer specs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerSpecsProvider(lbs))
}

func ecloudLoadBalancerSpecShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <spec: id>...",
		Short:   "Shows an load balancer spec",
		Long:    "This command shows one or more load balancer specs",
		Example: "ukfast ecloud loadbalancerspec show lbn-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing load balancer spec")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudLoadBalancerSpecShow),
	}
}

func ecloudLoadBalancerSpecShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var lbs []ecloud.LoadBalancerSpec
	for _, arg := range args {
		lb, err := service.GetLoadBalancerSpec(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving load balancer spec [%s]: %s", arg, err)
			continue
		}

		lbs = append(lbs, lb)
	}

	return output.CommandOutput(cmd, OutputECloudLoadBalancerSpecsProvider(lbs))
}
