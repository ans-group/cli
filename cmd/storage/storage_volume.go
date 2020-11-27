package storage

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	flaghelper "github.com/ukfast/cli/internal/pkg/helper/flag"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/storage"
)

func storageVolumeRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "sub-commands relating to volumes",
	}

	// Child commands
	cmd.AddCommand(storageVolumeListCmd(f))
	cmd.AddCommand(storageVolumeShowCmd(f))

	return cmd
}

func storageVolumeListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists volumes",
		Long:    "This command lists volumes",
		Example: "ukfast storage volume list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return storageVolumeList(c.StorageService(), cmd, args)
		},
	}
}

func storageVolumeList(service storage.StorageService, cmd *cobra.Command, args []string) error {
	params, err := flaghelper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	volumes, err := service.GetVolumes(params)
	if err != nil {
		return fmt.Errorf("Error retrieving volumes: %s", err)
	}

	return output.CommandOutput(cmd, OutputStorageVolumesProvider(volumes))
}

func storageVolumeShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <volume: id>...",
		Short:   "Shows a volume",
		Long:    "This command shows one or more volumes",
		Example: "ukfast storage volume show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing volume")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return storageVolumeShow(c.StorageService(), cmd, args)
		},
	}
}

func storageVolumeShow(service storage.StorageService, cmd *cobra.Command, args []string) error {
	var volumes []storage.Volume
	for _, arg := range args {
		volumeID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid volume ID [%s]", arg)
			continue
		}

		volume, err := service.GetVolume(volumeID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving volume [%s]: %s", arg, err)
			continue
		}

		volumes = append(volumes, volume)
	}

	return output.CommandOutput(cmd, OutputStorageVolumesProvider(volumes))
}
