package ecloudflex

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
)

func ECloudFlexRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ecloudflex",
		Short: "Commands relating to eCloud Flex service",
	}

	// Child root commands
	cmd.AddCommand(ecloudflexProjectRootCmd(f))

	return cmd
}
