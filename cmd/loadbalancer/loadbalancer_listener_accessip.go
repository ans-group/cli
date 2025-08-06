package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/loadbalancer"
	"github.com/spf13/cobra"
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
		Example: "ans loadbalancer listener accessip list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing listener")
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
		return fmt.Errorf("invalid listener ID")
	}

	accessips, err := service.GetListenerAccessIPs(listenerID, params)
	if err != nil {
		return fmt.Errorf("error retrieving access IPs: %s", err)
	}

	return output.CommandOutput(cmd, AccessIPCollection(accessips))
}

func loadbalancerListenerAccessIPCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <listener: id>",
		Short:   "Creates an access IP",
		Long:    "This command creates an access IP",
		Example: "ans loadbalancer listener accessip create 123 --ip 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing listener")
			}

			return nil
		},
		RunE: loadbalancerCobraRunEFunc(f, loadbalancerListenerAccessIPCreate),
	}

	cmd.Flags().String("ip", "", "IP address for access IP")
	_ = cmd.MarkFlagRequired("ip")

	return cmd
}

func loadbalancerListenerAccessIPCreate(service loadbalancer.LoadBalancerService, cmd *cobra.Command, args []string) error {
	listenerID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("invalid listener ID")
	}

	ip, _ := cmd.Flags().GetString("ip")
	createRequest := loadbalancer.CreateAccessIPRequest{
		IP: connection.IPAddress(ip),
	}

	accessipID, err := service.CreateListenerAccessIP(listenerID, createRequest)
	if err != nil {
		return fmt.Errorf("error creating access IP: %s", err)
	}

	accessip, err := service.GetAccessIP(accessipID)
	if err != nil {
		return fmt.Errorf("error retrieving new access IP: %s", err)
	}

	return output.CommandOutput(cmd, AccessIPCollection([]loadbalancer.AccessIP{accessip}))
}
