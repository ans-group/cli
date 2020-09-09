package threatmonitoring

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func ThreatMonitoringRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "threatmonitoring",
		Short: "Commands relating to Threat Monitoring service",
	}

	// Child root commands
	cmd.AddCommand(threatmonitoringAgentRootCmd(f))
	cmd.AddCommand(threatmonitoringAlertRootCmd(f))

	return cmd
}
