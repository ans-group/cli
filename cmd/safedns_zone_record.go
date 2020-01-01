package cmd

import (
	"errors"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/ptr"
	"github.com/ukfast/sdk-go/pkg/service/safedns"
)

func safednsZoneRecordRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "sub-commands relating to zone records",
	}

	// Child commands
	cmd.AddCommand(safednsZoneRecordListCmd())
	cmd.AddCommand(safednsZoneRecordShowCmd())
	cmd.AddCommand(safednsZoneRecordCreateCmd())
	cmd.AddCommand(safednsZoneRecordUpdateCmd())
	cmd.AddCommand(safednsZoneRecordDeleteCmd())

	return cmd
}
func safednsZoneRecordListCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneRecordList(getClient().SafeDNSService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Zone record name for filtering")
	cmd.Flags().String("type", "", "Zone record type for filtering")
	cmd.Flags().String("content", "", "Zone record content for filtering")

	return cmd
}

func safednsZoneRecordList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	if cmd.Flags().Changed("name") {
		filterName, _ := cmd.Flags().GetString("name")
		params.WithFilter(helper.GetFilteringInferOperator("name", filterName))
	}

	if cmd.Flags().Changed("type") {
		filterType, _ := cmd.Flags().GetString("type")
		params.WithFilter(helper.GetFilteringInferOperator("type", filterType))
	}

	if cmd.Flags().Changed("content") {
		filterContent, _ := cmd.Flags().GetString("content")
		params.WithFilter(helper.GetFilteringInferOperator("content", filterContent))
	}

	zoneRecords, err := service.GetZoneRecords(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving records for zone: %s", err)
		return
	}

	outputSafeDNSRecords(zoneRecords)
}

func safednsZoneRecordShowCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneRecordShow(getClient().SafeDNSService(), cmd, args)
		},
	}
}

func safednsZoneRecordShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
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

	outputSafeDNSRecords(zoneRecords)
}

func safednsZoneRecordCreateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneRecordCreate(getClient().SafeDNSService(), cmd, args)
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

func safednsZoneRecordCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	recordName, _ := cmd.Flags().GetString("name")
	recordType, _ := cmd.Flags().GetString("type")
	recordContent, _ := cmd.Flags().GetString("content")

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
		output.Fatalf("Error creating record: %s", err)
		return
	}

	zoneRecord, err := service.GetZoneRecord(args[0], id)
	if err != nil {
		output.Fatalf("Error retrieving new record: %s", err)
		return
	}

	outputSafeDNSRecords([]safedns.Record{zoneRecord})
}

func safednsZoneRecordUpdateCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneRecordUpdate(getClient().SafeDNSService(), cmd, args)
		},
	}
	cmd.Flags().String("name", "", "Name of record")
	cmd.Flags().String("type", "", "Type of record")
	cmd.Flags().String("content", "", "Record content")
	cmd.Flags().Int("priority", 0, "Record priority. Only applicable with MX type records")

	return cmd
}

func safednsZoneRecordUpdate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
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

	outputSafeDNSRecords(zoneRecords)
}

func safednsZoneRecordDeleteCmd() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			safednsZoneRecordDelete(getClient().SafeDNSService(), cmd, args)
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
