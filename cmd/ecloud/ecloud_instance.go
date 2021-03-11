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

func ecloudInstanceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "sub-commands relating to instances",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceListCmd(f))
	cmd.AddCommand(ecloudInstanceShowCmd(f))
	cmd.AddCommand(ecloudInstanceCreateCmd(f))
	cmd.AddCommand(ecloudInstanceUpdateCmd(f))
	cmd.AddCommand(ecloudInstanceDeleteCmd(f))
	cmd.AddCommand(ecloudInstanceLockCmd(f))
	cmd.AddCommand(ecloudInstanceUnlockCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudInstanceVolumeRootCmd(f))
	cmd.AddCommand(ecloudInstanceCredentialRootCmd(f))
	cmd.AddCommand(ecloudInstanceNICRootCmd(f))

	return cmd
}

func ecloudInstanceListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists instances",
		Long:    "This command lists instances",
		Example: "ukfast ecloud instance list",
		RunE:    ecloudCobraRunEFunc(f, ecloudInstanceList),
	}

	cmd.Flags().String("name", "", "Instance name for filtering")

	return cmd
}

func ecloudInstanceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	instances, err := service.GetInstances(params)
	if err != nil {
		return fmt.Errorf("Error retrieving instances: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudInstancesProvider(instances))
}

func ecloudInstanceShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <instance: id>...",
		Short:   "Shows a instance",
		Long:    "This command shows one or more instances",
		Example: "ukfast ecloud instance show i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceShow),
	}
}

func ecloudInstanceShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var instances []ecloud.Instance
	for _, arg := range args {
		instance, err := service.GetInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving instance [%s]: %s", arg, err)
			continue
		}

		instances = append(instances, instance)
	}

	return output.CommandOutput(cmd, OutputECloudInstancesProvider(instances))
}

func ecloudInstanceCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an instance",
		Long:    "This command creates an instance",
		Example: "ukfast ecloud instance create --vpc vpc-abcdef12 --az az-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudInstanceCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of instance")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().Int("vcpu", 0, "Number of vCPU cores to allocate")
	cmd.MarkFlagRequired("vcpu")
	cmd.Flags().Int("ram", 0, "Amount of RAM (in MB) to allocate")
	cmd.MarkFlagRequired("ram")
	cmd.Flags().Int("volume", 0, "Size of volume to allocate")
	cmd.MarkFlagRequired("volume")
	cmd.Flags().String("image", "", "ID of image to deploy from")
	cmd.MarkFlagRequired("image")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance has been completely created before continuing on")

	return cmd
}

func ecloudInstanceCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateInstanceRequest{}
	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.VCPUCores, _ = cmd.Flags().GetInt("vcpu")
	createRequest.RAMCapacity, _ = cmd.Flags().GetInt("ram")
	createRequest.VolumeCapacity, _ = cmd.Flags().GetInt("volume")
	createRequest.ImageID, _ = cmd.Flags().GetString("image")

	instanceID, err := service.CreateInstance(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating instance: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(InstanceResourceSyncStatusWaitFunc(service, instanceID, ecloud.SyncStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for instance sync: %s", err)
		}
	}

	instance, err := service.GetInstance(instanceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new instance: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudInstancesProvider([]ecloud.Instance{instance}))
}

func ecloudInstanceUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <instance: id>...",
		Short:   "Updates an instance",
		Long:    "This command updates one or more instances",
		Example: "ukfast ecloud instance update i-abcdef12 --name \"my instance\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceUpdate),
	}

	cmd.Flags().String("name", "", "Name of instance")
	cmd.Flags().Int("vcpu", 0, "Number of vCPU cores to allocate")
	cmd.Flags().Int("ram", 0, "Amount of RAM (in MB) to allocate")

	return cmd
}

func ecloudInstanceUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchInstanceRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	if cmd.Flags().Changed("vcpu") {
		vcpu, _ := cmd.Flags().GetInt("vcpu")
		patchRequest.VCPUCores = vcpu
	}

	if cmd.Flags().Changed("ram") {
		ram, _ := cmd.Flags().GetInt("ram")
		patchRequest.RAMCapacity = ram
	}

	var instances []ecloud.Instance
	for _, arg := range args {
		err := service.PatchInstance(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating instance [%s]: %s", arg, err)
			continue
		}

		instance, err := service.GetInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated instance [%s]: %s", arg, err)
			continue
		}

		instances = append(instances, instance)
	}

	return output.CommandOutput(cmd, OutputECloudInstancesProvider(instances))
}

func ecloudInstanceDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <instance: id...>",
		Short:   "Removes an instance",
		Long:    "This command removes one or more instances",
		Example: "ukfast ecloud instance delete i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceDelete),
	}
}

func ecloudInstanceDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing instance [%s]: %s", arg, err)
		}
	}
	return nil
}

func ecloudInstanceLockCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "lock <instance: id...>",
		Short:   "Locks an instance",
		Long:    "This command locks one or more instances",
		Example: "ukfast ecloud instance lock i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceLock),
	}
}

func ecloudInstanceLock(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.LockInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error locking instance [%s]: %s", arg, err)
		}
	}
	return nil
}

func ecloudInstanceUnlockCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "unlock <instance: id...>",
		Short:   "Unlocks an instance",
		Long:    "This command unlocks one or more instances",
		Example: "ukfast ecloud instance unlock i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceUnlock),
	}
}

func ecloudInstanceUnlock(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.UnlockInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error unlocking instance [%s]: %s", arg, err)
		}
	}
	return nil
}

func ecloudInstanceStartCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "start <instance: id...>",
		Short:   "Starts an instance",
		Long:    "This command powers on one or more instances",
		Example: "ukfast ecloud instance start i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceStart),
	}
}

func ecloudInstanceStart(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.PowerOnInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error starting instance [%s]: %s", arg, err)
		}
	}
	return nil
}

func ecloudInstanceStopCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop <instance: id...>",
		Short:   "Stops an instance",
		Long:    "This command powers off one or more instances",
		Example: "ukfast ecloud instance stop i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceStop),
	}

	cmd.Flags().Bool("force", false, "Specifies that instance should be forcefully powered off")

	return cmd
}

func ecloudInstanceStop(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	for _, arg := range args {
		if force {
			err := service.PowerOffInstance(arg)
			if err != nil {
				output.OutputWithErrorLevelf("Error stopping instance [%s] (forced): %s", arg, err)
			}
		} else {
			err := service.PowerShutdownInstance(arg)
			if err != nil {
				output.OutputWithErrorLevelf("Error stopping instance [%s]: %s", arg, err)
			}
		}
	}
	return nil
}

func ecloudInstanceRestartCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restart <instance: id...>",
		Short:   "Restarts an instance",
		Long:    "This command restarts one or more instances",
		Example: "ukfast ecloud instance restart i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceRestart),
	}

	cmd.Flags().Bool("force", false, "Specifies that instance should be forcefully reset")

	return cmd
}

func ecloudInstanceRestart(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	for _, arg := range args {
		if force {
			err := service.PowerResetInstance(arg)
			if err != nil {
				output.OutputWithErrorLevelf("Error restarting instance [%s] (forced): %s", arg, err)
			}
		} else {
			err := service.PowerRestartInstance(arg)
			if err != nil {
				output.OutputWithErrorLevelf("Error restarting instance [%s]: %s", arg, err)
			}
		}
	}
	return nil
}

func InstanceResourceSyncStatusWaitFunc(service ecloud.ECloudService, instanceID string, status ecloud.SyncStatus) helper.WaitFunc {
	return ResourceSyncStatusWaitFunc(func() (ecloud.SyncStatus, error) {
		instance, err := service.GetInstance(instanceID)
		if err != nil {
			return "", err
		}
		return instance.Sync, nil
	}, status)
}
