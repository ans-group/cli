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

func threatmonitoringAgentRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agent",
		Short: "sub-commands relating to agents",
	}

	// Child commands
	cmd.AddCommand(threatmonitoringAgentListCmd(f))
	cmd.AddCommand(threatmonitoringAgentShowCmd(f))

	return cmd
}

func threatmonitoringAgentListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists agents",
		Long:    "This command lists agents",
		Example: "ukfast threatmonitoring agent list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return threatmonitoringAgentList(c.ThreatMonitoringService(), cmd, args)
		},
	}
}

func threatmonitoringAgentList(service threatmonitoring.ThreatMonitoringService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	agents, err := service.GetAgents(params)
	if err != nil {
		return fmt.Errorf("Error retrieving agents: %s", err)
	}

	return output.CommandOutput(cmd, OutputThreatMonitoringAgentsProvider(agents))
}

func threatmonitoringAgentShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <agent: id>...",
		Short:   "Shows a agent",
		Long:    "This command shows one or more agents",
		Example: "ukfast threatmonitoring agent show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing agent")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return threatmonitoringAgentShow(c.ThreatMonitoringService(), cmd, args)
		},
	}
}

func threatmonitoringAgentShow(service threatmonitoring.ThreatMonitoringService, cmd *cobra.Command, args []string) error {
	var agents []threatmonitoring.Agent
	for _, arg := range args {
		agent, err := service.GetAgent(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving agent [%s]: %s", arg, err)
			continue
		}

		agents = append(agents, agent)
	}

	return output.CommandOutput(cmd, OutputThreatMonitoringAgentsProvider(agents))
}
