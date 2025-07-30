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

func ecloudMonitoringGatewayRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "monitoringgateway",
		Short: "sub-commands relating to monitoring gateways",
	}

	// Child commands
	cmd.AddCommand(ecloudMonitoringGatewayListCmd(f))
	cmd.AddCommand(ecloudMonitoringGatewayShowCmd(f))
	cmd.AddCommand(ecloudMonitoringGatewayCreateCmd(f))
	cmd.AddCommand(ecloudMonitoringGatewayUpdateCmd(f))
	cmd.AddCommand(ecloudMonitoringGatewayDeleteCmd(f))

	return cmd
}

func ecloudMonitoringGatewayListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists monitoring gateways",
		Example: "ans ecloud monitoringgateway list",
		RunE:    ecloudCobraRunEFunc(f, ecloudMonitoringGatewayList),
	}

	cmd.Flags().String("name", "", "Monitoring gateway name for filtering")

	return cmd
}

func ecloudMonitoringGatewayList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	gateways, err := service.GetMonitoringGateways(params)
	if err != nil {
		return fmt.Errorf("Error retrieving monitoring gateways: %s", err)
	}

	return output.CommandOutput(cmd, MonitoringGatewayCollection(gateways))
}

func ecloudMonitoringGatewayShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <gateway: id>...",
		Short:   "Show details of a monitoring gateway",
		Example: "ans ecloud monitoringgateway show mgw-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing monitoring gateway ID")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudMonitoringGatewayShow),
	}
}

func ecloudMonitoringGatewayShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var monitoringGateways []ecloud.MonitoringGateway
	for _, arg := range args {
		monitoringGateway, err := service.GetMonitoringGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving monitoring gateway [%s]: %s", arg, err)
			continue
		}

		monitoringGateways = append(monitoringGateways, monitoringGateway)
	}

	return output.CommandOutput(cmd, MonitoringGatewayCollection(monitoringGateways))
}

func ecloudMonitoringGatewayCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a monitoring gateway",
		Example: "ans ecloud monitoringgateway create --router rtr-abcdef12 --specification mgws-abcdef12",
		RunE:    ecloudCobraRunEFunc(f, ecloudMonitoringGatewayCreate),
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of gateway")
	cmd.Flags().String("router", "", "ID of router")
	cmd.MarkFlagRequired("router")
	cmd.Flags().String("specification", "", "ID of monitoring gateway specification")
	cmd.MarkFlagRequired("specification")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the monitoring gateway has been completely created")

	return cmd
}

func ecloudMonitoringGatewayCreate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	createRequest := ecloud.CreateMonitoringGatewayRequest{}
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.RouterID, _ = cmd.Flags().GetString("router")
	createRequest.SpecID, _ = cmd.Flags().GetString("specification")

	taskRef, err := service.CreateMonitoringGateway(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating monitoring gateway: %s", err)
	}

	waitFlag, _ := cmd.Flags().GetBool("wait")
	if waitFlag {
		err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskRef.TaskID, ecloud.TaskStatusComplete))
		if err != nil {
			return fmt.Errorf("Error waiting for monitoring gateway task to complete: %s", err)
		}
	}

	monitoringGateway, err := service.GetMonitoringGateway(taskRef.ResourceID)
	if err != nil {
		return fmt.Errorf("Error retrieving new monitoring gateway: %s", err)
	}

	return output.CommandOutput(cmd, MonitoringGatewayCollection([]ecloud.MonitoringGateway{monitoringGateway}))
}

func ecloudMonitoringGatewayUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <gateway: id>...",
		Short:   "Updates a monitoring gateway",
		Long:    "Update the name of a monitoring gateway",
		Example: "ans ecloud monitoringgateway update mgw-abcdef12 --name \"my gateway\"",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing monitoring gateway")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudMonitoringGatewayUpdate),
	}

	cmd.Flags().String("name", "", "Name of gateway")
	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the monitoring gateway has been completely updated")

	return cmd
}

func ecloudMonitoringGatewayUpdate(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	patchRequest := ecloud.PatchMonitoringGatewayRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}

	var monitoringGateways []ecloud.MonitoringGateway
	for _, arg := range args {
		task, err := service.PatchMonitoringGateway(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating monitoring gateway [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, task.TaskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for monitoring gateway [%s]: %s", arg, err)
				continue
			}
		}

		monitoringGateway, err := service.GetMonitoringGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated monitoring gateway [%s]: %s", arg, err)
			continue
		}

		monitoringGateways = append(monitoringGateways, monitoringGateway)
	}

	return output.CommandOutput(cmd, MonitoringGatewayCollection(monitoringGateways))
}

func ecloudMonitoringGatewayDeleteCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "delete <gateway: id>...",
		Short:   "Removes a monitoring gateway",
		Example: "ans ecloud monitoringgateway delete mgw-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing monitoring gateway")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudMonitoringGatewayDelete),
	}

	cmd.Flags().Bool("wait", false, "Specifies that the command should wait until the monitoring gateway has been completely removed")

	return cmd
}

func ecloudMonitoringGatewayDelete(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		taskID, err := service.DeleteMonitoringGateway(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing monitoring gateway [%s]: %s", arg, err)
			continue
		}

		waitFlag, _ := cmd.Flags().GetBool("wait")
		if waitFlag {
			err := helper.WaitForCommand(TaskStatusWaitFunc(service, taskID, ecloud.TaskStatusComplete))
			if err != nil {
				output.OutputWithErrorLevelf("Error waiting for task to complete for monitoring gateway [%s]: %s", arg, err)
				continue
			}
		}
	}
	return nil
}
