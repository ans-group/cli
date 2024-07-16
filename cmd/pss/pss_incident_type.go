package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssIncidentTypeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "type",
		Short: "sub-commands relating to incident type case options",
	}

	// Child commands
	cmd.AddCommand(pssIncidentTypeListCmd(f))

	return cmd
}

func pssIncidentTypeListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists incident type case options",
		Long:    "This command lists incident type case options",
		Example: "ans pss request list",
		RunE:    pssCobraRunEFunc(f, pssIncidentTypeList),
	}
}

func pssIncidentTypeList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	options, err := service.GetIncidentTypeCaseOptions(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSCaseOptionsProvider(options))
}
