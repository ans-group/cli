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
	cmd.AddCommand(managedcloudflareZoneCreateCmd(f))
	cmd.AddCommand(managedcloudflareZoneDeleteCmd(f))

	return cmd
}

func managedcloudflareZoneListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists zones",
		Long:    "This command lists zones",
		Example: "ukfast managedcloudflare zone list",
		RunE:    managedcloudflareCobraRunEFunc(f, managedcloudflareZoneList),
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
		RunE: managedcloudflareCobraRunEFunc(f, managedcloudflareZoneShow),
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

func managedcloudflareZoneCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <zone: id>",
		Short:   "Creates a zone",
		Long:    "This command creates a zone",
		Example: "ukfast managedcloudflare zone create --cluster 123 --default-target-group 456 --name \"test-zone\" --mode http",
		RunE:    managedcloudflareCobraRunEFunc(f, managedcloudflareZoneCreate),
	}

	cmd.Flags().String("account", "", "ID of account")
	cmd.MarkFlagRequired("account")
	cmd.Flags().String("name", "", "Name of zone")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("subscription-type", "", "Type of subscription")
	cmd.MarkFlagRequired("subscription-type")

	return cmd
}

func managedcloudflareZoneCreate(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := managedcloudflare.CreateZoneRequest{}
	createRequest.AccountID, _ = cmd.Flags().GetString("account")
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.SubscriptionType, _ = cmd.Flags().GetString("subscription-type")

	zoneID, err := service.CreateZone(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating zone: %s", err)
	}

	zone, err := service.GetZone(zoneID)
	if err != nil {
		return fmt.Errorf("Error retrieving new zone: %s", err)
	}

	return output.CommandOutput(cmd, OutputManagedCloudflareZonesProvider([]managedcloudflare.Zone{zone}))
}

func managedcloudflareZoneDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <zone: id>...",
		Short:   "Removes a zone",
		Long:    "This command removes one or more zones",
		Example: "ukfast managedcloudflare zone delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		RunE: managedcloudflareCobraRunEFunc(f, managedcloudflareZoneDelete),
	}
}

func managedcloudflareZoneDelete(service managedcloudflare.ManagedCloudflareService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing zone [%s]: %s", arg, err)
			continue
		}
	}

	return nil
}
