package cmd

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/ltaas"
)

func loadtestAgreementRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "agreement",
		Short: "sub-commands relating to agreements",
	}

	// Child commands
	cmd.AddCommand(loadtestAgreementShowCmd())

	return cmd
}

func loadtestAgreementShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show <agreement: type>...",
		Short:   "Shows a agreement",
		Long:    "This command shows one or more agreements",
		Example: "ukfast loadtest agreement show single",
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("Missing agreement type")
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return loadtestAgreementShow(getClient().LTaaSService(), cmd, args)
		},
	}
}

func loadtestAgreementShow(service ltaas.LTaaSService, cmd *cobra.Command, args []string) error {
	var agreements []ltaas.Agreement
	for _, arg := range args {
		parsedAgreementType, err := ltaas.ParseAgreementType(arg)
		if err != nil {
			return err
		}

		agreement, err := service.GetLatestAgreement(parsedAgreementType)
		if err != nil {
			output.OutputWithErrorLevelf("Error retrieving agreement [%s]: %s", arg, err)
			continue
		}

		agreements = append(agreements, agreement)
	}

	return outputLoadTestAgreements(agreements)
}
