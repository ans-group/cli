package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ddosx"
)

func ddosxDomainCDNRuleRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rule",
		Short: "sub-commands relating to domain CDN rules",
	}

	// Child commands
	cmd.AddCommand(ddosxDomainCDNRuleListCmd())

	return cmd
}

func ddosxDomainCDNRuleListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list <domain: name>",
		Short:   "Lists domain CDN rules",
		Long:    "This command lists CDN rules",
		Example: "ukfast ddosx domain cdn rule list",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing domain")
			}

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ddosxDomainCDNRuleList(getClient().DDoSXService(), cmd, args)
		},
	}
}

func ddosxDomainCDNRuleList(service ddosx.DDoSXService, cmd *cobra.Command, args []string) {
	params, err := GetAPIRequestParametersFromFlags()
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	domains, err := service.GetDomainCDNRules(args[0], params)
	if err != nil {
		output.Fatalf("Error retrieving domain CDN rules: %s", err)
		return
	}

	outputDDoSXCDNRules(domains)
}
