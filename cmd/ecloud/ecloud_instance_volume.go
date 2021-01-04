package ecloud

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ecloud"
)

func ecloudInstanceVolumeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "sub-commands relating to instance volumes",
	}

	// Child commands
	cmd.AddCommand(ecloudInstanceVolumeListCmd(f))

	return cmd
}

func ecloudInstanceVolumeListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists instance volumes",
		Long:    "This command lists instance volumes",
		Example: "ukfast ecloud instance volume list i-abcdef12",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing instance")
			}

			return nil
		},
		RunE: ecloudCobraRunEFunc(f, ecloudInstanceVolumeList),
	}

	cmd.Flags().String("name", "", "Volume name for filtering")

	return cmd
}

func ecloudInstanceVolumeList(service ecloud.ECloudService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	flaghelper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, flaghelper.NewStringFilterFlag("name", "name"))

	volumes, err := service.GetInstanceVolumes(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving instance volumes: %s", err)
	}

	return output.CommandOutput(cmd, OutputECloudVolumesProvider(volumes))
}
