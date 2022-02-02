package managedcloudflare

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/managedcloudflare"
)

func managedcloudflareZoneRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone",
		Short: "sub-commands relating to zones",
	}

	// Child commands
	cmd.AddCommand(managedcloudflareZoneListCmd(f))
	cmd.AddCommand(managedcloudflareZoneShowCmd(f))

	return cmd
}

func managedcloudflareZoneListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists zones",
		Long:    "This command lists zones",
		Example: "ukfast managedcloudflare zone list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return managedcloudflareZoneList(c.ManagedCloudflareService(), cmd, args)
		},
	}
}

func managedcloudflareZoneList(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	zones, err := service.GetZones(params)
	if err != nil {
		return fmt.Errorf("Error retrieving zones: %s", err)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareZonesProvider(zones))
}

func managedcloudflareZoneShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: id>...",
		Short:   "Shows a zone",
		Long:    "This command shows one or more zones",
		Example: "ukfast managedcloudflare zone show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return managedcloudflareZoneShow(c.ManagedCloudflareService(), cmd, args)
		},
	}
}

func managedcloudflareZoneShow(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	var zones []managedcloudflare.Zone
	for _, arg := range args {
		zone, err := service.GetZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving zone [%s]: %s", arg, err)
			continue
		}

		zones = append(zones, zone)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareZonesProvider(zones))
}
