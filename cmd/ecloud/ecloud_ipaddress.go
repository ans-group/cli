package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudIPAddressRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ipaddress",
		Short: "sub-commands relating to IP addresses",
	}

	// Child commands
	cmd.AddCommand(ecloudIPAddressListCmd(f))
	cmd.AddCommand(ecloudIPAddressShowCmd(f))
	cmd.AddCommand(ecloudIPAddressCreateCmd(f))
	cmd.AddCommand(ecloudIPAddressUpdateCmd(f))
	cmd.AddCommand(ecloudIPAddressDeleteCmd(f))

	return cmd
}

func ecloudIPAddressListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists IP addresses",
		Long:    "This command lists IP addresses",
		Example: "ukfast ecloud ipaddress list",
		RunE:    ecloudCobraRunEFunc(f, ecloudIPAddressList),
	}

	cmd.Flags().String("name", "", "IP address name for filtering")

	return cmd
}

func ecloudIPAddressList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	ips, err := service.GetIPAddresses(params)
	if err != nil {
		return fmt.Errorf("Error retrieving IP addresses: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudIPAddressesProvider(ips))
}

func ecloudIPAddressShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <ip: id>...",
		Short:   "Shows a IP address",
		Long:    "This command shows one or more IP addresses",
		Example: "ukfast ecloud ipaddress show ip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing IP address")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudIPAddressShow),
	}
}

func ecloudIPAddressShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var ips []ecloud.IPAddress
	for _, arg := range args {
		ip, err := service.GetIPAddress(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving IP address [%s]: %s", arg, err)
			continue
		}

		ips = append(ips, ip)
	}

	return output.CommandOutput(cmd, OutputECloudIPAddressesProvider(ips))
}

func ecloudIPAddressCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a IP address",
		Long:    "This command creates a IP address",
		Example: "ukfast ecloud ipaddress create --network net-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudIPAddressCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of ip")
	cmd.Flags().String("ip-address", "", "IP address to allocate")
	cmd.Flags().String("network", "", "ID of network")
	cmd.MarkFlagRequired("network")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the IP address has been completely created")

	return cmd
}

func ecloudIPAddressCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateIPAddressRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	ipAddress, _ := cmd.Flags().GetString("ip-address")
	createRequest.IPAddress = connection.IPAddress(ipAddress)
	createRequest.NetworkID, _ = cmd.Flags().GetString("network")

	taskRef, err := service.CreateIPAddress(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating IP address: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for IP address task to complete: %s", err)
		}
	}

	ip, err := service.GetIPAddress(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new IP address: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudIPAddressesProvider([]ecloud.IPAddress{ip}))
}

func ecloudIPAddressUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <ip: id>...",
		Short:   "Updates a IP address",
		Long:    "This command updates one or more IP addresses",
		Example: "ukfast ecloud ipaddress update ip-abcdef12 --name \"my ip\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing IP address")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudIPAddressUpdate),
	}

	cmd.Flags().String("name", "", "Name of ip")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the IP address has been completely updated")

	return cmd
}

func ecloudIPAddressUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchIPAddressRequest{}
	patchRequest.Name, _ = cmd.Flags().GetString("name")

	var ips []ecloud.IPAddress
	for _, arg := range args {
		task, err := service.PatchIPAddress(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating IP address [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for IP address [%s]: %s", arg, err)
				continue
			}
		}

		ip, err := service.GetIPAddress(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated IP address [%s]: %s", arg, err)
			continue
		}

		ips = append(ips, ip)
	}

	return output.CommandOutput(cmd, OutputECloudIPAddressesProvider(ips))
}

func ecloudIPAddressDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <ip: id>...",
		Short:   "Removes a IP address",
		Long:    "This command removes one or more IP addresses",
		Example: "ukfast ecloud ipaddress delete ip-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing IP address")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudIPAddressDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the IP address has been completely removed")

	return cmd
}

func ecloudIPAddressDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteIPAddress(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing IP address [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for IP address [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
