package ecloud

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/connection"
	"github.com/ukfast/sdk-go/pkg/ptr"
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
	cmd.AddCommand(ecloudInstanceStartCmd(f))
	cmd.AddCommand(ecloudInstanceStopCmd(f))
	cmd.AddCommand(ecloudInstanceRestartCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudInstanceVolumeRootCmd(f))
	cmd.AddCommand(ecloudInstanceCredentialRootCmd(f))
	cmd.AddCommand(ecloudInstanceNICRootCmd(f))
	cmd.AddCommand(ecloudInstanceConsoleRootCmd(f))
	cmd.AddCommand(ecloudInstanceTaskRootCmd(f))

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
		Example: "ukfast ecloud instance create --vpc vpc-abcdef12 --vcpu 2 --ram 2048 --volume 20 --image \"CentOS 7\"",
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
	cmd.Flags().String("network", "", "ID of network to use for instance")
	cmd.MarkFlagRequired("network")
	cmd.Flags().String("image", "", "ID or name of image to deploy from")
	cmd.MarkFlagRequired("image")
	cmd.Flags().StringSlice("ssh-key-pair", []string{}, "ID of SSH key pair, can be repeated")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance has been completely created")

	return cmd
}

func ecloudInstanceCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateInstanceRequest{}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.VCPUCores, _ = cmd.Flags().GetInt("vcpu")
	createRequest.RAMCapacity, _ = cmd.Flags().GetInt("ram")
	createRequest.VolumeCapacity, _ = cmd.Flags().GetInt("volume")
	createRequest.NetworkID, _ = cmd.Flags().GetString("network")

	if cmd.Flags().Changed("name") {
		createRequest.Name, _ = cmd.Flags().GetString("name")
	}
	if cmd.Flags().Changed("ssh-key-pair") {
		createRequest.SSHKeyPairIDs, _ = cmd.Flags().GetStringSlice("ssh-key-pair")
	}

	imageFlag, _ := cmd.Flags().GetString("image")

	if strings.HasPrefix(imageFlag, "img-") {
		createRequest.ImageID = imageFlag
	} else {
		images, err := service.GetImages(connection.APIRequestParameters{
			Filtering: []connection.APIRequestFiltering{
				{
					Property: "name",
					Operator: connection.EQOperator,
					Value:    []string{imageFlag},
				},
			},
		})
		if err != nil {
			return fmt.Errorf("Error retrieving images: %s", err)
		}

		if len(images) != 1 {
			return fmt.Errorf("Expected 1 image, got %d images", len(images))
		}

		createRequest.ImageID = images[0].ID
		if err != nil {
			return err
		}
	}

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
	cmd.Flags().String("volume-group", "", "ID of volume-group to use for instance")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance has been completely updated")

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

	if cmd.Flags().Changed("volume-group") {
		volGroup, _ := cmd.Flags().GetString("volume-group")
		patchRequest.VolumeGroupID = ptr.String(volGroup)
	}

	var instances []ecloud.Instance
	for _, arg := range args {
		err := service.PatchInstance(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating instance [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(InstanceResourceSyncStatusWaitFunc(service, arg, ecloud.SyncStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for instance [%s] sync: %s", arg, err)
				continue
			}
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
	cmd := &cobra.Command{
		Use:     "delete <instance: id>...",
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

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance has been completely removed")

	return cmd
}

func ecloudInstanceDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing instance [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(InstanceNotFoundWaitFunc(service, arg))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for removal of instance [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}

func ecloudInstanceLockCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "lock <instance: id>...",
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
		Use:     "unlock <instance: id>...",
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
	cmd := &cobra.Command{
		Use:     "start <instance: id>...",
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

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance power on task has been completed")

	return cmd
}

func ecloudInstanceStart(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.PowerOnInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error starting instance [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for instance [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}

func ecloudInstanceStopCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "stop <instance: id>...",
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
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance power off task has been completed")

	return cmd
}

func ecloudInstanceStop(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	for _, arg := range args {
		var taskID string
		var err error
		if force {
			taskID, err = service.PowerOffInstance(arg)
			if err != nil {
				output.OutputWithErrorLevelf("Error stopping instance [%s] (forced): %s", arg, err)
				continue
			}
		} else {
			taskID, err = service.PowerShutdownInstance(arg)
			if err != nil {
				output.OutputWithErrorLevelf("Error stopping instance [%s]: %s", arg, err)
				continue
			}
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for instance [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}

func ecloudInstanceRestartCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "restart <instance: id>...",
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
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance restart task has been completed")

	return cmd
}

func ecloudInstanceRestart(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	force, _ := cmd.Flags().GetBool("force")

	for _, arg := range args {
		var taskID string
		var err error
		if force {
			taskID, err = service.PowerResetInstance(arg)
			if err != nil {
				output.OutputWithErrorLevelf("Error restarting instance [%s] (forced): %s", arg, err)
				continue
			}
		} else {
			taskID, err = service.PowerRestartInstance(arg)
			if err != nil {
				output.OutputWithErrorLevelf("Error restarting instance [%s]: %s", arg, err)
				continue
			}
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for instance [%s]: %s", arg, err)
				continue
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
		return instance.Sync.Status, nil
	}, status)
}

func InstanceNotFoundWaitFunc(service ecloud.ECloudService, instanceID string) helper.WaitFunc {
	return func() (finished bool, err error) {
		_, err = service.GetInstance(instanceID)
		if err != nil {
			switch err.(type) {
			case *ecloud.InstanceNotFoundError:
				return true, nil
			default:
				return false, fmt.Errorf("Failed to retrieve instance [%s]: %s", instanceID, err)
			}
		}

		return false, nil
	}
}
