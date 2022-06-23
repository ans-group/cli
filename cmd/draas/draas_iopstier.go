package draas

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/draas"
	"github.com/spf13/cobra"
)

func draasIOPSTierRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iopstier",
		Short: "sub-commands relating to IOPS tiers",
	}

	// Child commands
	cmd.AddCommand(draasIOPSTierListCmd(f))
	cmd.AddCommand(draasIOPSTierShowCmd(f))

	return cmd
}

func draasIOPSTierListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists available IOPS tiers",
		Long:    "This command lists available IOPS tiers",
		Example: "ans draas iopstier list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasIOPSTierList(c.DRaaSService(), cmd, args)
		},
	}
}

func draasIOPSTierList(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	iopstiers, err := service.GetIOPSTiers(params)
	if err != nil {
		return fmt.Errorf("Error retrieving IOPS tiers: %s", err)
	}

	return output.CommandOutput(cmd, OutputDRaaSIOPSTiersProvider(iopstiers))
}

func draasIOPSTierShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <iopstier: id>...",
		Short:   "Shows an IOPS tier",
		Long:    "This command shows one or more IOPS tiers",
		Example: "ans draas iopstier show 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing IOPS tier")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return draasIOPSTierShow(c.DRaaSService(), cmd, args)
		},
	}
}

func draasIOPSTierShow(service draas.DRaaSService, cmd *cobra.Command, args []string) error {
	var iopstiers []draas.IOPSTier
	for _, arg := range args {
		iopstier, err := service.GetIOPSTier(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving IOPS tier [%s]: %s", arg, err)
			continue
		}

		iopstiers = append(iopstiers, iopstier)
	}

	return output.CommandOutput(cmd, OutputDRaaSIOPSTiersProvider(iopstiers))
}
