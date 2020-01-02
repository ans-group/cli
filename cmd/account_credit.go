package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/helper"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func accountCreditRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "credit",
		Short: "sub-commands relating to credits",
	}

	// Child commands
	cmd.AddCommand(accountCreditListCmd())

	return cmd
}

func accountCreditListCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "list",
		Short:   "Lists credits",
		Long:    "This command lists credits",
		Example: "ukfast account credit list",
		Run: func(cmd *cobra.Command, args []string) {
			accountCreditList(getClient().AccountService(), cmd, args)
		},
	}
}

func accountCreditList(service account.AccountService, cmd *cobra.Command, args []string) {
	params, err := helper.GetAPIRequestParametersFromFlags(cmd)
	if err != nil {
		output.Fatal(err.Error())
		return
	}

	credits, err := service.GetCredits(params)
	if err != nil {
		output.Fatalf("Error retrieving credits: %s", err)
		return
	}

	outputAccountCredits(credits)
}
