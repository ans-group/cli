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

func ecloudResourceTierRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resourcetier",
		Short: "sub-commands relating to resource tiers",
	}

	// Child commands
	cmd.AddCommand(ecloudResourceTierListCmd(f))
	cmd.AddCommand(ecloudResourceTierShowCmd(f))

	return cmd
}

func ecloudResourceTierListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists resource tiers",
		Long:    "This command lists resource tiers",
		Example: "ans ecloud resourcetier list",
		RunE:    ecloudCobraRunEFunc(f, ecloudResourceTierList),
	}

	cmd.Flags().String("name", "", "Resource tier name for filtering")
	cmd.Flags().String("az", "", "Availability zone ID for filtering")

	return cmd
}

func ecloudResourceTierList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("az", "availability_zone_id"),
	)
	if err != nil {
		return err
	}

	tiers, err := service.GetResourceTiers(params)
	if err != nil {
		return fmt.Errorf("Error retrieving resource tiers: %s", err)
	}

	return output.CommandOutput(cmd, ResourceTierCollection(tiers))
}

func ecloudResourceTierShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <resourcetier: id>...",
		Short:   "Shows a resource tier",
		Long:    "This command shows one or more resource tiers",
		Example: "ans ecloud resourcetier show rt-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing resource tier")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudResourceTierShow),
	}
}

func ecloudResourceTierShow(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	var tiers []ecloud.ResourceTier
	for _, arg := range args {
		tier, err := service.GetResourceTier(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving resource tier [%s]: %s", arg, err)
			continue
		}

		tiers = append(tiers, tier)
	}

	return output.CommandOutput(cmd, ResourceTierCollection(tiers))
}
