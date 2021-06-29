package loadbalancer

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/loadbalancer"
)

func loadbalancerBindRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bind",
		Short: "sub-commands relating to binds",
	}

	// Child commands
	cmd.AddCommand(loadbalancerBindListCmd(f))

	return cmd
}

func loadbalancerBindListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists binds",
		Long:    "This command lists binds",
		Example: "ukfast loadbalancer bind list",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerBindList),
	}
}

func loadbalancerBindList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	binds, err := service.GetBinds(params)
	if err != nil {
		return fmt.Errorf("Error retrieving binds: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerBindsProvider(binds))
}
