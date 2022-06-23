package cloudflare

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/cloudflare"
	"github.com/spf13/cobra"
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
	cmd.AddCommand(cloudflareZoneUpdateCmd(f))
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
		Example: "ukfast cloudflare zone show 00000000-0000-0000-0000-000000000000",
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
		Example: "ukfast cloudflare zone create --account 621e88d4-c401-4063-bdcf-07ca3c09efed --name \"test-zone\" --subscription a144257d-df53-414e-a44d-3dd84ac90395",
		RunE:    cloudflareCobraRunEFunc(f, cloudflareZoneCreate),
	}

	cmd.Flags().String("account", "", "ID of account")
	cmd.MarkFlagRequired("account")
	cmd.Flags().String("name", "", "Name of zone")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("subscription", "", "ID of plan subscription")
	cmd.MarkFlagRequired("subscription")

	return cmd
}

func cloudflareZoneCreate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	createRequest := cloudflare.CreateZoneRequest{}
	createRequest.AccountID, _ = cmd.Flags().GetString("account")
	createRequest.Name, _ = cmd.Flags().GetString("name")
	createRequest.SubscriptionID, _ = cmd.Flags().GetString("subscription")

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

func cloudflareZoneUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <zone: id>...",
		Short:   "Removes a zone",
		Long:    "This command removes one or more zones",
		Example: "ukfast cloudflare zone update 83d70af6-80ba-4463-abda-2880613efbc1",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		RunE: cloudflareCobraRunEFunc(f, cloudflareZoneUpdate),
	}

	cmd.Flags().String("subscription", "", "ID of plan subscription")

	return cmd
}

func cloudflareZoneUpdate(service cloudflare.CloudflareService, cmd *cobra.Command, args []string) error {
	req := cloudflare.PatchZoneRequest{}
	req.SubscriptionID, _ = cmd.Flags().GetString("subscription")

	for _, arg := range args {
		err := service.PatchZone(arg, req)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating zone [%s]: %s", arg, err)
			continue
		}
	}

	return nil
}

func cloudflareZoneDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <zone: id>...",
		Short:   "Removes a zone",
		Long:    "This command removes one or more zones",
		Example: "ukfast cloudflare zone delete 1c3081b2-d65e-41d1-8077-c86f21759366",
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
