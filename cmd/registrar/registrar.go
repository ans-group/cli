package registrar

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
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
