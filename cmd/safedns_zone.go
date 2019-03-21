package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func safednsZoneRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone",
		Short: "sub-commands relating to zones",
	}

	// Child commands
	cmd.AddCommand(safednsZoneListCmd())
	cmd.AddCommand(safednsZoneShowCmd())
	cmd.AddCommand(safednsZoneCreateCmd())
	cmd.AddCommand(safednsZoneDeleteCmd())

	// Child root commands
	cmd.AddCommand(safednsZoneRecordRootCmd())

	return cmd
}

func safednsZoneListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists zones",
		Long:    "This command lists zones",
		Example: "ukfast safedns zone list",
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneList(getClient().SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Zone name for filtering")

	return cmd
}

func safednsZoneList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	if cmd.Flags().Changed("name") {
		filterName, _ := cmd.Flags().GetString("name")
		params.WithFilter(helper.GetFilteringInferOperator("name", filterName))
	}

	zones, err := service.GetZones(params)
	if err != nil {
		output.Fatalf("Error retrieving zones: %s", err)
		return
	}

	outputSafeDNSZones(zones)
}

func safednsZoneShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: name>...",
		Short:   "Shows a zone",
		Long:    "This command shows one or more zones",
		Example: "ukfast safedns zone show ukfast.co.uk\nukfast safedns zone show 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneShow(getClient().SafeDNSService(), cmd, args)
		},
	}
}

func safednsZoneShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	var zones []safedns.Zone
	for _, arg := range args {
		zone, err := service.GetZone(arg)
		if err != nil {
			OutputWithErrorLevelf("Error retrieving zone [%s]: %s", arg, err)
			continue
		}

		zones = append(zones, zone)
	}

	outputSafeDNSZones(zones)
}

func safednsZoneCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a zone",
		Long:    "This command creates a zone",
		Example: "ukfast safedns zone create",
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneCreate(getClient().SafeDNSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of zone")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("description", "", "Description for zone")

	return cmd
}

func safednsZoneCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")

	createRequest := safedns.CreateZoneRequest{
		Name:        name,
		Description: description,
	}

	err := service.CreateZone(createRequest)
	if err != nil {
		output.Fatalf("Error creating zone: %s", err)
		return
	}

	zone, err := service.GetZone(name)
	if err != nil {
		output.Fatalf("Error retrieving new zone: %s", err)
		return
	}

	outputSafeDNSZones([]safedns.Zone{zone})
}

func safednsZoneDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <zone: name...>",
		Short:   "Removes a zone",
		Long:    "This command removes one or more zones",
		Example: "ukfast safedns zone delete ukfast.co.uk\nukfast safedns zone delete 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneDelete(getClient().SafeDNSService(), cmd, args)
		},
	}
}

func safednsZoneDelete(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteZone(arg)
		if err != nil {
			OutputWithErrorLevelf("Error removing zone [%s]: %s", arg, err)
		}
	}
}
