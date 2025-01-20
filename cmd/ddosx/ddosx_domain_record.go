package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxDomainRecordRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "sub-commands relating to domain records",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainRecordListCmd(f))
	cmd.AddCommand(ddosxDomainRecordShowCmd(f))
	cmd.AddCommand(ddosxDomainRecordCreateCmd(f))
	cmd.AddCommand(ddosxDomainRecordUpdateCmd(f))
	cmd.AddCommand(ddosxDomainRecordDeleteCmd(f))

	return cmd
}

func ddosxDomainRecordListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists domain records",
		Long:    "This command lists domain record",
		Example: "ans ddosx domain record list example.com",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainRecordList(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainRecordList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	records, err := service.GetDomainRecords(args[0], params)
	if err != nil {
		return fmt.Errorf("Error retrieving domain records: %s", err)
	}

	return output.CommandOutput(cmd, RecordCollection(records))
}

func ddosxDomainRecordShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <domain: name> <record: id>...",
		Short:   "Shows domain records",
		Long:    "This command shows a domain record",
		Example: "ans ddosx domain record show example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
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

			return ddosxDomainRecordShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainRecordShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {

	var records []ddosx.Record

	for _, arg := range args[1:] {
		record, err := service.GetDomainRecord(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving domain record [%s]: %s", arg, err.Error())
			continue
		}

		records = append(records, record)
	}

	return output.CommandOutput(cmd, RecordCollection(records))
}

func ddosxDomainRecordCreateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "create <domain: name>",
		Short:   "Creates a domain record",
		Long:    "This command creates a new domain record",
		Example: "ans ddosx domain record create example.com --name sub.example.com --type A",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxDomainRecordCreate(c.DDoSXService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of record")
	cmd.MarkFlagRequired("name")
	cmd.Flags().String("type", "", "Type of record")
	cmd.Flags().String("content", "", "Content of record")
	cmd.Flags().String("ssl-id", "", "ID of SSL to use for record")
	cmd.Flags().Int("safedns-record-id", 0, "ID of SafeDNS record")

	return cmd
}

func ddosxDomainRecordCreate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	name, _ := cmd.Flags().GetString("name")
	recordType, _ := cmd.Flags().GetString("type")
	content, _ := cmd.Flags().GetString("content")
	sslID, _ := cmd.Flags().GetString("ssl-id")
	safednsRecordID, _ := cmd.Flags().GetInt("safedns-record-id")

	createRequest := ddosx.CreateRecordRequest{
		Name:            name,
		Type:            ddosx.RecordType(recordType),
		Content:         content,
		SSLID:           sslID,
		SafeDNSRecordID: safednsRecordID,
	}

	id, err := service.CreateDomainRecord(args[0], createRequest)
	if err != nil {
		return fmt.Errorf("Error creating domain record: %s", err)
	}

	record, err := service.GetDomainRecord(args[0], id)
	if err != nil {
		return fmt.Errorf("Error retrieving new domain record [%s]: %s", id, err)
	}

	return output.CommandOutput(cmd, RecordCollection([]ddosx.Record{record}))
}

func ddosxDomainRecordUpdateCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "update <domain: name>...",
		Short:   "Updates a domain record",
		Long:    "This command updates one or more domain records",
		Example: "ans ddosx domain record update example.com 00000000-0000-0000-0000-000000000000 --content 1.2.3.4",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
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

			return ddosxDomainRecordUpdate(c.DDoSXService(), cmd, args)
		},
	}

	// Setup flags
	cmd.Flags().String("name", "", "Name of record")
	cmd.Flags().String("type", "", "Type of record")
	cmd.Flags().String("content", "", "Content of record")
	cmd.Flags().String("ssl-id", "", "ID of SSL to use for record")
	cmd.Flags().Int("safedns-record-id", 0, "ID of SafeDNS record")

	return cmd
}

func ddosxDomainRecordUpdate(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	patchRequest := ddosx.PatchRecordRequest{}

	if cmd.Flags().Changed("name") {
		patchRequest.Name, _ = cmd.Flags().GetString("name")
	}
	if cmd.Flags().Changed("type") {
		recordType, _ := cmd.Flags().GetString("type")
		patchRequest.Type = ddosx.RecordType(recordType)
	}
	if cmd.Flags().Changed("content") {
		patchRequest.Content, _ = cmd.Flags().GetString("content")
	}
	if cmd.Flags().Changed("ssl-id") {
		patchRequest.SSLID, _ = cmd.Flags().GetString("ssl-id")
	}
	if cmd.Flags().Changed("safedns-record-id") {
		patchRequest.SafeDNSRecordID, _ = cmd.Flags().GetInt("safedns-record-id")
	}

	var records []ddosx.Record
	for _, arg := range args[1:] {
		err := service.PatchDomainRecord(args[0], arg, patchRequest)
		if err != nil {
			output.OutputWithErrorLevelf("Error updating domain record [%s]: %s", arg, err.Error())
			continue
		}

		record, err := service.GetDomainRecord(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving updated domain record [%s]: %s", arg, err)
			continue
		}

		records = append(records, record)
	}

	return output.CommandOutput(cmd, RecordCollection(records))
}

func ddosxDomainRecordDeleteCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "delete <domain: name>...",
		Short:   "Deletes a domain record",
		Long:    "This command deletes one or more domain records",
		Example: "ans ddosx domain record delete example.com 00000000-0000-0000-0000-000000000000",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
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

			ddosxDomainRecordDelete(c.DDoSXService(), cmd, args)
			return nil
		},
	}
}

func ddosxDomainRecordDelete(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	for _, arg := range args[1:] {
		err := service.DeleteDomainRecord(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error removing domain record [%s]: %s", arg, err.Error())
			continue
		}
	}
}
