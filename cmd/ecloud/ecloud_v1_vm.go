package ecloud

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudVirtualMachineRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "vm",
		Short: "sub-commands relating to virtual machines",
	}

	// Child commands
	cmd.AddCommand(ecloudVirtualMachineListCmd(f))
	cmd.AddCommand(ecloudVirtualMachineShowCmd(f))
	cmd.AddCommand(ecloudVirtualMachineCreateCmd(f))
	cmd.AddCommand(ecloudVirtualMachineUpdateCmd(f))
	cmd.AddCommand(ecloudVirtualMachineStartCmd(f))
	cmd.AddCommand(ecloudVirtualMachineStopCmd(f))
	cmd.AddCommand(ecloudVirtualMachineRestartCmd(f))
	cmd.AddCommand(ecloudVirtualMachineDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudVirtualMachineTagRootCmd(f))
	cmd.AddCommand(ecloudVirtualMachineDiskRootCmd(f))
	cmd.AddCommand(ecloudVirtualMachineTemplateRootCmd(f))
	cmd.AddCommand(ecloudVirtualMachineConsoleRootCmd(f))

	return cmd
}

func ecloudVirtualMachineListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists virtual machines",
		Long:    "This command lists virtual machines",
		Example: "ukfast ecloud vm list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVirtualMachineList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "VM name for filtering")

	return cmd
}

func ecloudVirtualMachineList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	vms, err := service.GetVirtualMachines(params)
	if err != nil {
		return fmt.Errorf("Error retrieving virtual machines: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVirtualMachinesProvider(vms))
}

func ecloudVirtualMachineShowCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVirtualMachineShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudVirtualMachineShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
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

	return output.CommandOutput(cmd, OutputECloudVirtualMachinesProvider(vms))
}

func ecloudVirtualMachineCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a virtual machine",
		Long:    "This command creates a virtual machine",
		Example: "ukfast ecloud vm create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVirtualMachineCreate(c.ECloudService(), cmd, args)
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
	cmd.Flags().Int("pod", 0, "Pod ID for virtual machine")

	return cmd
}

func ecloudVirtualMachineCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
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
	createRequest.PodID, _ = cmd.Flags().GetInt("pod")

	if cmd.Flags().Changed("tag") {
		tagsFlag, _ := cmd.Flags().GetStringSlice("tag")
		tagsReq, err := GetCreateTagRequestFromStringArrayFlag(tagsFlag)
		if err != nil {
			return fmt.Errorf("Invalid tag data: %s", err)
		}

		createRequest.Tags = tagsReq
	}

	if cmd.Flags().Changed("parameter") {
		parametersFlag, _ := cmd.Flags().GetStringSlice("parameter")
		parametersReq, err := GetCreateVirtualMachineRequestParameterFromStringArrayFlag(parametersFlag)
		if err != nil {
			return fmt.Errorf("Invalid parameter data: %s", err)
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
		return fmt.Errorf("Error creating virtual machine: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(VirtualMachineStatusWaitFunc(service, id, ecloud.VirtualMachineStatusComplete))
		if err != nil {
			return err
		}
	}

	vm, err := service.GetVirtualMachine(id)
	if err != nil {
		return fmt.Errorf("Error retrieving new virtual machine: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVirtualMachinesProvider([]ecloud.VirtualMachine{vm}))
}

func ecloudVirtualMachineUpdateCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudVirtualMachineUpdate(c.ECloudService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().Int("cpu", 0, "Amount of CPU cores for virtual machine")
	cmd.Flags().Int("ram", 0, "Amount of RAM (GB) for virtual machine")
	cmd.Flags().String("name", "", "Name of virtual machine")

	return cmd
}

func ecloudVirtualMachineUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {

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

		err = helper.WaitForCommand(VirtualMachineStatusWaitFunc(service, vmID, ecloud.VirtualMachineStatusComplete))
		if err != nil {
			output.OutputWithErrorLevelf("Error updating virtual machine [%d]: %s", vmID, err.Error())
			continue
		}

		vm, err := service.GetVirtualMachine(vmID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated virtual machine [%d]: %s", vmID, err.Error())
			continue
		}

		vms = append(vms, vm)
	}

	return output.CommandOutput(cmd, OutputECloudVirtualMachinesProvider(vms))
}

func ecloudVirtualMachineStartCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ecloudVirtualMachineStart(c.ECloudService(), cmd, args)
			return nil
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

func ecloudVirtualMachineStopCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ecloudVirtualMachineStop(c.ECloudService(), cmd, args)
			return nil
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

func ecloudVirtualMachineRestartCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ecloudVirtualMachineRestart(c.ECloudService(), cmd, args)
			return nil
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

func ecloudVirtualMachineDeleteCmd(f factory.ClientFactory) *cobra.Command {
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
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			ecloudVirtualMachineDelete(c.ECloudService(), cmd, args)
			return nil
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
			err := helper.WaitForCommand(VirtualMachineNotFoundWaitFunc(service, vmID))
			if err != nil {
				output.OutputWithErrorLevelf("Error removing virtual machine [%d]: %s", vmID, err)
				continue
			}
		}
	}
}

func VirtualMachineNotFoundWaitFunc(service ecloud.ECloudService, vmID int) helper.WaitFunc {
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

func VirtualMachineStatusWaitFunc(service ecloud.ECloudService, vmID int, status ecloud.VirtualMachineStatus) helper.WaitFunc {
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
