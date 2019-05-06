package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func accountDetailsRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "details",
		Short: "sub-commands relating to account details",
	}

	// Child commands
	cmd.AddCommand(accountDetailsShowCmd())

	return cmd
}

func accountDetailsShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows account details",
		Long:    "This command shows account details",
		Example: "ukfast account detail show",
		Run: func(cmd *cobra.Command, args []string) {
			accountDetailsShow(getClient().AccountService(), cmd, args)
		},
	}
}

func accountDetailsShow(service account.AccountService, cmd *cobra.Command, args []string) {
	details, err := service.GetDetails()
	if err != nil {
		output.Fatalf("Error retrieving details: %s", err)
		return
	}

	outputAccountDetails([]account.Details{details})
}
