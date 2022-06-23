package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"

	"github.com/spf13/cobra"
)

func ecloudVirtualMachineDiskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk",
		Short: "sub-commands relating to virtual machine disks",
	}

	// Child commands
	cmd.AddCommand(ecloudVirtualMachineDiskListCmd(f))
	cmd.AddCommand(ecloudVirtualMachineDiskUpdateCmd(f))

	return cmd
}

func ecloudVirtualMachineDiskListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <virtualmachine: id>",
		Short:   "lists virtual machine disks",
		Long:    "This command lists virtual machine disks",
		Example: "ans ecloud vm disk list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVirtualMachineDiskList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Disk name for filtering")

	return cmd
}

func ecloudVirtualMachineDiskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	vm, err := service.GetVirtualMachine(vmID)
	if err != nil {
		return fmt.Errorf("Error retrieving virtual machine [%s]: %s", args[0], err)
	}

	return output.CommandOutput(cmd, OutputECloudVirtualMachineDisksProvider(vm.Disks))
}

func ecloudVirtualMachineDiskUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <virtualmachine: id> <disk: id>",
		Short:   "Updates a virtual machine disk",
		Long:    "This command updates a virtual machine disk",
		Example: "ans ecloud vm disk update 123 00000000-0000-0000-0000-000000000000 --capacity 25",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}
			if len(args) < 2 {
				return errors.New("Missing disk")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVirtualMachineDiskUpdate(c.ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().Int("capacity", 0, "Capacity of virtual machine disk in GB")

	return cmd
}

func ecloudVirtualMachineDiskUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {

	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		return fmt.Errorf("Invalid virtual machine ID [%s]", args[0])
	}

	diskPatch := ecloud.PatchVirtualMachineRequestDisk{
		UUID: args[1],
	}

	if cmd.Flags().Changed("capacity") {
		capacity, _ := cmd.Flags().GetInt("capacity")
		diskPatch.Capacity = capacity
	}

	patchRequest := ecloud.PatchVirtualMachineRequest{
		Disks: []ecloud.PatchVirtualMachineRequestDisk{
			diskPatch,
		},
	}

	err = service.PatchVirtualMachine(vmID, patchRequest)
	if err != nil {
		return fmt.Errorf("Error updating virtual machine [%d]: %s", vmID, err.Error())
	}

	err = helper.WaitForCommand(VirtualMachineStatusWaitFunc(service, vmID, ecloud.VirtualMachineStatusComplete))
	if err != nil {
		return fmt.Errorf("Error updating virtual machine [%d]: %s", vmID, err.Error())
	}

	return nil
}
