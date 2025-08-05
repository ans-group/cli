package ddosx

import (
	"errors"
	"fmt"

	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/ans-group/cli/internal/pkg/helper"
	"github.com/ans-group/cli/internal/pkg/output"
	"github.com/ans-group/sdk-go/pkg/connection"
	"github.com/ans-group/sdk-go/pkg/service/ddosx"
	"github.com/spf13/cobra"
)

func ddosxWAFLogMatchRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "match",
		Short: "sub-commands relating to domain web application firewall log matches",
	}

	// Child commands
	cmd.AddCommand(ddosxWAFLogMatchListCmd(f))
	cmd.AddCommand(ddosxWAFLogMatchShowCmd(f))

	return cmd
}

func ddosxWAFLogMatchListCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "Lists WAF log matches",
		Long:    "This command lists WAF log matches",
		Example: "ans ddosx waf log match list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxWAFLogMatchList(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("log", "", "Show matches for specific log")

	return cmd
}

func ddosxWAFLogMatchList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var err error

	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	var paginatedMatches *connection.Paginated[ddosx.WAFLogMatch]

	if cmd.Flags().Changed("log") {
		log, _ := cmd.Flags().GetString("log")
		paginatedMatches, err = service.GetWAFLogRequestMatchesPaginated(log, params)
	} else {
		paginatedMatches, err = service.GetWAFLogMatchesPaginated(params)
	}
	if err != nil {
		return fmt.Errorf("error retrieving WAF log matches: %s", err)
	}

	return output.CommandOutputPaginated(cmd, WAFLogMatchCollection(paginatedMatches.Items()), paginatedMatches)
}

func ddosxWAFLogMatchShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <log: id> <match: id>...",
		Short:   "Shows WAF log matches",
		Long:    "This command shows a WAF log matches",
		Example: "ans ddosx waf log match show 2d8556677081cecf112b555c359a78c6 123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("missing log")
			}

			if len(args) < 2 {
				return errors.New("missing match")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxWAFLogMatchShow(c.DDoSXService(), cmd, args)
		},
	}
}

func ddosxWAFLogMatchShow(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	var logs []ddosx.WAFLogMatch

	for _, arg := range args[1:] {
		log, err := service.GetWAFLogRequestMatch(args[0], arg)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving WAF log matches [%s]: %s", arg, err.Error())
			continue
		}

		logs = append(logs, log)
	}

	return output.CommandOutput(cmd, WAFLogMatchCollection(logs))
}
