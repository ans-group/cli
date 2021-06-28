package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func loadbalancerClusterACLTemplateRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acltemplate",
		Short: "sub-commands relating to ACL templates",
	}

	// Child commands
	cmd.AddCommand(loadbalancerClusterACLTemplateShowCmd(f))

	return cmd
}

func loadbalancerClusterACLTemplateShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <cluster: id>",
		Short:   "Shows ACL templates",
		Long:    "This command shows ACL templates",
		Example: "ukfast loadbalancer cluster acltemplate show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing cluster")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerClusterACLTemplateShow),
	}
}

func loadbalancerClusterACLTemplateShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	clusterID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid cluster ID")
	}

	aclTemplates, err := service.GetClusterACLTemplates(clusterID)
	if err != nil {
		return fmt.Errorf("Error retrieving ACL templates: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLTemplatesProvider([]loadbalancer.ACLTemplates{aclTemplates}))
}
