package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/storage"
)

func storageVolumeRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "volume",
		Short: "sub-commands relating to volumes",
	}

	// Child commands
	cmd.AddCommand(storageVolumeListCmd())
	cmd.AddCommand(storageVolumeShowCmd())

	return cmd
}

func storageVolumeListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists volumes",
		Long:    "This command lists volumes",
		Example: "ukfast storage volume list",
		Run: func(cmd *cobra.Command, args []string) {
			storageVolumeList(getClient().StorageService(), cmd, args)
		},
	}
}

func storageVolumeList(service storage.StorageService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	volumes, err := service.GetVolumes(params)
	if err != nil {
		output.Fatalf("Error retrieving volumes: %s", err)
		return
	}

	outputStorageVolumes(volumes)
}

func storageVolumeShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			storageVolumeShow(getClient().StorageService(), cmd, args)
		},
	}
}

func storageVolumeShow(service storage.StorageService, cmd *cobra.Command, args []string) {
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

	outputStorageVolumes(volumes)
}
