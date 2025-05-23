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

func loadbalancerListenerACLRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "acl",
		Short: "sub-commands relating to ACLs",
	}

	// Child commands
	cmd.AddCommand(loadbalancerListenerACLListCmd(f))

	return cmd
}

func loadbalancerListenerACLListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <listener: id>",
		Short:   "Lists ACLs",
		Long:    "This command lists ACLs",
		Example: "ans loadbalancer listener acl list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerACLList),
	}
}

func loadbalancerListenerACLList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	acls, err := service.GetListenerACLs(listenerID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving ACLs: %s", err)
	}

	return output.CommandOutput(cmd, ACLCollection(acls))
}
