package ddosx

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxRecordRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "record",
		Short: "sub-commands relating to records",
	}

	// Child commands
	cmd.AddCommand(ddosxRecordListCmd(f))

	return cmd
}

func ddosxRecordListCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists records",
		Long:    "This command lists records",
		Example: "ukfast ddosx record list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return ddosxRecordList(f.NewClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxRecordList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	records, err := service.GetRecords(params)
	if err != nil {
		return fmt.Errorf("Error retrieving records: %s", err)
	}

	return output.CommandOutput(cmd, OutputDDoSXRecordsProvider(records))
}
