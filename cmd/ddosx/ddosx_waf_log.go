package ddosx

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxWAFLogRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "sub-commands relating to domain web application firewall logs",
	}

	// Child root commands
	cmd.AddCommand(ddosxWAFLogMatchRootCmd(f))

	// Child commands
	cmd.AddCommand(ddosxWAFLogListCmd(f))
	cmd.AddCommand(ddosxWAFLogShowCmd(f))

	return cmd
}

func ddosxWAFLogListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists WAF logs",
		Long:    "This command lists WAF logs",
		Example: "ukfast ddosx waf log list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxWAFLogList(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("domain", "", "Domain name for filtering")

	return cmd
}

func ddosxWAFLogList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	helper.HydrateAPIRequestParametersWithStringFilterFlag(&params, cmd, helper.NewStringFilterFlag("domain", "domain"))

	paginatedLogs, err := service.GetWAFLogsPaginated(params)
	if err != nil {
		return fmt.Errorf("Error retrieving WAF logs: %s", err)
	}

	return output.CommandOutputPaginated(cmd, OutputDDoSXWAFLogsProvider(paginatedLogs.Items), paginatedLogs)
}

func ddosxWAFLogShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <log: id>...",
		Short:   "Shows WAF logs",
		Long:    "This command shows a WAF log",
		Example: "ukfast ddosx waf log show 2d8556677081cecf112b555c359a78c6",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing log")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxWAFLogShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxWAFLogShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var logs []ddosx.WAFLog

	for _, arg := range args {
		log, err := service.GetWAFLog(arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving WAF log [%s]: %s", arg, err.Error())
			continue
		}

		logs = append(logs, log)
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFLogsProvider(logs))
}
