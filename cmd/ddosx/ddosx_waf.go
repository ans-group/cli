package ddosx

import (
	"github.com/spf13/cobra"
	"github.com/ukfast/cli/internal/pkg/factory"
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
