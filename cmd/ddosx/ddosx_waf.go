package ddosx

import (
	"github.com/ans-group/cli/internal/pkg/factory"
	"github.com/spf13/cobra"
)

func ddosxWAFRootCmd(f factory.ClientFactory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "waf",
		Short: "sub-commands relating to web application filewalls",
	}

	// Child root commands
	cmd.AddCommand(ddosxWAFLogRootCmd(f))

	return cmd
}
