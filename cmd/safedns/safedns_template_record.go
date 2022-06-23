package safedns

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/ptr"
	"github.com/ans-group/sdk-go/pkg/service/safedns"
	"github.com/spf13/cobra"
)

func safednsTemplateRecordRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "sub-commands relating to template records",
	}

	// Child commands
	cmd.AddCommand(safednsTemplateRecordListCmd(f))
	cmd.AddCommand(safednsTemplateRecordShowCmd(f))
	cmd.AddCommand(safednsTemplateRecordCreateCmd(f))
	cmd.AddCommand(safednsTemplateRecordUpdateCmd(f))
	cmd.AddCommand(safednsTemplateRecordDeleteCmd(f))

	return cmd
}

func safednsTemplateRecordListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <template: id/name>",
		Short:   "Lists template records",
		Long:    "This command lists template records",
		Example: "ans safedns template record list \"main template\"\nans safedns template record list 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsTemplateRecordList(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Template record name for filtering")
	cmd.Flags().String("type", "", "Template record type for filtering")
	cmd.Flags().String("content", "", "Template record content for filtering")

	return cmd
}

func safednsTemplateRecordList(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		return err
	}

	params, err := helper.GetAPIRequestParametersFromFlags(cmd,
		helper.NewStringFilterFlagOption("name", "name"),
		helper.NewStringFilterFlagOption("type", "type"),
		helper.NewStringFilterFlagOption("content", "content"))
	if err != nil {
		return err
	}

	templateRecords, err := service.GetTemplateRecords(templateID, params)
	if err != nil {
		return fmt.Errorf("Error retrieving records for template: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSRecordsProvider(templateRecords))
}

func safednsTemplateRecordShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <template: id/name> <record: id>...",
		Short:   "Shows a template record",
		Long:    "This command shows one or more template records",
		Example: "ans safedns template record show \"main template\" 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
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

			return safednsTemplateRecordShow(c.SafeDNSService(), cmd, args)
		},
	}
}

func safednsTemplateRecordShow(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		return err
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

	return output.CommandOutput(cmd, OutputSafeDNSRecordsProvider(templateRecords))
}

func safednsTemplateRecordCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <template: id/name>",
		Short:   "Creates a template record",
		Long:    "This command creates a template record",
		Example: "ans safedns template record create \"main template\" --name subdomain.ans.co.uk --type A --content 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return safednsTemplateRecordCreate(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of record")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("type", "", "Type of record")
	cmd.MarkFlagRequired("type")
	cmd.Flags().String("content", "", "Record content")
	cmd.MarkFlagRequired("content")
	cmd.Flags().Int("ttl", 0, "Record TTL")
	cmd.MarkFlagRequired("ttl")
	cmd.Flags().Int("priority", 0, "Record priority")

	return cmd
}

func safednsTemplateRecordCreate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		return err
	}

	recordName, _ := cmd.Flags().GetString("name")
	recordType, _ := cmd.Flags().GetString("type")
	recordContent, _ := cmd.Flags().GetString("content")

	createRequest := safedns.CreateRecordRequest{
		Name:    recordName,
		Type:    recordType,
		Content: recordContent,
	}

	if cmd.Flags().Changed("ttl") {
		recordTTLRaw, _ := cmd.Flags().GetInt("ttl")
		recordTTL := safedns.RecordTTL(recordTTLRaw)
		createRequest.TTL = &recordTTL
	}

	if cmd.Flags().Changed("priority") {
		recordPriority, _ := cmd.Flags().GetInt("priority")
		createRequest.Priority = ptr.Int(recordPriority)
	}

	id, err := service.CreateTemplateRecord(templateID, createRequest)
	if err != nil {
		return fmt.Errorf("Error creating record: %s", err)
	}

	templateRecord, err := service.GetTemplateRecord(templateID, id)
	if err != nil {
		return fmt.Errorf("Error retrieving new record: %s", err)
	}

	return output.CommandOutput(cmd, OutputSafeDNSRecordsProvider([]safedns.Record{templateRecord}))
}

func safednsTemplateRecordUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <template: id/name> <record: id>...",
		Short:   "Updates a template record",
		Long:    "This command updates one or more template records",
		Example: "ans safedns template record update \"main template\" 123 --content 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
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

			return safednsTemplateRecordUpdate(c.SafeDNSService(), cmd, args)
		},
	}

	cmd.Flags().String("name", "", "Name of record")
	cmd.Flags().String("type", "", "Type of record")
	cmd.Flags().String("content", "", "Record content")
	cmd.Flags().Int("ttl", 0, "Record TTL")
	cmd.Flags().Int("priority", 0, "Record priority")

	return cmd
}

func safednsTemplateRecordUpdate(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		return err
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
	if cmd.Flags().Changed("ttl") {
		recordTTLRaw, _ := cmd.Flags().GetInt("ttl")
		recordTTL := safedns.RecordTTL(recordTTLRaw)
		patchRequest.TTL = &recordTTL
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

	return output.CommandOutput(cmd, OutputSafeDNSRecordsProvider(templateRecords))
}

func safednsTemplateRecordDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <template: id/name> <record: id>...",
		Short:   "Removes a template record",
		Long:    "This command deletes one or more template records",
		Example: "ans safedns template record delete \"main template\" 123",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing template")
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

			return safednsTemplateRecordDelete(c.SafeDNSService(), cmd, args)
		},
	}
}

func safednsTemplateRecordDelete(service safedns.SafeDNSService, cmd *cobra.Command, args []string) error {
	templateID, err := getSafeDNSTemplateIDByNameOrID(service, args[0])
	if err != nil {
		return err
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

	return nil
}
