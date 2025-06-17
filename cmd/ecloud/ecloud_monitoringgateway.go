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
