package cloudflare

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/cloudflare"
)

func cloudflareZoneRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone",
		Short: "sub-commands relating to zones",
	}

	// Child commands
	cmd.AddCommand(cloudflareZoneListCmd(f))
	cmd.AddCommand(cloudflareZoneShowCmd(f))
	cmd.AddCommand(cloudflareZoneCreateCmd(f))
	cmd.AddCommand(cloudflareZoneDeleteCmd(f))

	return cmd
}

func cloudflareZoneListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists zones",
		Long:    "This command lists zones",
		Example: "ukfast cloudflare zone list",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareZoneList),
	}
}

func cloudflareZoneList(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	zones, err := service.GetZones(params)
	if err != nil {
		return fmt.Errorf("Error retrieving zones: %s", err)
	}

	return output.CommandOutput(cmd, OutputCloudflareZonesProvider(zones))
}

func cloudflareZoneShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: id>...",
		Short:   "Shows a zone",
		Long:    "This command shows one or more zones",
		Example: "ukfast cloudflare zone show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		RunE: cloudflareCobraRunEFunc(f, cloudflareZoneShow),
	}
}

func cloudflareZoneShow(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	var zones []cloudflare.Zone
	for _, arg := range args {
		zone, err := service.GetZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving zone [%s]: %s", arg, err)
			continue
		}

		zones = append(zones, zone)
	}

	return output.CommandOutput(cmd, OutputCloudflareZonesProvider(zones))
}

func cloudflareZoneCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <zone: id>",
		Short:   "Creates a zone",
		Long:    "This command creates a zone",
		Example: "ukfast cloudflare zone create --cluster 123 --default-target-group 456 --name \"test-zone\" --mode http",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareZoneCreate),
	}

	cmd.Flags().String("account", "", "ID of account")
	cmd.MarkFlagRequired("account")
	cmd.Flags().String("name", "", "Name of zone")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("subscription-type", "", "Type of subscription")
	cmd.MarkFlagRequired("subscription-type")

	return cmd
}

func cloudflareZoneCreate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := cloudflare.CreateZoneRequest{}
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

	return output.CommandOutput(cmd, OutputCloudflareZonesProvider([]cloudflare.Zone{zone}))
}

func cloudflareZoneDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <zone: id>...",
		Short:   "Removes a zone",
		Long:    "This command removes one or more zones",
		Example: "ukfast cloudflare zone delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		RunE: cloudflareCobraRunEFunc(f, cloudflareZoneDelete),
	}
}

func cloudflareZoneDelete(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	for _, arg := range args {
		err := service.DeleteZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing zone [%s]: %s", arg, err)
			continue
		}
	}

	return nil
}
