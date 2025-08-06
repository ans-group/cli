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

func ecloudAvailabilityZoneIOPSRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "iops",
		Short: "sub-commands relating to availability zone IOPS tiers",
		Aliases: []string{
			"az",
		},
	}

	// Child commands
	cmd.AddCommand(ecloudAvailabilityZoneIOPSTierListCmd(f))

	return cmd
}

func ecloudAvailabilityZoneIOPSTierListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists availability zone IOPS tiers",
		Long:    "This command lists availability zone IOPS tiers",
		Example: "ans ecloud availabilityzone iops list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing availability zone")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudAvailabilityZoneIOPSTierList),
	}

	cmd.Flags().String("name", "", "IOPS name for filtering")

	return cmd
}

func ecloudAvailabilityZoneIOPSTierList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
	)
	if err != nil {
		return err
	}

	zones, err := service.GetAvailabilityZoneIOPSTiers(args[0], params)
	if err != nil {
		return fmt.Errorf("error retrieving availability zone IOPS tiers: %s", err)
	}

	return output.CommandOutput(cmd, IOPSTierCollection(zones))
}
