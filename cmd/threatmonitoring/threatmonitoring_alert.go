package threatmonitoring

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/threatmonitoring"
)

func threatmonitoringAlertRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "alert",
		Short: "sub-commands relating to alerts",
	}

	// Child commands
	cmd.AddCommand(threatmonitoringAlertListCmd(f))
	cmd.AddCommand(threatmonitoringAlertShowCmd(f))

	return cmd
}

func threatmonitoringAlertListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists alerts",
		Long:    "This command lists alerts, paginated",
		Example: "ukfast threatmonitoring alert list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return threatmonitoringAlertList(c.ThreatMonitoringService(), cmd, args)
		},
	}
}

func threatmonitoringAlertList(service threatmonitoring.ThreatMonitoringService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	alerts, err := service.GetAlertsPaginated(params)
	if err != nil {
		return fmt.Errorf("Error retrieving alerts: %s", err)
	}

	return output.CommandOutputPaginated(cmd, OutputThreatMonitoringAlertsProvider(alerts.Items), alerts)
}

func threatmonitoringAlertShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <alert: id>...",
		Short:   "Shows a alert",
		Long:    "This command shows one or more alerts",
		Example: "ukfast threatmonitoring alert show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing alert")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return threatmonitoringAlertShow(c.ThreatMonitoringService(), cmd, args)
		},
	}
}

func threatmonitoringAlertShow(service threatmonitoring.ThreatMonitoringService, cmd *cobra.Command, args []string) error {
	var alerts []threatmonitoring.Alert
	for _, arg := range args {
		alert, err := service.GetAlert(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving alert [%s]: %s", arg, err)
			continue
		}

		alerts = append(alerts, alert)
	}

	return output.CommandOutput(cmd, OutputThreatMonitoringAlertsProvider(alerts))
}
