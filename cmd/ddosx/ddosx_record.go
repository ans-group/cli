package ddosx

import (
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
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
		Example: "ans ddosx record list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxRecordList(c.DDoSXService(), cmd, args)
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

	return output.CommandOutput(cmd, RecordCollection(records))
}
