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

func safednsTemplateRecordRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "sub-commands relating to template records",
	}

	// Child commands
	cmd.AddCommand(safednsTemplateRecordListCmd())
	cmd.AddCommand(safednsTemplateRecordShowCmd())
	cmd.AddCommand(safednsTemplateRecordCreateCmd())
	cmd.AddCommand(safednsTemplateRecordUpdateCmd())
	cmd.AddCommand(safednsTemplateRecordDeleteCmd())

	return cmd
}

func safednsTemplateRecordListCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <template: id/name>",
		Short:   "Lists template records",
		Long:    "This command lists template records",
		Example: "ukfast safedns template record list \"main template\"\nukfast safedns template record list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsTemplateRecordList(getClient().SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Template record name for filtering")
	cmd.Flags().String("type", "", "Template record type for filtering")
	cmd.Flags().String("content", "", "Template record content for filtering")

	return cmd
}

func safednsTemplateRecordList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	params, err := GetAPIRequestParametersFromFlags()
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

	templateRecords, err := service.GetTemplateRecords(templateID, params)
	if err != nil {
		output.Fatalf("Error retrieving records for template: %s", err)
	}

	outputSafeDNSRecords(templateRecords)
}

func safednsTemplateRecordShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <template: id/name> <record: id>...",
		Short:   "Shows a template record",
		Long:    "This command shows one or more template records",
		Example: "ukfast safedns template record show \"main template\" 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}
			if len(args) < 2 {
				return errors.New("Missing record")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsTemplateRecordShow(getClient().SafeDNSService(), cmd, args)
		},
	}
}

func safednsTemplateRecordShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	var templateRecords []safedns.Record

	for _, arg := range args[1:] {
		recordID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid record ID [%s]", arg)
			continue
		}

		templateRecord, err := service.GetTemplateRecord(templateID, recordID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving record [%d]: %s", recordID, err)
			continue
		}

		templateRecords = append(templateRecords, templateRecord)
	}

	outputSafeDNSRecords(templateRecords)
}

func safednsTemplateRecordCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <template: id/name>",
		Short:   "Creates a template record",
		Long:    "This command creates a template record",
		Example: "ukfast safedns template record create \"main template\" --name subdomain.ukfast.co.uk --type A --content 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsTemplateRecordCreate(getClient().SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of record")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("type", "", "Type of record")
	cmd.MarkFlagRequired("type")
	cmd.Flags().String("content", "", "Record content")
	cmd.MarkFlagRequired("content")
	cmd.Flags().Int("ttl", 0, "Record priority")
	cmd.MarkFlagRequired("ttl")
	cmd.Flags().Int("priority", 0, "Record priority")

	return cmd
}

func safednsTemplateRecordCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		output.Fatal(err.Error())
		return
	}

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

	id, err := service.CreateTemplateRecord(templateID, createRequest)
	if err != nil {
		output.Fatalf("Error creating record: %s", err)
		return
	}

	templateRecord, err := service.GetTemplateRecord(templateID, id)
	if err != nil {
		output.Fatalf("Error retrieving new record: %s", err)
		return
	}

	outputSafeDNSRecords([]safedns.Record{templateRecord})
}

func safednsTemplateRecordUpdateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <template: id/name> <record: id>...",
		Short:   "Updates a template record",
		Long:    "This command updates one or more template records",
		Example: "ukfast safedns template record update \"main template\" 123 --content 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}
			if len(args) < 2 {
				return errors.New("Missing record")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsTemplateRecordUpdate(getClient().SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of record")
	cmd.Flags().String("type", "", "Type of record")
	cmd.Flags().String("content", "", "Record content")
	cmd.Flags().Int("priority", 0, "Record priority")

	return cmd
}

func safednsTemplateRecordUpdate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		output.Fatal(err.Error())
		return
	}

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

	var templateRecords []safedns.Record

	for _, arg := range args[1:] {
		recordID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid record ID [%s]", arg)
			continue
		}

		_, err = service.PatchTemplateRecord(templateID, recordID, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating record [%d]: %s", recordID, err)
			continue
		}

		templateRecord, err := service.GetTemplateRecord(templateID, recordID)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated record [%d]: %s", recordID, err)
			continue
		}

		templateRecords = append(templateRecords, templateRecord)
	}

	outputSafeDNSRecords(templateRecords)
}

func safednsTemplateRecordDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "delete <template: id/name> <record: id>...",
		Short:   "Removes a template record",
		Long:    "This command removes one or more template records",
		Example: "ukfast safedns template record remove \"main template\" 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}
			if len(args) < 2 {
				return errors.New("Missing record")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			safednsTemplateRecordDelete(getClient().SafeDNSService(), cmd, args)
		},
	}
}

func safednsTemplateRecordDelete(service safedns.SafeDNSService, cmd *cobra.Command, args []string) {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	for _, arg := range args[1:] {
		recordID, err := strconv.Atoi(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Invalid record ID [%s]", arg)
			continue
		}

		err = service.DeleteTemplateRecord(templateID, recordID)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing record [%d]: %s", recordID, err)
		}
	}
}
