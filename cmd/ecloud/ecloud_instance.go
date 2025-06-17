package ecloud

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
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
	cmd.AddCommand(ecloudInstanceSSHCmd(f))
	cmd.AddCommand(ecloudInstanceMigrateCmd(f))
	cmd.AddCommand(ecloudInstanceEncryptCmd(f))
	cmd.AddCommand(ecloudInstanceDecryptCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudInstanceVolumeRootCmd(f))
	cmd.AddCommand(ecloudInstanceCredentialRootCmd(f))
	cmd.AddCommand(ecloudInstanceNICRootCmd(f))
	cmd.AddCommand(ecloudInstanceConsoleRootCmd(f))
	cmd.AddCommand(ecloudInstanceTaskRootCmd(f))
	cmd.AddCommand(ecloudInstanceFloatingIPRootCmd(f))
	cmd.AddCommand(ecloudInstanceImageRootCmd(f))

	return cmd
}

func ecloudInstanceListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists instances",
		Long:    "This command lists instances",
		Example: "ans ecloud instance list",
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

	return output.CommandOutput(cmd, InstanceCollection(instances))
}

func ecloudInstanceShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <instance: id>...",
		Short:   "Shows a instance",
		Long:    "This command shows one or more instances",
		Example: "ans ecloud instance show i-abcdef12",
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

	return output.CommandOutput(cmd, InstanceCollection(instances))
}

func ecloudInstanceCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates an instance",
		Long:    "This command creates an instance",
		Example: "ans ecloud instance create --vpc vpc-abcdef12 --network net-abcdef12 --vcpu-sockets 2 --vcpu-cores-per-socket 2 --ram 2048 --volume 20 --image \"CentOS 7\"",
		RunE:    ecloudCobraRunEFunc(f, ecloudInstanceCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of instance")
	cmd.Flags().String("vpc", "", "ID of VPC")
	cmd.MarkFlagRequired("vpc")
	cmd.Flags().Int("vcpu", 0, "Number of vCPU sockets to allocate")
	cmd.Flags().MarkDeprecated("vcpu", "use --vcpu-sockets / --vcpu-cores-per-socket flags instead")
	cmd.Flags().Int("vcpu-sockets", 1, "Number of vCPU sockets to allocate")
	cmd.Flags().Int("vcpu-cores-per-socket", 1, "Number of vCPU cores to allocate per socket")
	cmd.Flags().Int("ram", 0, "Amount of RAM (in MB) to allocate")
	cmd.MarkFlagRequired("ram")
	cmd.Flags().Int("volume", 0, "Size of volume to allocate")
	cmd.MarkFlagRequired("volume")
	cmd.Flags().String("network", "", "ID of network to use for instance")
	cmd.MarkFlagRequired("network")
	cmd.Flags().String("image", "", "ID or name of image to deploy from")
	cmd.MarkFlagRequired("image")
	cmd.Flags().StringSlice("ssh-key-pair", []string{}, "ID of SSH key pair, can be repeated")
	cmd.Flags().String("host-group", "", "ID of host group to deploy to")
	cmd.Flags().String("resource-tier", "", "ID of resource tier to deploy to. A default tier is chosen if not specified")
	cmd.Flags().String("ip-address", "", "IP address to allocate for DHCP")
	cmd.Flags().Bool("enable-vm-backups", false, "Enable VM-level backups")
	cmd.Flags().String("backup-gateway-id", "", "Backup gateway ID, enables agent-level backups")
	cmd.Flags().Bool("enable-monitoring", false, "Enable monitoring")
	cmd.Flags().String("monitoring-gateway-id", "", "Monitoring gateway ID")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance has been completely created")

	return cmd
}

func ecloudInstanceCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateInstanceRequest{}
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.RAMCapacity, _ = cmd.Flags().GetInt("ram")
	createRequest.VolumeCapacity, _ = cmd.Flags().GetInt("volume")
	createRequest.NetworkID, _ = cmd.Flags().GetString("network")
	createRequest.HostGroupID, _ = cmd.Flags().GetString("host-group")
	createRequest.ResourceTierID, _ = cmd.Flags().GetString("resource-tier")
	createRequest.BackupEnabled, _ = cmd.Flags().GetBool("enable-vm-backups")
	createRequest.BackupGatewayID, _ = cmd.Flags().GetString("backup-gateway-id")
	createRequest.MonitoringEnabled, _ = cmd.Flags().GetBool("enable-monitoring")
	createRequest.MonitoringGatewayID, _ = cmd.Flags().GetString("monitoring-gateway-id")
	createRequest.Name, _ = cmd.Flags().GetString("name")

	if cmd.Flags().Changed("vcpu") {
		if cmd.Flags().Changed("vcpu-sockets") || cmd.Flags().Changed("vcpu-cores-per-socket") {
			return fmt.Errorf("Flag --vcpu is mutually exclusive with --vcpu-sockets and --vcpu-cores-per-socket")
		}
		createRequest.VCPUCores, _ = cmd.Flags().GetInt("vcpu")
	} else {
		createRequest.VCPUSockets, _ = cmd.Flags().GetInt("vcpu-sockets")
		createRequest.VCPUCoresPerSocket, _ = cmd.Flags().GetInt("vcpu-cores-per-socket")
	}

	if cmd.Flags().Changed("ssh-key-pair") {
		createRequest.SSHKeyPairIDs, _ = cmd.Flags().GetStringSlice("ssh-key-pair")
	}

	if cmd.Flags().Changed("ip-address") {
		ipAddress, _ := cmd.Flags().GetString("ip-address")
		createRequest.CustomIPAddress = connection.IPAddress(ipAddress)
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

		if len(images) == 0 {
			return fmt.Errorf("Image not found with name '%s'", imageFlag)
		}

		if len(images) > 1 {
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

	return output.CommandOutput(cmd, InstanceCollection([]ecloud.Instance{instance}))
}

func ecloudInstanceUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <instance: id>...",
		Short:   "Updates an instance",
		Long:    "This command updates one or more instances",
		Example: "ans ecloud instance update i-abcdef12 --name \"my instance\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceUpdate),
	}

	cmd.Flags().String("name", "", "Name of instance")
	cmd.Flags().Int("vcpu", 0, "Number of vCPU sockets to allocate")
	cmd.Flags().MarkDeprecated("vcpu", "use --vcpu-sockets / --vcpu-cores-per-socket flags instead")
	cmd.Flags().Int("vcpu-sockets", 0, "Number of vCPU sockets to allocate")
	cmd.Flags().Int("vcpu-cores-per-socket", 0, "Number of vCPU cores to allocate per socket")
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
		if cmd.Flags().Changed("vcpu-sockets") || cmd.Flags().Changed("vcpu-cores-per-socket") {
			return fmt.Errorf("Flag --vcpu is mutually exclusive with --vcpu-sockets and --vcpu-cores-per-socket")
		}
		patchRequest.VCPUCores, _ = cmd.Flags().GetInt("vcpu")
	} else {
		if cmd.Flags().Changed("vcpu-sockets") {
			patchRequest.VCPUSockets, _ = cmd.Flags().GetInt("vcpu-sockets")
		}

		if cmd.Flags().Changed("vcpu-cores-per-socket") {
			patchRequest.VCPUCoresPerSocket, _ = cmd.Flags().GetInt("vcpu-cores-per-socket")
		}
	}

	if cmd.Flags().Changed("ram") {
		patchRequest.RAMCapacity, _ = cmd.Flags().GetInt("ram")
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

	return output.CommandOutput(cmd, InstanceCollection(instances))
}

func ecloudInstanceDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <instance: id>...",
		Short:   "Removes an instance",
		Long:    "This command removes one or more instances",
		Example: "ans ecloud instance delete i-abcdef12",
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
		Example: "ans ecloud instance lock i-abcdef12",
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
		Example: "ans ecloud instance unlock i-abcdef12",
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
		Example: "ans ecloud instance start i-abcdef12",
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
		Example: "ans ecloud instance stop i-abcdef12",
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
		Example: "ans ecloud instance restart i-abcdef12",
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

func ecloudInstanceSSHCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "ssh <instance: id>",
		Short:   "Invokes SSH for an instance",
		Long:    "This command invokes SSH for an instance",
		Example: "ans ecloud instance ssh i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceSSH),
	}

	cmd.Flags().Int("port", 2020, "Specifies port to connect to")
	cmd.Flags().Bool("internal", false, "Specifies internal IP should be used")
	cmd.Flags().String("user", "root", "Specifies user to connect with")
	cmd.Flags().String("args", "", "Specifies additional arguments to pass to SSH")

	return cmd
}

func ecloudInstanceSSH(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var ipAddress string
	internal, _ := cmd.Flags().GetBool("internal")
	if internal {
		nics, err := service.GetInstanceNICs(args[0], connection.APIRequestParameters{})
		if err != nil {
			return fmt.Errorf("Error retrieving instance NICs: %s", err)
		}

		if len(nics) < 1 {
			return fmt.Errorf("No floating IPs found for instance")
		}

		ipAddress = nics[0].IPAddress
	} else {
		fips, err := service.GetInstanceFloatingIPs(args[0], connection.APIRequestParameters{})
		if err != nil {
			return fmt.Errorf("Error retrieving instance floating IPs: %s", err)
		}

		if len(fips) < 1 {
			return fmt.Errorf("No floating IPs found for instance")
		}

		ipAddress = fips[0].IPAddress
	}

	user, _ := cmd.Flags().GetString("user")
	port, _ := cmd.Flags().GetInt("port")
	sshArgs, _ := cmd.Flags().GetString("args")

	sshCmd := exec.Command("ssh", fmt.Sprintf("%s@%s", user, ipAddress), "-p", strconv.Itoa(port), sshArgs)
	sshCmd.Stdout = os.Stdout
	sshCmd.Stdin = os.Stdin
	sshCmd.Stderr = os.Stderr

	sshCmd.Start()
	return sshCmd.Wait()
}

func ecloudInstanceMigrateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "migrate <instance: id>",
		Short:   "Migrates an instance",
		Long:    "This command migrates an instance to another resource tier or dedicated host group",
		Example: "ans ecloud instance migrate i-abcdef12 --resource-tier rt-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		PreRun: func(cmd *cobra.Command, args []string) {
			hostgroup, _ := cmd.Flags().GetString("host-group")
			if hostgroup == "" {
				cmd.MarkFlagRequired("resource-tier")
			}
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceMigrate),
	}

	cmd.Flags().String("resource-tier", "", "Specifies the resource-tier to migrate the instance to")
	cmd.Flags().String("host-group", "", "Specifies the host-group to migrate the instance to")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance migrate task has been completed")

	return cmd
}

func ecloudInstanceMigrate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	migrateRequest := ecloud.MigrateInstanceRequest{}

	if cmd.Flags().Changed("resource-tier") {
		migrateRequest.ResourceTierID, _ = cmd.Flags().GetString("resource-tier")
	}
	if cmd.Flags().Changed("host-group") {
		migrateRequest.HostGroupID, _ = cmd.Flags().GetString("host-group")
	}

	for _, arg := range args {
		taskID, err := service.MigrateInstance(arg, migrateRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error migrating instance [%s]: %s", arg, err)
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

func ecloudInstanceEncryptCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "encrypt <instance: id>",
		Short:   "Encrypts an instance",
		Long:    "This command encrypts an instance.",
		Example: "ans ecloud instance encrypt i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceEncrypt),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance encrypt task has been completed")

	return cmd
}

func ecloudInstanceEncrypt(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.EncryptInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error encrypting instance [%s]: %s", arg, err)
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

func ecloudInstanceDecryptCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "decrypt <instance: id>",
		Short:   "Decrypts an encrypted instance",
		Long:    "This command decrypts a previously encrypted instance..",
		Example: "ans ecloud instance decrypt i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceDecrypt),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the instance decrypt task has been completed")

	return cmd
}

func ecloudInstanceDecrypt(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DecryptInstance(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error decrypting instance [%s]: %s", arg, err)
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
