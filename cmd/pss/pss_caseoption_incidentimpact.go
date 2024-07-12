package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssCaseOptionIncidentImpactRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "incidentimpact",
		Short: "sub-commands relating to incident impact case options",
	}

	// Child commands
	cmd.AddCommand(pssCaseOptionIncidentImpactListCmd(f))

	return cmd
}

func pssCaseOptionIncidentImpactListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists incident impact case options",
		Long:    "This command lists incident impact case options",
		Example: "ans pss request list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssCaseOptionIncidentImpactList(c.PSSService(), cmd, args)
		},
	}
}

func pssCaseOptionIncidentImpactList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	requests, err := service.GetIncidentImpactCaseOptions(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSCaseOptionsProvider(requests))
}
