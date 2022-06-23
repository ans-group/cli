package registrar

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
)

func RegistrarRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "registrar",
		Short: "Commands relating to Registrar service",
	}

	// Child root commands
	cmd.AddCommand(registrarDomainRootCmd(f))
	cmd.AddCommand(registrarWhoisRootCmd(f))

	return cmd
}
