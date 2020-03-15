package cmd

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudVirtualMachineRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "sub-commands relating to virtual machines",
	}

	// Child commands
	cmd.AddCommand(ecloudVirtualMachineListCmd())
	cmd.AddCommand(ecloudVirtualMachineShowCmd())
	cmd.AddCommand(ecloudVirtualMachineCreateCmd())
	cmd.AddCommand(ecloudVirtualMachineUpdateCmd())
	cmd.AddCommand(ecloudVirtualMachineStartCmd())
	cmd.AddCommand(ecloudVirtualMachineStopCmd())
	cmd.AddCommand(ecloudVirtualMachineRestartCmd())
	cmd.AddCommand(ecloudVirtualMachineDeleteCmd())

	// Child root commands
	cmd.AddCommand(ecloudVirtualMachineTagRootCmd())
	cmd.AddCommand(ecloudVirtualMachineDiskRootCmd())
	cmd.AddCommand(ecloudVirtualMachineTemplateRootCmd())

	return cmd
}

func ecloudVirtualMachineListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists virtual machines",
		Long:    "This command lists virtual machines",
		Example: "ukfast ecloud vm list",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineList(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "VM name for filtering")

	return cmd
}

func ecloudVirtualMachineList(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	if cmd.Flags().Changed("name") {
		filterName, _ := cmd.Flags().GetString("name")
		params.WithFilter(helper.GetFilteringInferOperator("name", filterName))
	}

	vms, err := service.GetVirtualMachines(params)
	if err != nil {
		output.Fatalf("Error retrieving virtual machines: %s", err)
		return
	}

	outputECloudVirtualMachines(vms)
}

func ecloudVirtualMachineShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <virtualmachine: id>...",
		Short:   "Shows a virtual machine",
		Long:    "This command shows one or more virtual machines",
		Example: "ukfast ecloud vm show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineShow(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	var vms []ecloud.VirtualMachine
	for _, arg := range args {
		vmID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid virtual machine ID [%s]", arg)
			continue
		}

		vm, err := service.GetVirtualMachine(vmID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving virtual machine [%s]: %s", arg, err)
			continue
		}

		vms = append(vms, vm)
	}

	outputECloudVirtualMachines(vms)
}

func ecloudVirtualMachineCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a virtual machine",
		Long:    "This command creates a virtual machine",
		Example: "ukfast ecloud vm create",
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineCreate(getClient().ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("environment", "", "Environment for virtual machine (Public, Hybrid, Private)")
	cmd.MarkFlagRequired("environment")
	cmd.Flags().String("template", "", "Template to use for virtual machine. Must be specified if appliance-id is omitted")
	cmd.Flags().String("appliance", "", "Appliance ID to use for virtual machine. Must be specified if template is omitted")
	cmd.Flags().Int("cpu", 0, "Amount of CPU cores for virtual machine")
	cmd.MarkFlagRequired("cpu")
	cmd.Flags().Int("ram", 0, "Amount of RAM (GB) for virtual machine")
	cmd.MarkFlagRequired("ram")
	cmd.Flags().Int("hdd", 0, "Primary disk size (GB) for virtual machine")
	cmd.MarkFlagRequired("hdd")
	cmd.Flags().String("name", "", "Name of virtual machine")
	cmd.Flags().String("computername", "", "Computer name of virtual machine")
	cmd.Flags().Bool("backup", false, "Specifies backup should be applied for virtual machine")
	cmd.Flags().Bool("support", false, "Specifies support should be applied for virtual machine")
	cmd.Flags().Bool("monitoring", false, "Specifies monitoring should be applied for virtual machine")
	cmd.Flags().IntSlice("monitoring-contact", []int{}, "Monitoring contact ID for monitoring, can be repeated")
	cmd.Flags().Int("solution", 0, "Solution ID for virtual machine")
	cmd.Flags().Int("datastore", 0, "Datastore ID for virtual machine")
	cmd.Flags().Int("site", 0, "Site ID for virtual machine")
	cmd.Flags().Int("network", 0, "Network ID for virtual machine")
	cmd.Flags().Bool("external-ip", false, "Specifies an external IP address should be created for virtual machine")
	cmd.Flags().StringSlice("tag", []string{}, "Tag for virtual machine, can be repeated, e.g. key=value")
	cmd.Flags().StringSlice("ssh-key", []string{}, "SSH public key for virtual machine, can be repeated")
	cmd.Flags().StringSlice("parameter", []string{}, "Parameters for virtual machine, can be repeated, e.g. key=value")
	cmd.Flags().Bool("encrypt", false, "Specifies that the virtual machine should be encrypted")
	cmd.Flags().String("role", "", "Specifies role that VM should be created with")
	cmd.Flags().String("bootstrap-script", "", "Specifies boot script that should be executed on first boot")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VM has been completely created before continuing on")

	return cmd
}

func ecloudVirtualMachineCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	createRequest := ecloud.CreateVirtualMachineRequest{}
	createRequest.Environment, _ = cmd.Flags().GetString("environment")
	createRequest.Template, _ = cmd.Flags().GetString("template")
	createRequest.ApplianceID, _ = cmd.Flags().GetString("appliance")
	createRequest.CPU, _ = cmd.Flags().GetInt("cpu")
	createRequest.RAM, _ = cmd.Flags().GetInt("ram")
	createRequest.HDD, _ = cmd.Flags().GetInt("hdd")
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.ComputerName, _ = cmd.Flags().GetString("computername")
	createRequest.Backup, _ = cmd.Flags().GetBool("backup")
	createRequest.Support, _ = cmd.Flags().GetBool("support")
	createRequest.Monitoring, _ = cmd.Flags().GetBool("monitoring")
	createRequest.MonitoringContacts, _ = cmd.Flags().GetIntSlice("monitoring-contact")
	createRequest.SolutionID, _ = cmd.Flags().GetInt("solution")
	createRequest.DatastoreID, _ = cmd.Flags().GetInt("datastore")
	createRequest.SiteID, _ = cmd.Flags().GetInt("site")
	createRequest.NetworkID, _ = cmd.Flags().GetInt("network")
	createRequest.ExternalIPRequired, _ = cmd.Flags().GetBool("external-ip")
	createRequest.Role, _ = cmd.Flags().GetString("role")
	createRequest.BootstrapScript, _ = cmd.Flags().GetString("bootstrap-script")

	if cmd.Flags().Changed("tag") {
		tagsFlag, _ := cmd.Flags().GetStringSlice("tag")
		tagsReq, err := GetCreateTagRequestFromStringArrayFlag(tagsFlag)
		if err != nil {
			output.Fatalf("Invalid tag data: %s", err)
			return
		}

		createRequest.Tags = tagsReq
	}

	if cmd.Flags().Changed("parameter") {
		parametersFlag, _ := cmd.Flags().GetStringSlice("parameter")
		parametersReq, err := GetCreateVirtualMachineRequestParameterFromStringArrayFlag(parametersFlag)
		if err != nil {
			output.Fatalf("Invalid parameter data: %s", err)
			return
		}

		createRequest.Parameters = parametersReq
	}

	if cmd.Flags().Changed("ssh-key") {
		sshKeys, _ := cmd.Flags().GetStringSlice("ssh-key")

		createRequest.SSHKeys = sshKeys
	}

	if cmd.Flags().Changed("encrypt") {
		encrypt, _ := cmd.Flags().GetBool("encrypt")
		createRequest.Encrypt = &encrypt
	}

	id, err := service.CreateVirtualMachine(createRequest)
	if err != nil {
		output.Fatalf("Error creating virtual machine: %s", err)
		return
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := WaitForCommand(VirtualMachineStatusWaitFunc(service, id, ecloud.VirtualMachineStatusComplete))
		if err != nil {
			output.Fatalf(err.Error())
			return
		}
	}

	vm, err := service.GetVirtualMachine(id)
	if err != nil {
		output.Fatalf("Error retrieving new virtual machine: %s", err)
		return
	}

	outputECloudVirtualMachines([]ecloud.VirtualMachine{vm})
}

func ecloudVirtualMachineUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <virtualmachine: id>...",
		Short:   "Updates a virtual machine",
		Long:    "This command updates one or more virtual machines",
		Example: "ukfast ecloud vm update 123 --name \"test vm 1\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineUpdate(getClient().ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().Int("cpu", 0, "Amount of CPU cores for virtual machine")
	cmd.Flags().Int("ram", 0, "Amount of RAM (GB) for virtual machine")
	cmd.Flags().String("name", "", "Name of virtual machine")

	return cmd
}

func ecloudVirtualMachineUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) {

	patchRequest := ecloud.PatchVirtualMachineRequest{}

	if cmd.Flags().Changed("cpu") {
		cpu, _ := cmd.Flags().GetInt("cpu")
		patchRequest.CPU = cpu
	}
	if cmd.Flags().Changed("ram") {
		ram, _ := cmd.Flags().GetInt("ram")
		patchRequest.RAM = ram
	}
	if cmd.Flags().Changed("name") {
		name, _ := cmd.Flags().GetString("name")
		patchRequest.Name = ptr.String(name)
	}

	var vms []ecloud.VirtualMachine

	for _, arg := range args {
		vmID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid virtual machine ID [%s]", arg)
			continue
		}

		err = service.PatchVirtualMachine(vmID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating virtual machine [%d]: %s", vmID, err.Error())
			continue
		}

		err = WaitForCommand(VirtualMachineStatusWaitFunc(service, vmID, ecloud.VirtualMachineStatusComplete))
		if err != nil {
			output.OutputWithErrorLevelf("Error updating virtual machine [%d]: %s", vmID, err.Error())
			continue
		}

		vm, err := service.GetVirtualMachine(vmID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated virtual machine [%d]: %s", vmID, err.Error())
			return
		}

		vms = append(vms, vm)
	}

	outputECloudVirtualMachines(vms)
}

func ecloudVirtualMachineStartCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "start <virtualmachine: id>...",
		Short:   "Starts a virtual machine",
		Long:    "This command starts one or more virtual machines",
		Example: "ukfast ecloud vm start 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineStart(getClient().ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineStart(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		vmID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid virtual machine ID [%s]", arg)
			continue
		}

		err = service.PowerOnVirtualMachine(vmID)
		if err != nil {
			output.OutputWithErrorLevelf("Error powering on virtual machine [%s]: %s", arg, err)
			continue
		}
	}
}

func ecloudVirtualMachineStopCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop <virtualmachine: id>...",
		Short:   "Stops a virtual machine",
		Long:    "This command stops one or more virtual machines",
		Example: "ukfast ecloud vm stop 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineStop(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().Bool("force", false, "Specifies that VM should be forcefully powered off")

	return cmd
}

func ecloudVirtualMachineStop(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	force, _ := cmd.Flags().GetBool("force")

	for _, arg := range args {
		vmID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid virtual machine ID [%s]", arg)
			continue
		}

		if force {
			err = service.PowerOffVirtualMachine(vmID)
			if err != nil {
				output.OutputWithErrorLevelf("Error powering off (forced) virtual machine [%s]: %s", arg, err)
				continue
			}
		} else {
			err = service.PowerShutdownVirtualMachine(vmID)
			if err != nil {
				output.OutputWithErrorLevelf("Error powering off virtual machine [%s]: %s", arg, err)
				continue
			}
		}
	}
}

func ecloudVirtualMachineRestartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restart <virtualmachine: id>...",
		Short:   "Restarts a virtual machine",
		Long:    "This command restarts one or more virtual machines",
		Example: "ukfast ecloud vm restart 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineRestart(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().Bool("force", false, "Specifies that VM should be forcefully powered off")

	return cmd
}

func ecloudVirtualMachineRestart(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	force, _ := cmd.Flags().GetBool("force")

	for _, arg := range args {
		vmID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid virtual machine ID [%s]", arg)
			continue
		}

		if force {
			err = service.PowerResetVirtualMachine(vmID)
			if err != nil {
				output.OutputWithErrorLevelf("Error restarting (forced) virtual machine [%s]: %s", arg, err)
				continue
			}
		} else {
			err = service.PowerRestartVirtualMachine(vmID)
			if err != nil {
				output.OutputWithErrorLevelf("Error restarting virtual machine [%s]: %s", arg, err)
				continue
			}
		}
	}
}

func ecloudVirtualMachineDeleteCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <virtualmachine: id>...",
		Short:   "Removes a virtual machine",
		Long:    "This command removes one or more virtual machines",
		Example: "ukfast ecloud vm delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing virtual machine")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ecloudVirtualMachineDelete(getClient().ECloudService(), cmd, args)
		},
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the VM has been completely deleted before continuing on")

	return cmd
}

func ecloudVirtualMachineDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		vmID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid virtual machine ID [%s]", arg)
			continue
		}

		err = service.DeleteVirtualMachine(vmID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing virtual machine [%d]: %s", vmID, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := WaitForCommand(VirtualMachineNotFoundWaitFunc(service, vmID))
			if err != nil {
				output.OutputWithErrorLevelf("Error removing virtual machine [%d]: %s", vmID, err)
				continue
			}
		}
	}
}

func VirtualMachineNotFoundWaitFunc(service ecloud.ECloudService, vmID int) WaitFunc {
	return func() (finished bool, err error) {
		vm, err := service.GetVirtualMachine(vmID)
		if err != nil {
			switch err.(type) {
			case *ecloud.VirtualMachineNotFoundError:
				return true, nil
			default:
				return false, fmt.Errorf("Failed to retrieve virtual machine [%d]: %s", vmID, err)
			}
		}

		if vm.Status == ecloud.VirtualMachineStatusFailed {
			return false, fmt.Errorf("Virtual machine [%d] in [%s] state", vmID, ecloud.VirtualMachineStatusFailed.String())
		}

		return false, nil
	}
}

func VirtualMachineStatusWaitFunc(service ecloud.ECloudService, vmID int, status ecloud.VirtualMachineStatus) WaitFunc {
	return func() (finished bool, err error) {
		vm, err := service.GetVirtualMachine(vmID)
		if err != nil {
			return false, fmt.Errorf("Failed to retrieve virtual machine [%d]: %s", vmID, err)
		}
		if vm.Status == ecloud.VirtualMachineStatusFailed {
			return false, fmt.Errorf("Virtual machine [%d] in [%s] state", vmID, ecloud.VirtualMachineStatusFailed.String())
		}
		if vm.Status == status {
			return true, nil
		}

		return false, nil
	}
}
