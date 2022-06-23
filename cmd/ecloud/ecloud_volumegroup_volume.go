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

func ecloudVolumeGroupVolumeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "sub-commands relating to volumegroup volumes",
	}

	// Child commands
	cmd.AddCommand(ecloudVolumeGroupVolumeListCmd(f))

	return cmd
}

func ecloudVolumeGroupVolumeListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists volumegroup volumes",
		Long:    "This command lists volumegroup volumes",
		Example: "ukfast ecloud volumegroup volume list volgroup-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume-group")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVolumeGroupVolumeList),
	}

	cmd.Flags().String("name", "", "Volume name for filtering")

	return cmd
}

func ecloudVolumeGroupVolumeList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	volumes, err := service.GetVolumeGroupVolumes(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving volume-group volumes: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVolumesProvider(volumes))
}
