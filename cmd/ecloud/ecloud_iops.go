package ecloud

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ecloud"
	"github.com/spf13/cobra"
)

func ecloudIOPSRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iops",
		Short: "sub-commands relating to IOPS tiers",
	}

	// Child commands
	cmd.AddCommand(ecloudIOPSTierListCmd(f))
	cmd.AddCommand(ecloudIOPSTierShowCmd(f))

	return cmd
}

func ecloudIOPSTierListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists IOPS tiers",
		Long:    "This command lists IOPS tiers",
		Example: "ans ecloud iops list",
		RunE:    ecloudCobraRunEFunc(f, ecloudIOPSTierList),
	}

	cmd.Flags().String("name", "", "IOPS name for filtering")

	return cmd
}

func ecloudIOPSTierList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	tiers, err := service.GetIOPSTiers(params)
	if err != nil {
		return fmt.Errorf("Error retrieving IOPS tiers: %s", err)
	}

	return output.CommandOutput(cmd, IOPSTierCollection(tiers))
}

func ecloudIOPSTierShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <iops: id>...",
		Short:   "Shows an IOPS tier",
		Long:    "This command shows one or more IOPS tiers",
		Example: "ans ecloud IOPS show iops-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing IOPS tier")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudIOPSTierShow),
	}
}

func ecloudIOPSTierShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var tiersList []ecloud.IOPSTier
	for _, arg := range args {
		tiers, err := service.GetIOPSTier(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving IOPS tier [%s]: %s", arg, err)
			continue
		}

		tiersList = append(tiersList, tiers)
	}

	return output.CommandOutput(cmd, IOPSTierCollection(tiersList))
}
