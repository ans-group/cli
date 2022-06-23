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

func loadbalancerClusterRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cluster",
		Short: "sub-commands relating to clusters",
	}

	// Child commands
	cmd.AddCommand(loadbalancerClusterListCmd(f))
	cmd.AddCommand(loadbalancerClusterShowCmd(f))
	cmd.AddCommand(loadbalancerClusterUpdateCmd(f))
	cmd.AddCommand(loadbalancerClusterDeployCmd(f))
	cmd.AddCommand(loadbalancerClusterValidateCmd(f))

	// Child root commands
	cmd.AddCommand(loadbalancerClusterACLTemplateRootCmd(f))

	return cmd
}

func loadbalancerClusterListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists clusters",
		Long:    "This command lists clusters",
		Example: "ans loadbalancer cluster list",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerClusterList),
	}
}

func loadbalancerClusterList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	clusters, err := service.GetClusters(params)
	if err != nil {
		return fmt.Errorf("Error retrieving clusters: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerClustersProvider(clusters))
}

func loadbalancerClusterShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <cluster: id>...",
		Short:   "Shows a cluster",
		Long:    "This command shows one or more clusters",
		Example: "ans loadbalancer cluster show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing cluster")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerClusterShow),
	}
}

func loadbalancerClusterShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var clusters []loadbalancer.Cluster
	for _, arg := range args {
		clusterID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid cluster ID [%s]", arg)
			continue
		}

		cluster, err := service.GetCluster(clusterID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving cluster [%s]: %s", arg, err)
			continue
		}

		clusters = append(clusters, cluster)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerClustersProvider(clusters))
}

func loadbalancerClusterUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <cluster: id>...",
		Short:   "Updates a cluster",
		Long:    "This command updates one or more clusters",
		Example: "ans loadbalancer cluster update 123 --name mycluster",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing cluster")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerClusterUpdate),
	}

	cmd.Flags().String("name", "", "Name of cluster")

	return cmd
}

func loadbalancerClusterUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchClusterRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var clusters []loadbalancer.Cluster
	for _, arg := range args {
		clusterID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid cluster ID [%s]", arg)
			continue
		}

		err = service.PatchCluster(clusterID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating cluster [%s]: %s", arg, err)
			continue
		}

		cluster, err := service.GetCluster(clusterID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated cluster [%s]: %s", arg, err)
			continue
		}

		clusters = append(clusters, cluster)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerClustersProvider(clusters))
}

func loadbalancerClusterDeployCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "deploy <cluster: id>...",
		Short:   "Deploys a cluster",
		Long:    "This command deploys one or more clusters",
		Example: "ans loadbalancer cluster deploy 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing cluster")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerClusterDeploy),
	}
}

func loadbalancerClusterDeploy(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		clusterID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid cluster ID [%s]", arg)
			continue
		}

		err = service.DeployCluster(clusterID)
		if err != nil {
			output.OutputWithErrorLevelf("Error deploying cluster [%s]: %s", arg, err)
			continue
		}
	}

	return nil
}

func loadbalancerClusterValidateCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "validate <cluster: id>...",
		Short:   "Validates a cluster",
		Long:    "This command validates one or more clusters",
		Example: "ans loadbalancer cluster validate 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing cluster")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerClusterValidate),
	}
}

func loadbalancerClusterValidate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		clusterID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid cluster ID [%s]", arg)
			continue
		}

		err = service.ValidateCluster(clusterID)
		if err != nil {
			output.OutputWithErrorLevelf("Error validating cluster [%s]: %s", arg, err)
			continue
		}
	}

	return nil
}
