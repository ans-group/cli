package account

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func AccountRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Commands relating to Account service",
	}

	// Child root commands
	cmd.AddCommand(accountContactRootCmd(f))
	cmd.AddCommand(accountDetailsRootCmd(f))
	cmd.AddCommand(accountCreditRootCmd(f))
	cmd.AddCommand(accountInvoiceRootCmd(f))
	cmd.AddCommand(accountInvoiceQueryRootCmd(f))

	return cmd
}
