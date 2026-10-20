package ecloud

import (
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/spf13/cobra"
)

var (
	instanceTagColumns        = []string{"tags"}
	instanceBackupColumns     = []string{"backup_enabled", "backup_agent_enabled", "backup_gateway_id"}
	instanceMonitoringColumns = []string{"monitoring_enabled", "monitoring_gateway_id"}
)

// instanceOutputHelp documents the optional instance output columns, for use in command help text
const instanceOutputHelp = `
Additional columns can be included in table output with the following flags:
  --with-tags        tags
  --with-backup      backup_enabled, backup_agent_enabled, backup_gateway_id
  --with-monitoring  monitoring_enabled, monitoring_gateway_id
  --with-all         all of the above

These properties can also be selected directly with the --property flag, e.g.
  --property id,name,backup_enabled,backup_agent_enabled,backup_gateway_id`

// addInstanceOutputFlags registers flags for including optional columns in instance output
func addInstanceOutputFlags(cmd *cobra.Command) {
	cmd.Flags().Bool("with-tags", false, "Include tags column in output")
	cmd.Flags().Bool("with-backup", false, "Include backup columns (backup_enabled, backup_agent_enabled, backup_gateway_id) in output")
	cmd.Flags().Bool("with-monitoring", false, "Include monitoring columns (monitoring_enabled, monitoring_gateway_id) in output")
	cmd.Flags().Bool("with-all", false, "Include tags, backup and monitoring columns in output")
}

// instanceAdditionalColumns returns the optional instance columns requested via flags
func instanceAdditionalColumns(cmd *cobra.Command) []string {
	withAll, _ := cmd.Flags().GetBool("with-all")

	var columns []string
	for _, group := range []struct {
		flag    string
		columns []string
	}{
		{"with-tags", instanceTagColumns},
		{"with-backup", instanceBackupColumns},
		{"with-monitoring", instanceMonitoringColumns},
	} {
		if enabled, _ := cmd.Flags().GetBool(group.flag); enabled || withAll {
			columns = append(columns, group.columns...)
		}
	}

	return columns
}

// instanceOutputOptions returns output handler options for the optional instance columns requested via flags
func instanceOutputOptions(cmd *cobra.Command) []output.OutputHandlerOption {
	columns := instanceAdditionalColumns(cmd)
	if len(columns) == 0 {
		return nil
	}

	return []output.OutputHandlerOption{output.WithAdditionalColumns(columns...)}
}
