package safedns

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/spf13/cobra"
)

func safednsZoneRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "zone",
		Short: "sub-commands relating to zones",
	}

	// Child commands
	cmd.AddCommand(safednsZoneListCmd(f))
	cmd.AddCommand(safednsZoneShowCmd(f))
	cmd.AddCommand(safednsZoneCreateCmd(f))
	cmd.AddCommand(safednsZoneUpdateCmd(f))
	cmd.AddCommand(safednsZoneDeleteCmd(f))

	// Child root commands
	cmd.AddCommand(safednsZoneRecordRootCmd(f))

	return cmd
}

func safednsZoneListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists zones",
		Long:    "This command lists zones",
		Example: "ans safedns zone list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsZoneList(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Zone name for filtering")

	return cmd
}

func safednsZoneList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd, helper.NewStringFilterFlagOption("name", "name"))
	if err != nil {
		return err
	}

	zones, err := service.GetZones(params)
	if err != nil {
		return fmt.Errorf("Error retrieving zones: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSZonesProvider(zones))
}

func safednsZoneShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: name>...",
		Short:   "Shows a zone",
		Long:    "This command shows one or more zones",
		Example: "ans safedns zone show ans.co.uk\nans safedns zone show 123",
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

			return safednsZoneShow(c.SafeDNSService(), cmd, args)
		},
	}
}

func safednsZoneShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	var zones []safedns.Zone
	for _, arg := range args {
		zone, err := service.GetZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving zone [%s]: %s", arg, err)
			continue
		}

		zones = append(zones, zone)
	}

	return output.CommandOutput(cmd, OutputSafeDNSZonesProvider(zones))
}

func safednsZoneCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create",
		Short:   "Creates a zone",
		Long:    "This command creates a zone",
		Example: "ans safedns zone create",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsZoneCreate(c.SafeDNSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of zone")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("description", "", "Description for zone")

	return cmd
}

func safednsZoneCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	description, _ := cmd.Flags().GetString("description")

	createRequest := safedns.CreateZoneRequest{
		Name:        name,
		Description: description,
	}

	err := service.CreateZone(createRequest)
	if err != nil {
		return fmt.Errorf("Error creating zone: %s", err)
	}

	zone, err := service.GetZone(name)
	if err != nil {
		return fmt.Errorf("Error retrieving new zone: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSZonesProvider([]safedns.Zone{zone}))
}

func safednsZoneUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <zone: name>...",
		Short:   "Updates a zone",
		Long:    "This command updates one or more zones",
		Example: "ans safedns zone update ans.co.uk --description \"some description\"",
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

			return safednsZoneUpdate(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("description", "", "Description for zone")

	return cmd
}

func safednsZoneUpdate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	patchRequest := safedns.PatchZoneRequest{}
	patchRequest.Description, _ = cmd.Flags().GetString("description")

	var zones []safedns.Zone
	for _, arg := range args {
		err := service.PatchZone(arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating zone [%s]: %s", arg, err)
			continue
		}

		zone, err := service.GetZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated zone [%s]: %s", arg, err)
			continue
		}

		zones = append(zones, zone)
	}

	return output.CommandOutput(cmd, OutputSafeDNSZonesProvider(zones))
}

func safednsZoneDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <zone: name>...",
		Short:   "Removes a zone",
		Long:    "This command removes one or more zones",
		Example: "ans safedns zone delete ans.co.uk\nans safedns zone delete 123",
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

			safednsZoneDelete(c.SafeDNSService(), cmd, args)
			return nil
		},
	}
}

func safednsZoneDelete(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	for _, arg := range args {
		err := service.DeleteZone(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing zone [%s]: %s", arg, err)
		}
	}
}
