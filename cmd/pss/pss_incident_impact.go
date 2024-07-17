package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssIncidentImpactRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "impact",
		Short: "sub-commands relating to incident impact case options",
	}

	// Child commands
	cmd.AddCommand(pssIncidentImpactListCmd(f))

	return cmd
}

func pssIncidentImpactListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists incident impact case options",
		Long:    "This command lists incident impact case options",
		Example: "ans pss request list",
		RunE:    pssCobraRunEFunc(f, pssIncidentImpactList),
	}
}

func pssIncidentImpactList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	options, err := service.GetIncidentImpactCaseOptions(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSCaseOptionsProvider(options))
}
