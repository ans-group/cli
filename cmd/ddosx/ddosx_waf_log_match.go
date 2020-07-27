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
		Example: "ukfast ddosx waf log match list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return ddosxWAFLogMatchList(c.DDoSXService(), cmd, args)
		},
	}

	cmd.Flags().String("request", "", "Show matches for specific request")

	return cmd
}

func ddosxWAFLogMatchList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) error {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		return err
	}

	var matches []ddosx.WAFLogMatch

	if cmd.Flags().Changed("request") {
		request, _ := cmd.Flags().GetString("request")
		matches, err = service.GetWAFLogRequestMatches(request, params)
		if err != nil {
			return fmt.Errorf("Error retrieving WAF log request matches: %s", err)
		}
	} else {
		matches, err = service.GetWAFLogMatches(params)
		if err != nil {
			return fmt.Errorf("Error retrieving WAF log matches: %s", err)
		}
	}

	return output.CommandOutput(cmd, OutputDDoSXWAFLogMatchesProvider(matches))
}

func ddosxWAFLogMatchShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show <request: id> <match: id>...",
		Short:   "Shows WAF log request matches",
		Long:    "This command shows a WAF log request matches",
		Example: "ukfast ddosx waf log match show 2d8556677081cecf112b555c359a78c6 123456",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing request")
			}

			if len(args) < 2 {
				return errors.New("Missing match")
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

	return output.CommandOutput(cmd, OutputDDoSXWAFLogMatchesProvider(logs))
}
