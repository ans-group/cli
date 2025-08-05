package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
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
		Example: "ans loadbalancer cluster acltemplate show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing cluster")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerClusterACLTemplateShow),
	}
}

func loadbalancerClusterACLTemplateShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	clusterID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid cluster ID")
	}

	aclTemplates, err := service.GetClusterACLTemplates(clusterID)
	if err != nil {
		return fmt.Errorf("error retrieving ACL templates: %s", err)
	}

	return output.CommandOutput(cmd, ACLTemplatesCollection([]loadbalancer.ACLTemplates{aclTemplates}))
}
