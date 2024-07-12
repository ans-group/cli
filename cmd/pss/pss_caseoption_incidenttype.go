package pss

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/pss"
	"github.com/spf13/cobra"
)

func pssCaseOptionIncidentTypeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "incidenttype",
		Short: "sub-commands relating to incident type case options",
	}

	// Child commands
	cmd.AddCommand(pssCaseOptionIncidentTypeListCmd(f))

	return cmd
}

func pssCaseOptionIncidentTypeListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists incident type case options",
		Long:    "This command lists incident type case options",
		Example: "ans pss request list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return pssCaseOptionIncidentTypeList(c.PSSService(), cmd, args)
		},
	}
}

func pssCaseOptionIncidentTypeList(service pss.PSSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	requests, err := service.GetIncidentTypeCaseOptions(params)
	if err != nil {
		return err
	}

	return output.CommandOutput(cmd, OutputPSSCaseOptionsProvider(requests))
}
