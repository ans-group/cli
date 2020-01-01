package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxRecordRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "sub-commands relating to records",
	}

	// Child commands
	cmd.AddCommand(ddosxRecordListCmd())

	return cmd
}

func ddosxRecordListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists records",
		Long:    "This command lists records",
		Example: "ukfast ddosx record list",
		Run: func(cmd *cobra.Command, args []string) {
			ddosxRecordList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxRecordList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	records, err := service.GetRecords(params)
	if err != nil {
		output.Fatalf("Error retrieving records: %s", err)
		return
	}

	outputDDoSXRecords(records)
}
