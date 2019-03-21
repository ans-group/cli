package cmd

import (
	"errors"
	"strconv"

	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"

	"github.com/spf13/cobra"
)

func ecloudVirtualMachineDiskRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "disk",
		Short: "sub-commands relating to virtual machine disks",
	}

	// Child commands
	cmd.AddCommand(ecloudVirtualMachineDiskListCmd())
	cmd.AddCommand(ecloudVirtualMachineDiskUpdateCmd())

	return cmd
}

func ecloudVirtualMachineDiskListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <virtualmachine: id>",
		Short:   "lists virtual machine disks",
		Long:    "This command lists virtual machine disks",
		Example: "ukfast ecloud vm disk list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineDiskList(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Disk name for filtering")

	return cmd
}

func ecloudVirtualMachineDiskList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid virtual machine ID [%s]", args[0])
		return
	}

	vm, err := service.GetVirtualMachine(vmID)
	if err != nil {
		output.Fatalf("Error retrieving virtual machine [%s]: %s", args[0], err)
		return
	}

	outputECloudVirtualMachineDisks(vm.Disks)
}

func ecloudVirtualMachineDiskUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <virtualmachine: id> <disk: id>",
		Short:   "Updates a virtual machine disk",
		Long:    "This command updates a virtual machine disk",
		Example: "ukfast ecloud vm disk update 123 00000000-0000-0000-0000-000000000000 --capacity 25",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}
			if len(args) < 2 {
				return errors.New("Missing disk")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineDiskUpdate(getClient().ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().Int("capacity", 0, "Capacity of virtual machine disk in GB")

	return cmd
}

func ecloudVirtualMachineDiskUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {

	vmID, err := strconv.Atoi(args[0])
	if err != nil {
		output.Fatalf("Invalid virtual machine ID [%s]", args[0])
		return
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
		output.Fatalf("Error updating virtual machine [%d]: %s", vmID, err.Error())
		return
	}

	err = WaitForCommand(VirtualMachineStatusWaitFunc(service, vmID, ecloud.VirtualMachineStatusComplete))
	if err != nil {
		output.Fatalf("Error updating virtual machine [%d]: %s", vmID, err.Error())
		return
	}
}
