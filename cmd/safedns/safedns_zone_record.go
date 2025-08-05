package safedns

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/spf13/cobra"
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
		Example: "ans safedns zone record list ans.co.uk\nans safedns zone record list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing zone")
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
		return fmt.Errorf("error retrieving records for zone: %s", err)
	}

	return output.CommandOutput(cmd, RecordCollection(zoneRecords))
}

func safednsZoneRecordShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <zone: name> <record: id>...",
		Short:   "Shows a zone record",
		Long:    "This command shows one or more zone records",
		Example: "ans safedns zone record show ans.co.uk 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing zone")
			}
			if len(args) < 2 {
				return errors.New("missing record")
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

	return output.CommandOutput(cmd, RecordCollection(zoneRecords))
}

func safednsZoneRecordCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <zone: name>",
		Short:   "Creates a zone record",
		Long:    "This command creates a zone record",
		Example: "ans safedns zone record create ans.co.uk --name subdomain.ans.co.uk --type A --content 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing zone")
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
	_ = cmd.MarkFlagRequired("name")
	cmd.Flags().String("type", "", "Type of record")
	_ = cmd.MarkFlagRequired("type")
	cmd.Flags().String("content", "", "Record content")
	_ = cmd.MarkFlagRequired("content")
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
		return fmt.Errorf("error creating record: %s", err)
	}

	zoneRecord, err := service.GetZoneRecord(args[0], id)
	if err != nil {
		return fmt.Errorf("error retrieving new record: %s", err)
	}

	return output.CommandOutput(cmd, RecordCollection([]safedns.Record{zoneRecord}))
}

func safednsZoneRecordUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <zone: name> <record: id>...",
		Short:   "Updates a zone record",
		Long:    "This command updates one or more zone records",
		Example: "ans safedns zone record update ans.co.uk 123 --content 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing zone")
			}
			if len(args) < 2 {
				return errors.New("missing record")
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

	return output.CommandOutput(cmd, RecordCollection(zoneRecords))
}

func safednsZoneRecordDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <zone: name> <record: id>...",
		Short:   "Removes a zone record",
		Long:    "This command deletes one or more zone records",
		Example: "ans safedns zone record delete ans.co.uk 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing zone")
			}
			if len(args) < 2 {
				return errors.New("missing record")
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
