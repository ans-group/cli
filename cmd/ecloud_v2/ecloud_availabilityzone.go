package ecloud_v2

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
		Use:   "vpc",
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

	return cmd
}

func ecloudAvailabilityZoneList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	if cmd.Flags().Changed("name") {
		filterName, _ := cmd.Flags().GetString("name")
		params.WithFilter(helper.GetFilteringInferOperator("name", filterName))
	}

	vpcs, err := service.GetAvailabilityZones(params)
	if err != nil {
		return fmt.Errorf("Error retrieving vpcs: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudAvailabilityZonesProvider(vpcs))
}

func ecloudAvailabilityZoneShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <vpc: id>...",
		Short:   "Shows a vpc",
		Long:    "This command shows one or more vpcs",
		Example: "ukfast ecloud vpc show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing vpc")
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
	var vpcs []ecloud.AvailabilityZone
	for _, arg := range args {
		vpc, err := service.GetAvailabilityZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving vpc [%s]: %s", arg, err)
			continue
		}

		vpcs = append(vpcs, vpc)
	}

	return output.CommandOutput(cmd, OutputECloudAvailabilityZonesProvider(vpcs))
}
