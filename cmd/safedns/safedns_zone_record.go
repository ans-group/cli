package safedns

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func safednsZoneRecordRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "sub-commands relating to zone records",
	}

	// Child commands
	cmd.AddCommand(safednsZoneRecordListCmd(f))
	cmd.AddCommand(safednsZoneRecordShowCmd(f))
	cmd.AddCommand(safednsZoneRecordCreateCmd(f))
	cmd.AddCommand(safednsZoneRecordUpdateCmd(f))
	cmd.AddCommand(safednsZoneRecordDeleteCmd(f))

	return cmd
}

func safednsZoneRecordListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <zone: name>",
		Short:   "Lists zone records",
		Long:    "This command lists zone records",
		Example: "ukfast safedns zone record list ukfast.co.uk\nukfast safedns zone record list 123",
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

			return safednsZoneRecordList(c.SafeDNSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Zone record name for filtering")
	cmd.Flags().String("type", "", "Zone record type for filtering")
	cmd.Flags().String("content", "", "Zone record content for filtering")

	return cmd
}

func safednsZoneRecordList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("type", "type"),
		helper.NewStringFilterFlagOption("content", "content"))
	if err != nil {
		return err
	}

	zoneRecords, err := service.GetZoneRecords(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving records for zone: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSRecordsProvider(zoneRecords))
}

func safednsZoneRecordShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: name> <record: id>...",
		Short:   "Shows a zone record",
		Long:    "This command shows one or more zone records",
		Example: "ukfast safedns zone record show ukfast.co.uk 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}
			if len(args) < 2 {
				return errors.New("Missing record")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsZoneRecordShow(c.SafeDNSService(), cmd, args)
		},
	}
}

func safednsZoneRecordShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	var zoneRecords []safedns.Record

	for _, arg := range args[1:] {
		recordID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid record ID [%s]", arg)
			continue
		}

		zoneRecord, err := service.GetZoneRecord(args[0], recordID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving record [%d]: %s", recordID, err)
			continue
		}

		zoneRecords = append(zoneRecords, zoneRecord)
	}

	return output.CommandOutput(cmd, OutputSafeDNSRecordsProvider(zoneRecords))
}

func safednsZoneRecordCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <zone: name>",
		Short:   "Creates a zone record",
		Long:    "This command creates a zone record",
		Example: "ukfast safedns zone record create ukfast.co.uk --name subdomain.ukfast.co.uk --type A --content 1.2.3.4",
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

			return safednsZoneRecordCreate(c.SafeDNSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of record")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("type", "", "Type of record")
	cmd.MarkFlagRequired("type")
	cmd.Flags().String("content", "", "Record content")
	cmd.MarkFlagRequired("content")
	cmd.Flags().Int("priority", 0, "Record priority. Only applicable with MX and SRV type records")

	return cmd
}

func safednsZoneRecordCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	recordName, _ := cmd.Flags().GetString("name")
	recordType, _ := cmd.Flags().GetString("type")
	recordContent, _ := cmd.Flags().GetString("content")

	if strings.ToUpper(recordType) == "TXT" || strings.ToUpper(recordType) == "SPF" {
		recordContent = fmt.Sprintf("%q", strings.Trim(recordContent, "\""))
	}

	createRequest := safedns.CreateRecordRequest{
		Name:    recordName,
		Type:    recordType,
		Content: recordContent,
	}

	if cmd.Flags().Changed("priority") {
		recordPriority, _ := cmd.Flags().GetInt("priority")
		createRequest.Priority = ptr.Int(recordPriority)
	}

	id, err := service.CreateZoneRecord(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating record: %s", err)
	}

	zoneRecord, err := service.GetZoneRecord(args[0], id)
	if err != nil {
		return fmt.Errorf("Error retrieving new record: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSRecordsProvider([]safedns.Record{zoneRecord}))
}

func safednsZoneRecordUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <zone: name> <record: id>...",
		Short:   "Updates a zone record",
		Long:    "This command updates one or more zone records",
		Example: "ukfast safedns zone record update ukfast.co.uk 123 --content 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}
			if len(args) < 2 {
				return errors.New("Missing record")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsZoneRecordUpdate(c.SafeDNSService(), cmd, args)
		},
	}
	cmd.Flags().String("name", "", "Name of record")
	cmd.Flags().String("type", "", "Type of record")
	cmd.Flags().String("content", "", "Record content")
	cmd.Flags().Int("priority", 0, "Record priority. Only applicable with MX type records")

	return cmd
}

func safednsZoneRecordUpdate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	patchRequest := safedns.PatchRecordRequest{}

	if cmd.Flags().Changed("name") {
		recordName, _ := cmd.Flags().GetString("name")
		patchRequest.Name = recordName
	}
	if cmd.Flags().Changed("type") {
		recordType, _ := cmd.Flags().GetString("type")
		patchRequest.Type = recordType
	}
	if cmd.Flags().Changed("content") {
		recordContent, _ := cmd.Flags().GetString("content")
		patchRequest.Content = recordContent
	}
	if cmd.Flags().Changed("priority") {
		recordPriority, _ := cmd.Flags().GetInt("priority")
		patchRequest.Priority = ptr.Int(recordPriority)
	}

	var zoneRecords []safedns.Record

	for _, arg := range args[1:] {
		recordID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid record ID [%s]", arg)
			continue
		}

		id, err := service.PatchZoneRecord(args[0], recordID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating record [%d]: %s", recordID, err)
			continue
		}

		zoneRecord, err := service.GetZoneRecord(args[0], id)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated record [%d]: %s", recordID, err)
			continue
		}

		zoneRecords = append(zoneRecords, zoneRecord)
	}

	return output.CommandOutput(cmd, OutputSafeDNSRecordsProvider(zoneRecords))
}

func safednsZoneRecordDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <zone: name> <record: id>...",
		Short:   "Removes a zone record",
		Long:    "This command removes one or more zone records",
		Example: "ukfast safedns zone record remove ukfast.co.uk 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing zone")
			}
			if len(args) < 2 {
				return errors.New("Missing record")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			safednsZoneRecordDelete(c.SafeDNSService(), cmd, args)
			return nil
		},
	}
}

func safednsZoneRecordDelete(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		recordID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid record ID [%s]", arg)
			continue
		}

		err = service.DeleteZoneRecord(args[0], recordID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing record [%d]: %s", recordID, err)
		}
	}
}
