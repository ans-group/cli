package account

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
)

func AccountRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "account",
		Short: "Commands relating to Account service",
	}

	// Child root commands
	cmd.AddCommand(accountClientRootCmd(f))
	cmd.AddCommand(accountContactRootCmd(f))
	cmd.AddCommand(accountDetailsRootCmd(f))
	cmd.AddCommand(accountCreditRootCmd(f))

	return cmd
}
