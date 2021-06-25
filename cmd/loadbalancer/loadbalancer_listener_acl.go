package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
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
		Example: "ukfast loadbalancer listener acl list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			_, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("Invalid listener ID")
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

	params.WithFilter(connection.APIRequestFiltering{
		Property: "listener_id",
		Operator: connection.EQOperator,
		Value:    []string{args[0]},
	})

	acls, err := service.GetACLs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving ACLs: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerACLsProvider(acls))
}
