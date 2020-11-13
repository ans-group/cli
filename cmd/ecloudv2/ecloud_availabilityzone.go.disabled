package ecloudv2

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudAvailabilityZoneRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "az",
		Short: "sub-commands relating to availability zones",
	}

	// Child commands
	cmd.AddCommand(ecloudAvailabilityZoneListCmd(f))
	cmd.AddCommand(ecloudAvailabilityZoneShowCmd(f))

	return cmd
}

func ecloudAvailabilityZoneListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists availability zones",
		Long:    "This command lists availability zones",
		Example: "ukfast ecloud az list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudAvailabilityZoneList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Availability zone name for filtering")
	cmd.Flags().String("code", "", "Availability zone code for filtering")

	return cmd
}

func ecloudAvailabilityZoneList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd,
		helper.NewStringFilterFlag("name", "name"),
		helper.NewStringFilterFlag("code", "code"),
	)

	azs, err := service.GetAvailabilityZones(params)
	if err != nil {
		return fmt.Errorf("Error retrieving availability zones: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudAvailabilityZonesProvider(azs))
}

func ecloudAvailabilityZoneShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <az: id>...",
		Short:   "Shows an availability zone",
		Long:    "This command shows one or more availability zones",
		Example: "ukfast ecloud az show az-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing availability zone")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudAvailabilityZoneShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudAvailabilityZoneShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var azs []ecloud.AvailabilityZone
	for _, arg := range args {
		az, err := service.GetAvailabilityZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving availability zone [%s]: %s", arg, err)
			continue
		}

		azs = append(azs, az)
	}

	return output.CommandOutput(cmd, OutputECloudAvailabilityZonesProvider(azs))
}
