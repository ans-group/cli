package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
)

func loadbalancerDeploymentRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deployment",
		Short: "sub-commands relating to deployments",
	}

	// Child commands
	cmd.AddCommand(loadbalancerDeploymentListCmd(f))
	cmd.AddCommand(loadbalancerDeploymentShowCmd(f))

	return cmd
}

func loadbalancerDeploymentListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists deployments",
		Long:    "This command lists deployments",
		Example: "ans loadbalancer deployment list",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerDeploymentList),
	}
}

func loadbalancerDeploymentList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	deployments, err := service.GetDeployments(params)
	if err != nil {
		return fmt.Errorf("Error retrieving deployments: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerDeploymentsProvider(deployments))
}

func loadbalancerDeploymentShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <deployment: id>...",
		Short:   "Shows a deployment",
		Long:    "This command shows one or more deployments",
		Example: "ans loadbalancer deployment show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing deployment")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerDeploymentShow),
	}
}

func loadbalancerDeploymentShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var deployments []loadbalancer.Deployment
	for _, arg := range args {
		deploymentID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid deployment ID [%s]", arg)
			continue
		}

		deployment, err := service.GetDeployment(deploymentID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving deployment [%d]: %s", deploymentID, err)
			continue
		}

		deployments = append(deployments, deployment)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerDeploymentsProvider(deployments))
}
