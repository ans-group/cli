package account

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
	"github.com/ukfast/cli/internal/pkg/output"
	"github.com/ukfast/sdk-go/pkg/service/account"
)

func accountDetailsRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "details",
		Short: "sub-commands relating to account details",
	}

	// Child commands
	cmd.AddCommand(accountDetailsShowCmd(f))

	return cmd
}

func accountDetailsShowCmd(f factory.ClientFactory) *cobra.Command {
	return &cobra.Command{
		Use:     "show",
		Short:   "Shows account details",
		Long:    "This command shows account details",
		Example: "ukfast account detail show",
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := f.NewClient()
			if err != nil {
				return err
			}

			return accountDetailsShow(c.AccountService(), cmd, args)
		},
	}
}

func accountDetailsShow(service account.AccountService, cmd *cobra.Command, args []string) error {
	details, err := service.GetDetails()
	if err != nil {
		return fmt.Errorf("Error retrieving details: %s", err)
	}

	return output.CommandOutput(cmd, OutputAccountDetailsProvider([]account.Details{details}))
}
