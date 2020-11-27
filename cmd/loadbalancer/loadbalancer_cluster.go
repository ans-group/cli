package loadbalancer

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
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
	cmd.AddCommand(loadbalancerClusterDeleteCmd(f))

	return cmd
}

func loadbalancerClusterListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists clusters",
		Long:    "This command lists clusters",
		Example: "ukfast loadbalancer cluster list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerClusterList(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Cluster name for filtering")

	return cmd
}

func loadbalancerClusterList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	flaghelper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, flaghelper.NewStringFilterFlag("name", "name"))

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
		Example: "ukfast loadbalancer cluster show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing cluster")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerClusterShow(c.LoadBalancerService(), cmd, args)
		},
	}
}

func loadbalancerClusterShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var clusters []loadbalancer.Cluster
	for _, arg := range args {
		cluster, err := service.GetCluster(arg)
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
		Example: "ukfast loadbalancer cluster update 00000000-0000-0000-0000-000000000000 --name \"my cluster\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing cluster")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return loadbalancerClusterUpdate(c.LoadBalancerService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of cluster")

	return cmd
}

func loadbalancerClusterUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchClusterRequest{}

	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest.Name = ptr.String(name)
	}

	var clusters []loadbalancer.Cluster
	for _, arg := range args {
		err := service.PatchCluster(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating cluster [%s]: %s", arg, err)
			continue
		}

		cluster, err := service.GetCluster(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated cluster [%s]: %s", arg, err)
			continue
		}

		clusters = append(clusters, cluster)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerClustersProvider(clusters))
}

func loadbalancerClusterDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <cluster: id...>",
		Short:   "Removes a cluster",
		Long:    "This command removes one or more clusters",
		Example: "ukfast loadbalancer cluster delete 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing cluster")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			loadbalancerClusterDelete(c.LoadBalancerService(), cmd, args)
			return nil
		},
	}
}

func loadbalancerClusterDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteCluster(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing cluster [%s]: %s", arg, err)
		}
	}
}
