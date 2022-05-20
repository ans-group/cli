package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudNICIPAddressRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ipaddress",
		Short: "sub-commands relating to IP addresses",
	}

	// Child commands
	cmd.AddCommand(ecloudNICIPAddressListCmd(f))
	cmd.AddCommand(ecloudNICIPAddressAssignCmd(f))
	cmd.AddCommand(ecloudNICIPAddressUnassignCmd(f))

	return cmd
}

func ecloudNICIPAddressListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists IP addresses for NIC",
		Long:    "This command lists IP addresses for NIC",
		Example: "ukfast ecloud nic ipaddress list ip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing NIC")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNICIPAddressList),
	}

	cmd.Flags().String("name", "", "IP address name for filtering")

	return cmd
}

func ecloudNICIPAddressList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	ips, err := service.GetNICIPAddresses(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving NIC IP addresses: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudIPAddressesProvider(ips))
}

func ecloudNICIPAddressAssignCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "assign",
		Short:   "Assigns an IP address to a NIC",
		Long:    "This command assigns an IP address to one or more NICs",
		Example: "ukfast ecloud nic ipaddress assign nic-abcdef12 --ip-address ip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing NIC")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNICIPAddressAssign),
	}

	cmd.Flags().String("ip-address", "", "IP address ID to assign")
	cmd.MarkFlagRequired("ip-address")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the IP address has been completely assigned")

	return cmd
}

func ecloudNICIPAddressAssign(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	assignRequest := ecloud.AssignIPAddressRequest{}
	assignRequest.IPAddressID, _ = cmd.Flags().GetString("ip-address")

	for _, arg := range args {
		taskID, err := service.AssignNICIPAddress(arg, assignRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error assigning IP address to NIC [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for NIC [%s]: %s", arg, err)
				continue
			}
		}
	}

	return nil
}

func ecloudNICIPAddressUnassignCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "unassign",
		Short:   "Unassigns an IP address from a NIC",
		Long:    "This command unassigns an IP address from one or more NICs",
		Example: "ukfast ecloud nic ipaddress unassign nic-abcdef12 --ip-address ip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing NIC")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudNICIPAddressUnassign),
	}

	cmd.Flags().String("ip-address", "", "IP address ID to unassign")
	cmd.MarkFlagRequired("ip-address")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the IP address has been completely unassigned")

	return cmd
}

func ecloudNICIPAddressUnassign(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	ipAddressID, _ := cmd.Flags().GetString("ip-address")

	for _, arg := range args {
		taskID, err := service.UnassignNICIPAddress(arg, ipAddressID)
		if err != nil {
			output.OutputWithErrorLevelf("Error unassigning IP address from NIC [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for NIC [%s]: %s", arg, err)
				continue
			}
		}
	}

	return nil
}
