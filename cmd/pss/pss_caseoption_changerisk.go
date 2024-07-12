package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssCaseOptionChangeRiskRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changerisk",
		Short: "sub-commands relating to change risk case options",
	}

	// Child commands
	cmd.AddCommand(pssCaseOptionChangeRiskListCmd(f))

	return cmd
}

func pssCaseOptionChangeRiskListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists change risk case options",
		Long:    "This command lists change risk case options",
		Example: "ans pss request list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssCaseOptionChangeRiskList(c.PSSService(), cmd, args)
		},
	}
}

func pssCaseOptionChangeRiskList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	requests, err := service.GetChangeRiskCaseOptions(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSCaseOptionsProvider(requests))
}
