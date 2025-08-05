package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudBackupGatewayRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "backupgateway",
		Short: "sub-commands relating to backup gateways",
	}

	// Child commands
	cmd.AddCommand(ecloudBackupGatewaySpecificationRootCmd(f))
	cmd.AddCommand(ecloudBackupGatewayListCmd(f))
	cmd.AddCommand(ecloudBackupGatewayShowCmd(f))
	cmd.AddCommand(ecloudBackupGatewayCreateCmd(f))
	cmd.AddCommand(ecloudBackupGatewayUpdateCmd(f))
	cmd.AddCommand(ecloudBackupGatewayDeleteCmd(f))

	return cmd
}

func ecloudBackupGatewayListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists backup gateways",
		Example: "ans ecloud backupgateway list",
		RunE:    ecloudCobraRunEFunc(f, ecloudBackupGatewayList),
	}

	cmd.Flags().String("name", "", "Backup gateway name for filtering")

	return cmd
}

func ecloudBackupGatewayList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	gateways, err := service.GetBackupGateways(params)
	if err != nil {
		return fmt.Errorf("error retrieving backup gateways: %s", err)
	}

	return output.CommandOutput(cmd, BackupGatewayCollection(gateways))
}

func ecloudBackupGatewayShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <gateway: id>...",
		Short:   "Show details of a backup gateway",
		Example: "ans ecloud backupgateway show bgw-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing backup gateway ID")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudBackupGatewayShow),
	}
}

func ecloudBackupGatewayShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var backupGateways []ecloud.BackupGateway
	for _, arg := range args {
		backupGateway, err := service.GetBackupGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving backup gateway [%s]: %s", arg, err)
			continue
		}

		backupGateways = append(backupGateways, backupGateway)
	}

	return output.CommandOutput(cmd, BackupGatewayCollection(backupGateways))
}

func ecloudBackupGatewayCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a backup gateway",
		Example: "ans ecloud backupgateway create --router rtr-abcdef12 --vpc vpc-abcd1234 --specification bgws-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudBackupGatewayCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of gateway")
	cmd.Flags().String("vpc", "", "ID of VPC")
	_ = cmd.MarkFlagRequired("vpc")
	cmd.Flags().String("router", "", "ID of router")
	_ = cmd.MarkFlagRequired("router")
	cmd.Flags().String("specification", "", "ID of backup gateway specification")
	_ = cmd.MarkFlagRequired("specification")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the backup gateway has been completely created")

	return cmd
}

func ecloudBackupGatewayCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateBackupGatewayRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.VPCID, _ = cmd.Flags().GetString("vpc")
	createRequest.RouterID, _ = cmd.Flags().GetString("router")
	createRequest.GatewaySpecID, _ = cmd.Flags().GetString("specification")

	taskRef, err := service.CreateBackupGateway(createRequest)
	if err != nil {
		return fmt.Errorf("error creating backup gateway: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("error waiting for backup gateway task to complete: %s", err)
		}
	}

	backupGateway, err := service.GetBackupGateway(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("error retrieving new backup gateway: %s", err)
	}

	return output.CommandOutput(cmd, BackupGatewayCollection([]ecloud.BackupGateway{backupGateway}))
}

func ecloudBackupGatewayUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <gateway: id>...",
		Short:   "Updates a backup gateway",
		Long:    "Update the name of a backup gateway",
		Example: "ans ecloud backupgateway update bgw-abcdef12 --name \"my gateway\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing backup gateway")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudBackupGatewayUpdate),
	}

	cmd.Flags().String("name", "", "Name of gateway")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the backup gateway has been completely updated")

	return cmd
}

func ecloudBackupGatewayUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchBackupGatewayRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var backupGateways []ecloud.BackupGateway
	for _, arg := range args {
		task, err := service.PatchBackupGateway(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating backup gateway [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for backup gateway [%s]: %s", arg, err)
				continue
			}
		}

		backupGateway, err := service.GetBackupGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated backup gateway [%s]: %s", arg, err)
			continue
		}

		backupGateways = append(backupGateways, backupGateway)
	}

	return output.CommandOutput(cmd, BackupGatewayCollection(backupGateways))
}

func ecloudBackupGatewayDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <gateway: id>...",
		Short:   "Removes a backup gateway",
		Example: "ans ecloud backupgateway delete bgw-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing backup gateway")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudBackupGatewayDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the backup gateway has been completely removed")

	return cmd
}

func ecloudBackupGatewayDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteBackupGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing backup gateway [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for backup gateway [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
