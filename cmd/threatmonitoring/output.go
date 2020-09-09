package threatmonitoring

import (
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/threatmonitoring"
)

func OutputThreatMonitoringAgentsProvider(agents []threatmonitoring.Agent) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(agents).WithDefaultFields([]string{"id", "friendly_name", "status", "platform", "created_at"})
}

func OutputThreatMonitoringAlertsProvider(alerts []threatmonitoring.Alert) output.OutputHandlerProvider {
	return output.NewSerializedOutputHandlerProvider(alerts).WithDefaultFields([]string{"id", "agent_id", "agent_friendly_name", "level", "description", "timestamp"})
}
