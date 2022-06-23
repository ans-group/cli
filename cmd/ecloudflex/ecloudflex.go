package ecloudflex

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
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
