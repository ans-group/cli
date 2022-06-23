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
		Example: "ans ecloud region list",
		RunE:    ecloudCobraRunEFunc(f, ecloudRegionList),
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
		Example: "ans ecloud region show reg-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing region")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudRegionShow),
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
