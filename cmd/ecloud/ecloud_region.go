package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudRegionRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "region",
		Short: "sub-commands relating to regions",
	}

	// Child commands
	cmd.AddCommand(ecloudRegionListCmd(f))
	cmd.AddCommand(ecloudRegionShowCmd(f))

	return cmd
}

func ecloudRegionListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists regions",
		Long:    "This command lists regions",
		Example: "ukfast ecloud region list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudRegionList(c.ECloudService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Region name for filtering")

	return cmd
}

func ecloudRegionList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	regions, err := service.GetRegions(params)
	if err != nil {
		return fmt.Errorf("Error retrieving regions: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudRegionsProvider(regions))
}

func ecloudRegionShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <region: id>...",
		Short:   "Shows a region",
		Long:    "This command shows one or more regions",
		Example: "ukfast ecloud region show reg-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing region")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ecloudRegionShow(c.ECloudService(), cmd, args)
		},
	}
}

func ecloudRegionShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var regions []ecloud.Region
	for _, arg := range args {
		region, err := service.GetRegion(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving region [%s]: %s", arg, err)
			continue
		}

		regions = append(regions, region)
	}

	return output.CommandOutput(cmd, OutputECloudRegionsProvider(regions))
}
