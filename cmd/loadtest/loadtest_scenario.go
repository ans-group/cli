package loadtest

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestScenarioRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "scenario",
		Short: "sub-commands relating to scenarios",
	}

	// Child commands
	cmd.AddCommand(loadtestScenarioListCmd(f))

	return cmd
}

func loadtestScenarioListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists scenarios",
		Long:    "This command lists scenarios",
		Example: "ukfast loadtest scenario list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestScenarioList(f.NewClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestScenarioList(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	scenarios, err := service.GetScenarios(params)
	if err != nil {
		return fmt.Errorf("Error retrieving scenarios: %s", err)
	}

	return output.CommandOutput(cmd, OutputLoadTestScenariosProvider(scenarios))
}
