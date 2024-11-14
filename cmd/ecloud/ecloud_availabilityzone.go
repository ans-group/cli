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

func ecloudAvailabilityZoneRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "availabilityzone",
		Short: "sub-commands relating to availability zones",
		Aliases: []string{
			"az",
		},
	}

	// Child commands
	cmd.AddCommand(ecloudAvailabilityZoneListCmd(f))
	cmd.AddCommand(ecloudAvailabilityZoneShowCmd(f))

	// Child root commands
	cmd.AddCommand(ecloudAvailabilityZoneIOPSRootCmd(f))

	return cmd
}

func ecloudAvailabilityZoneListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists availability zones",
		Long:    "This command lists availability zones",
		Example: "ans ecloud availabilityzone list",
		RunE:    ecloudCobraRunEFunc(f, ecloudAvailabilityZoneList),
	}

	cmd.Flags().String("name", "", "Availability zone name for filtering")
	cmd.Flags().String("region", "", "Region ID for filtering")

	return cmd
}

func ecloudAvailabilityZoneList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("region", "region_id"),
	)
	if err != nil {
		return err
	}

	zones, err := service.GetAvailabilityZones(params)
	if err != nil {
		return fmt.Errorf("Error retrieving availability zones: %s", err)
	}

	return output.CommandOutput(cmd, AvailabilityZoneCollection(zones))
}

func ecloudAvailabilityZoneShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: id>...",
		Short:   "Shows an availability zone",
		Long:    "This command shows one or more availability zones",
		Example: "ans ecloud availabilityzone show az-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing availability zone")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudAvailabilityZoneShow),
	}
}

func ecloudAvailabilityZoneShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var zones []ecloud.AvailabilityZone
	for _, arg := range args {
		zone, err := service.GetAvailabilityZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving availability zone [%s]: %s", arg, err)
			continue
		}

		zones = append(zones, zone)
	}

	return output.CommandOutput(cmd, AvailabilityZoneCollection(zones))
}
