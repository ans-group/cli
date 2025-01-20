package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssChangeImpactRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "impact",
		Short: "sub-commands relating to change impact case options",
	}

	// Child commands
	cmd.AddCommand(pssChangeImpactListCmd(f))

	return cmd
}

func pssChangeImpactListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists change impact case options",
		Long:    "This command lists change impact case options",
		Example: "ans pss request list",
		RunE:    pssCobraRunEFunc(f, pssChangeImpactList),
	}
}

func pssChangeImpactList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	options, err := service.GetChangeImpactCaseOptions(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, CaseOptionCollection(options))
}
