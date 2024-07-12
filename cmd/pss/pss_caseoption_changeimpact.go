package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssCaseOptionChangeImpactRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "changeimpact",
		Short: "sub-commands relating to change impact case options",
	}

	// Child commands
	cmd.AddCommand(pssCaseOptionChangeImpactListCmd(f))

	return cmd
}

func pssCaseOptionChangeImpactListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists change impact case options",
		Long:    "This command lists change impact case options",
		Example: "ans pss request list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssCaseOptionChangeImpactList(c.PSSService(), cmd, args)
		},
	}
}

func pssCaseOptionChangeImpactList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	requests, err := service.GetChangeImpactCaseOptions(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSCaseOptionsProvider(requests))
}
