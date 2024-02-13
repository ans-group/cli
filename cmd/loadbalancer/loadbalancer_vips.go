package loadbalancer

import (
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
)

func loadbalancerVipsCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vips",
		Short: "sub-commands relating to VIPs",
	}

	cmd.AddCommand(loadbalancerVipsListCmd(f))
	cmd.AddCommand(loadbalancerVipsShowCmd(f))

	return cmd
}

func loadbalancerVipsListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists VIPs assigned to your LB clusters",
		Example: "ans loadbalancer vips list",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerVipsList),
	}
}

func loadbalancerVipsList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	vips, err := service.GetVIPs(params)
	if err != nil {
		return fmt.Errorf("Error retrieving VIPs: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerVIPsProvider(vips))
}

func loadbalancerVipsShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Show details about a VIP",
		Example: "ans loadbalancer vips show 12345",
		RunE:    loadbalancerCobraRunEFunc(f, loadbalancerVipsShow),
	}
}

func loadbalancerVipsShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var vips []loadbalancer.VIP
	for _, arg := range args {
		vipID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid VIP ID [%s]", arg)
			continue
		}

		vip, err := service.GetVIP(vipID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving VIP [%d]: %s", vip, err)
			continue
		}

		vips = append(vips, vip)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerVIPsProvider(vips))
}
