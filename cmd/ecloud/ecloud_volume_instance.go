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

func ecloudVolumeInstanceRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "instance",
		Short: "sub-commands relating to volume instances",
	}

	// Child commands
	cmd.AddCommand(ecloudVolumeInstanceListCmd(f))

	return cmd
}

func ecloudVolumeInstanceListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists volume instances",
		Long:    "This command lists volume instances",
		Example: "ans ecloud volume instance list vol-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeInstanceList),
	}

	cmd.Flags().String("name", "", "Instance name for filtering")

	return cmd
}

func ecloudVolumeInstanceList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	instances, err := service.GetVolumeInstances(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving volume instances: %s", err)
	}

	return output.CommandOutput(cmd, InstanceCollection(instances))
}
