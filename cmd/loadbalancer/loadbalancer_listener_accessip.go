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

func loadbalancerListenerAccessIPRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accessip",
		Short: "sub-commands relating to accessips",
	}

	// Child commands
	cmd.AddCommand(loadbalancerListenerAccessIPListCmd(f))
	cmd.AddCommand(loadbalancerListenerAccessIPCreateCmd(f))

	return cmd
}

func loadbalancerListenerAccessIPListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <listener: id>",
		Short:   "Lists access IPs",
		Long:    "This command lists access IPs",
		Example: "ukfast loadbalancer listener accessip list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerAccessIPList),
	}
}

func loadbalancerListenerAccessIPList(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	accessips, err := service.GetListenerAccessIPs(listenerID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving access IPs: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerAccessIPsProvider(accessips))
}

func loadbalancerListenerAccessIPCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <listener: id>",
		Short:   "Creates an access IP",
		Long:    "This command creates an access IP",
		Example: "ukfast loadbalancer listener accessip create 123 --ip 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerAccessIPCreate),
	}

	cmd.Flags().String("ip", "", "IP address for access IP")
	cmd.MarkFlagRequired("ip")

	return cmd
}

func loadbalancerListenerAccessIPCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid listener ID")
	}

	ip, _ := cmd.Flags().GetString("ip")
	createRequest := loadbalancer.CreateAccessIPRequest{
		IP: connection.IPAddress(ip),
	}

	accessipID, err := service.CreateListenerAccessIP(listenerID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating access IP: %s", err)
	}

	accessip, err := service.GetAccessIP(accessipID)
	if err != nil {
		return fmt.Errorf("Error retrieving new access IP: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadBalancerAccessIPsProvider([]loadbalancer.AccessIP{accessip}))
}
