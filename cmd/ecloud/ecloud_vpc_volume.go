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

func ecloudVPCVolumeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "sub-commands relating to VPC volumes",
	}

	// Child commands
	cmd.AddCommand(ecloudVPCVolumeListCmd(f))

	return cmd
}

func ecloudVPCVolumeListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists VPC volumes",
		Long:    "This command lists VPC volumes",
		Example: "ukfast ecloud vpc volume list vpc-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing VPC")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudVPCVolumeList),
	}

	cmd.Flags().String("name", "", "Volume name for filtering")

	return cmd
}

func ecloudVPCVolumeList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	volumes, err := service.GetVPCVolumes(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving VPC volumes: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVolumesProvider(volumes))
}
