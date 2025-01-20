package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssChangeRiskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "risk",
		Short: "sub-commands relating to change risk case options",
	}

	// Child commands
	cmd.AddCommand(pssChangeRiskListCmd(f))

	return cmd
}

func pssChangeRiskListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists change risk case options",
		Long:    "This command lists change risk case options",
		Example: "ans pss request list",
		RunE:    pssCobraRunEFunc(f, pssChangeRiskList),
	}
}

func pssChangeRiskList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	options, err := service.GetChangeRiskCaseOptions(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, CaseOptionCollection(options))
}
