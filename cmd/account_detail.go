package cmd

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func accountDetailRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "detail",
		Short: "sub-commands relating to details",
	}

	// Child commands
	cmd.AddCommand(accountDetailShowCmd())

	return cmd
}

func accountDetailShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows account details",
		Long:    "This command shows account details",
		Example: "ukfast account detail show",
		Run: func(cmd *cobra.Command, args []string) {
			accountDetailShow(getClient().AccountService(), cmd, args)
		},
	}
}

func accountDetailShow(service account.AccountService, cmd *cobra.Command, args []string) {
	details, err := service.GetDetails()
	if err != nil {
		output.Fatalf("Error retrieving details: %s", err)
		return
	}

	outputAccountDetails([]account.Details{details})
}
