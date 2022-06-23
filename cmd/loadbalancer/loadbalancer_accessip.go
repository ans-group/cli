package loadbalancer

import (
	"errors"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
)

func loadbalancerAccessIPRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accessip",
		Short: "sub-commands relating to access IPs",
	}

	// Child commands
	cmd.AddCommand(loadbalancerAccessIPShowCmd(f))
	cmd.AddCommand(loadbalancerAccessIPUpdateCmd(f))
	cmd.AddCommand(loadbalancerAccessIPDeleteCmd(f))

	return cmd
}

func loadbalancerAccessIPShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <accessip: id>...",
		Short:   "Shows an access IP",
		Long:    "This command shows one or more access IPs",
		Example: "ans loadbalancer accessip show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing access IP")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerAccessIPShow),
	}
}

func loadbalancerAccessIPShow(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	var accessips []loadbalancer.AccessIP
	for _, arg := range args {
		accessipID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid access IP ID [%s]", arg)
			continue
		}

		accessip, err := service.GetAccessIP(accessipID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving access IP [%d]: %s", accessipID, err)
			continue
		}

		accessips = append(accessips, accessip)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerAccessIPsProvider(accessips))
}

func loadbalancerAccessIPUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <accessip: id>...",
		Short:   "Updates an access IP",
		Long:    "This command updates one or more access IPs",
		Example: "ans loadbalancer accessip update 123 --ip 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing access IP")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerAccessIPUpdate),
	}

	cmd.Flags().String("ip", "", "IP address for access IP")

	return cmd
}

func loadbalancerAccessIPUpdate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	patchRequest := loadbalancer.PatchAccessIPRequest{}
	if cmd.Flags().Changed("ip") {
		ip, _ := cmd.Flags().GetString("ip")
		patchRequest.IP = connection.IPAddress(ip)
	}

	var accessips []loadbalancer.AccessIP
	for _, arg := range args {
		accessipID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid access IP ID [%s]", arg)
			continue
		}

		err = service.PatchAccessIP(accessipID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating access IP [%d]: %s", accessipID, err)
			continue
		}

		accessip, err := service.GetAccessIP(accessipID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated access IP [%d]: %s", accessipID, err)
			continue
		}

		accessips = append(accessips, accessip)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerAccessIPsProvider(accessips))
}

func loadbalancerAccessIPDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <accessip: id>...",
		Short:   "Removes an access IP",
		Long:    "This command removes one or more access IPs",
		Example: "ans loadbalancer accessip delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing access IP")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerAccessIPDelete),
	}
}

func loadbalancerAccessIPDelete(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		accessipID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid access IP ID [%s]", arg)
			continue
		}

		err = service.DeleteAccessIP(accessipID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing access IP [%d]: %s", accessipID, err)
			continue
		}
	}

	return nil
}
